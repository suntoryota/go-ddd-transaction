package transaction

import (
	"context"
	"testing"

	"onlineShop/config"
	"onlineShop/database"

	"github.com/stretchr/testify/require"
)

var svc service

func init() {
	filename := "../cmd/api/config.yaml"
	err := config.LoadConfig(filename)
	if err != nil {
		panic(err)
	}

	db, err := database.ConnectPostgres(config.Cfg.DB)
	if err != nil {
		panic(err)
	}
	repo := newRepository(db)
	svc = newService(repo)
}

func TestCreateTransaction(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := CreateTransactionRequestPayload{
			ProductSKU:   "122083dc-0e92-45ee-a098-59aa8b91a62d",
			Amount:       2,
			UserPublicId: "15e95ba0-352c-4346-8e28-af0aa591693f",
		}
		err := svc.CreateTransaction(context.Background(), req)
		require.Nil(t, err)
	})
}
