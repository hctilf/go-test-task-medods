package tokens

import (
	"errors"
	"fmt"

	"github.com/hctilf/go-test-task-medods/internal/entity"
	uc "github.com/hctilf/go-test-task-medods/internal/usecase"
	bt "github.com/hctilf/go-test-task-medods/pkg/bcrypt_tools"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokensRepository struct {
	DB *gorm.DB
}

var _ uc.TokensModelInterface = (*TokensRepository)(nil)

func NewTokensRepository(db *gorm.DB) *TokensRepository {
	return &TokensRepository{
		DB: db,
	}
}

func (r *TokensRepository) GetTokenByUserGUID(ipAddress string, guid string) (*entity.RefreshToken, error) {
	var token entity.RefreshToken

	q := r.DB.Model(&entity.RefreshToken{}).
		Where("user_guid = ?", guid).
		First(&token)

	if q.Error != nil {
		if errors.Is(q.Error, gorm.ErrRecordNotFound) {
			newRefreshToken, err := bt.HashToken(ipAddress)
			if err != nil {
				return &token, err
			}
			token = entity.RefreshToken{
				UserGUID:  uuid.New(),
				IpAddress: ipAddress,
				Token:     newRefreshToken,
			}

			if err := r.DB.Create(&token).Error; err != nil {
				return &token, fmt.Errorf("cannot create token: %w", err)
			}
		}

		return nil, q.Error
	}

	return &token, nil
}

func (r *TokensRepository) RefreshToken(refreshToken, ipAddress string) (bool, error) {
	var (
		isIpChanged bool
		token       entity.RefreshToken
	)
	q := r.DB.Model(&entity.RefreshToken{}).
		Where("token = ?", refreshToken).
		First(&token)

	if q.Error != nil {
		return isIpChanged, q.Error
	}

	if token.IpAddress != ipAddress {
		token.IpAddress = ipAddress
		isIpChanged = true
	}

	newRefreshToken, err := bt.HashToken(ipAddress)
	if err != nil {
		return isIpChanged, err
	}

	token.Token = newRefreshToken
	if err := r.DB.Table("refresh_tokens").
		Where("id = ?", token.Id).
		Updates(map[string]interface{}{"ip_address": token.IpAddress, "token": token.Token}).Error; err != nil {
		return isIpChanged, fmt.Errorf("cannot update token: %w", err)
	}

	return isIpChanged, nil
}

func (r *TokensRepository) CreateTestToken(ipAddress string) (*entity.RefreshToken, error) {
	newRefreshToken, err := bt.HashToken(ipAddress)
	if err != nil {
		return nil, err
	}
	token := entity.RefreshToken{
		UserGUID:  uuid.New(),
		IpAddress: ipAddress,
		Token:     newRefreshToken,
	}

	if err := r.DB.Create(&token).Error; err != nil {
		return nil, fmt.Errorf("cannot create token: %w", err)
	}

	return &token, nil
}
