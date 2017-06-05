package main

import (
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	"fmt"
)

func TestBalancePage(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/balance", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	fmt.Println(balancePage(c))
	// Assertions
	if assert.NoError(t, balancePage(c)) {
		//"assert.Equal(t, http.StatusOK, rec.Code)
		//assert.Equal(t, userJSON, rec.Body.String())
	}
}