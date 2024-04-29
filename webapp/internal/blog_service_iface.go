package internal

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/sparkymat/oxgen/webapp/internal/dbx"
	"github.com/sparkymat/oxgen/webapp/internal/service/blog"
)

type BlogService interface {
	CreateUser(ctx context.Context, params blog.CreateUserParams) (dbx.User, error)
	SearchUsers(ctx context.Context, query string, pageSize int32, pageNumber int32) ([]dbx.User, int64, error)
	FetchRecentUsers(ctx context.Context, pageSize int32, pageNumber int32) ([]dbx.User, int64, error)
	FetchUser(ctx context.Context, id uuid.UUID) (dbx.User, error)
	DestroyUser(ctx context.Context, id uuid.UUID) error
	UpdateUserName(ctx context.Context, id uuid.UUID, value string) (dbx.User, error)
	UpdateUserAge(ctx context.Context, id uuid.UUID, valuePtr *int32) (dbx.User, error)
	UpdateUserDob(ctx context.Context, id uuid.UUID, valuePtr *time.Time) (dbx.User, error)
	UploadUserPhoto(ctx context.Context, id uuid.UUID, filename string, attachmentFile io.Reader) (dbx.User, error)

	CreateUser(ctx context.Context, params blog.CreateUserParams) (dbx.User, error)
	SearchUsers(ctx context.Context, query string, pageSize int32, pageNumber int32) ([]dbx.User, int64, error)
	FetchRecentUsers(ctx context.Context, pageSize int32, pageNumber int32) ([]dbx.User, int64, error)
	FetchUser(ctx context.Context, id uuid.UUID) (dbx.User, error)
	DestroyUser(ctx context.Context, id uuid.UUID) error
	UpdateUserName(ctx context.Context, id uuid.UUID, value string) (dbx.User, error)
	UpdateUserAge(ctx context.Context, id uuid.UUID, valuePtr *int32) (dbx.User, error)
	UpdateUserDob(ctx context.Context, id uuid.UUID, valuePtr *time.Time) (dbx.User, error)
	UploadUserPhoto(ctx context.Context, id uuid.UUID, filename string, attachmentFile io.Reader) (dbx.User, error)

	CreateUser(ctx context.Context, params blog.CreateUserParams) (dbx.User, error)
	SearchUsers(ctx context.Context, query string, pageSize int32, pageNumber int32) ([]dbx.User, int64, error)
	FetchRecentUsers(ctx context.Context, pageSize int32, pageNumber int32) ([]dbx.User, int64, error)
	FetchUser(ctx context.Context, id uuid.UUID) (dbx.User, error)
	DestroyUser(ctx context.Context, id uuid.UUID) error
	UpdateUserName(ctx context.Context, id uuid.UUID, value string) (dbx.User, error)
	UpdateUserAge(ctx context.Context, id uuid.UUID, valuePtr *int32) (dbx.User, error)
	UpdateUserDob(ctx context.Context, id uuid.UUID, valuePtr *time.Time) (dbx.User, error)
	UploadUserPhoto(ctx context.Context, id uuid.UUID, filename string, attachmentFile io.Reader) (dbx.User, error)

	CreateUser(ctx context.Context, params blog.CreateUserParams) (dbx.User, error)
	SearchUsers(ctx context.Context, query string, pageSize int32, pageNumber int32) ([]dbx.User, int64, error)
	FetchRecentUsers(ctx context.Context, pageSize int32, pageNumber int32) ([]dbx.User, int64, error)
	FetchUser(ctx context.Context, id uuid.UUID) (dbx.User, error)
	DestroyUser(ctx context.Context, id uuid.UUID) error
	UpdateUserName(ctx context.Context, id uuid.UUID, value string) (dbx.User, error)
	UpdateUserAge(ctx context.Context, id uuid.UUID, valuePtr *int32) (dbx.User, error)
	UpdateUserDob(ctx context.Context, id uuid.UUID, valuePtr *time.Time) (dbx.User, error)
	UploadUserPhoto(ctx context.Context, id uuid.UUID, filename string, attachmentFile io.Reader) (dbx.User, error)
}
