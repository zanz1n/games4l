package sqli

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type QuestionStyle string

const (
	QuestionStyleImage QuestionStyle = "image"
	QuestionStyleAudio QuestionStyle = "audio"
	QuestionStyleVideo QuestionStyle = "video"
	QuestionStyleText  QuestionStyle = "text"
)

func (e *QuestionStyle) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = QuestionStyle(s)
	case string:
		*e = QuestionStyle(s)
	default:
		return fmt.Errorf("unsupported scan type for QuestionStyle: %T", src)
	}
	return nil
}

type NullQuestionStyle struct {
	QuestionStyle QuestionStyle
	Valid         bool // Valid is true if QuestionStyle is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullQuestionStyle) Scan(value interface{}) error {
	if value == nil {
		ns.QuestionStyle, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.QuestionStyle.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullQuestionStyle) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.QuestionStyle), nil
}

type QuestionType string

const (
	QuestionType2Alt QuestionType = "2Alt"
	QuestionType4Alt QuestionType = "4Alt"
)

func (e *QuestionType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = QuestionType(s)
	case string:
		*e = QuestionType(s)
	default:
		return fmt.Errorf("unsupported scan type for QuestionType: %T", src)
	}
	return nil
}

type NullQuestionType struct {
	QuestionType QuestionType
	Valid        bool // Valid is true if QuestionType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullQuestionType) Scan(value interface{}) error {
	if value == nil {
		ns.QuestionType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.QuestionType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullQuestionType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.QuestionType), nil
}

type Question struct {
	ID            int32
	CreatedAt     pgtype.Timestamp
	UpdatedAt     pgtype.Timestamp
	Question      string
	Answer1       string
	Answer2       string
	Answer3       pgtype.Text
	Answer4       pgtype.Text
	CorrectAnswer int16
	Type          QuestionType
	Style         QuestionStyle
	Difficulty    int16
}
