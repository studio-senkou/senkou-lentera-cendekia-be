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

	gomail "github.com/studio-senkou/lentera-cendekia-be/utils/mail"
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

func (ac *AuthController) LoginUser(c *fiber.Ctx) error {
	return ac.Login(c, false)
}

func (ac *AuthController) LoginAdmin(c *fiber.Ctx) error {
	return ac.Login(c, true)
}

func (ac *AuthController) Login(c *fiber.Ctx, isAdmin bool) error {
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
	if user != nil && user.Role != "admin" && isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "Access denied",
			"error":   "You do not have permission to access this resource",
		})
	}

	if err != nil || user == nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "fail",
			"message": "Login failed",
			"error":   "Invalid email or password",
		})
	} else if !user.CheckPassword(loginRequest.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Login failed",
			"error":   "Invalid email or password",
		})
	}

	if !user.IsEmailVerified() || !user.IsActive {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Email not verified",
			"error":   "Please verify your email before logging in",
		})
	}

	accessToken, err := ac.jwtManager.GenerateToken(auth.Payload{
		UserID: user.ID,
		Role:   user.Role,
	}, time.Now().Add(24*time.Hour))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot generate access token",
			"error":   err.Error(),
		})
	}

	refreshToken, err2 := ac.jwtManager.GenerateToken(auth.Payload{
		UserID: user.ID,
		Role:   user.Role,
	}, time.Now().Add(30*24*time.Hour))
	if err2 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot generate refresh token",
			"error":   err2.Error(),
		})
	}

	if err := ac.authRepo.UpdateOrCreate(&models.UserHasToken{
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
			"active_role":          user.Role,
			"access_token":         accessToken.Token,
			"access_token_expiry":  accessToken.ExpiresAt,
			"refresh_token":        refreshToken.Token,
			"refresh_token_expiry": refreshToken.ExpiresAt,
		},
	})
}

func (ac *AuthController) VerifyAccount(c *fiber.Ctx) error {
	verifyAccountRequest := new(requests.VerifyAccountRequest)
	if validationError, err := validator.ValidateRequest(c, verifyAccountRequest); err != nil {
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

	user, err := ac.userRepo.GetByEmail(verifyAccountRequest.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot verify account",
			"error":   "Your account is not registered as user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Account verification successful",
		"data": fiber.Map{
			"email": verifyAccountRequest.Email,
			"role":  user.Role,
		},
	})
}

func (ac *AuthController) VerifyOneTimeToken(c *fiber.Ctx) error {
	verifyTokenRequest := new(requests.VerifyTokenRequest)
	if validationError, err := validator.ValidateRequest(c, verifyTokenRequest); err != nil {
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

	oneTimeToken, err := auth.CheckOneTimeTokenStatus(verifyTokenRequest.Token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid token",
			"error":   err.Error(),
		})
	}

	if oneTimeToken.ExpiresAt.Before(time.Now()) || oneTimeToken.Used {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Token has expired",
			"error":   "The provided token is no longer valid or has already been used",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Token is valid",
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid refresh token",
			"error":   err.Error(),
		})
	}

	payloadMap, ok := claims["payload"].(map[string]any)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid token payload format",
		})
	}

	var parsedUserID int
	if id, ok := payloadMap["user_id"].(float64); ok {
		parsedUserID = int(id)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user ID in token",
			"error":   "Invalid user ID",
		})
	}

	userIDStr := fmt.Sprintf("%v", parsedUserID)
	userID, err := strconv.ParseInt(userIDStr, 10, 32)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid refresh token",
			"error":   "Invalid user ID in token",
		})
	}

	if ok, err := ac.authRepo.ValidateSessionExist(int(userID), refreshTokenRequest.Token); err != nil || !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid refresh token",
			"error":   "Session does not exist or has been invalidated",
		})
	}

	newPayload := auth.Payload{
		UserID: uint(userID),
		Role:   payloadMap["role"].(string),
	}

	newAccessToken, err := ac.jwtManager.GenerateToken(newPayload, time.Now().Add(24*time.Hour))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot generate new access token",
			"error":   err.Error(),
		})
	}

	newRefreshToken, err2 := ac.jwtManager.GenerateToken(newPayload, time.Now().Add(30*24*time.Hour))
	if err2 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Cannot generate new refresh token",
			"error":   err2.Error(),
		})
	}

	if err := ac.authRepo.UpdateOrCreate(&models.UserHasToken{
		UserID: uint(userID),
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
			"active_role":          newPayload.Role,
			"access_token":         newAccessToken.Token,
			"access_token_expiry":  newAccessToken.ExpiresAt,
			"refresh_token":        newRefreshToken.Token,
			"refresh_token_expiry": newRefreshToken.ExpiresAt,
		},
	})
}

func (ac *AuthController) Logout(c *fiber.Ctx) error {
	userIDStr := fmt.Sprintf("%v", c.Locals("userID"))
	userID, err := strconv.ParseInt(userIDStr, 10, 32)

	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unable to log out",
			"error":   "Invalid user ID",
		})
	}

	if err := ac.authRepo.InvalidateToken(int(userID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unable to log out",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully logged out",
	})
}

func (ac *AuthController) ResetPasswordRequest(c *fiber.Ctx) error {
	requestPasswordRequest := new(requests.ResetPasswordRequest)
	if validationError, err := validator.ValidateRequest(c, requestPasswordRequest); err != nil {
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

	user, err := ac.userRepo.GetByEmail(requestPasswordRequest.Email)
	if err != nil || user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Reset password request failed",
			"error":   "User not found",
		})
	}

	oneTimeToken, err := auth.GenerateOneTimeToken(user.ID, "password_reset", 15*time.Minute)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to generate reset password token",
			"error":   err.Error(),
		})
	}

	email, err := gomail.NewMailFromTemplate(
		requestPasswordRequest.Email,
		"Reset Password Request",
		"templates/emails/reset_password.html",
		fiber.Map{
			"Name":      user.Name,
			"ResetLink": fmt.Sprintf("%s/reset-password?token=%s", app.GetEnv("APP_FE_URL", "http://localhost:3000"), oneTimeToken.Token),
		},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to create reset password email",
			"error":   err.Error(),
		})
	}

	if err := email.Send(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to send reset password email",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Reset password request is valid",
	})
}
