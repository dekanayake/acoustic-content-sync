package api

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
)

func (element TextElement) Update(new Element) (Element, error) {
	oldValue := element.Value
	newValue := new.(TextElement).Value
	if oldValue != newValue {
		element.Value = newValue
		return element, nil
	} else {
		return nil, nil
	}
}

func (element DateElement) Update(new Element) (Element, error) {
	return nil, errors.ErrorMessageWithStack("not implemented")
}

func (element FormattedTextElement) Update(new Element) (Element, error) {
	return nil, errors.ErrorMessageWithStack("not implemented")
}

func (element BooleanElement) Update(new Element) (Element, error) {
	return nil, errors.ErrorMessageWithStack("not implemented")
}

func (element MultiTextElement) Update(new Element) (Element, error) {
	return nil, errors.ErrorMessageWithStack("not implemented")
}

func (element NumberElement) Update(new Element) (Element, error) {
	return nil, errors.ErrorMessageWithStack("not implemented")
}

func (element LinkElement) Update(new Element) (Element, error) {
	return nil, errors.ErrorMessageWithStack("not implemented")
}

func (element CategoryElement) Update(new Element) (Element, error) {
	return nil, errors.ErrorMessageWithStack("not implemented")
}

func (element CategoryPartElement) Update(new Element) (Element, error) {
	return nil, errors.ErrorMessageWithStack("not implemented")
}

func (element ImageElement) Update(new Element) (Element, error) {
	return nil, errors.ErrorMessageWithStack("not implemented")
}

func (element FileElement) Update(new Element) (Element, error) {
	return nil, errors.ErrorMessageWithStack("not implemented")
}

func (element GroupElement) Update(new Element) (Element, error) {
	newValue := new.(GroupElement)
	for k, v := range element.Value {
		updatedOldVal, err := v.(Element).Update(newValue.Value[k].(Element))
		if err != nil {
			return nil, err
		}
		element.Value[k] = updatedOldVal
	}
	return element, nil
}

func (element MultiGroupElement) Update(new Element) (Element, error) {
	return nil, errors.ErrorMessageWithStack("not implemented")
}

func (element ReferenceElement) Update(new Element) (Element, error) {
	return nil, errors.ErrorMessageWithStack("not implemented")
}

func (element MultiReferenceElement) Update(new Element) (Element, error) {
	oldValues := element.Values
	newValues := new.(MultiReferenceElement).Values
	oldValues = append(oldValues, newValues...)

	allKeys := make(map[string]bool)
	updatedList := []ReferenceValue{}
	for _, item := range oldValues {
		if _, value := allKeys[item.ID]; !value {
			allKeys[item.ID] = true
			updatedList = append(updatedList, item)
		}
	}
	element.Values = updatedList
	return element, nil
}
