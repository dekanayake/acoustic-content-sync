package api

import (
	"bufio"
	"fmt"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/dekanayake/acoustic-content-sync/pkg/image"
	"github.com/monmohan/xferspdy"
	"github.com/wesovilabs/koazee"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
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
	SearchValues           map[string]string
	SearchKeys             []string
	SearchType             string
	NameFields             []string
	Values                 []GenericData
	Tags                   []string
}

type GenericData struct {
	Name    string
	Type    string
	Ignore  bool
	Value   interface{}
	Context Context
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
	SearchType   string
	Type         string
	AlwaysNew    bool
	Data         []GenericData
	SearchTerm   string
	SearchValues []string
	NameFields   []string
	Tags         []string
}

type AcousticFileAsset struct {
	AssetName             map[string]string
	AcousticAssetBasePath string
	AssetLocation         string
	Tags                  []string
	IsWebUrl              bool
	Value                 string
}

type AcousticImageAsset struct {
	Profiles              []string
	EnforceImageDimension bool
	ImageWidth            uint
	ImageHeight           uint
	AcousticFileAsset
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
		AssetName:             acousticImageAsset.AssetName,
		Tags:                  acousticImageAsset.Tags,
		Value:                 acousticImageAsset.Value,
		AcousticAssetBasePath: acousticImageAsset.AcousticAssetBasePath,
		IsWebUrl:              acousticImageAsset.IsWebUrl,
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

func (acousticDataRecord AcousticDataRecord) searchQuerytoGetTheContentToUpdate() (string, error) {
	if !acousticDataRecord.Update {
		errors.ErrorMessageWithStack("Search term is available only for updatable contents")
	}
	searchValuesMap := acousticDataRecord.SearchValues
	searchValues := make([]interface{}, 0)
	for _, searchKey := range acousticDataRecord.SearchKeys {
		searchValues = append(searchValues, searchValuesMap[searchKey])
	}
	return fmt.Sprintf(acousticDataRecord.SearchTerm, searchValues...), nil
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

func (element NumberElement) Convert(data interface{}) (Element, error) {
	numValue, err := strconv.ParseInt(data.(GenericData).Value.(string), 0, 64)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	element.Value = numValue
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
	cats := strings.Split(data.(GenericData).Value.(string), env.MultipleItemsSeperator())
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
	return element, nil
}

func (element CategoryPartElement) Convert(data interface{}) (Element, error) {
	cats := strings.Split(data.(GenericData).Value.(string), env.MultipleItemsSeperator())
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

func getWebAssetFile(imgData GenericData) (string, string, error) {
	assetUrl := imgData.Value.(string)
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

func (element ImageElement) Convert(data interface{}) (Element, error) {
	imgData := data.(GenericData)
	imageValue := imgData.Value.(AcousticImageAsset)
	var assetFile *os.File
	var assetExtension string
	var err error
	if imageValue.IsWebUrl {
		assetFilePath, assetExt, err := getWebAssetFile(imgData)
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
		assetFile, assetExtension, err = getLocalAssetFile(imageValue.GetFileAsset())
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		defer assetFile.Close()
	}

	if imageValue.EnforceImageDimension {
		ok, err := image.GetImageService().IsImageInExpectedDimension(imageValue.ImageWidth, imageValue.ImageHeight, assetFile)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		if !ok {
			resizedAsset, err := image.GetImageService().Resize(imageValue.ImageWidth, imageValue.ImageHeight, assetFile)
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

	assetNameValue := getAssetName(imageValue.AssetName) + assetExtension
	acousticAssetPath := imageValue.AcousticAssetBasePath + "/" + assetNameValue
	profileValues := imageValue.Profiles
	if profileValues == nil {
		profileValues = []string{}
	}
	resp, err := NewAssetClient(env.AcousticAPIUrl()).Create(bufio.NewReader(assetFile), assetNameValue, imageValue.Tags,
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

func (element FileElement) Convert(data interface{}) (Element, error) {
	assetCreationFn := func() (Element, error) {
		fileData := data.(GenericData)
		fileValue := fileData.Value.(AcousticFileAsset)
		var assetFile *os.File
		var assetExtension string
		var err error
		if fileValue.IsWebUrl {
			assetFilePath, assetExt, err := getWebAssetFile(fileData)
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

		assetNameValue := getAssetName(fileValue.AssetName) + assetExtension

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

	assetUpdateFn := func(updatedElement Element) (Element, error) {
		fileData := data.(GenericData)
		fileValue := fileData.Value.(AcousticFileAsset)
		var assetFile *os.File
		var assetExtension string
		var err error
		if fileValue.IsWebUrl {
			assetFilePath, assetExt, err := getWebAssetFile(fileData)
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

		existingAsset, err := NewAssetClient(env.AcousticAPIUrl()).Get(updatedElement.(FileElement).Asset.ID)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		existingAssetFile, err := getExistingAssetFile(existingAsset.Path)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		newFingerprint := xferspdy.NewFingerprint(assetFile.Name(), 1024)
		oldFingerprint := xferspdy.NewFingerprint(existingAssetFile.Name(), 1024)
		isAssetsSame := newFingerprint.DeepEqual(oldFingerprint)

		if !isAssetsSame {
			assetNameValue := getAssetName(fileValue.AssetName) + "_update_" + strconv.FormatInt(time.Now().Unix(), 10) + assetExtension
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
		} else {
			element.Asset = Asset{
				ID: updatedElement.(FileElement).Asset.ID,
			}
			return element, nil
		}
	}

	assetDeleteFn := func(oldElement Element) error {
		err := NewAssetClient(env.AcousticAPIUrl()).Delete(oldElement.(FileElement).Asset.ID)
		if err != nil {
			errors.ErrorWithStack(err)
		}
		return nil
	}
	element.PreContentCreateFunctionList = []PreContentCreateFunc{assetCreationFn}
	element.PreContentUpdateFunctionList = []PreContentUpdateFunc{assetUpdateFn}
	element.PostContentUpdateFunctionList = []PostContentUpdateFunc{assetDeleteFn}
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
			Term:           query,
			ContentTypes:   []string{referenceValue.SearchType},
			Classification: "content",
		}
		searchResponse, err := NewSearchClient(env.AcousticAPIUrl()).Search(env.LibraryID(), searchRequest, Pagination{Start: 0, Rows: 1})
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
			Term:           query,
			ContentTypes:   []string{referenceValue.SearchType},
			Classification: "content",
		}
		searchResponse, err := NewSearchClient(env.AcousticAPIUrl()).Search(env.LibraryID(), searchRequest, Pagination{Start: 0, Rows: 1})
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		if searchResponse.Count == 0 {
			return nil, errors.ErrorMessageWithStack("No existing content available . content type : " + referenceValue.Type)
		}
		value.ID = searchResponse.Documents[0].Document.ID
	}
	values := make([]ReferenceValue, 0, 1)
	values = append(values, value)
	element.Values = values
	return element, nil
}
