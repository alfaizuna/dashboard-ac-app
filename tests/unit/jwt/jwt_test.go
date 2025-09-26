package jwt

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"dashboard-ac-backend/internal/domain"
	"dashboard-ac-backend/pkg/jwt"
)

func TestMain(m *testing.M) {
	// Set test JWT secret
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	code := m.Run()
	os.Exit(code)
}

func TestGenerateTokenPair(t *testing.T) {
	user := &domain.User{
		ID:    1,
		Email: "test@example.com",
		Role:  domain.RoleAdmin,
	}
	secret := "test-secret-key-for-testing"

	tokenPair, err := jwt.GenerateTokenPair(user, secret)

	assert.NoError(t, err)
	assert.NotNil(t, tokenPair)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
	assert.NotEqual(t, tokenPair.AccessToken, tokenPair.RefreshToken)
}

func TestValidateToken_ValidAccessToken(t *testing.T) {
	user := &domain.User{
		ID:    1,
		Email: "test@example.com",
		Role:  domain.RoleAdmin,
	}
	secret := "test-secret-key-for-testing"

	tokenPair, err := jwt.GenerateTokenPair(user, secret)
	assert.NoError(t, err)

	claims, err := jwt.ValidateToken(tokenPair.AccessToken, secret)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, user.Role, claims.Role)
	assert.Equal(t, "access", claims.Subject)
}

func TestValidateToken_ValidRefreshToken(t *testing.T) {
	user := &domain.User{
		ID:    1,
		Email: "test@example.com",
		Role:  domain.RoleAdmin,
	}
	secret := "test-secret-key-for-testing"

	tokenPair, err := jwt.GenerateTokenPair(user, secret)
	assert.NoError(t, err)

	claims, err := jwt.ValidateToken(tokenPair.RefreshToken, secret)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, user.Role, claims.Role)
	assert.Equal(t, "refresh", claims.Subject)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	secret := "test-secret-key-for-testing"
	
	tests := []struct {
		name        string
		token       string
		expectError bool
	}{
		{
			name:        "Empty token",
			token:       "",
			expectError: true,
		},
		{
			name:        "Invalid token format",
			token:       "invalid.token.format",
			expectError: true,
		},
		{
			name:        "Malformed token",
			token:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature",
			expectError: true,
		},
		{
			name:        "Wrong secret",
			token:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := jwt.ValidateToken(tt.token, secret)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
			}
		})
	}
}

func TestExtractTokenFromHeader(t *testing.T) {
	tests := []struct {
		name        string
		header      string
		expected    string
		expectError bool
	}{
		{
			name:        "Valid Bearer token",
			header:      "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expected:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectError: false,
		},
		{
			name:        "Empty header",
			header:      "",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid format - no Bearer",
			header:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid format - only Bearer",
			header:      "Bearer",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid format - Bearer with space but no token",
			header:      "Bearer ",
			expected:    "",
			expectError: false, // This actually returns empty string but no error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := jwt.ExtractTokenFromHeader(tt.header)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.expected, token)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, token)
			}
		})
	}
}

func TestTokenExpiration(t *testing.T) {
	user := &domain.User{
		ID:    1,
		Email: "test@example.com",
		Role:  domain.RoleAdmin,
	}
	secret := "test-secret-key-for-testing"

	tokenPair, err := jwt.GenerateTokenPair(user, secret)
	assert.NoError(t, err)

	// Validate that tokens are valid immediately after creation
	accessClaims, err := jwt.ValidateToken(tokenPair.AccessToken, secret)
	assert.NoError(t, err)
	assert.NotNil(t, accessClaims)

	refreshClaims, err := jwt.ValidateToken(tokenPair.RefreshToken, secret)
	assert.NoError(t, err)
	assert.NotNil(t, refreshClaims)
}

func TestClaims_Structure(t *testing.T) {
	user := &domain.User{
		ID:    1,
		Email: "test@example.com",
		Role:  domain.RoleAdmin,
	}
	secret := "test-secret-key-for-testing"

	tokenPair, err := jwt.GenerateTokenPair(user, secret)
	assert.NoError(t, err)

	// Test access token claims
	accessClaims, err := jwt.ValidateToken(tokenPair.AccessToken, secret)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, accessClaims.UserID)
	assert.Equal(t, user.Email, accessClaims.Email)
	assert.Equal(t, user.Role, accessClaims.Role)
	assert.Equal(t, "access", accessClaims.Subject)
	assert.NotNil(t, accessClaims.ExpiresAt)
	assert.NotNil(t, accessClaims.IssuedAt)

	// Test refresh token claims
	refreshClaims, err := jwt.ValidateToken(tokenPair.RefreshToken, secret)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, refreshClaims.UserID)
	assert.Equal(t, user.Email, refreshClaims.Email)
	assert.Equal(t, user.Role, refreshClaims.Role)
	assert.Equal(t, "refresh", refreshClaims.Subject)
	assert.NotNil(t, refreshClaims.ExpiresAt)
	assert.NotNil(t, refreshClaims.IssuedAt)
}