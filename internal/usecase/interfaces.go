package usecase

import (
	"github.com/hctilf/go-test-task-medods/internal/entity"
)

type (
	TokensModelInterface interface {
		GetTokenByUserGUID(ipAddress string, guid string) (*entity.RefreshToken, error)
		RefreshToken(refreshToken, ipAddress string) (bool, error)
		CreateTestToken(ipAddress string) (*entity.RefreshToken, error)
	}
)
