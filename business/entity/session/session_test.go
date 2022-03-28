package session_test

import (
	"testing"
	"time"

	"github.com/appinesshq/caservice/business/entity/session"
	"github.com/appinesshq/caservice/business/entity/user"
	"github.com/appinesshq/caservice/foundation/tests"
)

var (
	past   = time.Date(2022, time.March, 3, 0, 0, 0, 0, time.UTC)
	now    = time.Now()
	future = time.Now().Add(1 * time.Hour)
)

func TestSessionEntity(t *testing.T) {
	t.Log("Given the need to work with Session entities.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single Session.", testID)
		{
			t.Run("testUserRoles", testUserRoles)
			t.Run("testSessionInPastIsExpiredAndInvalid", testSessionInPastIsExpiredAndInvalid)
			t.Run("testEmptySessionIsExpiredAndInvalid", testEmptySessionIsExpiredAndInvalid)
		}
	}
}

func testUserRoles(t *testing.T) {
	u, _ := user.New("Test User", "my@email.com", "testpassword", []user.Role{user.RoleUser}, now)
	session := session.New(u, now)

	if !session.UserHasRole(user.RoleUser) {
		t.Fatalf("\t%s\tShould be authorized with user role.", tests.Failed)
	}
	t.Logf("\t%s\tShould be authorized with user role.", tests.Success)

	if session.UserHasRole(user.RoleAdmin) {
		t.Fatalf("\t%s\tShould not be authorized with admin role.", tests.Failed)
	}
	t.Logf("\t%s\tShould not be authorized with admin role.", tests.Success)
}

func testSessionInPastIsExpiredAndInvalid(t *testing.T) {
	u, _ := user.New("Test User", "my@email.com", "testpassword", []user.Role{user.RoleUser}, now)
	session := session.New(u, past)

	if !session.IsExpired() {
		t.Fatalf("\t%s\tShould be expired.", tests.Failed)
	}
	t.Logf("\t%s\tShould be expired.", tests.Success)

	if session.IsValid() {
		t.Fatalf("\t%s\tShould not be valid.", tests.Failed)
	}
	t.Logf("\t%s\tShould not be valid.", tests.Success)
}

func testEmptySessionIsExpiredAndInvalid(t *testing.T) {
	session := session.Session{}

	if !session.IsExpired() {
		t.Fatalf("\t%s\tShould be expired.", tests.Failed)
	}
	t.Logf("\t%s\tShould be expired.", tests.Success)

	if session.IsValid() {
		t.Fatalf("\t%s\tShould not be valid.", tests.Failed)
	}
	t.Logf("\t%s\tShould not be valid.", tests.Success)
}

func testValidSession(t *testing.T) {
	u, _ := user.New("Test User", "my@email.com", "testpassword", []user.Role{user.RoleUser}, now)
	session := session.New(u, future)

	if session.IsExpired() {
		t.Fatalf("\t%s\tShould not be expired.", tests.Failed)
	}
	t.Logf("\t%s\tShould not be expired.", tests.Success)

	if !session.IsValid() {
		t.Fatalf("\t%s\tShould be valid.", tests.Failed)
	}
	t.Logf("\t%s\tShould be valid.", tests.Success)
}
