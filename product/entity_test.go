package product

import (
	"testing"

	"onlineShop/response"

	"github.com/stretchr/testify/require"
)

func TestValidateProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		product := Product{
			Name:  "Notebook v1",
			Stock: 100,
			Price: 10_000_000,
		}

		err := product.Validate()
		require.Nil(t, err)
	})

	t.Run("product required", func(t *testing.T) {
		product := Product{
			Name:  "",
			Stock: 100,
			Price: 10_000_000,
		}

		err := product.Validate()
		require.NotNil(t, err)
		require.Equal(t, response.ErrProductRequired, err)
	})
	t.Run("product invalid", func(t *testing.T) {
		product := Product{
			Name:  "No",
			Stock: 100,
			Price: 10_000_000,
		}

		err := product.Validate()
		require.NotNil(t, err)
		require.Equal(t, response.ErrProductInvalid, err)
	})

	t.Run("stock invalid", func(t *testing.T) {
		product := Product{
			Name:  "Notebook v1",
			Stock: 0,
			Price: 10_000_000,
		}

		err := product.Validate()
		require.NotNil(t, err)
		require.Equal(t, response.ErrStockInvalid, err)
	})

	t.Run("price invalid", func(t *testing.T) {
		product := Product{
			Name:  "Notebook v1",
			Stock: 100,
			Price: 0,
		}

		err := product.Validate()
		require.NotNil(t, err)
		require.Equal(t, response.ErrPriceInvalid, err)
	})
}
