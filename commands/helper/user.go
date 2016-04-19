package user

import (
	"fmt"
	"io/ioutil"

	"github.com/HPI-BP2015H/go-travis/client"
)

// User from Travis
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
}

// CurrentUser returns the user currently logged in into Travis
func CurrentUser(client *client.Client) (User, error) {
	user := User{}
	res, err := client.PerformAction("user", "current", map[string]string{})
	if err != nil {
		return user, fmt.Errorf("Error: Could not get current user! \n%s", err.Error())
	}
	if res.StatusCode > 299 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return user, err
		}
		return user, fmt.Errorf("Unexpected HTTP status: %d\n%s\n", res.StatusCode, string(body))
	}
	defer res.Body.Close()
	res.Unmarshal(&user)
	return user, nil
}

func (u User) String() string {
	if u.Name != "" && u.Login != u.Name {
		return fmt.Sprintf("%s (%s)", u.Login, u.Name)
	}
	return u.Login
}
