package utils

import "golang.org/x/crypto/bcrypt"

// GenerateHash creates a bcrypt hash from the provided password
// Returns the hashed password as a string or an error if hashing fails
func GenerateHash(pwd string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// CompareHash compares a password with its hash
// Returns true if the password matches the hash, false otherwise
func CompareHash(pwd, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)); err != nil {
		return false
	}
	return true
}
