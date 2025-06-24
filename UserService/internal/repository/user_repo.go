package repository

import (
	"context"

	"github.com/DevisArya/BE-challenge-syn/UserService/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, tx *gorm.DB, user *entity.User) (uint, error)
	Update(ctx context.Context, tx *gorm.DB, user *entity.User) error
	Delete(ctx context.Context, tx *gorm.DB, userId uint) error
	FindByName(ctx context.Context, name string) (*entity.User, error)
	FindById(ctx context.Context, userId uint) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindAll(ctx context.Context, limit, offset int) ([]*entity.User, *int64, error)
}

type UserRepositoryImpl struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewUserRepository(log *logrus.Logger, db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		Log: log,
		DB:  db,
	}
}

func (r *UserRepositoryImpl) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.DB
}

// Save implements UserRepository
func (r *UserRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, user *entity.User) (uint, error) {

	db := r.getDB(tx)

	if err := db.WithContext(ctx).
		Create(user).
		Error; err != nil {
		r.Log.Errorf("Failed to create user: %v", err)
		return 0, err
	}

	r.Log.Infof("DARIII USECASEEE %d ", user.Id)

	return user.Id, nil
}

// Update implements UserRepository
func (r *UserRepositoryImpl) Update(ctx context.Context, tx *gorm.DB, userData *entity.User) error {

	db := r.getDB(tx)
	if err := db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", userData.Id).
		Updates(userData).
		Error; err != nil {
		r.Log.Errorf("Failed to update user: %v", err)
		return err
	}

	return nil
}

// Delete implements UserRepository
func (r *UserRepositoryImpl) Delete(ctx context.Context, tx *gorm.DB, userId uint) error {
	db := r.getDB(tx)

	if err := db.WithContext(ctx).
		Delete(&entity.User{}, userId).
		Error; err != nil {
		r.Log.Errorf("Failed to delete user: %v", err)
		return err
	}

	return nil
}

// FindByEmail implements UserRepository
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User

	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {

		if err != gorm.ErrRecordNotFound {
			r.Log.Errorf("Error finding user by email: %v", err)

		}
		return nil, err
	}

	return &user, nil
}

// FindById implements userRepository
func (r *UserRepositoryImpl) FindByName(ctx context.Context, name string) (*entity.User, error) {
	var user entity.User

	if err := r.DB.WithContext(ctx).
		Where("name = ?", name).
		First(&user).Error; err != nil {
		r.Log.Errorf("Error finding user by name: %v", err)
		return nil, err
	}

	return &user, nil
}

// FindById implements UserRepository
func (r *UserRepositoryImpl) FindById(ctx context.Context, userId uint) (*entity.User, error) {
	var user entity.User

	if err := r.DB.WithContext(ctx).First(&user, userId).Error; err != nil {
		r.Log.Errorf("Error finding user by id: %v", err)
		return nil, err
	}

	return &user, nil
}

// FindAll implements UserRepository
func (r *UserRepositoryImpl) FindAll(ctx context.Context, limit, offset int) ([]*entity.User, *int64, error) {

	var users []*entity.User
	var count int64

	query := r.DB.WithContext(ctx).Model(&entity.User{})

	if err := query.Count(&count).Error; err != nil {
		r.Log.Errorf("Error counting users: %v", err)
		return nil, nil, err
	}
	if err := query.
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		r.Log.Errorf("Error fetching users: %v", err)
		return nil, nil, err
	}

	return users, &count, nil
}
