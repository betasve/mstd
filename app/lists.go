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

package app

import (
	api "github.com/betasve/mstd/todoapi"
	"github.com/olekukonko/tablewriter" // TODO: Replace with the abstracted ext/tablewriter
	"os"
	"reflect"
	"strings"
)

// Maps the column values returned from the MS API to the ones we need to
// display in the CLI.
var ColumnsToKeysMap map[string]string = map[string]string{
	"display name": "Name",
	"shared":       "Shared",
	"owner":        "Owner",
	"system name":  "System",
	"id":           "Id",
}

// Constructs a list of strings representing the needed headers for the table
// thats printed as a result of the List operations.
var ListItemHeaders []string = func(m map[string]string) []string {
	keys := []string{}

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}(ColumnsToKeysMap)

// Prints a formatted table with all the lists, contains only the columns,
// listed in the `columns []string`.
func ListsIndex(columns []string) error {
	apiClient.SetToken(config.ClientAccessToken())

	lists, err := apiClient.ListsIndex()

	if err != nil {
		return err
	}

	printResults(lists, columns)
	return nil
}

// Creates a new list item and prints it back to output, formatted with the
// list of columns mentioned in the `columns []string`.
func ListsCreate(name string, columns []string) error {
	apiClient.SetToken(config.ClientAccessToken())
	newList, err := apiClient.ListsCreate(name)

	if err != nil {
		return err
	}

	printResults(&[]api.ListsItem{*newList}, columns)

	return nil
}

// Updates the name of a list. Upon success it returns the updated list with its
// attributes in columns to the CLI.
// TODO: Extend the update posibilities to other attributes too (e.g. set as a
// default list)
func ListsUpdate(id, name string, columns []string) error {
	apiClient.SetToken(config.ClientAccessToken())
	list, err := apiClient.ListsUpdate(id, name)

	if err != nil {
		return err
	}

	printResults(&[]api.ListsItem{*list}, columns)

	return nil
}

// With the received params [ListItem]s and columns it renders a table with
// the `columns` as headers of the table, and each ListItem's attributes for
// that column.
func printResults(lists *[]api.ListsItem, columns []string) {
	columnsToHeaders := columnsHeadersIntersection(columns)
	keys := columnsToKeys(columnsToHeaders, ColumnsToKeysMap)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(columnsToHeaders)

	for _, item := range *lists {
		table.Append(listStrValuesForKeys(item, keys))
	}

	table.Render()
}

// Converts a boolean value to a `yes` or `no` string.
func boolToStr(b bool) string {
	if b {
		return "yes"
	} else {
		return "no"
	}
}

// Filters the original list of [ListItemHeaders] columns with the ones passed
// in `columns`. Or returns all of them if `columns` is `[]string{"all"}`.
func columnsHeadersIntersection(columns []string) []string {
	if len(columns) == 1 &&
		columns[0] == "all" {
		return ListItemHeaders
	}

	result := []string{}

	headers := strings.Join(ListItemHeaders, ",")

	for _, c := range columns {
		if strings.Contains(headers, c) {
			result = append(result, c)
		}
	}

	return result
}

// Filters the selected columns to their respective keys in [ListItem] enabling
// accessing them dynamically in [listStrValuesForKeys].
func columnsToKeys(columns []string, columnsKeysMap map[string]string) []string {
	resKeys := []string{}

	for _, c := range columns {
		val := columnsKeysMap[c]
		if val != "" {
			resKeys = append(resKeys, val)
		}
	}

	return resKeys
}

// Composes a `[]string` of the values for a particular `ListItem` base on it
// and a list of its keys that are requested.
func listStrValuesForKeys(listItem api.ListsItem, keys []string) []string {
	row := []string{}

	for _, k := range keys {
		val := reflect.ValueOf(listItem).FieldByName(k)
		switch val.Kind() {
		case reflect.Bool:
			row = append(row, boolToStr(val.Interface().(bool)))
		case reflect.String:
			row = append(row, val.Interface().(string))
		}
	}

	return row
}
