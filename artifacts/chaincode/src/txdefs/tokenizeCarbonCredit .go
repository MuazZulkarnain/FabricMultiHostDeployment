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

// Updates the tenant of a Book
// POST Method
var TokenizeCarbonCredit = tx.Transaction{
	Tag:         "tokenizeCarbonCredit",
	Label:       "Tokenize Carbon Offset",
	Description: "Add the current amount with new tokenized amount",
	Method:      "PUT",
	Callers:     []string{`$org\dMSP`, "Org1MSP", "orgMSP"}, // Any orgs can call this transaction

	Args: []tx.Argument{
		{
			Tag:         "project",
			Label:       "Name of the carbon offset project",
			Description: "Targeted Project for tokenization",
			DataType:    "->project", // Assuming you have a "Project" asset type
			Required:    true,
		},
		{
			Tag:         "amount",
			Label:       "Amount",
			Description: "New value of the carbon credit",
			DataType:    "carbonCredit", // Allow both integer and float64
			Required:    true,
		},
	},

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		projectKey, ok := req["project"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter project must be an asset")
		}

		// Retrieve project asset from the ledger
		projectAsset, err := projectKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to get project from the ledger")
		}

		// Update the carbon credit amount in the project
		var amount float64
		switch val := req["amount"].(type) {
		case float64:
			amount = val
		case int:
			amount = float64(val)
		default:
			return nil, errors.WrapError(nil, "Parameter amount must be a number")
		}

		// Assuming your project asset has an "amount" field
		projectMap := (map[string]interface{})(*projectAsset)
		currentAmount, ok := projectMap["amount"].(float64)
		if !ok {
			return nil, errors.WrapError(nil, "Project asset does not have a valid amount field")
		}

		// Update the amount
		projectMap["amount"] = currentAmount + amount

		// Update lastAction with the last action made and the timestamp
		malaysiaTime := time.Now().UTC().Add(8 * time.Hour)
		projectMap["lastAction"] = fmt.Sprintf("Tokenized %.2f tokens at %s", amount, malaysiaTime.Format(time.RFC3339))

		projectMap, err = projectAsset.Update(stub, projectMap)
		if err != nil {
			return nil, errors.WrapError(err, "failed to update asset")
		}

		projectJSON, nerr := json.Marshal(projectMap)
		if nerr != nil {
			return nil, errors.WrapError(err, "failed to marshal response")
		}

		return projectJSON, nil
	},
}
