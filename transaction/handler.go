package transaction

import (
	"fmt"
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

func (h handler) CreateTransaction(ctx *fiber.Ctx) error {
	req := CreateTransactionRequestPayload{}

	if err := ctx.BodyParser(&req); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(err),
			response.WithHttpCode(http.StatusBadRequest),
		).Send(ctx)
	}

	userPublicId := ctx.Locals("PUBLIC_ID")
	req.UserPublicId = fmt.Sprintf("%v", userPublicId)

	if err := h.svc.CreateTransaction(ctx.UserContext(), req); err != nil {
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
		response.WithMessage("create transactions success"),
	).Send(ctx)
}

func (h handler) GetTransactionByUser(ctx *fiber.Ctx) error {
	userPublicId := fmt.Sprintf("%v", ctx.Locals("PUBLIC_ID"))

	trxs, err := h.svc.TransactionHistory(ctx.UserContext(), userPublicId)
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

	responses := []TransactionHistoryResponse{}

	for _, trx := range trxs {
		responses = append(responses, trx.ToTransactionHistoryResponse())
	}

	return response.NewResponse(
		response.WithHttpCode(http.StatusCreated),
		response.WithPayload(responses),
		response.WithMessage("get transaction histories success"),
	).Send(ctx)
}
