package api

import (
	"bufio"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/wesovilabs/koazee"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ContextKey string

type AcousticDataRecord struct {
	NameFields []string
	Values     []GenericData
	Tags       []string
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
)

func (context Context) getValue(key ContextKey) (interface{}, error) {
	if val, ok := context.Data[key]; ok {
		return val, nil
	} else {
		return nil, errors.ErrorMessageWithStack("no context value found for key :" + string(key))
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

func (element CategoryElement) Convert(data interface{}) (Element, error) {
	catItems := strings.Split(data.(GenericData).Value, "/")
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

	element.CategoryIds = catIds
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

func (element ImageElement) Convert(data interface{}) (Element, error) {
	imgData := data.(GenericData)
	imgDataContext := imgData.Context
	assetLocation, err := imgDataContext.getValue(AssetLocation)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	assetFullPath := assetLocation.(string) + "/" + imgData.Value
	assetExtension := filepath.Ext(assetFullPath)
	assetFile, err := os.Open(assetFullPath)
	defer assetFile.Close()
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
		acousticAssetPath, env.ContentStatus(), profileValues)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	element.Asset = Asset{
		ID: resp.Id,
	}
	element.Mode = "shared"
	return element, nil
}
