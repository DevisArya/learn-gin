package delivery

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/DevisArya/BE-challenge-syn/UserService/internal/dto"
	"github.com/DevisArya/BE-challenge-syn/UserService/internal/helper"
	"github.com/DevisArya/BE-challenge-syn/UserService/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserController interface {
	Login(c *gin.Context)
	// Logout(c *gin.Context)
	Create(c *gin.Context)
	// UpdatePassword(c *gin.Context)
	// UpdateEmail(c *gin.Context)
	// UpdateNameProfile(c *gin.Context)
	// Delete(c *gin.Context)
	// FindById(c *gin.Context)
	FindAll(c *gin.Context)
}

type UserControllerImpl struct {
	UserUseCase usecase.UserUseCase
	Log         *logrus.Logger
}

func NewUserController(UserUseCase usecase.UserUseCase, log *logrus.Logger) UserController {
	return &UserControllerImpl{
		UserUseCase: UserUseCase,
		Log:         log,
	}
}

// Login implements AuthHandler.
func (ctrl *UserControllerImpl) Login(c *gin.Context) {

	var req dto.UserLoginRequest

	if err := c.ShouldBind(&req); err != nil {
		ctrl.Log.Warnf("[Login] Failed to parse request body: %v", err)
		helper.SendErrorResponse(c, http.StatusBadRequest, []string{"invalid payload"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	token, err := ctrl.UserUseCase.Login(ctx, &req)
	if err != nil {
		helper.SendCustomErrorResponse(c, err)
		return
	}

	data := dto.UserLoginResponse{
		Token: token,
	}

	helper.SendSuccessResponseWithData(c, http.StatusOK, "Login Successfully", data)
}

// // Logout implements AuthHandler.
// func (ctrl *UserControllerImpl) Logout(c *gin.Context) {

// 	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
// 	defer cancel()

// 	userID, err := helper.GetUserId(c)
// 	if err != nil {
// 		ctrl.Log.Errorf("[Logout] Failed to get user id: %v", err)
// 		helper.SendErrorResponse(c, http.StatusUnauthorized, []string{err.Error()})
// 	}

// 	if err := ctrl.UserUseCase.Logout(ctx, userID); err != nil {
// 		helper.SendCustomErrorResponse(c, err)
// 	}

// 	helper.SendSuccessResponse(c, http.StatusOK, "Logout Successfully")
// }

// Create implements CustomerController.
func (ctrl *UserControllerImpl) Create(c *gin.Context) {
	var req dto.UserCreateRequest

	//bind ke struct req
	if err := c.ShouldBind(&req); err != nil {
		ctrl.Log.Warnf("[Create] Failed to parse request body: %v", err)
		helper.SendErrorResponse(c, http.StatusBadRequest, []string{"invalid request payload"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	data, err := ctrl.UserUseCase.Create(ctx, &req)

	if err != nil {
		helper.SendCustomErrorResponse(c, err)
		return
	}

	helper.SendSuccessResponseWithData(c, http.StatusCreated, "user created successfully", data)
}

// // Delete implements CustomerController.
// func (ctrl *UserControllerImpl) Delete(c *gin.Context) {

// 	userID, err := helper.GetUserId(c)
// 	if err != nil {
// 		ctrl.Log.Errorf("[Delete] Failed to get user id: %v", err)
// 		helper.SendErrorResponse(c, http.StatusUnauthorized, []string{err.Error()})
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
// 	defer cancel()

// 	if err := ctrl.UserUseCase.Delete(ctx, userID); err != nil {
// 		helper.SendCustomErrorResponse(c, err)
// 	}

// 	helper.SendSuccessResponse(c, http.StatusOK, "customer deleted sucessfully")
// }

// FindAll implements CustomerController.
func (ctrl *UserControllerImpl) FindAll(c *gin.Context) {

	var limit, page int

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Parse query parameters
	if limitParam := c.Query("limit"); limitParam != "" {
		if val, err := strconv.Atoi(limitParam); err == nil {
			limit = val
		} else {
			ctrl.Log.Warnf("[FindAll] invalid limit parameter: %v", err)
			helper.SendErrorResponse(c, http.StatusBadRequest, []string{"invalid limit parameter"})
		}
	}
	if pageParam := c.Query("offset"); pageParam != "" {
		if val, err := strconv.Atoi(pageParam); err == nil {
			page = val
		} else {
			ctrl.Log.Warnf("[FindAll] invalid page parameter: %v", err)
			helper.SendErrorResponse(c, http.StatusBadRequest, []string{"invalid page parameter"})
		}
	}

	data, paging, err := ctrl.UserUseCase.FindAll(ctx, limit, page)

	if err != nil {
		helper.SendCustomErrorResponse(c, err)
		return
	}

	helper.SendSuccessResponseWithPagination(c, http.StatusOK, "find users succesfully", data, paging)
}

// // FindById implements CustomerController.
// func (ctrl *UserControllerImpl) FindById(c *gin.Context) {
// 	userID, err := helper.GetUserId(c)
// 	if err != nil {
// 		ctrl.Log.Errorf("[FindById] Failed to get user id: %v", err)
// 		helper.SendErrorResponse(c, http.StatusUnauthorized, []string{err.Error()})
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
// 	defer cancel()

// 	data, err := ctrl.UserUseCase.FindById(ctx, userID)

// 	if err != nil {
// 		helper.SendCustomErrorResponse(c, err)
// 	}

// 	helper.SendSuccessResponseWithData(c, http.StatusOK, "find user succesfully", data)
// }

// // UpdateEmail implements CustomerController.
// func (ctrl *UserControllerImpl) UpdateEmail(c *gin.Context) {
// 	var req dto.UserupdateEmailRequest

// 	if err := c.ShouldBind(&req); err != nil {
// 		ctrl.Log.Warnf("[UpdateEmail] Failed to parse request body: %v", err)
// 		helper.SendErrorResponse(c, http.StatusBadRequest, []string{"invalid request payload"})
// 	}

// 	userID, err := helper.GetUserId(c)
// 	if err != nil {
// 		ctrl.Log.Errorf("[UpdateEmail] Failed to get user id: %v", err)
// 		helper.SendErrorResponse(c, http.StatusUnauthorized, []string{err.Error()})
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
// 	defer cancel()

// 	if err := ctrl.UserUseCase.UpdateEmail(ctx, &req, userID); err != nil {
// 		helper.SendCustomErrorResponse(c, err)
// 	}

// 	helper.SendSuccessResponse(c, http.StatusOK, "email updated successfully")
// }

// UpdateNameProfile implements CustomerController.
// func (ctrl *UserControllerImpl) UpdateNameProfile(c *gin.Context) {
// 	var req dto.UserUpdateNameProfileRequest

// 	if err := c.ShouldBind(&req); err != nil {
// 		ctrl.Log.Warnf("[UpdateName] Failed to parse request body: %v", err)
// 		helper.SendErrorResponse(c, http.StatusBadRequest, []string{"invalid request payload"})
// 	}

// 	userID, err := helper.GetUserId(c)
// 	if err != nil {
// 		ctrl.Log.Errorf("[UpdateName] Failed to get user id: %v", err)
// 		helper.SendErrorResponse(c, http.StatusUnauthorized, []string{err.Error()})
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
// 	defer cancel()

// 	if err := ctrl.UserUseCase.UpdateNameProfile(ctx, &req, userID); err != nil {
// 		helper.SendCustomErrorResponse(c, err)
// 	}

// 	helper.SendSuccessResponse(c, http.StatusOK, "name profile updated successfully")
// }

// UpdatePassword implements CustomerController.
// func (ctrl *UserControllerImpl) UpdatePassword(c *gin.Context) {
// 	var req dto.UserupdatePasswordRequest

// 	if err := c.ShouldBind(&req); err != nil {
// 		ctrl.Log.Warnf("[UpdatePassword] Failed to parse request body: %v", err)
// 		helper.SendErrorResponse(c, http.StatusBadRequest, []string{"invalid request payload"})
// 	}

// 	userID, err := helper.GetUserId(c)
// 	if err != nil {
// 		ctrl.Log.Errorf("[UpdatePassword] Failed to get user id: %v", err)
// 		helper.SendErrorResponse(c, http.StatusUnauthorized, []string{err.Error()})
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
// 	defer cancel()

// 	if err := ctrl.UserUseCase.UpdatePassword(ctx, &req, userID); err != nil {
// 		helper.SendCustomErrorResponse(c, err)
// 	}

// 	helper.SendSuccessResponse(c, http.StatusOK, "password updated successfully")
// }
