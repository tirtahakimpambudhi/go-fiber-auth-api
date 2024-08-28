package mapper_test

import (
	"github.com/stretchr/testify/require"
	"github.com/tirtahakimpambudhi/restful_api/internal/entity"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/mapper"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/request"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/response"
	"testing"
)

func TestEntityUserToResponseConversion(t *testing.T) {
	user := &entity.Users{
		ID:        "123",
		Username:  "testuser",
		Email:     "testuser@example.com",
		CreatedAt: 1625097600,
		UpdatedAt: 1625097601,
	}

	expected := &response.User{
		ID:        "123",
		Username:  "testuser",
		Email:     "testuser@example.com",
		CreatedAt: 1625097600,
		UpdatedAt: 1625097601,
	}

	result := mapper.EntityUserToResponse(user)

	require.Equal(t, expected, result)
}

func TestEntitiesUserToResponsesConversion(t *testing.T) {
	user := []*entity.Users{
		{
			ID:        "123",
			Username:  "testuser",
			Email:     "testuser@example.com",
			CreatedAt: 1625097600,
			UpdatedAt: 1625097601,
		},
	}

	expected := []*response.User{
		{
			ID:        "123",
			Username:  "testuser",
			Email:     "testuser@example.com",
			CreatedAt: 1625097600,
			UpdatedAt: 1625097601,
		},
	}

	result := mapper.EntitiesUserToResponses(user)

	require.Equal(t, expected, result)
	require.Equal(t, len(expected), len(result))
}

func TestRequestUserToEntityConversion(t *testing.T) {
	reqUser := request.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	id := "12345"
	entityUser := mapper.RequestUserToEntity(id, reqUser)

	require.Equal(t, id, entityUser.ID)
	require.Equal(t, reqUser.Username, entityUser.Username)
	require.Equal(t, reqUser.Email, entityUser.Email)
	require.Equal(t, reqUser.Password, entityUser.Password)
}

func TestRequestUserEditToEntityConversion(t *testing.T) {
	reqUser := request.UserEdit{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	id := "12345"
	entityUser := mapper.RequestUserEditToEntity(id, reqUser)

	require.Equal(t, id, entityUser.ID)
	require.Equal(t, reqUser.Username, entityUser.Username)
	require.Equal(t, reqUser.Email, entityUser.Email)
	require.Equal(t, reqUser.Password, entityUser.Password)
}
