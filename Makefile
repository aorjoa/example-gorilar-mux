init:
	sqlite3 thaichana.db < _scripts/init.sql
run:
	PORT=9000 DB_CONN=thaichana.db go run main.go