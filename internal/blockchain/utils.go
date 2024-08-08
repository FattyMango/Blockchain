package blockchain

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"os"
)

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		panic(err)
	}
	return buff.Bytes()
}

// check if the database exists by checking if the directory exists
func DBExists() bool {
	if _, err := os.Stat(DB_FILE); os.IsNotExist(err) {
		return false
	}
	return true
}

func CheckInputsExistinPreviousTransactions(inputs []TxInput, prevTXs map[string]Transaction) bool {
	for _, in := range inputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil {
			return false
		}
	}
	return true

}
