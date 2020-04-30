package config

import (
	"encoding/json"
	"errors"
	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
	"os"
)

type MainConfig struct {

	EndpointData struct {
		ApiKey string `env:"NASA_API_KEY" json:"api_key"`
		RoversEndPoint string `env:"ROVERS_ENDPOINT" json:"rover_endpoint"`
	} `json:"endpoint_data"`
	ServerData struct {
		Interface string `env:"INTERFACE" json:"interface"`
		Port string `env:"PORT" json:"port""`
	} `json:"server_data"`
}

func jsonDefault() []byte {
	if jsonConfig := os.Getenv("JSON_CONFIG"); jsonConfig == "" {
		return []byte(`{
			"server_data": {
				"interface": "0.0.0.0",
				"port": "8081"
			},
			"endpoint_data": {
				"rover_endpoint": "https://api.nasa.gov/mars-photos/api/v1/rovers"
			  }
		}`)
	} else {
		return []byte(jsonConfig)
	}
}
var Config MainConfig = MainConfig{}

func LoadConfig(filenames... string) {

	if err := json.Unmarshal(jsonDefault(), &Config); err != nil {
		panic(errors.New(err.Error()))
	}

	godotenv.Load(filenames...)
	if _, err := env.UnmarshalFromEnviron(&Config); err != nil {
		panic(errors.New(err.Error()))
	}

}
