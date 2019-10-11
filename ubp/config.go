package ubp

import (
	"encoding/json"
	"log"
	"os"
)

// var UbpConfig Config = LoadConfiguration("config.dev.json")

type Config struct {
	BaseUrl               string `json:"baseUrl"`
	PartnerAuthPath       string `json:"partnerAuthPath"`
	InstapayPath          string `json:"instapayPath"`
	InstapayGetBanksPath  string `json:"instapayGetBanksPath"`
	PesonetGetBanksPath   string `json:"pesonetGetBanksPath"`
	GetTransferStatusPath string `json:"getTransferStatusPath"`
	PesonetPath           string `json:"pesonetPath"`
	ClientId              string `json:"clientId"`
	ClientSecret          string `json:"clientSecret"`
	PartnerId             string `json:"partnerId"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	Scope                 string `json:"scope"`
}

// LoadConfiguration to load json config file
func (c *Config) LoadConfiguration(file string) {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Fatalln(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&c)

	return
}
