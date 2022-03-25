package user_test

import (
	"strings"
	"testing"
	"time"

	"github.com/appinesshq/service/business/entity/user"
	"github.com/appinesshq/service/foundation/tests"
	"github.com/appinesshq/service/foundation/validation"
)

func TestUserEntity(t *testing.T) {
	t.Log("Given the need to work with User entities.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single User.", testID)
		{
			// ctx := context.Background()
			now := time.Date(2022, time.March, 3, 0, 0, 0, 0, time.UTC)

			//=====================================================================
			// Edge cases
			//=====================================================================

			// Invalid email
			_, err := user.New("notanemail", "testpassword", []string{"ADMIN"}, now)
			if err == nil {
				t.Fatalf("\t%s\tTest %d:\tShould not be able to create new user with wrong email: %v.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould not be able to create a new user with wrong email.", tests.Success, testID)

			if !validation.IsValidationErrors(err) {
				t.Fatalf("\t%s\tTest %d:\tShould get an error of type IsValidationErrors.", tests.Failed, testID)
			}
			t.Logf("\t%s\tTest %d:\tShould get an error of type IsValidationErrors.", tests.Success, testID)

			expected := "Email must be a valid email address"
			if !strings.Contains(err.Error(), expected) {
				t.Fatalf("\t%s\tTest %d:\tShould get an error containing %q, but got %q.", tests.Failed, testID, expected, err.Error())
			}
			t.Logf("\t%s\tTest %d:\tShould get an error containing %q.", tests.Success, testID, expected)

			// Invalid password
			_, err = user.New("my@email.com", "", []string{"ADMIN"}, now)
			if !validation.IsValidationErrors(err) {
				t.Fatalf("\t%s\tTest %d:\tShould not be able to create new user without a password: %v.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould not be able to create a new user without a password.", tests.Success, testID)

			//=====================================================================
			// Other cases
			//=====================================================================
			o, err := user.New("my@email.com", "testpassword", []string{"ADMIN"}, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create new user with valid details: %v.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create new user with valid details.", tests.Success, testID)

			if o.ID == "" {
				t.Fatalf("\t%s\tTest %d:\tShould have a generated id in the new user.", tests.Failed, testID)
			}
			t.Logf("\t%s\tTest %d:\tShould have a generated id in the new user.", tests.Success, testID)

			if o.Password == "" {
				t.Fatalf("\t%s\tTest %d:\tShould have a generated password in the new user.", tests.Failed, testID)
			}
			t.Logf("\t%s\tTest %d:\tShould have a generated password in the new user.", tests.Success, testID)

			if err := o.CheckPassword("testpassword"); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould pass check with valid password: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould pass check with valid password.", tests.Success, testID)

			if err := o.CheckPassword("invalidpassword"); err == nil {
				t.Fatalf("\t%s\tTest %d:\tShould not pass check with invalid password.", tests.Failed, testID)
			}
			t.Logf("\t%s\tTest %d:\tShould not pass check with invalid password.", tests.Success, testID)
		}
	}
}
