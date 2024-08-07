package config_test

import (
	"os"
	"testing"

	"micros/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	dir := os.Getenv("GOPATH") + "/src/micros/"
	conf, err := config.New(dir+"configs/conf.ini", dir+"configs/conf.yaml")
	assert.Nil(t, err)
	assert.Equal(t, "mysql", conf.DB.DBType)
	assert.Equal(t, "DEBUG", conf.Log.Level)
	assert.Equal(t, "true", conf.Opt.Profile)
	assert.Equal(t, "micros", conf.Project)
}
