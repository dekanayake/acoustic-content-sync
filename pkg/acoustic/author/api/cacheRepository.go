package api

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var cacheOnce sync.Once

var cacheInstanceMap map[CacheType]*cache.Cache = make(map[CacheType]*cache.Cache)

type CacheType string

const (
	AssetCache CacheType = "Asset"
)

type CacheRepository interface {
	PutCache(cacheType CacheType, key string, value interface{}) error
	GetCache(cacheType CacheType, key string) (interface{}, error)
}

var cacheRepositoryInstanceOnce sync.Once

var cacheRepositoryInstance *cacheRepository

type cacheRepository struct {
	mux *sync.RWMutex
}

func (c cacheRepository) PutCache(cacheType CacheType, key string, value interface{}) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	cacheInstance, err := getCache(cacheType)
	if err != nil {
		return err
	}
	cacheInstance.Add(key, value, cache.NoExpiration)
	return nil
}

func (c cacheRepository) GetCache(cacheType CacheType, key string) (interface{}, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	cacheInstance, err := getCache(cacheType)
	if err != nil {
		return nil, err
	}
	cachedVal, cachedValExist := cacheInstance.Get(key)
	if cachedValExist {
		return cachedVal, nil
	} else {
		return nil, nil
	}
}

func NewCacheRepository() CacheRepository {
	cacheRepositoryInstanceOnce.Do(func() {
		cacheRepositoryInstance = &cacheRepository{
			mux: &sync.RWMutex{},
		}
	})
	return cacheRepositoryInstance
}

func getCache(cacheType CacheType) (*cache.Cache, error) {
	var cacheInstance *cache.Cache = nil
	if cacheRepo, ok := cacheInstanceMap[cacheType]; ok {
		cacheInstance = cacheRepo
	} else {
		cacheInstance = cache.New(24*time.Hour, 24*time.Hour)
		cacheInstanceMap[cacheType] = cacheInstance
	}
	if cacheInstance == nil {
		return nil, errors.ErrorMessageWithStack("Cache is not available")
	}
	return cacheInstance, nil
}
