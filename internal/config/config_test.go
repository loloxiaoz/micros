package config_test

import (
	"os"
	"testing"

	"micros/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	var conf config.Conf

	dir := os.Getenv("GOPATH") + "/src/micros/"
	conf.Init(dir + "configs/conf.ini", dir + "configs/conf.yaml")
	assert.Equal(t, "mysql", conf.DB.DBType)
	assert.Equal(t, "DEBUG", conf.Log.Level)
	assert.Equal(t, "micros", conf.Project)
}