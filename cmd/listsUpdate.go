/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"strings"
)

var name string

var listsUpdateCmd = &cobra.Command{
	Use:   "update [ID]",
	Short: "Update a list",
	Long:  `Update a list's properties in To Do app`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if app.LoginNeeded() {
			app.Login()
		}

		return app.ListsUpdate(
			strings.Join(args, " "),
			name,
			parseStringToList(showColumns, ListSeparator, noSpaceLowerCase),
		)
	},
}

func init() {
	listsCmd.AddCommand(listsUpdateCmd)
	listsUpdateCmd.Flags().StringVarP(
		&name,
		"name", "", "",
		"Set the name property of a list",
	)
}
