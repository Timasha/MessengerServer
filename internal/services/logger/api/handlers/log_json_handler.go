package handlers

import (
	"messengerServer/internal/api_objects/logger/responses"

	"github.com/gofiber/fiber/v2"
)

func LogsFileHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var resp responses.LogFileResponse
		if string(c.Request().Header.ContentType()) != "plain/text" {
			resp.Err = "wrong content type"
			return nil
		}
		return nil
	}
}
