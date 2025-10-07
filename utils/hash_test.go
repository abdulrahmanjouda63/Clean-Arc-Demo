package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateHash(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false, // bcrypt handles empty passwords
		},
		{
			name:     "long password",
			password: "this_is_a_very_long_password_that_should_work_fine_with_bcrypt",
			wantErr:  false,
		},
		{
			name:     "password with special characters",
			password: "p@ssw0rd!@#$%^&*()",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := GenerateHash(tt.password)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
				assert.NotEqual(t, tt.password, hash) // Hash should be different from original password
			}
		})
	}
}

func TestCompareHash(t *testing.T) {
	tests := []struct {
		name     string
		password string
		hash     string
		expected bool
	}{
		{
			name:     "correct password",
			password: "password123",
			hash:     "$2a$10$vnz04c9pQOhKP3lc7p4LLOZYHapMZBdodhQdv5TYw/4gL3.xpGv4m", // "password123"
			expected: true,
		},
		{
			name:     "incorrect password",
			password: "wrongpassword",
			hash:     "$2a$10$vnz04c9pQOhKP3lc7p4LLOZYHapMZBdodhQdv5TYw/4gL3.xpGv4m", // "password123"
			expected: false,
		},
		{
			name:     "empty password",
			password: "",
			hash:     "$2a$10$vnz04c9pQOhKP3lc7p4LLOZYHapMZBdodhQdv5TYw/4gL3.xpGv4m", // "password123"
			expected: false,
		},
		{
			name:     "invalid hash",
			password: "password123",
			hash:     "invalid_hash",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareHash(tt.password, tt.hash)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateHashAndCompareHash(t *testing.T) {
	// Test that generated hash can be compared correctly
	password := "testpassword123"
	
	hash, err := GenerateHash(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	
	// The generated hash should match the original password
	result := CompareHash(password, hash)
	assert.True(t, result)
	
	// Wrong password should not match
	wrongResult := CompareHash("wrongpassword", hash)
	assert.False(t, wrongResult)
}

func TestHashConsistency(t *testing.T) {
	// Test that the same password generates different hashes (due to salt)
	password := "consistentpassword"
	
	hash1, err1 := GenerateHash(password)
	assert.NoError(t, err1)
	
	hash2, err2 := GenerateHash(password)
	assert.NoError(t, err2)
	
	// Hashes should be different due to random salt
	assert.NotEqual(t, hash1, hash2)
	
	// But both should verify against the original password
	assert.True(t, CompareHash(password, hash1))
	assert.True(t, CompareHash(password, hash2))
}
