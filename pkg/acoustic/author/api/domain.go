package api

import "errors"

type Type string

const (
	Text Type = "text"
	Number Type = "number"
	Boolean Type = "toggle"
	Link Type = "link"
	Date Type = "date"
	Category Type = "category"
	File Type = "file"
	Video Type = "video"
	Image Type = "image"
)


type Content struct {
	Name string `json:"name"`
	TypeId string `json:"typeId"`
	Status string `json:"status"`
	Elements  map[string] interface{} `json:"elements"`
	LibraryID string `json:"libraryId"`
}

type Convertor func(data interface{}, element Element) (Element,error)

type Element interface {
	Enrich(element interface{}, converter Convertor) (Element,error)
}

type element struct {
	ElementType Type `json:"elementType"`
}



type TextElement struct {
	Value string `json:"value"`
	element
}

type NumberElement struct {
	Value int `json:"value"`
	element
}

type BooleanElement struct {
	Value bool `json:"value"`
	element
}

type LinkElement struct {
	LinkURL int `json:"linkURL"`
	LinkText int `json:"linkText"`
	LinkTitle int `json:"linkTitle"`
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

type ContentCreateResponse struct {
	Id string `json:"id"`
	Rev string `json:"rev"`
	Name string `json:"name"`
	TypeId string `json:"typeId"`
	Type string `json:"type"`
}

type ContentCreateErrorResponse struct {
	RequestId string `json:"requestId"`
	Service string `json:"service"`
	RequestMethod string `json:"requestMethod"`
	RequestUri string `json:"requestUri"`
	Type string `json:"type"`
	Errors []ContentError `json:"errors"`
}

type ContentError struct {
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

func (acousticElement element) Enrich(element interface{}, converter Convertor) (Element,error) {
	convertedElement,err := converter(element,&acousticElement)
	return convertedElement,err
}

func Build(fieldType string) (Element,error) {
	switch fieldType {
	case Text:
		element := TextElement{}
		element.ElementType = fieldType
		return element,nil
	case Number:
		element := NumberElement{}
		element.ElementType = fieldType
		return element,nil
	case Boolean:
		element := BooleanElement{}
		element.ElementType = fieldType
		return element,nil
	case Link:
		element := LinkElement{}
		element.ElementType = fieldType
		return element,nil
	case Date:
		element := DateElement{}
		element.ElementType = fieldType
		return element,nil
	case Category:
		element := CategoryElement{}
		element.ElementType = fieldType
		return element,nil
	default:
		return nil,errors.New("No element found for property type" + fieldType)
	}
}


