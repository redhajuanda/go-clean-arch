package auth

import (
	"context"
	"go-clean-arch/configs"
	"go-clean-arch/internal/domain"
	"go-clean-arch/internal/ierr"
	"go-clean-arch/internal/repository/port"
	"go-clean-arch/pkg/otel"
	"go-clean-arch/shared/auth"
	"go-clean-arch/shared/password"
	"go-clean-arch/shared/times"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// Service encapsulates the authentication logic.
type Service struct {
	cfg         *configs.Config
	repoRegitry port.RepositoryRegistry
}

// NewService creates and returns a new auth service
func NewService(cfg *configs.Config, repoRegitry port.RepositoryRegistry) Service {
	return Service{cfg, repoRegitry}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s Service) Login(ctx context.Context, req RequestLogin) (LoginResponse, error) {

	ctx, span := otel.Start(ctx)
	defer span.End()

	var res LoginResponse

	err := req.Validate()
	if err != nil {
		return res, err
	}

	identity, err := s.authenticate(ctx, req.Username, req.Password)
	if err != nil {
		return res, err
	}

	accessToken, expiresAt, refreshToken, err := s.generateJWT(ctx, identity)
	return LoginResponse{
		AccessToken:  accessToken,
		ExpiresAt:    expiresAt.Format(time.RFC3339),
		RefreshToken: refreshToken,
	}, err

}

// RefreshToken refresh the access token
func (s Service) RefreshToken(ctx context.Context, req RefreshTokenRequest) (LoginResponse, error) {

	ctx, span := otel.Start(ctx)
	defer span.End()

	var res LoginResponse

	err := req.Validate()
	if err != nil {
		return res, err
	}

	token, err := auth.VerifyToken(req.RefreshToken, s.cfg.JWT.SigningKey)
	if err != nil {
		return res, ierr.ErrInvalidToken
	}
	claims := token.Claims.(jwt.MapClaims)
	var tokenType string
	if val, ok := claims["token_type"].(string); ok {
		tokenType = val
	}

	if tokenType != "refresh" {
		return res, ierr.ErrInvalidToken
	}

	var id string
	if val, ok := claims["id"].(string); ok {
		id = val
	}

	repoUser := s.repoRegitry.GetUserRepository()
	user, err := repoUser.GetByID(ctx, id)
	if err != nil {
		return res, err
	}

	if *user.RefreshToken != req.RefreshToken {
		return res, ierr.ErrExpiredToken
	}

	accessToken, expiresAt, refreshToken, err := s.generateJWT(ctx, user)
	return LoginResponse{
		AccessToken:  accessToken,
		ExpiresAt:    expiresAt.Format(time.RFC3339),
		RefreshToken: refreshToken,
	}, err
}

// authenticate authenticates a user using username and password.
// if username and password are correct, an identity is returned. Otherwise, nil is returned.
func (s Service) authenticate(ctx context.Context, username, plainPwd string) (Identity, error) {

	ctx, span := otel.Start(ctx)
	defer span.End()

	repoUser := s.repoRegitry.GetUserRepository()
	user, err := repoUser.GetByUsername(ctx, username)
	if err != nil {
		if err == ierr.ErrResourceNotFound {
			return nil, ierr.ErrInvalidCreds
		}
		return nil, err
	}

	if username == user.GetUsername() && password.ComparePasswords(user.GetPassword(), []byte(plainPwd)) {
		// user is not active
		if !user.IsActive {
			return nil, ierr.ErrUserIsNotActive
		}
		// authentication successful
		return user, nil
	}

	// authentication failed
	return nil, ierr.ErrInvalidCreds

}

// generateJWT generates a JWT that encodes an iddomain.
func (s Service) generateJWT(ctx context.Context, identity Identity) (accessToken string, expiresAt time.Time, refreshToken string, err error) {

	ctx, span := otel.Start(ctx)
	defer span.End()

	//generate access token
	accessToken, expiresAt, err = s.generateAccessToken(ctx, identity)
	if err != nil {
		return
	}
	// generate refresh token
	refreshToken, err = s.generateRefreshToken(ctx, identity)
	if err != nil {
		return
	}
	user := domain.User{
		ID:           identity.GetID(),
		RefreshToken: &refreshToken,
	}
	repoUser := s.repoRegitry.GetUserRepository()
	err = repoUser.Update(ctx, identity.GetID(), user)
	return
}

func (s Service) generateAccessToken(ctx context.Context, identity Identity) (accessToken string, expiresAt time.Time, err error) {

	_, span := otel.Start(ctx)
	defer span.End()

	expiresAt = times.Now().Add(time.Duration(s.cfg.JWT.TokenExpiration) * time.Minute)
	expiresAtUnix := times.Now().Add(time.Duration(s.cfg.JWT.TokenExpiration) * time.Minute).Unix()
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         identity.GetID(),
		"username":   identity.GetUsername(),
		"user_type":  identity.GetType(),
		"exp":        expiresAtUnix,
		"token_type": "access",
	}).SignedString([]byte(s.cfg.JWT.SigningKey))
	err = errors.Wrap(err, "cannot generate token")
	return
}

func (s Service) generateRefreshToken(ctx context.Context, identity Identity) (refreshToken string, err error) {

	_, span := otel.Start(ctx)
	defer span.End()

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         identity.GetID(),
		"user_type":  identity.GetType(),
		"exp":        times.Now().AddDate(1000, 0, 0).Unix(),
		"token_type": "refresh",
	}).SignedString([]byte(s.cfg.JWT.SigningKey))
	err = errors.Wrap(err, "cannot generate token")
	return
}
