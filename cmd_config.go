package main

var config models.ConfigOptions

func init() {
	parser.AddCommand("config",
		"Manage the mergo configuration",
		`The config command helps you manage the mergo configuration
Do note that the options get, set, unset and delete-credentials can be set several times to be applied several times`,
		&models.configOptions)
}
