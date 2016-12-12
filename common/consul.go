package common

import (
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/jacobstr/confer"
)

// NewConsulFromConfig initializes Consul client based on config + environment overrides
func NewConsulFromConfig(config *confer.Config, configPrefix string) (client *api.Client, consulPrefix string, err error) {
	upperPrefix := strings.ToUpper(configPrefix)

	config.BindEnv(upperPrefix+"_ENABLED", configPrefix+".enabled")
	config.BindEnv(upperPrefix+"_SCHEME", configPrefix+".scheme")
	config.BindEnv(upperPrefix+"_ADDRESS", configPrefix+".address")
	config.BindEnv(upperPrefix+"_PREFIX", configPrefix+".prefix")
	config.BindEnv(upperPrefix+"_TOKEN", configPrefix+".token")

	if !config.GetBool(configPrefix + ".enabled") {
		// Consul is disabled
		return
	}

	conf := api.DefaultConfig()
	conf.Scheme = config.GetString(configPrefix + ".scheme")
	conf.Address = config.GetString(configPrefix + ".address")
	conf.Token = config.GetString(configPrefix + ".token")

	consulPrefix = config.GetString(configPrefix + ".prefix")
	client, err = api.NewClient(conf)

	return
}
