package blockchain

import (
	"bytes"
	"encoding/binary"
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
