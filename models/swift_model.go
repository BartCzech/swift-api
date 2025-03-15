package models

// The representation of a SWIFT code record
type SwiftCode struct {
	CountryISO2   string `json:"country_iso2"`
	SwiftCode     string `json:"swift_code"`
	CodeType      string `json:"code_type"`
	BankName      string `json:"bank_name"`
	Address       string `json:"address"`
	TownName      string `json:"town_name"`   
	CountryName   string `json:"country_name"`
	TimeZone      string `json:"time_zone"`   
	IsHeadquarter bool   `json:"is_headquarter"`
}
