package documentdb

// Resource
type Resource struct {
	Id	string	`json:"id,omitempty"`
	Self	string	`json:"_self,omitempty"`
	Etag	string	`json:"_etag,omitempty"`
	Rid	string	`json:"_rid,omitempty"`
	Ts	int	`json:"_ts,omitempty"`
}

// Indexing policy
// TODO: Ex/IncludePaths
type IndexingPolicy struct {
	IndexingMode	string	`json: "indexingMode,omitempty"`
	Automatic	string	`json: "automatic,omitempty"`
}

// Database
type Database struct {
	Resource
	Colls	string	`json:"_colls,omitempty"`
	Users	string	`json:"_users,omitempty"`
}

// Collection
type Collection struct {
	Resource
	IndexingPolicy	IndexingPolicy	`json:"indexingPolicy,omitempty"`
	docs		string		`json:"_docs,omitempty"`
	udf		string		`json:"_udf,omitempty"`
	sporcs		string		`json:"_sporcs,omitempty"`
	triggers	string		`json:"_triggers,omitempty"`
	conflicts	string		`json:"_conflicts,omitempty"`
}

// Document
type Documents struct {
	Resource
	attachments	string	`json: "attachments,omitempty"`
}
