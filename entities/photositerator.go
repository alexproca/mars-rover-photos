package entities

import (
	"log"
	"time"
)

type PhotosIterator struct {
	RoverName   string
	CameraName  string
	CurrentDate time.Time
	MaxDate     time.Time
	NextPhotos  []Photo
}

func NewPhotosIterator(roverName, cameraName, startDate, maxDate string) (*PhotosIterator, error) {

	result := &PhotosIterator{
		RoverName: roverName,
		CameraName: cameraName,
	}

	if currentDate, err := time.Parse("2006-01-02", startDate); err == nil {
		result.CurrentDate = currentDate
	} else {
		return nil, err
	}

	if maxDate, err := time.Parse("2006-01-02", maxDate); err == nil {
		result.MaxDate = maxDate
	} else {
		return nil, err
	}

	return result, nil
}

func (self *PhotosIterator) Next() (*Photo) {
	arr := self.NextPhotos
	head, arr := arr[0], arr[1:]
	self.NextPhotos = arr

	return &head
}

func (self *PhotosIterator) hasNext() bool {

	for len(self.NextPhotos) == 0 && self.CurrentDate.Unix() <= self.MaxDate.Unix() {

		if photos, err := GetPhotos(self.RoverName, self.CameraName, self.getDateString()); err == nil {
			self.CurrentDate = self.CurrentDate.AddDate(0,0,1)
			if len(photos) != 0 {
				self.NextPhotos = append(self.NextPhotos, photos...)
				break
			}
		} else {
			log.Println(err)
			return false
		}
	}

	return len(self.NextPhotos) != 0
}

func (self *PhotosIterator) getDateString() string {

	if self == nil {
		return ""
	}

	return self.CurrentDate.Format("2006-01-02")
}