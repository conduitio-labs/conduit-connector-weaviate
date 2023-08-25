// Code generated by paramgen. DO NOT EDIT.
// Source: github.com/ConduitIO/conduit-connector-sdk/tree/main/cmd/paramgen

package destination

import (
	sdk "github.com/conduitio/conduit-connector-sdk"
)

func (DestinationConfig) Parameters() map[string]sdk.Parameter {
	return map[string]sdk.Parameter{
		"apiKey": {
			Default:     "",
			Description: "TODO: support additional auth schemes __sL__",
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
		"generateUUID": {
			Default:     "",
			Description: "",
			Type:        sdk.ParameterTypeBool,
			Validations: []sdk.Validation{},
		},
		"moduleAPIKey.name": {
			Default:     "",
			Description: "",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{},
		},
		"moduleAPIKey.value": {
			Default:     "",
			Description: "",
			Type:        sdk.ParameterTypeString,
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
