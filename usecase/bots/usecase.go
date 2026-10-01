package bots

import "mindfulBot/repository/bots"

// Usecase - юзкейс для ботов.
type Usecase struct {
	botsRepo bots.Repository
}

// New - конструктор для *Usecase.
func New(
	botsRepo bots.Repository,
) *Usecase {
	return &Usecase{
		botsRepo: botsRepo,
	}
}
