package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const WALLET_FILE = "./tmp/wallets_%s.json"

func newKeyPair() (*ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return private, pubKey
}

// Wallets stores a collection of wallets.
type Wallets struct {
	Wallets map[string]*Wallet
}

func CreateWallets(nodeID string) (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFile(nodeID)

	return &wallets, err
}

func (ws *Wallets) AddWallet() string {
	wallet := MakeWallet()
	address := string(wallet.Address())

	ws.Wallets[address] = wallet

	return address
}

func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

func (ws Wallets) GetWallet(address string) Wallet {

	fmt.Printf("Getting : %+v\n", ws.Wallets)
	return *ws.Wallets[address]
}

func (ws *Wallets) LoadFile(nodeID string) error {
	if _, err := os.Stat(FormatWalletFile(nodeID)); os.IsNotExist(err) {
		// Create a new wallet file
		os.WriteFile(FormatWalletFile(nodeID), []byte{}, 0644)
		return nil
	}

	fileContent, err := os.ReadFile(FormatWalletFile(nodeID))
	if err != nil {
		log.Panic(err)
	}

	if len(fileContent) == 0 {
		// If the file is empty, initialize an empty Wallets struct
		ws.Wallets = make(map[string]*Wallet)
		return nil
	}

	var wallets Wallets

	err = json.Unmarshal(fileContent, &wallets)
	if err != nil {
		log.Panic(err)
	}

	// set the curve for all the private keys
	// we need to this since we cannot save the curve in the json
	for address, wallet := range wallets.Wallets {
		wallet.PrivateKey.Curve = elliptic.P256()
		ws.Wallets[address] = wallet
	}

	ws.Wallets = wallets.Wallets

	return nil
}

func (ws *Wallets) SaveFile(nodeID string) {

	wsJson, err := ws.MarshalJSON()
	if err != nil {
		log.Panic(err)
	}

	// check if file exists, if not create it
	if _, err := os.Stat(FormatWalletFile(nodeID)); os.IsNotExist(err) {
		// Create a new wallet file
		os.WriteFile(FormatWalletFile(nodeID), []byte{}, 0644)
	}

	err = os.WriteFile(FormatWalletFile(nodeID), wsJson, 0644)
	if err != nil {
		log.Panic(err)
	}
}

// This is not the best way to do it but gob encoder is not working with ecdsa.PrivateKey with error:
//
//	https://stackoverflow.com/questions/73762677/panic-gob-type-elliptic-p256curve-has-no-exported-fields
//
// so this is the easiest way to do it.
func (ws *Wallets) MarshalJSON() ([]byte, error) {
	walletsMap := make(map[string]interface{})
	for address, wallet := range ws.Wallets {

		walletsMap[address] = wallet.MarshalJSON()
	}
	final := map[string]interface{}{
		"Wallets": walletsMap,
	}
	return json.Marshal(final)

}
