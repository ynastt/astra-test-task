package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Структура, которая хранит данные уязвимости
type Warning struct {
	Id        int
	RuleId    string
	Uri       string
	StartLine int
	XSeverity int
}

// Структура, которая хранит информацию о подключении к базе данных.
type Postgres struct {
	db *pgxpool.Pool
}

// функция-конструктор для структуры Postgres
func NewPGXPool(db *pgxpool.Pool) Postgres {
	return Postgres{
		db: db,
	}
}

// Функция подключения к бд
func ConnectDB(dbConnURL string) (*pgxpool.Pool, error) {
	// Подключаемся к базе данных
	dbpool, err := pgxpool.New(context.Background(), dbConnURL)
	if err != nil {
		return nil, err
	}

	// Проверяем соединение с помощью пинга
	err = dbpool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return dbpool, nil
}

// Метод создания таблицы warnings
func (p *Postgres) CreateWarningsTable() error {
	createquery := `CREATE TABLE IF NOT EXISTS warnings (
    warning_id SERIAL PRIMARY KEY, 
    ruleId text NOT NULL,
    uri text NOT NULL,
    startLine int NOT NULL CHECK (startLine >= 0),
    xseverity int NOT NULL CHECK (xseverity >= 0 AND xseverity <= 2));`

	_, err := p.db.Exec(context.Background(), createquery)
	return err
}

// Метод вставки строки в таблицу warnings
func (p *Postgres) InsertWarning(w *Warning) error {
	query := `INSERT INTO warnings (ruleId, uri, startLine, xseverity) VALUES (@ruleId, @uri, @startLine, @xseverity)`
	args := pgx.NamedArgs{
		"ruleId":    &w.RuleId,
		"uri":       &w.Uri,
		"startLine": &w.StartLine,
		"xseverity": &w.XSeverity,
	}
	_, err := p.db.Exec(context.Background(), query, args)
	return err
}
