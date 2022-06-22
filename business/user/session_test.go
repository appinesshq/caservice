package user_test

import (
	"context"
	"testing"
	"time"

	user "github.com/appinesshq/caservice/business/user"
	"github.com/appinesshq/caservice/foundation/tests"
	"github.com/google/go-cmp/cmp"
)

var (
	past   = time.Date(2022, time.March, 3, 0, 0, 0, 0, time.UTC)
	now    = time.Now()
	future = time.Now().Add(1 * time.Hour)
)

func TestSessionEntity(t *testing.T) {
	t.Parallel()

	t.Log("Given the need to work with Session entities.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single Session.", testID)
		{
			t.Run("testUserRoles", testUserRoles)
			t.Run("testSessionInPastIsExpiredAndInvalid", testSessionInPastIsExpiredAndInvalid)
			t.Run("testEmptySessionIsExpiredAndInvalid", testEmptySessionIsExpiredAndInvalid)
			t.Run("testValidSession", testValidSession)
			t.Run("testSessionContext", testSessionContext)
		}
	}
}

func testUserRoles(t *testing.T) {
	u, _ := user.NewWithID("Test User", "my@email.com", "testpassword", []string{"USER"}, now)
	session := user.NewSession(u, now)

	if !session.UserHasRole("USER") {
		t.Fatalf("\t%s\tShould be authorized with user role.", tests.Failed)
	}
	t.Logf("\t%s\tShould be authorized with user role.", tests.Success)

	if session.UserHasRole("ADMIN") {
		t.Fatalf("\t%s\tShould not be authorized with admin role.", tests.Failed)
	}
	t.Logf("\t%s\tShould not be authorized with admin role.", tests.Success)
}

func testSessionInPastIsExpiredAndInvalid(t *testing.T) {
	u, _ := user.NewWithID("Test User", "my@email.com", "testpassword", []string{"USER"}, now)
	session := user.NewSession(u, past)

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
	session := user.Session{}

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
	u, _ := user.NewWithID("Test User", "my@email.com", "testpassword", []string{"USER"}, now)
	session := user.NewSession(u, future)

	if session.IsExpired() {
		t.Fatalf("\t%s\tShould not be expired.", tests.Failed)
	}
	t.Logf("\t%s\tShould not be expired.", tests.Success)

	if !session.IsValid() {
		t.Fatalf("\t%s\tShould be valid.", tests.Failed)
	}
	t.Logf("\t%s\tShould be valid.", tests.Success)
}

func testSessionContext(t *testing.T) {
	u, _ := user.NewWithID("Test User", "my@email.com", "testpassword", []string{"USER"}, now)
	s := user.NewSession(u, now)
	ctx := user.ContextWithSession(context.Background(), s)

	got, err := user.GetSession(ctx)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get session from context, but got: %v.", tests.Failed, err)
	}
	t.Logf("\t%s\tShould be able to get session from context.", tests.Success)

	if diff := cmp.Diff(got, s); diff != "" {
		t.Fatalf("\t%s\tShould get the expected result. Diff:\n%s", tests.Failed, diff)
	}
	t.Logf("\t%s\tShould get the expected result.", tests.Success)

	if _, err := user.GetSession(context.Background()); err == nil {
		t.Fatalf("\t%s\tShould not be able to get session from empty context.", tests.Failed)
	}
	t.Logf("\t%s\tShould not be able to get session from empty context.", tests.Success)

}
