// Copyright © 2017 Michael Shindle <mshindle@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/mshindle/smithy/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/mshindle/smithy/data"
)

const defaultBaseDir = "$HOME/.smithy"
const defaultSystemDir = "/etc/smithy"

var cfgFile string
var argAsString bool
var processor data.Processor

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "smithy",
	Short: "Provide per-value encryption of sensitive data with JSON files",
	Long: `
smithy provides a way for sensitive data to be encrypted and stored
in a secure way. Instead of encrypting an entire file and losing the 
ability to modify non-sensitive data unless sensitive data is exposed, 
smithy will detect if a field is encrypted and decrypt as appropriate.`,

	// ensure that the base dir exists
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// set logging level
		config.UpdateLogging()

		// create our basedir if not existent
		config.CreateBaseDir()
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	// initialize logging
	log.SetOutput(os.Stderr)
	log.SetLevel(log.WarnLevel)

	// initialize cobra
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.smithy/smithy.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("smithy")         // name of config file (without extension)
	viper.AddConfigPath(defaultBaseDir)   // adding home directory as first search path
	viper.AddConfigPath(defaultSystemDir) // adding system directory as second search path
	viper.AutomaticEnv()                  // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fileUsed := viper.ConfigFileUsed()
		log.WithField("file", fileUsed).Info("using config file")
		viper.SetDefault("baseDir", filepath.Dir(fileUsed))
	} else {
		viper.SetDefault("baseDir", defaultBaseDir)
	}

	// set configuration defaults
	viper.SetDefault("encryptMethod", "rsa")
	viper.SetDefault("publicKey", "public.key")
	viper.SetDefault("privateKey", "private.key")
	viper.SetDefault("logging.level", "warn")

	// load into settings
	config.Initialize()
}
