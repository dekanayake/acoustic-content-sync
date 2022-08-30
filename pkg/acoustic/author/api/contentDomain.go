package api

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
)

type Operation string

const (
	DELETE            Operation = "delete"
	UPDATE            Operation = "update"
	CREATE            Operation = "create"
	DEFAULT_OPERATION Operation = "-"
)

type AssetType string

const (
	DOCUMENT AssetType = "document"
	FILE     AssetType = "file"
	IMAGE    AssetType = "image"
	VIDEO    AssetType = "video"
)

type FieldType string

const (
	Text           FieldType = "text"
	MultiText      FieldType = "multi-text"
	FormattedText  FieldType = "formatted-text"
	Number         FieldType = "number"
	Float          FieldType = "float"
	Boolean        FieldType = "toggle"
	Link           FieldType = "link"
	Date           FieldType = "date"
	Category       FieldType = "category"
	CategoryPart   FieldType = "category-part"
	File           FieldType = "file"
	Video          FieldType = "video"
	Image          FieldType = "image"
	Group          FieldType = "group"
	MultiGroup     FieldType = "multi-group"
	Reference      FieldType = "reference"
	MultiReference FieldType = "multi-reference"
)

func (ft FieldType) Convert() (AcousticFieldType, error) {
	switch ft {
	case Text, MultiText:
		return AcousticFieldType(AcousticFieldText), nil
	case FormattedText:
		return AcousticFieldType(AcousticFieldFormattedText), nil
	case Number:
		return AcousticFieldType(AcousticFieldNumber), nil
	case Float:
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
	case Group, MultiGroup:
		return AcousticFieldType(AcousticFieldGroup), nil
	case Reference, MultiReference:
		return AcousticFieldType(AcousticFieldReference), nil
	default:
		return AcousticFieldType("no mapping"), errors.ErrorMessageWithStack("No Acoustic field type found for property type" + string(ft))
	}
}

type AcousticFieldType string

const (
	AcousticFieldText          FieldType = "text"
	AcousticFieldFormattedText FieldType = "formattedtext"
	AcousticFieldNumber        FieldType = "number"
	AcousticFieldBoolean       FieldType = "toggle"
	AcousticFieldLink          FieldType = "link"
	AcousticFieldDate          FieldType = "date"
	AcousticFieldCategory      FieldType = "category"
	AcousticFieldFile          FieldType = "file"
	AcousticFieldVideo         FieldType = "video"
	AcousticFieldImage         FieldType = "image"
	AcousticFieldGroup         FieldType = "group"
	AcousticFieldReference     FieldType = "reference"
)

type Tags struct {
	Values []string `json:"values"`
}

type Content struct {
	ID        string                 `json:"id,omitempty"`
	REV       string                 `json:"rev,omitempty"`
	Name      string                 `json:"name"`
	TypeId    string                 `json:"typeId"`
	Status    string                 `json:"status"`
	Elements  map[string]interface{} `json:"elements"`
	LibraryID string                 `json:"libraryId"`
	Tags      []string               `json:"tags"`
}

type PreContentCreateFunc func() (Element, error)
type PreContentUpdateFunc func(updatedElement Element) (Element, []PostContentUpdateFunc, error)
type PostContentUpdateFunc func() error

type CSVValues struct {
	Name        string
	Value       string
	ChildValues map[string]CSVValues
}

func (csvValues CSVValues) hasChildren() bool {
	return len(csvValues.ChildValues) > 0
}

func (csvValues CSVValues) GetValue(fieldName string) (string, error) {
	if !csvValues.hasChildren() {
		if csvValues.Value != "" {
			return csvValues.Value, nil
		}
	} else {
		for _, childValue := range csvValues.ChildValues {
			val, _ := childValue.GetValue(fieldName)
			if val != "" {
				return val, nil
			}
		}
	}
	return "", errors.ErrorMessageWithStack("no value for the field name :" + fieldName)
}

type Element interface {
	Convert(data interface{}) (Element, error)
	Update(new Element) (Element, error)
	PreContentCreateFunctions() []PreContentCreateFunc
	PreContentUpdateFunctions() []PreContentUpdateFunc
	ChildElements() map[string]Element
	UpdateChildElement(key string, updatedElement Element) error
	ToCSV(childFields map[string]interface{}) (CSVValues, error)
	GetOperation() Operation
}

type element struct {
	ElementType                  AcousticFieldType      `json:"elementType"`
	PreContentCreateFunctionList []PreContentCreateFunc `json:"-"`
	PreContentUpdateFunctionList []PreContentUpdateFunc `json:"-"`
	Operation                    Operation              `json:"-"`
}

type TextElement struct {
	Value string `json:"value"`
	element
}

func (element element) PreContentCreateFunctions() []PreContentCreateFunc {
	if element.PreContentCreateFunctionList == nil {
		return []PreContentCreateFunc{}
	} else {
		return element.PreContentCreateFunctionList
	}
}

func (element element) PreContentUpdateFunctions() []PreContentUpdateFunc {
	if element.PreContentCreateFunctionList == nil {
		return []PreContentUpdateFunc{}
	} else {
		return element.PreContentUpdateFunctionList
	}
}

func (element element) GetOperation() Operation {
	return element.Operation
}

func (element element) ChildElements() map[string]Element {
	return nil
}

func (element element) UpdateChildElement(key string, updatedElement Element) error {
	return nil
}

type FormattedTextElement struct {
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

type FloatElement struct {
	Value float64 `json:"value"`
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
	Categories  []string `json:"categories"`
	element
}

type CategoryPartElement struct {
	CategoryIds []string `json:"categoryIds"`
	element
}

type ImageElement struct {
	Mode  string `json:"mode"`
	Asset Asset  `json:"asset"`
	URL   string `json:"url"`
	element
}

type FileElement struct {
	Asset Asset `json:"asset"`
	element
}

type GroupElement struct {
	TypeRef map[string]string      `json:"typeRef"`
	Value   map[string]interface{} `json:"value"`
	element
}

func (groupElement GroupElement) ChildElements() map[string]Element {
	elementMap := make(map[string]Element)
	for key, value := range groupElement.Value {
		elementMap[key] = value.(Element)
	}
	return elementMap
}

func (groupElement GroupElement) UpdateChildElement(key string, updatedElement Element) error {
	if _, ok := groupElement.Value[key]; ok {
		groupElement.Value[key] = updatedElement
		return nil
	} else {
		return errors.ErrorMessageWithStack("key does not exist :" + key)
	}

}

type MultiGroupElement struct {
	TypeRef map[string]string        `json:"typeRef"`
	Values  []map[string]interface{} `json:"values"`
	element
}

type ReferenceElement struct {
	Value ReferenceValue `json:"value"`
	element
}

type MultiReferenceElement struct {
	Values []ReferenceValue `json:"values"`
	element
}

type ReferenceValue struct {
	ID string `json:"id"`
}

type Asset struct {
	ID string `json:"id"`
}

type ContentAutheringResponse struct {
	Id     string `json:"id"`
	Rev    string `json:"rev"`
	Name   string `json:"name"`
	TypeId string `json:"typeId"`
	Type   string `json:"type"`
}

type ContentUpdateResponse struct {
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
	case FormattedText:
		element := FormattedTextElement{}
		element.ElementType = acousticFieldType
		return element, nil
	case Number:
		element := NumberElement{}
		element.ElementType = acousticFieldType
		return element, nil
	case Float:
		element := FloatElement{}
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
	case File:
		element := FileElement{}
		element.ElementType = acousticFieldType
		return element, nil
	case Group:
		element := GroupElement{}
		element.ElementType = acousticFieldType
		return element, nil
	case MultiGroup:
		element := MultiGroupElement{}
		element.ElementType = acousticFieldType
		return element, nil
	case Reference:
		element := ReferenceElement{}
		element.ElementType = acousticFieldType
		return element, nil
	case MultiReference:
		element := MultiReferenceElement{}
		element.ElementType = acousticFieldType
		return element, nil
	default:
		return nil, errors.ErrorMessageWithStack("No element found for property type " + fieldType)
	}
}

type FeedType string

const (
	CSV      FeedType = "CSV"
	ACOUSTIC FeedType = "Acoustic"
)
