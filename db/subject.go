package db

import (
	"errors"

	"github.com/jackc/pgtype"
	"github.com/solderneer/axiom-backend/graph/model"
)

type Subject struct {
	Name  string
	Level string
}

func (dst *Subject) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	if src == nil {
		return errors.New("NULL values can't be decoded. Scan into a &*MyType to handle NULLs")
	}

	if err := (pgtype.CompositeFields{&dst.Name, &dst.Level}).DecodeBinary(ci, src); err != nil {
		return err
	}

	return nil
}

func (src *Subject) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) (newBuf []byte, err error) {
	a := pgtype.Text{String: src.Name, Status: pgtype.Present}
	b := pgtype.Text{String: src.Level, Status: pgtype.Present}

	return (pgtype.CompositeFields{&a, &b}).EncodeBinary(ci, buf)
}

func (sb *Subject) ToSubjectModel() model.Subject {
	subject := model.Subject{Name: model.SubjectName(sb.Name), Level: model.SubjectLevel(sb.Level)}
	return subject
}
