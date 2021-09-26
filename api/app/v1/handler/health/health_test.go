package health_test

import (
	"context"
	"io/ioutil"
	v1 "jokes-bapak2-api/app/v1"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/assert"
)

var jokesData = []interface{}{1, "https://via.placeholder.com/300/06f/fff.png", 1, 2, "https://via.placeholder.com/300/07f/fff.png", 1, 3, "https://via.placeholder.com/300/08f/fff.png", 1}
var app *fiber.App = v1.New()
var db *pgxpool.Pool

func cleanup() {
	poolConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln("Unable to create pool config", err)
	}
	poolConfig.MaxConns = 5
	poolConfig.MinConns = 2

	db, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalln("Unable to create connection", err)
	}

	j, err := db.Query(context.Background(), "DROP TABLE \"jokesbapak2\"")
	if err != nil {
		panic(err)
	}
	defer j.Close()

	a, err := db.Query(context.Background(), "DROP TABLE \"administrators\"")
	if err != nil {
		panic(err)
	}

	defer a.Close()
}

func setup() error {
	poolConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln("Unable to create pool config", err)
	}
	poolConfig.MaxConns = 5
	poolConfig.MinConns = 2

	db, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalln("Unable to create connection", err)
	}

	a, err := db.Query(context.Background(), "INSERT INTO \"administrators\" (id, key, token, last_used) VALUES ($1, $2, $3, $4);", 1, "very secure", "not the real one", time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}
	defer a.Close()

	j, err := db.Query(context.Background(), "INSERT INTO \"jokesbapak2\" (id, link, creator) VALUES ($1, $2, $3), ($4, $5, $6), ($7, $8, $9);", jokesData...)
	if err != nil {
		return err
	}

	defer j.Close()

	return nil
}

func TestHealth(t *testing.T) {
	err := setup()
	if err != nil {
		t.Fatal(err)
	}

	defer cleanup()

	req, _ := http.NewRequest("GET", "/health", nil)
	res, err := app.Test(req, int(time.Minute * 2))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equalf(t, false, err != nil, "health")
	assert.Equalf(t, 200, res.StatusCode, "health")
	assert.NotEqualf(t, 0, res.ContentLength, "health")
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	
	assert.Nilf(t, err, "health")
}
