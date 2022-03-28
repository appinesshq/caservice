package user_test

import (
	"strings"
	"testing"
	"time"

	"github.com/appinesshq/caservice/business/entity/user"
	"github.com/appinesshq/caservice/foundation/tests"
	"github.com/appinesshq/caservice/foundation/validation"
)

var now = time.Date(2022, time.March, 3, 0, 0, 0, 0, time.UTC)

func TestUserEntity(t *testing.T) {
	t.Log("Given the need to work with User entities.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single User.", testID)
		{
			t.Run("invalidEmail", testNewUserWithInvalidEmailFails)
			t.Run("invalidPassword", testNewUserWithInvalidPasswordFails)
			t.Run("newUser", testNewUserSucceeds)
		}
	}
}

func testNewUserWithInvalidEmailFails(t *testing.T) {
	_, err := user.New("Test User", "notanemail", "testpassword", []string{"ADMIN"}, now)
	if err == nil {
		t.Fatalf("\t%s\tShould not be able to create new user with wrong email: %v.", tests.Failed, err)
	}
	t.Logf("\t%s\tShould not be able to create a new user with wrong email.", tests.Success)

	if !validation.IsValidationErrors(err) {
		t.Fatalf("\t%s\tShould get an error of type IsValidationErrors.", tests.Failed)
	}
	t.Logf("\t%s\tShould get an error of type IsValidationErrors.", tests.Success)

	expected := "Email must be a valid email address"
	if !strings.Contains(err.Error(), expected) {
		t.Fatalf("\t%s\tShould get an error containing %q, but got %q.", tests.Failed, expected, err.Error())
	}
	t.Logf("\t%s\tShould get an error containing %q.", tests.Success, expected)
}

func testNewUserWithInvalidPasswordFails(t *testing.T) {
	_, err := user.New("Test User", "my@email.com", "", []string{"ADMIN"}, now)
	if !validation.IsValidationErrors(err) {
		t.Fatalf("\t%s\tShould not be able to create new user without a password: %v.", tests.Failed, err)
	}
	t.Logf("\t%s\tShould not be able to create a new user without a password.", tests.Success)

	expected := "Password is a required field"
	if !strings.Contains(err.Error(), expected) {
		t.Fatalf("\t%s\tShould get an error containing %q, but got %q.", tests.Failed, expected, err.Error())
	}
	t.Logf("\t%s\tShould get an error containing %q.", tests.Success, expected)
}

func testNewUserSucceeds(t *testing.T) {
	o, err := user.New("Test User", "my@email.com", "testpassword", []string{"ADMIN"}, now)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to create new user with valid details: %v.", tests.Failed, err)
	}
	t.Logf("\t%s\tShould be able to create new user with valid details.", tests.Success)

	if o.ID == "" {
		t.Fatalf("\t%s\tShould have a generated id in the new user.", tests.Failed)
	}
	t.Logf("\t%s\tShould have a generated id in the new user.", tests.Success)

	if o.Password == "" {
		t.Fatalf("\t%s\tShould have a generated password in the new user.", tests.Failed)
	}
	t.Logf("\t%s\tShould have a generated password in the new user.", tests.Success)

	if err := o.CheckPassword("testpassword"); err != nil {
		t.Fatalf("\t%s\tShould pass check with valid password: %s.", tests.Failed, err)
	}
	t.Logf("\t%s\tShould pass check with valid password.", tests.Success)

	if err := o.CheckPassword("invalidpassword"); err == nil {
		t.Fatalf("\t%s\tShould not pass check with invalid password.", tests.Failed)
	}
	t.Logf("\t%s\tShould not pass check with invalid password.", tests.Success)
}
