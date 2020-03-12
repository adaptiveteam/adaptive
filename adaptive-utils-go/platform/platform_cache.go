package platform

import (
	"fmt"
	"github.com/ReneKroon/ttlcache"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"log"
	"time"
)

// UserPlatformToken contains authentication information for communicating
// with a user
type UserPlatformToken struct {
	PlatformName        models.PlatformName
	AuthenticationToken string
}

type GetUserPlatformToken func(userID string) UserPlatformToken

// InitLocalCache Initializing a local cache for key-value pairs
// TTL can be set globally and at the key level
func InitLocalCache(cache *ttlcache.Cache) *ttlcache.Cache {
	if cache == nil {
		cache = ttlcache.NewCache()
		var expiration time.Duration 
		expiration = 10 * time.Second
		cache.SetTTL(expiration)
		expirationCallback := func(key string, value interface{}) {
			fmt.Printf("PlatformTokenCache: This key(%s) has expired after duration %v\n", key, expiration)
		}
		cache.SetExpirationCallback(expirationCallback)
	}
	return cache
}

func UserPlatformTokenFromCache(userID string, cache *ttlcache.Cache, token GetUserPlatformToken,
	ttl time.Duration) (upt UserPlatformToken) {
	value, exists := cache.Get(userID)
	if !exists {
		log.Println(fmt.Sprintf("Profile info for %s doesn't exist in local cache. Querying table", userID))
		upt = token(userID)
		cache.SetWithTTL(userID, upt, ttl)
	} else {
		log.Println(fmt.Sprintf("Profile info for %s exists in local cache", userID))
		upt = value.(UserPlatformToken)
	}
	return
}
