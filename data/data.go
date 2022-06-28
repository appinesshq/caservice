// Package data provides functionality for data interaction.
package data

import user "github.com/appinesshq/caservice/business/user/usecases"

type Repositories struct {
	UserRepo user.UserRepository
}
