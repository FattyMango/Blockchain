package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var (
	// Difficulty is the number of leading bytes that must be zero in the hash
	Difficulty = 4 * 4
)

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	return &ProofOfWork{b, target}
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.HashTransactions(),
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)

}

func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte
	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			fmt.Printf("Nonce: %d\n", nonce)
			fmt.Printf("Hash: %x\n", hash)
			fmt.Printf("Target: %x\n", pow.Target)
			fmt.Println()
			break
		}

		nonce++
	}

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)

	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}
