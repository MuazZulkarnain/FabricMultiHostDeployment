package datatypes

import (
	"fmt"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
)

var carbonCredit = assets.DataType{
	AcceptedFormats: []string{"number"},
	Parse: func(data interface{}) (string, interface{}, errors.ICCError) {
		// Use fmt.Sprintf for formatting the float to a string with 3 decimal places
		formattedCredit := fmt.Sprintf("%.3f", data)

		return formattedCredit, data, nil
	},
}
