package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

const (
	host        = "http://localhost"
	port        = ":4000"
	ytBaseURL   = "https://www.youtube.com/watch?v=%v"
	ytPlBaseURL = "https://www.youtube.com/playlist?list=%v"
)

var (
	useAuth       = flag.Bool("USEAUTH", false, "(Optional) Pass this argument if you want to use Youtube authentication API.")
	clientSecret  = flag.String("CLSECRET", "", "(Required if USEAUTH passed) Youtube authentication API client_secrets json file path.")
	authTokenFile = flag.String("OAUTH2", "", "(Required if USEAUTH passed) Youtube authentication API oauth2 json file path.")
	argYtdlPath   = flag.String("YTDL", "", "(Optional) Youtube dl binary path. Default './youtube-dl'.")
	argOuthPath   = flag.String("OUT", "", "(Optional) Output directory path. Default './media-files'.")
	ytdlPath      = "./youtube-dl"
	localDir      = "./media-files/"
)

func main() {

	flag.Parse()

	initArgs(*argYtdlPath, *argOuthPath)

	if *useAuth {
		// Authenticate on youtube using youtube upload api
		// TODO: Connect this to sentry and create a retry mechanism later
		fmt.Printf("Authenticating yt client...\n")
		getAuthService(*clientSecret, *authTokenFile)
		fmt.Printf("Authentication successful\n")
	}

	mux := http.NewServeMux()

	for _, route := range getRoutes() {
		mux.HandleFunc(
			route.path,
			methodHandler(
				route.method,
				route.handler,
			),
		)
	}

	fmt.Printf("Listening on %v%v\n", host, port)

	if err := http.ListenAndServe(port, mux); err != nil {
		panic(err)
	}
}

func getRoutes() []route {

	return []route{
		route{
			method:  http.MethodGet,
			path:    "/",
			handler: welcome,
		},
		route{
			method:  http.MethodGet,
			path:    "/info",
			handler: infoHandler,
		},
		route{
			method:  http.MethodGet,
			path:    "/download",
			handler: downloadHandler,
		},
		route{
			method:  http.MethodGet,
			path:    "/playlist",
			handler: playlistHandler,
		},
	}
}

func initArgs(argYtdlPath string, argOuthPath string) {
	if len(argYtdlPath) != 0 {
		ytdlPath = argYtdlPath
	}

	if len(argOuthPath) != 0 {
		localDir = argOuthPath
	}
}

func getAuthService(clientSecret string, authFile string) {

	if len(clientSecret) == 0 {
		log.Fatalln("You have to provide CLSECRET json while using USEAUTH flag!")
	}

	if len(authFile) == 0 {
		log.Fatalln("You have to provide OAUTH2 json while using USEAUTH flag!")
	}

	ctx := context.Background()

	b, err := ioutil.ReadFile(clientSecret)
	if err != nil {
		fmt.Printf("Unable to read client secret file: %v\n", err)
		// TODO: Sentry
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeUploadScope)
	if err != nil {
		fmt.Printf("Unable to parse client secret file to config: %v\n", err)
		// TODO: Sentry
	}
	client := getClient(ctx, config)
	_, err = youtube.New(client)

	handleError(err, "Error creating authenticated Youtube client")
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		fmt.Printf("%v: %v\n", message, err.Error())
		// TODO: Sentry
	}
}
