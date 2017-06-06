package main

import (
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

var joinTournamentErrorsTests = []string{
	"/joinTournament",
	"/joinTournament?tournamentId=dfd",
	"/joinTournament?tournamentId=0",
}

var joinTournamentBackerIdErrorsTests = []string{
	"/joinTournament?tournamentId=1&playerId=1&backerId=0",
	"/joinTournament?tournamentId=1&playerId=1&backerId=dfd",
}

func TestJoinTournamentErrors(t *testing.T) {
	mysqlConnect()
	defer db.Close()
	conf.initConfig()

	echoServer := echo.New()

	for _, url := range joinTournamentErrorsTests {
		req := httptest.NewRequest(echo.GET, url, nil)
		rec := httptest.NewRecorder()
		c := echoServer.NewContext(req, rec)

		assert.Error(t, takePage(c))
	}
}

func TestJoinTournamentNotAvailable(t *testing.T) {
	mysqlConnect()
	defer db.Close()
	conf.initConfig()

	db.Query("UPDATE tournaments SET status = 1 WHERE id = 1")

	echoServer := echo.New()
	req := httptest.NewRequest(echo.GET, "/joinTournament?tournamentId=1", nil)
	rec := httptest.NewRecorder()
	context := echoServer.NewContext(req, rec)

	assert.Error(t, joinTournamentPage(context))
}

func TestJoinTournamentBackerIdErrors(t *testing.T) {
	mysqlConnect()
	defer db.Close()
	conf.initConfig()

	db.Query("UPDATE tournaments SET status = 1 WHERE id = 1")

	echoServer := echo.New()

	for _, url := range joinTournamentBackerIdErrorsTests {
		req := httptest.NewRequest(echo.GET, url, nil)
		rec := httptest.NewRecorder()
		c := echoServer.NewContext(req, rec)

		assert.Error(t, takePage(c))
	}
}

func TestJoinTournamentPlayerBalanceProblem(t *testing.T) {
	mysqlConnect()
	defer db.Close()
	conf.initConfig()

	db.Query("UPDATE tournaments SET status = 1, deposit = 1000 WHERE id = 1")
	db.Query("DELETE FROM tournament_players WHERE tournament_id = 1")
	db.Query("DELETE FROM tournament_player_backers WHERE tournament_id = 1")
	db.Query("UPDATE players SET balance = 10 WHERE id = 1")

	echoServer := echo.New()
	req := httptest.NewRequest(echo.GET, "/joinTournament?tournamentId=1&playerId=1", nil)
	rec := httptest.NewRecorder()
	context := echoServer.NewContext(req, rec)

	assert.Error(t, joinTournamentPage(context))
}

func TestJoinTournamentBackerPlayerBalanceProblem(t *testing.T) {
	mysqlConnect()
	defer db.Close()
	conf.initConfig()

	db.Query("UPDATE tournaments SET status = 1, deposit = 1000 WHERE id = 1")
	db.Query("DELETE FROM tournament_players WHERE tournament_id = 1")
	db.Query("DELETE FROM tournament_player_backers WHERE tournament_id = 1")
	db.Query("UPDATE players SET balance = 10 WHERE id = 1")
	db.Query("UPDATE players SET balance = 10 WHERE id = 2")
	db.Query("UPDATE players SET balance = 10000 WHERE id = 3")

	echoServer := echo.New()
	req := httptest.NewRequest(echo.GET, "/joinTournament?tournamentId=1&playerId=1&backerId=2&backerId=3", nil)
	rec := httptest.NewRecorder()
	context := echoServer.NewContext(req, rec)

	assert.Error(t, joinTournamentPage(context))
}

func TestJoinTournamentPlayerAlreadyIn(t *testing.T) {
	mysqlConnect()
	defer db.Close()
	conf.initConfig()

	db.Query("UPDATE tournaments SET status = 1, deposit = 1000 WHERE id = 1")
	db.Query("DELETE FROM tournament_players WHERE tournament_id = 1")
	db.Query("INSERT INTO tournament_players SET tournament_id = 1, player_id = 1")
	db.Query("DELETE FROM tournament_player_backers WHERE tournament_id = 1")
	db.Query("UPDATE players SET balance = 10000 WHERE id = 1")

	echoServer := echo.New()
	req := httptest.NewRequest(echo.GET, "/joinTournament?tournamentId=1&playerId=1", nil)
	rec := httptest.NewRecorder()
	context := echoServer.NewContext(req, rec)

	assert.Error(t, joinTournamentPage(context))
}

func TestJoinTournamentBackerPlayerAlreadyIn(t *testing.T) {
	mysqlConnect()
	defer db.Close()
	conf.initConfig()

	db.Query("UPDATE tournaments SET status = 1, deposit = 1000 WHERE id = 1")
	db.Query("DELETE FROM tournament_players WHERE tournament_id = 1")
	db.Query("DELETE FROM tournament_player_backers WHERE tournament_id = 1")
	db.Query("INSERT INTO tournament_player_backers SET tournament_id = 1, player_id = 1, backer_id=2")
	db.Query("UPDATE players SET balance = 10 WHERE id = 1")
	db.Query("UPDATE players SET balance = 10000 WHERE id = 2")

	echoServer := echo.New()
	req := httptest.NewRequest(echo.GET, "/joinTournament?tournamentId=1&playerId=1", nil)
	rec := httptest.NewRecorder()
	context := echoServer.NewContext(req, rec)

	assert.Error(t, joinTournamentPage(context))
}

func TestJoinTournamentPlayerOnly(t *testing.T) {
	mysqlConnect()
	defer db.Close()
	conf.initConfig()

	db.Query("UPDATE tournaments SET status = 1, deposit = 1000 WHERE id = 1")
	db.Query("DELETE FROM tournament_players WHERE tournament_id = 1")
	db.Query("DELETE FROM tournament_player_backers WHERE tournament_id = 1")
	db.Query("UPDATE players SET balance = 10000 WHERE id = 1")

	echoServer := echo.New()
	req := httptest.NewRequest(echo.GET, "/joinTournament?tournamentId=1&playerId=1", nil)
	rec := httptest.NewRecorder()
	context := echoServer.NewContext(req, rec)

	if assert.NoError(t, joinTournamentPage(context)) {
		player := Player{}
		player.initPlayer(1)
		assert.Equal(t, player.Balance, 9000.00)
	}
}

func TestJoinTournamentPlayerWithBackers(t *testing.T) {
	mysqlConnect()
	defer db.Close()
	conf.initConfig()

	db.Query("UPDATE tournaments SET status = 1, deposit = 1000 WHERE id = 1")
	db.Query("DELETE FROM tournament_players WHERE tournament_id = 1")
	db.Query("DELETE FROM tournament_player_backers WHERE tournament_id = 1")
	db.Query("UPDATE players SET balance = 10 WHERE id = 1")
	db.Query("UPDATE players SET balance = 1000 WHERE id = 2")
	db.Query("UPDATE players SET balance = 600 WHERE id = 3")

	echoServer := echo.New()
	req := httptest.NewRequest(echo.GET, "/joinTournament?tournamentId=1&playerId=1&backerId=2&backerId=3", nil)
	rec := httptest.NewRecorder()
	context := echoServer.NewContext(req, rec)

	if assert.NoError(t, joinTournamentPage(context)) {
		player := Player{}
		player.initPlayer(2)
		assert.Equal(t, player.Balance, 500.00)

		player = Player{}
		player.initPlayer(3)
		assert.Equal(t, player.Balance, 100.00)
	}
}
