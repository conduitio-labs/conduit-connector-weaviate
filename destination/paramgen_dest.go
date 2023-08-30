// Code generated by paramgen. DO NOT EDIT.
// Source: github.com/ConduitIO/conduit-connector-sdk/tree/main/cmd/paramgen

package destination

import (
	sdk "github.com/conduitio/conduit-connector-sdk"
)

func (Config) Parameters() map[string]sdk.Parameter {
	return map[string]sdk.Parameter{
		"apiKey": {
			Default:     "",
			Description: "Weaviate API key.",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{},
		},
		"auth.mechanism": {
			Default:     "none",
			Description: "auth.mechanism specifies in which way the connector will authenticate to Weaviate.",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{
				sdk.ValidationInclusion{List: []string{"none", "apiKey", "wcsCredentials"}},
			},
		},
		"class": {
			Default:     "",
			Description: "The class name as defined in the schema. A record will be saved under this class unless it has the `weaviate.class` metadata field.",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{
				sdk.ValidationRequired{},
			},
		},
		"endpoint": {
			Default:     "",
			Description: "Host of the Weaviate instance.",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{
				sdk.ValidationRequired{},
			},
		},
		"generateUUID": {
			Default:     "",
			Description: "Whether a UUID for records should be automatically generated. The generated UUIDs are MD5 sums of record keys.",
			Type:        sdk.ParameterTypeBool,
			Validations: []sdk.Validation{},
		},
		"moduleHeader.name": {
			Default:     "",
			Description: "name of the header configuring a module (e.g. `X-OpenAI-Api-Key`)",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{},
		},
		"moduleHeader.value": {
			Default:     "",
			Description: "value for header given in `moduleHeader.name`.",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{},
		},
		"scheme": {
			Default:     "https",
			Description: "scheme of the Weaviate instance.",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{
				sdk.ValidationInclusion{List: []string{"http", "https"}},
			},
		},
		"wcs.password": {
			Default:     "",
			Description: "WCS password",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{},
		},
		"wcs.username": {
			Default:     "",
			Description: "WCS username",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{},
		},
	}
}
