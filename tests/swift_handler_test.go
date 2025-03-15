package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BartCzech/swift-api/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetSwiftCodeDetails_InvalidLength(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/v1/swift-codes/:swift-code", handlers.GetSwiftCodeDetails)

	req, _ := http.NewRequest("GET", "/v1/swift-codes/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Swift code must be either 8 or 11 characters long")
}
