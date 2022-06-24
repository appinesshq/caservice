package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// Authenticate authenticates a user at the service and returns a JWT token.
func Authenticate(host string, email string, password string, file string) error {
	if email == "" || password == "" {
		fmt.Println("help: authenticate <email> <password>")
		return ErrHelp
	}

	url := host + "/v1/users/token"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.SetBasicAuth(email, password)

	res, err := new(http.Client).Do(req)
	if err != nil {
		return fmt.Errorf("sending request to %q: %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		msg := struct {
			Error string `json:"error"`
		}{}
		if err := json.NewDecoder(res.Body).Decode(&msg); err != nil {
			return fmt.Errorf("request failed: %d %s", res.StatusCode, res.Status)
		}
		return fmt.Errorf("request failed: %s", msg.Error)
	}

	a := struct {
		Token string `json:"token"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(&a); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	if file == "" {
		fmt.Printf("user has authenticated succesfully, token: %q\n", a.Token)
		return nil
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(file), 0700) // Create your file
	}

	if err := ioutil.WriteFile(file, []byte(a.Token), 0600); err != nil {
		return fmt.Errorf("writing token file: %w", err)
	}

	fmt.Printf("user has authenticated succesfully, token is stored in: %q\n", file)

	return nil
}
