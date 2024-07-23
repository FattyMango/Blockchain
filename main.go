package main

import (
	"blockchain/internal/blockchain"
	"fmt"
)

func main() {
	chain := blockchain.NewBlockChain()

	chain.AddBlock("First Block after Genesis")
	chain.AddBlock("Second Block after Genesis")
	chain.AddBlock("Third Block after Genesis")
	iter := chain.Iterator()
	for block := iter.Next(); block != nil; block = iter.Next() {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %v\n", pow.Validate())
		fmt.Println()

		if len(block.PrevHash) == 0 {
			fmt.Println("Genesis Block")
			break
		}

	}
}
