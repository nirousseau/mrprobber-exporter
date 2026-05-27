package cmd

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"mrprober/conf"
	"mrprober/probes"
)

const defaultConfigurationDirectory = "/etc/mrprober/"
const defaultConfigurationFile = "mrprober.yaml"

var (
	rootCmd = &cobra.Command{
		Use:   "mrprober",
		Short: "Probe scheduler and reporter for K8S.",
		Long: `Mr Prober runs system probes and exposes the results on a Prometheus endpoint.
This program is meant to be run as a DaemonSet on all nodes in Kubernetes/Openshift clusters.
You can also trigger a "one-shot" run.

Maintainer : Nikita ROUSSEAU
License : MIT License`,
	}

	cfgFile string
)

func Execute(appVersion string) {
	rootCmd.Version = appVersion
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", fmt.Sprintf("config file (default is %s%s)", defaultConfigurationDirectory, defaultConfigurationFile))

	// Subscribe commands to root
	rootCmd.AddCommand(
		newRunCommand(),
		newDaemonCommand(),
	)
}

func initConfig() {

	// Read only yaml files
	viper.SetConfigType("yaml")

	// Read config either from cfgFile or from configuration directory
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Read from configuration repository
		viper.AddConfigPath(defaultConfigurationDirectory)
		viper.SetConfigName(defaultConfigurationFile)
	}

	err := conf.SafeConfiguration.Update()
	if err != nil {
		log.Fatal(err)
	}

	// Add watcher
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Reloaded configuration:", e.Name)

		err := conf.SafeConfiguration.Update()
		if err != nil {
			log.Fatal(err)
		}

		probes.UnregisterAllMetrics()
	})
	viper.WatchConfig()

	// Notify metadata probe about the current version
	probes.VERSION = rootCmd.Version
}
