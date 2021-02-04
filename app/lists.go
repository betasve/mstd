package app

import (
	api "github.com/betasve/mstd/todoapi"
	"github.com/olekukonko/tablewriter"
	"os"
	"reflect"
	"strings"
)

var ColumnsToKeysMap map[string]string = map[string]string{
	"display name": "Name",
	"shared":       "Shared",
	"owner":        "Owner",
	"system name":  "System",
	"id":           "Id",
}

var ListItemHeaders []string = func(m map[string]string) []string {
	keys := []string{}

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}(ColumnsToKeysMap)

func ListsIndex(columns []string) error {
	apiClient.SetToken(config.ClientAccessToken())

	lists, err := apiClient.ListsIndex()

	if err != nil {
		return err
	}

	printResults(lists, columns)

	return nil
}

func ListsCreate(name string, columns []string) error {
	apiClient.SetToken(config.ClientAccessToken())
	newList, err := apiClient.ListsCreate(name)

	if err != nil {
		return err
	}

	printResults(&[]api.ListsItem{*newList}, columns)

	return nil
}

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

func boolToStr(b bool) string {
	if b {
		return "yes"
	} else {
		return "no"
	}
}

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
