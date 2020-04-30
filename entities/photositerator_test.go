package entities

import (
	"nasa-api/config"
	"os"
	"sync"
	"testing"
)

var once sync.Once

func TestNewPhotosIteratorStart(t *testing.T) {

	loadTestsConfig(t)

	it, err := NewPhotosIterator("Curiosity", "FHAZ", "2012-08-06", "2020-04-28")

	if err != nil {
		t.Error("Could not create new Photos Iterator")
	}

	if it == nil {
		t.Error("New Iterator is nil")
	}

	next := it.HasNext()

	if next == false {
		t.Error("We should have next photos")
	}

	image := it.Next()

	if image == nil {
		t.Error("Image is nil")
	}

	firstUrl := image.ImageURL

	i := 0
	for it.HasNext() && i< 100 {
		i++
		image = it.Next()
		t.Logf("Date: %s i: %d link: %s", it.CurrentDate, i, image.ImageURL)
	}

	if i != 100 {
		t.Error("Could not read 100 images from endpoint")
	}

	if image == nil {
		t.Error("Image is nil after 20 iterations")
	}

	secondUrl := image.ImageURL

	if firstUrl == secondUrl {
		t.Error("Iterator does not advance")
	}

}

func TestNewPhotosIteratorEnd(t *testing.T) {

	loadTestsConfig(t)

	it, err := NewPhotosIterator("Spirit", "NAVCAM", "2010-02-01", "2010-03-21")

	if err != nil {
		t.Error("Could not create new Photos Iterator")
	}

	if it == nil {
		t.Error("New Iterator is nil")
	}

	next := it.HasNext()

	image := it.Next()

	if image == nil {
		t.Error("Image is nil")
	}

	if next == false {
		t.Error("We should have next photos")
	}

	i := 0
	for it.HasNext() && i < 100 {
		i++
		image = it.Next()
		t.Logf("Date: %s i: %d link: %s", it.CurrentDate, i, image.ImageURL)
	}

	if i != 43 {
		t.Error("We should have exactly 43 photos starting with 2010-02-01")
	}

}

func loadTestsConfig(t *testing.T) {
	once.Do(func() {
		config.LoadConfig("../.env")

		if config.Config.EndpointData.ApiKey == "" {
			t.Log("Please add a NASA_API_KEY in .env file")
			os.Exit(1)
		}
	})
}