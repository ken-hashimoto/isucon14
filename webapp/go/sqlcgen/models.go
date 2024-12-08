// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package sqlcgen

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type RideStatusesStatus string

const (
	RideStatusesStatusMATCHING  RideStatusesStatus = "MATCHING"
	RideStatusesStatusENROUTE   RideStatusesStatus = "ENROUTE"
	RideStatusesStatusPICKUP    RideStatusesStatus = "PICKUP"
	RideStatusesStatusCARRYING  RideStatusesStatus = "CARRYING"
	RideStatusesStatusARRIVED   RideStatusesStatus = "ARRIVED"
	RideStatusesStatusCOMPLETED RideStatusesStatus = "COMPLETED"
)

func (e *RideStatusesStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = RideStatusesStatus(s)
	case string:
		*e = RideStatusesStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for RideStatusesStatus: %T", src)
	}
	return nil
}

type NullRideStatusesStatus struct {
	RideStatusesStatus RideStatusesStatus
	Valid              bool // Valid is true if RideStatusesStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRideStatusesStatus) Scan(value interface{}) error {
	if value == nil {
		ns.RideStatusesStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.RideStatusesStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRideStatusesStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.RideStatusesStatus), nil
}

// 椅子情報テーブル
type Chair struct {
	// 椅子ID
	ID string
	// オーナーID
	OwnerID string
	// 椅子の名前
	Name string
	// 椅子のモデル
	Model string
	// 配椅子受付中かどうか
	IsActive bool
	// アクセストークン
	AccessToken string
	// 登録日時
	CreatedAt time.Time
	// 更新日時
	UpdatedAt time.Time
}

// 椅子の現在位置情報テーブル
type ChairLocation struct {
	ID string
	// 椅子ID
	ChairID string
	// 経度
	Latitude int32
	// 緯度
	Longitude int32
	// 登録日時
	CreatedAt time.Time
}

// 椅子モデルテーブル
type ChairModel struct {
	// 椅子モデル名
	Name string
	// 移動速度
	Speed int32
}

// クーポンテーブル
type Coupon struct {
	// 所有しているユーザーのID
	UserID string
	// クーポンコード
	Code string
	// 割引額
	Discount int32
	// 付与日時
	CreatedAt time.Time
	// クーポンが適用されたライドのID
	UsedBy sql.NullString
}

// 椅子のオーナー情報テーブル
type Owner struct {
	// オーナーID
	ID string
	// オーナー名
	Name string
	// アクセストークン
	AccessToken string
	// 椅子登録トークン
	ChairRegisterToken string
	// 登録日時
	CreatedAt time.Time
	// 更新日時
	UpdatedAt time.Time
}

// 決済トークンテーブル
type PaymentToken struct {
	// ユーザーID
	UserID string
	// 決済トークン
	Token string
	// 登録日時
	CreatedAt time.Time
}

// ライド情報テーブル
type Ride struct {
	// ライドID
	ID string
	// ユーザーID
	UserID string
	// 割り当てられた椅子ID
	ChairID sql.NullString
	// 配車位置(経度)
	PickupLatitude int32
	// 配車位置(緯度)
	PickupLongitude int32
	// 目的地(経度)
	DestinationLatitude int32
	// 目的地(緯度)
	DestinationLongitude int32
	// 評価
	Evaluation sql.NullInt32
	// 要求日時
	CreatedAt time.Time
	// 状態更新日時
	UpdatedAt time.Time
}

// ライドステータスの変更履歴テーブル
type RideStatus struct {
	ID string
	// ライドID
	RideID string
	// 状態
	Status RideStatusesStatus
	// 状態変更日時
	CreatedAt time.Time
	// ユーザーへの状態通知日時
	AppSentAt sql.NullTime
	// 椅子への状態通知日時
	ChairSentAt sql.NullTime
}

// システム設定テーブル
type Setting struct {
	// 設定名
	Name string
	// 設定値
	Value string
}

// 利用者情報テーブル
type User struct {
	// ユーザーID
	ID string
	// ユーザー名
	Username string
	// 本名(名前)
	Firstname string
	// 本名(名字)
	Lastname string
	// 生年月日
	DateOfBirth string
	// アクセストークン
	AccessToken string
	// 招待トークン
	InvitationCode string
	// 登録日時
	CreatedAt time.Time
	// 更新日時
	UpdatedAt time.Time
}
