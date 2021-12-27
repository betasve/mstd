# app
--
    import "."

The `app` package is holding logic for coordinating the features of the app
itself. Holding its configuration too. It's faily isolated from the logic of
doing the actual communication (part of the `todoapi` package). The `app`
package aims to be the link between the commands in `cmd` and the communication
results provided from the `todoapi`.

## Usage

```go
var CfgFilePath string
```

```go
var ColumnsToKeysMap map[string]string = map[string]string{
	"display name": "Name",
	"shared":       "Shared",
	"owner":        "Owner",
	"system name":  "System",
	"id":           "Id",
}
```
Maps the column values returned from the MS API to the ones we need to display
in the CLI.

```go
var ListItemHeaders []string = func(m map[string]string) []string {
	keys := []string{}

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}(ColumnsToKeysMap)
```
Constructs a list of strings representing the needed headers for the table thats
printed as a result of the List operations.

#### func  InitAppConfig

```go
func InitAppConfig()
```
This is app's entry point. It's being invoked by the command-line tool that is
being used. Here we read the config file from the path that's being set for it
and initializing the configuration for the app.

#### func  ListsCreate

```go
func ListsCreate(name string, columns []string) error
```
Creates a new list item and prints it back to output, formatted with the list of
columns mentioned in the `columns []string`.

#### func  ListsIndex

```go
func ListsIndex(columns []string) error
```
Prints a formatted table with all the lists, contains only the columns, listed
in the `columns []string`.

#### func  ListsUpdate

```go
func ListsUpdate(id, name string, columns []string) error
```
Updates the name of a list. Upon success it returns the updated list with its
attributes in columns to the CLI. TODO: Extend the update posibilities to other
attributes too (e.g. set as a default list)

#### func  Login

```go
func Login()
```
Facilitates the login procedure for the app. Mainly - reads credentials from the
config file and saves them in a place they can be easily accessed.

#### func  LoginNeeded

```go
func LoginNeeded() bool
```
Checks if the user needs to be logged in (again) or his current session is still
active.
