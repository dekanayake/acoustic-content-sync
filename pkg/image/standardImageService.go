//go:build standard
// +build standard

package image

import (
	"os"
)

type imageService struct {
}

func (i imageService) IsImageInExpectedDimension(width uint, height uint, asset *os.File) (bool, error) {
	//TODO implement me
	panic("Not supported use content sync with image magick support")
}

func (i imageService) Resize(width uint, height uint, asset *os.File) (*os.File, error) {
	//TODO implement me
	panic("Not supported use content sync with image magick support")
}

func initImageService() *imageService {
	return &imageService{}
}
