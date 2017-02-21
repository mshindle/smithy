// Copyright © 2017 Michael Shindle <mshindle@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"errors"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/mshindle/smithy/config"
	"github.com/mshindle/smithy/crypt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var argAsString bool
var processor config.Processor

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "encrypts a string or file with a public key",
	Long:  `encrypts a string or file with a public key`,
	PreRun: func(cmd *cobra.Command, args []string) {
		argAsString = viper.GetBool("string")
		switch viper.GetString("format") {
		case "yaml", "yml":
			processor = config.NewYamlProcessor()
		default:
			log.WithField("format", viper.GetBool("format")).Fatal("unsupported format")
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
	var data [][]byte
	var err error
	var encryptedValues = make(map[string]interface{})
	var label = viper.GetString("label")

	data, err = parseEncryptArgs(args)
	if err != nil {
		log.WithError(err).Error("could not encrypt args")
		return
	}

	if len(data) == 1 {
		encryptedValues[label], err = crypt.EncryptToString(data[0], label, config.PublicKey())
		if err != nil {
			log.WithError(err).Fatal("encryption failed")
			return
		}
	} else {
		strings := make([]string, len(data))
		for i := range data {
			strings[i], err = crypt.EncryptToString(data[i], label, config.PublicKey())
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
	var data [][]byte
	var err error

	switch len(args) {
	case 0:
		if argAsString {
			return nil, errors.New("must specify at least one argument with --string flag")
		}
		data = make([][]byte, 1)
		data[0], err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Error("could not read from stdin")
		}
	default:
		data = make([][]byte, len(args))
		for i, arg := range args {
			data[i], err = parseArg(arg)
			if err != nil {
				log.WithField("arg", arg).Error("could not parse arg")
				break
			}
		}
	}
	return data, err
}

func parseArg(arg string) ([]byte, error) {
	if argAsString {
		return []byte(arg), nil
	}
	return ioutil.ReadFile(arg)
}
