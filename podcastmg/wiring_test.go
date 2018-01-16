package podcastmg

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"testing"
)

func init() {
	store = DBStore{
		*dbDialect,
		*dbConnectionString,
		nil,
	}
	if *dbDialect == "sqlite3" {
		os.Remove(*dbConnectionString)
	}
}

func TestUserCreationAndPersistence(t *testing.T) {
	type userTestCase struct {
		email string
		admin bool
		err   bool
	}
	testCases := []userTestCase{
		{"tc@test.com", false, false},
		{"bg@example.com", false, false},
		{"tt32@ghh.com", false, false},
		{"", false, true},
		{"admin@tt.com", true, false},
	}
	for _, testCase := range testCases {
		user, err := NewUser(testCase.email, testCase.admin)
		if err != nil {
			if !testCase.err {
				t.Errorf("Error in creating user: %v", err)
			} else {
				continue
			}
		}
		if user.UserEmail != testCase.email {
			t.Errorf("Email Want:%s\t Have:%s", testCase.email, user.UserEmail)
		}
		if user.Admin != testCase.admin {
			if user.Admin {
				t.Errorf("Incorrect Admin Flag:%v", testCase)
			}
		}
	}
}
