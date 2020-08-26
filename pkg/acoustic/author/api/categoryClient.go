package api

import (
	"errors"
	"github.com/wesovilabs/koazee"
	"gopkg.in/resty.v1"
	"strconv"
)

type Categories struct {
	Offset int `json:"offset"`
	Limit int `json:"limit"`
	Next string `json:"next"`
	Items []CategoryItem `json:"items"`
}

type CategoryItem struct {
	Id string `json:"id"`
	NamePath []string `json:"namePath"`
}

func (categoryItem *CategoryItem) IsMatchingCategory(categoryName string) bool {
	return categoryName == categoryItem.NamePath[0] &&
		len(categoryItem.NamePath) > 1
}

func (categoryItem *CategoryItem) FullNamePath() string {
	return koazee.StreamOf(categoryItem.NamePath).
		Reduce(func(acc string, namePath string) string {
			if acc == "" {
				acc += namePath
			} else {
				acc += "/" + namePath
			}
			return acc
	}).String()
}

type CategoryClient interface {
	Categories(categoryName string) ([]CategoryItem,error)
}

type categoryClient struct {
	c *resty.Client
	acousticApiUrl string
}

func NewCategoryClient(acousticApiUrl string) CategoryClient {
	return &categoryClient{
		c: Connect(),
		acousticApiUrl: acousticApiUrl,
	}
}

func (categoryClient *categoryClient) 	Categories(categoryName string) ([]CategoryItem,error) {



	categoryItems := make([]CategoryItem,0,10)
	offSet := 0
	for {
		req := categoryClient.c.NewRequest().
			SetResult(&Categories{}).SetQueryParam("offset",strconv.Itoa(offSet)).SetQueryParam("limit","100")
		if resp, err := req.Get(categoryClient.acousticApiUrl + "/authoring/v2/categories") ; err != nil {
			return nil,err
		} else if resp.IsSuccess() {
			categories := resp.Result().(*Categories)

			offSet += categories.Limit
			koazee.StreamOf(categories.Items).
				ForEach(func(categoryItem CategoryItem)  {
					if categoryItem.IsMatchingCategory(categoryName) {
						categoryItems = append(categoryItems,categoryItem)
					}
			}).Do()
			if categories.Next == "" {
				break
			}
		} else {
			return nil, errors.New("error in creating content : " + resp.Status())
		}
	}
	return categoryItems,nil

}
