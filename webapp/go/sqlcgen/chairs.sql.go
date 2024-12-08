// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: chairs.sql

package sqlcgen

import (
	"context"
)

const createChairs = `-- name: CreateChairs :execrows
INSERT INTO chairs (id, owner_id, name, model, is_active, access_token) VALUES (?, ?, ?, ?, ?, ?)
`

type CreateChairsParams struct {
	ID          string
	OwnerID     string
	Name        string
	Model       string
	IsActive    bool
	AccessToken string
}

func (q *Queries) CreateChairs(ctx context.Context, arg CreateChairsParams) (int64, error) {
	result, err := q.db.Exec(ctx, createChairs,
		arg.ID,
		arg.OwnerID,
		arg.Name,
		arg.Model,
		arg.IsActive,
		arg.AccessToken,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
