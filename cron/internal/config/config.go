package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	SMTP     SMTPConfig     `yaml:"smtp"`
	Reminder ReminderConfig `yaml:"reminder"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}

type ReminderConfig struct {
	Recipients []string `yaml:"recipients"`
	CronSpec   string   `yaml:"cron_spec"`
	Timezone   string   `yaml:"timezone"`
}

func Load(path string) (*Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return nil, err
	}

	cfg.applyDefaults()
	return &cfg, nil
}

func (c *Config) applyDefaults() {
	if strings.TrimSpace(c.Server.Port) == "" {
		c.Server.Port = "8080"
	}

	if c.SMTP.Port == 0 {
		c.SMTP.Port = 587
	}

	if strings.TrimSpace(c.Reminder.CronSpec) == "" {
		c.Reminder.CronSpec = "0 8 * * *"
	}

	if strings.TrimSpace(c.Reminder.Timezone) == "" {
		c.Reminder.Timezone = "Asia/Shanghai"
	}
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.Database.DSN) == "" {
		return errors.New("database.dsn is required")
	}

	if strings.TrimSpace(c.SMTP.Host) == "" {
		return errors.New("smtp.host is required")
	}

	if c.SMTP.Port <= 0 {
		return errors.New("smtp.port must be greater than 0")
	}

	if strings.TrimSpace(c.SMTP.Username) == "" {
		return errors.New("smtp.username is required")
	}

	if strings.TrimSpace(c.SMTP.Password) == "" {
		return errors.New("smtp.password is required")
	}

	if strings.TrimSpace(c.SMTP.From) == "" {
		return errors.New("smtp.from is required")
	}

	if len(c.Reminder.Recipients) == 0 {
		return errors.New("reminder.recipients is required")
	}

	for index, recipient := range c.Reminder.Recipients {
		if strings.TrimSpace(recipient) == "" {
			return fmt.Errorf("reminder.recipients[%d] cannot be empty", index)
		}
	}

	if strings.TrimSpace(c.Reminder.CronSpec) == "" {
		return errors.New("reminder.cron_spec is required")
	}

	if strings.TrimSpace(c.Reminder.Timezone) == "" {
		return errors.New("reminder.timezone is required")
	}

	return nil
}
