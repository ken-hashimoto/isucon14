// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: payments.sql

package sqlcgen

import (
	"context"
)

const createPaymentToken = `-- name: CreatePaymentToken :execrows
INSERT INTO payment_tokens (user_id, token) VALUES (?, ?)
`

type CreatePaymentTokenParams struct {
	UserID string
	Token  string
}

func (q *Queries) CreatePaymentToken(ctx context.Context, arg CreatePaymentTokenParams) (int64, error) {
	result, err := q.db.Exec(ctx, createPaymentToken, arg.UserID, arg.Token)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const updatePaymentGatewayURLValue = `-- name: UpdatePaymentGatewayURLValue :execrows
UPDATE settings SET value = ? WHERE name = 'payment_gateway_url'
`

func (q *Queries) UpdatePaymentGatewayURLValue(ctx context.Context, value string) (int64, error) {
	result, err := q.db.Exec(ctx, updatePaymentGatewayURLValue, value)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
