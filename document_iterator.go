package documentdb

type DocumentIterator struct {
	coll              string
	query             *Query
	docs              interface{}
	opts              []CallOption
	continuationToken string
}

func NewDocumentIterator(coll string, query *Query, docs interface{}, opts ...CallOption) *DocumentIterator {
	return &DocumentIterator{
		coll:  coll,
		query: query,
		docs:  docs,
		opts:  opts,
	}
}

func (di *DocumentIterator) Next(db *DocumentDB) (bool, error) {
	rp, err := db.QueryDocuments(di.coll, di.query, di.docs, append(di.opts, Continuation(di.continuationToken))...)
	if err != nil {
		return false, err
	}
	di.continuationToken = rp.Continuation()
	return di.continuationToken != "", nil
}
