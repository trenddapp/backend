package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/trenddapp/backend/pkg/auth"
	internalhttp "github.com/trenddapp/backend/pkg/http"
	"github.com/trenddapp/backend/service/wordle/pkg/model"
	"github.com/trenddapp/backend/service/wordle/pkg/repository/word"
	"github.com/trenddapp/backend/service/wordle/pkg/repository/wordle"
	"github.com/trenddapp/backend/service/wordle/pkg/workflow"
)

type Server struct {
	workflowEngine   *workflow.Engine
	wordRepository   word.Repository
	wordleRepository wordle.Repository
}

func NewServer(
	workflowEngine *workflow.Engine,
	wordRepository word.Repository,
	wordleRepository wordle.Repository,
) *Server {
	return &Server{
		workflowEngine:   workflowEngine,
		wordRepository:   wordRepository,
		wordleRepository: wordleRepository,
	}
}

func (s Server) GetWordle(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		internalhttp.NewError(http.StatusBadRequest, "invalid argument id").WriteJSON(ctx)
		return
	}

	userID := ctx.Param("user_id")
	if userID == "" {
		internalhttp.NewError(http.StatusBadRequest, "invalid argument user_id").WriteJSON(ctx)
		return
	}

	if !isAuthorized(ctx, userID) {
		internalhttp.NewError(http.StatusUnauthorized, "invalid token").WriteJSON(ctx)
		return
	}

	out, err := s.workflowEngine.GetWordle(ctx, id, wordle.WithUserID(userID))
	if err != nil {
		if errors.Is(err, wordle.ErrWordleNotFound) {
			internalhttp.NewError(http.StatusNotFound, "wordle not found").WriteJSON(ctx)
			return
		}

		internalhttp.NewError(http.StatusInternalServerError, "unknown error occurred").WriteJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"wordle": out})
}

func (s Server) CreateWordle(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		internalhttp.NewError(http.StatusBadRequest, "invalid argument user_id").WriteJSON(ctx)
		return
	}

	if !isAuthorized(ctx, userID) {
		internalhttp.NewError(http.StatusUnauthorized, "invalid token").WriteJSON(ctx)
		return
	}

	var in model.Wordle

	if err := ctx.ShouldBindJSON(&in); err != nil {
		internalhttp.NewError(http.StatusBadRequest, "invalid argument body").WriteJSON(ctx)
		return
	}

	in.UserID = userID

	out, err := s.workflowEngine.CreateWordle(ctx, in)
	if err != nil {
		internalhttp.NewError(http.StatusInternalServerError, "unknown error occurred").WriteJSON(ctx)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"wordle": out})
}

func (s Server) UpdateWordle(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		internalhttp.NewError(http.StatusBadRequest, "invalid argument id").WriteJSON(ctx)
		return
	}

	userID := ctx.Param("user_id")
	if userID == "" {
		internalhttp.NewError(http.StatusBadRequest, "invalid argument user_id").WriteJSON(ctx)
		return
	}

	if !isAuthorized(ctx, userID) {
		internalhttp.NewError(http.StatusUnauthorized, "invalid token").WriteJSON(ctx)
		return
	}

	var in model.Wordle

	if err := ctx.ShouldBindJSON(&in); err != nil {
		internalhttp.NewError(http.StatusBadRequest, "invalid argument body").WriteJSON(ctx)
		return
	}

	out, err := s.workflowEngine.UpdateWordle(ctx, in, wordle.WithUserID(userID))
	if err != nil {
		if errors.Is(err, wordle.ErrWordleNotFound) {
			internalhttp.NewError(http.StatusNotFound, "wordle not found").WriteJSON(ctx)
			return
		}

		internalhttp.NewError(http.StatusInternalServerError, err.Error()).WriteJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"wordle": out})
}

func (s Server) DeleteWordle(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		internalhttp.NewError(http.StatusBadRequest, "invalid argument id").WriteJSON(ctx)
		return
	}

	userID := ctx.Param("user_id")
	if userID == "" {
		internalhttp.NewError(http.StatusBadRequest, "invalid argument user_id").WriteJSON(ctx)
		return
	}

	if !isAuthorized(ctx, userID) {
		internalhttp.NewError(http.StatusUnauthorized, "invalid token").WriteJSON(ctx)
		return
	}

	if _, err := s.wordleRepository.DeleteWordle(ctx, id, wordle.WithUserID(userID)); err != nil {
		if errors.Is(err, wordle.ErrWordleNotFound) {
			internalhttp.NewError(http.StatusNotFound, "wordle not found")
			return
		}

		internalhttp.NewError(http.StatusInternalServerError, "unknown error occurred").WriteJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (s Server) ListWordles(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		internalhttp.NewError(http.StatusBadRequest, "invalid argument user_id").WriteJSON(ctx)
		return
	}

	if !isAuthorized(ctx, userID) {
		internalhttp.NewError(http.StatusUnauthorized, "invalid token").WriteJSON(ctx)
		return
	}

	options := []wordle.Option{
		wordle.WithPageToken(ctx.Query("pageToken")),
		wordle.WithUserID(userID),
	}

	if ctx.Query("pageSize") != "" {
		pageSize, err := strconv.Atoi(ctx.Query("pageSize"))
		if err != nil {
			internalhttp.NewError(http.StatusBadRequest, "invalid argument page_size").WriteJSON(ctx)
			return
		}

		options = append(options, wordle.WithPageSize(int32(pageSize)))
	}

	// TODO: Use "filter" as query parameter.
	if ctx.Query("status") != "" {
		statusToModel := map[string]model.Status{
			"open":     model.StatusOpen,
			"complete": model.StatusComplete,
			"canceled": model.StatusCanceled,
		}

		status, ok := statusToModel[ctx.Query("status")]
		if !ok {
			internalhttp.NewError(http.StatusBadRequest, "invalid argument status").WriteJSON(ctx)
			return
		}

		options = append(options, wordle.WithStatus(status))
	}

	wordles, nextPageToken, err := s.workflowEngine.ListWordles(ctx, options...)
	if err != nil {
		if errors.Is(err, wordle.ErrWordleNotFound) {
			internalhttp.NewError(http.StatusNotFound, "wordles not found").WriteJSON(ctx)
			return
		}

		internalhttp.NewError(http.StatusInternalServerError, err.Error()).WriteJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"next_page_token": nextPageToken,
		"wordles":         wordles,
	})
}

func isAuthorized(ctx *gin.Context, userID string) bool {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		return false
	}

	claims, err := auth.ParseToken(token)
	if err != nil {
		return false
	}

	return userID == claims["id"].(string)
}
