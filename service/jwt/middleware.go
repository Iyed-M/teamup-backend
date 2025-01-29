package jwt_service

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (j JwtService) Middleware(c *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SuccessHandler: jwtSuccess,
		ErrorHandler:   jwtError,
		SigningKey:     jwtware.SigningKey{Key: j.Secret},
	})(c)
}

func jwtError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Unauthorized",
	})
}

func jwtSuccess(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userId, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		log.Errorw("JwtMiddleware error", "err", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	c.Locals("userId", userId)
	return c.Next()
}
