package image

import (
	"github.com/pkg/errors"
	"gopkg.in/gographics/imagick.v3/imagick"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var imageServiceOnce sync.Once
var imageServiceInstance *imageService

type ImageService interface {
	IsImageInExpectedDimension(width uint, height uint, asset *os.File) (bool, error)
	Resize(width uint, height uint, asset *os.File) (*os.File, error)
}

type imageService struct {
}

func GetImageService() ImageService {
	imageServiceOnce.Do(func() {
		imagick.Initialize()
		imageServiceInstance = &imageService{}
		runtime.SetFinalizer(imageServiceInstance, func(imageService *ImageService) { imagick.Terminate() })
	})
	return imageServiceInstance
}

func tmpAssetFile(asset *os.File) (*os.File, error) {
	tmpImageFile, err := ioutil.TempFile("/tmp", "tempImage_*"+filepath.Ext(asset.Name()))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	input, err := ioutil.ReadFile(asset.Name())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = ioutil.WriteFile(tmpImageFile.Name(), input, 0644)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tmpImageFile, nil
}

func (i imageService) IsImageInExpectedDimension(width uint, height uint, asset *os.File) (bool, error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	tmpImageAsset, err := tmpAssetFile(asset)
	if err != nil {
		return false, errors.WithStack(err)
	}
	defer tmpImageAsset.Close()
	defer os.Remove(tmpImageAsset.Name())

	err = mw.ReadImageFile(tmpImageAsset)
	if err != nil {
		return false, errors.WithStack(err)
	}
	return mw.GetImageWidth() >= width && mw.GetImageHeight() >= height, nil
}

func (i imageService) Resize(width uint, height uint, asset *os.File) (*os.File, error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	tmpImageAsset, err := tmpAssetFile(asset)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer tmpImageAsset.Close()
	defer os.Remove(tmpImageAsset.Name())

	err = mw.ReadImageFile(tmpImageAsset)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = mw.ResizeImage(width, height, imagick.FILTER_LANCZOS)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resizedImageFile, err := ioutil.TempFile("/tmp", "resized_*"+filepath.Ext(asset.Name()))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = mw.WriteImageFile(resizedImageFile)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resizedImageFile, err = os.Open(resizedImageFile.Name())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return resizedImageFile, nil
}
