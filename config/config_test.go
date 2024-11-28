package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("CONF:Not nil", func(t *testing.T) {
		t.Setenv("CONFIG_PATH", "./")
		t.Setenv("CONFIG_NAME", "config")
		t.Setenv("CONFIG_EXT", "yaml")
		conf := GetConfig()
		assert.NotNil(t, conf, "Config should not be nil")
	})

	t.Run("CONF:Same conf", func(t *testing.T) {
		t.Setenv("CONFIG_PATH", "./")
		t.Setenv("CONFIG_NAME", "config")
		t.Setenv("CONFIG_EXT", "yaml")
		conf1 := GetConfig()
		conf2 := GetConfig()
		assert.Equal(t, conf1, conf2)
	})
}
