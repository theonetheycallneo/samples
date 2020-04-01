// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	AuthToken string        `config:"token"`
	MaxEvents int           `config:"max_events"`
	Period    time.Duration `config:"period"`
}

var DefaultConfig = Config{
	AuthToken: "",
	MaxEvents: 10000,
	Period:    1 * time.Second,
}
