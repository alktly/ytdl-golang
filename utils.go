package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func execCmd(process string, args ...string) error {
	cmd := exec.Command(process, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("cmd.Run() failed with %s\n", err)
	}

	return err
}

func execCmdRead(process string, args ...string) ([]byte, error) {
	cmd := exec.Command(process, args...)
	out, err := cmd.CombinedOutput()

	//fmt.Printf(string(out))

	return out, err
}

func start(process string, args ...string) (p *os.Process, err error) {
	if args[0], err = exec.LookPath(args[0]); err == nil {
		var procAttr os.ProcAttr
		procAttr.Files = []*os.File{os.Stdin,
			os.Stdout, os.Stderr}
		p, err := os.StartProcess(process, args, &procAttr)
		if err == nil {
			return p, nil
		}
	}
	return nil, err
}

func hashAndRenameMediaFile(ytid string) (int, string, error) {
	localpath := localDir + ytid

	data, err := ioutil.ReadFile(localpath)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	hash := fmt.Sprintf("%x", sha1.Sum(data))

	if err := os.Rename(localpath, localDir+hash); err != nil {
		return http.StatusInternalServerError, "", err
	}

	return http.StatusOK, hash, nil
}

func filterUnusedExtensions(meta metaData, formatToKeep string) metaData {
	var tempMeta metaData
	tempMeta = meta
	tempMeta.Formats = []format{}

	for _, element := range meta.Formats {
		if element.Extension == formatToKeep {
			tempMeta.Formats = append(tempMeta.Formats, element)
		}
	}

	return tempMeta
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func emptyDir(directory string) error {
	// Open the directory and read all its files.
	dirRead, err := os.Open(directory)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return err
	}

	dirFiles, err := dirRead.Readdir(0)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return err
	}

	fmt.Printf("There are %v files to delete.\n", len(dirFiles))

	// Loop over the directory's files.
	for index := range dirFiles {
		fileHere := dirFiles[index]

		// Get name of file and its full path.
		nameHere := fileHere.Name()
		fullPath := directory + nameHere

		// Remove the file.
		os.Remove(fullPath)
		//fmt.Printf("Removed file:", fullPath)
	}

	return nil
}
