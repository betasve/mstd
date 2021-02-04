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
	"github.com/spf13/cobra"
)

var showColumns string

var listsCmd = &cobra.Command{
	Use:   "lists",
	Short: "Perform operations over To-Do Lists",
	Long: `A command that provides the capability of listing, creating, editing
	and deleting lists in your Microsoft To-Do account.`,
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(listsCmd)

	listsCmd.PersistentFlags().StringVarP(
		&showColumns,
		"columns", "c", "all",
		"Which columns to show, default `all`. E.g. -col=\"Display Name, Id\"",
	)
}
