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
	"github.com/wesovilabs/koazee"
	"io/ioutil"
	"reflect"
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
}

type ContentTypeMapping struct {
	Type                   string                `yaml:"type"`
	FieldMapping           []ContentFieldMapping `yaml:"fieldMapping"`
	Name                   []string              `yaml:"name"`
	Tags                   []string              `yaml:"tags"`
	CsvRecordKey           string                `yaml:"csvRecordKey"`
	Update                 bool                  `yaml:"update"`
	CreateNonExistingItems bool                  `yaml:"createNonExistingItems"`
	SearchTerm             string                `yaml:"searchTerm"`
	SearchKeys             []string              `yaml:"searchKeys"`
	SearchType             string                `yaml:"searchType"`
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

type ContentFieldMapping struct {
	CsvProperty           string               `yaml:"csvProperty"`
	Ignore                bool                 `yaml:"ignore"`
	ValuePattern          string               `yaml:"valuePattern"`
	Mandatory             bool                 `yaml:"mandatory"`
	StaticValue           string               `yaml:"staticValue"`
	AcousticProperty      string               `yaml:"acousticProperty"`
	PropertyType          string               `yaml:"propertyType"`
	CategoryName          string               `yaml:"categoryName"`
	AssetName             []RefPropertyMapping `yaml:"assetName"`
	Profiles              []string             `yaml:"profiles"`
	AcousticAssetBasePath string               `yaml:"acousticAssetBasePath"`
	AssetLocation         string               `yaml:"assetLocation"`
	IsWebUrl              bool                 `yaml:"isWebUrl"`
	ImageWidth            uint                 `yaml:"imageWidth"`
	ImageHeight           uint                 `yaml:"imageHeight"`
	EnforceImageDimension bool                 `yaml:"enforceImageDimension"`
	// configuration related to group
	Type         string                `yaml:"type"`
	FieldMapping []ContentFieldMapping `yaml:"fieldMapping"`
	// configuration related to category part
	LinkToParents bool `yaml:"linkToParents"`
	// configuration related to reference
	RefContentTypeMapping ContentTypeMapping `yaml:"refContentTypeMapping"`
	AlwaysNew             bool               `yaml:"alwaysNew"`
	SearchTerm            string             `yaml:"searchTerm"`
	SearchKeys            []string           `yaml:"searchKeys"`
	// configuration related the column value in
	ValueAsJSON bool   `yaml:"valueAsJSON"`
	JSONKey     string `yaml:"JSONKey"`
	SearchType  string `yaml:"searchType"`
	// if json is list for multi groups or list of references
	JSONListIndex   int
	ValueAsJSONList bool
}

type RefPropertyMapping struct {
	RefCSVProperty string `yaml:"refCSVProperty"`
	PropertyName   string `yaml:"propertyName"`
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
	if contentFieldMapping.StaticValue != "" {
		return contentFieldMapping.StaticValue, nil
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
		return strings.Join(catsWithRootCat, env.MultipleItemsSeperator()), nil
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
		assetName, err := assetName(contentFieldMapping.AssetName, dataRow)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		value, err := contentFieldMapping.getCsvValueOrStaticValue(dataRow)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		if value == "" {
			return nil, nil
		}
		asset := api.AcousticFileAsset{
			AssetName:             assetName,
			AcousticAssetBasePath: contentFieldMapping.AcousticAssetBasePath,
			AssetLocation:         contentFieldMapping.AssetLocation,
			Tags:                  configTypeMapping.Tags,
			IsWebUrl:              contentFieldMapping.IsWebUrl,
			Value:                 value,
		}
		return asset, nil
	case api.Image:
		assetName, err := assetName(contentFieldMapping.AssetName, dataRow)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		value, err := contentFieldMapping.getCsvValueOrStaticValue(dataRow)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		if value == "" {
			return nil, nil
		}
		image := api.AcousticImageAsset{
			Profiles:              contentFieldMapping.Profiles,
			EnforceImageDimension: contentFieldMapping.EnforceImageDimension,
			ImageHeight:           contentFieldMapping.ImageHeight,
			ImageWidth:            contentFieldMapping.ImageWidth,
		}
		image.AssetName = assetName
		image.AcousticAssetBasePath = contentFieldMapping.AcousticAssetBasePath
		image.AssetLocation = contentFieldMapping.AssetLocation
		image.Tags = append(contentFieldMapping.RefContentTypeMapping.Tags, configTypeMapping.Tags...)
		image.IsWebUrl = contentFieldMapping.IsWebUrl
		image.Value = value
		return image, nil
	case api.Reference, api.MultiReference:
		reference := api.AcousticReference{}
		reference.Type = contentFieldMapping.RefContentTypeMapping.Type
		reference.AlwaysNew = contentFieldMapping.AlwaysNew
		reference.NameFields = contentFieldMapping.RefContentTypeMapping.Name
		reference.Tags = append(contentFieldMapping.RefContentTypeMapping.Tags, configTypeMapping.Tags...)
		reference.SearchType = contentFieldMapping.SearchType

		if !contentFieldMapping.AlwaysNew {
			value, err := contentFieldMapping.getCsvValueOrStaticValue(dataRow)
			if err != nil {
				return nil, errors.ErrorWithStack(err)
			}
			if value == "" {
				return nil, nil
			}
			reference.SearchValues = make([]string, 0)
			for _, _ = range contentFieldMapping.SearchKeys {
				reference.SearchValues = append(reference.SearchValues, value)
			}
			reference.SearchTerm = contentFieldMapping.SearchTerm
		} else {
			dataList := make([]api.GenericData, 0, len(contentFieldMapping.RefContentTypeMapping.FieldMapping))
			for _, fieldMapping := range contentFieldMapping.FieldMapping {
				data, err := fieldMapping.ConvertToGenericData(dataRow, configTypeMapping)
				if err != nil {
					return nil, errors.ErrorWithStack(err)
				}
				dataList = append(dataList, data)
			}
			reference.Data = dataList
		}
		return reference, nil
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

func (refPropertyMapping RefPropertyMapping) Context(dataRow DataRow) (map[string]string, error) {
	val, err := dataRow.Get(refPropertyMapping.RefCSVProperty)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	return map[string]string{
		refPropertyMapping.PropertyName: val,
	}, nil
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

type Config interface {
	GetContentType(contentModel string) (*ContentTypeMapping, error)
	GetCategory(categoryName string) (*CategoryMapping, error)
	GetDeleteMapping(name string) (*DeleteMapping, error)
}

type config struct {
	mappings *ContentTypesMapping
}

func InitConfig(configPath string) (Config, error) {
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
