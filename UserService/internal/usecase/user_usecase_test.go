package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/DevisArya/BE-challenge-syn/UserService/internal/dto"
	"github.com/DevisArya/BE-challenge-syn/UserService/internal/entity"
	"github.com/DevisArya/BE-challenge-syn/UserService/internal/helper"
	mr "github.com/DevisArya/BE-challenge-syn/UserService/internal/repository/mock"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type MockHelperTx struct {
	mock.Mock
}

func (m *MockHelperTx) Begin() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockHelperTx) CommitOrRollback(tx *gorm.DB) {
	m.Called(tx)
}

func hashPassword(raw string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	return string(hash)
}

type MockNotifClient struct {
	mock.Mock
}

func (m *MockNotifClient) SendUserNotification(ctx context.Context, request *dto.SendUserNotificationRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}

func (m *MockNotifClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestLogin(t *testing.T) {
	type args struct {
		request *dto.UserLoginRequest
	}

	tests := []struct {
		name          string
		args          args
		mockSetup     func(m *mr.UserRepositoryMock)
		expectError   bool
		expectMessage string
	}{
		{
			name: "success login",
			args: args{
				request: &dto.UserLoginRequest{
					Name:     "devis",
					Password: "123456",
				},
			},
			mockSetup: func(m *mr.UserRepositoryMock) {
				hashed := hashPassword("123456")
				m.On("FindByName", mock.Anything, "devis").
					Return(&entity.User{Id: 1, Name: "devis", Password: hashed}, nil)
				m.On("Update", mock.Anything, mock.Anything, mock.AnythingOfType("*entity.User")).
					Return(nil)
			},
			expectError: false,
		},
		{
			name: "user not found",
			args: args{
				request: &dto.UserLoginRequest{
					Name:     "unknown",
					Password: "123456",
				},
			},
			mockSetup: func(m *mr.UserRepositoryMock) {
				m.On("FindByName", mock.Anything, mock.Anything).
					Return((*entity.User)(nil), errors.New("user not found"))
			},
			expectError:   true,
			expectMessage: "user not found",
		},
		{
			name: "wrong password",
			args: args{
				request: &dto.UserLoginRequest{
					Name:     "devis",
					Password: "wrongpass",
				},
			},
			mockSetup: func(m *mr.UserRepositoryMock) {
				hashed := hashPassword("123456")
				m.On("FindByName", mock.Anything, "devis").
					Return(&entity.User{Id: 1, Name: "devis", Password: hashed}, nil)
			},

			expectError:   true,
			expectMessage: "Wrong password",
		},
		{
			name: "invalid input",
			args: args{
				request: &dto.UserLoginRequest{},
			},
			mockSetup: func(m *mr.UserRepositoryMock) {
			},
			expectError: true,
		},
		{
			name: "fail update token",
			args: args{
				request: &dto.UserLoginRequest{
					Name:     "devis",
					Password: "123456",
				},
			},
			mockSetup: func(m *mr.UserRepositoryMock) {
				hashed := hashPassword("123456")
				m.On("FindByName", mock.Anything, "devis").Return(&entity.User{
					Id: 1, Name: "devis", Password: hashed}, nil)
				m.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("fail update token"))
			},
			expectError:   true,
			expectMessage: "fail update token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mr.UserRepositoryMock)
			mockHelperTx := new(MockHelperTx)
			tt.mockSetup(mockRepo)

			validate := validator.New()
			logger := logrus.New()
			errLoad := godotenv.Load("../../../.env")
			helper.PanicIfError(errLoad)

			uc := NewUserUseCase(mockRepo, validate, nil, logger, mockHelperTx)

			token, err := uc.Login(context.Background(), tt.args.request)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectMessage != "" {
					assert.Contains(t, err.Error(), tt.expectMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}

			mockRepo.AssertExpectations(t)
			mockHelperTx.AssertExpectations(t)
		})
	}
}

func TestLogout(t *testing.T) {

	tests := []struct {
		name          string
		id            uint
		mockSetup     func(m *mr.UserRepositoryMock)
		expectError   bool
		expectMessage string
	}{
		{
			name: "success logout",
			id:   1,
			mockSetup: func(m *mr.UserRepositoryMock) {
				m.On("Update", mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
			},
			expectError: false,
		},
		{
			name: "invalid user id",
			id:   0,
			mockSetup: func(m *mr.UserRepositoryMock) {
			},
			expectError: true,
		},
		{
			name: "fail update token",
			id:   1,
			mockSetup: func(m *mr.UserRepositoryMock) {
				m.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("fail update token"))
			},
			expectError:   true,
			expectMessage: "fail update token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mr.UserRepositoryMock)
			mockHelperTx := new(MockHelperTx)
			tt.mockSetup(mockRepo)

			validate := validator.New()
			logger := logrus.New()

			uc := NewUserUseCase(mockRepo, validate, nil, logger, mockHelperTx)

			err := uc.Logout(context.Background(), tt.id)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectMessage != "" {
					assert.Contains(t, err.Error(), tt.expectMessage)
				}
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
			mockHelperTx.AssertExpectations(t)
		})
	}
}

func TestVerifyToken(t *testing.T) {
	type args struct {
		request dto.VerifyTokenRequest
		id      uint
	}
	test := []struct {
		name          string
		args          args
		setupMock     func(m *mr.UserRepositoryMock)
		expectError   bool
		expectMessage string
	}{
		{
			name: "succes verify token",
			args: args{
				request: dto.VerifyTokenRequest{
					Token: "token123",
				},
				id: 1,
			},
			setupMock: func(m *mr.UserRepositoryMock) {
				m.On("FindById", mock.Anything, mock.Anything).
					Return(&entity.User{Token: "token123"}, nil)
			},
			expectError: false,
		},
		{
			name: "invalid input",
			args: args{
				request: dto.VerifyTokenRequest{},
				id:      1,
			},
			setupMock:   func(m *mr.UserRepositoryMock) {},
			expectError: true,
		},
		{
			name: "user not found",
			args: args{
				request: dto.VerifyTokenRequest{
					Token: "token123",
				},
				id: 1,
			},
			setupMock: func(m *mr.UserRepositoryMock) {
				m.On("FindById", mock.Anything, mock.Anything).
					Return((*entity.User)(nil), errors.New("user not found"))
			},
			expectError:   true,
			expectMessage: "user not found",
		},
		{
			name: "invalid token",
			args: args{
				request: dto.VerifyTokenRequest{
					Token: "token123",
				},
				id: 1,
			},
			setupMock: func(m *mr.UserRepositoryMock) {
				m.On("FindById", mock.Anything, mock.Anything).
					Return(&entity.User{Token: "token12345"}, nil)
			},
			expectError: true,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mr.UserRepositoryMock)
			mockHelperTx := new(MockHelperTx)
			tt.setupMock(mockRepo)

			validate := validator.New()
			logger := logrus.New()

			uc := NewUserUseCase(mockRepo, validate, nil, logger, mockHelperTx)

			err := uc.VerifyToken(context.Background(), &tt.args.request, tt.args.id)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectMessage != "" {
					assert.Contains(t, err.Error(), tt.expectMessage)
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockHelperTx.AssertExpectations(t)
		})
	}
}

func TestCreate(t *testing.T) {
	type args struct {
		request *dto.UserCreateRequest
	}
	tests := []struct {
		name          string
		args          args
		setupMock     func(m *mr.UserRepositoryMock, n *MockNotifClient, x *MockHelperTx)
		expectError   bool
		expectMessage string
	}{
		{
			name: "success create user",
			args: args{
				request: &dto.UserCreateRequest{
					Email:    "test@gmail.com",
					Name:     "testname",
					Password: "test12345",
				},
			},
			setupMock: func(m *mr.UserRepositoryMock, n *MockNotifClient, x *MockHelperTx) {
				m.On("FindByEmail", mock.Anything, mock.Anything).Return((*entity.User)(nil), gorm.ErrRecordNotFound)
				x.On("Begin").Return((*gorm.DB)(nil))
				x.On("CommitOrRollback", (*gorm.DB)(nil)).Return((*gorm.DB)(nil))
				m.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
				n.On("SendUserNotification", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "validation failed",
			args: args{
				request: &dto.UserCreateRequest{},
			},
			setupMock:   func(m *mr.UserRepositoryMock, n *MockNotifClient, x *MockHelperTx) {},
			expectError: true,
		},
		{
			name: "internal error",
			args: args{
				request: &dto.UserCreateRequest{
					Email:    "test@gmail.com",
					Name:     "testname",
					Password: "test12345",
				},
			},
			setupMock: func(m *mr.UserRepositoryMock, n *MockNotifClient, x *MockHelperTx) {
				m.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, errors.New("internal error"))
			},
			expectError:   true,
			expectMessage: "internal error",
		},
		{
			name: "email already in use",
			args: args{
				request: &dto.UserCreateRequest{
					Email:    "test@gmail.com",
					Name:     "testname",
					Password: "test12345",
				},
			},
			setupMock: func(m *mr.UserRepositoryMock, n *MockNotifClient, x *MockHelperTx) {
				m.On("FindByEmail", mock.Anything, mock.Anything).Return(&entity.User{
					Id:   1,
					Name: "test",
				}, nil)
			},
			expectError: true,
		},
		{
			name: "failed to save data",
			args: args{
				request: &dto.UserCreateRequest{
					Email:    "test@gmail.com",
					Name:     "testname",
					Password: "test12345",
				},
			},
			setupMock: func(m *mr.UserRepositoryMock, n *MockNotifClient, x *MockHelperTx) {
				m.On("FindByEmail", mock.Anything, mock.Anything).Return(&entity.User{
					Id:   1,
					Name: "test",
				}, nil)
				m.On("Create", mock.Anything, mock.Anything).Return(0, errors.New("failed to save data"))
			},
			expectError:   true,
			expectMessage: "fail to save data",
		},
		{
			name: "failed to send notification",
			args: args{
				request: &dto.UserCreateRequest{
					Email:    "test@gmail.com",
					Name:     "testname",
					Password: "test12345",
				},
			},
			setupMock: func(m *mr.UserRepositoryMock, n *MockNotifClient, x *MockHelperTx) {
				m.On("FindByEmail", mock.Anything, mock.Anything).Return(&entity.User{
					Id:   1,
					Name: "test",
				}, nil)
				m.On("Create", mock.Anything, mock.Anything).Return(1, nil)
				n.On("SendUserNotification", mock.Anything, mock.Anything).Return(errors.New("failed to send notification"))
			},
			expectError:   true,
			expectMessage: "failed to send notification",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mr.UserRepositoryMock)
			mockNotifClient := new(MockNotifClient)
			mockHelperTx := new(MockHelperTx)
			tt.setupMock(mockRepo, mockNotifClient, mockHelperTx)

			validate := validator.New()
			logger := logrus.New()

			uc := NewUserUseCase(mockRepo, validate, nil, logger, mockHelperTx)

			res, err := uc.Create(context.Background(), tt.args.request)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectMessage != "" {
					assert.Contains(t, err.Error(), tt.expectMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
			}

			mockRepo.AssertExpectations(t)
			mockHelperTx.AssertExpectations(t)
		})
	}
}

// func TestUserUseCaseImpl_UpdatePassword(t *testing.T) {
// 	type args struct {
// 		request *dto.UserupdatePasswordRequest
// 		id      uint
// 	}
// 	tests := []struct {
// 		name          string
// 		args          args
// 		setupMock     func(m *mr.UserRepositoryMock)
// 		expectError   bool
// 		expectMessage string
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockRepo := new(mr.UserRepositoryMock)
// 			mockHelperTx := new(MockHelperTx)
// 			tt.setupMock(mockRepo)

// 			validate := validator.New()
// 			logger := logrus.New()

// 			uc := NewUserUseCase(mockRepo, validate, nil, logger, mockHelperTx)

// 			err := uc.VerifyToken(context.Background(), &tt.args.request, tt.args.id)

// 			if tt.expectError {
// 				assert.Error(t, err)
// 				if tt.expectMessage != "" {
// 					assert.Contains(t, err.Error(), tt.expectMessage)
// 				}
// 			} else {
// 				assert.NoError(t, err)
// 			}

// 			mockRepo.AssertExpectations(t)
// 			mockHelperTx.AssertExpectations(t)
// 		})
// 	}
// }

// func TestUserUseCaseImpl_UpdateEmail(t *testing.T) {
// 	type args struct {
// 		request *dto.UserupdateEmailRequest
// 	}
// 	tests := []struct {
// 		name          string
// 		args          args
// 		setupMock     func(m *mr.UserRepositoryMock)
// 		expectError   bool
// 		expectMessage string
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockRepo := new(mr.UserRepositoryMock)
// 			mockHelperTx := new(MockHelperTx)
// 			tt.setupMock(mockRepo)

// 			validate := validator.New()
// 			logger := logrus.New()

// 			uc := NewUserUseCase(mockRepo, validate, nil, logger, mockHelperTx)

// 			err := uc.VerifyToken(context.Background(), &tt.args.request, tt.args.id)

// 			if tt.expectError {
// 				assert.Error(t, err)
// 				if tt.expectMessage != "" {
// 					assert.Contains(t, err.Error(), tt.expectMessage)
// 				}
// 			} else {
// 				assert.NoError(t, err)
// 			}

// 			mockRepo.AssertExpectations(t)
// 			mockHelperTx.AssertExpectations(t)
// 		})
// 	}
// }

// func TestUserUseCaseImpl_UpdateNameProfile(t *testing.T) {
// 	type args struct {
// 		request *dto.UserUpdateNameProfileRequest
// 	}
// 	tests := []struct {
// 		name          string
// 		args          args
// 		setupMock     func(m *mr.UserRepositoryMock)
// 		expectError   bool
// 		expectMessage string
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockRepo := new(mr.UserRepositoryMock)
// 			mockHelperTx := new(MockHelperTx)
// 			tt.setupMock(mockRepo)

// 			validate := validator.New()
// 			logger := logrus.New()

// 			uc := NewUserUseCase(mockRepo, validate, nil, logger, mockHelperTx)

// 			err := uc.VerifyToken(context.Background(), &tt.args.request, tt.args.id)

// 			if tt.expectError {
// 				assert.Error(t, err)
// 				if tt.expectMessage != "" {
// 					assert.Contains(t, err.Error(), tt.expectMessage)
// 				}
// 			} else {
// 				assert.NoError(t, err)
// 			}

// 			mockRepo.AssertExpectations(t)
// 			mockHelperTx.AssertExpectations(t)
// 		})
// 	}
// }

// func TestUserUseCaseImpl_Delete(t *testing.T) {
// 	type args struct {
// 		id uint
// 	}
// 	tests := []struct {
// 		name          string
// 		args          args
// 		setupMock     func(m *mr.UserRepositoryMock)
// 		expectError   bool
// 		expectMessage string
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockRepo := new(mr.UserRepositoryMock)
// 			mockHelperTx := new(MockHelperTx)
// 			tt.setupMock(mockRepo)

// 			validate := validator.New()
// 			logger := logrus.New()

// 			uc := NewUserUseCase(mockRepo, validate, nil, logger, mockHelperTx)

// 			err := uc.VerifyToken(context.Background(), &tt.args.request, tt.args.id)

// 			if tt.expectError {
// 				assert.Error(t, err)
// 				if tt.expectMessage != "" {
// 					assert.Contains(t, err.Error(), tt.expectMessage)
// 				}
// 			} else {
// 				assert.NoError(t, err)
// 			}

// 			mockRepo.AssertExpectations(t)
// 			mockHelperTx.AssertExpectations(t)
// 		})
// 	}
// }

// func TestUserUseCaseImpl_FindById(t *testing.T) {
// 	type args struct {
// 		id uint
// 	}
// 	tests := []struct {
// 		name          string
// 		args          args
// 		setupMock     func(m *mr.UserRepositoryMock)
// 		expectError   bool
// 		expectMessage string
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockRepo := new(mr.UserRepositoryMock)
// 			mockHelperTx := new(MockHelperTx)
// 			tt.setupMock(mockRepo)

// 			validate := validator.New()
// 			logger := logrus.New()

// 			uc := NewUserUseCase(mockRepo, validate, nil, logger, mockHelperTx)

// 			err := uc.VerifyToken(context.Background(), &tt.args.request, tt.args.id)

// 			if tt.expectError {
// 				assert.Error(t, err)
// 				if tt.expectMessage != "" {
// 					assert.Contains(t, err.Error(), tt.expectMessage)
// 				}
// 			} else {
// 				assert.NoError(t, err)
// 			}

// 			mockRepo.AssertExpectations(t)
// 			mockHelperTx.AssertExpectations(t)
// 		})
// 	}
// }

// func TestUserUseCaseImpl_FindAll(t *testing.T) {
// 	type args struct {
// 		limit int
// 		page  int
// 	}
// 	tests := []struct {
// 		name          string
// 		args          args
// 		setupMock     func(m *mr.UserRepositoryMock)
// 		expectError   bool
// 		expectMessage string
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockRepo := new(mr.UserRepositoryMock)
// 			mockHelperTx := new(MockHelperTx)
// 			tt.setupMock(mockRepo)

// 			validate := validator.New()
// 			logger := logrus.New()

// 			uc := NewUserUseCase(mockRepo, validate, nil, logger, mockHelperTx)

// 			err := uc.VerifyToken(context.Background(), &tt.args.request, tt.args.id)

// 			if tt.expectError {
// 				assert.Error(t, err)
// 				if tt.expectMessage != "" {
// 					assert.Contains(t, err.Error(), tt.expectMessage)
// 				}
// 			} else {
// 				assert.NoError(t, err)
// 			}

// 			mockRepo.AssertExpectations(t)
// 			mockHelperTx.AssertExpectations(t)
// 		})
// 	}
// }
