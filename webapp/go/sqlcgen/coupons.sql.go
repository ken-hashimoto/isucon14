// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: coupons.sql

package sqlcgen

import (
	"context"
	"database/sql"
)

const createCoupons = `-- name: CreateCoupons :execrows
INSERT INTO coupons (user_id, code, discount) VALUES (?, ?, ?)
`

type CreateCouponsParams struct {
	UserID   string
	Code     string
	Discount int32
}

func (q *Queries) CreateCoupons(ctx context.Context, arg CreateCouponsParams) (int64, error) {
	result, err := q.db.Exec(ctx, createCoupons, arg.UserID, arg.Code, arg.Discount)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const createCouponsWihtCalc = `-- name: CreateCouponsWihtCalc :execrows
INSERT INTO coupons (user_id, code, discount) VALUES (?, CONCAT(?, '_', FLOOR(UNIX_TIMESTAMP(NOW(3))*1000)), ?)
`

type CreateCouponsWihtCalcParams struct {
	UserID   string
	CONCAT   interface{}
	Discount int32
}

func (q *Queries) CreateCouponsWihtCalc(ctx context.Context, arg CreateCouponsWihtCalcParams) (int64, error) {
	result, err := q.db.Exec(ctx, createCouponsWihtCalc, arg.UserID, arg.CONCAT, arg.Discount)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getCouponByUsedBy = `-- name: GetCouponByUsedBy :one
SELECT user_id, code, discount, created_at, used_by FROM coupons WHERE used_by = ? LIMIT 1
`

func (q *Queries) GetCouponByUsedBy(ctx context.Context, usedBy sql.NullString) (Coupon, error) {
	row := q.db.QueryRow(ctx, getCouponByUsedBy, usedBy)
	var i Coupon
	err := row.Scan(
		&i.UserID,
		&i.Code,
		&i.Discount,
		&i.CreatedAt,
		&i.UsedBy,
	)
	return i, err
}

const getCouponsByCode = `-- name: GetCouponsByCode :many
SELECT user_id, code, discount, created_at, used_by FROM coupons WHERE code = ? FOR UPDATE
`

func (q *Queries) GetCouponsByCode(ctx context.Context, code string) ([]Coupon, error) {
	rows, err := q.db.Query(ctx, getCouponsByCode, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Coupon
	for rows.Next() {
		var i Coupon
		if err := rows.Scan(
			&i.UserID,
			&i.Code,
			&i.Discount,
			&i.CreatedAt,
			&i.UsedBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCouponsByUserIDAndCode = `-- name: GetCouponsByUserIDAndCode :one
SELECT user_id, code, discount, created_at, used_by FROM coupons WHERE user_id = ? AND code = ? AND used_by IS NULL LIMIT 1
`

type GetCouponsByUserIDAndCodeParams struct {
	UserID string
	Code   string
}

func (q *Queries) GetCouponsByUserIDAndCode(ctx context.Context, arg GetCouponsByUserIDAndCodeParams) (Coupon, error) {
	row := q.db.QueryRow(ctx, getCouponsByUserIDAndCode, arg.UserID, arg.Code)
	var i Coupon
	err := row.Scan(
		&i.UserID,
		&i.Code,
		&i.Discount,
		&i.CreatedAt,
		&i.UsedBy,
	)
	return i, err
}

const getOldestUnusedCouponByUserID = `-- name: GetOldestUnusedCouponByUserID :one
SELECT user_id, code, discount, created_at, used_by FROM coupons WHERE user_id = ? AND used_by IS NULL ORDER BY created_at LIMIT 1
`

func (q *Queries) GetOldestUnusedCouponByUserID(ctx context.Context, userID string) (Coupon, error) {
	row := q.db.QueryRow(ctx, getOldestUnusedCouponByUserID, userID)
	var i Coupon
	err := row.Scan(
		&i.UserID,
		&i.Code,
		&i.Discount,
		&i.CreatedAt,
		&i.UsedBy,
	)
	return i, err
}

const getUnusedCouponByUserIDAndCode = `-- name: GetUnusedCouponByUserIDAndCode :one
SELECT user_id, code, discount, created_at, used_by FROM coupons WHERE user_id = ? AND code = ? AND used_by IS NULL FOR UPDATE
`

type GetUnusedCouponByUserIDAndCodeParams struct {
	UserID string
	Code   string
}

func (q *Queries) GetUnusedCouponByUserIDAndCode(ctx context.Context, arg GetUnusedCouponByUserIDAndCodeParams) (Coupon, error) {
	row := q.db.QueryRow(ctx, getUnusedCouponByUserIDAndCode, arg.UserID, arg.Code)
	var i Coupon
	err := row.Scan(
		&i.UserID,
		&i.Code,
		&i.Discount,
		&i.CreatedAt,
		&i.UsedBy,
	)
	return i, err
}
