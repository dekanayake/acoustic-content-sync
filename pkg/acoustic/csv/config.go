package csv

import (
	"encoding/json"
	"fmt"
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/author/api"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/goccy/go-yaml"
	"github.com/jinzhu/copier"
	log "github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"github.com/wesovilabs/koazee"
	"io/ioutil"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ContentType interface {
	GetFieldMapping(csvField string) (ContentFieldMapping, error)
}

type ContentTypesMapping struct {
	ContentType     []ContentTypeMapping `yaml:"contentType"`
	CategoryMapping []CategoryMapping    `yaml:"category"`
	DeleteMapping   []DeleteMapping      `yaml:"delete"`
	SiteMapping     []SiteMapping        `yaml:"site"`
}

type ContentTypeMapping struct {
	Type                   string                `yaml:"type"`
	FieldMapping           []ContentFieldMapping `yaml:"fieldMapping"`
	Name                   []string              `yaml:"name"`
	Tags                   []string              `yaml:"tags"`
	CsvRecordKey           string                `yaml:"csvRecordKey"`
	Update                 bool                  `yaml:"update"`
	FeedType               api.FeedType          `yaml:"feedType"`
	CreateNonExistingItems bool                  `yaml:"createNonExistingItems"`
	SearchOnLibrary        bool                  `yaml:"searchOnLibrary"`
	SearchTerm             string                `yaml:"searchTerm"`
	SearchTerms            map[string]string     `yaml:"searchTerms"`
	SearchKeys             []string              `yaml:"searchKeys"`
	SearchType             string                `yaml:"searchType"`
	SearchOnDeliveryAPI    bool                  `yaml:"searchOnDeliveryAPI"`
	PaginationRows         int                   `yaml:"paginationRows"`
	// This config allows to filter records in the data csv
	FilterRecords      bool     `yaml:"filterRecords"`
	FilterType         string   `yaml:"filterType"`
	FilterColumns      []string `yaml:"filterColumns"`
	FilterFileLocation string   `yaml:"filterFileLocation"`
}

type SiteMapping struct {
	ContentTypeMapping    `yaml:",inline"`
	DontCreatePageIfExist bool `yaml:"dontCreatePageIfExist"`
	UpdatePageIfExist     bool `yaml:"updatePageIfExist"`
}

type CategoryMapping struct {
	Parent string `yaml:"parent"`
	Column string `yaml:"column"`
}

type DeleteMapping struct {
	Name          string        `yaml:"name"`
	AssetType     api.AssetType `yaml:"assetType"`
	SearchMapping SearchMapping `yaml:"search"`
}

type SearchMapping struct {
	ContentType    string `yaml:"contentType"`
	Classification string `yaml:"classification"`
	SearchTerm     string `yaml:"searchTerm"`
}

type AssetNameConfig struct {
	AssetName               []RefPropertyMapping `yaml:"assetName"`
	AppendOriginalAssetName bool                 `yaml:"appendOriginalAssetName"`
	UseOnlyAssetName        bool                 `yaml:"useOnlyAssetName"`
}

type SanitizeConfig struct {
	Sanitize bool     `yaml:"sanitize"`
	Regx     []string `yaml:"regx"`
}

type ContentFieldMapping struct {
	CsvProperty      string         `yaml:"csvProperty"`
	Ignore           bool           `yaml:"ignore"`
	ValuePattern     string         `yaml:"valuePattern"`
	SanitizeConfig   SanitizeConfig `yaml:"sanitizeConfig"`
	RegxMatchedGroup uint           `yaml:"regXMatchedGroup"`
	Regx             []string       `yaml:"regx"`
	Mandatory        bool           `yaml:"mandatory"`
	StaticValue      string         `yaml:"staticValue"`
	JoinedValue      string         `yaml:"joinedValue"`
	AcousticProperty string         `yaml:"acousticProperty"`
	AcousticID       bool           `yaml:"acousticID"`
	PropertyType     string         `yaml:"propertyType"`
	CategoryName     string         `yaml:"categoryName"`
	LoadFromFile     bool           `yaml:"loadFromFile"`

	AssetName                          AssetNameConfig `yaml:"assetNameConfig"`
	Profiles                           []string        `yaml:"profiles"`
	AcousticAssetBasePath              string          `yaml:"acousticAssetBasePath"`
	AssetLocation                      string          `yaml:"assetLocation"`
	IsWebUrl                           bool            `yaml:"isWebUrl"`
	ImageWidth                         uint            `yaml:"imageWidth"`
	UseExistingAsset                   bool            `yaml:"useExistingAsset"`
	ImageHeight                        uint            `yaml:"imageHeight"`
	EnforceImageDimension              bool            `yaml:"enforceImageDimension"`
	Operation                          api.Operation   `yaml:"operation"`
	DontCreateAssetIfAssetNotAvailable bool            `yaml:"dontCreateAssetIfAssetNotAvailable"`
	// configuration related to group
	Type         string                `yaml:"type"`
	FieldMapping []ContentFieldMapping `yaml:"fieldMapping"`
	// configuration related to category part
	LinkToParents bool `yaml:"linkToParents"`
	// configuration related to reference
	RefContentTypeMapping ContentTypeMapping `yaml:"refContentTypeMapping"`
	AlwaysNew             bool               `yaml:"alwaysNew"`
	SearchTerm            string             `yaml:"searchTerm"`
	SearchOnLibrary       bool               `yaml:"searchOnLibrary"`
	SearchKeys            []string           `yaml:"searchKeys"`
	SearchOnDeliveryAPI   bool               `yaml:"searchOnDeliveryAPI"`
	// configuration related the column value in
	ValueAsJSON bool   `yaml:"valueAsJSON"`
	JSONKey     string `yaml:"JSONKey"`
	SearchType  string `yaml:"searchType"`
	// if json is list for multi groups or list of references
	JSONListIndex   int
	ValueAsJSONList bool
	// link related data

}

type LinkMapping struct {
	LinkUrl         string `yaml:"linkUrl"`
	LinkText        string `yaml:"linkText"`
	LinkDescription string `yaml:"linkDescription"`
}

const JOIN_VALUE_VAR_REGX string = "\\${1}\\{{1}\\w+\\}{1}"
const JOIN_VALUE_VAR_SYMBOL_REGX string = "\\$*\\{*\\}*"

type RefPropertyMapping struct {
	PropertyName        string `yaml:"propertyName"`
	ContentFieldMapping `yaml:",inline"`
}

func (contentFieldMapping ContentFieldMapping) ToJSONListIndexableContentFieldMapping(jsonValueListIndex int) ContentFieldMapping {
	clonedContentFieldMapping := contentFieldMapping.Clone()
	clonedContentFieldMapping.JSONListIndex = jsonValueListIndex
	clonedContentFieldMapping.ValueAsJSONList = true
	return clonedContentFieldMapping
}

func (contentFieldMapping ContentFieldMapping) Clone() ContentFieldMapping {
	clonedContentFieldMapping := ContentFieldMapping{}
	copier.Copy(&clonedContentFieldMapping, &contentFieldMapping)
	return clonedContentFieldMapping
}

func (contentFieldMapping ContentFieldMapping) ConvertToGenericData(dataRow DataRow, configTypeMapping *ContentTypeMapping) (api.GenericData, error) {
	data := api.GenericData{}
	data.Name = contentFieldMapping.AcousticProperty
	data.Type = contentFieldMapping.PropertyType
	data.Ignore = contentFieldMapping.Ignore
	err := contentFieldMapping.Validate()
	if err != nil {
		return api.GenericData{}, errors.ErrorWithStack(err)
	}
	val, err := contentFieldMapping.Value(dataRow, configTypeMapping)
	if err != nil {
		return api.GenericData{}, errors.ErrorWithStack(err)
	}
	if val == nil && contentFieldMapping.Mandatory {
		return api.GenericData{}, errors.ErrorMessageWithStack("empty value for mandatory field : " + contentFieldMapping.CsvProperty)
	}
	data.Value = val
	context, err := contentFieldMapping.Context(dataRow, configTypeMapping)
	if err != nil {
		return api.GenericData{}, errors.ErrorWithStack(err)
	}
	data.Context = context
	return data, nil
}

func (contentFieldMapping ContentFieldMapping) getCsvValueOrStaticValue(dataRow DataRow) (string, error) {
	if contentFieldMapping.JoinedValue != "" {
		variableRegx, _ := regexp.Compile(JOIN_VALUE_VAR_REGX)
		variableSymbolRegx, _ := regexp.Compile(JOIN_VALUE_VAR_SYMBOL_REGX)
		matchedVariables := variableRegx.FindAllString(contentFieldMapping.JoinedValue, -1)
		transformedValue := contentFieldMapping.JoinedValue
		for _, matchedVariable := range matchedVariables {
			matchedVariableName := variableSymbolRegx.ReplaceAllString(matchedVariable, "")
			variableValue, err := dataRow.Get(matchedVariableName)
			if err != nil {
				return "", err
			}
			transformedValue = strings.Replace(transformedValue, matchedVariable, variableValue, 1)
		}
		return transformedValue, nil
	} else if contentFieldMapping.StaticValue != "" {
		return contentFieldMapping.StaticValue, nil
	} else if contentFieldMapping.Regx != nil {
		value, err := dataRow.Get(contentFieldMapping.CsvProperty)
		if err != nil {
			return "", err
		}
		var extractedValue = value
		for _, regx := range contentFieldMapping.Regx {
			compiledRegx, err := regexp.Compile(regx)
			if err != nil {
				return "", err
			}
			extractedValue = compiledRegx.FindString(extractedValue)
		}

		if extractedValue == "" {
			return "", errors.ErrorMessageWithStack("value was not extracted using regx , check  the regx is correct , check the regx configs. value :" + value)
		}
		return extractedValue, nil
	} else {
		value, err := dataRow.Get(contentFieldMapping.CsvProperty)
		if err != nil {
			return "", err
		}

		if contentFieldMapping.ValueAsJSONList {
			var jsonValue []map[string]interface{}
			json.Unmarshal([]byte(value), &jsonValue)
			mapWithStringVal, err := toStringValueMapArray(jsonValue)
			if err != nil {
				return "", errors.ErrorWithStack(err)
			}
			return mapWithStringVal[contentFieldMapping.JSONListIndex][contentFieldMapping.JSONKey], nil
		}

		if contentFieldMapping.ValueAsJSON {
			var jsonValue map[string]interface{}
			json.Unmarshal([]byte(value), &jsonValue)
			mapWithStringVal, err := toStringValueMap(jsonValue)
			if err != nil {
				return "", errors.ErrorWithStack(err)
			}
			return mapWithStringVal[contentFieldMapping.JSONKey], nil
		} else {
			return value, nil
		}
	}
}

func assetName(refPropertyMappings []RefPropertyMapping, dataRow DataRow) (map[string]string, error) {
	acc := make(map[string]string, 0)
	for _, refPropertyMapping := range refPropertyMappings {
		val, err := refPropertyMapping.Context(dataRow)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		for k, v := range val {
			acc[k] = v
		}
	}
	return acc, nil
}

func (contentFieldMapping ContentFieldMapping) getJSONACSVColumnValue(dataRow DataRow) ([]map[string]string, error) {
	value, err := dataRow.Get(contentFieldMapping.CsvProperty)
	if err != nil {
		return nil, err
	}
	var jsonArray []map[string]interface{}
	json.Unmarshal([]byte(value), &jsonArray)
	jsonValueAsStringArray, err := toStringValueMapArray(jsonArray)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	return jsonValueAsStringArray, nil
}

func toStringValueMapArray(jsonValueList []map[string]interface{}) ([]map[string]string, error) {
	stringValueMapArray := make([]map[string]string, 0)
	for _, jsonValue := range jsonValueList {
		stringValueMap, err := toStringValueMap(jsonValue)
		if err != nil {
			return nil, err
		}
		stringValueMapArray = append(stringValueMapArray, stringValueMap)
	}
	return stringValueMapArray, nil
}

func toStringValueMap(jsonValue map[string]interface{}) (map[string]string, error) {
	jsonValueAsString := make(map[string]string)
	for key, value := range jsonValue {
		valueAsString := ""
		switch reflect.TypeOf(value).Kind() {
		case reflect.Int:
			valueAsString = strconv.Itoa(value.(int))
		case reflect.String:
			valueAsString = value.(string)
		case reflect.Slice:
			slice := reflect.ValueOf(value)
			stringSlice := make([]string, 0, slice.Len())
			for i := 0; i < slice.Len(); i++ {
				stringSlice = append(stringSlice, fmt.Sprint(slice.Index(i)))
			}
			valueAsString = "[" + strings.Join(stringSlice, ",") + "]"
		default:
			return nil, errors.ErrorMessageWithStack("type not available. ")
		}
		jsonValueAsString[key] = valueAsString
	}
	return jsonValueAsString, nil
}

func (contentFieldMapping ContentFieldMapping) Validate() error {
	switch propType := api.FieldType(contentFieldMapping.PropertyType); propType {
	case api.MultiGroup:
		if !contentFieldMapping.ValueAsJSON || contentFieldMapping.CsvProperty == "" {
			return errors.ErrorMessageWithStack(string(api.MultiGroup + " should use in one single column as a array of json"))
		}
		fieldMappings := contentFieldMapping.FieldMapping
		if len(fieldMappings) == 0 {
			return errors.ErrorMessageWithStack(string(api.MultiGroup + " should have field mappings"))
		}

		for _, fieldMapping := range fieldMappings {
			if fieldMapping.JSONKey == "" {
				return errors.ErrorMessageWithStack(string(api.MultiGroup + " should have attached json key in each field mappings"))
			}
		}
	}
	return nil
}

func (contentFieldMapping ContentFieldMapping) Value(dataRow DataRow, configTypeMapping *ContentTypeMapping) (interface{}, error) {
	switch propType := api.FieldType(contentFieldMapping.PropertyType); propType {
	case api.Category, api.CategoryPart:
		category := api.AcousticCategory{}
		value, err := contentFieldMapping.getCsvValueOrStaticValue(dataRow)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		if value == "" {
			return nil, nil
		}
		categoryItems := strings.Split(value, env.MultipleItemsSeperator())
		catsWithRootCat := make([]string, 0, len(categoryItems))
		for _, categoryItem := range categoryItems {
			catsWithRootCat = append(catsWithRootCat, contentFieldMapping.CategoryName+env.CategoryHierarchySeperator()+categoryItem)
		}
		category.Value = strings.Join(catsWithRootCat, env.MultipleItemsSeperator())
		category.Operation = contentFieldMapping.Operation
		return category, nil
	case api.Group:
		group := api.AcousticGroup{}
		group.Type = contentFieldMapping.Type
		dataList := make([]api.GenericData, 0, len(contentFieldMapping.FieldMapping))
		for _, fieldMapping := range contentFieldMapping.FieldMapping {
			data, err := fieldMapping.ConvertToGenericData(dataRow, configTypeMapping)
			if err != nil {
				return nil, errors.ErrorWithStack(err)
			}
			dataList = append(dataList, data)
		}
		group.Data = dataList
		return group, nil
	case api.MultiGroup:
		multiGroup := api.AcousticMultiGroup{}
		multiGroup.Type = contentFieldMapping.Type
		valueAsJson, err := contentFieldMapping.getJSONACSVColumnValue(dataRow)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		if valueAsJson == nil {
			return nil, nil
		}
		group_data_list := make([][]api.GenericData, 0, len(valueAsJson))
		for json_index, _ := range valueAsJson {
			dataList := make([]api.GenericData, 0, len(contentFieldMapping.FieldMapping))
			for _, fieldMapping := range contentFieldMapping.FieldMapping {
				jSONListIndexableField := fieldMapping.ToJSONListIndexableContentFieldMapping(json_index)
				jSONListIndexableField.ValueAsJSON = contentFieldMapping.ValueAsJSON
				jSONListIndexableField.CsvProperty = contentFieldMapping.CsvProperty
				data, err := jSONListIndexableField.ConvertToGenericData(dataRow, configTypeMapping)
				if err != nil {
					return nil, errors.ErrorWithStack(err)
				}
				dataList = append(dataList, data)
			}
			group_data_list = append(group_data_list, dataList)
		}

		multiGroup.Data = group_data_list
		return multiGroup, nil
	case api.File:
		value, err := contentFieldMapping.getCsvValueOrStaticValue(dataRow)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		if value == "" {
			return nil, nil
		}
		assetNameMappings, err := assetName(contentFieldMapping.AssetName.AssetName, dataRow)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		asset := api.AcousticFileAsset{
			AssetNameConfig: api.AssetNameConfig{
				UseOnlyAssetName:        contentFieldMapping.AssetName.UseOnlyAssetName,
				AppendOriginalAssetName: contentFieldMapping.AssetName.AppendOriginalAssetName,
				AssetName:               assetNameMappings,
			},
			AcousticAssetBasePath:              contentFieldMapping.AcousticAssetBasePath,
			AssetLocation:                      contentFieldMapping.AssetLocation,
			Tags:                               configTypeMapping.Tags,
			IsWebUrl:                           contentFieldMapping.IsWebUrl,
			DontCreateAssetIfAssetNotAvailable: contentFieldMapping.DontCreateAssetIfAssetNotAvailable,
			Value:                              value,
		}
		return asset, nil
	case api.Image:
		value, err := contentFieldMapping.getCsvValueOrStaticValue(dataRow)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		if value == "" {
			return nil, nil
		}
		assetNameMappings, err := assetName(contentFieldMapping.AssetName.AssetName, dataRow)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		image := api.AcousticImageAsset{
			Profiles:              contentFieldMapping.Profiles,
			EnforceImageDimension: contentFieldMapping.EnforceImageDimension,
			ImageHeight:           contentFieldMapping.ImageHeight,
			ImageWidth:            contentFieldMapping.ImageWidth,
		}
		image.AssetNameConfig = api.AssetNameConfig{
			UseOnlyAssetName:        contentFieldMapping.AssetName.UseOnlyAssetName,
			AppendOriginalAssetName: contentFieldMapping.AssetName.AppendOriginalAssetName,
			AssetName:               assetNameMappings,
		}
		image.DontCreateAssetIfAssetNotAvailable = contentFieldMapping.DontCreateAssetIfAssetNotAvailable
		image.AcousticAssetBasePath = contentFieldMapping.AcousticAssetBasePath
		image.AssetLocation = contentFieldMapping.AssetLocation
		image.Tags = append(contentFieldMapping.RefContentTypeMapping.Tags, configTypeMapping.Tags...)
		image.UseExistingAsset = contentFieldMapping.UseExistingAsset
		image.IsWebUrl = contentFieldMapping.IsWebUrl
		image.Value = value
		return image, nil
	case api.MultiImage:
		value, err := contentFieldMapping.getCsvValueOrStaticValue(dataRow)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		if value == "" {
			return nil, nil
		}
		assetNameMappings, err := assetName(contentFieldMapping.AssetName.AssetName, dataRow)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		imageAssets := strings.Split(value, env.MultipleItemsSeperator())
		convertedImageAssets := funk.Map(imageAssets, func(imageAsset string) api.AcousticImageAsset {
			image := api.AcousticImageAsset{
				Profiles:              contentFieldMapping.Profiles,
				EnforceImageDimension: contentFieldMapping.EnforceImageDimension,
				ImageHeight:           contentFieldMapping.ImageHeight,
				ImageWidth:            contentFieldMapping.ImageWidth,
			}
			image.AssetNameConfig = api.AssetNameConfig{
				UseOnlyAssetName:        contentFieldMapping.AssetName.UseOnlyAssetName,
				AppendOriginalAssetName: contentFieldMapping.AssetName.AppendOriginalAssetName,
				AssetName:               assetNameMappings,
			}
			image.DontCreateAssetIfAssetNotAvailable = contentFieldMapping.DontCreateAssetIfAssetNotAvailable
			image.AcousticAssetBasePath = contentFieldMapping.AcousticAssetBasePath
			image.AssetLocation = contentFieldMapping.AssetLocation
			image.Tags = append(contentFieldMapping.RefContentTypeMapping.Tags, configTypeMapping.Tags...)
			image.IsWebUrl = contentFieldMapping.IsWebUrl
			image.UseExistingAsset = contentFieldMapping.UseExistingAsset
			image.Value = imageAsset
			return image
		}).([]api.AcousticImageAsset)
		return api.AcousticMultiImageAsset{
			Assets: convertedImageAssets,
		}, nil
	case api.Reference:
		value, err := contentFieldMapping.getCsvValueOrStaticValue(dataRow)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		if value == "" {
			return nil, errors.ErrorWithStack(err)
		}
		reference, err := createAcousticReference(value, dataRow, configTypeMapping, contentFieldMapping)
		if err != nil {
			return nil, err
		}
		return reference, nil
	case api.MultiReference:
		value, err := contentFieldMapping.getCsvValueOrStaticValue(dataRow)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		if value == "" {
			return nil, errors.ErrorWithStack(err)
		}
		referenceKeys := strings.Split(value, env.MultipleItemsSeperator())
		references := make([]api.AcousticReference, 0)
		for _, referenceKey := range referenceKeys {
			reference, err := createAcousticReference(referenceKey, dataRow, configTypeMapping, contentFieldMapping)
			if err != nil {
				return nil, err
			}
			references = append(references, reference)
		}
		return api.AcousticMultiReference{
			References: references,
			Operation:  contentFieldMapping.Operation,
		}, nil
	case api.Text:
		value, err := contentFieldMapping.getCsvValueOrStaticValue(dataRow)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		if value == "" {
			return nil, nil
		}
		if contentFieldMapping.SanitizeConfig.Sanitize {
			for _, regx := range contentFieldMapping.SanitizeConfig.Regx {
				compiledRegx, err := regexp.Compile(regx)
				if err != nil {
					return nil, err
				}
				value = compiledRegx.ReplaceAllString(value, "")
			}
		}
		return api.AcousticValue{
			Value:        value,
			LoadFromFile: contentFieldMapping.LoadFromFile,
		}, nil
	default:
		value, err := contentFieldMapping.getCsvValueOrStaticValue(dataRow)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		if value == "" {
			return nil, nil
		}
		return value, nil
	}
}

func createAcousticReference(value string, dataRow DataRow, configTypeMapping *ContentTypeMapping, contentFieldMapping ContentFieldMapping) (api.AcousticReference, error) {
	reference := api.AcousticReference{}
	reference.Type = contentFieldMapping.RefContentTypeMapping.Type
	reference.AlwaysNew = contentFieldMapping.AlwaysNew
	reference.Operation = contentFieldMapping.Operation
	reference.NameFields = contentFieldMapping.RefContentTypeMapping.Name
	reference.Tags = append(contentFieldMapping.RefContentTypeMapping.Tags, configTypeMapping.Tags...)
	reference.SearchType = contentFieldMapping.SearchType
	reference.SearchOnDeliveryAPI = contentFieldMapping.SearchOnDeliveryAPI

	if !contentFieldMapping.AlwaysNew {
		reference.SearchValues = make([]string, 0)
		for _, _ = range contentFieldMapping.SearchKeys {
			reference.SearchValues = append(reference.SearchValues, value)
		}
		reference.SearchTerm = contentFieldMapping.SearchTerm
		reference.SearchOnLibrary = contentFieldMapping.SearchOnLibrary
	} else {
		dataList := make([]api.GenericData, 0, len(contentFieldMapping.RefContentTypeMapping.FieldMapping))
		for _, fieldMapping := range contentFieldMapping.FieldMapping {
			data, err := fieldMapping.ConvertToGenericData(dataRow, configTypeMapping)
			if err != nil {
				return api.AcousticReference{}, errors.ErrorWithStack(err)
			}
			dataList = append(dataList, data)
		}
		reference.Data = dataList
	}
	return reference, nil
}

func (refPropertyMapping RefPropertyMapping) Context(dataRow DataRow) (map[string]string, error) {
	value, err := refPropertyMapping.getCsvValueOrStaticValue(dataRow)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	return map[string]string{
		refPropertyMapping.PropertyName: value,
	}, nil
}

func (contentFieldMapping ContentFieldMapping) Context(dataRow DataRow, configTypeMapping *ContentTypeMapping) (api.Context, error) {
	switch fieldType := api.FieldType(contentFieldMapping.PropertyType); fieldType {
	case api.CategoryPart:
		return api.Context{Data: map[api.ContextKey]interface{}{
			api.LinkToParents: contentFieldMapping.LinkToParents,
		}}, nil
	default:
		return api.Context{}, nil
	}
}

func (csvContentTypesMapping *ContentTypesMapping) GetContentTypeMapping(contentType string) (*ContentTypeMapping, error) {
	contentTypeMapping := koazee.StreamOf(csvContentTypesMapping.ContentType).
		Filter(func(contentTypeMapping ContentTypeMapping) bool {
			return contentTypeMapping.Type == contentType
		}).
		First().Val().(ContentTypeMapping)

	if &contentTypeMapping != nil {
		return &contentTypeMapping, nil
	} else {
		return nil, errors.ErrorMessageWithStack("No mapping found for content type :" + contentType)
	}
}

func (contentFieldMapping ContentFieldMapping) GetAllChildCSVFields() ([]string, error) {
	childFields := make([]string, 0)
	if contentFieldMapping.FieldMapping != nil {
		for _, childFieldMapping := range contentFieldMapping.FieldMapping {
			childFields = append(childFields, childFieldMapping.CsvProperty)
			childFieldsOfChildField, err := childFieldMapping.GetAllChildCSVFields()
			if err != nil {
				return nil, errors.ErrorWithStack(err)
			}
			childFields = append(childFields, childFieldsOfChildField...)
		}

	}
	return childFields, nil
}

func (contentFieldMapping ContentFieldMapping) GetAcousticChildFields() (map[string]interface{}, error) {
	if contentFieldMapping.FieldMapping != nil {
		childFields := make(map[string]interface{}, 0)
		for _, childFieldMapping := range contentFieldMapping.FieldMapping {
			childFieldsOfChildField, err := childFieldMapping.GetAcousticChildFields()
			if err != nil {
				return nil, errors.ErrorWithStack(err)
			}
			childFields[childFieldMapping.AcousticProperty] = childFieldsOfChildField
		}
		return childFields, nil
	} else {
		return nil, nil
	}
}

func (csvContentTypeMapping *ContentTypeMapping) GetFieldMappingByAcousticField(acousticField string) (*ContentFieldMapping, error) {
	fieldMapping := koazee.StreamOf(csvContentTypeMapping.FieldMapping).
		Filter(func(contentFieldMapping ContentFieldMapping) bool {
			return contentFieldMapping.AcousticProperty == acousticField
		}).
		First().Val().(ContentFieldMapping)

	if &fieldMapping != nil {
		return &fieldMapping, nil
	} else {
		return nil, errors.ErrorMessageWithStack("No mapping found for field :" + acousticField)
	}
}

func (csvContentTypeMapping *ContentTypeMapping) GetAcousticFields() []string {
	return koazee.StreamOf(csvContentTypeMapping.FieldMapping).
		Map(func(contentFieldMapping ContentFieldMapping) string {
			return contentFieldMapping.AcousticProperty
		}).
		Out().Val().([]string)
}

type CSVToAcousticFieldMapping struct {
	CSVField           string
	AcousticField      string
	ChildFieldMappings []CSVToAcousticFieldMapping
}

type AcousticField struct {
	Name  string
	Child *AcousticField
}

func (acousticField *AcousticField) GetFieldNameHierarchy() []string {
	fieldNameHierarchy := make([]string, 0)
	child := acousticField.Child
	for child != nil {
		fieldNameHierarchy = append(fieldNameHierarchy, child.Name)
		child = child.Child
	}
	return fieldNameHierarchy
}

func (csvToAcousticFieldMapping *CSVToAcousticFieldMapping) HasChildren() bool {
	return csvToAcousticFieldMapping.CSVField == "" && len(csvToAcousticFieldMapping.ChildFieldMappings) > 0
}

func (csvToAcousticFieldMapping *CSVToAcousticFieldMapping) AllChildCSVFields() []string {
	allChildCSVFields := make([]string, 0)
	if csvToAcousticFieldMapping.HasChildren() {
		for _, childFieldMapping := range csvToAcousticFieldMapping.ChildFieldMappings {
			if childFieldMapping.CSVField != "" {
				allChildCSVFields = append(allChildCSVFields, childFieldMapping.CSVField)
			}
			allChildCSVFields = append(allChildCSVFields, childFieldMapping.AllChildCSVFields()...)
		}
	}
	return allChildCSVFields
}

func GetAcousticField(csvToAcousticFieldMappings []CSVToAcousticFieldMapping, csvField string) (*AcousticField, error) {
	for _, mapping := range csvToAcousticFieldMappings {
		mappedAcousticField := mapping.getAcousticField(csvField)
		if mappedAcousticField != nil {
			return mappedAcousticField, nil
		}
	}
	return nil, errors.ErrorMessageWithStack("No matching acoustic field found for csv field : " + csvField)
}

func (csvToAcousticFieldMapping *CSVToAcousticFieldMapping) getAcousticField(csvField string) *AcousticField {
	if csvToAcousticFieldMapping.CSVField == csvField {
		return &AcousticField{
			Name: csvToAcousticFieldMapping.AcousticField,
		}
	} else if csvToAcousticFieldMapping.HasChildren() {
		for _, childCsvToAcousticFieldMapping := range csvToAcousticFieldMapping.ChildFieldMappings {
			childMappedAcousticField := childCsvToAcousticFieldMapping.getAcousticField(csvField)
			if childMappedAcousticField != nil {
				return &AcousticField{
					Name:  csvToAcousticFieldMapping.AcousticField,
					Child: childMappedAcousticField,
				}
			}
		}
	}
	return nil
}

func (fieldMapping *ContentFieldMapping) GetChildCSVToAcousticFieldMapping() []CSVToAcousticFieldMapping {
	csvToAcousticFieldMappings := make([]CSVToAcousticFieldMapping, 0)
	for _, childFieldMapping := range fieldMapping.FieldMapping {
		csvToAcousticFieldMappings = append(csvToAcousticFieldMappings, CSVToAcousticFieldMapping{
			AcousticField:      childFieldMapping.AcousticProperty,
			CSVField:           childFieldMapping.CsvProperty,
			ChildFieldMappings: childFieldMapping.GetChildCSVToAcousticFieldMapping(),
		})
	}
	if len(csvToAcousticFieldMappings) > 0 {
		return csvToAcousticFieldMappings
	} else {
		return nil
	}
}

func (csvContentTypeMapping *ContentTypeMapping) GetCSVToAcousticFieldMapping() []CSVToAcousticFieldMapping {
	csvToAcousticFieldMappings := make([]CSVToAcousticFieldMapping, 0)
	for _, fieldMapping := range csvContentTypeMapping.FieldMapping {
		csvToAcousticFieldMappings = append(csvToAcousticFieldMappings, CSVToAcousticFieldMapping{
			AcousticField:      fieldMapping.AcousticProperty,
			CSVField:           fieldMapping.CsvProperty,
			ChildFieldMappings: fieldMapping.GetChildCSVToAcousticFieldMapping(),
		})
	}
	return csvToAcousticFieldMappings
}

type Config interface {
	GetContentType(contentModel string) (*ContentTypeMapping, error)
	GetCategory(categoryName string) (*CategoryMapping, error)
	GetDeleteMapping(name string) (*DeleteMapping, error)
	GetSiteMapping(pageContentModel string) (*SiteMapping, error)
}

type config struct {
	mappings *ContentTypesMapping
}

func InitContentTypeMappingConfig(configPath string) (Config, error) {
	if configContent, err := ioutil.ReadFile(configPath); err != nil {
		return nil, errors.ErrorWithStack(err)
	} else {
		mappings := &ContentTypesMapping{}
		if err := yaml.Unmarshal(configContent, mappings); err != nil {
			return nil, errors.ErrorWithStack(err)
		} else {
			log.Info("csv config loaded config path : " + configPath)
			return &config{
				mappings: mappings,
			}, nil
		}
	}
}

func (config config) GetCategory(categoryName string) (*CategoryMapping, error) {
	categoryMapping := koazee.StreamOf(config.mappings.CategoryMapping).
		Filter(func(mapping CategoryMapping) bool {
			return mapping.Parent == categoryName
		}).
		First().Val().(CategoryMapping)
	if &categoryMapping == nil {
		return nil, errors.ErrorMessageWithStack("No category mapping found for provided category :" + categoryName)
	}
	return &categoryMapping, nil
}

func (config *config) GetContentType(contentType string) (*ContentTypeMapping, error) {
	if config.mappings != nil {
		if contentTypeMapping, err := config.mappings.GetContentTypeMapping(contentType); err != nil {
			return nil, errors.ErrorWithStack(err)
		} else {
			return contentTypeMapping, nil
		}
	} else {
		return nil, errors.ErrorMessageWithStack("config not yet set")
	}
}

func (config *config) GetDeleteMapping(name string) (*DeleteMapping, error) {
	deleteMapping := koazee.StreamOf(config.mappings.DeleteMapping).
		Filter(func(deleteMapping DeleteMapping) bool {
			return deleteMapping.Name == name
		}).
		First().Val().(DeleteMapping)
	if &deleteMapping == nil {
		return nil, errors.ErrorMessageWithStack("No delete mapping found for provided name :" + name)
	}
	return &deleteMapping, nil
}

func (config config) GetSiteMapping(pageContentType string) (*SiteMapping, error) {
	siteMapping := koazee.StreamOf(config.mappings.SiteMapping).
		Filter(func(siteMapping SiteMapping) bool {
			return siteMapping.Type == pageContentType
		}).
		First().Val().(SiteMapping)
	if &siteMapping == nil {
		return nil, errors.ErrorMessageWithStack("No site mapping found for provided page content type :" + pageContentType)
	}
	return &siteMapping, nil
}
