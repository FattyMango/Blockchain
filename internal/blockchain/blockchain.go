package blockchain

import (
	"github.com/dgraph-io/badger"
)

const (
	DB_PATH       = "./tmp/blocks"
	LAST_HASH_KEY = "lh"
)

type BlockChain struct {
	LastHash []byte
	DB       *badger.DB
}

func NewBlockChain() *BlockChain {

	opts := badger.DefaultOptions(DB_PATH)
	db, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}

	var lastHash []byte

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(LAST_HASH_KEY))
		if err != nil {
			genesis := Genesis()
			println("No existing blockchain found")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			if err != nil {
				panic(err)
			}
			err = txn.Set([]byte(LAST_HASH_KEY), genesis.Hash)
			if err != nil {
				panic(err)
			}

			lastHash = genesis.Hash

			return err
		}
		item, err = txn.Get([]byte(LAST_HASH_KEY))
		if err != nil {
			panic(err)
		}
		err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			panic(err)
		}

		return err
	})

	return &BlockChain{LastHash: lastHash, DB: db}

}

func (chain *BlockChain) AddBlock(data string) {
	// var lastHash []byte

	// err := chain.DB.View(func(txn *badger.Txn) error {
	// 	item, err := txn.Get([]byte(LAST_HASH_KEY))
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	err = item.Value(func(val []byte) error {
	// 		lastHash = append([]byte{}, val...)
	// 		return nil
	// 	})
	// 	return err
	// })
	// if err != nil {
	// 	panic(err)
	// }

	newBlock := CreateBlock(data, chain.LastHash)

	err := chain.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			panic(err)
		}
		err = txn.Set([]byte(LAST_HASH_KEY), newBlock.Hash)
		if err != nil {
			panic(err)
		}

		chain.LastHash = newBlock.Hash

		return err
	})

	if err != nil {
		panic(err)
	}
}

type BlockChainIterator struct {
	CurrentHash []byte
	DB          *badger.DB
}

func newBlockChainIterator(lh []byte, db *badger.DB) *BlockChainIterator {
	return &BlockChainIterator{lh, db}
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	return newBlockChainIterator(chain.LastHash, chain.DB)
}

func (i *BlockChainIterator) Next() *Block {

	var block *Block

	err := i.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(i.CurrentHash)
		if err != nil {
			panic(err)
		}
		err = item.Value(func(val []byte) error {
			block = Deserialize(val)
			return nil
		})
		return err
	})

	if err != nil {
		panic(err)
	}

	i.CurrentHash = block.PrevHash

	return block
}
