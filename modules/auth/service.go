package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redukasquad/be-reduka/database/entities"
	"github.com/redukasquad/be-reduka/modules/users"
	"github.com/redukasquad/be-reduka/packages/dto"
	"github.com/redukasquad/be-reduka/packages/utils"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(input dto.RegisterInput) (entities.User, error)
	Login(input dto.LoginInput) (string, error)
	VerifyEmail(code string) error
}

type authService struct {
	repo users.Repository
}

func NewService(repo users.Repository) Service {
	return &authService{repo: repo}
}

func (s *authService) Register(input dto.RegisterInput) (entities.User, error) {
	user := entities.User{
		Username:   input.Username,
		Email:      input.Email,
		Role:       input.Role,
		IsVerified: false,
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.Password = string(passwordHash)

	if _, err := s.repo.FindByEmail(input.Email); err == nil {
		return user, errors.New("email already registered")
	}
	code := utils.GenerateVerificationCode()
	user.VerificationCode = code

	err = s.repo.Create(user)
	if err != nil {
		return user, err
	}

	emailBody := "Your verification code is: " + code
	go utils.SendEmail(user.Email, "Email Verification", emailBody)

	return user, nil
}

func (s *authService) VerifyEmail(code string) error {
	user, err := s.repo.FindByVerificationCode(code)
	if err != nil {
		return errors.New("invalid verification code")
	}

	if user.IsVerified {
		return errors.New("email already verified")
	}

	user.IsVerified = true
	user.VerificationCode = ""

	return s.repo.Update(user)
}

func (s *authService) Login(input dto.LoginInput) (string, error) {
	email := input.Email
	password := input.Password

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	return generateToken(int(user.ID))
}

func generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
