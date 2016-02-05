package main

import (
	"fmt"
	"os"

	"weaselbot"
	"weaselbot/config"
	"weaselbot/server"
)

func fail(message string) {
	fmt.Fprintf(os.Stderr, message+"\n")
	os.Exit(1)
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fail("Failed to load config: " + err.Error())
	}

	// Load the words
	f, err := os.Open(cfg.Words_File)
	if err != nil {
		fail("Failed to open words file: " + err.Error())
	}
	defer f.Close()

	words, err := weaselbot.WordsFromReader(f)

	// Load the users
	f, err = os.Open(cfg.Users_File)
	if err != nil {
		fail("Failed to open users file: " + err.Error())
	}
	defer f.Close()

	users, err := weaselbot.UsersFromReader(f)

	srv := server.NewServer(cfg, words, users)
	if err := srv.Run(); err != nil {
		fail("Failed to run server: " + err.Error())
	}
}
