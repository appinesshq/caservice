package user_test

import (
	"errors"
	"testing"

	// "time"
	"github.com/appinesshq/caservice/business/user"
	"github.com/appinesshq/caservice/foundation/tests"
	"github.com/appinesshq/caservice/foundation/validation"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

var (
	id       = uuid.New().String()
	name     = "Test User"
	email    = "test@example.com"
	password = "password"
	roles    = []string{"USER"}
)

func TestUserEntity(t *testing.T) {
	t.Parallel()

	t.Log("Given the need to work with User entities.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single User.", testID)
		{
			t.Run("testEmptyUser", testEmptyUser)
			t.Run("testInvalidUser", testInvalidUser)
			t.Run("testValidUser", testValidUser)
		}
	}
}

func testEmptyUser(t *testing.T) {
	err := user.User{}.Validate()

	var got validation.ValidationError
	if !errors.As(err, &got) {
		t.Fatalf("\t%s\tShould get a validation.ValidationError, but got %T.", tests.Failed, err)
	}
	t.Logf("\t%s\tShould get a validation.ValidationError.", tests.Success)

	exp := validation.ValidationError{
		Err: "data validation error",
		Fields: map[string]string{
			"ID":           "ID is a required field",
			"Email":        "Email is a required field",
			"Name":         "Name is a required field",
			"PasswordHash": "PasswordHash is a required field",
			"Roles":        "Roles is a required field",
			"DateCreated":  "DateCreated is a required field",
			"DateUpdated":  "DateUpdated is a required field"}}

	if diff := cmp.Diff(got, exp); diff != "" {
		t.Fatalf("\t%s\tShould get the expected result. Diff:\n%s", tests.Failed, diff)
	}
	t.Logf("\t%s\tShould get the expected result.", tests.Success)
}

func testInvalidUser(t *testing.T) {
	_, err := user.New("123", name, "invalid", "", roles, now)

	var got validation.ValidationError
	if !errors.As(err, &got) {
		t.Fatalf("\t%s\tShould get a validation.ValidationError, but got %T.", tests.Failed, err)
	}
	t.Logf("\t%s\tShould get a validation.ValidationError.", tests.Success)

	exp := validation.ValidationError{
		Err: "data validation error",
		Fields: map[string]string{
			"ID":           "ID must be a valid UUID",
			"Email":        "Email must be a valid email address",
			"PasswordHash": "PasswordHash can't be an empty string",
		}}

	if diff := cmp.Diff(got, exp); diff != "" {
		t.Fatalf("\t%s\tShould get the expected result. Diff:\n%s", tests.Failed, diff)
	}
	t.Logf("\t%s\tShould get the expected result.", tests.Success)
}

func testValidUser(t *testing.T) {
	got, err := user.New(id, name, email, password, roles, now)
	if err != nil {
		t.Fatalf("\t%s\tShould not be able to create a user with valid data.", tests.Failed)
	}
	t.Logf("\t%s\tShould not be able to create a user with valid data.", tests.Success)

	exp := got
	exp.ID = id
	exp.Name = name
	exp.Email = email
	exp.Roles = roles
	exp.DateCreated = now
	exp.DateUpdated = now

	if diff := cmp.Diff(got, exp); diff != "" {
		t.Fatalf("\t%s\tShould get the expected result. Diff:\n%s", tests.Failed, diff)
	}
	t.Logf("\t%s\tShould get the expected result.", tests.Success)
}
