package gamestrg

import (
	"context"

	"github.com/ghazlabs/hex-mathrush/internal/core"
)

type Storage struct {
	gameMap map[string]core.Game
}

func New() *Storage {
	return &Storage{gameMap: map[string]core.Game{}}
}

func (s *Storage) PutGame(ctx context.Context, g core.Game) error {
	s.gameMap[g.GameID] = g
	return nil
}

func (s *Storage) GetGame(ctx context.Context, gameID string) (*core.Game, error) {
	game := s.gameMap[gameID]
	return &game, nil
}
