package auth

import (
	"encoding/base64"

	"github.com/hctilf/go-test-task-medods/internal/app"
	"github.com/hctilf/go-test-task-medods/pkg/jwt_tools"

	"github.com/gofiber/fiber/v2"
)

type (
	AuthHandler struct {
		app *app.Application
	}

	requestPayload struct {
		RefreshToken string `json:"refreshToken"`
	}

	responsePayload struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
)

func NewAuthHandler(app *app.Application) *AuthHandler {
	return &AuthHandler{
		app: app,
	}
}

func (h *AuthHandler) GetAccessNRefreshToken(c *fiber.Ctx) error {
	userGUID := c.Query("userGUID")
	if userGUID == "" {
		h.app.Logger.Errorf("Error getting userGUID from query: empty guid")

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
		})
	}

	userIpAddress := c.IP()

	refreshToken, err := h.app.Tokens.GetTokenByUserGUID(userIpAddress, userGUID)
	if err != nil {
		h.app.Logger.Errorf("Error getting token by userGUID: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
		})
	}

	accessToken, err := jwt_tools.GenerateToken(userGUID, userIpAddress)
	if err != nil {
		h.app.Logger.Errorf("Error generating access token: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
		})
	}

	resfreshBase64 := base64.StdEncoding.EncodeToString([]byte(refreshToken.Token))

	responsePayload := responsePayload{
		AccessToken:  accessToken,
		RefreshToken: resfreshBase64,
	}

	return c.JSON(responsePayload)
}

func (h *AuthHandler) PostRefreshToken(c *fiber.Ctx) error {
	var payload requestPayload

	reqIp := c.IP()

	if err := c.BodyParser(&payload); err != nil {
		h.app.Logger.Errorf("Error parsing request body: %v", err)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
		})
	}

	if payload.RefreshToken == "" {
		h.app.Logger.Errorf("Error refreshing token: empty refresh token")

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
		})
	}

	bytePayload, err := base64.StdEncoding.DecodeString(payload.RefreshToken)
	if err != nil {
		h.app.Logger.Errorf("Error decoding refresh token: %v", err)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
		})
	}

	payload.RefreshToken = string(bytePayload)

	isIpChanged, err := h.app.Tokens.RefreshToken(payload.RefreshToken, reqIp)
	if err != nil {
		h.app.Logger.Errorf("Error refreshing token: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
		})
	}

	if isIpChanged {
		h.app.Background(func() {
			h.app.Logger.Info("Sending notification to user about ip change")
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
	})
}

func (h *AuthHandler) TestTokenCreate(c *fiber.Ctx) error {
	tkn, err := h.app.Tokens.CreateTestToken(c.IP())

	if err != nil {
		h.app.Logger.Errorf("Error creating test token: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"guid":   tkn.UserGUID.String(),
	})
}
