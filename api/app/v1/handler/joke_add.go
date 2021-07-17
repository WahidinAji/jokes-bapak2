package handler

import (
	"context"

	"jokes-bapak2-api/app/v1/core"
	"jokes-bapak2-api/app/v1/models"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
)

func AddNewJoke(c *fiber.Ctx) error {
	var body models.Joke
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	sql, args, err := psql.Insert("jokesbapak2").Columns("link", "creator").Values(body.Link, c.Locals("userID")).ToSql()
	if err != nil {
		return err
	}

	_, err = db.Query(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	jokes, err := core.GetAllJSONJokes(db)
	if err != nil {
		return err
	}
	memory.Set("jokes", jokes, cache.NoExpiration)

	return c.Status(fiber.StatusCreated).JSON(models.ResponseJoke{
		Link: body.Link,
	})
}
