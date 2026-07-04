package caedral_test

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

const apiKeyPrefix = "cd_live_"

type testKeyFixture struct {
	UserID   string
	APIKeyID string
	RawKey   string
	conn     *pgx.Conn
}

func (f *testKeyFixture) Cleanup() {
	if f.conn == nil {
		return
	}
	ctx := context.Background()
	_, _ = f.conn.Exec(ctx, `DELETE FROM usage_logs WHERE user_id = $1`, f.UserID)
	_, _ = f.conn.Exec(ctx, `DELETE FROM api_keys WHERE id = $1`, f.APIKeyID)
	_, _ = f.conn.Exec(ctx, `DELETE FROM subscriptions WHERE user_id = $1`, f.UserID)
	_, _ = f.conn.Exec(ctx, `DELETE FROM "user" WHERE id = $1`, f.UserID)
	f.conn.Close(context.Background())
}

func generateAPIKeySecret() (string, error) {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return apiKeyPrefix + base64.RawURLEncoding.EncodeToString(buf), nil
}

func createTestAPIKey(t *testing.T) (*testKeyFixture, error) {
	t.Helper()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}

	userID := randomUUID()
	apiKeyID := randomUUID()
	subID := randomUUID()
	rawKey, err := generateAPIKeySecret()
	if err != nil {
		conn.Close(context.Background())
		return nil, err
	}
	keyPrefix := rawKey[:16]
	hash, err := bcrypt.GenerateFromPassword([]byte(rawKey), bcrypt.DefaultCost)
	if err != nil {
		conn.Close(context.Background())
		return nil, err
	}
	email := fmt.Sprintf("sdk-go-test-%s@example.com", userID)

	ctx := context.Background()
	_, err = conn.Exec(ctx, `
		INSERT INTO "user" (id, name, email, email_verified, balance_cents, account_status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, userID, "SDK Go Test", email, true, 0, "active")
	if err != nil {
		conn.Close(context.Background())
		return nil, err
	}

	_, err = conn.Exec(ctx, `
		INSERT INTO subscriptions (
			id, user_id, plan, status, weekly_pool_limit, weekly_pool_used,
			overage_enabled, overage_used_cents
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, subID, userID, "pro", "active", 1_000_000, 0, false, 0)
	if err != nil {
		conn.Close(context.Background())
		return nil, err
	}

	_, err = conn.Exec(ctx, `
		INSERT INTO api_keys (id, user_id, name, key_prefix, key_hash)
		VALUES ($1, $2, $3, $4, $5)
	`, apiKeyID, userID, "SDK Go test key", keyPrefix, string(hash))
	if err != nil {
		conn.Close(context.Background())
		return nil, err
	}

	return &testKeyFixture{
		UserID:   userID,
		APIKeyID: apiKeyID,
		RawKey:   rawKey,
		conn:     conn,
	}, nil
}

func randomUUID() string {
	buf := make([]byte, 16)
	_, _ = rand.Read(buf)
	buf[6] = (buf[6] & 0x0f) | 0x40
	buf[8] = (buf[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:16])
}
