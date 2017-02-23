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
	"errors"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/mshindle/smithy/config"
	"github.com/mshindle/smithy/crypt"
	"github.com/mshindle/smithy/data"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "encrypts a string or file with a public key",
	Long:  `encrypts a string or file with a public key`,
	PreRun: func(cmd *cobra.Command, args []string) {
		argAsString = viper.GetBool("string")
		switch viper.GetString("format") {
		case "yaml", "yml":
			processor = data.NewYamlProcessor()
		default:
			log.WithField("format", viper.GetString("format")).Fatal("unsupported format")
		}
	},
	Run: encrypt,
}

func init() {
	RootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().StringP("label", "l", "label", "label to use for each encrypted string")
	encryptCmd.Flags().BoolP("string", "s", false, "encrypt args as a string instead of a file")
	encryptCmd.Flags().StringP("format", "f", "yaml", "output data format (default: yaml)")
	viper.BindPFlag("label", encryptCmd.Flags().Lookup("label"))
	viper.BindPFlag("string", encryptCmd.Flags().Lookup("string"))
	viper.BindPFlag("format", encryptCmd.Flags().Lookup("format"))
}

func encrypt(cmd *cobra.Command, args []string) {
	var d [][]byte
	var err error
	var encryptedValues = make(map[string]interface{})
	var label = viper.GetString("label")

	d, err = parseEncryptArgs(args)
	if err != nil {
		log.WithError(err).Error("could not encrypt args")
		return
	}

	if len(d) == 1 {
		encryptedValues[label], err = crypt.EncryptToString(d[0], label, config.PublicKey())
		if err != nil {
			log.WithError(err).Fatal("encryption failed")
			return
		}
	} else {
		strings := make([]string, len(d))
		for i := range d {
			strings[i], err = crypt.EncryptToString(d[i], label, config.PublicKey())
			if err != nil {
				log.WithError(err).Fatal("encryption failed")
				return
			}
			encryptedValues[label] = strings
		}
	}

	buffer, err := processor.Marshal(encryptedValues)
	if err != nil {
		log.WithError(err).Error("cannot marshal data")
		return
	}
	_, err = buffer.WriteTo(os.Stdout)
	if err != nil {
		log.WithError(err).Error("cannot write out data")
	}
}

func parseEncryptArgs(args []string) ([][]byte, error) {
	var d [][]byte
	var err error

	switch len(args) {
	case 0:
		if argAsString {
			return nil, errors.New("must specify at least one argument with --string flag")
		}
		d = make([][]byte, 1)
		d[0], err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Error("could not read from stdin")
		}
	default:
		d = make([][]byte, len(args))
		for i, arg := range args {
			d[i], err = parseArg(arg)
			if err != nil {
				log.WithField("arg", arg).Error("could not parse arg")
				break
			}
		}
	}
	return d, err
}

func parseArg(arg string) ([]byte, error) {
	if argAsString {
		return []byte(arg), nil
	}
	return ioutil.ReadFile(arg)
}
