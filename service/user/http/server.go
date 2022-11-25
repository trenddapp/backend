package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/trenddapp/backend/pkg/auth"
	internalhttp "github.com/trenddapp/backend/pkg/http"
	"github.com/trenddapp/backend/service/user/model"
	"github.com/trenddapp/backend/service/user/repository/nonce"
	"github.com/trenddapp/backend/service/user/repository/user"
)

type Server struct {
	logger          *zap.Logger
	nonceRepository nonce.Repository
	userRepository  user.Repository
}

func NewServer(
	logger *zap.Logger,
	nonceRepository nonce.Repository,
	userRepository user.Repository,
) *Server {
	return &Server{
		logger:          logger,
		nonceRepository: nonceRepository,
		userRepository:  userRepository,
	}
}

func (s Server) GetSession(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		internalhttp.NewError(http.StatusBadRequest, "invalid argument user_id").WriteJSON(ctx)
		return
	}

	user, err := s.userRepository.GetUser(ctx, userID)
	if err != nil {
		// TODO: Check for not found error.
		internalhttp.NewError(http.StatusNotFound, "user not found").WriteJSON(ctx)
		return
	}

	nonceStr, err := auth.NewNonce()
	if err != nil {
		internalhttp.NewError(http.StatusInternalServerError, "failed to generate nonce").WriteJSON(ctx)
		return
	}

	in := model.Nonce{
		UserID: user.ID,
		Value:  nonceStr,
	}

	nonce, err := s.nonceRepository.CreateNonce(ctx, in)
	if err != nil {
		internalhttp.NewError(http.StatusInternalServerError, "failed to create nonce").WriteJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"session": gin.H{
			"nonce": nonce.Value,
		},
	})
}

func (s Server) CreateSession(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		internalhttp.NewError(http.StatusBadRequest, "invalid argument user_id").WriteJSON(ctx)
		return
	}

	var body struct {
		SignedMessage string `json:"signed_message"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		internalhttp.NewError(http.StatusBadRequest, "invalid argument body").WriteJSON(ctx)
		return
	}

	if body.SignedMessage == "" {
		internalhttp.NewError(http.StatusBadRequest, "invalid argument signed_message").WriteJSON(ctx)
		return
	}

	user, err := s.userRepository.GetUser(ctx, userID)
	if err != nil {
		s.logger.Error(
			"failed to get user",
			zap.String("id", userID),
			zap.Error(err),
		)
		internalhttp.NewError(http.StatusNotFound, "user not found").WriteJSON(ctx)
		return
	}

	nonce, err := s.nonceRepository.GetNonceByUserID(ctx, user.ID)
	if err != nil {
		s.logger.Error(
			"failed to get nonce by user id",
			zap.String("id", userID),
			zap.Error(err),
		)
		internalhttp.NewError(http.StatusNotFound, "nonce not found").WriteJSON(ctx)
		return
	}

	defer func() {
		if _, err := s.nonceRepository.DeleteNonce(ctx, nonce.ID); err != nil {
			s.logger.Error(
				"failed to delete nonce",
				zap.String("id", nonce.ID),
				zap.Error(err),
			)
		}
	}()

	isValid, err := auth.VerifySignature(user.Address, nonce.Value, body.SignedMessage)
	if !isValid || err != nil {
		s.logger.Error(
			"failed to verify wallet signature",
			zap.String("address", user.Address),
			zap.String("nonce", nonce.Value),
			zap.String("signature", body.SignedMessage),
			zap.String("user_id", userID),
			zap.Error(err),
		)
		internalhttp.NewError(http.StatusUnauthorized, "invalid credentials").WriteJSON(ctx)
		return
	}

	claims := auth.Claims{
		"id":      user.ID,
		"address": user.Address,
		"exp":     time.Now().AddDate(0, 0, 7).Unix(),
	}

	token, err := auth.NewToken(claims)
	if err != nil {
		s.logger.Error(
			"failed to create a new token",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		internalhttp.NewError(http.StatusInternalServerError, "failed to generate JWT").WriteJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"session": gin.H{
			"token": token,
		},
	})
}

func (s Server) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		internalhttp.NewError(http.StatusBadRequest, "invalid argument id").WriteJSON(ctx)
		return
	}

	var user model.User
	if strings.Contains(id, "address=") {
		var err error
		user, err = s.userRepository.GetUserByAddress(ctx, strings.TrimPrefix(id, "address="))
		if err != nil {
			// TODO: Check for not found error.
			internalhttp.NewError(http.StatusNotFound, "user not found").WriteJSON(ctx)
			return
		}
	} else {
		var err error
		user, err = s.userRepository.GetUser(ctx, id)
		if err != nil {
			// TODO: Check for not found error.
			internalhttp.NewError(http.StatusNotFound, "user not found").WriteJSON(ctx)
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (s Server) CreateUser(ctx *gin.Context) {
	var in model.User

	if err := ctx.ShouldBindJSON(&in); err != nil {
		internalhttp.NewError(http.StatusBadRequest, "invalid argument body").WriteJSON(ctx)
		return
	}

	if in.Address == "" {
		internalhttp.NewError(http.StatusBadRequest, "invalid argument address").WriteJSON(ctx)
		return
	}

	user, err := s.userRepository.CreateUser(ctx, in)
	if err != nil {
		internalhttp.NewError(http.StatusInternalServerError, "invalid operation").WriteJSON(ctx)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": user})
}

func (s Server) UpdateUser(ctx *gin.Context) {
}

func (s Server) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		internalhttp.NewError(http.StatusBadRequest, "invalid argument id").WriteJSON(ctx)
		return
	}

	if _, err := s.userRepository.DeleteUser(ctx, id); err != nil {
		// TODO: Check for not found error.
		internalhttp.NewError(http.StatusNotFound, "user not found").WriteJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
