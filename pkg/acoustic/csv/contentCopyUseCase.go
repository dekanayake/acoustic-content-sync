package csv

import (
	"fmt"
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/author/api"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/golang-collections/collections/stack"
	"sort"
)

type contentCopyUserCase struct {
	contentClient api.ContentClient
}

type ContentCopyUserCase interface {
	CopyContent(id string) (*ContentCreationStatus, error)
}

func NewContentCopyUserCase(acousticAuthApiUrl string) ContentCopyUserCase {
	return &contentCopyUserCase{
		contentClient: api.NewContentClient(acousticAuthApiUrl),
	}
}

type contentContainer struct {
	id       string                         `json:"id"`
	children map[string][]*contentContainer `json:"children"`
	parent   *contentContainer              `json:"parent"`
}

func (c contentCopyUserCase) getChildReference(elements map[string]interface{}) (map[string][]*api.Content, error) {
	result := make(map[string][]*api.Content)
	for name, element := range elements {
		_, isElementMap := element.(map[string]interface{})
		var existingElement api.Element
		if isElementMap {
			var err error
			existingElement, err = api.Convert(element.(map[string]interface{}))
			if err != nil {
				return nil, errors.ErrorWithStack(err)
			}
		} else {
			existingElement = element.(api.Element)
		}

		if existingElement.Type() == "ReferenceElement" {
			id := existingElement.(api.ReferenceElement).Value.ID
			referenceContent, err := c.contentClient.Get(id)
			if err != nil {
				return nil, errors.ErrorWithStack(err)
			}
			result[name] = []*api.Content{referenceContent}
		}

		if existingElement.Type() == "MultiReferenceElement" {
			idValues := existingElement.(api.MultiReferenceElement).Values
			childReferenceContentList := make([]*api.Content, 0)
			for _, idValue := range idValues {
				referenceContent, err := c.contentClient.Get(idValue.ID)
				if err != nil {
					return nil, errors.ErrorWithStack(err)
				}
				childReferenceContentList = append(childReferenceContentList, referenceContent)
			}
			result[name] = childReferenceContentList
		}

		if existingElement.Type() == "GroupElement" {
			groupElements := existingElement.(api.GroupElement).Value
			groupElementRefChilds, err := c.getChildReference(groupElements)
			if err != nil {
				return nil, errors.ErrorWithStack(err)
			}
			for groupElementName, groupElementRef := range groupElementRefChilds {
				groupRefName := name + "." + groupElementName
				result[groupRefName] = groupElementRef
			}
		}

		if existingElement.Type() == "MultiGroupElement" {
			groupElementsList := existingElement.(api.MultiGroupElement).Values
			for _, groupElement := range groupElementsList {
				groupElementRefChilds, err := c.getChildReference(groupElement)
				if err != nil {
					return nil, errors.ErrorWithStack(err)
				}
				for groupElementName, groupElementRef := range groupElementRefChilds {
					groupRefName := name + "." + groupElementName
					result[groupRefName] = append(groupElementRef, result[groupRefName]...)
				}
			}

		}
	}
	return result, nil
}

func (c contentCopyUserCase) print(content *contentContainer) {
	contentStack := stack.New()
	contentStack.Push([]interface{}{"", "", content})
	for contentStack.Peek() != nil {
		parentContent := contentStack.Pop().([]interface{})
		refSpace := parentContent[0].(string)
		elementFieldName := parentContent[1].(string)
		ref := parentContent[2].(*contentContainer)
		fmt.Printf("%s%s%s\n", refSpace, elementFieldName, ref.id)

		for name, children := range ref.children {
			elementFieldName = name + ":"
			refSpace = refSpace + "-"
			for _, child := range children {
				contentStack.Push([]interface{}{refSpace, elementFieldName, child})
			}
		}
	}
}

func (c contentCopyUserCase) getLevels(content *contentContainer) map[int][]*contentContainer {
	contentStack := stack.New()
	output := make(map[int][]*contentContainer, 0)
	level := 0
	contentStack.Push([]interface{}{level, content})
	for contentStack.Peek() != nil {
		parentContent := contentStack.Pop().([]interface{})
		level := parentContent[0].(int)
		ref := parentContent[1].(*contentContainer)
		_, levelExists := output[level]
		if !levelExists {
			output[level] = make([]*contentContainer, 0)
		}
		output[level] = append(output[level], ref)

		for _, children := range ref.children {
			level := level + 1
			for _, child := range children {
				contentStack.Push([]interface{}{level, child})
			}
		}
	}
	return output
}

func (c contentCopyUserCase) cloneContentElements(content *api.Content) (*api.Content, error) {
	elements := content.Elements
	for fieldName, element := range elements {
		existingElement, err := api.Convert(element.(map[string]interface{}))
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		clonedElement, err := existingElement.Clone()
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		content.Elements[fieldName] = clonedElement
	}
	return content, nil
}

func (c contentCopyUserCase) convertElementToInterfaceMap(elementsMap map[string]api.Element) map[string]interface{} {
	elementToInterfaceMap := make(map[string]interface{}, 0)
	for name, groupElement := range elementsMap {
		elementToInterfaceMap[name] = groupElement.(interface{})
	}
	return elementToInterfaceMap
}

func (c contentCopyUserCase) convertElementInterfaceToElementMap(elementsMap map[string]interface{}) map[string]api.Element {
	elementToInterfaceMap := make(map[string]api.Element, 0)
	for name, groupElement := range elementsMap {
		elementToInterfaceMap[name] = groupElement.(api.Element)
	}
	return elementToInterfaceMap
}

func (c contentCopyUserCase) updateDependencyRefs(elements map[string]api.Element, childReferences map[string]string) (map[string]api.Element, error) {
	for name, element := range elements {
		if element.Type() == "ReferenceElement" {
			referenceElementToUpdate := element.(api.ReferenceElement)
			refId := referenceElementToUpdate.Value.ID
			updatedRefID, mappingRefAvailable := childReferences[refId]
			if mappingRefAvailable {
				referenceElementToUpdate.Value.ID = updatedRefID
				elements[name] = referenceElementToUpdate
			}
		}

		if element.Type() == "MultiReferenceElement" {
			multiReferenceElementToUpdate := element.(api.MultiReferenceElement)
			referencesList := multiReferenceElementToUpdate.Values
			updatedReferenceList := make([]api.ReferenceValue, 0)
			for _, reference := range referencesList {
				updatedRefID, mappingRefAvailable := childReferences[reference.ID]
				if mappingRefAvailable {
					reference.ID = updatedRefID
					updatedReferenceList = append(updatedReferenceList, reference)
				}
			}
			multiReferenceElementToUpdate.Values = updatedReferenceList
			elements[name] = multiReferenceElementToUpdate
		}

		if element.Type() == "GroupElement" {
			groupElementToUpdate := element.(api.GroupElement)
			groupElements := make(map[string]api.Element, 0)
			for name, element := range groupElementToUpdate.Value {
				groupElements[name] = element.(api.Element)
			}
			updatedGroupElement, err := c.updateDependencyRefs(groupElements, childReferences)
			if err != nil {
				return nil, errors.ErrorWithStack(err)
			}
			groupElementToUpdate.Value = c.convertElementToInterfaceMap(updatedGroupElement)
			elements[name] = groupElementToUpdate
		}

		if element.Type() == "MultiGroupElement" {
			groupElementToUpdate := element.(api.MultiGroupElement)
			groupElements := make(map[string]api.Element, 0)
			groupElementList := groupElementToUpdate.Values
			updatedGroupElementList := make([]map[string]interface{}, 0)
			for _, groupValue := range groupElementList {
				for name, element := range groupValue {
					groupElements[name] = element.(api.Element)
				}
				updatedGroupElement, err := c.updateDependencyRefs(groupElements, childReferences)
				if err != nil {
					return nil, errors.ErrorWithStack(err)
				}
				updatedGroupElementList = append(updatedGroupElementList, c.convertElementToInterfaceMap(updatedGroupElement))
			}
			groupElementToUpdate.Values = updatedGroupElementList
			elements[name] = groupElementToUpdate
		}
	}
	return elements, nil
}

func (c contentCopyUserCase) cloneContent(content *api.Content, childReferences map[string]string) (*api.Content, error) {
	clonedContent, err := c.cloneContentElements(content)
	clonedContent.ID = ""
	clonedContent.REV = ""
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	updatedElements, err := c.updateDependencyRefs(c.convertElementInterfaceToElementMap(clonedContent.Elements), childReferences)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	clonedContent.Elements = c.convertElementToInterfaceMap(updatedElements)
	return &api.Content{
		Elements:  clonedContent.Elements,
		Status:    clonedContent.Status,
		TypeId:    clonedContent.TypeId,
		LibraryID: clonedContent.LibraryID,
		Tags:      clonedContent.Tags,
	}, nil
}

func (c contentCopyUserCase) clone(contentContainerList []*contentContainer, childReferences map[string]string) (map[string]string, error) {
	if childReferences == nil {
		childReferences = make(map[string]string, 0)
	}
	clonedParentReferences := make(map[string]string, 0)
	for _, contentContainer := range contentContainerList {
		content, err := c.contentClient.Get(contentContainer.id)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		originalContentID := content.ID
		clonedContent, err := c.cloneContent(content, childReferences)
		clonedContent.Name = content.Name + "_cloned"
		contentAuthoringResponse, err := c.contentClient.Create(*clonedContent)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		clonedParentReferences[originalContentID] = contentAuthoringResponse.Id
	}
	return clonedParentReferences, nil
}

func (c contentCopyUserCase) CopyContent(id string) (*ContentCreationStatus, error) {
	parentContent, err := c.contentClient.Get(id)
	if err != nil {
		return nil, err
	}
	contentStack := stack.New()
	parentContentContainer := contentContainer{
		id:       parentContent.ID,
		children: make(map[string][]*contentContainer),
	}
	contentStack.Push(&parentContentContainer)
	for contentStack.Peek() != nil {
		parentContentInStack := contentStack.Pop().(*contentContainer)
		parentContent, err := c.contentClient.Get(parentContentInStack.id)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		childReference, err := c.getChildReference(parentContent.Elements)
		if err != nil {
			return nil, err
		}
		for name, childRefContents := range childReference {
			childRefContentContainers := make([]*contentContainer, 0)
			for _, childRefContent := range childRefContents {
				childContentContainer := contentContainer{
					id:       childRefContent.ID,
					parent:   parentContentInStack,
					children: make(map[string][]*contentContainer),
				}
				childRefContentContainers = append(childRefContentContainers, &childContentContainer)
				contentStack.Push(&childContentContainer)
			}
			parentContentInStack.children[name] = childRefContentContainers
		}
	}
	levelsMap := c.getLevels(&parentContentContainer)
	var levels []int
	for level := range levelsMap {
		levels = append(levels, level)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(levels)))
	childRefMap := make(map[string]string, 0)
	for _, level := range levels {
		childRefMap, err = c.clone(levelsMap[level], childRefMap)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
