migration.up:
	sql-migrate up -config=dbconfig.yml -env="development"

migration.down:
	sql-migrate down -config=dbconfig.yml -env="development"

int.tests:
	gotestsum --format testdox -- -count=1 ./integration_tests/...