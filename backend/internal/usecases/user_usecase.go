package usecases

import (
	"fmt"
	"myapp/internal/config"
	"myapp/internal/entities"
	"myapp/internal/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordIncorrect = fmt.Errorf("password is incorrect")
var ErrUserNotFound = repositories.ErrUserNotFound

type PostSigninUsecase struct {
	repository entities.UserRepository
}

func NewPostSigninUsecase(r *repositories.UserRepository) *PostSigninUsecase {
	return &PostSigninUsecase{
		repository: r,
	}
}

func (u *PostSigninUsecase) Execute(username, password string) (*entities.User, string, error) {

	// Call repository
	user, err := u.repository.GetUserByName(username)
	if err != nil {
		return nil, "", err
	}
	passwordByte := []byte(password)
	storedPasswordByte := []byte(user.Password)

	err = bcrypt.CompareHashAndPassword(storedPasswordByte, passwordByte)

	if err != nil {
		return nil, "", ErrPasswordIncorrect
	}

	// トークンの発行（ヘッダー・ペイロード）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Name": user.Name,
		"Id":   user.Id,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"iat":  time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return nil, "", err
	}

	return user, tokenString, nil
}

type PostSignoutUsecase struct {
	repository entities.UserRepository
}

func NewPostSignoutUsecase(r *repositories.UserRepository) *PostSignoutUsecase {
	return &PostSignoutUsecase{
		repository: r,
	}
}

func (u *PostSignoutUsecase) Execute() (string, error) {

	// トークンの発行（ヘッダー・ペイロード）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return tokenString, err
	}

	return tokenString, nil
}

type GetUserUsecase struct {
	repository entities.UserRepository
}

func NewGetUserUsecase(r *repositories.UserRepository) *GetUserUsecase {
	return &GetUserUsecase{
		repository: r,
	}
}

func (u *GetUserUsecase) Execute(userId int) (*entities.User, error) {

	// Call repository
	user, err := u.repository.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}
