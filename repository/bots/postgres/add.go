package postgres

import "mindfulBot/models"

// Add - добавляет бота в репозиторий.
func (r *Repository) Add(bot models.Bot) error {
	_, err := r.db.Exec(`INSERT INTO bots(st_id, owner_id, bot_name, telegram_token)
			VALUES ($1, $2, $3, $4)
			`, bot.StID, bot.OwnerID, bot.Name, bot.Token)
	if err != nil {
		return err
	}

	return nil
}
