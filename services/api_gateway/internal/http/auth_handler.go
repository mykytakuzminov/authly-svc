package http

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"

	"github.com/mykytakuzminov/ridely-svc/services/api_gateway/internal/domain"
	c "github.com/mykytakuzminov/ridely-svc/shared/context"
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

// @Summary     Register new user
// @Description Register new user with credentials and return tokens
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body domain.RegisterInput true "Register credentials"
// @Success     200 {object} SuccessResponse{data=domain.TokensResponse}
// @Failure     400 {object} ErrorResponse "Bad request"
// @Failure     409 {object} ErrorResponse "Conflict"
// @Failure     500 {object} ErrorResponse "Internal server error"
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
		responseError(w, code, msg)
		return
	}

	responseSuccess(w, domain.TokensResponseFromProto(tokens))
}

// @Summary     Authenticate user
// @Description Authenticate user with credentials and return tokens
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body domain.LoginInput true "Login credentials"
// @Success     200 {object} SuccessResponse{data=domain.TokensResponse}
// @Failure     400 {object} ErrorResponse "Bad request"
// @Failure     401 {object} ErrorResponse "Unauthenticated"
// @Failure     500 {object} ErrorResponse "Internal server error"
// @Router      /auth/login [post]
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
		responseError(w, code, msg)
		return
	}

	responseSuccess(w, domain.TokensResponseFromProto(tokens))
}

// @Summary     Refresh tokens
// @Description Generate new tokens and return them to client
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body domain.RefreshInput true "Refresh token"
// @Success     200 {object} SuccessResponse{data=domain.TokensResponse}
// @Failure     400 {object} ErrorResponse "Bad request"
// @Failure     401 {object} ErrorResponse "Unauthenticated"
// @Failure     500 {object} ErrorResponse "Internal server error"
// @Router      /auth/refresh [post]
func (h *authHTTPHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	traceID := c.GetTraceID(r.Context())

	var reqBody domain.RefreshInput
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

	tokens, err := h.authClient.Refresh(r.Context(), reqBody.ToProto())
	if err != nil {
		code, msg := errors.ToHTTPError(err)
		if code >= 500 {
			log.Failed(h.logger, traceID, "token refreshing", err)
		} else {
			log.Declined(h.logger, traceID, "token refreshing", err)
		}
		responseError(w, code, msg)
		return
	}

	responseSuccess(w, domain.TokensResponseFromProto(tokens))
}

// @Summary     Logout user
// @Description Logout user and delete refresh token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body domain.LogoutInput true "Delete token"
// @Security    BearerAuth
// @Success     204
// @Failure     400 {object} ErrorResponse "Bad request"
// @Failure     401 {object} ErrorResponse "Unauthenticated"
// @Failure     500 {object} ErrorResponse "Internal server error"
// @Router      /auth/logout [post]
func (h *authHTTPHandler) Logout(w http.ResponseWriter, r *http.Request) {
	traceID := c.GetTraceID(r.Context())

	var reqBody domain.LogoutInput
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

	_, err := h.authClient.Logout(r.Context(), reqBody.ToProto())
	if err != nil {
		code, msg := errors.ToHTTPError(err)
		if code >= 500 {
			log.Failed(h.logger, traceID, "logging out", err)
		} else {
			log.Declined(h.logger, traceID, "logging out", err)
		}
		responseError(w, code, msg)
		return
	}

	responseNoContent(w)
}
