package db

import (
	"database/sql"
	"embed"
	"log"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func RunMigrations(pgDSN string) error {
	if pgDSN == "" {
		log.Println("PG_DSN não configurado, pulando migrations")
		return nil
	}

	log.Println("Executando migrations do banco de dados...")
	db, err := sql.Open("postgres", pgDSN)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrationsFS,
		Root:       "migrations",
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	if n > 0 {
		log.Printf("%d migration(s) aplicada(s) com sucesso", n)
	} else {
		log.Println("Migrations já estão atualizadas")
	}

	return nil
}
