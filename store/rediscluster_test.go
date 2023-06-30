package store

import (
	"context"
	"testing"
	"time"

	mocksStore "github.com/wlxwlxwlx/gocache/v2/test/mocks/store/clients"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewRedisCluster(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	client := mocksStore.NewMockRedisClusterClientInterface(ctrl)
	options := &Options{
		Expiration: 6 * time.Second,
	}

	// When
	store := NewRedisCluster(client, options)

	// Then
	assert.IsType(t, new(RedisClusterStore), store)
	assert.Equal(t, client, store.clusclient)
	assert.Equal(t, options, store.options)
}

func TestRedisClusterGet(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	client := mocksStore.NewMockRedisClusterClientInterface(ctrl)
	client.EXPECT().Get(ctx, "my-key").Return(&redis.StringCmd{})

	store := NewRedisCluster(client, nil)

	// When
	value, err := store.Get(ctx, "my-key")

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, value)
}

func TestRedisClusterSet(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	cacheKey := "my-key"
	cacheValue := "my-cache-value"
	options := &Options{
		Expiration: 6 * time.Second,
	}

	client := mocksStore.NewMockRedisClusterClientInterface(ctrl)
	client.EXPECT().Set(ctx, "my-key", cacheValue, 5*time.Second).Return(&redis.StatusCmd{})

	store := NewRedisCluster(client, options)

	// When
	err := store.Set(ctx, cacheKey, cacheValue, &Options{
		Expiration: 5 * time.Second,
	})

	// Then
	assert.Nil(t, err)
}

func TestRedisClusterSetWhenNoOptionsGiven(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	cacheKey := "my-key"
	cacheValue := "my-cache-value"
	options := &Options{
		Expiration: 6 * time.Second,
	}

	client := mocksStore.NewMockRedisClusterClientInterface(ctrl)
	client.EXPECT().Set(ctx, "my-key", cacheValue, 6*time.Second).Return(&redis.StatusCmd{})

	store := NewRedisCluster(client, options)

	// When
	err := store.Set(ctx, cacheKey, cacheValue, nil)

	// Then
	assert.Nil(t, err)
}

func TestRedisClusterSetWithTags(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	cacheKey := "my-key"
	cacheValue := "my-cache-value"

	client := mocksStore.NewMockRedisClusterClientInterface(ctrl)
	client.EXPECT().Set(ctx, cacheKey, cacheValue, time.Duration(0)).Return(&redis.StatusCmd{})
	client.EXPECT().SAdd(ctx, "gocache_tag_tag1", "my-key").Return(&redis.IntCmd{})
	client.EXPECT().Expire(ctx, "gocache_tag_tag1", 720*time.Hour).Return(&redis.BoolCmd{})

	store := NewRedisCluster(client, nil)

	// When
	err := store.Set(ctx, cacheKey, cacheValue, &Options{Tags: []string{"tag1"}})

	// Then
	assert.Nil(t, err)
}

func TestRedisClusterDelete(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	cacheKey := "my-key"

	client := mocksStore.NewMockRedisClusterClientInterface(ctrl)
	client.EXPECT().Del(ctx, "my-key").Return(&redis.IntCmd{})

	store := NewRedisCluster(client, nil)

	// When
	err := store.Delete(ctx, cacheKey)

	// Then
	assert.Nil(t, err)
}

func TestRedisClusterInvalidate(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	options := InvalidateOptions{
		Tags: []string{"tag1"},
	}

	cacheKeys := &redis.StringSliceCmd{}

	client := mocksStore.NewMockRedisClusterClientInterface(ctrl)
	client.EXPECT().SMembers(ctx, "gocache_tag_tag1").Return(cacheKeys)
	client.EXPECT().Del(ctx, "gocache_tag_tag1").Return(&redis.IntCmd{})

	store := NewRedisCluster(client, nil)

	// When
	err := store.Invalidate(ctx, options)

	// Then
	assert.Nil(t, err)
}

func TestRedisClusterClear(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	client := mocksStore.NewMockRedisClusterClientInterface(ctrl)
	client.EXPECT().FlushAll(ctx).Return(&redis.StatusCmd{})

	store := NewRedisCluster(client, nil)

	// When
	err := store.Clear(ctx)

	// Then
	assert.Nil(t, err)
}

func TestRedisClusterGetType(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	client := mocksStore.NewMockRedisClusterClientInterface(ctrl)

	store := NewRedisCluster(client, nil)

	// When - Then
	assert.Equal(t, RedisClusterType, store.GetType())
}
