package auth

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func VerifySignature(address, message, signature string) (bool, error) {
	signedMessage, err := hexutil.Decode(signature)
	if err != nil {
		return false, err
	}

	signedMessage[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1

	recovered, err := crypto.SigToPub(accounts.TextHash([]byte(message)), signedMessage)
	if err != nil {
		return false, err
	}

	return crypto.PubkeyToAddress(*recovered).Hex() == address, nil
}
