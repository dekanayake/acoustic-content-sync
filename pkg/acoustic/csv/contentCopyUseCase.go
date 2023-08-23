package csv

import (
	"fmt"
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/author/api"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/golang-collections/collections/stack"
	log "github.com/sirupsen/logrus"
	"reflect"
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

		if reflect.TypeOf(existingElement).Name() == "ReferenceElement" {
			id := existingElement.(api.ReferenceElement).Value.ID
			referenceContent, err := c.contentClient.Get(id)
			if err != nil {
				return nil, errors.ErrorWithStack(err)
			}
			result[name] = []*api.Content{referenceContent}
		}

		if reflect.TypeOf(existingElement).Name() == "MultiReferenceElement" {
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

		if reflect.TypeOf(existingElement).Name() == "GroupElement" {
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

		if reflect.TypeOf(existingElement).Name() == "MultiGroupElement" {
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
	log.Info(levelsMap)
	return nil, nil
}
