package store

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	mocksStore "github.com/wlxwlxwlx/gocache/v2/test/mocks/store/clients"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewRistretto(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	client := mocksStore.NewMockRistrettoClientInterface(ctrl)
	options := &Options{
		Cost: 8,
	}

	// When
	store := NewRistretto(client, options)

	// Then
	assert.IsType(t, new(RistrettoStore), store)
	assert.Equal(t, client, store.client)
	assert.Equal(t, options, store.options)
}

func TestRistrettoGet(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	cacheKey := "my-key"
	cacheValue := "my-cache-value"

	client := mocksStore.NewMockRistrettoClientInterface(ctrl)
	client.EXPECT().Get(cacheKey).Return(cacheValue, true)

	store := NewRistretto(client, nil)

	// When
	value, err := store.Get(ctx, cacheKey)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, cacheValue, value)
}

func TestRistrettoGetWhenError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	cacheKey := "my-key"

	client := mocksStore.NewMockRistrettoClientInterface(ctrl)
	client.EXPECT().Get(cacheKey).Return(nil, false)

	store := NewRistretto(client, nil)

	// When
	value, err := store.Get(ctx, cacheKey)

	// Then
	assert.Nil(t, value)
	assert.Equal(t, errors.New("Value not found in Ristretto store"), err)
}

func TestRistrettoGetWithTTL(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	cacheKey := "my-key"
	cacheValue := "my-cache-value"

	client := mocksStore.NewMockRistrettoClientInterface(ctrl)
	client.EXPECT().Get(cacheKey).Return(cacheValue, true)

	store := NewRistretto(client, nil)

	// When
	value, ttl, err := store.GetWithTTL(ctx, cacheKey)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, cacheValue, value)
	assert.Equal(t, 0*time.Second, ttl)
}

func TestRistrettoGetWithTTLWhenError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	cacheKey := "my-key"

	client := mocksStore.NewMockRistrettoClientInterface(ctrl)
	client.EXPECT().Get(cacheKey).Return(nil, false)

	store := NewRistretto(client, nil)

	// When
	value, ttl, err := store.GetWithTTL(ctx, cacheKey)

	// Then
	assert.Nil(t, value)
	assert.Equal(t, errors.New("Value not found in Ristretto store"), err)
	assert.Equal(t, 0*time.Second, ttl)
}

func TestRistrettoSet(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	cacheKey := "my-key"
	cacheValue := "my-cache-value"
	options := &Options{
		Cost: 7,
	}

	client := mocksStore.NewMockRistrettoClientInterface(ctrl)
	client.EXPECT().SetWithTTL(cacheKey, cacheValue, int64(4), 0*time.Second).Return(true)

	store := NewRistretto(client, options)

	// When
	err := store.Set(ctx, cacheKey, cacheValue, &Options{
		Cost: 4,
	})

	// Then
	assert.Nil(t, err)
}

func TestRistrettoSetWhenNoOptionsGiven(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	cacheKey := "my-key"
	cacheValue := "my-cache-value"
	options := &Options{
		Cost: 7,
	}

	client := mocksStore.NewMockRistrettoClientInterface(ctrl)
	client.EXPECT().SetWithTTL(cacheKey, cacheValue, int64(7), 0*time.Second).Return(true)

	store := NewRistretto(client, options)

	// When
	err := store.Set(ctx, cacheKey, cacheValue, nil)

	// Then
	assert.Nil(t, err)
}

func TestRistrettoSetWhenError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	cacheKey := "my-key"
	cacheValue := "my-cache-value"
	options := &Options{
		Cost: 7,
	}

	client := mocksStore.NewMockRistrettoClientInterface(ctrl)
	client.EXPECT().SetWithTTL(cacheKey, cacheValue, int64(7), 0*time.Second).Return(false)

	store := NewRistretto(client, options)

	// When
	err := store.Set(ctx, cacheKey, cacheValue, nil)

	// Then
	assert.Equal(t, fmt.Errorf("An error has occurred while setting value '%v' on key '%v'", cacheValue, cacheKey), err)
}

func TestRistrettoSetWithTags(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	cacheKey := "my-key"
	cacheValue := []byte("my-cache-value")

	client := mocksStore.NewMockRistrettoClientInterface(ctrl)
	client.EXPECT().SetWithTTL(cacheKey, cacheValue, int64(0), 0*time.Second).Return(true)
	client.EXPECT().Get("gocache_tag_tag1").Return(nil, true)
	client.EXPECT().SetWithTTL("gocache_tag_tag1", []byte("my-key"), int64(0), 720*time.Hour).Return(true)

	store := NewRistretto(client, nil)

	// When
	err := store.Set(ctx, cacheKey, cacheValue, &Options{Tags: []string{"tag1"}})

	// Then
	assert.Nil(t, err)
}

func TestRistrettoSetWithTagsWhenAlreadyInserted(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	cacheKey := "my-key"
	cacheValue := []byte("my-cache-value")

	client := mocksStore.NewMockRistrettoClientInterface(ctrl)
	client.EXPECT().SetWithTTL(cacheKey, cacheValue, int64(0), 0*time.Second).Return(true)
	client.EXPECT().Get("gocache_tag_tag1").Return([]byte("my-key,a-second-key"), true)
	client.EXPECT().SetWithTTL("gocache_tag_tag1", []byte("my-key,a-second-key"), int64(0), 720*time.Hour).Return(true)

	store := NewRistretto(client, nil)

	// When
	err := store.Set(ctx, cacheKey, cacheValue, &Options{Tags: []string{"tag1"}})

	// Then
	assert.Nil(t, err)
}

func TestRistrettoDelete(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	cacheKey := "my-key"

	client := mocksStore.NewMockRistrettoClientInterface(ctrl)
	client.EXPECT().Del(cacheKey)

	store := NewRistretto(client, nil)

	// When
	err := store.Delete(ctx, cacheKey)

	// Then
	assert.Nil(t, err)
}

func TestRistrettoInvalidate(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	options := InvalidateOptions{
		Tags: []string{"tag1"},
	}

	cacheKeys := []byte("a23fdf987h2svc23,jHG2372x38hf74")

	client := mocksStore.NewMockRistrettoClientInterface(ctrl)
	client.EXPECT().Get("gocache_tag_tag1").Return(cacheKeys, true)
	client.EXPECT().Del("a23fdf987h2svc23")
	client.EXPECT().Del("jHG2372x38hf74")
	client.EXPECT().Del("gocache_tag_tag1")

	store := NewRistretto(client, nil)

	// When
	err := store.Invalidate(ctx, options)

	// Then
	assert.Nil(t, err)
}

func TestRistrettoInvalidateWhenError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	options := InvalidateOptions{
		Tags: []string{"tag1"},
	}

	cacheKeys := []byte("a23fdf987h2svc23,jHG2372x38hf74")

	client := mocksStore.NewMockRistrettoClientInterface(ctrl)
	client.EXPECT().Get("gocache_tag_tag1").Return(cacheKeys, false)

	store := NewRistretto(client, nil)

	// When
	err := store.Invalidate(ctx, options)

	// Then
	assert.Nil(t, err)
}

func TestRistrettoClear(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	client := mocksStore.NewMockRistrettoClientInterface(ctrl)
	client.EXPECT().Clear()

	store := NewRistretto(client, nil)

	// When
	err := store.Clear(ctx)

	// Then
	assert.Nil(t, err)
}

func TestRistrettoGetType(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	client := mocksStore.NewMockRistrettoClientInterface(ctrl)

	store := NewRistretto(client, nil)

	// When - Then
	assert.Equal(t, RistrettoType, store.GetType())
}
