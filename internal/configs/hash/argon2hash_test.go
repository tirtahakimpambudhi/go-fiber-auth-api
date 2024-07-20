package hash_test

import (
	"fmt"
	"github.com/alexedwards/argon2id"
	"github.com/stretchr/testify/require"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/hash"

	"testing"
)

// Successfully create a new Argon2 hash instance with default parameters
func TestNewHashArgon2_Success(t *testing.T) {
	hashInstance, err := hash.NewHashArgon2()
	require.NoError(t, err)
	require.NotNil(t, hashInstance)
	require.Equal(t, hashInstance.Params, argon2id.DefaultParams)
}

func TestHashArgon2id(t *testing.T) {
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
			argon2idInstance, err := hash.NewHashArgon2()
			require.NoError(t, err)
			require.NotNil(t, argon2idInstance)

			hashPassword, err := argon2idInstance.Create(testcase.actualPassword)
			require.NoError(t, err)
			match, err := argon2idInstance.Match(testcase.expectPassword, hashPassword)
			require.NoError(t, err)
			require.Equal(t, testcase.isMatch, match)
		})
	}
}
