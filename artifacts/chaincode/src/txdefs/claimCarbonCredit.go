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

// ClaimCarbonCredit is a transaction to allow a project to claim carbon credits
var ClaimCarbonCredit = tx.Transaction{
	Tag:         "claimCarbonCredit",
	Label:       "Claim Carbon Credit",
	Description: "Allow a project to claim carbon credits",
	Method:      "PUT",
	Callers:     []string{`$org\dMSP`, "orgMSP"}, // Any orgs can call this transaction

	Args: []tx.Argument{
		{
			Tag:         "project",
			Label:       "Project",
			Description: "The project to claim carbon credits from",
			DataType:    "->project", // Assuming you have a "Project" asset type
			Required:    true,
		},
		{
			Tag:         "amount",
			Label:       "Amount",
			Description: "Amount of carbon credits to claim",
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

		projectMap := (map[string]interface{})(*projectAsset)
		projectAmount, ok := projectMap["amount"].(float64)
		if !ok {
			return nil, errors.WrapError(nil, "Project does not have a valid amount field")
		}

		// Check if the project has enough credits to claim
		if projectAmount < amount {
			return nil, errors.WrapError(nil, "Insufficient carbon credits in the project")
		}

		// Update the claimed amount in the project
		claimAmount, ok := projectMap["claimAmount"].(float64)
		if !ok {
			claimAmount = 0
		}
		projectMap["claimAmount"] = claimAmount + amount

		// Update the amount in the project
		projectMap["amount"] = projectAmount - amount

		// Update lastAction with the last action made and the timestamp
		malaysiaTime := time.Now().UTC().Add(8 * time.Hour)
		projectMap["lastAction"] = fmt.Sprintf("Claimed %.2f tokens at %s", amount, malaysiaTime.Format(time.RFC3339))

		// Save the updated project asset back to the ledger
		_, err = projectKey.Update(stub, projectMap)
		if err != nil {
			return nil, errors.WrapError(err, "failed to update project asset")
		}

		// Return a success response
		response := map[string]interface{}{
			"message": "Carbon credits claimed successfully",
		}
		responseBytes, nerr := json.Marshal(response)
		if nerr != nil {
			return nil, errors.WrapError(err, "failed to marshal response")
		}

		return responseBytes, nil
	},
}
