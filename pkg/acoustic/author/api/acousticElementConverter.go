package api

import (
	"encoding/json"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
)

type elementMatcherFunc func(fieldType AcousticFieldType) bool
type isMultiMatcherFunc func() bool
type convertFunc func(acousticElement map[string]interface{}) (Element, error)

type acousticElementConvertor struct {
	elementFactMatcher elementMatcherFunc
	isMultiMatcher     isMultiMatcherFunc
	convert            convertFunc
}

var textElementConverter = acousticElementConvertor{
	elementFactMatcher: elementMatcherFunc(func(fieldType AcousticFieldType) bool {
		return fieldType == AcousticFieldType(AcousticFieldText)
	}),
	isMultiMatcher: isMultiMatcherFunc(func() bool {
		return false
	}),
	convert: convertFunc(func(acousticElement map[string]interface{}) (Element, error) {
		jsonString, err := json.Marshal(acousticElement)
		if err != nil {
			return nil, err
		}
		element := TextElement{}
		err = json.Unmarshal(jsonString, &element)
		if err != nil {
			return nil, err
		}
		return element, nil
	}),
}

var multiReferenceElementConverter = acousticElementConvertor{
	elementFactMatcher: elementMatcherFunc(func(fieldType AcousticFieldType) bool {
		return fieldType == AcousticFieldType(AcousticFieldReference)
	}),
	isMultiMatcher: isMultiMatcherFunc(func() bool {
		return true
	}),
	convert: convertFunc(func(acousticElement map[string]interface{}) (Element, error) {
		jsonString, err := json.Marshal(acousticElement)
		if err != nil {
			return nil, err
		}
		element := MultiReferenceElement{}
		err = json.Unmarshal(jsonString, &element)
		if err != nil {
			return nil, err
		}
		return element, nil
	}),
}

var converterList = []acousticElementConvertor{
	textElementConverter,
	multiReferenceElementConverter,
}

func Convert(acousticElementData map[string]interface{}) (Element, error) {
	jsonString, err := json.Marshal(acousticElementData)
	if err != nil {
		return nil, err
	}
	element := element{}
	json.Unmarshal(jsonString, &element)
	if err != nil {
		return nil, err
	}
	for _, converter := range converterList {
		if converter.elementFactMatcher(element.ElementType) {
			_, multiOk := acousticElementData["values"]
			if multiOk && converter.isMultiMatcher() {
				converted, err := converter.convert(acousticElementData)
				if err != nil {
					return nil, err
				}
				return converted, nil
			} else {
				converted, err := converter.convert(acousticElementData)
				if err != nil {
					return nil, err
				}
				return converted, nil
			}
		}
	}
	return nil, errors.ErrorMessageWithStack("No converter found for element type :" + string(element.ElementType))
}
