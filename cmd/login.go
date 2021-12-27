/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/betasve/mstd/app"
	"github.com/spf13/cobra"
)

// Represents the login command. It is the entry point of the login mechanism.
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Performs a login procedure to login into your Microsoft account",
	Long: `Triggers a browser to open the login page for MS. If your OS is not
	recognized by the tool, you will be presented an url to copy/paste in
	your browser`,
	Run: func(cmd *cobra.Command, args []string) {
		app.Login()
	},
}

// Adds the `loginCmd` to the root command, thereby enabling it for invokation
// from in the command-line tool.
func init() {
	rootCmd.AddCommand(loginCmd)
}
