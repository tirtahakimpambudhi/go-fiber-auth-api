package sqlconfig_test

import (
	"github.com/stretchr/testify/require"
	sqlconfig "github.com/tirtahakimpambudhi/restful_api/internal/configs/sql"
	"os"
	"testing"
)

func TestNewConfig_Success(t *testing.T) {
	os.Setenv("DB_DRIVER", "postgres")
	os.Setenv("DB_PROTOCOL", "tcp")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASS", "password")

	defer os.Unsetenv("DB_DRIVER")
	defer os.Unsetenv("DB_PROTOCOL")
	defer os.Unsetenv("DB_NAME")
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_PORT")
	defer os.Unsetenv("DB_USER")
	defer os.Unsetenv("DB_PASS")

	config, err := sqlconfig.NewConfig()

	require.NoError(t, err)
	require.NotNil(t, config)
	require.Equal(t, "postgres", config.Driver)
	require.Equal(t, "tcp", config.Protocol)
	require.Equal(t, "testdb", config.Name)
	require.Equal(t, "localhost", config.Host)
	require.Equal(t, "5432", config.Port)
	require.Equal(t, "user", config.User)
	require.Equal(t, "password", config.Password)
}

func TestNewConfig_MissingEnvVars(t *testing.T) {
	os.Clearenv()

	config, err := sqlconfig.NewConfig()

	require.Error(t, err)
	require.Nil(t, config)
}
