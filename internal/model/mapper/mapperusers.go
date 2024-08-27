package mapper

import (
	"github.com/tirtahakimpambudhi/restful_api/internal/entity"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/request"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/response"
)

// Converts a request.User to an entity.Users entity.
func RequestUserToEntity(id string, user request.User) *entity.Users {
	return &entity.Users{
		ID:       id,            // User ID
		Username: user.Username, // User's username
		Email:    user.Email,    // User's email
		Password: user.Password, // User's password
	}
}

// Converts a request.UserEdit to an entity.Users entity.
func RequestUserEditToEntity(id string, user request.UserEdit) *entity.Users {
	return &entity.Users{
		ID:       id,            // User ID
		Username: user.Username, // User's username
		Email:    user.Email,    // User's email
		Password: user.Password, // User's password
	}
}

// Converts an entity.Users entity to a response.User for response formatting.
func EntityUserToResponse(users *entity.Users) *response.User {
	return &response.User{
		ID:        users.ID,        // User ID
		Username:  users.Username,  // User's username
		Email:     users.Email,     // User's email
		CreatedAt: users.CreatedAt, // Account creation timestamp
		UpdatedAt: users.UpdatedAt, // Last update timestamp
	}
}

// Converts an []entity.Users entity to a []response.User for response formatting.
func EntitiesUserToResponses(users []*entity.Users) []*response.User {
	responses := []*response.User{}
	for _, user := range users {
		responses = append(responses, &response.User{
			ID:        user.ID,        // User ID
			Username:  user.Username,  // User's username
			Email:     user.Email,     // User's email
			CreatedAt: user.CreatedAt, // Account creation timestamp
			UpdatedAt: user.UpdatedAt, // Last update timestamp
		})
	}
	return responses
}
