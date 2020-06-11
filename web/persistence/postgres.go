package persistence

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/eldelto/solvent"
	"github.com/eldelto/solvent/web/dto"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(host, port, user, password string) (*PostgresRepository, error) {
	connectUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/solvent", user, password, host, port)
	db, err := sql.Open("pgx", connectUrl)
	if err != nil {
		return nil, err
	}

	repo := PostgresRepository{db: db}

	_, err = repo.db.Exec(`CREATE TABLE IF NOT EXISTS notebooks(
			id VARCHAR(36) PRIMARY KEY NOT NULL, 
			data JSONB NOT NULL
		)`)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &repo, nil
}

func (r *PostgresRepository) Close() {
	r.db.Close()
}

func (r *PostgresRepository) Store(notebook *solvent.Notebook) error {
	id := notebook.ID.String()
	data, err := notebookToJson(notebook)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("INSERT INTO notebooks VALUES($1, $2)", id, data)

	return err
}

func (r *PostgresRepository) Update(notebook *solvent.Notebook) error {
	id := notebook.ID.String()
	data, err := notebookToJson(notebook)
	if err != nil {
		return err
	}

	result, err := r.db.Exec("UPDATE notebooks SET data = $2 WHERE id = $1", id, data)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	} else if count <= 0 {
		return fmt.Errorf("Notebook with ID '%v' could not be found", notebook.ID)
	}

	return nil
}

func (r *PostgresRepository) Fetch(id uuid.UUID) (*solvent.Notebook, error) {
	var data []byte
	err := r.db.QueryRow("SELECT data FROM notebooks WHERE id = $1", id.String()).Scan(&data)
	if err != nil {
		return nil, err
	}

	var notebookDto dto.NotebookDto
	json.Unmarshal(data, &notebookDto)
	notebook := dto.NotebookFromDto(&notebookDto)

	return notebook, nil
}

func (r *PostgresRepository) Remove(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM notebooks WHERE id = $1", id.String())
	return err
}

func notebookToJson(notebook *solvent.Notebook) ([]byte, error) {
	dto := dto.NotebookToDto(notebook)
	data, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	return data, nil
}
