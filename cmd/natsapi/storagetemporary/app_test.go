package storagetemporary_test

import (
	"testing"
	"time"

	"github.com/av-belyakov/enricher_geoip/cmd/natsapi/storagetemporary"
	"github.com/nats-io/nats.go"
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

	storage.SetReq("1", &nats.Msg{})
	storage.SetReq("2", &nats.Msg{})

	time.Sleep(1 * time.Second)
	assert.Equal(t, storage.DataSize(), 2)

	storage.SetReq("3", &nats.Msg{})

	time.Sleep(1 * time.Second)
	assert.Equal(t, storage.DataSize(), 3)

	storage, err = storagetemporary.New(
		ctx,
		storagetemporary.WithCacheTimeTick(2),
		storagetemporary.WithCacheTTL(3),
	)
	assert.NoError(t, err)
	assert.Equal(t, storage.DataSize(), 3)

	_, ok := storage.GetReq("2")
	assert.True(t, ok)
}
