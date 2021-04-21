package telemetry

import (
	"context"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"os/user"
	"strings"
	"sync"
)

type contextKey struct{}

var key = contextKey{}

// Properties maintain state of Telemetry event Properties.
type Properties struct {
	lock    sync.Mutex
	storage map[string]interface{}
}

// GetTelemetryConsent fires telemetry consent popup.
func GetTelemetryConsent() bool {
	prompt := promptui.Prompt{
		Label:       "Would you like to contribute anonymous usage statistics [y/n]",
		HideEntered: true,
	}
	userInput, err := prompt.Run()

	if err != nil {
		log.Fatal().Msgf(fmt.Sprintf("Prompt failed %v\n", err))
	}

	userResponse := Find(userInput)

	return userResponse
}

// Find compared user input string with accepted values
func Find(val string) bool {
	yes := []string{"y", "Y", "1"}
	no := []string{"n", "N", "0"}

	for _, item := range yes {
		if item == val {
			return true
		}
	}
	for _, item := range no {
		if item == val {
			return false
		}
	}
	return GetTelemetryConsent()
}

func (p *Properties) set(name string, value interface{}) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.storage[name] = value
}

func (p *Properties) values() map[string]interface{} {
	p.lock.Lock()
	defer p.lock.Unlock()
	ret := make(map[string]interface{})
	for k, v := range p.storage {
		ret[k] = v
	}
	return ret
}

func propertiesFromContext(ctx context.Context) *Properties {
	value := ctx.Value(key)
	if cast, ok := value.(*Properties); ok {
		return cast
	}
	return nil
}

// NewContext creates a New Context
func NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, &Properties{storage: make(map[string]interface{})})
}

// GetContextProperties returns current property state
func GetContextProperties(ctx context.Context) map[string]interface{} {
	properties := propertiesFromContext(ctx)
	if properties == nil {
		return make(map[string]interface{})
	}
	return properties.values()
}

func setContextProperty(ctx context.Context, key string, value interface{}) {
	properties := propertiesFromContext(ctx)
	if properties != nil {
		properties.set(key, value)
	}
}

// SetError replaces sensitive data from error recording
func SetError(err error) string {
	// Mask username if present in the error string
	currentUser, err1 := user.Current()
	if err1 != nil {
		return err1.Error()
	}
	configHome, _ := homedir.Dir()
	withoutHomeDir := strings.ReplaceAll(err.Error(), configHome, "$HOME")
	return strings.ReplaceAll(withoutHomeDir, currentUser.Username, "$USERNAME")
}

// SetFlag records Json, verbose flags
func SetFlag(ctx context.Context, flag string, value bool) {
	setContextProperty(ctx, flag, value)
}

// SetManifest sets manifest name property
func SetManifest(ctx context.Context, value string) {
	setContextProperty(ctx, "manifest", value)
}

// SetExitCode sets exit code property
func SetExitCode(ctx context.Context, value int) {
	setContextProperty(ctx, "exit-code", value)
}

// SetClient sets cient property: Ex: terminal, jenkins, etc
func SetClient(ctx context.Context, value string) {
	setContextProperty(ctx, "client", value)
}

// SetVulnerability sets total vulnerability found property
func SetVulnerability(ctx context.Context, value int) {
	setContextProperty(ctx, "total-vulnerabilities", value)
}

// SetSnykTokenAssociation sets synk-token-associated property
func SetSnykTokenAssociation(ctx context.Context, value bool) {
	setContextProperty(ctx, "snyk-token-associated", value)
}
