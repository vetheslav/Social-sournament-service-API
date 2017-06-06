package main

import (
	"errors"
	_ "github.com/ziutek/mymysql/native"
)

type Player struct {
	Id      int `json:"playerId"`
	Balance float64 `json:"balance"`
}

func (player *Player) initPlayer(playerID int) (err error) {
	player.Id = playerID
	rows, res, err := db.Query(`
SELECT balance FROM players WHERE id = %d`, player.Id)
	if err != nil {
		return
	}

	if len(rows) > 0 {
		balanceMap := res.Map("balance")
		player.Balance = rows[0].Float(balanceMap)
	} else {
		err = errors.New("Player not found")
	}

	return
}

func (player *Player) takePoints(points float64) (err error) {
	_, _, err = db.Query(`
UPDATE players SET balance = balance - %f WHERE id = %d`, points, player.Id)
	if err != nil {
		return
	}

	return player.initPlayer(player.Id)
}

func (player *Player) fundPoints(points float64) (err error) {
	_, _, err = db.Query(`
UPDATE players SET balance = balance + %f WHERE id = %d`, points, player.Id)
	if err != nil {
		return
	}

	return player.initPlayer(player.Id)
}
