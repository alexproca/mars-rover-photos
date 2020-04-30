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

//MainConfig - main config used across application (read config from .env file or from env.
//You can specify settings standalone as env variables (ex: NASA_API_KEY) or you can create a json
//respecting MainConfig structure above and store it in JSON_CONFIG env var
var Config MainConfig = MainConfig{}

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

func LoadConfig(filenames... string) {

	if err := json.Unmarshal(jsonDefault(), &Config); err != nil {
		panic(errors.New(err.Error()))
	}

	godotenv.Load(filenames...)
	if _, err := env.UnmarshalFromEnviron(&Config); err != nil {
		panic(errors.New(err.Error()))
	}

}
