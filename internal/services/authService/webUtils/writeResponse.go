package webUtils

import (
	"encoding/json"
	"messengerServer/internal/servicesApi/authService/responses"

	"github.com/gofiber/fiber/v2"
)

// Функция для отправки ответа на запрос в зависимости от его типа.
func WriteResponse[T responses.AuthorizationResponse | responses.CheckTokenResponse | responses.RegistrationResponse](resp T, respCode int, c *fiber.Ctx) {
	rawResp, _ := json.Marshal(resp)
	c.Status(respCode).Write(rawResp)
}
