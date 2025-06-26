package helper

import (
	"github.com/DevisArya/BE-challenge-syn/UserService/internal/dto"
	"github.com/DevisArya/BE-challenge-syn/UserService/internal/entity"
)

func ToUserResponse(user *entity.User) *dto.UserResponse {
	return &dto.UserResponse{
		Id:    user.Id,
		Email: user.Email,
		Name:  user.Name,
	}
}

func ToUsersResponses(user []*entity.User) []*dto.UserResponse {

	var userResponses []*dto.UserResponse

	for _, val := range user {
		userResponses = append(userResponses, ToUserResponse(val))
	}
	return userResponses
}
