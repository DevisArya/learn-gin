package mock

import (
	"context"
	"errors"

	"github.com/DevisArya/BE-challenge-syn/UserService/internal/entity"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) Create(ctx context.Context, tx *gorm.DB, user *entity.User) (uint, error) {
	args := m.Called(ctx, tx, user)

	id, ok := args.Get(0).(uint)
	if !ok {
		return 0, errors.New("error assertion")
	}
	return id, args.Error(1)
}
func (m *UserRepositoryMock) Update(ctx context.Context, tx *gorm.DB, user *entity.User) error {
	args := m.Called(ctx, tx, user)
	return args.Error(0)
}
func (m *UserRepositoryMock) Delete(ctx context.Context, tx *gorm.DB, userId uint) error {
	args := m.Called(ctx, tx, userId)
	return args.Error(0)
}
func (m *UserRepositoryMock) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	res, ok := args.Get(0).(*entity.User)
	if !ok {
		return nil, errors.New("error assertion")
	}

	return res, args.Error(1)
}
func (m *UserRepositoryMock) FindByName(ctx context.Context, name string) (*entity.User, error) {
	args := m.Called(ctx, name)
	res, ok := args.Get(0).(*entity.User)
	if !ok {
		return nil, errors.New("error assertion")
	}

	return res, args.Error(1)
}
func (m *UserRepositoryMock) FindById(ctx context.Context, userId uint) (*entity.User, error) {
	args := m.Called(ctx, userId)
	res, ok := args.Get(0).(*entity.User)
	if !ok {
		return nil, errors.New("error assertion")
	}

	return res, args.Error(1)
}
func (m *UserRepositoryMock) FindAll(ctx context.Context, limit, offset int) ([]*entity.User, *int64, error) {
	args := m.Called(ctx, limit, offset)
	entityUser, ok := args.Get(0).([]*entity.User)
	if !ok {
		return nil, nil, errors.New("error assertion entity user")
	}
	totalRecord, _ := args.Get(1).(*int64)
	if !ok {
		return nil, nil, errors.New("error assertion total record")
	}

	return entityUser, totalRecord, args.Error(2)
}
