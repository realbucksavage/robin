package vhost

import (
	"context"
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/realbucksavage/robin/pkg/types"
	"github.com/realbucksavage/robin/pkg/vhosts"
)

var (
	ErrNotFound = errors.New("not found")
)

type Service interface {
	ListVhosts(ctx context.Context) ([]types.Vhost, error)
	GetVhost(ctx context.Context, id uint) (types.Vhost, error)
	PostVhost(ctx context.Context, v types.Vhost) (types.Vhost, error)
	PutVhost(ctx context.Context, id uint, v types.Vhost) (types.Vhost, error)
	DeleteVhost(ctx context.Context, id uint) error
}

type defaultService struct {
	db    *gorm.DB
	vault vhosts.Vault
}

func NewService(db *gorm.DB, vault vhosts.Vault) Service {
	return &defaultService{db: db, vault: vault}
}

func (s *defaultService) ListVhosts(ctx context.Context) ([]types.Vhost, error) {
	var v []types.Vhost
	if err := s.db.Find(&v).Error; err != nil {
		return nil, err
	}

	return v, nil
}

func (s *defaultService) GetVhost(ctx context.Context, id uint) (types.Vhost, error) {
	var v types.Vhost
	err := s.db.Preload("Cert").First(&v, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return types.Vhost{}, ErrNotFound
		} else {
			return types.Vhost{}, err
		}
	}

	return v, nil
}

func (s *defaultService) PostVhost(ctx context.Context, v types.Vhost) (types.Vhost, error) {
	db := s.db.BeginTx(ctx, nil)
	if db.Error != nil {
		return types.Vhost{}, db.Error
	}

	if err := db.Save(&v).Error; err != nil {
		db.Rollback()
		return types.Vhost{}, err
	}

	if err := s.vault.Put(v.FQDN, vhosts.H{
		FQDN:       v.FQDN,
		Origin:     v.Origin,
		PrivateKey: v.Cert.RSAKey,
		X509Cert:   v.Cert.X509,
	}); err != nil {
		db.Rollback()
		return types.Vhost{}, err
	}

	db.Commit()

	return v, nil
}

func (s *defaultService) PutVhost(ctx context.Context, id uint, v types.Vhost) (types.Vhost, error) {
	return types.Vhost{}, errors.New("implement me")
}

func (s *defaultService) DeleteVhost(ctx context.Context, id uint) error {
	db := s.db.BeginTx(ctx, nil)
	if db.Error != nil {
		return db.Error
	}

	var v types.Vhost
	err := db.Preload("Cert").Find(&v, id).Error
	if err != nil {
		db.Rollback()
		if gorm.IsRecordNotFoundError(err) {
			return ErrNotFound
		} else {
			return err
		}
	}

	if err := db.Delete(&v.Cert).Error; err != nil {
		db.Rollback()
		return err
	}

	if err := db.Delete(&v).Error; err != nil {
		db.Rollback()
		return err
	}

	s.vault.Remove(v.FQDN)

	db.Commit()
	return nil
}
