package api

import (
	"bufio"
	"fmt"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/dekanayake/acoustic-content-sync/pkg/image"
	"github.com/monmohan/xferspdy"
	log "github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"github.com/wesovilabs/koazee"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ContextKey string

type AcousticDataRecord struct {
	CSVRecordKey           string
	Update                 bool
	CreateNonExistingItems bool
	SearchTerm             string
	SortQuery              string
	SearchTerms            map[string]string
	SearchOnLibrary        bool
	SearchOnDeliveryAPI    bool
	SearchValues           map[string]string
	SearchKeys             []string
	SearchType             string
	NameFields             []string
	Values                 []GenericData
	Tags                   []string
	// This config allows to filter records in the data csv
	FilterRecords      bool
	FilterType         string
	FilterColumns      []string
	FilterFileLocation string
}

type GenericData struct {
	Name    string
	Type    string
	Ignore  bool
	Value   interface{}
	Context Context
}

type AcousticCategory struct {
	Value     string
	Operation Operation
}

type AcousticGroup struct {
	Type string
	Data []GenericData
}

type AcousticMultiGroup struct {
	Type                          string
	ListUpdateAcousticPropertyKey string
	Data                          [][]GenericData
}

type AcousticReference struct {
	SearchType          string
	Type                string
	AlwaysNew           bool
	Data                []GenericData
	SearchTerm          string
	SearchOnLibrary     bool
	SearchOnDeliveryAPI bool
	SearchValues        []string
	NameFields          []string
	Tags                []string
	Operation           Operation
}

type AcousticMultiReference struct {
	References []AcousticReference
	Operation  Operation
}

type AssetNameConfig struct {
	AssetName               map[string]string
	AppendOriginalAssetName bool
	UseOnlyAssetName        bool
}

type AcousticFileAsset struct {
	AcousticAssetBasePath              string
	AssetLocation                      string
	Tags                               []string
	IsWebUrl                           bool
	UseExistingAsset                   bool
	AssetNameConfig                    AssetNameConfig
	Value                              string
	DontCreateAssetIfAssetNotAvailable bool
}

type AcousticImageAsset struct {
	Profiles              []string
	EnforceImageDimension bool
	ImageWidth            uint
	ImageHeight           uint
	AcousticFileAsset
}

type AcousticMultiImageAsset struct {
	Assets []AcousticImageAsset
}

func (accusticReference AcousticReference) searchQuery() (string, error) {
	values := make([]interface{}, 0)
	for _, searchValue := range accusticReference.SearchValues {
		values = append(values, searchValue)
	}
	return fmt.Sprintf(accusticReference.SearchTerm, values...), nil
}

func (acousticImageAsset AcousticImageAsset) GetFileAsset() AcousticFileAsset {
	return AcousticFileAsset{
		Tags:                  acousticImageAsset.Tags,
		Value:                 acousticImageAsset.Value,
		AcousticAssetBasePath: acousticImageAsset.AcousticAssetBasePath,
		IsWebUrl:              acousticImageAsset.IsWebUrl,
		UseExistingAsset:      acousticImageAsset.UseExistingAsset,
		AssetNameConfig:       acousticImageAsset.AssetNameConfig,
		AssetLocation:         acousticImageAsset.AssetLocation,
	}
}

type Context struct {
	Data map[ContextKey]interface{}
}

const (
	LinkToParents ContextKey = "LinkToParents"
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
			return 0, errors.ErrorMessageWithStack("Cannot cast the provided value " + reflect.TypeOf(val).String())
		}
	} else {
		return 0, errors.ErrorMessageWithStack("no context value found for key :" + string(key))
	}
}

func (context Context) getBoolValue(key ContextKey) (bool, error) {
	if val, ok := context.Data[key]; ok {
		if castVal, ok := val.(bool); ok {
			return castVal, nil
		} else {
			return false, errors.ErrorMessageWithStack("Cannot cast the provided value " + reflect.TypeOf(val).String())
		}
	} else {
		return false, errors.ErrorMessageWithStack("no context value found for key :" + string(key))
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
				acc += data.Value.(string)
			} else {
				acc += "__" + data.Value.(string)
			}
			return acc
		}).String()
}

func (acousticDataRecord AcousticDataRecord) CSVRecordKeyValue() string {
	return koazee.StreamOf(acousticDataRecord.Values).
		Filter(func(columnValue GenericData) bool {
			return columnValue.Name == acousticDataRecord.CSVRecordKey
		}).Map(func(columnValue GenericData) string {
		return columnValue.Value.(string)
	}).First().String()
}

func (acousticDataRecord AcousticDataRecord) GetValue(columnName string) interface{} {
	return koazee.StreamOf(acousticDataRecord.Values).
		Filter(func(columnValue GenericData) bool {
			return columnValue.Name == columnName
		}).Map(func(columnValue GenericData) interface{} {
		return columnValue.Value
	}).First().Val()
}

func (acousticDataRecord AcousticDataRecord) searchQuerytoGetTheContentToUpdate() (map[string]string, error) {
	result := make(map[string]string, 0)
	if !acousticDataRecord.Update {
		errors.ErrorMessageWithStack("Search term is available only for updatable contents")
	}
	searchValuesMap := acousticDataRecord.SearchValues
	searchValues := make([]interface{}, 0)
	for _, searchKey := range acousticDataRecord.SearchKeys {
		searchValues = append(searchValues, searchValuesMap[searchKey])
	}
	if acousticDataRecord.SearchTerm != "" {
		result["q"] = fmt.Sprintf(acousticDataRecord.SearchTerm, searchValues...)
	} else {
		for searchKey, searchTerm := range acousticDataRecord.SearchTerms {
			result[searchKey] = fmt.Sprintf(searchTerm, searchValues...)
		}
	}

	if acousticDataRecord.SearchTerm != "" {

	}

	return result, nil
}

func (element TextElement) Convert(data interface{}) (Element, error) {
	element.Value = data.(GenericData).Value.(string)
	return element, nil
}

func (element FormattedTextElement) Convert(data interface{}) (Element, error) {
	element.Value = data.(GenericData).Value.(string)
	return element, nil
}

func (element BooleanElement) Convert(data interface{}) (Element, error) {
	val, err := strconv.ParseBool(data.(GenericData).Value.(string))
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	element.Value = val
	return element, nil
}

func (element MultiTextElement) Convert(data interface{}) (Element, error) {
	element.Values = strings.Split(data.(GenericData).Value.(string), env.MultipleItemsSeperator())
	return element, nil
}

func (element OptionSelectionElement) Convert(data interface{}) (Element, error) {
	panic("implement me")
}

func (element MultiOptionSelectionElement) Convert(data interface{}) (Element, error) {
	options := strings.Split(data.(GenericData).Value.(string), env.MultipleItemsSeperator())
	options = funk.UniqString(options)
	optionSelections := make([]OptionSelectionValue, 0)
	for _, option := range options {
		optionSelections = append(optionSelections, OptionSelectionValue{
			Selection: option,
		})
	}
	element.Values = optionSelections
	return element, nil
}

func (element NumberElement) Convert(data interface{}) (Element, error) {
	numValue, err := strconv.ParseInt(data.(GenericData).Value.(string), 0, 64)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	element.Value = numValue
	return element, nil
}

func (element MultiNumberElement) Convert(data interface{}) (Element, error) {
	numbersInStrings := strings.Split(data.(GenericData).Value.(string), env.MultipleItemsSeperator())
	numValues := make([]int64, 0)
	for _, numberValInStr := range numbersInStrings {
		numValue, err := strconv.ParseInt(numberValInStr, 0, 64)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		numValues = append(numValues, numValue)
	}
	element.Values = numValues
	return element, nil
}

func (element FloatElement) Convert(data interface{}) (Element, error) {
	numValue, err := strconv.ParseFloat(data.(GenericData).Value.(string), 32)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	element.Value = math.Ceil(numValue*10000) / 10000
	return element, nil
}

func (element LinkElement) Convert(data interface{}) (Element, error) {
	element.LinkURL = data.(GenericData).Value.(string)
	return element, nil
}

func categoryIds(category string) ([]string, error) {
	catItems := strings.Split(category, env.CategoryHierarchySeperator())
	if len(catItems) == 1 {
		return nil, errors.ErrorMessageWithStack("empty category :" + catItems[0])
	}

	categoryItems, err := NewCachedCategoryClient(env.AcousticAPIUrl()).Categories(catItems[0])
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
					acc += env.CategoryHierarchySeperator() + catNamePath
				}
				return acc
			}).String()
		catNamePaths = append(catNamePaths, catNamePath)
	}

	catIds := koazee.StreamOf(categoryItems).
		Filter(func(categoryItem CategoryItem) bool {
			fullNamePath := categoryItem.FullNamePath()
			contains, _ := koazee.StreamOf(catNamePaths).Contains(fullNamePath)
			return strings.Contains(fullNamePath, env.CategoryHierarchySeperator()) && contains
		}).
		Map(func(categoryItem CategoryItem) string {
			return categoryItem.Id
		}).Out().Val().([]string)

	return catIds, nil
}

func catIdFromCatPart(catPart string, linkToParent bool) ([]string, error) {
	catItems := strings.Split(catPart, env.CategoryHierarchySeperator())
	if len(catItems) == 1 {
		return nil, errors.ErrorMessageWithStack("empty category :" + catItems[0])
	}
	categoryItems, err := NewCachedCategoryClient(env.AcousticAPIUrl()).Categories(catItems[0])
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	for _, catItem := range categoryItems {
		if strings.Contains(catItem.Name, catItems[1]) {
			if linkToParent {
				return categoryIds(catItem.FullNamePath())
			} else {
				return []string{catItem.Id}, nil
			}
		}
	}
	return nil, errors.ErrorMessageWithStack("no category matched with the given cat part :" + catPart)

}

func (element CategoryElement) Convert(data interface{}) (Element, error) {
	cats := strings.Split(data.(GenericData).Value.(AcousticCategory).Value, env.MultipleItemsSeperator())
	if len(cats) == 0 {
		return nil, errors.ErrorMessageWithStack("No categories :" + data.(GenericData).Value.(string))
	}
	categoryName := strings.Split(cats[0], env.CategoryHierarchySeperator())[0]
	cats = koazee.StreamOf(cats).
		Map(func(cat string) string {
			if !strings.Contains(cat, categoryName) {
				return categoryName + env.CategoryHierarchySeperator() + cat
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
	element.Operation = data.(GenericData).Value.(AcousticCategory).Operation
	return element, nil
}

func (element CategoryPartElement) Convert(data interface{}) (Element, error) {
	cats := strings.Split(data.(GenericData).Value.(AcousticCategory).Value, env.MultipleItemsSeperator())
	linkToParents, err := data.(GenericData).Context.getBoolValue(LinkToParents)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	if len(cats) == 0 {
		return nil, errors.ErrorMessageWithStack("No categories :" + data.(GenericData).Value.(string))
	}
	allCatIdsMap := make(map[string]bool)
	for _, cat := range cats {
		catIds, err := catIdFromCatPart(cat, linkToParents)
		if err != nil {
			return element, errors.ErrorWithStack(err)
		}
		for _, catId := range catIds {
			if _, ok := allCatIdsMap[catId]; !ok {
				allCatIdsMap[catId] = true
			}
		}
	}
	allCatIds := make([]string, 0, len(allCatIdsMap))
	for key := range allCatIdsMap {
		allCatIds = append(allCatIds, key)
	}
	element.CategoryIds = allCatIds
	element.Operation = data.(GenericData).Value.(AcousticCategory).Operation
	return element, nil
}

func extractAssetNameFromAssetPath(path string) (string, error) {
	extractedValue := path
	regxList := [2]string{"(\\w+-*_*\\s*)+\\.(png|jpg)", "(\\w+-*_*\\s*)+"}
	for _, regx := range regxList {
		compiledRegx, err := regexp.Compile(regx)
		if err != nil {
			return "", err
		}
		extractedValue = compiledRegx.FindString(extractedValue)
	}
	return extractedValue, nil
}

func getAssetName(asset AcousticFileAsset) (string, error) {
	if asset.AssetNameConfig.UseOnlyAssetName {
		assetName, err := extractAssetNameFromAssetPath(asset.Value)
		if err != nil {
			return "", errors.ErrorWithStack(err)
		}
		return assetName, nil
	}

	names := asset.AssetNameConfig.AssetName
	assetName := ""
	for _, v := range names {
		if assetName != "" {
			assetName += "_"
		}
		assetName += v
	}
	if asset.AssetNameConfig.AppendOriginalAssetName {
		assetNameFromPath, err := extractAssetNameFromAssetPath(asset.Value)
		if err != nil {
			return "", errors.ErrorWithStack(err)
		}
		assetName = assetName + "_" + assetNameFromPath
	}
	return assetName, nil
}

func getLocalAssetFile(imgData AcousticFileAsset) (*os.File, string, error) {

	assetFullPath := imgData.AssetLocation + "/" + imgData.Value
	assetExtension := filepath.Ext(assetFullPath)
	assetFile, err := os.Open(assetFullPath)
	if err != nil {
		return nil, "", errors.ErrorWithStack(err)
	} else {
		return assetFile, assetExtension, nil
	}
}

func checkAssetToUploadExists(asset AcousticFileAsset) (bool, error) {
	if asset.IsWebUrl {
		client := http.Client{
			Timeout: 10 * time.Second,
		}
		response, err := client.Get(asset.Value)
		if err != nil {
			return false, err
		}
		if response.StatusCode == 200 {
			return true, nil
		}

	} else {
		assetFullPath := asset.Value + "/" + asset.Value
		file, _ := os.Open(assetFullPath)
		if file != nil {
			return true, nil
		}

	}
	log.Info("the asset in the path not available , since ignoring the asset. path :" + asset.Value)
	return false, nil
}

func getWebAssetFile(fileAsset AcousticFileAsset) (string, string, error) {
	assetUrl := fileAsset.Value
	assetExtension := filepath.Ext(assetUrl)
	response, err := http.Get(assetUrl)
	if err != nil {
		return "", "", err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "", "", errors.ErrorMessageWithStack("Received non 200 response code")
	}
	file, err := ioutil.TempFile("", "acousticWebAsset")
	if err != nil {
		return "", "", errors.ErrorWithStack(err)
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", "", errors.ErrorWithStack(err)
	}

	return file.Name(), assetExtension, nil
}

func getExistingAssetFile(filePath string) (*os.File, error) {
	response, err := http.Get(env.AcousticBaseUrl() + filePath)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, errors.ErrorMessageWithStack("Received non 200 response code")
	}
	file, err := ioutil.TempFile("", "acousticWebAsset")
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	return file, nil
}

var isSameAsset = func(assetId string, newAssetName string) (bool, error) {
	existingAsset, err := NewAssetClient(env.AcousticAPIUrl()).Get(assetId)
	if err != nil {
		return false, errors.ErrorWithStack(err)
	}
	existingAssetFile, err := getExistingAssetFile(existingAsset.Path)
	if err != nil {
		return false, errors.ErrorWithStack(err)
	}
	newFingerprint := xferspdy.NewFingerprint(newAssetName, 1024)
	oldFingerprint := xferspdy.NewFingerprint(existingAssetFile.Name(), 1024)
	return newFingerprint.DeepEqual(oldFingerprint), nil
}

var getImageFunc = func(imageValue AcousticImageAsset) (*os.File, *os.File, string, error) {
	var assetFile *os.File
	var tmpFile *os.File
	var assetExtension string
	var err error
	if imageValue.IsWebUrl {
		assetFilePath, assetExt, err := getWebAssetFile(imageValue.AcousticFileAsset)
		assetExtension = assetExt
		if err != nil {
			return nil, nil, "", errors.ErrorWithStack(err)
		}
		assetFile, err = os.Open(assetFilePath)
		if err != nil {
			return nil, nil, "", errors.ErrorWithStack(err)
		}
		tmpFile = assetFile

	} else {
		assetFile, assetExtension, err = getLocalAssetFile(imageValue.GetFileAsset())
		if err != nil {
			return nil, nil, "", errors.ErrorWithStack(err)
		}
	}

	if imageValue.EnforceImageDimension {
		ok, err := image.GetImageService().IsImageInExpectedDimension(imageValue.ImageWidth, imageValue.ImageHeight, assetFile)
		if err != nil {
			return nil, nil, "", errors.ErrorWithStack(err)
		}
		if !ok {
			resizedAsset, err := image.GetImageService().Resize(imageValue.ImageWidth, imageValue.ImageHeight, assetFile)
			if err != nil {
				return nil, nil, "", errors.ErrorWithStack(err)
			}
			assetFile = resizedAsset
			tmpFile = resizedAsset
		}
	}

	if err != nil {
		return nil, nil, "", errors.ErrorWithStack(err)
	}
	return assetFile, tmpFile, assetExtension, nil
}

func (element MultiImageElement) Convert(data interface{}) (Element, error) {
	imgData := data.(GenericData)
	multiImage := imgData.Value.(AcousticMultiImageAsset)
	imageCreationFn := func() (Element, error) {
		assets := make([]ImageElementItem, 0)
		for _, imageAsset := range multiImage.Assets {
			getOrCreateImageAssetFn := func() (string, bool, error) {
				id, assetExist, cleanUpFunc, err := getOrCreateImageAsset(imageAsset, data)
				if !assetExist {
					return "", false, nil
				}
				if err != nil {
					return "", false, err
				}
				if cleanUpFunc != nil {
					defer cleanUpFunc()
				}
				return id, true, nil
			}

			id, assetExist, err := getOrCreateImageAssetFn()
			if err != nil {
				return nil, err
			}
			if !assetExist {
				continue
			}

			assets = append(assets, ImageElementItem{
				Asset: Asset{
					ID: id,
				},
				Mode: "shared",
			})
		}
		if len(assets) == 0 {
			return nil, nil
		}
		element.Values = assets
		return element, nil
	}

	imageUpdateFn := func(updatedElement Element) (Element, []PostContentUpdateFunc, error) {
		assets := make([]ImageElementItem, 0)
		postContentUpdateFunctions := make([]PostContentUpdateFunc, 0)
		for _, imageAsset := range multiImage.Assets {
			getOrUpdateImageAssetFunc := func() (string, []PostContentUpdateFunc, error) {
				id, cleanUpFunc, postContentUpdateFuncs, err := getOrUpdateImageAsset(imageAsset, updatedElement, data)
				if err != nil {
					return "", nil, err
				}
				if cleanUpFunc != nil {
					defer cleanUpFunc()
				}
				return id, postContentUpdateFuncs, nil
			}

			id, postContentUpdateFuncs, err := getOrUpdateImageAssetFunc()
			if err != nil {
				return nil, nil, err
			}
			if postContentUpdateFunctions != nil {
				postContentUpdateFunctions = append(postContentUpdateFunctions, postContentUpdateFuncs...)
			}

			assets = append(assets, ImageElementItem{
				Asset: Asset{
					ID: id,
				},
				Mode: "shared",
			})
		}
		element.Values = assets
		return element, postContentUpdateFunctions, nil
	}
	element.PreContentCreateFunctionList = []PreContentCreateFunc{imageCreationFn}
	element.PreContentUpdateFunctionList = []PreContentUpdateFunc{imageUpdateFn}
	return element, nil
}

func (element ImageElement) Convert(data interface{}) (Element, error) {
	imageCreationFn := func() (Element, error) {
		imgData := data.(GenericData)
		imageValue := imgData.Value.(AcousticImageAsset)
		id, assetExist, cleanUpFunc, err := getOrCreateImageAsset(imageValue, data)
		if err != nil {
			return nil, err
		}
		if cleanUpFunc != nil {
			defer cleanUpFunc()
		}
		if !assetExist {
			return nil, nil
		}

		element.Asset = Asset{
			ID: id,
		}
		element.Mode = "shared"
		return element, nil
	}

	imageUpdateFn := func(updatedElement Element) (Element, []PostContentUpdateFunc, error) {
		imgData := data.(GenericData)
		imageValue := imgData.Value.(AcousticImageAsset)
		id, cleanUpFunc, postContentUpdateFuncs, err := getOrUpdateImageAsset(imageValue, updatedElement, data)
		if err != nil {
			return nil, nil, err
		}
		if cleanUpFunc != nil {
			defer cleanUpFunc()
		}

		element.Asset = Asset{
			ID: id,
		}
		element.Mode = "shared"
		return element, postContentUpdateFuncs, nil
	}

	element.PreContentCreateFunctionList = []PreContentCreateFunc{imageCreationFn}
	element.PreContentUpdateFunctionList = []PreContentUpdateFunc{imageUpdateFn}
	return element, nil

}

func getOrCreateImageAsset(imageValue AcousticImageAsset, data interface{}) (string, bool, func(), error) {
	var isAssetExist = false
	var id = ""
	var cleanUpFunc func() = nil
	if imageValue.UseExistingAsset {
		var err error = nil
		isAssetExist, id, err = checkAssetExist(imageValue)
		if err != nil {
			return "", false, cleanUpFunc, err
		}
	}

	if !isAssetExist {
		if imageValue.DontCreateAssetIfAssetNotAvailable {
			exist, err := checkAssetToUploadExists(imageValue.AcousticFileAsset)
			if err != nil {
				return "", false, cleanUpFunc, err
			}
			if !exist {
				return "", false, cleanUpFunc, nil
			}
		}
		assetFile, tmpFile, assetExtension, err := getImageFunc(imageValue)
		cleanUpFunc = func() {
			assetFile.Close()
			if tmpFile != nil {
				os.Remove(tmpFile.Name())
			}
		}

		if err != nil {
			return "", false, cleanUpFunc, err
		}
		assetName, err := getAssetName(imageValue.AcousticFileAsset)
		if err != nil {
			return "", false, cleanUpFunc, err
		}
		assetNameValue := assetName + assetExtension
		acousticAssetPath := imageValue.AcousticAssetBasePath + "/" + assetNameValue
		profileValues := imageValue.Profiles
		if profileValues == nil {
			profileValues = []string{}
		}
		resp, err := NewAssetClient(env.AcousticAPIUrl()).Create(bufio.NewReader(assetFile), assetNameValue, imageValue.Tags,
			acousticAssetPath, env.ContentStatus(), profileValues, env.LibraryID())
		if err != nil {
			return "", false, cleanUpFunc, errors.ErrorWithStack(err)
		}
		NewCacheRepository().PutCache(AssetCache, resp.Path, resp.Id)
		id = resp.Id
	}
	return id, true, cleanUpFunc, nil
}

func getOrUpdateImageAsset(imageValue AcousticImageAsset, updatedElement Element, data interface{}) (string, func(), []PostContentUpdateFunc, error) {
	var isAssetExist = false
	var id = ""
	var cleanUpFunc func() = nil
	var postContentUpdateFuncs []PostContentUpdateFunc = nil
	if imageValue.UseExistingAsset {
		var err error = nil
		isAssetExist, id, err = checkAssetExist(imageValue)
		if err != nil {
			return "", cleanUpFunc, nil, err
		}
	}

	if !isAssetExist {
		assetFile, tmpFile, assetExtension, err := getImageFunc(imageValue)
		cleanUpFunc = func() {
			assetFile.Close()
			if tmpFile != nil {
				os.Remove(tmpFile.Name())
			}
		}
		if err != nil {
			return "", cleanUpFunc, nil, err
		}
		oldAssetId := updatedElement.(ImageElement).Asset.ID
		isSameImage, err := isSameAsset(oldAssetId, assetFile.Name())
		if err != nil {
			return "", cleanUpFunc, nil, err
		}

		if !isSameImage {
			imgData := data.(GenericData)
			imageValue := imgData.Value.(AcousticImageAsset)
			assetName, err := getAssetName(imageValue.AcousticFileAsset)
			if err != nil {
				return "", cleanUpFunc, nil, errors.ErrorWithStack(err)
			}
			assetNameValue := assetName + "_update_" + strconv.FormatInt(time.Now().Unix(), 10) + assetExtension
			acousticAssetPath := imageValue.AcousticAssetBasePath + "/" + assetNameValue
			resp, err := NewAssetClient(env.AcousticAPIUrl()).Create(bufio.NewReader(assetFile), assetNameValue, imageValue.Tags,
				acousticAssetPath, env.ContentStatus(), []string{}, env.LibraryID())
			NewCacheRepository().PutCache(AssetCache, resp.Path, resp.Id)
			postUpdateFunc := func() error {
				err := NewAssetClient(env.AcousticAPIUrl()).Delete(oldAssetId)
				if err != nil {
					return errors.ErrorWithStack(err)
				}
				return nil
			}
			if err != nil {
				return "", cleanUpFunc, nil, errors.ErrorWithStack(err)
			}
			id = resp.Id
			postContentUpdateFuncs = []PostContentUpdateFunc{postUpdateFunc}
		} else {
			id = updatedElement.(FileElement).Asset.ID
		}
	}
	return id, cleanUpFunc, postContentUpdateFuncs, nil
}

func checkAssetExist(imageValue AcousticImageAsset) (bool, string, error) {
	assetFile, tmpFile, assetExtension, err := getImageFunc(imageValue)
	defer assetFile.Close()
	if tmpFile != nil {
		defer os.Remove(tmpFile.Name())
	}
	if err != nil {
		return false, "", err
	}
	assetName, err := getAssetName(imageValue.AcousticFileAsset)
	if err != nil {
		return false, "", err
	}
	assetNameValue := assetName + assetExtension
	path := imageValue.AcousticAssetBasePath + "/" + assetNameValue
	var assetResponse *AssetResponse = nil
	var isAssetExist = false
	var id = ""
	assetId, err := NewCacheRepository().GetCache(AssetCache, path)
	if err != nil {
		return false, "", err
	}
	if assetId != nil {
		isAssetExist = true
		id = assetId.(string)
	}
	if !isAssetExist {
		isAssetExist, assetResponse, err = NewAssetClient(env.AcousticAPIUrl()).GetByPath(path)
		if err != nil {
			return false, "", err
		}
		if isAssetExist {
			id = assetResponse.ID
		}
	}
	return isAssetExist, id, nil
}

func (element FileElement) Convert(data interface{}) (Element, error) {
	assetCreationFn := func() (Element, error) {
		fileData := data.(GenericData)
		fileValue := fileData.Value.(AcousticFileAsset)
		var assetFile *os.File
		var assetExtension string
		var err error
		if fileValue.IsWebUrl {
			assetFilePath, assetExt, err := getWebAssetFile(fileValue)
			assetExtension = assetExt
			if err != nil {
				return nil, errors.ErrorWithStack(err)
			}
			assetFile, err = os.Open(assetFilePath)
			if err != nil {
				return nil, errors.ErrorWithStack(err)
			}
			defer assetFile.Close()
			defer os.Remove(assetFile.Name())
		} else {
			assetFile, assetExtension, err = getLocalAssetFile(fileValue)
			if err != nil {
				return nil, errors.ErrorWithStack(err)
			}
			defer assetFile.Close()
		}

		assetName, err := getAssetName(fileValue)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		assetNameValue := assetName + assetExtension

		acousticAssetPath := fileValue.AcousticAssetBasePath + "/" + assetNameValue
		resp, err := NewAssetClient(env.AcousticAPIUrl()).Create(bufio.NewReader(assetFile), assetNameValue, fileValue.Tags,
			acousticAssetPath, env.ContentStatus(), []string{}, env.LibraryID())
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		element.Asset = Asset{
			ID: resp.Id,
		}
		return element, nil
	}

	assetUpdateFn := func(updatedElement Element) (Element, []PostContentUpdateFunc, error) {
		fileData := data.(GenericData)
		fileValue := fileData.Value.(AcousticFileAsset)
		var assetFile *os.File
		var assetExtension string
		var err error
		if fileValue.IsWebUrl {
			assetFilePath, assetExt, err := getWebAssetFile(fileValue)
			assetExtension = assetExt
			if err != nil {
				return nil, nil, errors.ErrorWithStack(err)
			}
			assetFile, err = os.Open(assetFilePath)
			if err != nil {
				return nil, nil, errors.ErrorWithStack(err)
			}
			defer assetFile.Close()
			defer os.Remove(assetFile.Name())
		} else {
			assetFile, assetExtension, err = getLocalAssetFile(fileValue)
			if err != nil {
				return nil, nil, errors.ErrorWithStack(err)
			}
			defer assetFile.Close()
		}
		oldAssetId := updatedElement.(FileElement).Asset.ID
		existingAsset, err := NewAssetClient(env.AcousticAPIUrl()).Get(updatedElement.(FileElement).Asset.ID)
		if err != nil {
			return nil, nil, errors.ErrorWithStack(err)
		}
		existingAssetFile, err := getExistingAssetFile(existingAsset.Path)
		if err != nil {
			return nil, nil, errors.ErrorWithStack(err)
		}
		newFingerprint := xferspdy.NewFingerprint(assetFile.Name(), 1024)
		oldFingerprint := xferspdy.NewFingerprint(existingAssetFile.Name(), 1024)
		isAssetsSame := newFingerprint.DeepEqual(oldFingerprint)

		if !isAssetsSame {
			assetName, err := getAssetName(fileValue)
			if err != nil {
				return nil, nil, errors.ErrorWithStack(err)
			}
			assetNameValue := assetName + "_update_" + strconv.FormatInt(time.Now().Unix(), 10) + assetExtension
			acousticAssetPath := fileValue.AcousticAssetBasePath + "/" + assetNameValue
			resp, err := NewAssetClient(env.AcousticAPIUrl()).Create(bufio.NewReader(assetFile), assetNameValue, fileValue.Tags,
				acousticAssetPath, env.ContentStatus(), []string{}, env.LibraryID())
			if err != nil {
				return nil, nil, errors.ErrorWithStack(err)
			}
			postUpdateFunc := func() error {
				err := NewAssetClient(env.AcousticAPIUrl()).Delete(oldAssetId)
				if err != nil {
					errors.ErrorWithStack(err)
				}
				return nil
			}
			element.Asset = Asset{
				ID: resp.Id,
			}
			return element, []PostContentUpdateFunc{postUpdateFunc}, nil
		} else {
			element.Asset = Asset{
				ID: updatedElement.(FileElement).Asset.ID,
			}
			return element, nil, nil
		}
	}

	element.PreContentCreateFunctionList = []PreContentCreateFunc{assetCreationFn}
	element.PreContentUpdateFunctionList = []PreContentUpdateFunc{assetUpdateFn}
	return element, nil
}

func (element GroupElement) Convert(data interface{}) (Element, error) {
	groupData := data.(GenericData)
	groupValue := groupData.Value.(AcousticGroup)
	element.TypeRef = map[string]string{
		"id": groupValue.Type,
	}
	values := make(map[string]interface{}, len(groupValue.Data))
	for _, dataItem := range groupValue.Data {
		if dataItem.Ignore {
			continue
		}
		if dataItem.Value == nil {
			continue
		}
		element, err := Build(dataItem.Type)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		element, err = element.Convert(dataItem)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		values[dataItem.Name] = element
	}
	element.Value = values
	return element, nil
}

func (element MultiGroupElement) Convert(data interface{}) (Element, error) {
	groupData := data.(GenericData)
	groupValue := groupData.Value.(AcousticMultiGroup)
	element.TypeRef = map[string]string{
		"id": groupValue.Type,
	}
	groupValues := make([]map[string]interface{}, 0, len(groupValue.Data))
	for _, dataItemGroup := range groupValue.Data {
		values := make(map[string]interface{}, len(dataItemGroup))
		for _, dataItem := range dataItemGroup {
			if dataItem.Ignore {
				continue
			}
			if dataItem.Value == nil {
				continue
			}
			element, err := Build(dataItem.Type)
			if err != nil {
				return nil, errors.ErrorWithStack(err)
			}
			element, err = element.Convert(dataItem)
			if err != nil {
				return nil, errors.ErrorWithStack(err)
			}
			values[dataItem.Name] = element
		}
		groupValues = append(groupValues, values)
	}
	element.Values = groupValues
	return element, nil
}

func (element ReferenceElement) Convert(data interface{}) (Element, error) {
	referenceData := data.(GenericData)
	referenceValue := referenceData.Value.(AcousticReference)
	value := ReferenceValue{}
	if referenceValue.AlwaysNew {
		acousticDataRecord := AcousticDataRecord{
			Values:     referenceValue.Data,
			NameFields: referenceValue.NameFields,
			Tags:       referenceValue.Tags,
		}
		contentCreateResponse, err := NewContentService(env.AcousticAuthUrl(), env.LibraryID()).CreateOrUpdateContentWithRetry(acousticDataRecord, referenceValue.Type)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		value.ID = contentCreateResponse.Id
	} else {
		query, err := referenceValue.searchQuery()
		if err != nil {
			return nil, err
		}
		searchRequest := SearchRequest{
			Terms:          map[string]string{"q": query},
			ContentTypes:   []string{referenceValue.SearchType},
			Classification: "content",
		}
		searchResponse, err := NewSearchClient(env.AcousticAPIUrl()).Search(env.LibraryID(), referenceValue.SearchOnLibrary, referenceValue.SearchOnDeliveryAPI, searchRequest, Pagination{Start: 0, Rows: 1})
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		if searchResponse.Count == 0 {
			return nil, errors.ErrorMessageWithStack("No existing content available . content type : " + referenceValue.Type)
		}
		value.ID = searchResponse.Documents[0].Document.ID
	}

	element.Value = value
	return element, nil
}

func (element MultiReferenceElement) Convert(data interface{}) (Element, error) {
	referenceData := data.(GenericData)
	acousticMultiReference := referenceData.Value.(AcousticMultiReference)

	values := make([]ReferenceValue, 0)
	for _, referenceValue := range acousticMultiReference.References {
		value := ReferenceValue{}
		if referenceValue.AlwaysNew {
			acousticDataRecord := AcousticDataRecord{
				Values:     referenceValue.Data,
				NameFields: referenceValue.NameFields,
				Tags:       referenceValue.Tags,
			}
			contentCreateResponse, err := NewContentService(env.AcousticAuthUrl(), env.LibraryID()).CreateOrUpdateContentWithRetry(acousticDataRecord, referenceValue.Type)
			if err != nil {
				return nil, errors.ErrorWithStack(err)
			}
			value.ID = contentCreateResponse.Id
		} else {
			query, err := referenceValue.searchQuery()
			if err != nil {
				return nil, err
			}
			searchRequest := SearchRequest{
				Terms:          map[string]string{"q": query},
				ContentTypes:   []string{referenceValue.SearchType},
				Classification: "content",
			}
			searchResponse, err := NewSearchClient(env.AcousticAPIUrl()).Search(env.LibraryID(), true, referenceValue.SearchOnDeliveryAPI, searchRequest, Pagination{Start: 0, Rows: 1})
			if err != nil {
				return nil, errors.ErrorWithStack(err)
			}
			if searchResponse.Count == 0 {
				return nil, errors.ErrorMessageWithStack("No existing content available . content type : " + referenceValue.Type)
			}
			value.ID = searchResponse.Documents[0].Document.ID
		}
		values = append(values, value)
	}

	element.Values = values
	element.Operation = acousticMultiReference.Operation
	return element, nil
}
