package config

import (
	"fmt"
	"os"

	"github.com/ceralena/envconf"
)

type Config struct {
	Server_Addr          string
	Words_File           string
	Users_File           string
	Message_Queue_Length int
	Slack_Token          string
}

func (c Config) Validate() error {
	if _, err := os.Stat(c.Words_File); os.IsNotExist(err) {
		return fmt.Errorf("Bad config: words file does not exist: %s", c.Words_File)
	} else if err != nil {
		return err
	}
	if _, err := os.Stat(c.Users_File); os.IsNotExist(err) {
		return fmt.Errorf("Bad config: users file does not exist: %s", c.Users_File)
	} else if err != nil {
		return err
	}
	return nil
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	if err := envconf.ReadConfigEnv(&cfg); err != nil {
		return cfg, err
	}
	return cfg, cfg.Validate()
}
