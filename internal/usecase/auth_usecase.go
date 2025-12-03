package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redukasquad/be-reduka/configs"
	"github.com/redukasquad/be-reduka/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type authUsecase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
	config         *configs.Config
	oauthConfig    *oauth2.Config
}

func NewAuthUsecase(userRepo domain.UserRepository, timeout time.Duration, cfg *configs.Config) domain.AuthUsecase {
	return &authUsecase{
		userRepo:       userRepo,
		contextTimeout: timeout,
		config:         cfg,
		oauthConfig: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (u *authUsecase) Register(c context.Context, req *domain.RegisterRequest) error {
	_, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// Check if email exists
	if _, err := u.userRepo.FindByEmail(req.Email); err == nil {
		return errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Generate verification token
	tokenBytes := make([]byte, 16)
	rand.Read(tokenBytes)
	verificationToken := hex.EncodeToString(tokenBytes)

	user := &domain.User{
		Username:          req.Username,
		Email:             req.Email,
		Password:          string(hashedPassword),
		VerificationToken: verificationToken,
		IsVerified:        false,
		Role:              "Students", // Default role
		Kelas:             "Kelas 12", // Default kelas
	}

	if err := u.userRepo.CreateUser(user); err != nil {
		return err
	}

	// Mock send email
	fmt.Printf("Sending verification email to %s with token: %s\n", user.Email, user.VerificationToken)

	return nil
}

func (u *authUsecase) Login(c context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	_, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := u.generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (u *authUsecase) GoogleLogin(c context.Context) (string, error) {
	url := u.oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return url, nil
}

func (u *authUsecase) GoogleCallback(c context.Context, code string) (*domain.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	token, err := u.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := u.oauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var googleUser struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Picture  string `json:"picture"`
		Verified bool   `json:"verified_email"`
	}

	if err := json.Unmarshal(userData, &googleUser); err != nil {
		return nil, err
	}

	user, err := u.userRepo.FindByEmail(googleUser.Email)
	if err != nil {
		// Register new user
		user = &domain.User{
			Username:     googleUser.Name,
			Email:        googleUser.Email,
			Password:     "", // No password for OAuth users
			IsVerified:   googleUser.Verified,
			Role:         "Students",
			Kelas:        "Kelas 12",
			ProfileImage: json.RawMessage(fmt.Sprintf(`"%s"`, googleUser.Picture)),
		}
		if err := u.userRepo.CreateUser(user); err != nil {
			return nil, err
		}
	}

	jwtToken, err := u.generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token: jwtToken,
		User:  user,
	}, nil
}

func (u *authUsecase) VerifyEmail(c context.Context, token string) error {
	_, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.FindByVerificationToken(token)
	if err != nil {
		return errors.New("invalid verification token")
	}

	if user.IsVerified {
		return errors.New("user already verified")
	}

	user.IsVerified = true
	user.VerificationToken = "" // Clear token after verification

	return u.userRepo.UpdateUser(user)
}

func (u *authUsecase) generateJWT(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.UserID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.config.JWTSecret))
}