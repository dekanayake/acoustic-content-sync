package api

import "github.com/dekanayake/acoustic-content-sync/pkg/errors"

type AssetType string

const (
	DOCUMENT AssetType = "document"
	FILE     AssetType = "file"
	IMAGE    AssetType = "image"
	VIDEO    AssetType = "video"
)

type FieldType string

const (
	Text         FieldType = "text"
	MultiText    FieldType = "multi-text"
	Number       FieldType = "number"
	Boolean      FieldType = "toggle"
	Link         FieldType = "link"
	Date         FieldType = "date"
	Category     FieldType = "category"
	CategoryPart FieldType = "category-part"
	File         FieldType = "file"
	Video        FieldType = "video"
	Image        FieldType = "image"
	Group        FieldType = "group"
)

func (ft FieldType) Convert() (AcousticFieldType, error) {
	switch ft {
	case Text, MultiText:
		return AcousticFieldType(AcousticFieldText), nil
	case Number:
		return AcousticFieldType(AcousticFieldNumber), nil
	case Boolean:
		return AcousticFieldType(AcousticFieldBoolean), nil
	case Link:
		return AcousticFieldType(AcousticFieldLink), nil
	case Date:
		return AcousticFieldType(AcousticFieldDate), nil
	case Category, CategoryPart:
		return AcousticFieldType(AcousticFieldCategory), nil
	case File:
		return AcousticFieldType(AcousticFieldFile), nil
	case Video:
		return AcousticFieldType(AcousticFieldVideo), nil
	case Image:
		return AcousticFieldType(AcousticFieldImage), nil
	case Group:
		return AcousticFieldType(AcousticFieldGroup), nil
	default:
		return AcousticFieldType("no mapping"), errors.ErrorMessageWithStack("No Acoustic field type found for property type" + string(ft))
	}
}

type AcousticFieldType string

const (
	AcousticFieldText     FieldType = "text"
	AcousticFieldNumber   FieldType = "number"
	AcousticFieldBoolean  FieldType = "toggle"
	AcousticFieldLink     FieldType = "link"
	AcousticFieldDate     FieldType = "date"
	AcousticFieldCategory FieldType = "category"
	AcousticFieldFile     FieldType = "file"
	AcousticFieldVideo    FieldType = "video"
	AcousticFieldImage    FieldType = "image"
	AcousticFieldGroup    FieldType = "group"
)

type Tags struct {
	Values []string `json:"values"`
}

type Content struct {
	Name      string                 `json:"name"`
	TypeId    string                 `json:"typeId"`
	Status    string                 `json:"status"`
	Elements  map[string]interface{} `json:"elements"`
	LibraryID string                 `json:"libraryId"`
	Tags      []string               `json:"tags"`
}

type Element interface {
	Convert(data interface{}) (Element, error)
}

type element struct {
	ElementType AcousticFieldType `json:"elementType"`
}

type TextElement struct {
	Value string `json:"value"`
	element
}

type MultiTextElement struct {
	Values []string `json:"values"`
	element
}

type NumberElement struct {
	Value int64 `json:"value"`
	element
}

type BooleanElement struct {
	Value bool `json:"value"`
	element
}

type LinkElement struct {
	LinkURL   string `json:"linkURL"`
	LinkText  string `json:"linkText"`
	LinkTitle string `json:"linkTitle"`
	element
}

type DateElement struct {
	Value string `json:"value"`
	element
}

type CategoryElement struct {
	CategoryIds []string `json:"categoryIds"`
	element
}

type CategoryPartElement struct {
	CategoryIds []string `json:"categoryIds"`
	element
}

type ImageElement struct {
	Mode  string `json:"mode"`
	Asset Asset  `json:"asset"`
	element
}

type GroupElement struct {
	TypeRef map[string]string      `json:"typeRef"`
	Value   map[string]interface{} `json:"value"`
	element
}

type Asset struct {
	ID string `json:"id"`
}

type ContentCreateResponse struct {
	Id     string `json:"id"`
	Rev    string `json:"rev"`
	Name   string `json:"name"`
	TypeId string `json:"typeId"`
	Type   string `json:"type"`
}

type ContentAuthoringErrorResponse struct {
	RequestId     string                  `json:"requestId"`
	Service       string                  `json:"service"`
	RequestMethod string                  `json:"requestMethod"`
	RequestUri    string                  `json:"requestUri"`
	Type          string                  `json:"type"`
	Errors        []ContentAuthoringError `json:"errors"`
}

type ContentAuthoringError struct {
	Code        int64       `json:"code"`
	Key         string      `json:"key"`
	Message     string      `json:"message"`
	Description string      `json:"description"`
	MoreInfo    string      `json:"more_info"`
	Category    string      `json:"category"`
	Level       string      `json:"level"`
	Parameters  interface{} `json:"parameters"`
	Field       interface{} `json:"field"`
	Locale      interface{} `json:"locale"`
}

func (element element) Convert(data interface{}) (Element, error) {
	return nil, errors.ErrorMessageWithStack("Not implementd need to override in extending elements")
}

func Build(fieldType string) (Element, error) {
	fieldTypeConst := FieldType(fieldType)
	acousticFieldType, err := fieldTypeConst.Convert()
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	switch fieldTypeConst {
	case Text:
		element := TextElement{}
		element.ElementType = acousticFieldType
		return element, nil
	case MultiText:
		element := MultiTextElement{}
		element.ElementType = acousticFieldType
		return element, nil
	case Number:
		element := NumberElement{}
		element.ElementType = acousticFieldType
		return element, nil
	case Boolean:
		element := BooleanElement{}
		element.ElementType = acousticFieldType
		return element, nil
	case Link:
		element := LinkElement{}
		element.ElementType = acousticFieldType
		return element, nil
	case Date:
		element := DateElement{}
		element.ElementType = acousticFieldType
		return element, nil
	case Category:
		element := CategoryElement{}
		element.ElementType = acousticFieldType
		return element, nil
	case CategoryPart:
		element := CategoryPartElement{}
		element.ElementType = acousticFieldType
		return element, nil
	case Image:
		element := ImageElement{}
		element.ElementType = acousticFieldType
		return element, nil
	case Group:
		element := GroupElement{}
		element.ElementType = acousticFieldType
		return element, nil
	default:
		return nil, errors.ErrorMessageWithStack("No element found for property type " + fieldType)
	}
}
