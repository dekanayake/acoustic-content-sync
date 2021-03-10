package csv

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/author/api"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/goccy/go-yaml"
	log "github.com/sirupsen/logrus"
	"github.com/wesovilabs/koazee"
	"io/ioutil"
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
	Type         string                `yaml:"type"`
	FieldMapping []ContentFieldMapping `yaml:"fieldMapping"`
	Name         []string              `yaml:"name"`
	Tags         []string              `yaml:"tags"`
	CsvRecordKey string                `yaml:"csvRecordKey"`
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
}

type RefPropertyMapping struct {
	RefCSVProperty string `yaml:"refCSVProperty"`
	PropertyName   string `yaml:"propertyName"`
}

func (contentFieldMapping ContentFieldMapping) ConvertToGenericData(dataRow DataRow, configTypeMapping *ContentTypeMapping) (api.GenericData, error) {
	data := api.GenericData{}
	data.Name = contentFieldMapping.AcousticProperty
	data.Type = contentFieldMapping.PropertyType
	data.Ignore = contentFieldMapping.Ignore
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
		return dataRow.Get(contentFieldMapping.CsvProperty)
	}
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
	case api.Image:
		assetName, err := assetName(contentFieldMapping.AssetName, dataRow)
		if err != nil {
			return api.Context{}, errors.ErrorWithStack(err)
		}
		return api.Context{Data: map[api.ContextKey]interface{}{
			api.AssetName:             assetName,
			api.Profiles:              contentFieldMapping.Profiles,
			api.AcousticAssetBasePath: contentFieldMapping.AcousticAssetBasePath,
			api.AssetLocation:         contentFieldMapping.AssetLocation,
			api.TagList:               configTypeMapping.Tags,
			api.IsWebUrl:              contentFieldMapping.IsWebUrl,
			api.EnforceImageDimension: contentFieldMapping.EnforceImageDimension,
			api.ImageWidth:            contentFieldMapping.ImageWidth,
			api.ImageHeight:           contentFieldMapping.ImageHeight,
		}}, nil
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
