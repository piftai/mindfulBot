package bots

import (
	"context"
	"mindfulBot/models"

	"github.com/rs/zerolog/log"
)

func (u *Usecase) List(ctx context.Context) ([]models.Bot, error) {
	logger := log.Ctx(ctx)

	logger.Info().Msg("list of bots request")

	return u.botsRepo.List()
}
