package user_test

import (
	"strings"
	"testing"
	"time"

	entity "github.com/appinesshq/caservice/business/entity/user"
	"github.com/appinesshq/caservice/business/usecase/user"
	"github.com/appinesshq/caservice/business/usecase/user/mock"
	"github.com/appinesshq/caservice/foundation/tests"
	"github.com/appinesshq/caservice/foundation/validation"
	"github.com/golang/mock/gomock"
)

var now = time.Date(2022, time.March, 3, 0, 0, 0, 0, time.UTC)

func setupMockRepo(t *testing.T) (*mock.MockRepository, func()) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockRepository(ctrl)
	return repo, ctrl.Finish
}

func TestUserUseCase(t *testing.T) {
	t.Log("Given the need to work with User usecases.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single User.", testID)
		{
			t.Run("registerUserWithInvalidEmailFails", testRegisterUserWithInvalidEmailFails)
			t.Run("registerValidUserSucceeds", testRegisterValidUserSucceeds)
			t.Run("testAuthenticationWithWrongDetailsFails", testAuthenticationWithInvalidCredentialsFails)
			t.Run("testAuthenticationWithValidCredentialsSucceeds", testAuthenticationWithValidCredentialsSucceeds)
		}
	}
}

func testRegisterUserWithInvalidEmailFails(t *testing.T) {
	uc := user.UseCase{Repo: nil}

	_, err := uc.Register("Test User", "notanemail", "testpassword", []string{"ADMIN"})
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

func testRegisterValidUserSucceeds(t *testing.T) {
	repo, teardown := setupMockRepo(t)
	defer teardown()
	uc := user.UseCase{Repo: repo}

	repo.
		EXPECT().
		Create(gomock.Any()).
		Return(nil)

	o, err := uc.Register("Test User", "my@email.com", "testpassword", []string{"ADMIN"})
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

func testAuthenticationWithInvalidCredentialsFails(t *testing.T) {
	repo, teardown := setupMockRepo(t)
	defer teardown()
	uc := user.UseCase{Repo: repo}

	// Test authentication with an invalid non-existing email address.
	repo.
		EXPECT().
		GetByEmail(gomock.Eq("nonexisting@email.com")).
		Return(nil, user.ErrNotFound)

	_, err := uc.Authenticate("nonexisting@email.com", "somepassword")
	if err == nil {
		t.Fatalf("\t%s\tShould not be able to authenticate with a non-existing email.", tests.Failed)
	}
	t.Logf("\t%s\tShould not be able to authenticate with a non-existing email.", tests.Success)

	if err != user.ErrAuthenticationFailed {
		t.Fatalf("\t%s\tShould get ErrAuthenticationFailed.", tests.Failed)
	}
	t.Logf("\t%s\tShould get ErrAuthenticationFailed.", tests.Success)

	// Test authentication with an wrong password fails.
	repo.
		EXPECT().
		GetByEmail(gomock.Eq("my@email.com")).
		Return(&entity.User{
			ID:       "5cf37266-3473-4006-984f-9325122678b7",
			Name:     "Test User",
			Email:    "my@email.com",
			Password: "$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a",
			Roles:    []string{"ADMIN", "USER"},
			Created:  now,
			Modified: now,
		}, nil)

	_, err = uc.Authenticate("my@email.com", "wrongpassword")
	if err == nil {
		t.Fatalf("\t%s\tShould not be able to authenticate with a wrong password.", tests.Failed)
	}
	t.Logf("\t%s\tShould not be able to authenticate with a wrong password.", tests.Success)

	if err != user.ErrAuthenticationFailed {
		t.Fatalf("\t%s\tShould get ErrAuthenticationFailed.", tests.Failed)
	}
	t.Logf("\t%s\tShould get ErrAuthenticationFailed.", tests.Success)
}

func testAuthenticationWithValidCredentialsSucceeds(t *testing.T) {
	repo, teardown := setupMockRepo(t)
	defer teardown()
	uc := user.UseCase{Repo: repo}

	repo.
		EXPECT().
		GetByEmail(gomock.Eq("my@email.com")).
		Return(&entity.User{
			ID:       "5cf37266-3473-4006-984f-9325122678b7",
			Name:     "Test User",
			Email:    "my@email.com",
			Password: "$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a",
			Roles:    []string{"ADMIN", "USER"},
			Created:  now,
			Modified: now,
		}, nil)

	o, err := uc.Authenticate("my@email.com", "gophers")
	if err != nil {
		t.Fatalf("\t%s\tShould be able to authenticate with valid credentials.", tests.Failed)
	}
	t.Logf("\t%s\tShould be able to authenticate with valid credentials.", tests.Success)

	if len(o) == 0 {
		t.Fatalf("\t%s\tShould have received a valid token.", tests.Failed)
	}
	t.Logf("\t%s\tShould have received a valid token.", tests.Success)
}
