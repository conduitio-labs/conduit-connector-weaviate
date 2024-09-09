// Code generated by paramgen. DO NOT EDIT.
// Source: github.com/ConduitIO/conduit-commons/tree/main/paramgen

package destination

import (
	"github.com/conduitio/conduit-commons/config"
)

const (
	ConfigAuthApiKey           = "auth.apiKey"
	ConfigAuthMechanism        = "auth.mechanism"
	ConfigAuthWcsCredsPassword = "auth.wcsCreds.password"
	ConfigAuthWcsCredsUsername = "auth.wcsCreds.username"
	ConfigClass                = "class"
	ConfigEndpoint             = "endpoint"
	ConfigGenerateUUID         = "generateUUID"
	ConfigModuleHeaderName     = "moduleHeader.name"
	ConfigModuleHeaderValue    = "moduleHeader.value"
	ConfigScheme               = "scheme"
)

func (Config) Parameters() map[string]config.Parameter {
	return map[string]config.Parameter{
		ConfigAuthApiKey: {
			Default:     "",
			Description: "A Weaviate API key.",
			Type:        config.ParameterTypeString,
			Validations: []config.Validation{},
		},
		ConfigAuthMechanism: {
			Default:     "none",
			Description: "Mechanism specifies in which way the connector will authenticate to Weaviate.",
			Type:        config.ParameterTypeString,
			Validations: []config.Validation{
				config.ValidationInclusion{List: []string{"none", "apiKey", "wcsCreds"}},
			},
		},
		ConfigAuthWcsCredsPassword: {
			Default:     "",
			Description: "WCS password",
			Type:        config.ParameterTypeString,
			Validations: []config.Validation{},
		},
		ConfigAuthWcsCredsUsername: {
			Default:     "",
			Description: "WCS username",
			Type:        config.ParameterTypeString,
			Validations: []config.Validation{},
		},
		ConfigClass: {
			Default:     "",
			Description: "The class name as defined in the schema.\nA record will be saved under this class unless\nit has the `weaviate.class` metadata field.",
			Type:        config.ParameterTypeString,
			Validations: []config.Validation{
				config.ValidationRequired{},
			},
		},
		ConfigEndpoint: {
			Default:     "",
			Description: "Host of the Weaviate instance.",
			Type:        config.ParameterTypeString,
			Validations: []config.Validation{
				config.ValidationRequired{},
			},
		},
		ConfigGenerateUUID: {
			Default:     "",
			Description: "Whether a UUID for records should be automatically generated.\nThe generated UUIDs are MD5 sums of record keys.",
			Type:        config.ParameterTypeBool,
			Validations: []config.Validation{},
		},
		ConfigModuleHeaderName: {
			Default:     "",
			Description: "Name of the header configuring a module (e.g. `X-OpenAI-Api-Key`)",
			Type:        config.ParameterTypeString,
			Validations: []config.Validation{},
		},
		ConfigModuleHeaderValue: {
			Default:     "",
			Description: "Value for header given in `moduleHeader.name`.",
			Type:        config.ParameterTypeString,
			Validations: []config.Validation{},
		},
		ConfigScheme: {
			Default:     "https",
			Description: "Scheme of the Weaviate instance.",
			Type:        config.ParameterTypeString,
			Validations: []config.Validation{
				config.ValidationInclusion{List: []string{"http", "https"}},
			},
		},
	}
}
