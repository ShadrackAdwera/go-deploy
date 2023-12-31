// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateLog(ctx context.Context, arg CreateLogParams) (Log, error)
	CreateTenant(ctx context.Context, arg CreateTenantParams) (Tenant, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteLog(ctx context.Context, id uuid.UUID) error
	DeleteTenant(ctx context.Context, id uuid.UUID) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetLog(ctx context.Context, id uuid.UUID) (Log, error)
	GetTenant(ctx context.Context, id uuid.UUID) (Tenant, error)
	GetUser(ctx context.Context, id uuid.UUID) (interface{}, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	ListLogs(ctx context.Context, arg ListLogsParams) (Log, error)
	ListTenants(ctx context.Context, arg ListTenantsParams) (Tenant, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]interface{}, error)
	UpdateTenant(ctx context.Context, arg UpdateTenantParams) (Tenant, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
