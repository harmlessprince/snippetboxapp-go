package mysql

import (
	"database/sql"
	"errors"
	"github.com/harmlessprince/snippetboxapp/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (sm *SnippetModel) Insert(title, content, expires string) (int, error) {
	statement := `
	INSERT INTO snippets (title, content,created, expires) 
	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL  ? DAY ))
`
	result, err := sm.DB.Exec(statement, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (sm *SnippetModel) Get(id int) (*models.Snippet, error) {
	statement := `SELECT * FROM snippets WHERE  expires > UTC_TIMESTAMP() AND id = ?`
	row := sm.DB.QueryRow(statement, id)
	snippet := &models.Snippet{}
	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return snippet, nil
}
func (sm *SnippetModel) Latest() ([]*models.Snippet, error) {
	statement := `SELECT * FROM snippets WHERE  expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	rows, err := sm.DB.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var snippets []*models.Snippet
	for rows.Next() {
		snippet := &models.Snippet{}
		err := rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, snippet)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
