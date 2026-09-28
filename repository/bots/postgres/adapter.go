package postgres

import "mindfulBot/models"

func toBot(b bot) models.Bot {
	return models.Bot{
		ID:        b.ID,
		StID:      b.StID,
		OwnerID:   b.OwnerID,
		Name:      b.Name,
		Token:     b.Token,
		CreatedAt: b.CreatedAt,
	}
}

func fromBot(b models.Bot) bot {
	return bot{
		ID:        b.ID,
		StID:      b.StID,
		OwnerID:   b.OwnerID,
		Name:      b.Name,
		Token:     b.Token,
		CreatedAt: b.CreatedAt,
	}
}

func toBots(in []bot) []models.Bot {
	out := make([]models.Bot, 0, len(in))

	for _, b := range in {
		out = append(out, toBot(b))
	}

	return out
}

func fromBots(in []models.Bot) []bot {
	out := make([]bot, 0, len(in))

	for _, b := range in {
		out = append(out, fromBot(b))
	}

	return out
}
