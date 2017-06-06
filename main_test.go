package main

import (
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

type balancePagePair struct {
	url string
	waitError bool
	response string
}

var balancePageTests = []balancePagePair{
	{"/balance", true, ""},
	{"/balance?playerId=hello", true, ""},
	{"/balance?playerId=0", true, ""},
	{"/balance?playerId=1", false, ""},
}

func TestBalancePage(t *testing.T) {
	mysqlConnect()
	defer db.Close()
	conf.initConfig()

	echoServer := echo.New()

	for _, test := range balancePageTests {
		req := httptest.NewRequest(echo.GET, test.url, nil)
		rec := httptest.NewRecorder()
		c := echoServer.NewContext(req, rec)

		if !test.waitError && assert.NoError(t, balancePage(c)) {
			assert.Regexp(t, "{\"playerId\":1,\"balance\":.*", rec.Body.String())
		} else {
			assert.Error(t, balancePage(c))
		}
	}
}