package entities

import (
	"encoding/json"
	"fmt"
	"nasa-api/config"
	"nasa-api/util"
)

//Camera - camera of the rover
type Camera struct {
	FullName string `json:"full_name"`
	Name     string `json:"name"`
}

//Rover - rover
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

type Photo struct {
	Sol int `json:"sol"`
	EarthDay string `json:"earth_date"`
	ImageURL string `json:"img_src"`
	Camera Camera `json:"camera"`
	Rover Rover `json:"rover"`
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

func GetAllRovers() ([]Rover, error) {

	endPointData := config.Config.EndpointData

	url := fmt.Sprintf("%s?api_key=%s", endPointData.RoversEndPoint, endPointData.ApiKey)

	body := util.GetBody(url)

	content := AllRoversContent{}

	if unmarshallError := json.Unmarshal(body, &content); unmarshallError != nil {
		fmt.Println(unmarshallError)
		return nil, unmarshallError
	}

	return content.Rovers, nil
}

func GetRover(roverName string) (*Rover, error) {

	endPointData := config.Config.EndpointData

	url := fmt.Sprintf("%s/%s?api_key=%s", endPointData.RoversEndPoint, roverName, endPointData.ApiKey)

	body := util.GetBody(url)

	content := SingleRoversContent{}

	if unmarshallError := json.Unmarshal(body, &content); unmarshallError != nil {
		fmt.Println(unmarshallError)
		return nil, unmarshallError
	}

	return &content.Rover, nil

}

func GetPhotos(roverName, cameraName, earthDate string) ([]Photo, error) {

	endPointData := config.Config.EndpointData
	pageNumber := 1

	//https://api.nasa.gov/mars-photos/api/v1/rovers/Spirit/photos?page=1&earth_date=2004-01-10&camera=fhaz&api_key=SijUgWPLeaBPwiMNl8w7Ce1jKnud42GRhKW6O7Ro
	url := fmt.Sprintf("%s/%s/photos?page=%d&earth_date=%s&camera=%s&api_key=%s",
		endPointData.RoversEndPoint,
		roverName,
		pageNumber,
		earthDate,
		cameraName,
		endPointData.ApiKey,
		)

	body := util.GetBody(url)

	content := PhotosFromSingleRoverCamera{}

	if unmarshallError := json.Unmarshal(body, &content); unmarshallError != nil {
		fmt.Println(unmarshallError)
		return nil, unmarshallError
	}

	return content.Photos, nil
}