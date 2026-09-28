package teachsimpleclient

import "mindfulBot/models"

func toBots(listBots []BotResponse) []models.Bot {
	out := make([]models.Bot, 0, len(listBots))

	for _, bot := range listBots {
		out = append(out, models.Bot{
			ID:        bot.ID,
			OwnerID:   bot.OwnerID,
			Name:      bot.Name,
			Token:     bot.Token,
			CreatedAt: bot.CreatedAt,
		})
	}

	return out
}
