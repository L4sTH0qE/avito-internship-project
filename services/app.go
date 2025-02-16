package services

import (
	"awesomeProject/dao"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/jinzhu/gorm"
	_ "github.com/joho/godotenv"
	_ "golang.org/x/crypto/bcrypt"
	"os"
	_ "strings"
	"time"
)

type Token struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateJWT - Генерация jwt-токена по username.
func GenerateJWT(username string) (tokenString string, err error) {

	tokenPassword := []byte(os.Getenv("TOKEN_PASSWORD"))

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Token{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(tokenPassword)
	return
}

// GetUser - Получение пользователя из БД по username.
func GetUser(username string) (user dao.User, err error) {
	user = dao.User{}
	err = nil
	result := GetDB().Where("username = ?", username).First(&user)
	if result.Error != nil {
		err = fmt.Errorf("no user with username %s", username)
	}
	return
}

// GetMerch - Получение мерча из БД по type.
func GetMerch(merchType string) (merch dao.Merch, err error) {
	merch = dao.Merch{}
	err = nil
	result := GetDB().Where("type = ?", merchType).First(&merch)
	if result.Error != nil {
		err = fmt.Errorf("no merch with type %s", merchType)
	}
	return
}
