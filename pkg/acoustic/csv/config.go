package csv

import (
	"errors"
	"github.com/goccy/go-yaml"
	log "github.com/sirupsen/logrus"
	"github.com/wesovilabs/koazee"
	"io/ioutil"
)


type ContentType interface {
	GetFieldMapping(csvField string) (ContentFieldMapping, error)
}

type ContentFieldMapping struct {
	CsvProperty string `yaml:"csvProperty"`
	AcousticProperty string `yaml:"acousticProperty"`
	PropertyType  string `yaml:"propertyType"`
}

type ContentTypeMapping struct {
	Type string                        `yaml:"type"`
	FieldMapping []ContentFieldMapping `yaml:"fieldMapping"`
	Name []string `yaml:"name"`
}

type ContentTypesMapping struct {
	ContentType []ContentTypeMapping `yaml:"contentType"`
}

func (csvContentTypesMapping *ContentTypesMapping) GetContentTypeMapping(contentType string) (*ContentTypeMapping, error) {
	contentTypeMapping := koazee.StreamOf(csvContentTypesMapping.ContentType).
		Filter(func(contentTypeMapping ContentTypeMapping) bool{
			return contentTypeMapping.Type == contentType
		}).
		First().Val().(ContentTypeMapping)

	if &contentTypeMapping != nil {
		return &contentTypeMapping,nil
	} else {
		return nil, errors.New("No mapping found for content type :" + contentType)
	}
}

func (csvContentTypeMapping *ContentTypeMapping) GetFieldMapping(csvField string) (*ContentFieldMapping, error) {
	  fieldMapping := koazee.StreamOf(csvContentTypeMapping.FieldMapping).
		Filter(func(contentFieldMapping ContentFieldMapping) bool{
			return contentFieldMapping.CsvProperty == csvField
		}).
		First().Val().(ContentFieldMapping)

	  if &fieldMapping != nil {
	  	return &fieldMapping,nil
	  } else {
	  	return nil, errors.New("No mapping found for field :" + csvField)
	  }
}

type Config interface {
	Get(contentModel string) (*ContentTypeMapping,error)
}

type config struct {
	mappings *ContentTypesMapping
}

func InitConfig(configPath string) (Config,error) {
	if configContent, err := ioutil.ReadFile(configPath) ; err != nil {
		return nil,err
	} else {
		mappings := &ContentTypesMapping{}
		if err := yaml.Unmarshal(configContent, mappings); err != nil {
			return nil,err
		} else {
			log.Info("csv config loaded config path : " + configPath )
			return &config{
				mappings: mappings,
			},nil
		}
	}
}

func (config *config) Get(contentType string) (*ContentTypeMapping,error) {
	if config.mappings != nil {
		if contentTypeMapping,err := config.mappings.GetContentTypeMapping(contentType) ; err != nil {
			return nil,err
		} else {
			return contentTypeMapping,nil
		}
	} else {
		return nil,errors.New("config not yet set")
	}
}
