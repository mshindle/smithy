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

	log "github.com/Sirupsen/logrus"
	"github.com/mshindle/smithy/config"
	"github.com/mshindle/smithy/crypt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const force = "force"

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate public/private key pair for encryption",
	Long: `
Smithy generates a public/private key pair for use in encrypting/decrypting fields.
The keys will be written to the files identified by the publicKey & privateKey
configuration fields. The default names are public_key.pem and private_key.pem.
Unless absolute paths are specified, the keys will be written into the baseDir.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		log.WithField("force", viper.GetBool(force)).Info("overwriting of existing key files")
	},
	Run: func(cmd *cobra.Command, args []string) {
		pubFile := config.PublicKey()
		privateFile := config.PrivateKey()

		if checkExists(pubFile) || checkExists(privateFile) {
			return
		}

		err := crypt.Generate(pubFile, privateFile)
		if err != nil {
			log.WithError(err).Fatal("could not generate keys")
		}
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)
	generateCmd.Flags().BoolP(force, "f", false, "overwrite existing keys")
	viper.BindPFlag(force, generateCmd.Flags().Lookup(force))
}

func checkExists(file string) bool {
	_, err := os.Stat(file)
	if err == nil && !viper.GetBool(force) {
		fmt.Printf("cannot overwrite existing file: %s\nuse --force if necessary.\n\n", file)
		return true
	}
	return false
}
