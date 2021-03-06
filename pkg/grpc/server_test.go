package grpc

import (
	"sync"
	"testing"
	"time"

	"github.com/dapr/dapr/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestCertRenewal(t *testing.T) {
	t.Run("shouldn't renew", func(t *testing.T) {
		certExpiry := time.Now().Add(time.Hour * 2).UTC()
		certDuration := certExpiry.Sub(time.Now().UTC())

		renew := shouldRenewCert(certExpiry, certDuration)
		assert.False(t, renew)
	})

	t.Run("should renew", func(t *testing.T) {
		certExpiry := time.Now().Add(time.Second * 3).UTC()
		certDuration := certExpiry.Sub(time.Now().UTC())

		time.Sleep(time.Millisecond * 2200)
		renew := shouldRenewCert(certExpiry, certDuration)
		assert.True(t, renew)
	})
}

func TestGetMiddlewareOptions(t *testing.T) {
	t.Run("should enable two interceptors if tracing and metrics are enabled", func(t *testing.T) {
		fakeServer := &server{
			config: ServerConfig{
				EnableMetrics: false,
			},
			tracingSpec: config.TracingSpec{
				Enabled: true,
			},
			renewMutex: &sync.Mutex{},
		}

		serverOption := fakeServer.getMiddlewareOptions()

		assert.Equal(t, 2, len(serverOption))
	})

	t.Run("should disable middlewares", func(t *testing.T) {
		fakeServer := &server{
			config: ServerConfig{
				EnableMetrics: false,
			},
			tracingSpec: config.TracingSpec{
				Enabled: false,
			},
			renewMutex: &sync.Mutex{},
		}

		serverOption := fakeServer.getMiddlewareOptions()

		assert.Equal(t, 0, len(serverOption))
	})
}
