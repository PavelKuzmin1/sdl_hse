package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"os"
)

type Config struct {
	Host string `yaml:"Host"`
	Port string `yaml:"Port"`
	Name string `yaml:"Name"`
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.File("index.html")
	})
	e.POST("/connect", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		config, err := loadConfig("config.yaml")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to load config: %v\n", err)
			os.Exit(1)
		}

		urlExample := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			username, password, config.Host, config.Port, config.Name)

		conn, err := pgx.Connect(context.Background(), urlExample)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to connect to database: %v", err))
		}
		defer conn.Close(context.Background())

		var version string
		err = conn.QueryRow(context.Background(), "SELECT VERSION();").Scan(&version)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("QueryRow failed: %v", err))
		}

		return c.String(http.StatusOK, version)
	})

	e.Start(":1234")
}

func loadConfig(filename string) (Config, error) {
	var config Config
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(fileContent, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
