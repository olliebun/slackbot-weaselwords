package weaselbot

import (
	"bufio"
	"io"
)

// List of usernames that weaselbot is enabled for.
type Users []string

func (users Users) Matches(username string) bool {
	for _, user := range users {
		if user == username {
			return true
		}
	}
	return false
}

// Load line-delimited words from a stream.
func UsersFromReader(input io.Reader) (Users, error) {
	users := make(Users, 0)
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		users = append(users, scanner.Text())
	}

	return users, scanner.Err()
}
