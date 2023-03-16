package wordle

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"

	"github.com/trenddapp/backend/pkg/paginator"
	"github.com/trenddapp/backend/service/wordle/pkg/model"
)

var (
	ErrWordleNotFound = errors.New("wordle not found")
)

type Option interface {
	apply(*options)
}

type options struct {
	pageSize  int32
	pageToken string
	status    model.Status
	userID    string
}

type pageSizeOption int32

func (p pageSizeOption) apply(opts *options) {
	opts.pageSize = int32(p)
}

func WithPageSize(in int32) Option {
	return pageSizeOption(in)
}

type pageTokenOption string

func (p pageTokenOption) apply(opts *options) {
	opts.pageToken = string(p)
}

func WithPageToken(in string) Option {
	return pageTokenOption(in)
}

type statusOption model.Status

func (s statusOption) apply(opts *options) {
	opts.status = model.Status(s)
}

func WithStatus(in model.Status) Option {
	return statusOption(in)
}

type userIDOption string

func (s userIDOption) apply(opts *options) {
	opts.userID = string(s)
}

func WithUserID(in string) Option {
	return userIDOption(in)
}

type repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) Repository {
	return &repository{db: db}
}

func (r repository) GetWordle(ctx context.Context, id string, opts ...Option) (model.Wordle, error) {
	var out wordle

	options := options{}
	for _, opt := range opts {
		opt.apply(&options)
	}

	query := r.db.NewSelect().Model(&out).Where("id = ?", id).Limit(1)

	if options.userID != "" {
		query = query.Where("user_id = ?", options.userID)
	}

	if err := query.Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Wordle{}, ErrWordleNotFound
		}

		return model.Wordle{}, err
	}

	return convertToModel(out), nil
}

func (r repository) CreateWordle(ctx context.Context, in model.Wordle) (model.Wordle, error) {
	out := convertFromModel(in)

	_, err := r.db.
		NewInsert().
		Model(&out).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return model.Wordle{}, err
	}

	return convertToModel(out), nil
}

func (r repository) UpdateWordle(ctx context.Context, in model.Wordle, opts ...Option) (model.Wordle, error) {
	out := convertFromModel(in)

	options := options{}
	for _, opt := range opts {
		opt.apply(&options)
	}

	query := r.db.NewUpdate().Model(&out).OmitZero().Where("id = ? AND updated_at = ?", in.ID, in.UpdatedAt).Returning("*")
	if options.userID != "" {
		query.Where("user_id = ?", options.userID)
	}

	result, err := query.Exec(ctx)
	if err != nil {
		return model.Wordle{}, err
	}

	rowsEffected, err := result.RowsAffected()
	if err != nil {
		return model.Wordle{}, err
	}

	if rowsEffected == 0 {
		return model.Wordle{}, ErrWordleNotFound
	}

	return convertToModel(out), nil
}

func (r repository) DeleteWordle(ctx context.Context, id string, opts ...Option) (model.Wordle, error) {
	var out wordle

	options := options{}
	for _, opt := range opts {
		opt.apply(&options)
	}

	query := r.db.NewDelete().Model(&out).Where("id = ?", id).Returning("*")
	if options.userID != "" {
		query.Where("user_id = ?", options.userID)
	}

	result, err := query.Exec(ctx)
	if err != nil {
		return model.Wordle{}, err
	}

	rowsEffected, err := result.RowsAffected()
	if err != nil {
		return model.Wordle{}, err
	}

	if rowsEffected == 0 {
		return model.Wordle{}, ErrWordleNotFound
	}

	return convertToModel(out), nil
}

func (r repository) ListWordles(ctx context.Context, opts ...Option) ([]model.Wordle, string, error) {
	var wordles []wordle

	options := options{
		pageSize: 20,
		status:   model.StatusUnspecified,
	}

	for _, opt := range opts {
		opt.apply(&options)
	}

	query := r.db.NewSelect().Model(&wordles).OrderExpr("created_at, id").Limit(int(options.pageSize))

	if options.userID != "" {
		query = query.Where("user_id = ?", options.userID)
	}

	if options.status != model.StatusUnspecified {
		query = query.Where("status = ?", statusFromModel(options.status))
	}

	if options.pageToken != "" {
		id, param, err := paginator.ParsePageToken(options.pageToken)
		if err != nil {
			return nil, "", err
		}

		query = query.Where("(created_at, id) > (?, ?)", param, id)
	}

	if err := query.Scan(ctx); err != nil {
		return nil, "", err
	}

	if len(wordles) == 0 {
		return nil, "", ErrWordleNotFound
	}

	out := make([]model.Wordle, 0, len(wordles))

	for _, wordle := range wordles {
		out = append(out, convertToModel(wordle))
	}

	nextPageToken := ""
	if len(out) == int(options.pageSize) {
		lastItem := out[len(out)-1]
		nextPageToken = paginator.GeneratePageToken(lastItem.ID, lastItem.CreatedAt.Format(time.RFC3339Nano))
	}

	return out, nextPageToken, nil
}

type wordle struct {
	bun.BaseModel `bun:"wordles"`

	ID        uuid.UUID `bun:",pk,default:gen_random_uuid()"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UserID    uuid.UUID `bun:",notnull"`
	Status    string    `bun:",nullzero,notnull,default:'OPEN'"`
	Solution  string    `bun:",notnull"`
	Guesses   []string  `bun:",array"`
}

func convertFromModel(in model.Wordle) wordle {
	return wordle{
		UserID:   uuid.FromStringOrNil(in.UserID),
		Status:   statusFromModel(in.Status),
		Solution: in.Solution,
		Guesses:  in.Guesses,
	}
}

func convertToModel(in wordle) model.Wordle {
	charStatus := make([][]model.CharStatus, 0, len(in.Guesses))
	for i, guess := range in.Guesses {
		charStatus = append(charStatus, make([]model.CharStatus, 0, len(guess)))
		for j, char := range guess {
			if string(char) == string(in.Solution[j]) {
				charStatus[i] = append(charStatus[i], model.CharStatusCorrect)
				continue
			}

			if strings.ContainsRune(in.Solution, char) {
				charStatus[i] = append(charStatus[i], model.CharStatusPresent)
				continue
			}

			charStatus[i] = append(charStatus[i], model.CharStatusNotPresent)
		}
	}

	return model.Wordle{
		ID:         in.ID.String(),
		CreatedAt:  in.CreatedAt,
		UpdatedAt:  in.UpdatedAt,
		UserID:     in.UserID.String(),
		Status:     statusToModel(in.Status),
		Solution:   in.Solution,
		Guesses:    in.Guesses,
		CharStatus: charStatus,
	}
}

func statusFromModel(in model.Status) string {
	statusFromModel := map[model.Status]string{
		model.StatusOpen:     "OPEN",
		model.StatusComplete: "COMPLETE",
		model.StatusCanceled: "CANCELED",
	}

	return statusFromModel[in]
}

func statusToModel(in string) model.Status {
	statusToModel := map[string]model.Status{
		"OPEN":     model.StatusOpen,
		"COMPLETE": model.StatusComplete,
		"CANCELED": model.StatusCanceled,
	}

	return statusToModel[in]
}
