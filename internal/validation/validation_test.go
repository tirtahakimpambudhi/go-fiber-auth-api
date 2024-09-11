package validation_test

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
	"github.com/tirtahakimpambudhi/restful_api/internal/validation"
	"testing"
)

func TestValidator(t *testing.T) {
	validate := validator.New()
	english := en.New()
	universalTranslate := ut.New(english, english)
	translator, found := universalTranslate.GetTranslator("en")
	require.True(t, found)
	v := validation.NewValidator(validate, translator)

	type TestStruct struct {
		ID int `validate:"ksuid"`
	}

	testStruct := TestStruct{ID: 123}
	errs := v.Validate(testStruct)
	err := v.ValidateVars(ksuid.New().String(), "ksuid")
	require.Nil(t, err)
	require.NotNil(t, errs)
	require.Equal(t, "Unprocessable Entity", errs[0].Code)
	require.Equal(t, "Validation Error", errs[0].Title)
}
