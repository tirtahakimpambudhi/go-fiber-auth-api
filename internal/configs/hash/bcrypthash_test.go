// NewHashBcrypt handles missing environment variables gracefully
package hash_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/hash"
	"os"
	"testing"
)

func TestNewHashBcryptHandlesMissingEnvVars(t *testing.T) {
	os.Clearenv()
	bcryptInstance, err := hash.NewHashBcrypt()
	require.NoError(t, err)
	require.NotNil(t, bcryptInstance)
	require.Equal(t, 10, bcryptInstance.Salt)
}

func TestNewHashBcryptHandlesIncorrectEnvVars(t *testing.T) {
	os.Setenv("HASH_SALT", "abc")
	defer os.Unsetenv("HASH_SALT")

	bcryptInstance, err := hash.NewHashBcrypt()
	require.Error(t, err)
	require.Nil(t, bcryptInstance)
}

func TestHashBcrypt(t *testing.T) {
	testcases := []struct {
		name           string
		expectPassword string
		actualPassword string
		isMatch        bool
	}{
		{
			name:           "Successfully Create And Match Password",
			expectPassword: "Admin#1234",
			actualPassword: "Admin#1234",
			isMatch:        true,
		},
		{
			name:           "Failure Create And Match Password Because Wrong Password",
			expectPassword: "Admin#1234",
			actualPassword: "4321#nimdA",
			isMatch:        false,
		},
	}

	for i, testcase := range testcases {
		nameCase := fmt.Sprintf("#Case-%d : %s", i+1, testcase.name)
		t.Run(nameCase, func(t *testing.T) {
			bcryptInstance, err := hash.NewHashBcrypt()
			require.NoError(t, err)
			require.NotNil(t, bcryptInstance)

			hashPassword, err := bcryptInstance.Create(testcase.actualPassword)
			require.NoError(t, err)
			match, err := bcryptInstance.Match(testcase.expectPassword, hashPassword)
			require.NoError(t, err)
			require.Equal(t, testcase.isMatch, match)
		})
	}
}
