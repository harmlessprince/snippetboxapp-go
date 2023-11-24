package mock

import (
	"github.com/harmlessprince/snippetboxapp/pkg/models"
	"time"
)

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "The boy and the king",
	Content: "The boy and the king",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct {
}

func (sm *SnippetModel) Insert(title, content, expires string) (int, error) {
	return 2, nil
}

func (sm *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}
func (sm *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
