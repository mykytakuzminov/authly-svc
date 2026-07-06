package http

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"

	"github.com/mykytakuzminov/ridely-svc/services/api_gateway/internal/domain"
	c "github.com/mykytakuzminov/ridely-svc/shared/context"
	"github.com/mykytakuzminov/ridely-svc/shared/contracts"
	"github.com/mykytakuzminov/ridely-svc/shared/errors"
	log "github.com/mykytakuzminov/ridely-svc/shared/logging"
	authpb "github.com/mykytakuzminov/ridely-svc/shared/proto/auth"
)

type authHTTPHandler struct {
	authClient authpb.AuthServiceClient
	validator  *validator.Validate
	logger     *zap.SugaredLogger
}

func NewAuthHTTPHandler(
	authClient authpb.AuthServiceClient,
	validator *validator.Validate,
	logger *zap.SugaredLogger,
) *authHTTPHandler {
	return &authHTTPHandler{
		authClient: authClient,
		validator:  validator,
		logger:     logger,
	}
}

// Register godoc
// @Summary     Register a new user
// @Tags        auth
// @Accept      json
// @Produce     json
// @Router      /auth/register [post]
func (h *authHTTPHandler) Register(w http.ResponseWriter, r *http.Request) {
	traceID := c.GetTraceID(r.Context())

	var reqBody domain.RegisterInput
	if err := ReadJSON(w, r, &reqBody); err != nil {
		log.FailedParseJSON(h.logger, traceID, err)
		responseBadRequest(w, "failed to parse JSON data")
		return
	}

	if err := h.validator.Struct(&reqBody); err != nil {
		log.FailedValidateData(h.logger, traceID, err)
		responseBadRequest(w, "failed to validate data")
		return
	}

	tokens, err := h.authClient.Register(r.Context(), reqBody.ToProto())
	if err != nil {
		code, msg := errors.ToHTTPError(err)
		if code >= 500 {
			log.Failed(h.logger, traceID, "registration", err)
		} else {
			log.Declined(h.logger, traceID, "registration", err)
		}
		http.Error(w, msg, code)
		return
	}

	response := contracts.APIResponse{Data: tokens}
	if err := WriteJSON(w, http.StatusCreated, response); err != nil {
		log.FailedWriteResponse(h.logger, traceID, err)
	}
}

func (h *authHTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	traceID := c.GetTraceID(r.Context())

	var reqBody domain.LoginInput
	if err := ReadJSON(w, r, &reqBody); err != nil {
		log.FailedParseJSON(h.logger, traceID, err)
		responseBadRequest(w, "failed to parse JSON data")
		return
	}

	if err := h.validator.Struct(&reqBody); err != nil {
		log.FailedValidateData(h.logger, traceID, err)
		responseBadRequest(w, "failed to validate data")
		return
	}

	tokens, err := h.authClient.Login(r.Context(), reqBody.ToProto())
	if err != nil {
		code, msg := errors.ToHTTPError(err)
		if code >= 500 {
			log.Failed(h.logger, traceID, "login", err)
		} else {
			log.Declined(h.logger, traceID, "login", err)
		}
		http.Error(w, msg, code)
		return
	}

	response := contracts.APIResponse{Data: tokens}
	if err := WriteJSON(w, http.StatusOK, response); err != nil {
		log.FailedWriteResponse(h.logger, traceID, err)
	}
}
