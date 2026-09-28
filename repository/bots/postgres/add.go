package postgres

import "mindfulBot/models"

// Add - добавляет бота в репозиторий.
func (r *Repository) Add(bot models.Bot) error {
	_, err := r.db.Exec(`INSERT INTO bots(st_id, owner_id, bot_name, telegram_token, created_at)
			VALUES ($1, $2, $3, $4, $5)
			`, bot.StID, bot.OwnerID, bot.Name, bot.Token, bot.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
