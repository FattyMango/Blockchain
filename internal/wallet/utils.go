package wallet

import (
	"log"

	"github.com/mr-tron/base58"
)

// Base 58 is derived from Base 64, but has been modified to avoid some characters that look alike so addresses are not confused.
// // 0 O l I + / are removed from the base 58 alphabet
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}

func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	if err != nil {
		log.Panic(err)
	}

	return decode
}
