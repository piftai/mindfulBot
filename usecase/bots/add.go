package bots

import (
	"context"
	"mindfulBot/models"

	"github.com/rs/zerolog/log"
)

// Add - добавляет бота.
func (u *Usecase) Add(ctx context.Context, bot models.Bot) error {
	logger := log.Ctx(ctx)

	logger.Info().Msgf("add bot with owner_id:%v, bot_name: %v", bot.OwnerID, bot.Name)

	return u.botsRepo.Add(bot)
}
