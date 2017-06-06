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

	echoServer.Logger.Fatal(echoServer.Start(conf.httpServer.host + ":" + conf.httpServer.port))
}

func mysqlConnect() {
	conf.initConfig()

	mysqlHost := conf.mysql.host + ":" + conf.mysql.port
	db = mysql.New("tcp", "", mysqlHost, conf.mysql.user, conf.mysql.password, conf.mysql.table)
	err := db.Connect()
	CheckError(err, "Connecting DB")
}

func balancePage(c echo.Context) error {
	playerIDString := c.QueryParam("playerId")
	if playerIDString != "" {
		playerID, err := strconv.Atoi(playerIDString)
		if err == nil {
			player := Player{}
			err = player.initPlayer(playerID)
			if err == nil {
				return c.JSON(http.StatusOK, player)
			}
			return echo.NewHTTPError(http.StatusNotFound, echo.Map{"message": "player not found"})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "playerId in not number"})
	}
	return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "playerId not found"})
}

