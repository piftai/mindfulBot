package postgres

import "mindfulBot/models"

// List - выводит список ботов.
func (r *Repository) List() ([]models.Bot, error) {
	var rows []bot

	err := r.db.Select(&rows, `SELECT * FROM bots ORDER BY id`)
	if err != nil {
		return nil, err
	}

	return toBots(rows), nil
}
