package psg

import (
	"HW-1/pkg"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Psg struct {
	conn *pgxpool.Pool
}

func NewPsg(ctx context.Context, cfg *pkg.Config) (*Psg, error) {
	pool, err := pgxpool.New(ctx, generateDSN(cfg))
	if err != nil {
		return nil, err
	}
	sql := `CREATE TABLE
    IF NOT EXISTS address_book(
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL DEFAULT '',
        last_name TEXT NOT NULL DEFAULT '',
        middle_name TEXT NOT NULL DEFAULT '',
        address TEXT NOT NULL DEFAULT '',
        phone TEXT NOT NULL DEFAULT ''
    );`
	pool.Exec(ctx, sql)

	return &Psg{conn: pool}, nil
}

func generateDSN(cfg *pkg.Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
}

func (p *Psg) Close() {
	p.conn.Close()
}
