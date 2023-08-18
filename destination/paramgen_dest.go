// Code generated by paramgen. DO NOT EDIT.
// Source: github.com/ConduitIO/conduit-connector-sdk/tree/main/cmd/paramgen

package destination

import (
	sdk "github.com/conduitio/conduit-connector-sdk"
)

func (DestinationConfig) Parameters() map[string]sdk.Parameter {
	return map[string]sdk.Parameter{
		"api_key": {
			Default:     "",
			Description: "",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{},
		},
		"class": {
			Default:     "",
			Description: "",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{
				sdk.ValidationRequired{},
			},
		},
		"endpoint": {
			Default:     "",
			Description: "",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{
				sdk.ValidationRequired{},
			},
		},
		"generate_uuid": {
			Default:     "",
			Description: "",
			Type:        sdk.ParameterTypeBool,
			Validations: []sdk.Validation{},
		},
		"scheme": {
			Default:     "",
			Description: "",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{
				sdk.ValidationRequired{},
			},
		},
	}
}
