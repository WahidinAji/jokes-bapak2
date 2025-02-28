package joke_test

import (
	"context"
	"encoding/json"
	"jokes-bapak2-api/core/joke"
	"jokes-bapak2-api/core/schema"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
)

func TestGetAllJSONJokes(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()

	defer Flush()

	conn, err := db.Acquire(ctx)
	if err != nil {
		t.Error("an error was thrown:", err)
	}
	defer conn.Release()

	err = conn.BeginFunc(ctx, func(t pgx.Tx) error {
		_, err := t.Exec(
			ctx,
			`INSERT INTO "administrators"
				(id, key, token, last_used)
				VALUES
				($1, $2, $3, $4),
				($5, $6, $7, $8)`,
			administratorsData...,
		)
		if err != nil {
			return err
		}
		_, err = t.Exec(
			ctx,
			`INSERT INTO "jokesbapak2" 
				(id, link, creator)
				VALUES
				($1, $2, $3),
				($4, $5, $6),
				($7, $8, $9)`,
			jokesData...,
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	j, err := joke.GetAllJSONJokes(db, ctx)
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	if string(j) == "" {
		t.Error("j should not be empty")
	}
}

func TestGetRandomJokeFromDB(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()

	defer Flush()

	conn, err := db.Acquire(ctx)
	if err != nil {
		t.Error("an error was thrown:", err)
	}
	defer conn.Release()

	err = conn.BeginFunc(ctx, func(t pgx.Tx) error {
		_, err := t.Exec(
			ctx,
			`INSERT INTO "administrators"
				(id, key, token, last_used)
				VALUES
				($1, $2, $3, $4),
				($5, $6, $7, $8)`,
			administratorsData...,
		)
		if err != nil {
			return err
		}
		_, err = t.Exec(
			ctx,
			`INSERT INTO "jokesbapak2" 
				(id, link, creator)
				VALUES
				($1, $2, $3),
				($4, $5, $6),
				($7, $8, $9)`,
			jokesData...,
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	j, err := joke.GetRandomJokeFromDB(db, ctx)
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	if j == "" {
		t.Error("j should not be empty")
	}
}

func TestGetRandomJokeFromCache(t *testing.T) {
	defer Flush()

	jokes := []schema.Joke{
		{ID: 1, Link: "link1", Creator: 1},
		{ID: 2, Link: "link2", Creator: 1},
		{ID: 3, Link: "link3", Creator: 1},
	}
	data, err := json.Marshal(jokes)
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	err = memory.Set("jokes", data)
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	j, err := joke.GetRandomJokeFromCache(memory)
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	if j == "" {
		t.Error("j should not be empty")
	}
}

func TestCheckJokesCache_True(t *testing.T) {
	defer Flush()

	jokes := []schema.Joke{
		{ID: 1, Link: "link1", Creator: 1},
		{ID: 2, Link: "link2", Creator: 1},
		{ID: 3, Link: "link3", Creator: 1},
	}
	data, err := json.Marshal(jokes)
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	err = memory.Set("jokes", data)
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	j, err := joke.CheckJokesCache(memory)
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	if j == false {
		t.Error("j should not be false")
	}
}

func TestCheckJokesCache_False(t *testing.T) {
	defer Flush()

	j, err := joke.CheckJokesCache(memory)
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	if j == true {
		t.Error("j should not be true")
	}
}

func TestCheckTotalJokesCache_True(t *testing.T) {
	defer Flush()

	err := memory.Set("total", []byte("10"))
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	j, err := joke.CheckTotalJokesCache(memory)
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	if j == false {
		t.Error("j should not be false")
	}
}

func TestCheckTotalJokesCache_False(t *testing.T) {
	defer Flush()

	j, err := joke.CheckTotalJokesCache(memory)
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	if j == true {
		t.Error("j should not be true")
	}
}

func TestGetCachedJokeByID(t *testing.T) {
	defer Flush()

	jokes := []schema.Joke{
		{ID: 1, Link: "link1", Creator: 1},
		{ID: 2, Link: "link2", Creator: 1},
		{ID: 3, Link: "link3", Creator: 1},
	}
	data, err := json.Marshal(jokes)
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	err = memory.Set("jokes", data)
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	j, err := joke.GetCachedJokeByID(memory, 1)
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	if j != "link1" {
		t.Error("j should be link1, got:", j)
	}

	k, err := joke.GetCachedJokeByID(memory, 4)
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	if k != "" {
		t.Error("k should be empty, got:", k)
	}
}

func TestGetCachedTotalJokes(t *testing.T) {
	defer Flush()

	err := memory.Set("total", []byte("10"))
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	j, err := joke.GetCachedTotalJokes(memory)
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	if j != 10 {
		t.Error("j should be 10, got:", j)
	}
}

func TestCheckJokeExists(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()

	defer Flush()

	conn, err := db.Acquire(ctx)
	if err != nil {
		t.Error("an error was thrown:", err)
	}
	defer conn.Release()

	err = conn.BeginFunc(ctx, func(t pgx.Tx) error {
		_, err := t.Exec(
			ctx,
			`INSERT INTO "administrators"
				(id, key, token, last_used)
				VALUES
				($1, $2, $3, $4),
				($5, $6, $7, $8)`,
			administratorsData...,
		)
		if err != nil {
			return err
		}
		_, err = t.Exec(
			ctx,
			`INSERT INTO "jokesbapak2" 
				(id, link, creator)
				VALUES
				($1, $2, $3),
				($4, $5, $6),
				($7, $8, $9)`,
			jokesData...,
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	j, err := joke.CheckJokeExists(db, ctx, "1")
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	if j == false {
		t.Error("j should not be false")
	}

	k, err := joke.CheckJokeExists(db, ctx, "4")
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	if k == true {
		t.Error("k should not be true")
	}
}
