package badger

import "github.com/dgraph-io/badger"

type BadgerDB struct {
	DB *badger.DB
}

func NewBadgerDB(path, filePath string) (*BadgerDB, error) {
	opts := badger.DefaultOptions(path)
	opts.ValueDir = filePath
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &BadgerDB{DB: db}, nil
}

func (bdb *BadgerDB) Close() error {
	return bdb.DB.Close()
}
