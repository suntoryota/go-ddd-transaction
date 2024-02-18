package auth

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

func (h handler) register(ctx *fiber.Ctx) error {
	req := RegisterRequestPayload{}

	if err := ctx.BodyParser(&req); err != nil {
		myErr := response.ErrorBadRequest
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(myErr),
			response.WithHttpCode(http.StatusBadRequest),
			response.WithMessage("register fail"),
		).Send(ctx)
	}

	if err := h.svc.register(ctx.UserContext(), req); err != nil {
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
		response.WithMessage("register success"),
	).Send(ctx)
}

func (h handler) login(ctx *fiber.Ctx) error {
	req := LoginRequestPayload{}

	if err := ctx.BodyParser(&req); err != nil {
		myErr := response.ErrorBadRequest
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(myErr),
			response.WithMessage("login fail"),
		).Send(ctx)
	}
	token, err := h.svc.login(ctx.UserContext(), req)
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

	return response.NewResponse(
		response.WithHttpCode(http.StatusCreated),
		response.WithPayload(map[string]interface{}{
			"access_token": token,
		}),
		response.WithMessage("login success"),
	).Send(ctx)
}
