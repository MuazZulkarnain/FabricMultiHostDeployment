package txdefs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

// TransferCarbonCredit is a transaction to transfer carbon credits from one project to another
var TransferCarbonCredit = tx.Transaction{
	Tag:         "transferCarbonCredit",
	Label:       "Transfer Carbon Credit",
	Description: "Transfer carbon credits from one project to another",
	Method:      "PUT",
	Callers:     []string{`$org\dMSP`, "Org1MSP", "orgMSP"}, // Any orgs can call this transaction

	Args: []tx.Argument{
		{
			Tag:         "fromProject",
			Label:       "From Project",
			Description: "The project from which to transfer carbon credits",
			DataType:    "->project", // Assuming you have a "Project" asset type
			Required:    true,
		},
		{
			Tag:         "toProject",
			Label:       "To Project",
			Description: "The project to which to transfer carbon credits",
			DataType:    "->project", // Assuming you have a "Project" asset type
			Required:    true,
		},
		{
			Tag:         "amount",
			Label:       "Amount",
			Description: "Amount of carbon credits to transfer",
			DataType:    "carbonCredit", // Allow both integer and float64
			Required:    true,
		},
	},

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		fromProjectKey, ok := req["fromProject"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter fromProject must be an asset")
		}

		toProjectKey, ok := req["toProject"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter toProject must be an asset")
		}

		// Retrieve project assets from the ledger
		fromProjectAsset, err := fromProjectKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to get fromProject from the ledger")
		}

		toProjectAsset, err := toProjectKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to get toProject from the ledger")
		}

		// Update the carbon credit amount in the fromProject
		var amount float64
		switch val := req["amount"].(type) {
		case float64:
			amount = val
		case int:
			amount = float64(val)
		default:
			return nil, errors.WrapError(nil, "Parameter amount must be a number")
		}

		fromProjectMap := (map[string]interface{})(*fromProjectAsset)
		fromProjectAmount, ok := fromProjectMap["amount"].(float64)
		if !ok {
			return nil, errors.WrapError(nil, "fromProject does not have a valid amount field")
		}

		// Check if the fromProject has enough credits to transfer
		if fromProjectAmount < amount {
			return nil, errors.WrapError(nil, "Insufficient carbon credits in fromProject")
		}

		// Update the amount in the toProject
		toProjectMap := (map[string]interface{})(*toProjectAsset)
		toProjectAmount, ok := toProjectMap["amount"].(float64)
		if !ok {
			return nil, errors.WrapError(nil, "toProject does not have a valid amount field")
		}

		// Update the latest amount in both sender and receiver
		fromProjectMap["amount"] = fromProjectAmount - amount
		toProjectMap["amount"] = toProjectAmount + amount

		// Update lastAction with the last action made and the timestamp
		malaysiaTime := time.Now().UTC().Add(8 * time.Hour)
		fromProjectMap["lastAction"] = fmt.Sprintf("Sent %.2f tokens to %s at %s", amount, toProjectMap["project"], malaysiaTime.Format(time.RFC3339))
		toProjectMap["lastAction"] = fmt.Sprintf("Received %.2f tokens from %s at %s", amount, fromProjectMap["project"], malaysiaTime.Format(time.RFC3339))

		// Save the updated fromProject asset back to the ledger
		_, err = fromProjectKey.Update(stub, fromProjectMap)
		if err != nil {
			return nil, errors.WrapError(err, "failed to update fromProject asset")
		}

		// Save the updated toProject asset back to the ledger
		_, err = toProjectKey.Update(stub, toProjectMap)
		if err != nil {
			return nil, errors.WrapError(err, "failed to update toProject asset")
		}

		// Return a success response
		response := map[string]interface{}{
			"message": "Carbon credits transferred successfully",
		}
		responseBytes, nerr := json.Marshal(response)
		if nerr != nil {
			return nil, errors.WrapError(err, "failed to marshal response")
		}

		return responseBytes, nil
	},
}
