package documentdb

// Resource
type Resource struct {
	Id	string	`json: "id"`
	Self	string	`json: "_self"`
	etag	string	`json: "_etag"`
	rid	string	`json: "_rid"`
	ts	string	`json: "_ts"`
}

// Indexing policy
// TODO: Ex/IncludePaths
type IndexingPolicy struct {
	IndexingMode	string	`json: "indexingMode"`
	Automatic	string	`json: "automatic"`
}

// Database
type Database struct {
	Resource
	colls	string	`json: "_colls"`
}

// Collection
type Collection struct {
	Resource
	IndexingPolicy	IndexingPolicy	`json: "indexingPolicy"`
	docs		string		`json: "_docs"`
	udf		string		`json: "_udf"`
	sporcs		string		`json: "_sporcs"`
	triggers	string		`json: "_triggers"`
	conflicts	string		`json: "_conflicts"`
}

// Document
type Documents struct {
	Resource
	attachments	string	`json: "attachments"`
}
