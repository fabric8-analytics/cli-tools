package segment

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	//cmdConfig "github.com/code-ready/crc/cmd/crc/cmd/config"
	//crcConfig "github.com/code-ready/crc/pkg/crc/config"
	//crcErr "github.com/code-ready/crc/pkg/crc/errors"
	"github.com/fabric8-analytics/cli-tools/pkg/telemetry"
	"github.com/stretchr/testify/require"
)

func mockServer() (chan []byte, *httptest.Server) {
	done := make(chan []byte, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		bin, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error().Msgf(err.Error())
			return
		}
		done <- bin
	}))

	return done, server
}

func TestClientUploadWithContext(t *testing.T) {
	body, server := mockServer()
	defer server.Close()
	defer close(body)

	dir, err := ioutil.TempDir("", "cfg")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	c, err := newCustomClient(filepath.Join(dir, "telemetry"), server.URL)
	require.NoError(t, err)

	ctx := telemetry.NewContext(context.Background())
	require.NoError(t, c.Upload(ctx, "start", time.Minute, nil))
	require.NoError(t, c.Close())
}

func TestClientUploadWithOutConsent(t *testing.T) {
	body, server := mockServer()
	defer server.Close()
	defer close(body)

	dir, err := ioutil.TempDir("", "cfg")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	c, err := newCustomClient(filepath.Join(dir, "telemetry"), server.URL)
	require.NoError(t, err)

	require.NoError(t, c.Upload(context.Background(), "start", time.Second, errors.New("an error occurred")))
	require.NoError(t, c.Close())

	select {
	case <-body:
		require.Fail(t, "server should not receive data")
	default:
	}
}
