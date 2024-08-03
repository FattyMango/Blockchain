package blockchain

type TxOutput struct {
	// The amount of tokens
	Value int
	// Needed to unlock tokens in Value, here its the name of the user receiving the tokens
	PubKey string
}

type TxInput struct {
	// The ID of the transaction that contains the output
	ID []byte
	// The index of the output in the transaction
	Out int
	// The signature to unlock the output, here its the name of the user sending the tokens
	Sig string
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}
