package blockchain

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
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
func DBExists(path string) bool {
	fmt.Println("Checking if database exists for" + path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
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

func FormatDBPath(nodeID string) string {
	return fmt.Sprintf(DB_PATH, nodeID)
}

func FormatDBFilePath(nodeID string) string {
	return FormatDBPath(nodeID) + "/MANIFEST"
}
