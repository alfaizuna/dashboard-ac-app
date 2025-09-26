package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dashboard-ac-backend/pkg/hash"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  false, // bcrypt can hash empty strings
		},
		{
			name:     "Long password",
			password: "this-is-a-very-long-password-that-should-still-work-fine-with-bcrypt",
			wantErr:  false,
		},
		{
			name:     "Special characters",
			password: "p@ssw0rd!@#$%^&*()",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := hash.HashPassword(tt.password)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hashedPassword)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hashedPassword)
				assert.NotEqual(t, tt.password, hashedPassword)
				assert.True(t, len(hashedPassword) > 0)
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"
	hashedPassword, err := hash.HashPassword(password)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		password       string
		hashedPassword string
		expectError    bool
	}{
		{
			name:           "Correct password",
			password:       password,
			hashedPassword: hashedPassword,
			expectError:    false,
		},
		{
			name:           "Incorrect password",
			password:       "wrongpassword",
			hashedPassword: hashedPassword,
			expectError:    true,
		},
		{
			name:           "Empty password",
			password:       "",
			hashedPassword: hashedPassword,
			expectError:    true,
		},
		{
			name:           "Empty hash",
			password:       password,
			hashedPassword: "",
			expectError:    true,
		},
		{
			name:           "Invalid hash format",
			password:       password,
			hashedPassword: "invalid-hash",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := hash.CheckPassword(tt.hashedPassword, tt.password)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHashPassword_Consistency(t *testing.T) {
	password := "testpassword"
	
	// Hash the same password multiple times
	hash1, err1 := hash.HashPassword(password)
	hash2, err2 := hash.HashPassword(password)
	
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEmpty(t, hash1)
	assert.NotEmpty(t, hash2)
	
	// Hashes should be different (due to salt)
	assert.NotEqual(t, hash1, hash2)
	
	// But both should verify correctly
	assert.NoError(t, hash.CheckPassword(hash1, password))
	assert.NoError(t, hash.CheckPassword(hash2, password))
}

func TestCheckPassword_CaseSensitive(t *testing.T) {
	password := "TestPassword"
	hashedPassword, err := hash.HashPassword(password)
	assert.NoError(t, err)
	
	// Password checking should be case sensitive
	assert.NoError(t, hash.CheckPassword(hashedPassword, "TestPassword"))
	assert.Error(t, hash.CheckPassword(hashedPassword, "testpassword"))
	assert.Error(t, hash.CheckPassword(hashedPassword, "TESTPASSWORD"))
}