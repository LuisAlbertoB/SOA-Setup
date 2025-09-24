package services

import "golang.org/x/crypto/bcrypt"

// HashPassword genera un hash bcrypt de una contraseña.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compara una contraseña en texto plano con su hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// Si err es nil, la contraseña es correcta.
	return err == nil
}
