package models

type RequestBody struct {
	Url string
	Alias string
}

type RequestBodyTest struct {
	Url string `mapstructure:"url" json:"url"`
	Alias string `mapstructure:"alias" json:"alias"`
}