package entities

import (
	"encoding/json"
	"fmt"
	"log"
	"nasa-api/config"
	"nasa-api/util"
)

//Camera - camera of the rover
type Camera struct {
	FullName string `json:"full_name"`
	Name     string `json:"name"`
}

//Rover - rover entity
type Rover struct {
	Cameras     []Camera `json:"cameras"`
	ID          int      `json:"id"`
	LandingDate string   `json:"landing_date"`
	LaunchDate  string   `json:"launch_date"`
	MaxDate     string   `json:"max_date"`
	MaxSol      int      `json:"max_sol"`
	Name        string   `json:"name"`
	Status      string   `json:"status"`
	TotalPhotos int      `json:"total_photos"`
}

//Photo - photo entity
type Photo struct {
	ID       int    `json:"id"`
	Sol      int    `json:"sol"`
	EarthDay string `json:"earth_date"`
	ImageURL string `json:"img_src"`
	Camera   Camera `json:"camera"`
	Rover    Rover  `json:"rover"`
}

//Content - content in multiple rovers endpoint
type AllRoversContent struct {
	Rovers []Rover `json:"rovers"`
}

//Content - content in single rover endpoint
type SingleRoversContent struct {
	Rover Rover `json:"rover"`
}

//Content - content in single rover endpoint
type PhotosFromSingleRoverCamera struct {
	Photos []Photo `json:"photos"`
}

//GetAllRovers - get a list of rovers
func GetAllRovers() ([]Rover, error) {

	endPointData := config.Config.EndpointData

	url := fmt.Sprintf("%s?api_key=%s", endPointData.RoversEndPoint, endPointData.ApiKey)

	body := util.HTTPGet(url)

	content := AllRoversContent{}

	if unmarshallError := json.Unmarshal(body, &content); unmarshallError != nil {
		log.Println("Something went wrong unpacking all rover data from NASA API", body, unmarshallError)
		return nil, unmarshallError
	}

	return content.Rovers, nil
}

//GetRover - get info about a specific rover
func GetRover(roverName string) (*Rover, error) {

	endPointData := config.Config.EndpointData

	url := fmt.Sprintf("%s/%s?api_key=%s", endPointData.RoversEndPoint, roverName, endPointData.ApiKey)

	body := util.HTTPGet(url)

	content := SingleRoversContent{}

	if unmarshallError := json.Unmarshal(body, &content); unmarshallError != nil {
		log.Printf("Something went wrong unpacking rover '%s' data from NASA API \n %s \n %s \n", roverName, body, unmarshallError)
		return nil, unmarshallError
	}

	return &content.Rover, nil
}

//GetPhotos - get photos for a specific camera on a specific rover
func GetPhotos(roverName, cameraName, earthDate string) ([]Photo, error) {

	endPointData := config.Config.EndpointData
	pageNumber := 1

	url := fmt.Sprintf("%s/%s/photos?page=%d&earth_date=%s&camera=%s&api_key=%s",
		endPointData.RoversEndPoint,
		roverName,
		pageNumber,
		earthDate,
		cameraName,
		endPointData.ApiKey,
	)

	body := util.HTTPGet(url)

	content := PhotosFromSingleRoverCamera{}

	if unmarshallError := json.Unmarshal(body, &content); unmarshallError != nil {
		log.Printf("Something went wrong unpacking rover '%s' and camera '%s' data from NASA API\n %s \n %s \n", roverName, cameraName, body, unmarshallError)
		return nil, unmarshallError
	}

	return content.Photos, nil
}
