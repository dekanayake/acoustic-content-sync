package image

import (
	"os"
	"sync"
)

var imageServiceOnce sync.Once
var imageServiceInstance *imageService

type ImageService interface {
	IsImageInExpectedDimension(width uint, height uint, asset *os.File) (bool, error)
	Resize(width uint, height uint, asset *os.File) (*os.File, error)
}

func GetImageService() ImageService {
	imageServiceOnce.Do(func() {
		imageServiceInstance = initImageService()
	})
	return imageServiceInstance
}
