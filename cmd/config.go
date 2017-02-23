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
	"github.com/mshindle/smithy/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print out smithy's configuration",
	Long: `
Prints out the configuration used by smithy. You can
easily verify command line overrides using this command
as all global flags will still be processed.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.Print()
	},
}

func init() {
	RootCmd.AddCommand(configCmd)
}
