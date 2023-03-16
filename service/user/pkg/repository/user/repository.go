package user

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
	ErrUserNotFound = errors.New("user not found")
)

type repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) Repository {
	return &repository{db: db}
}

func (r repository) GetUser(ctx context.Context, id string) (model.User, error) {
	var out user

	err := r.db.
		NewSelect().
		Model(&out).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}

		return model.User{}, err
	}

	return convertToModel(out), nil
}

func (r repository) GetUserByAddress(ctx context.Context, address string) (model.User, error) {
	var out user

	err := r.db.
		NewSelect().
		Model(&out).
		Where("address = ?", address).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}

		return model.User{}, err
	}

	return convertToModel(out), nil
}

func (r repository) CreateUser(ctx context.Context, in model.User) (model.User, error) {
	out := convertFromModel(in)

	_, err := r.db.
		NewInsert().
		Model(&out).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return model.User{}, err
	}

	return convertToModel(out), nil
}

func (r repository) UpdateUser(ctx context.Context, in model.User) (model.User, error) {
	out := convertFromModel(in)

	result, err := r.db.
		NewUpdate().
		Model(&out).
		OmitZero().
		Where("id = ? AND updated_at = ?", in.ID, in.UpdatedAt).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return model.User{}, err
	}

	rowsEffected, err := result.RowsAffected()
	if err != nil {
		return model.User{}, err
	}

	if rowsEffected == 0 {
		return model.User{}, ErrUserNotFound
	}

	return convertToModel(out), nil
}

func (r repository) DeleteUser(ctx context.Context, id string) (model.User, error) {
	var out user

	result, err := r.db.
		NewDelete().
		Model(&out).
		Where("id = ?", id).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return model.User{}, err
	}

	rowsEffected, err := result.RowsAffected()
	if err != nil {
		return model.User{}, err
	}

	if rowsEffected == 0 {
		return model.User{}, ErrUserNotFound
	}

	return convertToModel(out), nil
}

type user struct {
	bun.BaseModel `bun:"users"`

	ID        uuid.UUID `bun:",pk,default:gen_random_uuid()"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	Address   string    `bun:",unique,notnull"`
	Balance   int64     `bun:",notnull"`
}

func convertFromModel(in model.User) user {
	return user{
		Address: in.Address,
	}
}

func convertToModel(in user) model.User {
	return model.User{
		ID:        in.ID.String(),
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
		Address:   in.Address,
		Balance:   in.Balance,
	}
}
