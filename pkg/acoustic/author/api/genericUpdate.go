package api

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/thoas/go-funk"
	"golang.org/x/exp/slices"
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

func (element FloatElement) Update(new Element) (Element, error) {
	oldValue := element.Value
	newValue := new.(FloatElement).Value
	if oldValue != newValue {
		element.Value = newValue
		return element, nil
	} else {
		return nil, nil
	}
}

func (element FormattedTextElement) Update(new Element) (Element, error) {
	oldValue := element.Value
	newValue := new.(FormattedTextElement).Value
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

func (element BooleanElement) Update(new Element) (Element, error) {
	oldValue := element.Value
	newValue := new.(BooleanElement).Value
	if oldValue != newValue {
		element.Value = newValue
		return element, nil
	} else {
		return nil, nil
	}
}

func (element OptionSelectionElement) Update(new Element) (Element, error) {
	//TODO implement me
	panic("implement me")
}

func (element MultiTextElement) Update(new Element) (Element, error) {
	return nil, errors.ErrorMessageWithStack("not implemented")
}

func (element NumberElement) Update(new Element) (Element, error) {
	oldValue := element.Value
	newValue := new.(NumberElement).Value
	if oldValue != newValue {
		element.Value = newValue
		return element, nil
	} else {
		return nil, nil
	}
}

func (element LinkElement) Update(new Element) (Element, error) {
	oldLinkUrl := element.LinkURL
	newLinkUrl := new.(LinkElement).LinkURL
	if oldLinkUrl != newLinkUrl {
		element.LinkURL = newLinkUrl
	}

	oldLinkText := element.LinkText
	newLinkText := new.(LinkElement).LinkText
	if oldLinkText != newLinkText {
		element.LinkText = newLinkText
	}

	oldLinkTitle := element.LinkTitle
	newLinkTitle := new.(LinkElement).LinkTitle
	if oldLinkTitle != newLinkTitle {
		element.LinkTitle = newLinkTitle
	}

	return element, nil
}

func (element CategoryElement) Update(new Element) (Element, error) {
	oldCatIds := element.CategoryIds
	newElement := new.(CategoryElement)
	newCatIds := newElement.CategoryIds
	updatedCatIds := oldCatIds
	if new.GetOperation() == DELETE {
		for index, removeCat := range updatedCatIds {
			if slices.Contains(newCatIds, removeCat) {
				updatedCatIds = append(updatedCatIds[:index], updatedCatIds[index+1:]...)
			}
		}
	} else {
		updatedCatIds = append(newCatIds, updatedCatIds...)
	}

	newElement.CategoryIds = funk.UniqString(updatedCatIds)
	return newElement, nil
}

func (element CategoryPartElement) Update(new Element) (Element, error) {
	return nil, errors.ErrorMessageWithStack("not implemented")
}

func (element ImageElement) Update(new Element) (Element, error) {
	newElement := new.(ImageElement)
	newElement.Asset.ID = element.Asset.ID
	return newElement, nil
}

func (element FileElement) Update(new Element) (Element, error) {
	newElement := new.(FileElement)
	newElement.Asset.ID = element.Asset.ID
	return newElement, nil
}

func (element GroupElement) Update(new Element) (Element, error) {
	newValue := new.(GroupElement)
	for k, v := range element.Value {
		if newValue.Value[k] != nil {
			updatedOldVal, err := v.(Element).Update(newValue.Value[k].(Element))
			if err != nil {
				return nil, err
			}
			if updatedOldVal != nil {
				element.Value[k] = updatedOldVal
			}
		}
	}
	return element, nil
}

func (element MultiGroupElement) Update(new Element) (Element, error) {
	newValue := new.(MultiGroupElement)
	return newValue, nil
}

func (element ReferenceElement) Update(new Element) (Element, error) {
	return nil, errors.ErrorMessageWithStack("not implemented")
}

func (element MultiReferenceElement) Update(new Element) (Element, error) {
	oldValues := element.Values
	newValues := new.(MultiReferenceElement).Values
	if new.GetOperation() == DELETE {
		tempValues := make([]ReferenceValue, 0)
		for _, oldValue := range oldValues {
			if !slices.Contains(newValues, oldValue) {
				tempValues = append(tempValues, oldValue)
			}
		}
		oldValues = tempValues
	} else {
		oldValues = append(oldValues, newValues...)
	}

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

func (m MultiImageElement) Update(new Element) (Element, error) {
	//TODO implement me
	panic("implement me")
}

func (element MultiOptionSelectionElement) Update(new Element) (Element, error) {
	//TODO implement me
	panic("implement me")
}

func (m MultiNumberElement) Update(new Element) (Element, error) {
	//TODO implement me
	panic("implement me")
}
