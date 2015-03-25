## DocumentDB Go
> Go driver for Microsoft Azure DocumentDB 

### Note
This is a **WIP** project.  
I'm doing it on my spare time and hope to stabilize it soon. if you want to contribute, feel free to take some task [here](https://github.com/a8m/documentdb-go/issues/1)

## Table of contents:
- [Get Started](#get-started)
- [Databases](#databases)
  - [Get](#readdatabase)
  - [Query](#querydatabases)
  - [List](#readdatabases)
  - [Create](#createdatabase)
  - [Replace](#replacedatabase)
  - [Delete](#deletedatabase)
- [Collections](#collections)
  - [Get](#readcollection)
  - [Query](#querycollections)
  - [List](#readcollection)
  - [Create](#createcollection)
  - [Delete](#deletecollection)
- [Documents](#documents)
  - [Get](#readdocument)
  - [Query](#querydocuments)
  - [List](#readdocuments)
  - [Create](#createdocument)
  - [Replace](#replacedocument)
  - [Delete](#deletedocument)
- [StoredProcedures](#storedprocedures)
  - [Get](#readstoredprocedure)
  - [Query](#querystoredprocedures)
  - [List](#readstoredprocedures)
  - [Create](#createstoredprocedure)
  - [Replace](#replacestoredprocedure)
  - [Delete](#deletestoredprocedure)
- [UserDefinedFunctions](#userdefinedfunctions)
  - [Get](#readuserdefinedfunction)
  - [Query](#queryuserdefinedfunctions)
  - [List](#readuserdefinedfunctions)
  - [Create](#createuserdefinedfunction)
  - [Replace](#replaceuserdefinedfunction)
  - [Delete](#deleteuserdefinedfunction)

### Get Started
#### Installation
```sh
$ go get github.com/a8m/documentdb-go
```
#### Add to your project
```go
import (
	"github.com/a8m/documentdb"
)

func main() {
	client := documentdb.New("connection-url", documentdb.Config{"master-key"})
	// Start using DocumentDB
	dbs, err := client.ReadDatabases()
	if err != nill {
	  log.Fatal(err)
	}
	fmt.Println(dbs)
}
```

### MIT license.

