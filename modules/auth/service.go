package auth

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redukasquad/be-reduka/database/entities"
	"github.com/redukasquad/be-reduka/modules/users"
	"github.com/redukasquad/be-reduka/packages/utils"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/oauth2/v2"
	"gorm.io/gorm"
)

type Service interface {
	Register(input RegisterInput) (entities.User, error)
	Login(input LoginInput) (string, error)
	VerifyEmail(email, code string) error
	ResendVerificationCode(email string) error
	Me(user_id int) (*entities.User, error)
	LoginOrRegisterWithGoogle(googleUserInfo *oauth2.Userinfo) (*entities.User, string, error)
	ForgotPassword(input ForgotPasswordInput) error
	ResetPassword(input ResetPasswordInput) error
}

type authService struct {
	repo users.Repository
}

func NewService(repo users.Repository) Service {
	return &authService{repo: repo}
}

func (s *authService) Register(input RegisterInput) (entities.User, error) {
	defaultRole := "STUDENT"
	user := entities.User{
		Username:     input.Username,
		Email:        input.Email,
		AuthProvider: "PASSWORD",
		Role:         &defaultRole,
		IsVerified:   false,
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}
	user.Password = string(passwordHash)

	if _, err := s.repo.FindByEmail(input.Email); err == nil {
		return user, errors.New("email already registered")
	}
	code := utils.GenerateVerificationCode()
	hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}
	hashedCodeStr := string(hashedCode)
	user.VerificationCode = &hashedCodeStr

	err = s.repo.Create(&user)
	if err != nil {
		return user, err
	}

	emailBody := "Your verification code is: " + code
	go func() {
		err := utils.SendEmail(user.Email, "Email Verification", emailBody)
		if err != nil {
			log.Println("[EMAIL] FAILED:", err)
		}
	}()

	return user, nil
}

func (s *authService) VerifyEmail(email, code string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return errors.New("email not found")
	}

	if user.IsVerified {
		return errors.New("email already verified")
	}

	if user.VerificationCode == nil {
		return errors.New("invalid verification code")
	}
	err = bcrypt.CompareHashAndPassword([]byte(*user.VerificationCode), []byte(code))
	if err != nil {
		return errors.New("invalid verification code")
	}

	user.IsVerified = true
	user.VerificationCode = nil

	return s.repo.Update(user)
}

func (s *authService) ResendVerificationCode(email string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return errors.New("email not found")
	}

	if user.IsVerified {
		return errors.New("email already verified")
	}

	code := utils.GenerateVerificationCode()
	hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashedCodeStr := string(hashedCode)
	user.VerificationCode = &hashedCodeStr

	if err := s.repo.Update(user); err != nil {
		return err
	}

	emailBody := `
		<h2>Email Verification</h2>
		<p>Your verification code is: <strong>` + code + `</strong></p>
		<p>This code will expire in 15 minutes.</p>
	`
	go utils.SendEmail(user.Email, "Email Verification - Reduka", emailBody)

	return nil
}

func (s *authService) Login(input LoginInput) (string, error) {
	email := input.Email
	password := input.Password

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if user.AuthProvider == "GOOGLE" {
		return "", errors.New("this account uses Google login")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !user.IsVerified {
		return "", errors.New("email not verified")
	}

	return generateToken(int(user.ID))
}

func generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 1 ).Unix(), // 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *authService) ForgotPassword(input ForgotPasswordInput) error {
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

func (s *authService) ResetPassword(input ResetPasswordInput) error {
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

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
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
		defaultRole := "STUDENT"
		newUser := &entities.User{
			Username:     googleUserInfo.Name,
			Email:        googleUserInfo.Email,
			Password:     "",
			AuthProvider: "GOOGLE",
			Role:         &defaultRole,
			IsVerified:   true,
		}
		if err := s.repo.Create(newUser); err != nil {
			return nil, "", err
		}
		user = *newUser
	} else if err != nil {
		return nil, "", err
	} else {
		if user.AuthProvider == "PASSWORD" {
			return nil, "", errors.New("this account uses password login")
		}
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
