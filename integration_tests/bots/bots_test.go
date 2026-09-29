package bots

import (
	"mindfulBot/models"
	"mindfulBot/repository/bots/postgres"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// defaultTestDSN - dev-база из docker-compose.yml (совпадает с dbconfig.yml, env development).
const defaultTestDSN = "host=127.0.0.1 port=5432 user=admin password=admin dbname=mindfulBot-db sslmode=disable"

// connectDB - подключается к тестовой базе. DSN можно переопределить через TEST_DB_DSN.
func connectDB(t *testing.T) *sqlx.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = defaultTestDSN
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		t.Skipf("база недоступна (%v), запусти docker compose up -d и sql-migrate up", err)
	}
	t.Cleanup(func() { db.Close() })

	return db
}

func TestBotsRepo(t *testing.T) {
	t.Run("Проверяем подключение к базе данных.", func(t *testing.T) {
		db := connectDB(t)
		require.NoError(t, db.Ping())

		botsRepo := postgres.New(db)

		t.Run("Проверяем добавление нового бота", func(t *testing.T) {
			bot := models.Bot{
				StID:      3,
				OwnerID:   777,
				Name:      "egor_bot",
				Token:     "HSA:43021901219 598122 402109102 495100 5421004 10",
				CreatedAt: time.Now(),
			}

			err := botsRepo.Add(bot)
			assert.Nil(t, err)
		})
	})
}
