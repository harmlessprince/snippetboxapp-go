package mysql

import (
	"errors"
	"github.com/harmlessprince/snippetboxapp/pkg/models"
	"reflect"
	"testing"
	"time"
)

func TestUserModel_Get(t *testing.T) {
	//skip the test if --short flag is provided when running the test
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}
	// Set up a suite of table-driven tests and expected results.
	tests := []struct {
		name      string
		userID    int
		wantUser  *models.User
		wantError error
	}{
		{
			name:   "Valid ID",
			userID: 1,
			wantUser: &models.User{
				ID:      1,
				Name:    "Alice Jones",
				Email:   "alice@example.com",
				Created: time.Date(2018, 12, 23, 17, 25, 22, 0, time.UTC),
			},
			wantError: nil,
		},
		{
			name:      "Zero ID",
			userID:    0,
			wantUser:  nil,
			wantError: models.ErrNoRecord,
		},
		{
			name:      "Non-existent ID",
			userID:    2,
			wantUser:  nil,
			wantError: models.ErrNoRecord,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			m := UserModel{db}
			user, err := m.Get(test.userID)
			if !errors.Is(err, test.wantError) {
				t.Errorf("want %v; got %s", test.wantError, err)
			}
			if !reflect.DeepEqual(user, test.wantUser) {
				t.Errorf("want %v; got %v", test.wantUser, user)
			}
		})
	}

}
