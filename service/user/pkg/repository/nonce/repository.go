package nonce

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"

	"github.com/trenddapp/backend/service/user/pkg/model"
)

var (
	ErrNonceNotFound = errors.New("nonce not found")
)

type repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) Repository {
	return &repository{db: db}
}

func (r repository) GetNonceByUserID(ctx context.Context, userID string) (model.Nonce, error) {
	var out nonce

	err := r.db.
		NewSelect().
		Model(&out).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Nonce{}, ErrNonceNotFound
		}

		return model.Nonce{}, err
	}

	return convertToModel(out), nil
}

func (r repository) CreateNonce(ctx context.Context, in model.Nonce) (model.Nonce, error) {
	out := convertFromModel(in)

	_, err := r.db.
		NewInsert().
		Model(&out).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return model.Nonce{}, err
	}

	return convertToModel(out), nil
}

func (r repository) DeleteNonce(ctx context.Context, id string) (model.Nonce, error) {
	var out nonce

	result, err := r.db.
		NewDelete().
		Model(&out).
		Where("id = ?", id).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return model.Nonce{}, err
	}

	rowsEffected, err := result.RowsAffected()
	if err != nil {
		return model.Nonce{}, err
	}

	if rowsEffected == 0 {
		return model.Nonce{}, ErrNonceNotFound
	}

	return convertToModel(out), nil
}

type nonce struct {
	bun.BaseModel `bun:"nonces"`

	ID        uuid.UUID `bun:",pk,default:gen_random_uuid()"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UserID    uuid.UUID `bun:",notnull"`
	Value     string    `bun:",notnull"`
}

func convertFromModel(in model.Nonce) nonce {
	return nonce{
		UserID: uuid.FromStringOrNil(in.UserID),
		Value:  in.Value,
	}
}

func convertToModel(in nonce) model.Nonce {
	return model.Nonce{
		ID:        in.ID.String(),
		CreatedAt: in.CreatedAt,
		UserID:    in.UserID.String(),
		Value:     in.Value,
	}
}
