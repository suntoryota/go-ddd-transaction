package auth

import (
	"log"
	"testing"

	"onlineShop/response"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestValidateAuthEntity(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		authEntity := AuthEntity{
			Email:    "sy@gmail.com",
			Password: "P@ssw0rdsy",
		}

		err := authEntity.Validate()
		require.Nil(t, err)
	})

	t.Run("email is required", func(t *testing.T) {
		authEntity := AuthEntity{
			Email:    "",
			Password: "P@ssw0rdsy",
		}

		err := authEntity.Validate()
		require.NotNil(t, response.ErrEmailRequired, err)
	})

	t.Run("email is invalid", func(t *testing.T) {
		authEntity := AuthEntity{
			Email:    "sygamil.com",
			Password: "P@ssw0rdsy",
		}

		err := authEntity.Validate()
		require.NotNil(t, response.ErrEmailInvalid, err)
	})

	t.Run("password is required", func(t *testing.T) {
		authEntity := AuthEntity{
			Email:    "sy@gmail.com",
			Password: "",
		}

		err := authEntity.Validate()
		require.NotNil(t, response.ErrPasswordRequired, err)
	})

	t.Run("password must have minimum 6 character", func(t *testing.T) {
		authEntity := AuthEntity{
			Email:    "sy@gmail.com",
			Password: "passw",
		}

		err := authEntity.Validate()
		require.NotNil(t, response.ErrPasswordInvalidLength, err)
	})
}

func TestEncryptPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		authEntity := AuthEntity{
			Email:    "sy@gmail.com",
			Password: "P455w0rd",
		}

		err := authEntity.EncryptPassword(bcrypt.DefaultCost)
		require.Nil(t, err)

		log.Printf("%+v\n", authEntity)
	})
}
