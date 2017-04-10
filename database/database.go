package database

import "cloud.google.com/go/datastore"

type database struct {
	/*point to the datastore client*/
	DB *datastore.Client
}

/*DB Exports the db pointer*/
var DB database
