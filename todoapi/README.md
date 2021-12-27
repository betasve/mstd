# todoapi
--
    import "."

The `todoapi` is the package responsible for holding the layer of abstraction
that communicates with MS' API regarding the management of Lists and ToDos. It
strives to remaing as isolated and unaware of the rest of the environment as
possible. It aims to expose as little as possible. Namely only the structures
that are used to hold the data, as well as the methods who wrap its retrieving
from the API.

## Usage

#### type ContentType

```go
type ContentType string
```


#### type ListsItem

```go
type ListsItem struct {
	Id     string `json:"id"`
	Name   string `json:"displayName"`
	Owner  bool   `json:"isOwner"`
	Shared bool   `json:"isShared"`
	System string `json:"wellKnownListName"`
}
```


#### type ListsResponse

```go
type ListsResponse struct {
	Context string      `json:"@odata.context"`
	Lists   []ListsItem `json:"value"`
}
```


#### type TodoApi

```go
type TodoApi struct {
}
```


#### func (*TodoApi) ListsCreate

```go
func (ta *TodoApi) ListsCreate(name string) (*ListsItem, error)
```
Creates a ListItem setting its name.

#### func (*TodoApi) ListsIndex

```go
func (ta *TodoApi) ListsIndex() (*[]ListsItem, error)
```
Retrieves the collection of `ListItem`s.

#### func (*TodoApi) ListsUpdate

```go
func (ta *TodoApi) ListsUpdate(id, name string) (*ListsItem, error)
```
Updates a ListItem finding it by its id and changing its name.

#### func (*TodoApi) SetToken

```go
func (ta *TodoApi) SetToken(token string)
```
Sets a token to be used for the API communication.

#### func (*TodoApi) Token

```go
func (ta *TodoApi) Token() string
```
Gets the token that is used for the API communication.

#### type TodoApiClient

```go
type TodoApiClient interface {
	ListsIndex() (*[]ListsItem, error)
	ListsCreate(string) (*ListsItem, error)
	ListsUpdate(string, string) (*ListsItem, error)
	SetToken(string)
	Token() string
}
```
