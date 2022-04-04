package gamestrg

import (
	"context"
	"testing"

	"github.com/ghazlabs/hex-mathrush/internal/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPutGame(t *testing.T) {
	// initialize game
	game := core.Game{
		GameID:       uuid.NewString(),
		PlayerName:   "Riandy",
		Scenario:     "NEW_QUESTION",
		Score:        10,
		CountCorrect: 1,
		CurrentQuestion: &core.TimedQuestion{
			Question: core.Question{
				Problem: "1 + 2",
				Choices: []string{
					"3",
					"4",
					"5",
				},
				CorrectIndex: 1,
			},
			Timeout: 5,
		},
	}

	// initialize game storage
	gameStrg := New()
	// test the put game function
	err := gameStrg.PutGame(context.Background(), game)
	require.NoError(t, err)
}

func TestGetGame(t *testing.T) {
	// initialize game
	game := core.Game{
		GameID:       uuid.NewString(),
		PlayerName:   "Riandy",
		Scenario:     "NEW_QUESTION",
		Score:        20,
		CountCorrect: 2,
		CurrentQuestion: &core.TimedQuestion{
			Question: core.Question{
				Problem: "2 + 2",
				Choices: []string{
					"3",
					"4",
					"5",
				},
				CorrectIndex: 2,
			},
			Timeout: 5,
		},
	}

	// initialize game storage
	gameStrg := New()
	// set the game to game storage
	gameStrg.gameMap[game.GameID] = game
	// test the get game function
	actualGame, err := gameStrg.GetGame(context.Background(), game.GameID)
	require.NoError(t, err)

	// match the result
	assert.Equal(t, &game, actualGame, "mismatch game output")

}
