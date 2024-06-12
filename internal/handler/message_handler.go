package handler

import (
	"context"

	"github.com/anilsenay/message-sending-system/internal/model"
	"github.com/anilsenay/message-sending-system/pkg/logger"
	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/http/response"
	"github.com/go-swagno/swagno/components/parameter"
	"github.com/gofiber/fiber/v2"
)

type messageService interface {
	RetireveSentMessages(ctx context.Context) ([]model.Message, error)
	CreateMessage(ctx context.Context, m *model.Message) error
	StartMessageSending(ctx context.Context)
	StopMessageSending()
}

type MessageHandler struct {
	messageService messageService
}

func NewMessageHandler(messageService messageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

func (h *MessageHandler) handleGetSentMessages(ctx *fiber.Ctx) error {
	messages, err := h.messageService.RetireveSentMessages(ctx.UserContext())
	if err != nil {
		logger.Error().Str("url", ctx.BaseURL()).Msgf("error while getting sent messages: %s", err.Error())
		return handleError(ctx, fiber.StatusInternalServerError, "Some error occurred", err.Error())
	}

	return handleSuccess(ctx, fiber.StatusOK, messages)
}

func (h *MessageHandler) handleCreateMessage(ctx *fiber.Ctx) error {
	messageRequest := model.MessageCreateRequest{}

	if err := ctx.BodyParser(&messageRequest); err != nil {
		return handleError(ctx, fiber.StatusUnprocessableEntity, "invalid request body", err.Error())
	}

	if err := messageRequest.Validate(); err != nil {
		return handleError(ctx, fiber.StatusBadRequest, "validation error", err.Error())
	}

	message := model.Message{
		Content:              messageRequest.Content,
		RecipientPhoneNumber: messageRequest.RecipientPhoneNumber,
	}

	err := h.messageService.CreateMessage(ctx.UserContext(), &message)
	if err != nil {
		return handleError(ctx, fiber.StatusBadRequest, "error while creating message", err.Error())
	}

	return handleSuccess(ctx, fiber.StatusOK, message)
}

func (h *MessageHandler) handleStartOrStopWorker(ctx *fiber.Ctx) error {
	status := ctx.Params("status")

	if status == "start" {
		// use background context instead of UserContext to prevent cancellation after request
		h.messageService.StartMessageSending(context.Background())
	} else if status == "stop" {
		h.messageService.StopMessageSending()
	} else {
		return handleError(ctx, fiber.StatusBadRequest, "status parameter is missing", "")
	}

	return handleSuccess(ctx, fiber.StatusOK, "OK")
}

func (h *MessageHandler) SetRoutes(app *fiber.App) {
	messagesGroup := app.Group("/messages")

	messagesGroup.Get("/", h.handleGetSentMessages)
	messagesGroup.Post("/", h.handleCreateMessage)
	messagesGroup.Post("/worker/:status", h.handleStartOrStopWorker)
}

func (h *MessageHandler) GetSwaggerEndpoints() []*endpoint.EndPoint {
	return []*endpoint.EndPoint{
		endpoint.New(
			endpoint.GET, "/messages",
			endpoint.WithDescription("Retrieve sent messages"),
			endpoint.WithTags("Messages"),
			endpoint.WithSuccessfulReturns([]response.Response{
				response.New([]model.Message{}, "200", ""),
			}),
			endpoint.WithErrors([]response.Response{
				response.New(FailDetails{}, "500", ""),
			}),
		),
		endpoint.New(
			endpoint.POST, "/messages",
			endpoint.WithDescription("Create a new message"),
			endpoint.WithTags("Messages"),
			endpoint.WithBody(model.MessageCreateRequest{}),
			endpoint.WithSuccessfulReturns([]response.Response{
				response.New(model.Message{}, "200", ""),
			}),
			endpoint.WithErrors([]response.Response{
				response.New(FailDetails{}, "500", ""),
			}),
		),
		endpoint.New(
			endpoint.POST, "/messages/worker/{status}",
			endpoint.WithDescription("Start or Stop message sending worker"),
			endpoint.WithTags("Messages"),
			endpoint.WithParams(parameter.StrEnumParam("status", parameter.Path, []string{"start", "stop"}, parameter.WithRequired())),
			endpoint.WithSuccessfulReturns([]response.Response{
				response.New("", "200", ""),
			}),
			endpoint.WithErrors([]response.Response{
				response.New(FailDetails{}, "400", ""),
			}),
		),
	}
}
