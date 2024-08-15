package blockchain

import (
	"blockchain/pkg/merkle"
	"bytes"
	"encoding/gob"
	"time"
)

type Block struct {
	Timestamp    int64
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
	Height       int
}

func CreateBlock(transactions []*Transaction, prevHash []byte, height int) *Block {
	block := &Block{time.Now().Unix(), []byte{}, transactions, prevHash, 0, height}
	pow := NewProof(block)
	block.Nonce, block.Hash = pow.Run()

	return block
}

func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{}, 0)
}
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte

	for _, tx := range b.Transactions {
		ser, err := tx.Serialize()
		if err != nil {
			panic(err)
		}
		txHashes = append(txHashes, ser)
	}

	tree := merkle.NewMerkleTree(txHashes)

	return tree.RootNode.Data
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		panic(err)
	}

	return result.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		panic(err)
	}

	return &block
}

func (b *Block) IsGenesis() bool {
	return len(b.PrevHash) == 0
}
