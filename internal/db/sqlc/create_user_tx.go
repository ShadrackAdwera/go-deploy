package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateTenantInput struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	TenantName string `json:"tenant_name"`
	Logo       string `json:"logo"`
	Password   string `json:"password"`
}

type CreateTenantOutput struct {
	Message string
}

func (s *Store) CreateTenantTx(ctx context.Context, input CreateTenantInput) (CreateTenantOutput, error) {
	err := s.execTx(ctx, func(q *Queries) error {
		tenant, err := q.CreateTenant(ctx, CreateTenantParams{
			Name: input.TenantName,
			Logo: pgtype.Text{
				String: input.Logo,
				Valid:  true,
			},
		})

		if err != nil {
			return err
		}

		_, err = q.CreateUser(ctx, CreateUserParams{
			Username: input.Username,
			Email:    input.Email,
			TenantID: tenant.ID,
			Password: input.Password,
		})

		if err != nil {
			return err
		}
		// emit to queue - AfterCreate callback function
		return nil
	})

	if err != nil {
		return CreateTenantOutput{}, err
	}

	return CreateTenantOutput{
		Message: "Tenant Successfully Created",
	}, nil
}
