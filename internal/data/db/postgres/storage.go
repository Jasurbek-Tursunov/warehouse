package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Jasurbek-Tursunov/warehouse/pkg/config"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	pg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Storage struct {
	logger  *slog.Logger
	conn    *sql.DB
	timeout time.Duration
}

func NewStorage(logger *slog.Logger) *Storage {
	return &Storage{logger: logger}
}

func (s *Storage) MustConnect() {
	var err error

	if err = s.connect(); err != nil {
		s.logger.Error("Fail to connect", "error", err.Error())
		panic(err)
	}

	if err = s.migrateUp(); err != nil {
		s.logger.Error("Fail to migrate", "error", err.Error())
		panic(err)
	}
}

func (s *Storage) connect() error {
	const op = "postgres.storage.Connect"

	cfg := config.MustLoad[Config]()

	dst := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
	)

	db, err := sql.Open("postgres", dst)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	s.conn = db
	s.timeout = cfg.Timeout

	return nil
}

func (s *Storage) HealthCheck() {
	if err := s.conn.Ping(); err != nil {
		s.logger.Error("Fail heath check postgres storage", "error", err.Error())
		panic(err)
	}
}

func (s *Storage) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *Storage) migrateUp() error {
	const op = "postgres.storage.MigrateUp"

	driver, err := pg.WithInstance(s.conn, &pg.Config{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://migration", "postgres", driver)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
