package main

/*
Для старта без тестов
go run `ls *.go | grep -v _test.go`
*/

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/ziutek/mymysql/mysql"
	"log"
	"net/http"
	"strconv"
)

func CheckError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

var db mysql.Conn

func main() {
	mysqlConnect()
	defer db.Close()
	conf.initConfig()
	echoServer := echo.New()

	echoServer.GET("balance", balancePage)
	echoServer.GET("take", takePage)
	echoServer.GET("fund", fundPage)
	echoServer.GET("announceTournament", announceTournamentPage)
	echoServer.GET("joinTournament", joinTournamentPage)
	echoServer.GET("resultTournament", resultTournamentPage)

	echoServer.Logger.Fatal(echoServer.Start(conf.httpServer.host + ":" + conf.httpServer.port))
}

func mysqlConnect() {
	conf.initConfig()

	mysqlHost := conf.mysql.host + ":" + conf.mysql.port
	db = mysql.New("tcp", "", mysqlHost, conf.mysql.user, conf.mysql.password, conf.mysql.table)
	err := db.Connect()
	CheckError(err, "Connecting DB")
}

func balancePage(context echo.Context) error {
	player, err := getPlayerByPlayerID(context)
	if err == nil {
		return context.JSON(http.StatusOK, player)
	}

	return err
}

func takePage(context echo.Context) error {
	player, err := getPlayerByPlayerID(context)
	if err == nil {
		var points float64
		points, err = getPoints(context)
		if err == nil {
			if player.Balance >= points {
				err = player.takePoints(points)
				if err == nil {
					return context.JSON(http.StatusOK, player)
				}
			} else {
				err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "not enough points"})
			}
		}
	}

	return err
}

func fundPage(context echo.Context) error {
	player, err := getPlayerByPlayerID(context)
	if err == nil {
		var points float64
		points, err = getPoints(context)
		if err == nil {
			err = player.fundPoints(points)
			if err == nil {
				return context.JSON(http.StatusOK, player)
			}
		}
	}

	return err
}

func announceTournamentPage(context echo.Context) error {
	tournament := Tournament{}
	deposit, err := getDeposit(context)
	if err == nil {
		err = tournament.New(deposit)
		if err == nil {
			return context.JSON(http.StatusOK, tournament)
		}
		err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "there are problems with create tournament: " + err.Error()})
	}

	return err
}

func joinTournamentPage(context echo.Context) error {
	var needFee float64
	tournament, err := getTournamentByTournamentID(context)
	if err == nil {
		if tournament.isAvailable() {
			var player Player
			player, err = getPlayerByPlayerID(context)
			if err == nil {
				var backers []Player
				backers, err = getBakersByBakerID(context)
				if err == nil {
					if tournament.canPlayerToParticipateByBalance(player) || tournament.canBackersParticipateTournament(backers) {
						if !tournament.isPlayerInTournament(player) && !tournament.isPlayerInTournamentBackers(player) {
							err, needFee = tournament.addBackerToTournament(player, backers)
							if err == nil {
								err = tournament.addPlayerToTournament(player, backers, needFee)
								if err == nil {
									return context.JSON(http.StatusOK, tournament)
								}
							}
						} else {
							err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "player is already in the tournament"})
						}
					} else {
						err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "not enough points for participate"})
					}
				} else {
					err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "there are problems with add backers: " + err.Error()})
				}
			}
		} else {
			err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "tournament is not available"})
		}
	}
	return err
}

func resultTournamentPage(context echo.Context) error {
	tournament, err := getTournamentByTournamentID(context)
	if err == nil {
		if tournament.isAvailable() {
			var winnerPlayer TournamentPlayer
			winnerPlayer, err = tournament.getWinner()
			if err == nil {
				var prize float64
				prize, err = tournament.complete(winnerPlayer)
				if err == nil {
					err = context.JSON(http.StatusOK, echo.Map{"winners": echo.Map{"playerId": winnerPlayer.Player.Id, "prize": prize}})
				}
			} else {
				err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": err.Error()})
			}
		} else {
			err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "tournament is not available"})
		}
	}

	return err
}

func getPlayerByPlayerID(context echo.Context) (player Player, err error) {
	player = Player{}
	playerIDString := context.QueryParam("playerId")
	if playerIDString != "" {
		var playerID int
		playerID, err = strconv.Atoi(playerIDString)
		if err == nil {
			err = player.initPlayer(playerID)
			if err != nil {
				err = echo.NewHTTPError(http.StatusNotFound, echo.Map{"message": "player is not found"})
			}
		} else {
			err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "playerId is not number"})
		}
	} else {
		err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "playerId is not found"})
	}

	return
}

func getBakersByBakerID(context echo.Context) (players []Player, err error) {
	backerIDs := context.QueryParams()["backerId"]
	for _, backerIDString := range backerIDs {
		var backerID int
		backerID, err = strconv.Atoi(backerIDString)
		if err == nil {
			player := Player{}
			err = player.initPlayer(backerID)
			if err == nil {
				players = append(players, player)
			} else {
				err = echo.NewHTTPError(http.StatusNotFound, echo.Map{"message": "backer " + backerIDString + " is not found"})
			}
		} else {
			err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "backerId is not number"})
		}
	}

	return
}

func getTournamentByTournamentID(context echo.Context) (tournament Tournament, err error) {
	tournament = Tournament{}
	tournamentIDString := context.QueryParam("tournamentId")
	if tournamentIDString != "" {
		var tournamentID int
		tournamentID, err = strconv.Atoi(tournamentIDString)
		if err == nil {
			err = tournament.initTournament(tournamentID)
			if err == nil {
				err = tournament.initTournamentPlayers()
				if err != nil {
					err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "there are problems with init players: " + err.Error()})
				}
			} else {
				err = echo.NewHTTPError(http.StatusNotFound, echo.Map{"message": "tournament is not found"})
			}
		} else {
			err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "tournamentId is not number"})
		}
	} else {
		err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "tournamentId is not found"})
	}

	return
}

func getPoints(context echo.Context) (points float64, err error) {
	pointsString := context.QueryParam("points")
	if pointsString != "" {
		points, err = strconv.ParseFloat(pointsString, 64)
		if err != nil {
			err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "points is not number"})
		}
	} else {
		err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "points isn't found"})
	}

	return
}

func getDeposit(context echo.Context) (deposit float64, err error) {
	depositString := context.QueryParam("deposit")
	if depositString != "" {
		deposit, err = strconv.ParseFloat(depositString, 64)
		if err != nil {
			err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "deposit is not number"})
		}
	} else {
		err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "deposit is not found"})
	}

	return
}
