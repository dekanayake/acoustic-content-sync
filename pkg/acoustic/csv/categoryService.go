package csv

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/author/api"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/wesovilabs/koazee"
	"strings"
)

type CategoryService interface {
	Create(categoryName string, dataFeedPath string, configPath string) error
}

type categoryService struct {
	acousticAuthApiUrl string
	categoryClient     api.CategoryClient
}

func NewCategoryService(acousticAuthApiUrl string) CategoryService {
	return &categoryService{
		acousticAuthApiUrl: acousticAuthApiUrl,
		categoryClient:     api.NewCategoryClient(acousticAuthApiUrl),
	}
}

type category struct {
	rootCategory     string
	fullCategoryPath string
}

func (category category) fullNamePath() []string {
	return strings.Split(category.fullCategoryPath, "/")
}

func (category category) childCategory() string {
	namePaths := category.fullNamePath()
	return namePaths[len(namePaths)-1]
}

func (category category) parentCategory() string {
	namePaths := category.fullNamePath()
	if len(namePaths) == 1 {
		return namePaths[0]
	} else {
		namePaths := namePaths[0 : len(namePaths)-1]
		fullNamePath := koazee.StreamOf(namePaths).
			Reduce(func(acc string, namePath string) string {
				if acc == "" {
					acc += namePath
				} else {
					acc += "/" + namePath
				}
				return acc
			}).Val().(string)
		return fullNamePath
	}

}

func createCategory(newCategoryPath string, rootCategory string, existingCategories map[string]string, categoryClient api.CategoryClient) (map[string]string, error) {
	newCategory := &category{
		fullCategoryPath: newCategoryPath,
		rootCategory:     rootCategory,
	}
	parentCategoryPath := newCategory.parentCategory()
	parentCategoryID := existingCategories[parentCategoryPath]
	if parentCategoryID == "" {
		createdCategories, err := createCategory(parentCategoryPath, rootCategory, existingCategories, categoryClient)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		parentCategoryID := createdCategories[parentCategoryPath]
		if parentCategoryID == "" {
			return nil, errors.ErrorMessageWithStack("No created parent category id found : " + parentCategoryPath)
		}
		newCategory, err := categoryClient.CreateCategory(parentCategoryID, newCategory.childCategory())
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		createdCategories[newCategory.FullNamePath()] = newCategory.Id
		return createdCategories, nil
	} else {
		newCategory, err := categoryClient.CreateCategory(parentCategoryID, newCategory.childCategory())
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		createdCategories := make(map[string]string, 0)
		createdCategories[newCategory.FullNamePath()] = newCategory.Id
		return createdCategories, nil
	}
}

func (c categoryService) Create(categoryName string, dataFeedPath string, configPath string) error {
	categories, err := c.categoryClient.Categories(categoryName)
	if err != nil {
		return errors.ErrorWithStack(err)
	}

	config, err := InitConfig(configPath)
	if err != nil {
		return errors.ErrorWithStack(err)
	}
	categoryMapping, err := config.GetCategory(categoryName)
	if err != nil {
		return errors.ErrorWithStack(err)
	}
	dataFeed, err := LoadCSV(dataFeedPath)
	if err != nil {
		return errors.ErrorWithStack(err)
	}

	newCategories := make([]string, 0, dataFeed.RecordCount())
	for ok := true; ok; ok = dataFeed.HasNext() {
		dataRow := dataFeed.Next()
		val, err := dataRow.Get(categoryMapping.Column)
		if err != nil {
			return errors.ErrorWithStack(err)
		}
		for _, category := range strings.Split(val, ",") {
			newCategories = append(newCategories, categoryName+"/"+strings.TrimSpace(category))
		}

	}

	existingCategories := koazee.StreamOf(categories).
		Reduce(func(acc map[string]string, categoryItem api.CategoryItem) map[string]string {
			if acc == nil {
				acc = make(map[string]string, 0)
			}
			acc[categoryItem.FullNamePath()] = categoryItem.Id
			return acc
		}).Val().(map[string]string)

	newCategories = koazee.StreamOf(newCategories).
		Filter(func(newCategory string) bool {
			_, ok := existingCategories[newCategory]
			return !ok
		}).RemoveDuplicates().Out().Val().([]string)

	err = koazee.StreamOf(newCategories).
		ForEach(func(newCategoryPath string) error {
			createdCategories, err := createCategory(newCategoryPath, categoryName, existingCategories, c.categoryClient)
			if err != nil {
				return errors.ErrorWithStack(err)
			}
			for k, v := range createdCategories {
				existingCategories[k] = v
			}
			return nil
		}).Do().Out().Err().UserError()
	if err != nil {
		return errors.ErrorWithStack(err)
	}
	return nil

}
