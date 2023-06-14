package adapt

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jackc/tern/migrate"
	"go.uber.org/zap"
)

type Database struct {
	Config *ServerConfig
	Log    *zap.Logger
}

func (db *Database) ConnectDB() (*sql.DB, error) {
	dbConn, err := sql.Open("pgx", db.Config.GetUri())

	if err != nil {
		db.Log.Info("msg", zap.String("msg", "failed to connect to database"))
		log.Fatal(err)
	}

	return dbConn, nil
}

type EmbeddedF5 struct {
	dirname  string
	filename string
	glob     string
}

func (e EmbeddedF5) ReadDir(dirname string) ([]os.FileInfo, error) {
	entries, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	infos := make([]os.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}
	return infos, nil
}

func (e EmbeddedF5) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (e EmbeddedF5) Glob(pattern string) (matches []string, err error) {
	return filepath.Glob(pattern)
}

func NewEmbeddedFS() migrate.MigratorFS {
	return EmbeddedF5{}
}

func (db *Database) RunMigration(dbConn *sql.DB) error {
	conn, err := dbConn.Conn(context.Background())

	if err != nil {
		return err
	}

	err = conn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn) //conn is a *pgx.Conn
		opts := migrate.MigratorOptions{
			MigratorFS: NewEmbeddedFS(),
		}

		schema := "public"
		table := fmt.Sprintf("%s.schema_version", schema)

		dir, err := os.Getwd()

		if err != nil {
			db.Log.Info("msg", zap.String("msg", "incorrect migrations path"))
			return err
		}

		migrationPath := filepath.Join(dir, "/migrations")

		db.Log.Info("msg", zap.String("migrationPath", migrationPath))

		migrator, err := migrate.NewMigratorEx(context.Background(), conn.Conn(), table, &opts)

		err = migrator.LoadMigrations(migrationPath)
		if err != nil {
			return err
		}

		if err != nil {
			db.Log.Info("msg", zap.String("msg", "failed to connect to migrator"))
			return err
		}

		if err := migrator.Migrate(context.Background()); err != nil {
			db.Log.Info("msg", zap.String("msg", "failed to run migrations"))
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
