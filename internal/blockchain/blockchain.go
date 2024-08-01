package blockchain

import (
	badgerdb "blockchain/pkg/badger"
	"encoding/hex"
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	DB_PATH               = "./tmp/blocks"
	DB_FILE               = "./tmp/blocks/MANIFEST"
	GENISIS_COINBASE_DATA = "First transaction in the block, reward is 100"
	LAST_HASH_KEY         = "lh"
)

type BlockChain struct {
	LastHash []byte
	DB       *badgerdb.BadgerDB
}

func NewBlockChain(address string) (*BlockChain, error) {

	if DBExists() {
		println("Blockchain already exists")
		return nil, fmt.Errorf("blockchain already exists")
	}

	db, err := badgerdb.NewBadgerDB(DB_PATH, DB_PATH)

	var lastHash []byte

	err = db.DB.Update(func(txn *badger.Txn) error {

		genesis := Genesis(NewCoinbaseTX(address, GENISIS_COINBASE_DATA))
		println("No existing blockchain found")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		if err != nil {
			return err
		}
		err = txn.Set([]byte(LAST_HASH_KEY), genesis.Hash)
		if err != nil {
			return err
		}

		lastHash = genesis.Hash

		return nil
	})

	if err != nil {
		panic(err)
	}

	return &BlockChain{LastHash: lastHash, DB: db}, nil

}

func ContinueBlockChain(address string) (*BlockChain, error) {
	if !DBExists() {
		return nil, fmt.Errorf("no existing blockchain found")
	}

	var lastHash []byte

	db, err := badgerdb.NewBadgerDB(DB_PATH, DB_PATH)

	err = db.DB.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(LAST_HASH_KEY))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
		})
		return err
	})

	if err != nil {
		return nil, err
	}

	return &BlockChain{LastHash: lastHash, DB: db}, nil
}

func (chain *BlockChain) AddBlock(transactions []*Transaction) {
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

	newBlock := CreateBlock(transactions, chain.LastHash)

	err := chain.DB.DB.Update(func(txn *badger.Txn) error {
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
	return newBlockChainIterator(chain.LastHash, chain.DB.DB)
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

func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction
	// store spent txos [txid] -> [] of spent outputs indexes
	spentTXOs := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			// collect transactions that have outputs which haven't been used
			// as inputs in other transactions for the specified address

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					// if tx is spent, skip
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				// if tx is not spent, check if it belongs to address
				// if tx belongs to address, add to unspentTxs
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}

			// if tx is not coinbase, check if inputs are unlocked by address
			if !tx.IsCoinbase() {
				for _, in := range tx.Inputs {

					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}
	return unspentTxs
}

// find address's unspent transaction outputs
func (chain *BlockChain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}
