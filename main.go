package main

import (
	"fmt"
	"log"
	"nasa-api/config"
	"nasa-api/entities"
	"nasa-api/util"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
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
	http.HandleFunc("/ws/", wsHandler)
	http.HandleFunc("/ws", wsHandler)

	log.Printf("Starting server on interface '%s' and port '%s'\n", httpConfig.Interface, httpConfig.Port)
	http.ListenAndServe(fmt.Sprintf("%s:%s", httpConfig.Interface, httpConfig.Port), nil)

}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(writer http.ResponseWriter, request *http.Request) {
	ws, upgradeError := upgrader.Upgrade(writer, request, nil)

	if upgradeError != nil {
		log.Println(upgradeError)
		return
	}

	go handleConnection(ws)
}

func handleConnection(ws *websocket.Conn) {
	defer ws.Close()

	messageType, message, readErr := ws.ReadMessage()
	if readErr != nil {
		log.Println("Read Error:", readErr)
		return
	}

	if messageType == websocket.TextMessage {
		stringMessage := string(message)
		data := strings.Split(stringMessage, "/")

		if len(data) != 6 {
			log.Println("Bad request:", stringMessage)
			return
		}

		rover, camera, currentDate, maxDate := data[2], data[3], data[4], data[5]
		it, err := entities.NewPhotosIterator(rover, camera, currentDate, maxDate)

		if err != nil {
			log.Println("Cannot create iterator:", stringMessage)
			return
		}

		ticker := time.NewTicker(time.Second * 5)

		for it.HasNext() {

			photo := it.Next()
			response := photo.ImageURL

			err := ws.WriteMessage(websocket.TextMessage, []byte(response))
			if err != nil {
				log.Println("Write error:", err)
				break
			}

			log.Println("Sending: ", response)
			<-ticker.C
		}
	}
}

func cameraHandler(writer http.ResponseWriter, request *http.Request) {
	roverNameAndCamera := strings.TrimPrefix(request.URL.Path, "/camera/")

	tokens := strings.Split(roverNameAndCamera, "/")

	if len(tokens) != 4 {
		log.Println("Incorrect camera path: ", roverNameAndCamera)
		http.Redirect(writer, request, "/static/error.html", 302)
		return
	}

	roverName, cameraName, earthDate, maxDate := tokens[0], tokens[1], tokens[2], tokens[3]

	items := struct {
		RoverName  string
		CameraName string
		EarthDate  string
		MaxDate    string
	}{
		RoverName:  roverName,
		CameraName: cameraName,
		EarthDate:  earthDate,
		MaxDate:    maxDate,
	}
	if templateError := util.TemplateHandler("photos", items, writer); templateError != nil {
		log.Println(templateError)
		http.Redirect(writer, request, "/static/error.html", 302)
	}
}

func roverHandler(writer http.ResponseWriter, request *http.Request) {

	roverName := strings.TrimPrefix(request.URL.Path, "/rover/")
	if rover, retrieveRoverError := entities.GetRover(roverName); retrieveRoverError == nil {

		selectedDateValues := request.URL.Query()["selected_date"]

		if len(selectedDateValues) != 1 {
			log.Println("Something is not right with selected date", selectedDateValues)
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
			log.Println("Something is wrong with template: rover", templateError)
			http.Redirect(writer, request, "/static/error.html", 302)
		}
	} else {
		log.Printf("Something went wrong retrieving rover '%s' informatioin from NASA API: %s\n", roverName, retrieveRoverError)
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
			log.Println("Something is wrong with template: index", err)
			http.Redirect(writer, request, "/static/error.html", 302)
		}
	} else {
		log.Println("Something went wrong retrieving all rovers from NASA API: ", getRoversError)
		http.Redirect(writer, request, "/static/error.html", 302)
	}
}
