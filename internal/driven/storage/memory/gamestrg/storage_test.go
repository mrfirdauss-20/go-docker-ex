package gamestrg

import (
	"context"
	"testing"

	"github.com/ghazlabs/hex-mathrush/internal/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPutGameNew(t *testing.T) {
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
	// test the put game function with new game
	err := gameStrg.PutGame(context.Background(), game)
	require.NoError(t, err)
	// get the game from storage
	actualGame := gameStrg.gameMap[game.GameID]

	// match the game initialized with the game from storage
	assert.Equal(t, game, actualGame, "mismatch game")
}

func TestPutGameOverride(t *testing.T) {
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
	// insert new game to game storage
	gameStrg.gameMap[game.GameID] = game
	// change some fields (with the same game ID)
	changedGame := core.Game{
		GameID:       game.GameID,
		PlayerName:   game.PlayerName,
		Scenario:     "NEW_QUESTION",
		Score:        20, // this is the changed field
		CountCorrect: 2,  // this is the changed field
		CurrentQuestion: &core.TimedQuestion{
			Question: core.Question{
				Problem: "2 + 1", // this is the changed field
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
	// test the put game function with override game using changedGame
	err := gameStrg.PutGame(context.Background(), changedGame)
	require.NoError(t, err)
	// get the game from storage
	actualGame := gameStrg.gameMap[game.GameID]

	// match the changedGame with the game from storage
	// to make sure it is overridden
	assert.Equal(t, changedGame, actualGame, "mismatch game")
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
	// put the game to game storage
	gameStrg.gameMap[game.GameID] = game
	// test the get game function
	actualGame, err := gameStrg.GetGame(context.Background(), game.GameID)
	require.NoError(t, err)

	// match the result
	assert.Equal(t, &game, actualGame, "mismatch game output")

}
