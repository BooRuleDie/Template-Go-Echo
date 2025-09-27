package auth

import (
	sharedAuthService "go-echo-template/internal/shared/auth"
	"go-echo-template/internal/shared/log"
	"go-echo-template/internal/shared/utils"

	"github.com/labstack/echo/v4"
)

type authService interface {
	login(c echo.Context, lr *LoginRequest) error
	logout(c echo.Context) error
	refresh(c echo.Context) error
}

type service struct {
	logger log.CustomLogger

	sharedAuth sharedAuthService.AuthService
	repo       authRepository
}

func NewAuthService(
	logger log.CustomLogger,
	repo authRepository,
	sharedAuth sharedAuthService.AuthService,
) authService {
	return &service{logger: logger, repo: repo, sharedAuth: sharedAuth}
}

func (s *service) login(c echo.Context, lr *LoginRequest) error {
	user, err := s.repo.getUserByEmail(c.Request().Context(), lr.Email)
	if err != nil {
		return err
	}

	if ok := utils.CheckPasswordHash(lr.Password, user.Password); !ok {
		return errEmailOrPasswordWrong
	}

	return s.sharedAuth.Login(c, &sharedAuthService.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone.String,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (s *service) logout(c echo.Context) error {
	return s.sharedAuth.Logout(c)
}

func (s *service) refresh(c echo.Context) error {
	return s.sharedAuth.Refresh(c, nil)
}
