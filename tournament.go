package main

import (
	"errors"
	"math/rand"
	"time"
)

type TournamentPlayer struct {
	Player  Player
	Backers []Player
}

type Tournament struct {
	Id      int     `json:"tournamentId"`
	Deposit float64 `json:"deposit"`
	status  bool
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
SELECT deposit, status FROM tournaments WHERE id = %d`, tournament.Id)
	if err != nil {
		return
	}

	if len(rows) > 0 {
		depositMap := res.Map("deposit")
		statusMap := res.Map("status")
		tournament.Deposit = rows[0].Float(depositMap)
		tournament.status = rows[0].Bool(statusMap)
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
			tournamentPlayer.Player = player
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

func (tournament Tournament) isAvailable() bool {
	return tournament.status == true
}

func (tournament *Tournament) addPlayerToTournament(player Player, backers []Player, needBackerBalance float64) (err error) {
	_, _, err = db.Query(`
INSERT INTO tournament_players SET player_id = %d, tournament_id = %d`, player.Id, tournament.Id)
	if err != nil {
		return
	}

	if len(backers) == 0 {
		needBackerBalance = tournament.Deposit
	}
	err = player.takePoints(needBackerBalance)

	tournament.players = append(tournament.players, TournamentPlayer{Player: player, Backers: backers})

	return
}

func (tournament *Tournament) addBackerToTournament(player Player, backers []Player) (err error, needBackerBalance float64) {
	countBackers := float64(len(backers))
	if countBackers > 0 {
		needBackerBalance = tournament.Deposit / countBackers
		for _, backer := range backers {
			_, _, errQuery := db.Query(`
INSERT INTO tournament_player_backers SET player_id = %d, tournament_id = %d, backer_id = %d`, player.Id, tournament.Id, backer.Id)
			if errQuery != nil {
				return errQuery, needBackerBalance
			}

			err = backer.takePoints(needBackerBalance)
		}
	}
	return
}

func (tournament Tournament) isPlayerInTournament(player Player) bool {
	for _, players := range tournament.players {
		if players.Player.Id == player.Id {
			return true
		}
	}

	return false
}

func (tournament Tournament) isPlayerInTournamentBackers(player Player) bool {
	for _, players := range tournament.players {
		for _, backer := range players.Backers {
			if backer.Id == player.Id {
				return true
			}
		}
	}

	return false
}

func (tournament Tournament) canPlayerToParticipateByBalance(player Player) bool {
	return player.Balance >= tournament.Deposit
}

func (tournament Tournament) canBackersParticipateTournament(backers []Player) bool {
	countBackers := float64(len(backers))
	if countBackers > 0 {
		needBackerBalance := tournament.Deposit / countBackers
		for _, backer := range backers {
			if backer.Balance < needBackerBalance {
				return false
			}
		}

		return true
	}

	return false
}

func (tournament Tournament) getWinner() (player TournamentPlayer, err error) {
	countPlayers := len(tournament.players)
	if countPlayers >= 2 {
		rand.Seed(time.Now().Unix())
		player = tournament.players[rand.Intn(countPlayers)]
	} else {
		err = errors.New("Tournament needs more then 1 participates")
	}
	return
}

func (tournament *Tournament) complete(winner TournamentPlayer) (prize float64, err error) {
	countPlayers := len(tournament.players)
	countBackers := len(winner.Backers)
	prize = tournament.Deposit * float64(countPlayers)

	if countBackers == 0 {
		err = winner.Player.fundPoints(prize)
	} else {
		playerPrize := prize / float64(countBackers+1)
		err = winner.Player.fundPoints(playerPrize)
		if err == nil {
			for _, backerWinner := range winner.Backers {
				err = backerWinner.fundPoints(playerPrize)
				if err != nil {
					break
				}
			}
		}
	}

	if err == nil {
		err = tournament.setComplete()
	}

	return
}

func (tournament *Tournament) setComplete() error {
	_, _, err := db.Query(`
UPDATE tournaments SET status = false WHERE id = %d`, tournament.Id)
	tournament.status = false

	return err
}

func (tournamentPlayer *TournamentPlayer) initTournamentPlayerBackers(tournament_id int) (err error) {
	rows, res, err := db.Query(`
SELECT backer_id FROM tournament_player_backers WHERE tournament_id = %d AND player_id = %d`, tournament_id, tournamentPlayer.Player.Id)
	if err != nil {
		return
	}

	for _, row := range rows {
		backerIDMap := res.Map("backer_id")
		player := Player{}
		err = player.initPlayer(row.Int(backerIDMap))
		if err == nil {
			tournamentPlayer.Backers = append(tournamentPlayer.Backers, player)
		} else {
			break
		}
	}

	return
}
