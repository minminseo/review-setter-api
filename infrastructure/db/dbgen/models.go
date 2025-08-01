// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package dbgen

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type TargetWeightEnum string

const (
	TargetWeightEnumHeavy  TargetWeightEnum = "heavy"
	TargetWeightEnumNormal TargetWeightEnum = "normal"
	TargetWeightEnumLight  TargetWeightEnum = "light"
	TargetWeightEnumUnset  TargetWeightEnum = "unset"
)

func (e *TargetWeightEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TargetWeightEnum(s)
	case string:
		*e = TargetWeightEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for TargetWeightEnum: %T", src)
	}
	return nil
}

type NullTargetWeightEnum struct {
	TargetWeightEnum TargetWeightEnum `json:"target_weight_enum"`
	Valid            bool             `json:"valid"` // Valid is true if TargetWeightEnum is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTargetWeightEnum) Scan(value interface{}) error {
	if value == nil {
		ns.TargetWeightEnum, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TargetWeightEnum.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTargetWeightEnum) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TargetWeightEnum), nil
}

type ThemeColorEnum string

const (
	ThemeColorEnumDark  ThemeColorEnum = "dark"
	ThemeColorEnumLight ThemeColorEnum = "light"
)

func (e *ThemeColorEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ThemeColorEnum(s)
	case string:
		*e = ThemeColorEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for ThemeColorEnum: %T", src)
	}
	return nil
}

type NullThemeColorEnum struct {
	ThemeColorEnum ThemeColorEnum `json:"theme_color_enum"`
	Valid          bool           `json:"valid"` // Valid is true if ThemeColorEnum is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullThemeColorEnum) Scan(value interface{}) error {
	if value == nil {
		ns.ThemeColorEnum, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ThemeColorEnum.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullThemeColorEnum) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ThemeColorEnum), nil
}

type Category struct {
	ID           pgtype.UUID        `json:"id"`
	UserID       pgtype.UUID        `json:"user_id"`
	Name         string             `json:"name"`
	RegisteredAt pgtype.Timestamptz `json:"registered_at"`
	EditedAt     pgtype.Timestamptz `json:"edited_at"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
}

type EmailVerification struct {
	ID        pgtype.UUID        `json:"id"`
	UserID    pgtype.UUID        `json:"user_id"`
	CodeHash  string             `json:"code_hash"`
	ExpiresAt pgtype.Timestamptz `json:"expires_at"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

type PatternStep struct {
	ID           pgtype.UUID        `json:"id"`
	UserID       pgtype.UUID        `json:"user_id"`
	PatternID    pgtype.UUID        `json:"pattern_id"`
	StepNumber   int16              `json:"step_number"`
	IntervalDays int16              `json:"interval_days"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
}

type ReviewBox struct {
	ID           pgtype.UUID        `json:"id"`
	UserID       pgtype.UUID        `json:"user_id"`
	CategoryID   pgtype.UUID        `json:"category_id"`
	PatternID    pgtype.UUID        `json:"pattern_id"`
	Name         string             `json:"name"`
	RegisteredAt pgtype.Timestamptz `json:"registered_at"`
	EditedAt     pgtype.Timestamptz `json:"edited_at"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
}

type ReviewDate struct {
	ID                   pgtype.UUID        `json:"id"`
	UserID               pgtype.UUID        `json:"user_id"`
	CategoryID           pgtype.UUID        `json:"category_id"`
	BoxID                pgtype.UUID        `json:"box_id"`
	ItemID               pgtype.UUID        `json:"item_id"`
	StepNumber           int16              `json:"step_number"`
	InitialScheduledDate pgtype.Date        `json:"initial_scheduled_date"`
	ScheduledDate        pgtype.Date        `json:"scheduled_date"`
	IsCompleted          bool               `json:"is_completed"`
	CreatedAt            pgtype.Timestamptz `json:"created_at"`
	UpdatedAt            pgtype.Timestamptz `json:"updated_at"`
}

type ReviewItem struct {
	ID           pgtype.UUID        `json:"id"`
	UserID       pgtype.UUID        `json:"user_id"`
	CategoryID   pgtype.UUID        `json:"category_id"`
	BoxID        pgtype.UUID        `json:"box_id"`
	PatternID    pgtype.UUID        `json:"pattern_id"`
	Name         string             `json:"name"`
	Detail       pgtype.Text        `json:"detail"`
	LearnedDate  pgtype.Date        `json:"learned_date"`
	IsFinished   bool               `json:"is_finished"`
	RegisteredAt pgtype.Timestamptz `json:"registered_at"`
	EditedAt     pgtype.Timestamptz `json:"edited_at"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
}

type ReviewPattern struct {
	ID           pgtype.UUID        `json:"id"`
	UserID       pgtype.UUID        `json:"user_id"`
	Name         string             `json:"name"`
	TargetWeight TargetWeightEnum   `json:"target_weight"`
	RegisteredAt pgtype.Timestamptz `json:"registered_at"`
	EditedAt     pgtype.Timestamptz `json:"edited_at"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
}

type User struct {
	ID             pgtype.UUID        `json:"id"`
	EmailSearchKey string             `json:"email_search_key"`
	Email          string             `json:"email"`
	Password       string             `json:"password"`
	Timezone       string             `json:"timezone"`
	ThemeColor     ThemeColorEnum     `json:"theme_color"`
	Language       string             `json:"language"`
	VerifiedAt     pgtype.Timestamptz `json:"verified_at"`
	CreatedAt      pgtype.Timestamptz `json:"created_at"`
	UpdatedAt      pgtype.Timestamptz `json:"updated_at"`
}
