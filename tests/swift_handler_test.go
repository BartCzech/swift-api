package tests

import (
	"log"
	"os"
	"bytes"
	// "database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BartCzech/swift-api/handlers"
	"github.com/BartCzech/swift-api/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	// Try to load the .env file from the project root.
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: Error loading ../.env file. Make sure your .env file exists in the project root.")
	}
	// Initialize the database connection.
	// Ensure your .env file is available and points to a test database.
	repository.ConnectDB()

	// Run tests.
	code := m.Run()

	// Optionally, close the DB connection here if needed.
	// repository.DB.Close()

	os.Exit(code)
}

// getTestRouter initializes a Gin router with all routes for testing.
func getTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	router.GET("/v1/swift-codes", handlers.GetSwiftCodes)
	router.GET("/v1/swift-codes/:swift-code", handlers.GetSwiftCodeDetails)
	router.GET("/v1/swift-codes/country/:countryISO2code", handlers.GetSwiftCodesByCountry)
	router.POST("/v1/swift-codes", handlers.CreateSwiftCode)
	router.DELETE("/v1/swift-codes/:swift-code", handlers.DeleteSwiftCode)
	return router
}

// TestHealthCheck verifies the /ping endpoint.
func TestHealthCheck(t *testing.T) {
	router := getTestRouter()
	req, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

// Unit Test: GetSwiftCodeDetails with an invalid SWIFT code length should return 400.
func TestGetSwiftCodeDetails_InvalidLength(t *testing.T) {
	router := getTestRouter()
	req, _ := http.NewRequest("GET", "/v1/swift-codes/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Swift code must be either 8 or 11 characters long")
}

// Unit Test: Creating a SWIFT code with an invalid length should return 400.
func TestCreateSwiftCode_InvalidLength(t *testing.T) {
	router := getTestRouter()
	invalidPayload := map[string]interface{}{
		"address":       "123 Bank St, New York, NY",
		"bankName":      "Test Bank",
		"countryISO2":   "US",
		"countryName":   "United States",
		"isHeadquarter": true,
		"swiftCode":     "12345", // Invalid length
		"codeType":      "BIC11",
		"townName":      "New York",
		"timeZone":      "America/New_York",
	}
	body, _ := json.Marshal(invalidPayload)
	req, _ := http.NewRequest("POST", "/v1/swift-codes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "SwiftCode must be either 8 or 11 characters long")
}

// Integration Test: Create a new SWIFT code, retrieve it, then delete it, verifying the full flow.
func TestIntegration_CreateRetrieveDeleteSwiftCode(t *testing.T) {
	// Ensure repository.DB is connected to a test database.
	// In integration tests, it's best to run against a test instance.
	router := getTestRouter()

	// Define a valid SWIFT code payload.
	newSwift := map[string]interface{}{
		"address":       "456 Finance St, London, UK",
		"bankName":      "Finance Bank",
		"countryISO2":   "GB",
		"countryName":   "United Kingdom",
		"isHeadquarter": true,
		"swiftCode":     "FINBGB22XXX", // 11 characters
		"codeType":      "BIC11",
		"townName":      "London",
		"timeZone":      "Europe/London",
	}
	body, _ := json.Marshal(newSwift)

	// Step 1: Create the SWIFT code entry.
	createReq, _ := http.NewRequest("POST", "/v1/swift-codes", bytes.NewBuffer(body))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	router.ServeHTTP(createResp, createReq)
	assert.Equal(t, http.StatusCreated, createResp.Code)
	assert.Contains(t, createResp.Body.String(), "Swift code entry created successfully")

	// Step 2: Retrieve the created SWIFT code.
	getReq, _ := http.NewRequest("GET", "/v1/swift-codes/FINBGB22XXX", nil)
	getResp := httptest.NewRecorder()
	router.ServeHTTP(getResp, getReq)
	assert.Equal(t, http.StatusOK, getResp.Code)
	assert.Contains(t, getResp.Body.String(), `"swift_code":"FINBGB22XXX"`)

	// Step 3: Delete the SWIFT code.
	delReq, _ := http.NewRequest("DELETE", "/v1/swift-codes/FINBGB22XXX", nil)
	delResp := httptest.NewRecorder()
	router.ServeHTTP(delResp, delReq)
	assert.Equal(t, http.StatusOK, delResp.Code)
	assert.Contains(t, delResp.Body.String(), "Swift code entry deleted successfully")

	// Step 4: Confirm deletion by attempting to retrieve again.
	getReq2, _ := http.NewRequest("GET", "/v1/swift-codes/FINBGB22XXX", nil)
	getResp2 := httptest.NewRecorder()
	router.ServeHTTP(getResp2, getReq2)
	assert.Equal(t, http.StatusNotFound, getResp2.Code)
	assert.Contains(t, getResp2.Body.String(), "Swift code not found")
}

// Integration Test: Get swift codes by country endpoint.
func TestIntegration_GetSwiftCodesByCountry(t *testing.T) {
	router := getTestRouter()

	// For this test, first create an entry for country "US".
	newSwift := map[string]interface{}{
		"address":       "789 Finance Ave, New York, NY",
		"bankName":      "USA Bank",
		"countryISO2":   "US",
		"countryName":   "United States",
		"isHeadquarter": false,
		"swiftCode":     "USABUS12", // 8-character valid SWIFT code
		"codeType":      "BIC11",
		"townName":      "New York",
		"timeZone":      "America/New_York",
	}
	body, _ := json.Marshal(newSwift)
	createReq, _ := http.NewRequest("POST", "/v1/swift-codes", bytes.NewBuffer(body))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	router.ServeHTTP(createResp, createReq)
	assert.Equal(t, http.StatusCreated, createResp.Code)

	// Test the GET /v1/swift-codes/country/US endpoint.
	getReq, _ := http.NewRequest("GET", "/v1/swift-codes/country/US", nil)
	getResp := httptest.NewRecorder()
	router.ServeHTTP(getResp, getReq)
	assert.Equal(t, http.StatusOK, getResp.Code)
	// The response should include the created swift code.
	assert.Contains(t, getResp.Body.String(), "USABUS12")
}
