package handlers

import (
	"database/sql"
	"net/http"

	"github.com/BartCzech/swift-api/models"
	"github.com/BartCzech/swift-api/repository"
	"github.com/gin-gonic/gin"
)

// Test endpoint 0 - returns all SWIFT codes
func GetSwiftCodes(c *gin.Context) {
	rows, err := repository.DB.Query("SELECT * FROM swift_codes")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var swift models.SwiftCode
		err := rows.Scan(&swift.CountryISO2, &swift.SwiftCode, &swift.CodeType, &swift.BankName,
			&swift.Address, &swift.TownName, &swift.CountryName, &swift.TimeZone, &swift.IsHeadquarter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
			return
		}
		result := map[string]interface{}{
			"country_iso2":   swift.CountryISO2,
			"swift_code":     swift.SwiftCode,
			"code_type":      swift.CodeType,
			"bank_name":      swift.BankName,
			"address":        swift.Address,
			"town_name":      swift.TownName,
			"country_name":   swift.CountryName,
			"time_zone":      swift.TimeZone,
			"is_headquarter": swift.IsHeadquarter,
		}
		results = append(results, result)
	}
	c.JSON(http.StatusOK, results)
}

// Endpoint 1
func GetSwiftCodeDetails(c *gin.Context) {
	swiftCodeParam := c.Param("swift-code")

	row := repository.DB.QueryRow(
		`SELECT country_iso2, swift_code, code_type, bank_name, address, town_name, country_name, time_zone, is_headquarter
		 FROM swift_codes
		 WHERE swift_code = $1`, swiftCodeParam,
	)

	var swift models.SwiftCode
	err := row.Scan(&swift.CountryISO2, &swift.SwiftCode, &swift.CodeType, &swift.BankName,
		&swift.Address, &swift.TownName, &swift.CountryName, &swift.TimeZone, &swift.IsHeadquarter)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Swift code not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if swift.IsHeadquarter {
		if len(swift.SwiftCode) < 8 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid swift code format"})
			return
		}
		prefix := swift.SwiftCode[:8]
		branchRows, err := repository.DB.Query(
			`SELECT swift_code, bank_name, address, country_iso2, is_headquarter
			 FROM swift_codes
			 WHERE swift_code <> $1 AND substring(swift_code from 1 for 8) = $2`,
			swift.SwiftCode, prefix,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch branches"})
			return
		}
		defer branchRows.Close()

		var branches []map[string]interface{}
		for branchRows.Next() {
			var bSwiftCode, bBankName, bAddress, bCountryISO2 string
			var bIsHeadquarter bool
			err := branchRows.Scan(&bSwiftCode, &bBankName, &bAddress, &bCountryISO2, &bIsHeadquarter)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse branch data"})
				return
			}
			branch := map[string]interface{}{
				"swiftCode":     bSwiftCode,
				"bankName":      bBankName,
				"address":       bAddress,
				"countryISO2":   bCountryISO2,
				"isHeadquarter": bIsHeadquarter,
			}
			branches = append(branches, branch)
		}

		response := gin.H{
			"address":       swift.Address,
			"bankName":      swift.BankName,
			"country_iso2":  swift.CountryISO2,
			"country_name":  swift.CountryName,
			"is_headquarter": swift.IsHeadquarter,
			"swift_code":    swift.SwiftCode,
			"branches":      branches,
		}
		c.JSON(http.StatusOK, response)
		return
	}

	response := gin.H{
		"address":       swift.Address,
		"bankName":      swift.BankName,
		"country_iso2":  swift.CountryISO2,
		"country_name":  swift.CountryName,
		"is_headquarter": swift.IsHeadquarter,
		"swift_code":    swift.SwiftCode,
	}
	c.JSON(http.StatusOK, response)
}

// Endpoint 2
func GetSwiftCodesByCountry(c *gin.Context) {
	countryISO2 := c.Param("countryISO2code")
	rows, err := repository.DB.Query("SELECT * FROM swift_codes WHERE country_iso2 = $1", countryISO2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	defer rows.Close()

	var swiftCodes []map[string]interface{}
	var countryName string
	for rows.Next() {
		var swift models.SwiftCode
		err := rows.Scan(&swift.CountryISO2, &swift.SwiftCode, &swift.CodeType, &swift.BankName,
			&swift.Address, &swift.TownName, &swift.CountryName, &swift.TimeZone, &swift.IsHeadquarter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
			return
		}
		if len(swiftCodes) == 0 {
			countryName = swift.CountryName
		}
		entry := map[string]interface{}{
			"address":       swift.Address,
			"bankName":      swift.BankName,
			"country_iso2":  swift.CountryISO2,
			"is_headquarter": swift.IsHeadquarter,
			"swift_code":    swift.SwiftCode,
		}
		swiftCodes = append(swiftCodes, entry)
	}
	response := gin.H{
		"country_iso2": countryISO2,
		"country_name": countryName,
		"swiftCodes":   swiftCodes,
	}
	c.JSON(http.StatusOK, response)
}

// Endpoint 3
func CreateSwiftCode(c *gin.Context) {
	var newEntry struct {
		Address       string  `json:"address" binding:"required"`
		BankName      string  `json:"bankName" binding:"required"`
		CountryISO2   string  `json:"countryISO2" binding:"required"`
		CountryName   string  `json:"countryName" binding:"required"`
		IsHeadquarter bool    `json:"isHeadquarter" binding:"required"`
		SwiftCode     string  `json:"swiftCode" binding:"required"`
		// Decided to add the other fields in the following fashion:
		CodeType      string  `json:"codeType"`   // optional; default to "BIC11" if empty
		TownName      *string `json:"townName"`   // optional; allow nil for SQL NULL
		TimeZone      *string `json:"timeZone"`   // optional; allow nil for SQL NULL
	}
	if err := c.ShouldBindJSON(&newEntry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}
	if len(newEntry.SwiftCode) != 8 && len(newEntry.SwiftCode) != 11 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SwiftCode must be either 8 or 11 characters long"})
		return
	}
	if newEntry.CodeType == "" {
		newEntry.CodeType = "BIC11"
	}
	query := `
        INSERT INTO swift_codes 
            (country_iso2, swift_code, code_type, bank_name, address, town_name, country_name, time_zone, is_headquarter)
        VALUES
            ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
	_, err := repository.DB.Exec(query,
		newEntry.CountryISO2,
		newEntry.SwiftCode,
		newEntry.CodeType,
		newEntry.BankName,
		newEntry.Address,
		newEntry.TownName,
		newEntry.CountryName,
		newEntry.TimeZone,
		newEntry.IsHeadquarter,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert new swift code entry: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Swift code entry created successfully"})
}

// Endpoint 4
func DeleteSwiftCode(c *gin.Context) {
	swiftCodeParam := c.Param("swift-code")
	if len(swiftCodeParam) != 8 && len(swiftCodeParam) != 11 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Swift code must be either 8 or 11 characters long"})
		return
	}
	result, err := repository.DB.Exec("DELETE FROM swift_codes WHERE swift_code = $1", swiftCodeParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete swift code: " + err.Error()})
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not determine deletion result: " + err.Error()})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Swift code not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Swift code entry deleted successfully"})
}
