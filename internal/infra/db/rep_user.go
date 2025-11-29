package db

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"go-ai/internal/domain/user"
// 	"go-ai/internal/infra/sqlc"
// 	"go-ai/pkg/common"

// 	"github.com/google/uuid"
// 	"github.com/jackc/pgconn"
// 	"github.com/jackc/pgerrcode"
// 	"github.com/jackc/pgx/v5"
// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// type UserRepo struct {
// 	q *sqlc.Queries
// }

// func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
// 	return &UserRepo{q: sqlc.New(pool)}
// }

// func (r *UserRepo) GetByID(id uuid.UUID) (user.User, error) {
// 	u, err := r.q.GetUserByID(context.Background(), id)
// 	if errors.Is(err, pgx.ErrNoRows) {
// 		return user.User{}, common.ErrNotFound
// 	}
// 	if err != nil {
// 		return user.User{}, err
// 	}
// 	return user.User{ID: u.ID, Email: *u.Email, FullName: u.FullName}, nil
// }

// func (r *UserRepo) GetByEmail(email string) (user.User, error) {
// 	u, err := r.q.GetUserByEmail(context.Background(), &email)
// 	if errors.Is(err, pgx.ErrNoRows) {
// 		return user.User{}, common.ErrNotFound
// 	}
// 	if err != nil {
// 		return user.User{}, err
// 	}
// 	return user.User{ID: u.ID, Email: *u.Email, FullName: u.FullName}, nil
// }

// func (r *UserRepo) Create(in user.User) (uuid.UUID, error) {
// 	id, err := r.q.CreateUser(context.Background(), sqlc.CreateUserParams{
// 		Email:        &in.Email,
// 		FullName:     in.FullName,
// 		PasswordHash: in.Password,
// 	})
// 	if err != nil {
// 		// TODO: check unique_violation via pgerrcode
// 		var pgErr *pgconn.PgError
// 		if errors.As(err, &pgErr) {
// 			if pgErr.Code == pgerrcode.UniqueViolation {
// 				return uuid.Nil, fmt.Errorf("user already exists")
// 			}
// 		}
// 		return uuid.Nil, err
// 	}
// 	return id, nil
// }
