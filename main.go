package main

import (
	"errors"
	"fmt"
	"nasa-api/config"
	"nasa-api/entities"
	"nasa-api/util"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {

	config.LoadConfig()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGUSR1)

	util.SystemSignalsHandler(signals)

	httpConfig := config.Config.ServerData

	publicFs := http.FileServer(http.Dir("public"))

	http.Handle("/img/", publicFs)
	http.Handle("/js/", publicFs)
	http.Handle("/css/", publicFs)
	http.Handle("/static/", publicFs)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/rover/", roverHandler)
	http.HandleFunc("/camera/", cameraHandler)

	fmt.Printf("Starting server on interface '%s' and port '%s'\n", httpConfig.Interface, httpConfig.Port)
	http.ListenAndServe(fmt.Sprintf("%s:%s", httpConfig.Interface, httpConfig.Port), nil)

}

func cameraHandler(writer http.ResponseWriter, request *http.Request) {
	roverNameAndCamera := strings.TrimPrefix(request.URL.Path, "/camera/")

	tokens := strings.Split(roverNameAndCamera, "/")

	if len(tokens) != 3 {
		fmt.Println(errors.New(fmt.Sprintf("Incorrect path '%s'", roverNameAndCamera)))
		http.Redirect(writer, request, "/static/error.html", 302)
		return
	}

	roverName, cameraName, earthDate := tokens[0], tokens[1], tokens[2]

	if photos, photosError := entities.GetPhotos(roverName, cameraName, earthDate); photosError == nil {
		items := struct {
			Photos []entities.Photo
			RoverName string
			CameraName string
			EarthDate string
		}{
			Photos: photos,
			RoverName: roverName,
			CameraName: cameraName,
			EarthDate: earthDate,
		}
		if templateError := util.TemplateHandler("photos", items, writer); templateError != nil {
			fmt.Println(templateError)
			http.Redirect(writer, request, "/static/error.html", 302)
		}
	} else {
		fmt.Println(photosError)
		http.Redirect(writer, request, "/static/error.html", 302)
	}
}

func roverHandler(writer http.ResponseWriter, request *http.Request) {

	roverName := strings.TrimPrefix(request.URL.Path, "/rover/")
	if rover, retrieveRoverError := entities.GetRover(roverName); retrieveRoverError == nil {

		selectedDateValues := request.URL.Query()["selected_date"]

		if len(selectedDateValues) != 1 {
			fmt.Println(errors.New("Something is not right with selected date"))
			http.Redirect(writer, request, "/static/error.html", 302)
		}

		selectedDate := selectedDateValues[0]

		if selectedDate == "" {
			selectedDate = rover.LandingDate
		}

		items := struct {
			Rover       entities.Rover
			CurrentDate string
		}{
			Rover:       *rover,
			CurrentDate: selectedDate,
		}
		if templateError := util.TemplateHandler("rover", items, writer); templateError != nil {
			fmt.Println(templateError)
			http.Redirect(writer, request, "/static/error.html", 302)
		}
	} else {
		fmt.Println(retrieveRoverError)
		http.Redirect(writer, request, "/static/error.html", 302)
	}
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {

	if rovers, getRoversError := entities.GetAllRovers(); getRoversError == nil {
		items := struct {
			Rovers []entities.Rover
		}{
			Rovers: rovers,
		}

		if err := util.TemplateHandler("index", items, writer); err != nil {
			fmt.Println(err)
			http.Redirect(writer, request, "/static/error.html", 302)
		}
	} else {
		fmt.Println(getRoversError)
		http.Redirect(writer, request, "/static/error.html", 302)
	}
}

