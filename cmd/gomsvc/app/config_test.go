package app_test

import (
	"strings"
	"testing"

	"github.com/inquizarus/gomsvc/cmd/gomsvc/app"
	"github.com/stretchr/testify/assert"
)

func TestThatAddressWorksAsIntended(t *testing.T) {
	configWithPortSpecified := app.Config{
		Port: "8081",
	}
	configWithoutPortSpecified := app.Config{}

	assert.Equal(t, ":8081", configWithPortSpecified.Address())
	assert.Equal(t, ":8080", configWithoutPortSpecified.Address())
}

func TestThatLoadConfigFromPathWorksAsIntended(t *testing.T) {
	path := "./testdata/config.fixture.json"
	config, err := app.ConfigFromFilePath(path)

	assert.NoError(t, err)

	assert.Equal(t, "8081", config.Port)
}

func TestThatLoadConfigFromReaderWorksAsIntended(t *testing.T) {
	dataString := `{"port":"8081"}`
	config, err := app.ConfigFromReader(strings.NewReader(dataString))

	assert.NoError(t, err)

	assert.Equal(t, "8081", config.Port)
}

func TestThatLoadConfigFromReaderWorksAsIntendedWhenReaderIsNil(t *testing.T) {
	_, err := app.ConfigFromReader(nil)
	assert.Error(t, err)
}
