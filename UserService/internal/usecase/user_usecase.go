package usecase

import (
	"context"
	"net/http"

	"github.com/DevisArya/BE-challenge-syn/UserService/client/notification"
	"github.com/DevisArya/BE-challenge-syn/UserService/internal/dto"
	"github.com/DevisArya/BE-challenge-syn/UserService/internal/entity"
	"github.com/DevisArya/BE-challenge-syn/UserService/internal/helper"
	"github.com/DevisArya/BE-challenge-syn/UserService/internal/repository"
	"github.com/DevisArya/BE-challenge-syn/UserService/internal/security"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserUseCase interface {
	Login(ctx context.Context, request *dto.UserLoginRequest) (string, error)
	Logout(ctx context.Context, id uint) error
	VerifyToken(ctx context.Context, request *dto.VerifyTokenRequest, id uint) error
	Create(ctx context.Context, request *dto.UserCreateRequest) (*dto.UserCreateResponse, error)
	UpdatePassword(ctx context.Context, request *dto.UserupdatePasswordRequest, id uint) error
	UpdateEmail(ctx context.Context, request *dto.UserupdateEmailRequest, id uint) error
	UpdateNameProfile(ctx context.Context, request *dto.UserUpdateNameProfileRequest, id uint) error
	Delete(ctx context.Context, id uint) error
	FindById(ctx context.Context, id uint) (*dto.UserResponse, error)
	FindAll(ctx context.Context, limit, page int) ([]*dto.UserResponse, *dto.PaginationResponse, error)
}

type UserUseCaseImpl struct {
	UserRepository repository.UserRepository
	validate       *validator.Validate
	notifClient    notification.NotificationClient
	Log            *logrus.Logger
	Helper         helper.HelperTx
}

func NewUserUseCase(userRepository repository.UserRepository, validate *validator.Validate, notifClient notification.NotificationClient, log *logrus.Logger, Helper helper.HelperTx) UserUseCase {
	return &UserUseCaseImpl{
		UserRepository: userRepository,
		validate:       validate,
		notifClient:    notifClient,
		Log:            log,
		Helper:         Helper,
	}
}

// Login implements UserUseCase.
func (uc *UserUseCaseImpl) Login(ctx context.Context, request *dto.UserLoginRequest) (string, error) {

	if validationMessages := helper.ValidateStruct(uc.validate, request); len(validationMessages) > 0 {

		uc.Log.Warnf("[Login] Validation failed: %v", validationMessages)

		return "", helper.NewCustomError(http.StatusBadRequest, validationMessages)
	}

	user, err := uc.UserRepository.FindByName(ctx, request.Name)
	if err != nil {
		return "", helper.NewCustomError(http.StatusNotFound, []string{err.Error()})
	}

	if err := security.VerifyPassword(user.Password, request.Password); err != nil {
		uc.Log.Warnf("[Login] Invalid password for user: %s", request.Name)
		return "", helper.NewCustomError(http.StatusBadRequest, []string{"Wrong password"})
	}
	token, err := security.CreateToken(user.Id, user.Name)
	if err != nil {
		uc.Log.Errorf("[Login] Token creation failed: %v", err)
		return "", err
	}

	userData := entity.User{
		Id:    user.Id,
		Token: token,
	}

	if err := uc.UserRepository.Update(ctx, nil, &userData); err != nil {
		return "", err
	}

	uc.Log.Infof("[Login] User %s logged in successfully", request.Name)
	return token, nil
}

// Logout implements UserUseCase.
func (uc *UserUseCaseImpl) Logout(ctx context.Context, id uint) error {
	if id == 0 {
		uc.Log.Warn("[Logout] Invalid user ID")
		return helper.NewCustomError(http.StatusBadRequest, []string{"invalid user id"})
	}

	userData := entity.User{
		Id:    id,
		Token: "",
	}

	if err := uc.UserRepository.Update(ctx, nil, &userData); err != nil {
		return err
	}

	uc.Log.Infof("[Logout] User ID %d logged out successfully", id)

	return nil
}

// VerifyToken implements UserUseCase.
func (uc *UserUseCaseImpl) VerifyToken(ctx context.Context, request *dto.VerifyTokenRequest, id uint) error {

	if validationMessages := helper.ValidateStruct(uc.validate, request); len(validationMessages) > 0 {
		uc.Log.Warnf("[VerifyToken] Validation failed: %v", validationMessages)
		return helper.NewCustomError(http.StatusBadRequest, validationMessages)
	}
	user, err := uc.UserRepository.FindById(ctx, id)
	if err != nil {
		return err
	}

	if user.Token != request.Token {
		uc.Log.Warn("[VerifyToken] Invalid token")
		return helper.NewCustomError(http.StatusUnauthorized, []string{"invalid token"})
	}

	uc.Log.Infof("[VerifyToken] Token verified successfully for user ID %d", id)
	return nil
}

// Create implements UserUseCase
func (uc *UserUseCaseImpl) Create(ctx context.Context, request *dto.UserCreateRequest) (*dto.UserCreateResponse, error) {

	if validationMessages := helper.ValidateStruct(uc.validate, request); len(validationMessages) > 0 {
		uc.Log.Warnf("[Create] Validation failed: %v", validationMessages)
		return nil, helper.NewCustomError(http.StatusBadRequest, validationMessages)
	}

	//check used email
	user, err := uc.UserRepository.FindByEmail(ctx, request.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if user != nil {
		uc.Log.Warn("[Create] Email already in use")
		return nil, helper.NewCustomError(http.StatusConflict, []string{"email already use"})
	}

	//hash password
	hashedPassword, err := security.HashPassword(request.Password)
	if err != nil {
		uc.Log.Errorf("[Create] hash password failed: %v", err)
		return nil, err
	}

	userData := entity.User{
		Email:    request.Email,
		Name:     request.Name,
		Password: hashedPassword,
	}

	tx := uc.Helper.Begin()
	defer uc.Helper.CommitOrRollback(tx)

	id, err := uc.UserRepository.Create(ctx, tx, &userData)

	if err != nil {
		return nil, err
	}

	notifData := dto.SendUserNotificationRequest{
		UserID:      uint32(id),
		Action:      "CREATE",
		Description: "create new user",
	}

	uc.Log.Infof("DARIII USECASEEE %d ", notifData.UserID)

	//send notification
	if err := uc.notifClient.SendUserNotification(ctx, &notifData); err != nil {
		uc.Log.Errorf("[Create] send create user notification failed: %v", err)

		return nil, err
	}

	uc.Log.Infof("[Create] user created successfully with ID: %d", id)
	return &dto.UserCreateResponse{
		Id: int(id),
	}, nil
}

// UpdatePassword implements UserUseCase
func (uc *UserUseCaseImpl) UpdatePassword(ctx context.Context, request *dto.UserupdatePasswordRequest, id uint) error {

	if validationMessages := helper.ValidateStruct(uc.validate, request); len(validationMessages) > 0 {
		uc.Log.Warnf("[UpdatePassword] Validation failed: %v", validationMessages)
		return helper.NewCustomError(http.StatusBadRequest, validationMessages)
	}

	user, err := uc.UserRepository.FindById(ctx, id)
	if err != nil {
		return err
	}

	// hash new password
	hashedPassword, err := security.HashPassword(request.Password)
	if err != nil {
		uc.Log.Errorf("[UpdatePassword] hash password failed: %v", err)
		return err
	}

	//validate new password same or not
	if user.Password == hashedPassword {
		uc.Log.Warn("[UpdatePassword] new password must be different from the current password")
		return helper.NewCustomError(http.StatusBadRequest, []string{"new password must be different from the current password"})
	}

	userData := entity.User{
		Id:       id,
		Password: hashedPassword,
	}
	tx := uc.Helper.Begin()
	defer uc.Helper.CommitOrRollback(tx)

	if err := uc.UserRepository.Update(ctx, tx, &userData); err != nil {
		return err
	}

	notifData := dto.SendUserNotificationRequest{
		UserID:      uint32(id),
		Action:      "UPDATE",
		Description: "update password",
	}

	//send notification
	if err := uc.notifClient.SendUserNotification(ctx, &notifData); err != nil {
		uc.Log.Errorf("[UpdatePassword] send update password notification failed: %v", err)

		return err
	}

	uc.Log.Infof("[UpdatePassword] User ID %d update password successfully", id)
	return nil
}

// UpdateEmail implements UserUseCase
func (uc *UserUseCaseImpl) UpdateEmail(ctx context.Context, request *dto.UserupdateEmailRequest, id uint) error {

	if validationMessages := helper.ValidateStruct(uc.validate, request); len(validationMessages) > 0 {
		uc.Log.Warnf("[UpdateEmail] Validation failed: %v", validationMessages)
		return helper.NewCustomError(http.StatusBadRequest, validationMessages)
	}

	//check used email
	user, err := uc.UserRepository.FindByEmail(ctx, request.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if user != nil {
		uc.Log.Warn("[UpdateEmail] Email already in use")
		return helper.NewCustomError(http.StatusConflict, []string{"email already use"})
	}

	userData := entity.User{
		Id:    id,
		Email: request.Email,
	}

	tx := uc.Helper.Begin()
	defer uc.Helper.CommitOrRollback(tx)

	if err := uc.UserRepository.Update(ctx, tx, &userData); err != nil {
		return err
	}

	notifData := dto.SendUserNotificationRequest{
		UserID:      uint32(id),
		Action:      "UPDATE",
		Description: "update email",
	}

	//send notification
	if err := uc.notifClient.SendUserNotification(ctx, &notifData); err != nil {
		uc.Log.Errorf("[UpdateEmail] send update email notification failed: %v", err)

		return err
	}

	uc.Log.Infof("[UpdateEmail]  User ID %d update email successfully", id)

	return nil
}

// UpdateProfile implements UserUseCase
func (uc *UserUseCaseImpl) UpdateNameProfile(ctx context.Context, request *dto.UserUpdateNameProfileRequest, id uint) error {

	if validationMessages := helper.ValidateStruct(uc.validate, request); len(validationMessages) > 0 {
		uc.Log.Warnf("[UpdateName] Validation failed: %v", validationMessages)
		return helper.NewCustomError(http.StatusBadRequest, validationMessages)
	}

	tx := uc.Helper.Begin()
	defer uc.Helper.CommitOrRollback(tx)

	userData := entity.User{
		Id:   id,
		Name: request.Name,
	}

	if err := uc.UserRepository.Update(ctx, tx, &userData); err != nil {
		return err
	}

	notifData := dto.SendUserNotificationRequest{
		UserID:      uint32(id),
		Action:      "UPDATE",
		Description: "update name profile",
	}

	//send notification
	if err := uc.notifClient.SendUserNotification(ctx, &notifData); err != nil {
		uc.Log.Errorf("[UpdateName] send update name notification failed: %v", err)

		return err
	}

	uc.Log.Infof("[UpdateName]  User ID %d update name successfully", id)

	return nil
}

// Delete implements UserUseCase
func (uc *UserUseCaseImpl) Delete(ctx context.Context, id uint) error {

	if _, err := uc.UserRepository.FindById(ctx, id); err != nil {
		return err
	}

	tx := uc.Helper.Begin()
	defer uc.Helper.CommitOrRollback(tx)

	if err := uc.UserRepository.Delete(ctx, tx, id); err != nil {
		return err
	}

	notifData := dto.SendUserNotificationRequest{
		UserID:      uint32(id),
		Action:      "DELETE",
		Description: "delete user",
	}

	//send notification
	if err := uc.notifClient.SendUserNotification(ctx, &notifData); err != nil {
		uc.Log.Errorf("[Delete] send deleted user notification failed: %v", err)

		return err
	}

	uc.Log.Infof("[Delete]  User ID %d deleted successfully", id)

	return nil
}

// FindById implements UserUseCase
func (uc *UserUseCaseImpl) FindById(ctx context.Context, id uint) (*dto.UserResponse, error) {
	user, err := uc.UserRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	uc.Log.Infof("[FindById] find data User ID %d successfully", id)

	return helper.ToUserResponse(user), nil
}

// FindAll implements UserUseCase
func (uc *UserUseCaseImpl) FindAll(ctx context.Context, limit, page int) ([]*dto.UserResponse, *dto.PaginationResponse, error) {
	offset := (page - 1) * limit

	users, count, err := uc.UserRepository.FindAll(ctx, int(limit), int(offset))
	if err != nil {
		return nil, nil, err
	}

	totalRecord := int(*count)
	totalPage := totalRecord / (limit)

	uc.Log.Info("[FindAll] find all data user successfully")

	return helper.ToUsersResponses(users),
		&dto.PaginationResponse{
			CurrentPage: page,
			Limit:       limit,
			TotalRecord: totalRecord,
			TotalPage:   totalPage,
		},
		nil
}
