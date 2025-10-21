package fabric

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// EvaluateTransaction queries the ledger (read-only)
func EvaluateTransaction(network *client.Network, chaincodeName, fn string, args ...string) ([]byte, error) { // network jisse connect krna hai, chaincode ka naam, function jo invoke krna hai, and args passed
	contract := network.GetContract(chaincodeName) // Network se contract liya
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use EvaluateWithContext to actually apply the context
	result, err := contract.EvaluateWithContext(ctx, fn, client.WithArguments(args...)) //Uss contract se function invoke kiya and arguments pass kiye
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %v", err)
	}
	return result, nil //then result wapis krdiya
}

// SubmitTransaction submits a transaction to the ledger (write)
func SubmitTransaction(network *client.Network, chaincodeName, fn string, args ...string) ([]byte, error) { // network jisse connect krna hai, chaincode ka naam, function jo invoke krna hai, and args passed
	contract := network.GetContract(chaincodeName) //Network se chaincode uthaya
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// Use SubmitWithContext to actually apply the context
	result, err := contract.SubmitWithContext(ctx, fn, client.WithArguments(args...)) //Aur fir wo method invoke krdiya chaincode ke andar
	if err != nil {
		return nil, fmt.Errorf("failed to submit transaction: %v", err)
	}
	return result, nil //Fir result wapis krdiya
}

// GetAssetHistory fetches the history for a given asset
func GetAssetHistory(network *client.Network, chaincodeName, assetID string) ([]UserAssetHistory, error) {
	contract := network.GetContract(chaincodeName) //Same wohi network se chaincode uthaya
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bytes, err := contract.EvaluateWithContext(ctx, "GetHistoryForKey", client.WithArguments(assetID)) //Then specifically GetHistoryForKey wala function invoke kiya and arguments pass kiye
	if err != nil {
		return nil, fmt.Errorf("failed to get history: %v", err)
	}

	var history []UserAssetHistory
	if err := json.Unmarshal(bytes, &history); err != nil {
		return nil, fmt.Errorf("failed to parse history: %v", err)
	}
	return history, nil
}
