package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Contains all the keys for Wash's config
const (
	SocketKey = "socket"
)

type config struct {
	Socket string
}

// Fields contains the fields of Wash's config
var Fields config

// Load Wash's config.
func Load() error {
	// Set any defaults
	viper.SetDefault(SocketKey, "/tmp/wash-api.sock")

	// Tell viper that the config. can be read from WASH_<entry>
	// environment variables
	viper.SetEnvPrefix("WASH")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// TODO: Add any additional config files, then make sure to
	// invoke viper.ReadInConfig() to read-in their values

	// Load the config
	Fields = config{}
	Fields.Socket = viper.GetString(SocketKey)

	return nil
}