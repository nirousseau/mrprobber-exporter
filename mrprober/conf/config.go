package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"sync"
)

var (
	// SafeConfiguration is the unique configuration of the binary.
	// There are 2 ways to update this : either you update the file on
	// disk, either you change this programmatically.
	SafeConfiguration = SafeConfig{
		cfg: &Config{},
	}
)

type Config struct {
	Global Global   `yaml:"global"`
	Rules  []Rule   `yaml:"rules"`
	Alerts []string `yaml:"alerts"`
}

type Global struct {
	Tickrate int `yaml:"tickrate"`
	Web      Web `yaml:"web"`
}

type Web struct {
	ListenAddress string `yaml:"listenAddress"`
	MetricsPath   string `yaml:"metricsPath"`
}

type Rule struct {
	Name  string   `yaml:"name"`
	Probe string   `yaml:"probe"`
	Args  []string `yaml:"args"`
}

// SafeConfig is a wrapper around Config for concurrency access
type SafeConfig struct {
	mu  sync.RWMutex
	cfg *Config
}

// Update reloads the configuration.
func (sc *SafeConfig) Update() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	config, err := readConfig()
	sc.cfg = config
	if err != nil {
		return err
	}
	return nil
}

// Get returns a copy of the current configuration.
func (sc *SafeConfig) Get() Config {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return *sc.cfg
}

// readConfig reads and transforms the config file to a struct
// and performs a dummy schema checking
func readConfig() (*Config, error) {

	c := Config{}
	var err error

	if err = viper.ReadInConfig(); err != nil {
		log.Fatal(fmt.Errorf("can't read config : %w", err))
	}

	if err = viper.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration. %+v", err)
	}

	// Add metadata as a fake probe
	hookRules(&c)

	// Sanity check
	err = check(c)
	if err != nil {
		return nil, fmt.Errorf("failed to validate configuration: %+v", err)
	}

	return &c, nil
}

// hookRules add a metadata rule that is not explicitly set in the configuration file.
// The intention behind this is to add a "meta" metric for dashboard designers.
func hookRules(c *Config) {

	r := Rule{
		Name:  "Fake rule for self-testing and metadata information.",
		Probe: "fake",
		Args:  nil,
	}

	c.Rules = append(c.Rules, r)
}

// check analyze the configuration for errors
// raise an error when encountered
func check(c Config) error {

	var err error

	err = verifyDuplicatedRules(c)
	if err != nil {
		return err
	}

	return nil
}

// verifyDuplicatedRules check that rules from the configuration have unique names
func verifyDuplicatedRules(c Config) error {

	// https://stackoverflow.com/questions/65258003/memory-allocation-of-mapintinterface-vs-mapintstruct
	rulesRegistry := make(map[string]struct{}, len(c.Rules))
	for _, rule := range c.Rules {
		if _, hasRule := rulesRegistry[rule.Name]; hasRule {
			return fmt.Errorf("duplicated name identifier for `%s` rule", rule.Name)
		}
		rulesRegistry[rule.Name] = struct{}{}
	}

	return nil
}
