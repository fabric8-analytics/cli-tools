package segment

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fabric8-analytics/cli-tools/pkg/config"

	"github.com/rs/zerolog/log"

	"github.com/fabric8-analytics/cli-tools/pkg/telemetry"
	"github.com/fabric8-analytics/cli-tools/pkg/version"
	"github.com/mitchellh/go-homedir"
	"github.com/pborman/uuid"
	"github.com/segmentio/analytics-go"
)

var WriteKey = "MW6rAYP7Q6AAiSAZ3Ussk6eMebbVcchD" // test

type Client struct {
	segmentClient     analytics.Client
	telemetryFilePath string
}

func NewClient() (*Client, error) {
	configHome, _ := homedir.Dir()
	return newCustomClient(
		filepath.Join(configHome, ".redhat", "anonymousId"),
		analytics.DefaultEndpoint)
}

func newCustomClient(telemetryFilePath, segmentEndpoint string) (*Client, error) {
	client, err := analytics.NewWithConfig(WriteKey, analytics.Config{
		Endpoint: segmentEndpoint,
		DefaultContext: &analytics.Context{
			App: analytics.AppInfo{
				Name:    "crda cli",
				Version: version.GetCRDAVersion(),
			},
		},
		RetryAfter: func(attempt int) time.Duration {
			return time.Duration(attempt * 3)
		},
		BatchSize: 1,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		segmentClient:     client,
		telemetryFilePath: telemetryFilePath,
	}, nil
}

func (c *Client) Close() error {
	return c.segmentClient.Close()
}

func (c *Client) Upload(ctx context.Context, action string, duration time.Duration, err error) error {
	if config.ActiveConfigValues.ConsentTelemetry != "1" {
		return nil
	}
	log.Debug().Msgf("Sending Segment Info")
	userID, uerr := getUserIdentity(c.telemetryFilePath)
	if uerr != nil {
		return uerr
	}

	if err := c.segmentClient.Enqueue(analytics.Identify{
		UserId: userID,
	}); err != nil {
		return err
	}

	properties := analytics.NewProperties()
	for k, v := range telemetry.GetContextProperties(ctx) {
		properties = properties.Set(k, v)
	}

	properties = properties.
		Set("success", err == nil).
		Set("platform", runtime.GOOS).
		Set("duration", duration.Milliseconds())

	if err != nil {
		properties = properties.
			Set("error", telemetry.SetError(err)).
			Set("error-type", errorType(err))
	}

	return c.segmentClient.Enqueue(analytics.Track{
		UserId:     userID,
		Event:      action,
		Properties: properties,
	})
}

func getUserIdentity(telemetryFilePath string) (string, error) {
	var id []byte
	if err := os.MkdirAll(filepath.Dir(telemetryFilePath), 0750); err != nil {
		return "", err
	}
	if _, err := os.Stat(telemetryFilePath); !os.IsNotExist(err) {
		id, err = ioutil.ReadFile(telemetryFilePath)
		if err != nil {
			return "", err
		}
	}
	if uuid.Parse(strings.TrimSpace(string(id))) == nil {
		id = []byte(uuid.NewRandom().String())
		if err := ioutil.WriteFile(telemetryFilePath, id, 0600); err != nil {
			return "", err
		}
	}
	return strings.TrimSpace(string(id)), nil
}

func errorType(err error) string {
	wrappedErr := errors.Unwrap(err)
	if wrappedErr != nil {
		return fmt.Sprintf("%T", wrappedErr)
	}
	return fmt.Sprintf("%T", err)
}
