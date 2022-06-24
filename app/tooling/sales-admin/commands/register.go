package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	entity "github.com/appinesshq/caservice/business/user"
	user "github.com/appinesshq/caservice/business/user/usecases"
)

// Register registers a new user in the system.
func Register(host, name, email, password string) error {
	if name == "" || email == "" || password == "" {
		fmt.Println("help: register <name> <email> <password>")
		return ErrHelp
	}
	nu := user.NewUser{
		Name:            name,
		Email:           email,
		Password:        password,
		PasswordConfirm: password,
	}

	b := bytes.Buffer{}
	if err := json.NewEncoder(&b).Encode(nu); err != nil {
		return fmt.Errorf("json encoding payload: %w", err)
	}

	url := host + "/v1/users/register"

	req, err := http.NewRequest(http.MethodPost, url, &b)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	res, err := new(http.Client).Do(req)
	if err != nil {
		return fmt.Errorf("sending request to %q: %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		msg := struct {
			Error string `json:"error"`
		}{}
		if err := json.NewDecoder(res.Body).Decode(&msg); err != nil {
			return fmt.Errorf("request failed: %d %s", res.StatusCode, res.Status)
		}
		return fmt.Errorf("request failed: %s", msg.Error)
	}

	u := entity.User{}
	if err := json.NewDecoder(res.Body).Decode(&u); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	fmt.Printf("new user registered with id %q and roles %s\n", u.ID, strings.Join(u.Roles, ", "))
	return nil
}
