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

func getPlayerByPlayerID(context echo.Context) (player Player, err error) {
	player = Player{}
	playerIDString := context.QueryParam("playerId")
	if playerIDString != "" {
		var playerID int
		playerID, err = strconv.Atoi(playerIDString)
		if err == nil {
			err = player.initPlayer(playerID)
			if err != nil {
				err = echo.NewHTTPError(http.StatusNotFound, echo.Map{"message": "player not found"})
			}
		} else {
			err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "playerId in not number"})
		}
	} else {
		err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "playerId not found"})
	}

	return
}

func getPoints(context echo.Context) (points float64, err error) {
	pointsString := context.QueryParam("points")
	if pointsString != "" {
		points, err = strconv.ParseFloat(pointsString, 64)
		if err != nil {
			err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "points in not number"})
		}
	} else {
		err = echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "points not found"})
	}

	return
}
