package storagetemporary_test

import (
	"testing"
	"time"

	"github.com/av-belyakov/enricher_geoip/internal/storagetemporary"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	ctx := t.Context()

	storage, err := storagetemporary.New(
		ctx,
		storagetemporary.WithCacheTimeTick(2),
		storagetemporary.WithCacheTTL(3),
	)
	assert.NoError(t, err)

	storage.SetData("1", storagetemporary.IncomingRequest{
		TaskId: time.Now().Format("20060102-150405.000000"),
		IP:     []string{"213.13.3.64", "12.66.45.6"},
	})
	storage.SetData("2", storagetemporary.IncomingRequest{
		TaskId: time.Now().Format("20060102-150405.000000"),
		IP:     []string{"114.33.12.164", "132.66.45.16"},
	})

	time.Sleep(1 * time.Second)
	assert.Equal(t, storage.DataSize(), 2)

	storage.SetData("3", storagetemporary.IncomingRequest{
		TaskId: time.Now().Format("20060102-150405.000000"),
		IP:     []string{"33.43.15.224", "12.6.0.160"},
	})

	time.Sleep(1 * time.Second)
	assert.Equal(t, storage.DataSize(), 3)

	storage, err = storagetemporary.New(
		ctx,
		storagetemporary.WithCacheTimeTick(2),
		storagetemporary.WithCacheTTL(3),
	)
	assert.NoError(t, err)
	assert.Equal(t, storage.DataSize(), 3)

	value, ok := storage.GetData("2")
	assert.True(t, ok)
	assert.Equal(t, len(value.IP), 2)
}
