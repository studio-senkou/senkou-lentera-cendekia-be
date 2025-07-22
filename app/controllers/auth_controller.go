package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/models"
	"github.com/studio-senkou/lentera-cendekia-be/app/requests"
	"github.com/studio-senkou/lentera-cendekia-be/database"
	"github.com/studio-senkou/lentera-cendekia-be/utils/app"
	"github.com/studio-senkou/lentera-cendekia-be/utils/auth"
	"github.com/studio-senkou/lentera-cendekia-be/utils/validator"
)

type AuthController struct {
	jwtManager *auth.JwtManager
	userRepo   *models.UserRepository
	authRepo   *models.AuthenticationRepository
}

func NewAuthController() *AuthController {
	db := database.GetDB()
	authSecret := app.GetEnv("AUTH_SECRET", "")

	return &AuthController{
		jwtManager: auth.NewJwtManager(authSecret),
		userRepo:   models.NewUserRepository(db),
		authRepo:   models.NewAuthenticationRepository(db),
	}
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	loginRequest := new(requests.LoginRequest)
	if validationError, err := validator.ValidateRequest(c, loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if len(validationError) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Bad request",
			"errors":  validationError,
		})
	}

	user, err := ac.userRepo.GetByEmail(loginRequest.Email)

	if err != nil || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Login failed",
			"error":   "Invalid email or password",
		})
	} else if !user.CheckPassword(loginRequest.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Login failed",
			"error":   "Invalid email or password",
		})
	}

	accessToken, err := ac.jwtManager.GenerateToken(user.ID, time.Now().Add(24*time.Hour))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot generate access token",
			"error":   err.Error(),
		})
	}

	refreshToken, err2 := ac.jwtManager.GenerateToken(user.ID, time.Now().Add(30*24*time.Hour))
	if err2 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot generate refresh token",
			"error":   err2.Error(),
		})
	}

	if err := ac.authRepo.Create(&models.UserHasToken{
		UserID: user.ID,
		Token:  refreshToken.Token,
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to save authentication token",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Login successful",
		"data": fiber.Map{
			"access_token":         accessToken.Token,
			"access_token_expiry":  accessToken.ExpiresAt,
			"refresh_token":        refreshToken.Token,
			"refresh_token_expiry": refreshToken.ExpiresAt,
		},
	})
}

func (ac *AuthController) RefreshToken(c *fiber.Ctx) error {
	refreshTokenRequest := new(requests.RefreshTokenRequest)
	if validationError, err := validator.ValidateRequest(c, refreshTokenRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot parse request body",
			"error":   err.Error(),
		})
	} else if len(validationError) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Bad request",
			"errors":  validationError,
		})
	}

	claims, err := ac.jwtManager.ValidateToken(refreshTokenRequest.Token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid refresh token",
			"error":   err.Error(),
		})
	}

	userIDStr := fmt.Sprintf("%v", claims["payload"])
	userID, err := strconv.ParseInt(userIDStr, 10, 32)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid refresh token",
			"error":   "Invalid user ID in token",
		})
	}

	if ok, err := ac.authRepo.ValidateSessionExist(int(userID), refreshTokenRequest.Token); err != nil || !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid refresh token",
			"error":   "Session does not exist or has been invalidated",
		})
	}

	newAccessToken, err := ac.jwtManager.GenerateToken(int(userID), time.Now().Add(24*time.Hour))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot generate new access token",
			"error":   err.Error(),
		})
	}

	newRefreshToken, err2 := ac.jwtManager.GenerateToken(int(userID), time.Now().Add(30*24*time.Hour))
	if err2 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot generate new refresh token",
			"error":   err2.Error(),
		})
	}

	if err := ac.authRepo.UpdateOrCreate(&models.UserHasToken{
		UserID: int(userID),
		Token:  newRefreshToken.Token,
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to save new authentication token",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully refreshed token",
		"data": fiber.Map{
			"access_token":         newAccessToken.Token,
			"access_token_expiry":  newAccessToken.ExpiresAt,
			"refresh_token":        newRefreshToken.Token,
			"refresh_token_expiry": newRefreshToken.ExpiresAt,
		},
	})
}
