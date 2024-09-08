package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/token"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/mapper"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/response"
	"net/http"
	"reflect"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/phuslu/log"
	"github.com/segmentio/ksuid"

	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/hash"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/timeout"
	"github.com/tirtahakimpambudhi/restful_api/internal/entity"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/request"
	"github.com/tirtahakimpambudhi/restful_api/internal/repository"
	"github.com/tirtahakimpambudhi/restful_api/internal/usecase"
	"github.com/tirtahakimpambudhi/restful_api/internal/validation"
)

var (
	usersusecase  *usecase.UsersUsecase
	authusecase   *usecase.AuthUsecase
	usersRepoMock *repository.UsersRepositoryMock
	cacheRepoMock *repository.MockCacheRepository[*entity.Users]
	jwtToken      *token.JWTToken
	secretKey     *token.SecretKey
	argon2id      *hash.Argon2
)

func SetEnv() func() {
	os.Setenv("DB_TIMEOUT", "4")
	os.Setenv("CACHE_TIMEOUT", "4")
	os.Setenv("DOWN_STREAM_TIMEOUT", "4")
	os.Setenv("SECRET_KEY_ACCESS_TOKEN", "a_very_secret_key_access_is_32_byt")
	os.Setenv("SECRET_KEY_REFRESH_TOKEN", "a_very_secret_key_that_is_32_byt")
	os.Setenv("SECRET_KEY_FP_TOKEN", "a_very_secret_key_that_is_32_byt")
	unsetFunc := func() {
		os.Unsetenv("DB_TIMEOUT")
		os.Unsetenv("CACHE_TIMEOUT")
		os.Unsetenv("DOWN_STREAM_TIMEOUT")
		os.Unsetenv("SECRET_KEY_ACCESS_TOKEN")
		os.Unsetenv("SECRET_KEY_REFRESH_TOKEN")
		os.Unsetenv("SECRET_KEY_FP_TOKEN")
	}
	return unsetFunc
}

func TestMain(m *testing.M) {
	// Initialize the required mocks and structs
	unset := SetEnv()
	defer unset()
	validate := validator.New()
	english := en.New()
	universalTranslate := ut.New(english, english)
	translator, _ := universalTranslate.GetTranslator("en")
	usersRepoMock = new(repository.UsersRepositoryMock)
	cacheRepoMock = new(repository.MockCacheRepository[*entity.Users])
	validator := validation.NewValidator(validate, translator)
	timeoutConfig, _ := timeout.NewConfig()
	argon2id, _ = hash.NewHashArgon2()
	jwtToken, secretKey, _ = token.NewJWTToken()
	usersusecase = usecase.NewUsersUsecaseBuilder().WithLogger(&log.DefaultLogger).WithUsersRepository(usersRepoMock).WithCacheRepository(cacheRepoMock).WithHashing(argon2id).WithTimeoutConfig(timeoutConfig).WithValidator(validator).Build()
	authusecase = usecase.NewAuthUsecaseBuilder().WithLogger(&log.DefaultLogger).WithUsersRepository(usersRepoMock).WithToken(jwtToken).WithSecretKey(secretKey).WithHashing(argon2id).WithTimeoutConfig(timeoutConfig).WithValidator(validator).Build()
	m.Run()
}

// ================================================ LIST CASES ===================================================================
func TestUsersUsecase_List_WhenCache(t *testing.T) {

	// Prepare the request and expected response
	req := &request.Page{Size: 10, Before: ksuid.New().String(), After: ksuid.New().String()} // Ensure Before and After are initialized
	expectedUsers := []*entity.Users{
		{ID: "1", Username: "John Doe", Email: "john@example.com"},
	}

	// Define the behavior of the mocked methods
	cacheRepoMock.On("GetFromCache", mock.Anything, fmt.Sprintf("users:all:size[%d]:before[%s]:after[%s]", req.Size, req.Before, req.After)).Return(expectedUsers, nil).Once()
	usersRepoMock.On("Count", mock.Anything).Return(int64(1), nil).Once()

	// Call the List method
	resp, err := usersusecase.List(context.Background(), req)
	// Assertions
	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.Status)
	require.Equal(t, "STATUS_OK", resp.Code)
	require.Equal(t, mapper.EntitiesUserToResponses(expectedUsers), resp.Data)
	require.Equal(t, int64(1), resp.Meta["total_data"])
	require.Equal(t, float64(1), resp.Meta["total_page"])

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_List_WhenCacheMiss(t *testing.T) {

	// Prepare the request and expected response
	req := &request.Page{Size: 10, Before: ksuid.New().String(), After: ksuid.New().String()} // Ensure Before and After are initialized
	expectedUsers := []*entity.Users{
		{ID: "1", Username: "John Doe", Email: "john@example.com"},
	}
	expectedKey := fmt.Sprintf("users:all:size[%d]:before[%s]:after[%s]", req.Size, req.Before, req.After)

	// Define the behavior of the mocked methods
	cacheRepoMock.On("GetFromCache", mock.Anything, expectedKey).Return(nil, nil).Once()
	usersRepoMock.On("GetAll", mock.Anything, req).Return(expectedUsers, nil).Once()
	cacheRepoMock.On("SetToCache", mock.Anything, expectedKey, expectedUsers).Return(nil).Once()
	usersRepoMock.On("Count", mock.Anything).Return(int64(1), nil).Once()

	// Call the List method
	resp, err := usersusecase.List(context.Background(), req)
	// Assertions
	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.Status)
	require.Equal(t, "STATUS_OK", resp.Code)
	require.Equal(t, mapper.EntitiesUserToResponses(expectedUsers), resp.Data)
	require.Equal(t, int64(1), resp.Meta["total_data"])
	require.Equal(t, float64(1), resp.Meta["total_page"])

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_List_WhenCacheMiss_AndDBTimeout(t *testing.T) {

	// Prepare the request and expected response
	req := &request.Page{Size: 10, Before: ksuid.New().String(), After: ksuid.New().String()} // Ensure Before and After are initialized
	expectedKey := fmt.Sprintf("users:all:size[%d]:before[%s]:after[%s]", req.Size, req.Before, req.After)

	// Define the behavior of the mocked methods
	cacheRepoMock.On("GetFromCache", mock.Anything, expectedKey).Return(nil, nil).Once()
	usersRepoMock.On("GetAll", mock.Anything, req).Return(nil, context.DeadlineExceeded).Once()

	// Call the List method
	resp, err := usersusecase.List(context.Background(), req)
	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_List_WhenCacheMiss_AndSetCacheErr(t *testing.T) {

	// Prepare the request and expected response
	req := &request.Page{Size: 10, Before: ksuid.New().String(), After: ksuid.New().String()} // Ensure Before and After are initialized
	expectedUsers := []*entity.Users{
		{ID: "1", Username: "John Doe", Email: "john@example.com"},
	}
	expectedKey := fmt.Sprintf("users:all:size[%d]:before[%s]:after[%s]", req.Size, req.Before, req.After)

	// Define the behavior of the mocked methods
	cacheRepoMock.On("GetFromCache", mock.Anything, expectedKey).Return(nil, nil).Once()
	usersRepoMock.On("GetAll", mock.Anything, req).Return(expectedUsers, nil).Once()
	cacheRepoMock.On("SetToCache", mock.Anything, expectedKey, expectedUsers).Return(context.DeadlineExceeded).Once()
	usersRepoMock.On("Count", mock.Anything).Return(int64(1), nil).Once()

	// Call the List method
	resp, err := usersusecase.List(context.Background(), req)
	// Assertions
	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.Status)
	require.Equal(t, "STATUS_OK", resp.Code)
	require.Equal(t, mapper.EntitiesUserToResponses(expectedUsers), resp.Data)
	require.Equal(t, int64(1), resp.Meta["total_data"])
	require.Equal(t, float64(1), resp.Meta["total_page"])

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_List_WhenCacheMiss_AndCountTimeout(t *testing.T) {

	// Prepare the request and expected response
	req := &request.Page{Size: 10, Before: ksuid.New().String(), After: ksuid.New().String()} // Ensure Before and After are initialized
	expectedUsers := []*entity.Users{
		{ID: "1", Username: "John Doe", Email: "john@example.com"},
	}
	expectedKey := fmt.Sprintf("users:all:size[%d]:before[%s]:after[%s]", req.Size, req.Before, req.After)

	// Define the behavior of the mocked methods
	cacheRepoMock.On("GetFromCache", mock.Anything, expectedKey).Return(nil, nil).Once()
	usersRepoMock.On("GetAll", mock.Anything, req).Return(expectedUsers, nil).Once()
	cacheRepoMock.On("SetToCache", mock.Anything, expectedKey, expectedUsers).Return(nil).Once()
	usersRepoMock.On("Count", mock.Anything).Return(int64(0), context.DeadlineExceeded).Once()

	// Call the List method
	resp, err := usersusecase.List(context.Background(), req)
	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_List_WhenRequestIsInvalid(t *testing.T) {
	// Prepare the request
	req := new(request.Page)

	//	Call the Method
	res, errors := usersusecase.List(context.Background(), req)
	require.Nil(t, res)
	require.Error(t, errors)
	require.Equal(t, errors.Errors[0].Status, http.StatusUnprocessableEntity)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(errors).String())
}

func TestUsersUsecase_List_WhenCacheTimeout(t *testing.T) {
	// Prepare the request
	req := &request.Page{Size: 10, Before: ksuid.New().String(), After: ksuid.New().String()} // Ensure Before and After are initialized

	// Defined Mock Method Call
	cacheRepoMock.On("GetFromCache", mock.Anything, fmt.Sprintf("users:all:size[%d]:before[%s]:after[%s]", req.Size, req.Before, req.After)).Return(nil, context.DeadlineExceeded).Once()

	//	Call the Method
	res, errors := usersusecase.List(context.Background(), req)
	require.Nil(t, res)
	require.Error(t, errors)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(errors).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

// =============================================== END LIST CASES =================================================================

// =============================================== CREATE CASES ===================================================================

func TestUsersUsecase_Create(t *testing.T) {

	// Prepare the request and expected response
	req := &request.User{Username: "john doe", Email: "john@example.com", Password: "password123"}
	// Define the behavior of the mocked methods
	usersRepoMock.On("ExistByKeyValue", mock.Anything, map[string]any{"email": req.Email}).Return(false, nil).Once()
	usersRepoMock.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
	cacheRepoMock.On("DeleteToCacheByRegexKey", mock.Anything, "users:*").Return(nil).Once()
	usersRepoMock.On("GetById", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	// Call the Create method
	resp, err := usersusecase.Create(context.Background(), req)

	// Assertions
	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusCreated, resp.Status)

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Create_WhenExist(t *testing.T) {

	// Prepare the request and expected response
	req := &request.User{Username: "john doe", Email: "john@example.com", Password: "password123"}
	// Define the behavior of the mocked methods
	usersRepoMock.On("ExistByKeyValue", mock.Anything, map[string]any{"email": req.Email}).Return(true, nil).Once()

	// Call the Create method
	resp, err := usersusecase.Create(context.Background(), req)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Create_WhenExistDBTimeout(t *testing.T) {

	// Prepare the request and expected response
	req := &request.User{Username: "john doe", Email: "john@example.com", Password: "password123"}
	// Define the behavior of the mocked methods
	usersRepoMock.On("ExistByKeyValue", mock.Anything, map[string]any{"email": req.Email}).Return(false, context.DeadlineExceeded).Once()

	// Call the Create method
	resp, err := usersusecase.Create(context.Background(), req)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Create_WhenSaveDBTimeout(t *testing.T) {

	// Prepare the request and expected response
	req := &request.User{Username: "john doe", Email: "john@example.com", Password: "password123"}
	// Define the behavior of the mocked methods
	usersRepoMock.On("ExistByKeyValue", mock.Anything, map[string]any{"email": req.Email}).Return(false, nil).Once()
	usersRepoMock.On("Create", mock.Anything, mock.Anything).Return(context.DeadlineExceeded).Once()

	// Call the Create method
	resp, err := usersusecase.Create(context.Background(), req)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Create_WhenDeleteCacheErr(t *testing.T) {

	// Prepare the request and expected response
	req := &request.User{Username: "john doe", Email: "john@example.com", Password: "password123"}
	// Define the behavior of the mocked methods
	usersRepoMock.On("ExistByKeyValue", mock.Anything, map[string]any{"email": req.Email}).Return(false, nil).Once()
	usersRepoMock.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
	cacheRepoMock.On("DeleteToCacheByRegexKey", mock.Anything, "users:*").Return(errors.New("internal server")).Once()

	// Call the Create method
	resp, err := usersusecase.Create(context.Background(), req)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Create_WhenGetByIdErr(t *testing.T) {

	// Prepare the request and expected response
	req := &request.User{Username: "john doe", Email: "john@example.com", Password: "password123"}
	// Define the behavior of the mocked methods
	usersRepoMock.On("ExistByKeyValue", mock.Anything, map[string]any{"email": req.Email}).Return(false, nil).Once()
	usersRepoMock.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
	cacheRepoMock.On("DeleteToCacheByRegexKey", mock.Anything, "users:*").Return(nil).Once()
	usersRepoMock.On("GetById", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("internal server")).Once()

	// Call the Create method
	resp, err := usersusecase.Create(context.Background(), req)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

// =============================================== END CREATE CASES ===============================================================

// =============================================== UPDATE CASES ===================================================================

func TestUsersUsecase_Update(t *testing.T) {

	// Prepare the request and expected response
	id := ksuid.New().String()
	req := &request.User{Username: "john doe", Email: "john@example.com", Password: "password123"}
	// Define the behavior of the mocked methods
	usersRepoMock.On("CountById", mock.Anything, id).Return(int64(1), nil).Once()
	usersRepoMock.On("Update", mock.Anything, mock.Anything, id).Return(nil).Once()
	cacheRepoMock.On("DeleteToCacheByRegexKey", mock.Anything, "users:*").Return(nil).Once()
	usersRepoMock.On("GetById", mock.Anything, mock.Anything, id).Return(nil).Once()

	// Call the Create method
	resp, err := usersusecase.Update(context.Background(), req, id)

	// Assertions
	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.Status)

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Update_WhenReqErr(t *testing.T) {

	// Prepare the request and expected response
	id := ksuid.New().String()
	req := &request.User{Username: "john doe", Email: "john@example.com"}
	// Define the behavior of the mocked methods

	// Call the Create method
	resp, err := usersusecase.Update(context.Background(), req, id)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Update_WhenInvalidID(t *testing.T) {

	// Prepare the request and expected response
	id := "1"
	req := &request.User{Username: "john doe", Email: "john@example.com", Password: "password123"}
	// Define the behavior of the mocked methods

	// Call the Create method
	resp, err := usersusecase.Update(context.Background(), req, id)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Update_WhenNotExist(t *testing.T) {

	// Prepare the request and expected response
	id := ksuid.New().String()
	req := &request.User{Username: "john doe", Email: "john@example.com", Password: "password123"}
	// Define the behavior of the mocked methods
	usersRepoMock.On("CountById", mock.Anything, id).Return(int64(0), nil).Once()

	// Call the Create method
	resp, err := usersusecase.Update(context.Background(), req, id)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Update_WhenUpdateErr(t *testing.T) {

	// Prepare the request and expected response
	id := ksuid.New().String()
	req := &request.User{Username: "john doe", Email: "john@example.com", Password: "password123"}
	// Define the behavior of the mocked methods
	usersRepoMock.On("CountById", mock.Anything, id).Return(int64(1), nil).Once()
	usersRepoMock.On("Update", mock.Anything, mock.Anything, id).Return(errors.New("internal server")).Once()

	// Call the Create method
	resp, err := usersusecase.Update(context.Background(), req, id)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Update_WhenDeleteCacheErr(t *testing.T) {

	// Prepare the request and expected response
	id := ksuid.New().String()
	req := &request.User{Username: "john doe", Email: "john@example.com", Password: "password123"}
	// Define the behavior of the mocked methods
	usersRepoMock.On("CountById", mock.Anything, id).Return(int64(1), nil).Once()
	usersRepoMock.On("Update", mock.Anything, mock.Anything, id).Return(nil).Once()
	cacheRepoMock.On("DeleteToCacheByRegexKey", mock.Anything, "users:*").Return(context.DeadlineExceeded).Once()

	// Call the Create method
	resp, err := usersusecase.Update(context.Background(), req, id)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Update_WhenGetErr(t *testing.T) {

	// Prepare the request and expected response
	id := ksuid.New().String()
	req := &request.User{Username: "john doe", Email: "john@example.com", Password: "password123"}
	// Define the behavior of the mocked methods
	usersRepoMock.On("CountById", mock.Anything, id).Return(int64(1), nil).Once()
	usersRepoMock.On("Update", mock.Anything, mock.Anything, id).Return(nil).Once()
	cacheRepoMock.On("DeleteToCacheByRegexKey", mock.Anything, "users:*").Return(nil).Once()
	usersRepoMock.On("GetById", mock.Anything, mock.Anything, id).Return(errors.New("internal server")).Once()

	// Call the Create method
	resp, err := usersusecase.Update(context.Background(), req, id)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

// =============================================== END UPDATE CASES ================================================================

// ================================================ GET CASES ======================================================================

func TestUsersUsecase_Get(t *testing.T) {
	// Initialize id and users to find users
	id := ksuid.New().String()
	expectedUsers := new(entity.Users)
	// Define the behavior of the mocked methods
	usersRepoMock.On("CountById", mock.Anything, id).Return(int64(1), nil).Once()
	usersRepoMock.On("GetById", mock.Anything, expectedUsers, id).Return(nil).Once()

	// Call the Get method
	resp, err := usersusecase.Get(context.Background(), id)
	// Assertions
	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.Status)
	require.Equal(t, "STATUS_OK", resp.Code)
	require.Equal(t, mapper.EntityUserToResponse(expectedUsers), resp.Data)

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Get_WhenInvalidID(t *testing.T) {
	// Initialize id and users to find users
	id := "1"
	// Define the behavior of the mocked methods

	// Call the Get method
	resp, err := usersusecase.Get(context.Background(), id)
	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Get_WhenNotExist(t *testing.T) {
	// Initialize id and users to find users
	id := ksuid.New().String()
	// Define the behavior of the mocked methods
	usersRepoMock.On("CountById", mock.Anything, id).Return(int64(0), nil).Once()

	// Call the Get method
	resp, err := usersusecase.Get(context.Background(), id)
	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Get_WhenGetErr(t *testing.T) {
	// Initialize id and users to find users
	id := ksuid.New().String()
	expectedUsers := new(entity.Users)
	// Define the behavior of the mocked methods
	usersRepoMock.On("CountById", mock.Anything, id).Return(int64(1), nil).Once()
	usersRepoMock.On("GetById", mock.Anything, expectedUsers, id).Return(context.DeadlineExceeded).Once()

	// Call the Get method
	resp, err := usersusecase.Get(context.Background(), id)
	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
}

// ================================================ END GET CASES ===================================================================

// ================================================ EDIT CASES =====================================================================

func TestUsersUsecase_Edit(t *testing.T) {

	// Prepare the request and expected response
	id := ksuid.New().String()
	req := &request.UserEdit{Username: "john doe", Email: "john@example.com", Password: "password123"}
	// Define the behavior of the mocked methods
	usersRepoMock.On("CountById", mock.Anything, id).Return(int64(1), nil).Once()
	usersRepoMock.On("Update", mock.Anything, mock.Anything, id).Return(nil).Once()
	cacheRepoMock.On("DeleteToCacheByRegexKey", mock.Anything, "users:*").Return(nil).Once()
	usersRepoMock.On("GetById", mock.Anything, mock.Anything, id).Return(nil).Once()

	// Call the Create method
	resp, err := usersusecase.Edit(context.Background(), req, id)

	// Assertions
	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.Status)

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Edit_WhenInvalidID(t *testing.T) {

	// Prepare the request and expected response
	id := "1"
	req := &request.UserEdit{Username: "john doe", Email: "john@example.com", Password: "password123"}
	// Define the behavior of the mocked methods

	// Call the Create method
	resp, err := usersusecase.Edit(context.Background(), req, id)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Edit_WhenNotExist(t *testing.T) {

	// Prepare the request and expected response
	id := ksuid.New().String()
	req := &request.UserEdit{Username: "john doe", Email: "john@example.com", Password: "password123"}
	// Define the behavior of the mocked methods
	usersRepoMock.On("CountById", mock.Anything, id).Return(int64(0), nil).Once()

	// Call the Create method
	resp, err := usersusecase.Edit(context.Background(), req, id)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Edit_WhenUpdateErr(t *testing.T) {

	// Prepare the request and expected response
	id := ksuid.New().String()
	req := &request.UserEdit{Username: "john doe", Email: "john@example.com", Password: "password123"}
	// Define the behavior of the mocked methods
	usersRepoMock.On("CountById", mock.Anything, id).Return(int64(1), nil).Once()
	usersRepoMock.On("Update", mock.Anything, mock.Anything, id).Return(errors.New("internal server")).Once()

	// Call the Create method
	resp, err := usersusecase.Edit(context.Background(), req, id)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Edit_WhenDeleteCacheErr(t *testing.T) {

	// Prepare the request and expected response
	id := ksuid.New().String()
	req := &request.UserEdit{Username: "john doe", Email: "john@example.com", Password: "password123"}
	// Define the behavior of the mocked methods
	usersRepoMock.On("CountById", mock.Anything, id).Return(int64(1), nil).Once()
	usersRepoMock.On("Update", mock.Anything, mock.Anything, id).Return(nil).Once()
	cacheRepoMock.On("DeleteToCacheByRegexKey", mock.Anything, "users:*").Return(context.DeadlineExceeded).Once()

	// Call the Create method
	resp, err := usersusecase.Edit(context.Background(), req, id)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Edit_WhenGetErr(t *testing.T) {

	// Prepare the request and expected response
	id := ksuid.New().String()
	req := &request.UserEdit{Username: "john doe", Email: "john@example.com", Password: "password123"}
	// Define the behavior of the mocked methods
	usersRepoMock.On("CountById", mock.Anything, id).Return(int64(1), nil).Once()
	usersRepoMock.On("Update", mock.Anything, mock.Anything, id).Return(nil).Once()
	cacheRepoMock.On("DeleteToCacheByRegexKey", mock.Anything, "users:*").Return(nil).Once()
	usersRepoMock.On("GetById", mock.Anything, mock.Anything, id).Return(errors.New("internal server")).Once()

	// Call the Create method
	resp, err := usersusecase.Edit(context.Background(), req, id)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

// ================================================ END EDIT CASES =================================================================

// ================================================ DELETE CASES ===================================================================

func TestUsersUsecase_Delete(t *testing.T) {

	// Prepare the request and expected response
	id := ksuid.New().String()
	// Define the behavior of the mocked methods
	usersRepoMock.On("CountById", mock.Anything, id).Return(int64(1), nil).Once()
	usersRepoMock.On("Delete", mock.Anything, id).Return(nil).Once()
	cacheRepoMock.On("DeleteToCacheByRegexKey", mock.Anything, "users:*").Return(nil).Once()

	// Call the Create method
	resp, err := usersusecase.Delete(context.Background(), id)

	// Assertions
	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.Status)

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Delete_WhenInvalidID(t *testing.T) {

	// Prepare the request and expected response
	id := "1"
	// Define the behavior of the mocked methods

	// Call the Create method
	resp, err := usersusecase.Delete(context.Background(), id)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Delete_WhenNotExist(t *testing.T) {

	// Prepare the request and expected response
	id := ksuid.New().String()
	// Define the behavior of the mocked methods
	usersRepoMock.On("CountById", mock.Anything, id).Return(int64(0), nil).Once()

	// Call the Create method
	resp, err := usersusecase.Delete(context.Background(), id)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Delete_WhenDeleteErr(t *testing.T) {

	// Prepare the request and expected response
	id := ksuid.New().String()
	// Define the behavior of the mocked methods
	usersRepoMock.On("CountById", mock.Anything, id).Return(int64(1), nil).Once()
	usersRepoMock.On("Delete", mock.Anything, id).Return(errors.New("internal server")).Once()

	// Call the Create method
	resp, err := usersusecase.Delete(context.Background(), id)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Delete_WhenDeleteCacheErr(t *testing.T) {

	// Prepare the request and expected response
	id := ksuid.New().String()
	// Define the behavior of the mocked methods
	usersRepoMock.On("CountById", mock.Anything, id).Return(int64(1), nil).Once()
	usersRepoMock.On("Delete", mock.Anything, id).Return(nil).Once()
	cacheRepoMock.On("DeleteToCacheByRegexKey", mock.Anything, "users:*").Return(context.DeadlineExceeded).Once()

	// Call the Create method
	resp, err := usersusecase.Delete(context.Background(), id)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

// ================================================== END DELETE CASES =============================================================

// ==================================================== RESTORE CASES ==============================================================
func TestUsersUsecase_Restore(t *testing.T) {

	// Prepare the request and expected response
	id := ksuid.New().String()
	// Define the behavior of the mocked methods
	usersRepoMock.On("CountById", mock.Anything, id).Return(int64(0), nil).Once()
	usersRepoMock.On("Restore", mock.Anything, id).Return(nil).Once()
	cacheRepoMock.On("DeleteToCacheByRegexKey", mock.Anything, "users:*").Return(nil).Once()

	// Call the Create method
	resp, err := usersusecase.Restore(context.Background(), id)

	// Assertions
	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.Status)

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Restore_WhenInvalidID(t *testing.T) {

	// Prepare the request and expected response
	id := "1"
	// Define the behavior of the mocked methods

	// Call the Create method
	resp, err := usersusecase.Restore(context.Background(), id)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Restore_WhenExist(t *testing.T) {

	// Prepare the request and expected response
	id := ksuid.New().String()
	// Define the behavior of the mocked methods
	usersRepoMock.On("CountById", mock.Anything, id).Return(int64(1), nil).Once()

	// Call the Create method
	resp, err := usersusecase.Restore(context.Background(), id)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Restore_WhenDeleteErr(t *testing.T) {

	// Prepare the request and expected response
	id := ksuid.New().String()
	// Define the behavior of the mocked methods
	usersRepoMock.On("CountById", mock.Anything, id).Return(int64(0), nil).Once()
	usersRepoMock.On("Restore", mock.Anything, id).Return(errors.New("internal server")).Once()

	// Call the Create method
	resp, err := usersusecase.Restore(context.Background(), id)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

func TestUsersUsecase_Restore_WhenDeleteCacheErr(t *testing.T) {

	// Prepare the request and expected response
	id := ksuid.New().String()
	// Define the behavior of the mocked methods
	usersRepoMock.On("CountById", mock.Anything, id).Return(int64(0), nil).Once()
	usersRepoMock.On("Restore", mock.Anything, id).Return(nil).Once()
	cacheRepoMock.On("DeleteToCacheByRegexKey", mock.Anything, "users:*").Return(context.DeadlineExceeded).Once()

	// Call the Create method
	resp, err := usersusecase.Restore(context.Background(), id)

	// Assertions
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
	cacheRepoMock.AssertExpectations(t)
}

// ==================================================== END RESTORE CASES ==============================================================

// ===================================================== LOGIN CASES ===================================================================

func TestAuthUsecase_Login_WhenInvalidReq(t *testing.T) {
	// Prepare Request and mock arguments
	req := request.Auth{Email: "", Password: ""}
	// Define the behavior of the mocked methods
	// Call the Login methods
	resp, token, err := authusecase.Login(context.Background(), &req)
	// Assertions
	require.Empty(t, token)
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
}

func TestAuthUsecase_Login_WhenGetErr(t *testing.T) {
	// Prepare Request and mock arguments
	req := request.Auth{Email: "john@example.com", Password: "password123"}
	// Define the behavior of the mocked methods
	usersRepoMock.On("ExistByKeyValue", mock.Anything, mock.Anything, mock.Anything).Return(true, nil).Once()
	usersRepoMock.On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	// Call the Login methods
	resp, token, err := authusecase.Login(context.Background(), &req)
	// Assertions
	require.Empty(t, token)
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
}

func TestAuthUsecase_Login_WhenNotFound(t *testing.T) {
	// Prepare Request and mock arguments
	req := request.Auth{Email: "john@example.com", Password: "password123"}
	// Define the behavior of the mocked methods
	usersRepoMock.On("ExistByKeyValue", mock.Anything, mock.Anything, mock.Anything).Return(false, nil).Once()
	// Call the Login methods
	resp, token, err := authusecase.Login(context.Background(), &req)
	// Assertions
	require.Empty(t, token)
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(err).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
}

// ===================================================== END LOGIN CASES ===============================================================

// ===================================================== RESET PASSWORD CASES ==========================================================

func TestAuthUsecase_ResetPassword(t *testing.T) {
	// Prepare Request and mock arguments
	req := request.ResetPassword{Password: "password123", Confirm: "password123"}
	payload := token.NewTokenPayloadBuilder().WithEmail("john@example.com").WithUserID(ksuid.New()).WithExpiration(time.Now().Add(5 * time.Minute)).Build()
	// Define the behavior of the mocked methods
	usersRepoMock.On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	usersRepoMock.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	// Call the Login methods
	resp, errResp := authusecase.ResetPassword(context.Background(), payload, &req)
	// Assertions
	require.NotNil(t, resp)
	require.Nil(t, errResp)
	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
}

func TestAuthUsecase_ResetPassword_WhenInvalidReq(t *testing.T) {
	// Prepare Request and mock arguments
	req := request.ResetPassword{Password: "password12", Confirm: ""}
	payload := token.NewTokenPayloadBuilder().WithEmail("john@example.com").WithUserID(ksuid.New()).WithExpiration(time.Now().Add(5 * time.Minute)).Build()
	// Define the behavior of the mocked methods
	// Call the Login methods
	resp, errResp := authusecase.ResetPassword(context.Background(), payload, &req)
	// Assertions
	require.Nil(t, resp)
	require.Error(t, errResp)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(errResp).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
}

func TestAuthUsecase_ResetPassword_WhenGetErr(t *testing.T) {
	// Prepare Request and mock arguments
	req := request.ResetPassword{Password: "password123", Confirm: "password123"}
	payload := token.NewTokenPayloadBuilder().WithEmail("john@example.com").WithUserID(ksuid.New()).WithExpiration(time.Now().Add(5 * time.Minute)).Build()
	// Define the behavior of the mocked methods
	usersRepoMock.On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("internal server")).Once()
	// Call the Login methods
	resp, errResp := authusecase.ResetPassword(context.Background(), payload, &req)
	// Assertions
	require.Nil(t, resp)
	require.Error(t, errResp)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(errResp).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
}

func TestAuthUsecase_ResetPassword_WhenUpdateErr(t *testing.T) {
	// Prepare Request and mock arguments
	req := request.ResetPassword{Password: "password123", Confirm: "password123"}
	payload := token.NewTokenPayloadBuilder().WithEmail("john@example.com").WithUserID(ksuid.New()).WithExpiration(time.Now().Add(5 * time.Minute)).Build()
	// Define the behavior of the mocked methods
	usersRepoMock.On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	usersRepoMock.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("internal server")).Once()

	// Call the Login methods
	resp, errResp := authusecase.ResetPassword(context.Background(), payload, &req)
	// Assertions
	require.Nil(t, resp)
	require.Error(t, errResp)
	require.Equal(t, reflect.TypeOf(new(response.StandardErrors)).String(), reflect.TypeOf(errResp).String())

	// Assert that all expectations were met
	usersRepoMock.AssertExpectations(t)
}

// ===================================================== END RESET PASSWORD CASES =======================================================
