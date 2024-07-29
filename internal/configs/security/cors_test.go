package security_test

import (
	"github.com/stretchr/testify/require"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/security"
	"testing"
)

func TestCorsConfig_Success(t *testing.T) {
	t.Setenv("CORS_ALLOW_METHODS", "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS")
	t.Setenv("CORS_ALLOW_HEADERS", "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin, Authorization")
	t.Setenv("CORS_EXPOSE_HEADERS", "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin, Authorization")
	t.Setenv("CORS_ALLOW_ORIGINS", "localhost")
	t.Setenv("CORS_ALLOW_CREDENTIALS", "true")

	cors, err := security.NewCors()
	require.NotNil(t, cors)
	require.NoError(t, err)
}

func TestCorsConfig_Fail(t *testing.T) {
	cors, err := security.NewCors()
	require.Nil(t, cors)
	require.Error(t, err)
}
