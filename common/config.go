package common

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/jacobstr/confer"
)

// LoadConfig loads default config.yaml file
func LoadConfig() *confer.Config {
	return LoadConfigs("config.yaml")
}

// LoadGlobConfigs load configs by pattern
func LoadGlobConfigs(pattern string) *confer.Config {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalf("unable to glob config files: %s", err)
	}

	return LoadConfigs(matches...)
}

// LoadConfigs loads configs from the list of files
func LoadConfigs(configFiles ...string) *confer.Config {
	config := confer.NewConfig()

	log.Printf("Loading config from %v...", configFiles)
	err := config.ReadPaths(configFiles...)
	if err != nil {
		log.Fatalf("unable to read configuration: %s", err)
	}

	return config
}

// LoadConsulOverrides loads configuration from Consul (if any) and merges with current configuration
func LoadConsulOverrides(client *api.Client, prefix string) (map[string]interface{}, error) {
	kv := client.KV()

	prefix = strings.TrimPrefix(prefix, "/")

	var (
		pairs api.KVPairs
		err   error
	)

	for attempt := 0; attempt < 180; attempt++ {
		pairs, _, err = kv.List(prefix, nil)
		if err == nil {
			break
		}

		log.Printf("Failure loading config from Consul: %s, retrying in 1 second, attempt %d", err, attempt)
		time.Sleep(time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failure listing keys under %#v: %s", prefix, err)
	}

	configParts := map[string]interface{}{}

	for _, pair := range pairs {
		var val interface{}

		err = json.Unmarshal(pair.Value, &val)
		if err != nil {
			return nil, fmt.Errorf("unable to decode JSON for key %#v: %s", pair.Key, err)
		}

		configParts[pair.Key] = val
	}

	return configParts, nil
}

// MergeConfigOverrides merges current config with overrides loaded from other sources
func MergeConfigOverrides(config *confer.Config, configParts map[string]interface{}) error {
	keys := make([]string, len(configParts))
	i := 0
	for key := range configParts {
		keys[i] = key
		i++
	}

	sort.Strings(keys)

	for _, key := range keys {
		log.Printf("Merging overrides from key %s...", key)

		err := config.MergeAttributes(configParts[key])
		if err != nil {
			return fmt.Errorf("unable to merge data for key %s (data %#v): %s", key, configParts[key], err)
		}
	}

	return nil
}

// MergeConfigFromConsul loads and merges config from Consul
func MergeConfigFromConsul(config *confer.Config) {
	consulClient, prefix, err := NewConsulFromConfig(config, "consul")
	if err != nil {
		log.Fatalf("Failed building Consul client: %s", err)
	}

	if consulClient == nil {
		// Consul is disabled, nothing to do
		return
	}

	configParts, err := LoadConsulOverrides(consulClient, prefix)
	if err != nil {
		log.Fatalf("Failed loading overrides from Consul: %s", err)
	}

	err = MergeConfigOverrides(config, configParts)
	if err != nil {
		log.Fatalf("Failed applying overrides from Consul: %s", err)
	}
}
