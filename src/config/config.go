package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Config struct {
	Port         int          `yaml:"port"`
	Path         string       `yaml:"path"`
	SourceUrl    string       `yaml:"sourceUrl"`
	MongoUrl     string       `yaml:"mongoUrl"`
	Auth         AuthOptions  `yaml:"auth"`
	CacheOptions CacheOptions `yaml:"cacheOptions"`
}

type CacheOptions struct {
	Expiration int `yaml:"expiration"`
	Purge      int `yaml:"purge"`
}

type AuthOptions struct {
	Mode          string        `yaml:"mode"`
	JwtOptions    JwtOptions    `yaml:"jwtOptions"`
	HeaderOptions HeaderOptions `yaml:"headerOptions"`
}

type HeaderOptions struct {
	Name string `yaml:"name"`
}

type JwtOptions struct {
	SigningMethod string `yaml:"signingMethod"`
	JwksUrl       string `yaml:"jwksUrl"`
	Key           string `yaml:"key"`
	KeyPath       string `yaml:"keyPath"`
	KeyUrl        string `yaml:"keyUrl"`
	RoleClaim     string `yaml:"roleClaim"`
	AllowedSub    string `yaml:"allowedSub"`
	AllowedAud    string `yaml:"allowedAud"`
	MaxAgeSec     int64  `yaml:"maxAgeSec"`
}

func (c *Config) evaluateAndFillDefaults() error {
	if c.Port <= 0 {
		c.Port = 8080
	}
	if c.Path == "" {
		c.Path = "/graphql"
	}
	if c.SourceUrl == "" {
		return errors.New("no sourceUrl provided in config")
	}
	if c.MongoUrl == "" {
		return errors.New("no mongoUrl provided in config")
	}
	if c.Auth.Mode == "" {
		return errors.New("no auth provided in config")
	}
	if err := c.CacheOptions.evaluateAndFillDefaults(); err != nil {
		return err
	}
	switch c.Auth.Mode {
	case "jwt":
		err := c.Auth.JwtOptions.evaluateAndFillDefaults()
		if err != nil {
			return err
		}
		break
	case "header":
		err := c.Auth.HeaderOptions.evaluateAndFillDefaults()
		if err != nil {
			return err
		}
		break
	default:
		return errors.New("unknown auth mode provided")
	}
	return nil
}

func (c *JwtOptions) evaluateAndFillDefaults() error {
	if c.JwksUrl == "" && c.Key == "" && c.KeyPath == "" {
		return errors.New("none of key, keyPath, jwksUrl provided in config")
	}
	if c.RoleClaim == "" {
		return errors.New("no roleClaim provided in config")
	}
	return nil
}

func (c *HeaderOptions) evaluateAndFillDefaults() error {
	if c.Name == "" {
		return errors.New("no header name provided in options")
	}
	return nil
}

func (c *CacheOptions) evaluateAndFillDefaults() error {
	if c.Expiration <= 0 {
		c.Expiration = 5
	}
	if c.Purge <= 0 {
		c.Purge = 10
	}
	return nil
}

func GetConfig(path string) (Config, error) {
	var res Config

	file, err := os.Open(path)
	if err != nil {
		return res, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return res, err
	}

	err = yaml.Unmarshal(bytes, &res)
	if err != nil {
		return res, err
	}

	err = res.evaluateAndFillDefaults()
	return res, err
}