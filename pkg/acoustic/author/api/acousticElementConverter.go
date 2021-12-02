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

var formattedTextElementConverter = acousticElementConvertor{
	elementFactMatcher: elementMatcherFunc(func(fieldType AcousticFieldType) bool {
		return fieldType == AcousticFieldType(AcousticFieldFormattedText)
	}),
	isMultiMatcher: isMultiMatcherFunc(func() bool {
		return false
	}),
	convert: convertFunc(func(acousticElement map[string]interface{}) (Element, error) {
		jsonString, err := json.Marshal(acousticElement)
		if err != nil {
			return nil, err
		}
		element := FormattedTextElement{}
		err = json.Unmarshal(jsonString, &element)
		if err != nil {
			return nil, err
		}
		return element, nil
	}),
}

var booleanElementConverter = acousticElementConvertor{
	elementFactMatcher: elementMatcherFunc(func(fieldType AcousticFieldType) bool {
		return fieldType == AcousticFieldType(AcousticFieldBoolean)
	}),
	isMultiMatcher: isMultiMatcherFunc(func() bool {
		return false
	}),
	convert: convertFunc(func(acousticElement map[string]interface{}) (Element, error) {
		jsonBool, err := json.Marshal(acousticElement)
		if err != nil {
			return nil, err
		}
		element := BooleanElement{}
		err = json.Unmarshal(jsonBool, &element)
		if err != nil {
			return nil, err
		}
		return element, nil
	}),
}

var numberElementConverter = acousticElementConvertor{
	elementFactMatcher: elementMatcherFunc(func(fieldType AcousticFieldType) bool {
		return fieldType == AcousticFieldType(AcousticFieldNumber)
	}),
	isMultiMatcher: isMultiMatcherFunc(func() bool {
		return false
	}),
	convert: convertFunc(func(acousticElement map[string]interface{}) (Element, error) {
		jsonString, err := json.Marshal(acousticElement)
		if err != nil {
			return nil, err
		}
		element := NumberElement{}
		err = json.Unmarshal(jsonString, &element)
		if err != nil {
			//if unmarshelling fail one reason can be the value is stored as a float so , try unmarshelling to float
			element := FloatElement{}
			err = json.Unmarshal(jsonString, &element)
			if err != nil {
				return nil, err
			}
			return element, nil
		}
		return element, nil
	}),
}

var fileElementConverter = acousticElementConvertor{
	elementFactMatcher: elementMatcherFunc(func(fieldType AcousticFieldType) bool {
		return fieldType == AcousticFieldType(AcousticFieldFile)
	}),
	isMultiMatcher: isMultiMatcherFunc(func() bool {
		return false
	}),
	convert: convertFunc(func(acousticElement map[string]interface{}) (Element, error) {
		jsonString, err := json.Marshal(acousticElement)
		if err != nil {
			return nil, err
		}
		element := FileElement{}
		err = json.Unmarshal(jsonString, &element)
		if err != nil {
			return nil, err
		}
		return element, nil
	}),
}

var imageElementConverter = acousticElementConvertor{
	elementFactMatcher: elementMatcherFunc(func(fieldType AcousticFieldType) bool {
		return fieldType == AcousticFieldType(AcousticFieldImage)
	}),
	isMultiMatcher: isMultiMatcherFunc(func() bool {
		return false
	}),
	convert: convertFunc(func(acousticElement map[string]interface{}) (Element, error) {
		jsonString, err := json.Marshal(acousticElement)
		if err != nil {
			return nil, err
		}
		element := ImageElement{}
		err = json.Unmarshal(jsonString, &element)
		if err != nil {
			return nil, err
		}
		return element, nil
	}),
}

var groupElementConverter = acousticElementConvertor{
	elementFactMatcher: elementMatcherFunc(func(fieldType AcousticFieldType) bool {
		return fieldType == AcousticFieldType(AcousticFieldGroup)
	}),
	isMultiMatcher: isMultiMatcherFunc(func() bool {
		return false
	}),
	convert: convertFunc(func(acousticElement map[string]interface{}) (Element, error) {
		jsonString, err := json.Marshal(acousticElement)
		if err != nil {
			return nil, err
		}
		element := GroupElement{}
		err = json.Unmarshal(jsonString, &element)
		if err != nil {
			return nil, err
		}
		for k, v := range element.Value {
			convertedVal, err := Convert(v.(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			element.Value[k] = convertedVal
		}

		return element, nil
	}),
}

var multiGroupElementConverter = acousticElementConvertor{
	elementFactMatcher: elementMatcherFunc(func(fieldType AcousticFieldType) bool {
		return fieldType == AcousticFieldType(AcousticFieldGroup)
	}),
	isMultiMatcher: isMultiMatcherFunc(func() bool {
		return true
	}),
	convert: convertFunc(func(acousticElement map[string]interface{}) (Element, error) {
		jsonString, err := json.Marshal(acousticElement)
		if err != nil {
			return nil, err
		}
		element := MultiGroupElement{}
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

var converterList = make([]acousticElementConvertor, 0)

func init() {
	converterList = []acousticElementConvertor{
		textElementConverter,
		formattedTextElementConverter,
		numberElementConverter,
		multiReferenceElementConverter,
		groupElementConverter,
		multiGroupElementConverter,
		booleanElementConverter,
		fileElementConverter,
		imageElementConverter,
	}
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
			}

			if !multiOk {
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
