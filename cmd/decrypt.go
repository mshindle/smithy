// Copyright © 2017 Michael Shindle <mshindle@riotgames.com>
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
	"path/filepath"

	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/mshindle/smithy/config"
	"github.com/mshindle/smithy/data"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:     "decrypt",
	Short:   "decrypt a string or file with a private key",
	Long:    `decrypt a string or file with a private key`,
	PreRunE: preDecrypt,
	Run:     runDecrypt,
}

func init() {
	RootCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().StringP("label", "l", "label", "label to use for each encrypted string")
	decryptCmd.Flags().BoolP("string", "s", false, "decrypt args as a string instead of a file")
	viper.BindPFlag("decrypt.label", decryptCmd.Flags().Lookup("label"))
	viper.BindPFlag("decrypt.string", decryptCmd.Flags().Lookup("string"))
}

func preDecrypt(cmd *cobra.Command, args []string) error {
	argAsString = viper.GetBool("string")
	if len(args) != 1 {
		return errors.New("only one argument allowed for preDecrypt")
	}
	// determine format of input file
	ext := filepath.Ext(args[0])
	switch ext {
	case ".json":
		processor = data.NewJsonProcessor()
	case ".yaml", ".yml":
		processor = data.NewYamlProcessor()
	default:
		log.WithField("format", ext).Fatal("unsupported format")
		return errors.New("unsupported format")
	}
	return nil
}

func runDecrypt(cmd *cobra.Command, args []string) {
	object, err := processor.UnmarshalFile(args[0])
	if err != nil {
		log.WithError(err).WithField("file", args[0]).Error("cannot unmarshal and decrypt file")
		return
	}

	err = object.DecryptValues(viper.GetString("decrypt.label"), config.PrivateKey())
	if err != nil {
		log.WithError(err).WithField("object", object).Error("cannot decrypt object")
		return
	}

	b, err := processor.Marshal(object)
	if err != nil {
		log.WithError(err).Error("cannot marshal decrypted file")
		return
	}

	_, err = b.WriteTo(os.Stdout)
	if err != nil {
		log.WithError(err).Error("cannot write out data")
	}
}
