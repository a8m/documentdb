package documentdb

type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type P = Parameter

// func P(name, value string) Parameter {
// 	return Parameter{name, value}
// }

type Query struct {
	Query      string      `json:"query"`
	Parameters []Parameter `json:"parameters"`
}

func NewQuery(query string, parameters ...Parameter) *Query {
	return &Query{query, parameters}
}
