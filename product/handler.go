package product

import (
	"net/http"

	"onlineShop/response"

	"github.com/gofiber/fiber/v2"
)

type handler struct {
	svc service
}

func newHandler(svc service) handler {
	return handler{
		svc: svc,
	}
}

func (h handler) CreateProduct(ctx *fiber.Ctx) error {
	req := CreateProductRequestPayload{}

	if err := ctx.BodyParser(&req); err != nil {
		return response.NewResponse(
			response.WithMessage("invalid payload"),
			response.WithError(response.ErrorBadRequest),
		).Send(ctx)
	}

	if err := h.svc.CreateProduct(ctx.UserContext(), req); err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(myErr),
		).Send(ctx)
	}

	return response.NewResponse(
		response.WithHttpCode(http.StatusCreated),
		response.WithMessage("create product success"),
	).Send(ctx)
}

func (h handler) GetListProducts(ctx *fiber.Ctx) error {
	req := ListProductRequestPayload{}

	if err := ctx.QueryParser(&req); err != nil {
		return response.NewResponse(
			response.WithMessage("invalid payload"),
			response.WithError(response.ErrorBadRequest),
		).Send(ctx)
	}
	products, err := h.svc.ListProducts(ctx.UserContext(), req)
	if err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}
		return response.NewResponse(
			response.WithMessage("invalid payload"),
			response.WithError(myErr),
		).Send(ctx)
	}

	productListResponse := NewProductListResponseFromEntity(products)

	return response.NewResponse(
		response.WithHttpCode(http.StatusOK),
		response.WithMessage("get list products success"),
		response.WithPayload(productListResponse),
		response.WithQuery(req.GenerateDefaultValue()),
	).Send(ctx)
}

func (h handler) GetProductDetail(ctx *fiber.Ctx) error {
	sku := ctx.Params("sku", "")
	if sku == "" {
		return response.NewResponse(
			response.WithMessage("invalid payload"),
			response.WithError(response.ErrorBadRequest),
		).Send(ctx)
	}

	product, err := h.svc.ProductDetail(ctx.UserContext(), sku)
	if err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(myErr),
		).Send(ctx)
	}

	productDetail := ProductDetailResponse{
		Id:        product.Id,
		Name:      product.Name,
		SKU:       product.SKU,
		Stock:     product.Stock,
		Price:     product.Price,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}

	return response.NewResponse(
		response.WithHttpCode(http.StatusOK),
		response.WithMessage("get product detail success"),
		response.WithPayload(productDetail),
	).Send(ctx)
}
