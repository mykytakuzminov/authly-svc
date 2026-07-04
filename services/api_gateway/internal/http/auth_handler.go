package http

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"

	"github.com/mykytakuzminov/ridely-svc/services/api_gateway/internal/domain"
	"github.com/mykytakuzminov/ridely-svc/shared/contracts"
	"github.com/mykytakuzminov/ridely-svc/shared/errors"
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
func (h *authHTTPHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var reqBody domain.RegisterInput
	if err := ReadJSON(w, r, &reqBody); err != nil {
		h.logger.Errorw("failed to parse JSON data", "error", err)
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(&reqBody); err != nil {
		h.logger.Warnw("registration data validation failed", "error", err)
		http.Error(w, "registration data validation failed", http.StatusBadRequest)
		return
	}

	tokens, err := h.authClient.Register(r.Context(), reqBody.ToProto())
	if err != nil {
		code, msg := errors.ToHTTPError(err)
		if code >= 500 {
			h.logger.Errorw("register failed", "error", err)
		} else {
			h.logger.Warnw("register rejected", "error", err, "status", code)
		}
		http.Error(w, msg, code)
		return
	}

	response := contracts.APIResponse{Data: tokens}
	if err := WriteJSON(w, http.StatusCreated, response); err != nil {
		h.logger.Errorw("failed to write response", "error", err)
	}
}
