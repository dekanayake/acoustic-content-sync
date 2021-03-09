package api

import (
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var createOnce sync.Once

var cachedCategoryClientInstance *cachedCategoryClient

type cachedCategoryClient struct {
	categoryClient CategoryClient
	cache          *cache.Cache
}

func NewCachedCategoryClient(acousticApiUrl string) CategoryClient {
	createOnce.Do(func() {
		cachedCategoryClientInstance = &cachedCategoryClient{
			categoryClient: NewCategoryClient(acousticApiUrl),
			cache:          cache.New(1*time.Hour, 2*time.Hour),
		}
	})
	return cachedCategoryClientInstance
}

func (c cachedCategoryClient) Categories(categoryName string) ([]CategoryItem, error) {
	cached, found := c.cache.Get(categoryName)
	if found {
		return cached.([]CategoryItem), nil
	} else {
		categoryResponse, err := c.categoryClient.Categories(categoryName)
		if err != nil {
			return nil, err
		} else {
			c.cache.Set(categoryName, categoryResponse, cache.DefaultExpiration)
			return categoryResponse, nil
		}
	}
}

func (c cachedCategoryClient) CreateCategory(parentCategoryID string, categoryName string) (CategoryItem, error) {
	return c.categoryClient.CreateCategory(parentCategoryID, categoryName)
}

func (c cachedCategoryClient) DeleteCategory(categoryID string) error {
	return c.categoryClient.DeleteCategory(categoryID)
}
