package main

import (
	"errors"
)

type TournamentPlayer struct {
	player  Player
	backers []Player
}

type Tournament struct {
	Id      int
	Deposit float64
	players []TournamentPlayer
}

func (tournament *Tournament) New(deposit float64) (err error) {
	_, res, err := db.Query(`
INSERT INTO tournaments SET deposit = %f`, deposit)
	if err != nil {
		return
	}
	tournament.Deposit = deposit
	tournament.Id = int(res.InsertId())

	return
}

func (tournament *Tournament) initTournament(tournamentID int) (err error) {
	tournament.Id = tournamentID
	rows, res, err := db.Query(`
SELECT deposit FROM tournaments WHERE id = %d`, tournament.Id)
	if err != nil {
		return
	}

	if len(rows) > 0 {
		depositMap := res.Map("deposit")
		tournament.Deposit = rows[0].Float(depositMap)
	} else {
		err = errors.New("Tournament is not found")
	}

	return
}

func (tournament *Tournament) initTournamentPlayers() (err error) {
	rows, res, err := db.Query(`
SELECT player_id FROM tournament_players WHERE tournament_id = %d`, tournament.Id)
	if err != nil {
		return
	}

	for _, row := range rows {
		playerIDMap := res.Map("player_id")
		player := Player{}
		err = player.initPlayer(row.Int(playerIDMap))
		if err == nil {
			tournamentPlayer := TournamentPlayer{}
			tournamentPlayer.player = player
			err = tournamentPlayer.initTournamentPlayerBackers(tournament.Id)
			if err == nil {
				tournament.players = append(tournament.players, tournamentPlayer)
			} else {
				break
			}
		} else {
			break
		}
	}

	return
}

func (tournamentPlayer *TournamentPlayer) initTournamentPlayerBackers(tournament_id int) (err error) {
	rows, res, err := db.Query(`
SELECT backer_id FROM tournament_player_backers WHERE tournament_id = %d AND player_id = %d`, tournament_id, tournamentPlayer.player.Id)
	if err != nil {
		return
	}

	for _, row := range rows {
		backerIDMap := res.Map("backer_id")
		player := Player{}
		err = player.initPlayer(row.Int(backerIDMap))
		if err == nil {
			tournamentPlayer.backers = append(tournamentPlayer.backers, player)
		} else {
			break
		}
	}

	return
}
