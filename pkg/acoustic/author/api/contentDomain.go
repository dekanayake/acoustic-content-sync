package api

import "errors"

type FieldType string

const (
	Text     FieldType = "text"
	Number   FieldType = "number"
	Boolean  FieldType = "toggle"
	Link     FieldType = "link"
	Date     FieldType = "date"
	Category FieldType = "category"
	File     FieldType = "file"
	Video    FieldType = "video"
	Image    FieldType = "image"
)


type Tags struct {
	Values []string `json:"values"`
}

type Content struct {
	Name string `json:"name"`
	TypeId string `json:"typeId"`
	Status string `json:"status"`
	Elements  map[string] interface{} `json:"elements"`
	LibraryID string `json:"libraryId"`
}


type Element interface {
	Convert(data interface{}) (Element,error)
}

type element struct {
	ElementType FieldType `json:"elementType"`
}

type TextElement struct {
	Value string `json:"value"`
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
	LinkURL string `json:"linkURL"`
	LinkText string `json:"linkText"`
	LinkTitle string `json:"linkTitle"`
	element
}

type DateElement struct {
	Value string `json:"value"`
	element
}

type CategoryElement struct {
	CategoryIds []string `json:"categoryIds"`
	Categories []string `json:"categories"`
	element
}

type ImageElement struct {
	Mode string `json:"mode"`
	Asset Asset `json:"asset"`
	element
}

type Asset struct {
	ID string `json:"id"`
}

type ContentCreateResponse struct {
	Id string `json:"id"`
	Rev string `json:"rev"`
	Name string `json:"name"`
	TypeId string `json:"typeId"`
	Type string `json:"type"`
}

type ContentAuthoringErrorResponse struct {
	RequestId string               `json:"requestId"`
	Service string                 `json:"service"`
	RequestMethod string           `json:"requestMethod"`
	RequestUri string              `json:"requestUri"`
	Type string                    `json:"type"`
	Errors []ContentAuthoringError `json:"errors"`
}

type ContentAuthoringError struct {
	Code int64 `json:"code"`
	Key string `json:"key"`
	Message string `json:"message"`
	Description string `json:"description"`
	MoreInfo string `json:"more_info"`
	Category string `json:"category"`
	Level string `json:"level"`
	Parameters interface{} `json:"parameters"`
	Field interface{} `json:"field"`
	Locale interface{} `json:"locale"`
}

func (element element)  Convert(data interface{}) (Element,error) {
	return nil, errors.New("Not implementd need to override in extending elements")
}

func Build(fieldType string) (Element,error) {
	fieldTypeConst := FieldType(fieldType)
	switch fieldTypeConst {
	case Text:
		element := TextElement{}
		element.ElementType = fieldTypeConst
		return element,nil
	case Number:
		element := NumberElement{}
		element.ElementType = fieldTypeConst
		return element,nil
	case Boolean:
		element := BooleanElement{}
		element.ElementType = fieldTypeConst
		return element,nil
	case Link:
		element := LinkElement{}
		element.ElementType = fieldTypeConst
		return element,nil
	case Date:
		element := DateElement{}
		element.ElementType = fieldTypeConst
		return element,nil
	case Category:
		element := CategoryElement{}
		element.ElementType = fieldTypeConst
		return element,nil
	case Image:
		element := ImageElement{}
		element.ElementType = fieldTypeConst
		return element,nil
	default:
		return nil,errors.New("No element found for property type" + fieldType)
	}
}


