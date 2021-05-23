package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewApp(t *testing.T) {
	_, createErr := NewApp()
	assert.NoError(t, createErr)
}
