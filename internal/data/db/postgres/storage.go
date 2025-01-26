package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Jasurbek-Tursunov/warehouse/pkg/config"
	_ "github.com/lib/pq"
)

type Storage struct {
	conn    *sql.DB
	timeout time.Duration
}

func NewStorage() *Storage {
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
		panic(err)
	}
	return &Storage{conn: db, timeout: cfg.Timeout}
}

func (s *Storage) MustConnect() {
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
		panic(err)
	}
	s.conn = db
}

func (s *Storage) HealthCheck() {}

func (s *Storage) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}
