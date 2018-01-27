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
		email    string
		password string
		admin    bool
		err      bool
	}
	testCases := []userTestCase{
		{"tc@test.com", "what", false, false},
		{"bg@example.com", "should", false, false},
		{"tt32@ghh.com", "samplepassword", false, true},
		{"", "", false, true},
	}
	for _, testCase := range testCases {
		user, err := NewUser(testCase.email, testCase.password)
		if err != nil {
			if !testCase.err {
				t.Errorf("Error in creating user:%s: %v", testCase.email, err)
			} else {
				continue
			}
		}
		err = store.CreateUser(&user)
		if err != nil {
			t.Errorf("Error creating user in Databse:%v", err)
		}
		user = User{}
		user, err = store.GetUserByEmail(testCase.email)
		if err != nil {
			t.Errorf("Error reading user from Database:%v", err)
		}
		if user.UserEmail != testCase.email {
			t.Errorf("Email Want:%s\t Have:%s", testCase.email, user.UserEmail)
		}
		if user.admin != testCase.admin {
			if user.admin {
				t.Errorf("Incorrect admin Flag:%v", testCase)
			}
		}
	}

	// Test for primary key violation
	t.Run("Duplicate Email", func(t *testing.T) {
		user, err := NewUser(testCases[0].email, testCases[0].password)
		err = store.CreateUser(&user)
		if err == nil {
			t.Errorf("Should have errored in duplicate user creation")
		}
	})

	// Test for userUpdation
	t.Run("UserUpdation", func(t *testing.T) {
		user, err := store.GetUserByEmail(testCases[0].email)
		user.UserEmail = "new@shiny.com"
		err = store.UpdateUser(&user)
		if err != nil {
			t.Errorf("Failed to update user: %v", err)
		}
		user = User{}
		if user, err = store.GetUserByEmail("new@shiny.com"); err != nil {
			t.Errorf("Failed to get updated user: %v", err)
		}
	})

	// Test for deletion
	t.Run("UserDeletion", func(t *testing.T) {
		err := store.DeleteUserByEmail(testCases[2].email)
		if err != nil {
			t.Errorf("Failed to delete user:%v", err)
		}
		_, err = store.GetUserByEmail(testCases[2].email)
		if err == nil {
			t.Errorf("User not deleted succesfully, still found in Database")
		}
		// Retry deletion
		err = store.DeleteUserByEmail(testCases[2].email)
		if err == nil {
			t.Errorf("No error on re-deleting user, should error")
		}
	})

	// Test for non-existent user
	t.Run("Non Existent User", func(t *testing.T) {
		_, err := store.GetUserByEmail(testCases[0].email)
		if err == nil {
			t.Errorf("Should have errored and returned nil user on non existent user")
		}
	})
}
