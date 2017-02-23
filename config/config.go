package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// LogSettings holds the settings for logging
type LogSettings struct {
	Level string `yaml:"level"`
}

// Settings holds the global settings
type Settings struct {
	BaseDir       string      `yaml:"baseDir"`
	EncryptMethod string      `yaml:"encryptMethod"`
	PublicKey     string      `yaml:"publicKey"`
	PrivateKey    string      `yaml:"privateKey"`
	Logging       LogSettings `yaml:"logging"`
}

var config Settings

// Initialize our configuration by binding with viper
func Initialize() {
	err := viper.Unmarshal(&config)
	if err != nil {
		log.WithError(err).Fatal("unable to decode into struct")
		return
	}
}

// CreateBaseDir creates the base directory in the configuration
func CreateBaseDir() error {
	path := expandPath(config.BaseDir)
	logger := log.WithFields(log.Fields{"base_dir": config.BaseDir, "path": path})
	if _, err := os.Stat(path); err != nil {

		logger.Info("base_dir does not exist. creating.")
		err = os.MkdirAll(path, 0700)
		if err != nil {
			logger.Error("cannot create base_dir.")
			return err
		}
	} else {
		logger.Debug("base_dir exists. skipping.")
	}
	return nil
}

// PublicKey returns the absolute path to the public key
func PublicKey() string {
	return absPathToKey(config.PublicKey)
}

// PrivateKey returns the absolute path to the private key
func PrivateKey() string {
	return absPathToKey(config.PrivateKey)
}

func absPathToKey(key string) string {
	if filepath.IsAbs(key) {
		return key
	}
	return filepath.Join(expandPath(config.BaseDir), key)
}

// expandPath runs a crude shell expansion on the given path
func expandPath(path string) string {
	log.WithField("path", path).Debug("trying to resolve absolute path")

	if strings.HasPrefix(path, "$HOME") {
		path = UserHomeDir() + path[5:]
	}

	if strings.HasPrefix(path, "$") {
		end := strings.Index(path, string(os.PathSeparator))
		path = os.Getenv(path[1:end]) + path[end:]
	}

	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}

	p, err := filepath.Abs(path)
	if err == nil {
		return filepath.Clean(p)
	}

	log.WithError(err).Error("couldn't discover absolute path")
	return ""
}

// UserHomeDir returns the home directory for windows & *nix
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// Print dumps the current configuration to stdout
func Print() {
	log.Info("marshaling output as yaml")
	out, err := yaml.Marshal(&config)
	if err != nil {
		log.WithError(err).Fatal("could not marshal configuration")
	}
	fmt.Printf("configuration\n\n---\n%s\n\n", string(out))
}

// UpdateLogging modifies default logging based on configuration
func UpdateLogging() {
	level, err := log.ParseLevel(config.Logging.Level)
	if err != nil {
		log.WithField("level", config.Logging.Level).Error("not a valid level")
		return
	}
	log.SetLevel(level)
}
