// Handles error when loading Casbin configuration fails
package casbinconfig_test

import (
	"github.com/stretchr/testify/require"
	casbinconfig "github.com/tirtahakimpambudhi/restful_api/internal/configs/casbin"
	"testing"
)

func TestNewCasbinConfigDBError(t *testing.T) {
	middleware, err := casbinconfig.NewCasbin(nil, nil)
	require.Error(t, err)
	require.Nil(t, middleware)
}
