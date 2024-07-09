package bootstrap_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/bootstrap"
	pathhelper "github.com/tirtahakimpambudhi/restful_api/pkg/helper/path"
	"os"
	"testing"
	"time"
)

func SetUpEnv(keyvalue map[string]string) {
	for key, value := range keyvalue {
		os.Setenv(key, value)
	}
}

func SetupContainer(t *testing.T) []func(ctx context.Context) error {
	if testing.Short() {
		t.Skip()
		return nil
	}
	ctx := context.Background()

	// Redis setup
	redisReq := testcontainers.ContainerRequest{
		Image:       "redis:latest",
		WaitingFor:  wait.ForLog("Ready to accept connections"),
		NetworkMode: "host",
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: redisReq,
		Started:          true,
	})
	require.NoError(t, err)
	t.Log("Successfully Redis Setup")

	// Get the Redis port and host
	redisPort, err := redisC.MappedPort(ctx, "6379")
	require.NoError(t, err)
	t.Logf("Redis is running on port %s", redisPort.Port())

	// Postgres setup
	postgresReq := testcontainers.ContainerRequest{
		Image: "postgres:13-alpine",
		Env: map[string]string{
			"POSTGRES_DB":       "test",
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
		},
		WaitingFor:  wait.ForLog("database system is ready to accept connections"),
		NetworkMode: "host",
	}
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: postgresReq,
		Started:          true,
	})
	time.Sleep(10 * time.Second)
	require.NoError(t, err)
	t.Log("Successfully Postgres Setup")

	// Get the PostgreSQL port and connection string
	postgresPort, err := postgresC.MappedPort(ctx, "5432")
	require.NoError(t, err)
	t.Logf("Postgres is running on port %s", postgresPort.Port())

	// Print the connection string for PostgreSQL
	postgresHost, err := postgresC.Host(ctx)
	require.NoError(t, err)
	connectionString := fmt.Sprintf("postgres://postgres:postgres@%s:%s/test?sslmode=disable", postgresHost, postgresPort.Port())
	t.Logf("Postgres connection string: %s", connectionString)

	return []func(context.Context) error{redisC.Terminate, postgresC.Terminate}
}

func SetupCasbin(t *testing.T) {
	// Generate the file path
	filePath := pathhelper.AddWorkdirToSomePath("resource", "model", "rbac_model.conf")

	// Ensure the directory exists
	dirPath := pathhelper.AddWorkdirToSomePath("resource", "model")
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		require.NoError(t, err, "Failed to create directory: %v", err)
		return
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		require.NoError(t, err, "Failed to create file: %v", err)
		return
	}
	defer file.Close() // Defer closing only after checking for file creation error

	// Write the content to the file
	_, err = file.Write([]byte(`
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`))
	if err != nil {
		require.NoError(t, err, "Failed to write to file: %v", err)
	}
}

func ClearEnv(keyvalue map[string]string) {
	for key, _ := range keyvalue {
		os.Unsetenv(key)
	}
}

func TestMain(t *testing.M) {
	t.Run()

	defer func() {
		os.RemoveAll(pathhelper.AddWorkdirToSomePath("resource"))
	}()
}

func TestInitBootstrap(t *testing.T) {
	terminates := SetupContainer(t)
	SetupCasbin(t)

	defer func() {
		if terminates != nil || len(terminates) <= 0 {
			for _, terminate := range terminates {
				err := terminate(context.Background())
				require.NoError(t, err)
			}
		}
	}()
	testCases := []struct {
		name    string
		env     map[string]string
		isError bool
	}{
		{name: "Successfully Init Bootstrap", env: map[string]string{
			// RedisConfig
			"CACHE_DB_NAME":     "1",
			"CACHE_DB_HOST":     "localhost",
			"CACHE_DB_PORT":     "6379",
			"CACHE_DB_USER":     "default",
			"CACHE_DB_PASS":     "fakepassword",
			"CACHE_DB_MAX_CON":  "100",
			"CACHE_DB_MIN_CON":  "10",
			"CACHE_DB_MAX_TIME": "10",
			"CACHE_DB_MIN_TIME": "2",

			// Casbin
			"MODEL_PATH": "resource/model",

			// FiberConfig
			"FIBER_HOST":           "localhost",
			"FIBER_PORT":           "3000",
			"FIBER_PREFORK":        "false",
			"FIBER_STRICT_ROUTING": "true",
			"FIBER_CASE_SENSITIVE": "true",
			"FIBER_BODY_LIMIT":     "4",
			"FIBER_READ_TIMEOUT":   "4",
			"FIBER_WRITE_TIMEOUT":  "5",
			"FIBER_REDUCE_MEMU":    "true",
			"FIBER_JSON":           "json",

			// Bcrypt
			"HASH_SALT": "10",

			// LoggerConfig
			"LOG_PATH":             "resource/logs",
			"LOG_MAX_SIZE":         "10",
			"LOG_MAX_BACKUP":       "5",
			"LOG_MAX_SIZE_ROTATE":  "20",
			"LOG_TIME_FORMAT":      "2006-01-02",
			"LOG_COLOR_OUTPUT":     "false",
			"LOG_QUOTE_STR":        "false",
			"LOG_END_WITH_MESSAGE": "false",

			// SqlConfig
			"DB_DRIVER":   "postgres",
			"DB_PROTOCOL": "postgresql",
			"DB_NAME":     "test",
			"DB_HOST":     "localhost",
			"DB_PORT":     "5432",
			"DB_USER":     "postgres",
			"DB_PASS":     "postgres",
			"DB_MAX_CON":  "100",
			"DB_MIN_CON":  "10",
			"DB_MAX_TIME": "30",
			"DB_MIN_TIME": "5",

			// JWTToken
			"TOKEN_NAME": "RESTful_API_AUTH",

			// SecretKey
			"SECRET_KEY_ACCESS_TOKEN":  "Zawssh9t1IY50IlICrYpjrCDbq6G8UKL",
			"SECRET_KEY_REFRESH_TOKEN": "ZxywJoMMXIXwgcispeKs4L6Y65XgATqV",
			"SECRET_KEY_FP_TOKEN":      "BNXWuiMew8HhFHLirNw1zpOtO0aJW1cE",

			// Timeout
			"CACHE_TIMEOUT":       "8",
			"DB_TIMEOUT":          "20",
			"DOWN_STREAM_TIMEOUT": "30",
		}, isError: false},
		{name: "Failure Init Bootstrap Because FiberConfig Error", env: map[string]string{
			// RedisConfig
			"CACHE_DB_NAME":     "1",
			"CACHE_DB_HOST":     "localhost",
			"CACHE_DB_PORT":     "6379",
			"CACHE_DB_USER":     "default",
			"CACHE_DB_PASS":     "fakepassword",
			"CACHE_DB_MAX_CON":  "100",
			"CACHE_DB_MIN_CON":  "10",
			"CACHE_DB_MAX_TIME": "10",
			"CACHE_DB_MIN_TIME": "2",

			// Casbin
			"MODEL_PATH": "resource/model",
			// Bcrypt
			"HASH_SALT": "10",

			// LoggerConfig
			"LOG_PATH":             "resource/logs",
			"LOG_MAX_SIZE":         "10",
			"LOG_MAX_BACKUP":       "5",
			"LOG_MAX_SIZE_ROTATE":  "20",
			"LOG_TIME_FORMAT":      "2006-01-02",
			"LOG_COLOR_OUTPUT":     "false",
			"LOG_QUOTE_STR":        "false",
			"LOG_END_WITH_MESSAGE": "false",

			// SqlConfig
			"DB_DRIVER":   "postgres",
			"DB_PROTOCOL": "postgresql",
			"DB_NAME":     "test",
			"DB_HOST":     "localhost",
			"DB_PORT":     "5432",
			"DB_USER":     "postgres",
			"DB_PASS":     "fakepassword",
			"DB_MAX_CON":  "100",
			"DB_MIN_CON":  "10",
			"DB_MAX_TIME": "30",
			"DB_MIN_TIME": "5",

			// JWTToken
			"TOKEN_NAME": "RESTful_API_AUTH",

			// SecretKey
			"SECRET_KEY_ACCESS_TOKEN":  "Zawssh9t1IY50IlICrYpjrCDbq6G8UKL",
			"SECRET_KEY_REFRESH_TOKEN": "ZxywJoMMXIXwgcispeKs4L6Y65XgATqV",
			"SECRET_KEY_FP_TOKEN":      "BNXWuiMew8HhFHLirNw1zpOtO0aJW1cE",

			// Timeout
			"CACHE_TIMEOUT":       "8",
			"DB_TIMEOUT":          "20",
			"DOWN_STREAM_TIMEOUT": "30",
		}, isError: true},
		{name: "Failure Init Bootstrap Because RedisConfig Error", env: map[string]string{

			// Casbin
			"MODEL_PATH": "resource/model",

			// FiberConfig
			"FIBER_HOST":           "localhost",
			"FIBER_PORT":           "3000",
			"FIBER_PREFORK":        "false",
			"FIBER_STRICT_ROUTING": "true",
			"FIBER_CASE_SENSITIVE": "true",
			"FIBER_BODY_LIMIT":     "4",
			"FIBER_READ_TIMEOUT":   "4",
			"FIBER_WRITE_TIMEOUT":  "5",
			"FIBER_REDUCE_MEMU":    "true",
			"FIBER_JSON":           "json",

			// Bcrypt
			"HASH_SALT": "10",

			// LoggerConfig
			"LOG_PATH":             "resource/logs",
			"LOG_MAX_SIZE":         "10",
			"LOG_MAX_BACKUP":       "5",
			"LOG_MAX_SIZE_ROTATE":  "20",
			"LOG_TIME_FORMAT":      "2006-01-02",
			"LOG_COLOR_OUTPUT":     "false",
			"LOG_QUOTE_STR":        "false",
			"LOG_END_WITH_MESSAGE": "false",

			// SqlConfig
			"DB_DRIVER":   "postgres",
			"DB_PROTOCOL": "postgresql",
			"DB_NAME":     "test",
			"DB_HOST":     "localhost",
			"DB_PORT":     "5432",
			"DB_USER":     "postgres",
			"DB_PASS":     "postgres",
			"DB_MAX_CON":  "100",
			"DB_MIN_CON":  "10",
			"DB_MAX_TIME": "30",
			"DB_MIN_TIME": "5",

			// JWTToken
			"TOKEN_NAME": "RESTful_API_AUTH",

			// SecretKey
			"SECRET_KEY_ACCESS_TOKEN":  "Zawssh9t1IY50IlICrYpjrCDbq6G8UKL",
			"SECRET_KEY_REFRESH_TOKEN": "ZxywJoMMXIXwgcispeKs4L6Y65XgATqV",
			"SECRET_KEY_FP_TOKEN":      "BNXWuiMew8HhFHLirNw1zpOtO0aJW1cE",

			// Timeout
			"CACHE_TIMEOUT":       "8",
			"DB_TIMEOUT":          "20",
			"DOWN_STREAM_TIMEOUT": "30",
		}, isError: true},
		{name: "Failure Init Bootstrap Because LoggerConfig Error", env: map[string]string{
			// RedisConfig
			"CACHE_DB_NAME":     "1",
			"CACHE_DB_HOST":     "localhost",
			"CACHE_DB_PORT":     "6379",
			"CACHE_DB_USER":     "default",
			"CACHE_DB_PASS":     "fakepassword",
			"CACHE_DB_MAX_CON":  "100",
			"CACHE_DB_MIN_CON":  "10",
			"CACHE_DB_MAX_TIME": "10",
			"CACHE_DB_MIN_TIME": "2",

			// Casbin
			"MODEL_PATH": "resource/model",

			// FiberConfig
			"FIBER_HOST":           "localhost",
			"FIBER_PORT":           "3000",
			"FIBER_PREFORK":        "false",
			"FIBER_STRICT_ROUTING": "true",
			"FIBER_CASE_SENSITIVE": "true",
			"FIBER_BODY_LIMIT":     "4",
			"FIBER_READ_TIMEOUT":   "4",
			"FIBER_WRITE_TIMEOUT":  "5",
			"FIBER_REDUCE_MEMU":    "true",
			"FIBER_JSON":           "json",

			// Bcrypt
			"HASH_SALT": "10",

			// SqlConfig
			"DB_DRIVER":   "postgres",
			"DB_PROTOCOL": "postgresql",
			"DB_NAME":     "test",
			"DB_HOST":     "localhost",
			"DB_PORT":     "5432",
			"DB_USER":     "postgres",
			"DB_PASS":     "postgres",
			"DB_MAX_CON":  "100",
			"DB_MIN_CON":  "10",
			"DB_MAX_TIME": "30",
			"DB_MIN_TIME": "5",

			// JWTToken
			"TOKEN_NAME": "RESTful_API_AUTH",

			// SecretKey
			"SECRET_KEY_ACCESS_TOKEN":  "Zawssh9t1IY50IlICrYpjrCDbq6G8UKL",
			"SECRET_KEY_REFRESH_TOKEN": "ZxywJoMMXIXwgcispeKs4L6Y65XgATqV",
			"SECRET_KEY_FP_TOKEN":      "BNXWuiMew8HhFHLirNw1zpOtO0aJW1cE",

			// Timeout
			"CACHE_TIMEOUT":       "8",
			"DB_TIMEOUT":          "20",
			"DOWN_STREAM_TIMEOUT": "30",
		}, isError: true},
		{name: "Failure Init Bootstrap Because SQLConfig Error", env: map[string]string{
			// RedisConfig
			"CACHE_DB_NAME":     "1",
			"CACHE_DB_HOST":     "localhost",
			"CACHE_DB_PORT":     "6379",
			"CACHE_DB_USER":     "default",
			"CACHE_DB_PASS":     "fakepassword",
			"CACHE_DB_MAX_CON":  "100",
			"CACHE_DB_MIN_CON":  "10",
			"CACHE_DB_MAX_TIME": "10",
			"CACHE_DB_MIN_TIME": "2",

			// Casbin
			"MODEL_PATH": "resource/model",

			// FiberConfig
			"FIBER_HOST":           "localhost",
			"FIBER_PORT":           "3000",
			"FIBER_PREFORK":        "false",
			"FIBER_STRICT_ROUTING": "true",
			"FIBER_CASE_SENSITIVE": "true",
			"FIBER_BODY_LIMIT":     "4",
			"FIBER_READ_TIMEOUT":   "4",
			"FIBER_WRITE_TIMEOUT":  "5",
			"FIBER_REDUCE_MEMU":    "true",
			"FIBER_JSON":           "json",

			// Bcrypt
			"HASH_SALT": "10",

			// LoggerConfig
			"LOG_PATH":             "resource/logs",
			"LOG_MAX_SIZE":         "10",
			"LOG_MAX_BACKUP":       "5",
			"LOG_MAX_SIZE_ROTATE":  "20",
			"LOG_TIME_FORMAT":      "2006-01-02",
			"LOG_COLOR_OUTPUT":     "false",
			"LOG_QUOTE_STR":        "false",
			"LOG_END_WITH_MESSAGE": "false",

			// JWTToken
			"TOKEN_NAME": "RESTful_API_AUTH",

			// SecretKey
			"SECRET_KEY_ACCESS_TOKEN":  "Zawssh9t1IY50IlICrYpjrCDbq6G8UKL",
			"SECRET_KEY_REFRESH_TOKEN": "ZxywJoMMXIXwgcispeKs4L6Y65XgATqV",
			"SECRET_KEY_FP_TOKEN":      "BNXWuiMew8HhFHLirNw1zpOtO0aJW1cE",

			// Timeout
			"CACHE_TIMEOUT":       "8",
			"DB_TIMEOUT":          "20",
			"DOWN_STREAM_TIMEOUT": "30",
		}, isError: true},
		{name: "Failure Init Bootstrap Because SecretConfig Error", env: map[string]string{
			// RedisConfig
			"CACHE_DB_NAME":     "1",
			"CACHE_DB_HOST":     "localhost",
			"CACHE_DB_PORT":     "6379",
			"CACHE_DB_USER":     "default",
			"CACHE_DB_PASS":     "fakepassword",
			"CACHE_DB_MAX_CON":  "100",
			"CACHE_DB_MIN_CON":  "10",
			"CACHE_DB_MAX_TIME": "10",
			"CACHE_DB_MIN_TIME": "2",

			// Casbin
			"MODEL_PATH": "resource/model",

			// FiberConfig
			"FIBER_HOST":           "localhost",
			"FIBER_PORT":           "3000",
			"FIBER_PREFORK":        "false",
			"FIBER_STRICT_ROUTING": "true",
			"FIBER_CASE_SENSITIVE": "true",
			"FIBER_BODY_LIMIT":     "4",
			"FIBER_READ_TIMEOUT":   "4",
			"FIBER_WRITE_TIMEOUT":  "5",
			"FIBER_REDUCE_MEMU":    "true",
			"FIBER_JSON":           "json",

			// Bcrypt
			"HASH_SALT": "10",

			// LoggerConfig
			"LOG_PATH":             "resource/logs",
			"LOG_MAX_SIZE":         "10",
			"LOG_MAX_BACKUP":       "5",
			"LOG_MAX_SIZE_ROTATE":  "20",
			"LOG_TIME_FORMAT":      "2006-01-02",
			"LOG_COLOR_OUTPUT":     "false",
			"LOG_QUOTE_STR":        "false",
			"LOG_END_WITH_MESSAGE": "false",

			// SqlConfig
			"DB_DRIVER":   "postgres",
			"DB_PROTOCOL": "postgresql",
			"DB_NAME":     "test",
			"DB_HOST":     "localhost",
			"DB_PORT":     "5432",
			"DB_USER":     "postgres",
			"DB_PASS":     "postgres",
			"DB_MAX_CON":  "100",
			"DB_MIN_CON":  "10",
			"DB_MAX_TIME": "30",
			"DB_MIN_TIME": "5",

			// JWTToken
			"TOKEN_NAME": "RESTful_API_AUTH",

			// Timeout
			"CACHE_TIMEOUT":       "8",
			"DB_TIMEOUT":          "20",
			"DOWN_STREAM_TIMEOUT": "30",
		}, isError: true},
		{name: "Failure Init Bootstrap Because GORM Wrong Connection DB Source", env: map[string]string{
			// RedisConfig
			"CACHE_DB_NAME":     "1",
			"CACHE_DB_HOST":     "localhost",
			"CACHE_DB_PORT":     "6379",
			"CACHE_DB_USER":     "default",
			"CACHE_DB_PASS":     "fakepassword",
			"CACHE_DB_MAX_CON":  "100",
			"CACHE_DB_MIN_CON":  "10",
			"CACHE_DB_MAX_TIME": "10",
			"CACHE_DB_MIN_TIME": "2",

			// Casbin
			"MODEL_PATH": "resource/model",

			// FiberConfig
			"FIBER_HOST":           "localhost",
			"FIBER_PORT":           "3000",
			"FIBER_PREFORK":        "false",
			"FIBER_STRICT_ROUTING": "true",
			"FIBER_CASE_SENSITIVE": "true",
			"FIBER_BODY_LIMIT":     "4",
			"FIBER_READ_TIMEOUT":   "4",
			"FIBER_WRITE_TIMEOUT":  "5",
			"FIBER_REDUCE_MEMU":    "true",
			"FIBER_JSON":           "json",

			// Bcrypt
			"HASH_SALT": "10",

			// LoggerConfig
			"LOG_PATH":             "resource/logs",
			"LOG_MAX_SIZE":         "10",
			"LOG_MAX_BACKUP":       "5",
			"LOG_MAX_SIZE_ROTATE":  "20",
			"LOG_TIME_FORMAT":      "2006-01-02",
			"LOG_COLOR_OUTPUT":     "false",
			"LOG_QUOTE_STR":        "false",
			"LOG_END_WITH_MESSAGE": "false",

			// SqlConfig
			"DB_DRIVER":   "postgres",
			"DB_PROTOCOL": "postgresql",
			"DB_NAME":     "test",
			"DB_HOST":     "localhost",
			"DB_PORT":     "5432",
			"DB_USER":     "postgresql",
			"DB_PASS":     "wrongpassword",
			"DB_MAX_CON":  "100",
			"DB_MIN_CON":  "10",
			"DB_MAX_TIME": "30",
			"DB_MIN_TIME": "5",

			// JWTToken
			"TOKEN_NAME": "RESTful_API_AUTH",

			// SecretKey
			"SECRET_KEY_ACCESS_TOKEN":  "Zawssh9t1IY50IlICrYpjrCDbq6G8UKL",
			"SECRET_KEY_REFRESH_TOKEN": "ZxywJoMMXIXwgcispeKs4L6Y65XgATqV",
			"SECRET_KEY_FP_TOKEN":      "BNXWuiMew8HhFHLirNw1zpOtO0aJW1cE",

			// Timeout
			"CACHE_TIMEOUT":       "8",
			"DB_TIMEOUT":          "20",
			"DOWN_STREAM_TIMEOUT": "30",
		}, isError: true},
		{name: "Failure Init Bootstrap Because JWTConfig Error", env: map[string]string{
			// RedisConfig
			"CACHE_DB_NAME":     "1",
			"CACHE_DB_HOST":     "localhost",
			"CACHE_DB_PORT":     "6379",
			"CACHE_DB_USER":     "default",
			"CACHE_DB_PASS":     "fakepassword",
			"CACHE_DB_MAX_CON":  "100",
			"CACHE_DB_MIN_CON":  "10",
			"CACHE_DB_MAX_TIME": "10",
			"CACHE_DB_MIN_TIME": "2",

			// Casbin
			"MODEL_PATH": "resource/model",

			// FiberConfig
			"FIBER_HOST":           "localhost",
			"FIBER_PORT":           "3000",
			"FIBER_PREFORK":        "false",
			"FIBER_STRICT_ROUTING": "true",
			"FIBER_CASE_SENSITIVE": "true",
			"FIBER_BODY_LIMIT":     "4",
			"FIBER_READ_TIMEOUT":   "4",
			"FIBER_WRITE_TIMEOUT":  "5",
			"FIBER_REDUCE_MEMU":    "true",
			"FIBER_JSON":           "json",

			// Bcrypt
			"HASH_SALT": "10",

			// LoggerConfig
			"LOG_PATH":             "resource/logs",
			"LOG_MAX_SIZE":         "10",
			"LOG_MAX_BACKUP":       "5",
			"LOG_MAX_SIZE_ROTATE":  "20",
			"LOG_TIME_FORMAT":      "2006-01-02",
			"LOG_COLOR_OUTPUT":     "false",
			"LOG_QUOTE_STR":        "false",
			"LOG_END_WITH_MESSAGE": "false",

			// SqlConfig
			"DB_DRIVER":   "postgres",
			"DB_PROTOCOL": "postgresql",
			"DB_NAME":     "test",
			"DB_HOST":     "localhost",
			"DB_PORT":     "5432",
			"DB_USER":     "postgres",
			"DB_PASS":     "postgres",
			"DB_MAX_CON":  "100",
			"DB_MIN_CON":  "10",
			"DB_MAX_TIME": "30",
			"DB_MIN_TIME": "5",

			// SecretKey
			"SECRET_KEY_ACCESS_TOKEN":  "Zawssh9t1IY50IlICrYpjrCDbq6G8UKL",
			"SECRET_KEY_REFRESH_TOKEN": "ZxywJoMMXIXwgcispeKs4L6Y65XgATqV",
			"SECRET_KEY_FP_TOKEN":      "BNXWuiMew8HhFHLirNw1zpOtO0aJW1cE",

			// Timeout
			"CACHE_TIMEOUT":       "8",
			"DB_TIMEOUT":          "20",
			"DOWN_STREAM_TIMEOUT": "30",
		}, isError: true},
		{name: "Failure Init Bootstrap Because TimeOutConfig Error", env: map[string]string{
			// RedisConfig
			"CACHE_DB_NAME":     "1",
			"CACHE_DB_HOST":     "localhost",
			"CACHE_DB_PORT":     "6379",
			"CACHE_DB_USER":     "default",
			"CACHE_DB_PASS":     "fakepassword",
			"CACHE_DB_MAX_CON":  "100",
			"CACHE_DB_MIN_CON":  "10",
			"CACHE_DB_MAX_TIME": "10",
			"CACHE_DB_MIN_TIME": "2",

			// Casbin
			"MODEL_PATH": "resource/model",

			// FiberConfig
			"FIBER_HOST":           "localhost",
			"FIBER_PORT":           "3000",
			"FIBER_PREFORK":        "false",
			"FIBER_STRICT_ROUTING": "true",
			"FIBER_CASE_SENSITIVE": "true",
			"FIBER_BODY_LIMIT":     "4",
			"FIBER_READ_TIMEOUT":   "4",
			"FIBER_WRITE_TIMEOUT":  "5",
			"FIBER_REDUCE_MEMU":    "true",
			"FIBER_JSON":           "json",

			// Bcrypt
			"HASH_SALT": "10",

			// LoggerConfig
			"LOG_PATH":             "resource/logs",
			"LOG_MAX_SIZE":         "10",
			"LOG_MAX_BACKUP":       "5",
			"LOG_MAX_SIZE_ROTATE":  "20",
			"LOG_TIME_FORMAT":      "2006-01-02",
			"LOG_COLOR_OUTPUT":     "false",
			"LOG_QUOTE_STR":        "false",
			"LOG_END_WITH_MESSAGE": "false",

			// SqlConfig
			"DB_DRIVER":   "postgres",
			"DB_PROTOCOL": "postgresql",
			"DB_NAME":     "test",
			"DB_HOST":     "localhost",
			"DB_PORT":     "5432",
			"DB_USER":     "postgres",
			"DB_PASS":     "postgres",
			"DB_MAX_CON":  "100",
			"DB_MIN_CON":  "10",
			"DB_MAX_TIME": "30",
			"DB_MIN_TIME": "5",

			// JWTToken
			"TOKEN_NAME": "RESTful_API_AUTH",

			// SecretKey
			"SECRET_KEY_ACCESS_TOKEN":  "Zawssh9t1IY50IlICrYpjrCDbq6G8UKL",
			"SECRET_KEY_REFRESH_TOKEN": "ZxywJoMMXIXwgcispeKs4L6Y65XgATqV",
			"SECRET_KEY_FP_TOKEN":      "BNXWuiMew8HhFHLirNw1zpOtO0aJW1cE",
		}, isError: true},
	}
	for i, testCase := range testCases {
		nameCase := fmt.Sprintf("%d. Case : %s", i+1, testCase.name)
		t.Run(nameCase, func(t *testing.T) {
			SetUpEnv(testCase.env)
			defer ClearEnv(testCase.env)
			app, err := bootstrap.New()
			if err != nil {
				require.Equal(t, testCase.isError, err != nil)
				return
			}
			require.NotNil(t, app)
			require.NoError(t, err)
		})
	}
}
