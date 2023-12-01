package domain

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Selenium struct {
		SeleniumPort  int  `json:"seleniumPort"`
		SeleniumDebug bool `json:"seleniumDebug"`
	} `json:"selenium"`
	DataBase struct {
		CollectionInvestment      string `json:"collectionInvestment"`
		CollectionInvestmentGroup string `json:"collectionInvestmentGroup"`
	} `json:"dataBase"`
	Racional struct {
		FlagScroll       bool   `json:"flagScroll"`
		LoginPageURL     string `json:"loginPageUrl"`
		MovementsPageURL string `json:"movementsPageUrl"`
	} `json:"racional"`
	Renta4 struct {
		FlagTest         bool   `json:"flagTest"`
		LoginPageURL     string `json:"loginPageUrl"`
		MovementsPageURL string `json:"movementsPageUrl"`
	} `json:"renta4"`
	GoogleSheets struct {
		SpreadsheetsInvestment      string `json:"spreadsheetsInvestment"`
		SpreadsheetsCashFlow        string `json:"spreadsheetsCashFlow"`
		SpreadsheetsInvestmentGroup string `json:"spreadsheetsInvestmentGroup"`
	} `json:"googleSheets"`
}

// LoadConfig carga la configuración desde un archivo JSON
func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error abriendo el archivo: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("error decodificando el archivo JSON: %v", err)
	}

	return &config, nil
}
