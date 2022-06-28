package pg

import (
	"time"
	"unsafe"

	"github.com/appinesshq/caservice/business/user"
	"github.com/lib/pq"
)

// User represent the structure we need for moving data
// between the app and the database.
type User struct {
	ID           string         `db:"user_id"`
	Name         string         `db:"name"`
	Email        string         `db:"email"`
	Roles        pq.StringArray `db:"roles"`
	PasswordHash []byte         `db:"password_hash"`
	DateCreated  time.Time      `db:"date_created"`
	DateUpdated  time.Time      `db:"date_updated"`
}

// =============================================================================

func toEntity(dbUsr User) user.User {
	pu := (*user.User)(unsafe.Pointer(&dbUsr))
	return *pu
}

func toEntitySlice(dbUsrs []User) []user.User {
	users := make([]user.User, len(dbUsrs))
	for i, dbUsr := range dbUsrs {
		users[i] = toEntity(dbUsr)
	}
	return users
}

func toUser(u user.User) User {
	pu := (*User)(unsafe.Pointer(&u))
	return *pu
}

// func toUserSlice(usrs []user.User) []User {
// 	users := make([]User, len(usrs))
// 	for i, u := range usrs {
// 		users[i] = toUser(u)
// 	}
// 	return users
// }
