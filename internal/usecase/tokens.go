package usecase

import (
	"github.com/hctilf/go-test-task-medods/internal/entity"
)

type TokensUsecase struct {
	repo TokensModelInterface
}

func NewTokensUsecase(repo TokensModelInterface) *TokensUsecase {
	return &TokensUsecase{
		repo: repo,
	}
}

var _ TokensModelInterface = (*TokensUsecase)(nil)

func (u *TokensUsecase) GetTokenByUserGUID(ipAddress string, guid string) (*entity.RefreshToken, error) {
	return u.repo.GetTokenByUserGUID(ipAddress, guid)
}

func (u *TokensUsecase) RefreshToken(refreshToken, ipAddress string) (bool, error) {
	return u.repo.RefreshToken(refreshToken, ipAddress)
}

func (u *TokensUsecase) CreateTestToken(ipAddress string) (*entity.RefreshToken, error) {
	return u.repo.CreateTestToken(ipAddress)
}
