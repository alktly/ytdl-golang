package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type metaData struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Formats    []format `json:"formats"`
	Duration   int      `json:"duration"`
	UploadDate string   `json:"upload_date"`
}

type format struct {
	Format    string `json:"format"`
	FileSize  int    `json:"filesize"`
	Extension string `json:"ext"`
}

func welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("YTDownloader API"))
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	ytid := r.FormValue("yt_id")
	format := r.FormValue("fmt")

	if len(ytid) == 0 {
		badRequest(errors.New("Youtube ID cannot be empty"), w)
		return
	}

	fmt.Printf("Downloading video...\n")
	if status, err := downloadVideo(ytid, format); err != nil {
		writeErrorStatus(status, err, w)
		return
	}
	fmt.Printf("Done\n")

	fmt.Printf("Hashing and renaming media file...\n")
	status, hash, err := hashAndRenameMediaFile(ytid)
	if err != nil {
		writeErrorStatus(status, err, w)
		return
	}
	fmt.Printf("Done\n")

	success(hash, w)
}

func playlistHandler(w http.ResponseWriter, r *http.Request) {
	ytid := r.FormValue("yt_id")

	if len(ytid) == 0 {
		badRequest(errors.New("Youtube ID cannot be empty"), w)
		return
	}

	fmt.Printf("Downloading playlist...\n")
	if status, err := downloadPlaylist(ytid); err != nil {
		writeErrorStatus(status, err, w)
		return
	}
	fmt.Printf("Done\n")

	success("Done", w)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {

	ytid := r.FormValue("yt_id")

	if len(ytid) == 0 {
		badRequest(errors.New("Youtube ID cannot be empty"), w)
		return
	}

	fmt.Printf("Extracting video meta data...\n")
	status, res, err := extractInfo(ytid)
	if err != nil {
		writeErrorStatus(status, err, w)
	}
	fmt.Printf("Done\n")

	jsonResponse(res, w)
}

func downloadVideo(ytid string, format string) (int, error) {
	yturl := fmt.Sprintf(ytBaseURL, ytid)
	localpath := localDir + ytid

	opt := []string{
		"--no-progress",
		fmt.Sprintf("-o%v", localpath),
		yturl,
	}

	if len(format) > 0 {
		opt = append(opt, fmt.Sprintf("-f %v", format))
	}

	out, err := execCmdRead(ytdlPath, opt...)
	if err != nil {
		if catch404(out) {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, err
	}

	fmt.Print(string(out))

	return http.StatusOK, nil
}

func downloadPlaylist(ytid string) (int, error) {
	yturl := fmt.Sprintf(ytPlBaseURL, ytid)
	localpath := localDir + "%(id)s"

	opt := []string{
		"--no-progress",
		"--yes-playlist",
		"--ignore-errors",
		fmt.Sprintf("-o%v", localpath),
		yturl,
	}

	out, err := execCmdRead(ytdlPath, opt...)
	fmt.Printf("%v\n", string(out))
	if err != nil {
		if catch404(out) {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func extractInfo(ytid string) (int, []byte, error) {
	yturl := fmt.Sprintf(ytBaseURL, ytid)

	opt := []string{
		"--dump-json",
		"--skip-download",
		yturl,
	}

	out, err := execCmdRead(ytdlPath, opt...)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	//fmt.Printf(string(out))

	var meta metaData

	if err := json.Unmarshal(out, &meta); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	meta = filterUnusedExtensions(meta, "mp4")

	jsonRes, err := json.Marshal(meta)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, jsonRes, nil
}
