package main

import (
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

type balancePagePair struct {
	url       string
	waitError bool
	response  string
}

var balancePageTests = []balancePagePair{
	{"/balance", true, ""},
	{"/balance?playerId=hello", true, ""},
	{"/balance?playerId=0", true, ""},
	{"/balance?playerId=1", false, ""},
}

var takePointsErrorsTests = []string{
	"/take",
	"/take?playerId=1",
	"/take?playerId=1&points=df",
}

var fundPointsErrorsTests = []string{
	"/fund",
	"/fund?playerId=1",
	"/fund?playerId=1&points=df",
}

var announceTournamentErrorsTests = []string{
	"/announceTournament",
	"/announceTournament?deposit=hello",
}

func TestBalancePage(t *testing.T) {
	mysqlConnect()
	defer db.Close()
	conf.initConfig()

	echoServer := echo.New()

	for _, test := range balancePageTests {
		req := httptest.NewRequest(echo.GET, test.url, nil)
		rec := httptest.NewRecorder()
		context := echoServer.NewContext(req, rec)

		if !test.waitError && assert.NoError(t, balancePage(context)) {
			assert.Regexp(t, "{\"playerId\":1,\"balance\":.*", rec.Body.String())
		} else {
			assert.Error(t, balancePage(context))
		}
	}
}

func TestTakePointsErrors(t *testing.T) {
	mysqlConnect()
	defer db.Close()
	conf.initConfig()

	echoServer := echo.New()

	for _, url := range takePointsErrorsTests {
		req := httptest.NewRequest(echo.GET, url, nil)
		rec := httptest.NewRecorder()
		c := echoServer.NewContext(req, rec)

		assert.Error(t, takePage(c))
	}
}

func TestTakePointsNotEnough(t *testing.T) {
	mysqlConnect()
	defer db.Close()
	conf.initConfig()

	db.Query("UPDATE players SET balance = 0 WHERE id = 1")

	echoServer := echo.New()
	req := httptest.NewRequest(echo.GET, "/take?playerId=1&points=10", nil)
	rec := httptest.NewRecorder()
	context := echoServer.NewContext(req, rec)

	assert.Error(t, takePage(context))
}

func TestTakePoints(t *testing.T) {
	mysqlConnect()
	defer db.Close()
	conf.initConfig()

	db.Query("UPDATE players SET balance = 30.02 WHERE id = 1")

	echoServer := echo.New()
	req := httptest.NewRequest(echo.GET, "/take?playerId=1&points=10.01", nil)
	rec := httptest.NewRecorder()
	context := echoServer.NewContext(req, rec)

	if assert.NoError(t, takePage(context)) {
		assert.JSONEq(t, rec.Body.String(), `{"playerId":1,"balance":20.01}`)
	}
}

func TestFundPointsErrors(t *testing.T) {
	mysqlConnect()
	defer db.Close()
	conf.initConfig()

	echoServer := echo.New()

	for _, url := range fundPointsErrorsTests {
		req := httptest.NewRequest(echo.GET, url, nil)
		rec := httptest.NewRecorder()
		c := echoServer.NewContext(req, rec)

		assert.Error(t, takePage(c))
	}
}

func TestFundPoints(t *testing.T) {
	mysqlConnect()
	defer db.Close()
	conf.initConfig()

	db.Query("UPDATE players SET balance = 30.03 WHERE id = 1")

	echoServer := echo.New()
	req := httptest.NewRequest(echo.GET, "/fund?playerId=1&points=10.01", nil)
	rec := httptest.NewRecorder()
	context := echoServer.NewContext(req, rec)

	if assert.NoError(t, fundPage(context)) {
		assert.JSONEq(t, rec.Body.String(), `{"playerId":1,"balance":40.04}`)
	}
}

func TestAnnounceTournamentErrors(t *testing.T) {
	mysqlConnect()
	defer db.Close()
	conf.initConfig()

	echoServer := echo.New()

	for _, url := range announceTournamentErrorsTests {
		req := httptest.NewRequest(echo.GET, url, nil)
		rec := httptest.NewRecorder()
		c := echoServer.NewContext(req, rec)

		assert.Error(t, takePage(c))
	}
}

func TestAnnounceTournament(t *testing.T) {
	mysqlConnect()
	defer db.Close()
	conf.initConfig()

	echoServer := echo.New()

	req := httptest.NewRequest(echo.GET, "/announceTournament?deposit=1000", nil)
	rec := httptest.NewRecorder()
	context := echoServer.NewContext(req, rec)

	if assert.NoError(t, announceTournamentPage(context)) {
		assert.Regexp(t, "{\"tournamentId\":.*", rec.Body.String())
		db.Query("DELETE FROM tournaments WHERE id = LAST_INSERT_ID()")
	}
}
