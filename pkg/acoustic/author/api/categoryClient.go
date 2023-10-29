package api

import (
	"encoding/json"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/wesovilabs/koazee"
	"gopkg.in/resty.v1"
	"strconv"
)

type Categories struct {
	Offset int            `json:"offset"`
	Limit  int            `json:"limit"`
	Next   string         `json:"next"`
	Items  []CategoryItem `json:"items"`
}

type CategoryItem struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	NamePath []string `json:"namePath"`
}

type CategoryCreateRequest struct {
	Name   string `json:"name"`
	Parent string `json:"parent"`
}

func (categoryItem *CategoryItem) IsMatchingCategory(categoryName string) bool {
	return categoryName == categoryItem.NamePath[0]
}

func (categoryItem *CategoryItem) FullNamePath() string {
	return koazee.StreamOf(categoryItem.NamePath).
		Reduce(func(acc string, namePath string) string {
			if acc == "" {
				acc += namePath
			} else {
				acc += env.CategoryHierarchySeperator() + namePath
			}
			return acc
		}).String()
}

type CategoryClient interface {
	Categories(categoryName string) ([]CategoryItem, error)
	CreateCategory(parentCategoryID string, categoryName string) (CategoryItem, error)
	DeleteCategory(categoryID string) error
}

type categoryClient struct {
	c              *resty.Client
	acousticApiUrl string
}

func NewCategoryClient(acousticApiUrl string) CategoryClient {
	apiKey := env.AcousticAPIKey()
	return &categoryClient{
		c:              Connect(apiKey),
		acousticApiUrl: acousticApiUrl,
	}
}

func (categoryClient categoryClient) CreateCategory(parentCategoryID string, categoryName string) (CategoryItem, error) {
	req := categoryClient.c.NewRequest().
		SetBody(CategoryCreateRequest{Name: categoryName, Parent: parentCategoryID}).
		SetResult(&CategoryItem{})
	if resp, err := req.Post(categoryClient.acousticApiUrl + "/authoring/v1/categories"); err != nil {
		return CategoryItem{}, errors.ErrorWithStack(err)
	} else if resp.IsSuccess() {
		return *resp.Result().(*CategoryItem), nil
	} else if resp.IsError() && resp.StatusCode() == 400 {
		error := resp.Error()
		errorString, _ := json.MarshalIndent(error, "", "\t")
		return CategoryItem{}, errors.ErrorMessageWithStack("error in creating content : " + resp.Status() + "  " + string(errorString))
	} else {
		return CategoryItem{}, errors.ErrorMessageWithStack("error in creating content : " + resp.Status())
	}
}

func (categoryClient categoryClient) DeleteCategory(categoryID string) error {
	req := categoryClient.c.NewRequest().SetPathParams(map[string]string{
		"id": categoryID,
	})
	if resp, err := req.Delete(categoryClient.acousticApiUrl + "/authoring/v1/categories/{id}"); err != nil {
		return errors.ErrorWithStack(err)
	} else if resp.IsSuccess() {
		return nil
	} else if resp.IsError() && resp.StatusCode() == 400 {
		error := resp.Error()
		errorString, _ := json.MarshalIndent(error, "", "\t")
		return errors.ErrorMessageWithStack("error in deleting category : " + resp.Status() + "  " + string(errorString))
	} else {
		return errors.ErrorMessageWithStack("error in deleting category : " + resp.Status())
	}
}

func (categoryClient *categoryClient) Categories(categoryName string) ([]CategoryItem, error) {

	categoryItems := make([]CategoryItem, 0, 10)
	offSet := 0
	for {
		req := categoryClient.c.NewRequest().
			SetResult(&Categories{}).SetQueryParam("offset", strconv.Itoa(offSet)).SetQueryParam("limit", "10000")
		if resp, err := req.Get(categoryClient.acousticApiUrl + "/authoring/v2/categories"); err != nil {
			return nil, errors.ErrorWithStack(err)
		} else if resp.IsSuccess() {
			categories := resp.Result().(*Categories)

			offSet += categories.Limit
			koazee.StreamOf(categories.Items).
				ForEach(func(categoryItem CategoryItem) {
					if categoryItem.IsMatchingCategory(categoryName) {
						categoryItems = append(categoryItems, categoryItem)
					}
				}).Do()
			if categories.Next == "" {
				break
			}
		} else {
			return nil, errors.ErrorMessageWithStack("error in getting category : " + resp.Status())
		}
	}
	return categoryItems, nil

}
