package api

import (
	"bufio"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/dekanayake/acoustic-content-sync/pkg/image"
	"github.com/wesovilabs/koazee"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type ContextKey string

type AcousticDataRecord struct {
	CSVRecordKey string
	NameFields   []string
	Values       []GenericData
	Tags         []string
}

type GenericData struct {
	Name    string
	Type    string
	Value   string
	Context Context
}

type Context struct {
	Data map[ContextKey]interface{}
}

const (
	AssetName             ContextKey = "AssetName"
	Profiles              ContextKey = "Profiles"
	AcousticAssetBasePath ContextKey = "AcousticAssetBasePath"
	AssetLocation         ContextKey = "AssetLocation"
	TagList               ContextKey = "Tags"
	IsWebUrl              ContextKey = "IsWebUrl"
	EnforceImageDimension ContextKey = "EnforceImageDimension"
	ImageHeight           ContextKey = "ImageHeight"
	ImageWidth            ContextKey = "ImageWidth"
)

func (context Context) getValue(key ContextKey) (interface{}, error) {
	if val, ok := context.Data[key]; ok {
		return val, nil
	} else {
		return nil, errors.ErrorMessageWithStack("no context value found for key :" + string(key))
	}
}

func (context Context) getUintValue(key ContextKey) (uint, error) {
	if val, ok := context.Data[key]; ok {
		if castVal, ok := val.(uint); ok {
			return castVal, nil
		} else {
			return 0, errors.ErrorMessageWithStack("Cannot cast the provided value" + reflect.TypeOf(val).String())
		}
	} else {
		return 0, errors.ErrorMessageWithStack("no context value found for key :" + string(key))
	}
}

func (acousticDataRecord AcousticDataRecord) Name() string {
	values := make([]GenericData, 0, len(acousticDataRecord.NameFields))
	for _, nameField := range acousticDataRecord.NameFields {
		value := koazee.StreamOf(acousticDataRecord.Values).
			Filter(func(data GenericData) bool {
				return data.Name == nameField
			}).First().Val().(GenericData)
		values = append(values, value)
	}

	return koazee.StreamOf(values).
		Reduce(func(acc string, data GenericData) string {
			if acc == "" {
				acc += data.Value
			} else {
				acc += "__" + data.Value
			}
			return acc
		}).String()
}

func (acousticDataRecord AcousticDataRecord) CSVRecordKeyValue() string {
	return koazee.StreamOf(acousticDataRecord.Values).
		Filter(func(columnValue GenericData) bool {
			return columnValue.Name == acousticDataRecord.CSVRecordKey
		}).Map(func(columnValue GenericData) string {
		return columnValue.Value
	}).First().String()
}

func (element TextElement) Convert(data interface{}) (Element, error) {
	element.Value = data.(GenericData).Value
	return element, nil
}

func (element NumberElement) Convert(data interface{}) (Element, error) {
	numValue, err := strconv.ParseInt(data.(GenericData).Value, 0, 64)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	element.Value = numValue
	return element, nil
}

func (element LinkElement) Convert(data interface{}) (Element, error) {
	element.LinkURL = data.(GenericData).Value
	return element, nil
}

func categoryIds(category string) ([]string, error) {
	catItems := strings.Split(category, "/")
	if len(catItems) == 1 {
		return nil, errors.ErrorMessageWithStack("empty category :" + catItems[0])
	}

	categoryItems, err := NewCategoryClient(env.AcousticAPIUrl()).Categories(catItems[0])
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	categoryItems = koazee.StreamOf(categoryItems).
		Filter(func(categoryItem CategoryItem) bool {
			return len(categoryItem.NamePath) > 0
		}).Out().Val().([]CategoryItem)
	catNamePaths := make([]string, 0, 10)
	for i := 1; i <= len(catItems); i++ {
		catNamePathsSlice := catItems[0:i]
		catNamePath := koazee.StreamOf(catNamePathsSlice).
			Reduce(func(acc string, catNamePath string) string {
				if acc == "" {
					acc += catNamePath
				} else {
					acc += "/" + catNamePath
				}
				return acc
			}).String()
		catNamePaths = append(catNamePaths, catNamePath)
	}

	catIds := koazee.StreamOf(categoryItems).
		Filter(func(categoryItem CategoryItem) bool {
			fullNamePath := categoryItem.FullNamePath()
			contains, _ := koazee.StreamOf(catNamePaths).Contains(fullNamePath)
			return strings.Contains(fullNamePath, "/") && contains
		}).
		Map(func(categoryItem CategoryItem) string {
			return categoryItem.Id
		}).Out().Val().([]string)

	return catIds, nil
}

func (element CategoryElement) Convert(data interface{}) (Element, error) {
	cats := strings.Split(data.(GenericData).Value, ",")
	if len(cats) == 0 {
		return nil, errors.ErrorMessageWithStack("No categories :" + data.(GenericData).Value)
	}
	categoryName := strings.Split(cats[0], "/")[0]
	cats = koazee.StreamOf(cats).
		Map(func(cat string) string {
			if !strings.Contains(cat, categoryName) {
				return categoryName + "/" + cat
			} else {
				return cat
			}
		}).Do().Out().Val().([]string)
	allCatIds := make([]string, 0, 0)
	for _, cat := range cats {
		catIds, err := categoryIds(strings.TrimSpace(cat))
		if err != nil {
			return nil, err
		}
		for _, catId := range catIds {
			allCatIds = append(allCatIds, catId)
		}
	}
	element.CategoryIds = allCatIds
	return element, nil
}

func getAssetName(values map[string]string) string {
	assetName := ""
	for _, v := range values {
		if assetName != "" {
			assetName += "_"
		}
		assetName += v
	}
	return assetName
}

func getLocalAssetFile(imgData GenericData) (*os.File, string, error) {
	imgDataContext := imgData.Context
	assetLocation, err := imgDataContext.getValue(AssetLocation)
	if err != nil {
		return nil, "", errors.ErrorWithStack(err)
	}
	assetFullPath := assetLocation.(string) + "/" + imgData.Value
	assetExtension := filepath.Ext(assetFullPath)
	assetFile, err := os.Open(assetFullPath)
	if err != nil {
		return nil, "", errors.ErrorWithStack(err)
	} else {
		return assetFile, assetExtension, nil
	}
}

func getWebAssetFile(imgData GenericData) (*os.File, string, error) {
	assetUrl := imgData.Value
	assetExtension := filepath.Ext(assetUrl)
	response, err := http.Get(assetUrl)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, "", errors.ErrorMessageWithStack("Received non 200 response code")
	}
	file, err := ioutil.TempFile("", "acousticWebAsset")
	if err != nil {
		return nil, "", errors.ErrorWithStack(err)
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return nil, "", errors.ErrorWithStack(err)
	}

	return file, assetExtension, nil

}

func (element ImageElement) Convert(data interface{}) (Element, error) {
	imgData := data.(GenericData)
	imgDataContext := imgData.Context
	isWebUrl, err := imgDataContext.getValue(IsWebUrl)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	var assetFile *os.File
	var assetExtension string
	if isWebUrl.(bool) {
		assetFile, assetExtension, err = getWebAssetFile(imgData)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		defer assetFile.Close()
		defer os.Remove(assetFile.Name())
	} else {
		assetFile, assetExtension, err = getLocalAssetFile(imgData)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		defer assetFile.Close()
	}
	enforceImageDimension, err := imgDataContext.getValue(EnforceImageDimension)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	if enforceImageDimension.(bool) {
		imageWidth, err := imgDataContext.getUintValue(ImageWidth)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		imageHeight, err := imgDataContext.getUintValue(ImageHeight)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		ok, err := image.GetImageService().IsImageInExpectedDimension(imageWidth, imageHeight, assetFile)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		if !ok {
			resizedAsset, err := image.GetImageService().Resize(imageWidth, imageHeight, assetFile)
			if err != nil {
				return nil, errors.ErrorWithStack(err)
			}
			assetFile = resizedAsset
			defer resizedAsset.Close()
			defer os.Remove(resizedAsset.Name())
		}
	}

	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	assetName, err := imgDataContext.getValue(AssetName)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	assetNameValue := getAssetName(assetName.(map[string]string)) + assetExtension
	tags, err := imgDataContext.getValue(TagList)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	tagsValue := tags.([]string)
	acousticAssetBasePath, err := imgDataContext.getValue(AcousticAssetBasePath)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	acousticAssetPath := acousticAssetBasePath.(string) + "/" + assetNameValue
	profiles, err := imgDataContext.getValue(Profiles)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	profileValues := profiles.([]string)
	resp, err := NewAssetClient(env.AcousticAPIUrl()).Create(bufio.NewReader(assetFile), assetNameValue, tagsValue,
		acousticAssetPath, env.ContentStatus(), profileValues, env.LibraryID())
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	element.Asset = Asset{
		ID: resp.Id,
	}
	element.Mode = "shared"
	return element, nil
}
