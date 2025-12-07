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
	"google.golang.org/api/oauth2/v2"
	"gorm.io/gorm"
)

type Service interface {
	Register(input dto.RegisterInput) (entities.User, error)
	Login(input dto.LoginInput) (string, error)
	VerifyEmail(code string) error
	Me(user_id int) (*entities.User, error)
	LoginOrRegisterWithGoogle(googleUserInfo *oauth2.Userinfo) (*entities.User, string, error)
	ForgotPassword(input dto.ForgotPasswordInput) error
	ResetPassword(input dto.ResetPasswordInput) error
}

type authService struct {
	repo users.Repository
}

func NewService(repo users.Repository) Service {
	return &authService{repo: repo}
}

func (s *authService) Register(input dto.RegisterInput) (entities.User, error) {
	user := entities.User{
		Username:     input.Username,
		Email:        input.Email,
		Role:         input.Role,
		NoTelp:       input.NoTelp,
		JenisKelamin: input.JenisKelamin,
		Kelas:        input.Kelas,
		IsVerified:   false,
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

	err = s.repo.Create(&user)
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

func (s *authService) ForgotPassword(input dto.ForgotPasswordInput) error {
	user, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return errors.New("email not found")
	}

	token := utils.GenerateVerificationCode()
	user.ResetPasswordToken = token
	expiry := time.Now().Add(15 * time.Minute)
	user.ResetPasswordTokenExpiry = &expiry

	if err := s.repo.Update(user); err != nil {
		return err
	}

	emailBody := "Your reset password code is: " + token
	go utils.SendEmail(user.Email, "Reset Password", emailBody)

	return nil
}

func (s *authService) ResetPassword(input dto.ResetPasswordInput) error {
	if input.NewPassword != input.ConfirmPassword {
		return errors.New("passwords do not match")
	}

	user, err := s.repo.FindByResetToken(input.Token)
	if err != nil {
		return errors.New("invalid token")
	}

	if user.ResetPasswordTokenExpiry == nil || user.ResetPasswordTokenExpiry.Before(time.Now()) {
		return errors.New("token expired")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.MinCost)
	if err != nil {
		return err
	}

	user.Password = string(passwordHash)
	user.ResetPasswordToken = ""
	user.ResetPasswordTokenExpiry = nil

	return s.repo.Update(user)
}

func (s *authService) LoginOrRegisterWithGoogle(googleUserInfo *oauth2.Userinfo) (*entities.User, string, error) {
	user, err := s.repo.FindByEmail(googleUserInfo.Email)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		newUser := &entities.User{
			Username:   googleUserInfo.Name,
			Email:      googleUserInfo.Email,
			Password:   "",
			IsVerified: true,
			Kelas:      "Kelas 12",
			Role:       "Students",
		}
		if err := s.repo.Create(newUser); err != nil {
			return nil, "", err
		}
		user = *newUser
	} else if err != nil {
		return nil, "", err
	}

	token, err := generateToken(int(user.ID))
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func (s *authService) Me(user_id int) (*entities.User, error) {
	user, err := s.repo.FindByID(user_id)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return &user, nil
}