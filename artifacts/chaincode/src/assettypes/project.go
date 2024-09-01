package assettypes

import (
	"fmt"
	"github.com/hyperledger-labs/cc-tools/assets"
)

// CarbonCredit represents an asset type for carbon credits
var Project = assets.AssetType{
	Tag:         "project",
	Label:       "Carbon Offset Project",
	Description: "List of carbon offset projects",

	Props: []assets.AssetProp{
		{
			Required: true,
			IsKey:    true,
			Tag:      "project",
			Label:    "Name of the carbon offset project",
			DataType: "string",
			Writers : []string{`$org\dMSP`, "Org1MSP", "orgMSP"},  // Any orgs can call this transaction
		},
		{
			Required:     false,
			Tag:          "amount",
			Label:        "Amount of carbon credit available in the project",
			DataType:     "carbonCredit",
			DefaultValue: 0,
			Validate: func(amount interface{}) error {
				switch v := amount.(type) {
				case float64:
					if v < 0 {
						return fmt.Errorf("amount must be positive")
					}
				case int:
					if v < 0 {
						return fmt.Errorf("amount must be positive")
					}
					// You can handle other numeric types here
				default:
					return fmt.Errorf("amount must be a numeric type")
				}

				return nil
			},
		},
		{
			Required:     false,
			Tag:          "claimAmount",
			Label:        "Claimed carbon credit",
			DataType:     "carbonCredit",
			DefaultValue: 0,
			Validate: func(claimAmount interface{}) error {
				switch v := claimAmount.(type) {
				case float64:
					if v < 0 {
						return fmt.Errorf("amount must be positive")
					}
				case int:
					if v < 0 {
						return fmt.Errorf("amount must be positive")
					}
					// You can handle other numeric types here
				default:
					return fmt.Errorf("amount must be a numeric type")
				}

				return nil
			},
		},
		{
			Required: false,
			Tag:      "lastAction",
			Label:    "last action made onto this asset",
			DataType: "string",
		},
	},
}
