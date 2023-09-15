package api

import "github.com/dekanayake/acoustic-content-sync/pkg/errors"

func (element TextElement) Clone() (Element, error) {
	clonedElement := TextElement{}
	clonedElement.ElementType = element.ElementType
	clonedElement.Value = element.Value
	return clonedElement, nil
}

func (element MultiTextElement) Clone() (Element, error) {
	clonedElement := MultiTextElement{}
	clonedElement.ElementType = element.ElementType
	clonedElement.Values = element.Values
	return clonedElement, nil
}

func (element FormattedTextElement) Clone() (Element, error) {
	clonedElement := FormattedTextElement{}
	clonedElement.ElementType = element.ElementType
	clonedElement.Value = element.Value
	return clonedElement, nil
}

func (element NumberElement) Clone() (Element, error) {
	clonedElement := NumberElement{}
	clonedElement.ElementType = element.ElementType
	clonedElement.Value = element.Value
	return clonedElement, nil
}

func (element MultiNumberElement) Clone() (Element, error) {
	clonedElement := MultiNumberElement{}
	clonedElement.ElementType = element.ElementType
	clonedElement.Values = element.Values
	return clonedElement, nil
}

func (element FloatElement) Clone() (Element, error) {
	clonedElement := FloatElement{}
	clonedElement.ElementType = element.ElementType
	clonedElement.Value = element.Value
	return clonedElement, nil
}

func (element BooleanElement) Clone() (Element, error) {
	clonedElement := BooleanElement{}
	clonedElement.ElementType = element.ElementType
	clonedElement.Value = element.Value
	return clonedElement, nil
}

func (element LinkElement) Clone() (Element, error) {
	clonedElement := LinkElement{}
	clonedElement.ElementType = element.ElementType
	clonedElement.LinkURL = element.LinkURL
	clonedElement.LinkTitle = element.LinkTitle
	clonedElement.LinkText = element.LinkText
	return clonedElement, nil
}

func (element DateElement) Clone() (Element, error) {
	clonedElement := DateElement{}
	clonedElement.ElementType = element.ElementType
	clonedElement.Value = element.Value
	return clonedElement, nil
}

func (element CategoryElement) Clone() (Element, error) {
	clonedElement := CategoryElement{}
	clonedElement.ElementType = element.ElementType
	clonedElement.Categories = element.Categories
	clonedElement.CategoryIds = element.CategoryIds
	return clonedElement, nil
}

func (element CategoryPartElement) Clone() (Element, error) {
	clonedElement := CategoryPartElement{}
	clonedElement.ElementType = element.ElementType
	clonedElement.CategoryIds = element.CategoryIds
	return clonedElement, nil
}

func (element ImageElement) Clone() (Element, error) {
	clonedElement := ImageElement{}
	clonedElement.ElementType = element.ElementType
	clonedElement.Mode = element.Mode
	if element.Asset != nil {
		clonedElement.Asset = &Asset{
			ID: element.Asset.ID,
		}
	}
	return clonedElement, nil
}

func (element MultiImageElement) Clone() (Element, error) {
	clonedElement := MultiImageElement{}
	clonedElement.ElementType = element.ElementType
	imageElements := make([]ImageElementItem, 0)
	for _, imageElement := range element.Values {
		clonedImageElement := ImageElementItem{}
		clonedImageElement.Mode = imageElement.Mode
		if imageElement.Asset != nil {
			clonedImageElement.Asset = &Asset{
				ID: imageElement.Asset.ID,
			}
		}
		imageElements = append(imageElements, clonedImageElement)
	}
	clonedElement.Values = imageElements
	return clonedElement, nil
}

func (element FileElement) Clone() (Element, error) {
	clonedElement := FileElement{}
	clonedElement.ElementType = element.ElementType
	clonedElement.Asset = Asset{
		ID: element.Asset.ID,
	}
	return clonedElement, nil
}

func (element GroupElement) Clone() (Element, error) {
	clonedElement := GroupElement{}
	clonedElement.ElementType = element.ElementType
	clonedElement.TypeRef = element.TypeRef
	clonedElement.Value = make(map[string]interface{}, 0)
	for fieldName, childElement := range element.Value {
		clonedChildElement, err := childElement.(Element).Clone()
		if err != nil {
			errors.ErrorWithStack(err)
		}
		clonedElement.Value[fieldName] = clonedChildElement
	}
	if len(clonedElement.Value) == 0 {
		clonedElement.Value = nil
	}
	return clonedElement, nil
}

func (element MultiGroupElement) Clone() (Element, error) {
	clonedElement := MultiGroupElement{}
	clonedElement.ElementType = element.ElementType
	clonedElement.TypeRef = element.TypeRef
	clonedElement.Values = make([]map[string]interface{}, 0)
	for _, groupElements := range element.Values {
		clonedChildElementsMap := make(map[string]interface{})
		for fieldName, childElement := range groupElements {
			clonedChildElement, err := childElement.(Element).Clone()
			if err != nil {
				errors.ErrorWithStack(err)
			}
			clonedChildElementsMap[fieldName] = clonedChildElement
		}
		clonedElement.Values = append(clonedElement.Values, clonedChildElementsMap)
	}
	return clonedElement, nil
}

func (element ReferenceElement) Clone() (Element, error) {
	clonedElement := ReferenceElement{}
	clonedElement.ElementType = element.ElementType
	if element.Value != nil && element.Value.ID != "" {
		clonedElement.Value = &ReferenceValue{
			ID: element.Value.ID,
		}
	}
	return clonedElement, nil
}

func (element MultiReferenceElement) Clone() (Element, error) {
	clonedElement := MultiReferenceElement{}
	clonedElement.ElementType = element.ElementType
	clonedElement.Values = make([]ReferenceValue, 0)
	for _, value := range element.Values {
		clonedElement.Values = append(clonedElement.Values, ReferenceValue{
			ID: value.ID,
		})
	}
	return clonedElement, nil
}

func (element OptionSelectionElement) Clone() (Element, error) {
	clonedElement := OptionSelectionElement{}
	clonedElement.ElementType = element.ElementType
	if element.Value != nil && element.Value.Selection != "" {
		clonedElement.Value = &OptionSelectionValue{
			Selection: element.Value.Selection,
		}
	}
	return clonedElement, nil
}

func (element MultiOptionSelectionElement) Clone() (Element, error) {
	clonedElement := MultiOptionSelectionElement{}
	clonedElement.ElementType = element.ElementType
	clonedElement.Values = make([]*OptionSelectionValue, 0)
	for _, existingOptionSelectionValue := range clonedElement.Values {
		if existingOptionSelectionValue != nil && existingOptionSelectionValue.Selection != "" {
			clonedElement.Values = append(clonedElement.Values, &OptionSelectionValue{
				Selection: existingOptionSelectionValue.Selection,
			})
		}
	}
	return clonedElement, nil
}

func (element DateTimeElement) Clone() (Element, error) {
	clonedElement := DateTimeElement{}
	clonedElement.ElementType = element.ElementType
	clonedElement.Value = element.Value
	return clonedElement, nil
}

func (element MultiLinkElement) Clone() (Element, error) {
	clonedElement := MultiLinkElement{}
	clonedElement.ElementType = element.ElementType
	links := element.Values
	clonedLinks := make([]LinkElement, 0)
	for _, link := range links {
		clonedLink, err := link.Clone()
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		clonedLinks = append(clonedLinks, clonedLink.(LinkElement))
	}
	clonedElement.Values = clonedLinks
	return clonedElement, nil
}
