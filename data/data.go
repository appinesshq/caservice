// Package data provides functionality for data interaction.
package data

import "github.com/appinesshq/caservice/data/user"

type DataSources struct {
	UserRepo user.UserStorage
}
