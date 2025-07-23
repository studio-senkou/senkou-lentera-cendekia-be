package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/studio-senkou/lentera-cendekia-be/utils/cache"
)

type OneTimeToken struct {
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Purpose   string    `json:"purpose"` // "password_reset", "email_verification", "account_activation"
	Used      bool      `json:"used"`
}

func GenerateOneTimeToken(userID int, purpose string, expiry time.Duration) (*OneTimeToken, error) {
	log.Printf("[TOKEN] Generating one-time token for userID: %d, purpose: %s, expiry: %v", userID, purpose, expiry)

	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		log.Printf("[TOKEN] ERROR: Failed to generate random bytes for userID %d: %v", userID, err)
		return nil, fmt.Errorf("failed to generate random token: %w", err)
	}

	token := hex.EncodeToString(bytes)
	log.Printf("[TOKEN] Generated token for userID %d: %s... (showing first 8 chars)", userID, token[:8])

	oneTimeToken := &OneTimeToken{
		UserID:    userID,
		Token:     token,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(expiry),
		Used:      false,
	}

	ctx := context.Background()
	key := fmt.Sprintf("one_time_token:%s", token)

	if err := cache.Set(ctx, key, oneTimeToken, expiry); err != nil {
		log.Printf("[TOKEN] ERROR: Failed to cache token for userID %d: %v", userID, err)
		return nil, fmt.Errorf("failed to cache one-time token: %w", err)
	}

	log.Printf("[TOKEN] SUCCESS: One-time token cached successfully for userID %d, expires at: %v", userID, oneTimeToken.ExpiresAt)
	return oneTimeToken, nil
}

func ValidateOneTimeToken(token, purpose string) (*OneTimeToken, error) {
	tokenPreview := token
	if len(token) > 8 {
		tokenPreview = token[:8] + "..."
	}
	log.Printf("[TOKEN] Validating one-time token: %s, purpose: %s", tokenPreview, purpose)

	ctx := context.Background()
	key := fmt.Sprintf("one_time_token:%s", token)

	var oneTimeToken OneTimeToken
	if err := cache.Get(ctx, key, &oneTimeToken); err != nil {
		log.Printf("[TOKEN] ERROR: Failed to retrieve token %s from cache: %v", tokenPreview, err)
		return nil, fmt.Errorf("failed to retrieve one-time token: %w", err)
	}

	log.Printf("[TOKEN] Token found - UserID: %d, Purpose: %s, Used: %t, Expires: %v",
		oneTimeToken.UserID, oneTimeToken.Purpose, oneTimeToken.Used, oneTimeToken.ExpiresAt)

	if oneTimeToken.Used {
		log.Printf("[TOKEN] ERROR: Token %s has already been used by userID %d", tokenPreview, oneTimeToken.UserID)
		return nil, fmt.Errorf("one-time token has already been used")
	}

	if oneTimeToken.Purpose != purpose {
		log.Printf("[TOKEN] ERROR: Purpose mismatch for token %s - expected: %s, got: %s",
			tokenPreview, purpose, oneTimeToken.Purpose)
		return nil, fmt.Errorf("one-time token purpose mismatch: expected %s, got %s", purpose, oneTimeToken.Purpose)
	}

	if time.Now().After(oneTimeToken.ExpiresAt) {
		log.Printf("[TOKEN] ERROR: Token %s has expired at %v for userID %d",
			tokenPreview, oneTimeToken.ExpiresAt, oneTimeToken.UserID)

		// Delete expired token
		if delErr := cache.Delete(ctx, key); delErr != nil {
			log.Printf("[TOKEN] WARNING: Failed to delete expired token %s: %v", tokenPreview, delErr)
		} else {
			log.Printf("[TOKEN] Expired token %s deleted from cache", tokenPreview)
		}

		return nil, fmt.Errorf("one-time token has expired")
	}

	// Mark as used
	log.Printf("[TOKEN] Marking token %s as used for userID %d", tokenPreview, oneTimeToken.UserID)
	oneTimeToken.Used = true

	if err := cache.Set(ctx, key, oneTimeToken, time.Until(oneTimeToken.ExpiresAt)); err != nil {
		log.Printf("[TOKEN] ERROR: Failed to update token usage status for %s: %v", tokenPreview, err)
		return nil, fmt.Errorf("failed to update one-time token usage: %w", err)
	}

	log.Printf("[TOKEN] SUCCESS: Token %s validated and consumed for userID %d", tokenPreview, oneTimeToken.UserID)
	return &oneTimeToken, nil
}

func InvalidateOneTimeToken(token string) error {
	tokenPreview := token
	if len(token) > 8 {
		tokenPreview = token[:8] + "..."
	}
	log.Printf("[TOKEN] Invalidating token: %s", tokenPreview)

	ctx := context.Background()
	key := fmt.Sprintf("one_time_token:%s", token)

	if err := cache.Delete(ctx, key); err != nil {
		log.Printf("[TOKEN] ERROR: Failed to invalidate token %s: %v", tokenPreview, err)
		return err
	}

	log.Printf("[TOKEN] SUCCESS: Token %s invalidated successfully", tokenPreview)
	return nil
}

func CheckOneTimeTokenStatus(token string) (*OneTimeToken, error) {
	tokenPreview := token
	if len(token) > 8 {
		tokenPreview = token[:8] + "..."
	}
	log.Printf("[TOKEN] Checking status for token: %s", tokenPreview)

	ctx := context.Background()
	key := fmt.Sprintf("one_time_token:%s", token)

	var oneTimeToken OneTimeToken
	if err := cache.Get(ctx, key, &oneTimeToken); err != nil {
		log.Printf("[TOKEN] ERROR: Token %s not found in cache: %v", tokenPreview, err)
		return nil, fmt.Errorf("token not found")
	}

	log.Printf("[TOKEN] Token status - UserID: %d, Purpose: %s, Used: %t, Expires: %v, Valid: %t",
		oneTimeToken.UserID, oneTimeToken.Purpose, oneTimeToken.Used,
		oneTimeToken.ExpiresAt, time.Now().Before(oneTimeToken.ExpiresAt))

	return &oneTimeToken, nil
}
