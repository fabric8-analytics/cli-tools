package telemetry

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"os/user"
	"testing"
)

func TestSetError(t *testing.T) {
	err := errors.New("this is an error string")
	assert.Equal(t, err.Error(), SetError(err))

	currentUser, err := user.Current()
	assert.NoError(t, err)
	configHome, _ := homedir.Dir()
	err = fmt.Errorf("cannot access storage file '%s/.crc/machines/crc/crc.qcow2' (as uid:64055, gid:129): Permission denied')", configHome)
	assert.NotEqual(t, err.Error(), SetError(err))
	assert.NotContains(t, SetError(err), configHome)

	err = fmt.Errorf("user %s may not use sudo", currentUser.Username)
	assert.NotEqual(t, err.Error(), SetError(err))
	assert.NotContains(t, SetError(err), currentUser.Username)
}
