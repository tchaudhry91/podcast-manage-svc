package podcastmg

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"testing"
)

func TestUserCreationAndPersistence(t *testing.T) {
	store.Connect()
	store.Migrate()
	defer store.Close()
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
		err = store.CreateUser(&user)
		if err != nil {
			t.Errorf("Error creating user in Databse:%v", err)
		}
		user = User{}
		user, err = store.GetUserFromEmail(testCase.email)
		if err != nil {
			t.Errorf("Error reading user from Database:%v", err)
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
