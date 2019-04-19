// Code generated by yo. DO NOT EDIT.
// Package models contains the types.
package models

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
)

// CompositePrimaryKey represents a row from 'CompositePrimaryKeys'.
type CompositePrimaryKey struct {
	ID    int64  `spanner:"Id" json:"Id"`       // Id
	PKey1 string `spanner:"PKey1" json:"PKey1"` // PKey1
	PKey2 int64  `spanner:"PKey2" json:"PKey2"` // PKey2
	Error int64  `spanner:"Error" json:"Error"` // Error
	X     string `spanner:"X" json:"X"`         // X
	Y     string `spanner:"Y" json:"Y"`         // Y
	Z     string `spanner:"Z" json:"Z"`         // Z
}

func CompositePrimaryKeyPrimaryKeys() []string {
	return []string{
		"PKey1",
		"PKey2",
	}
}

func CompositePrimaryKeyColumns() []string {
	return []string{
		"Id",
		"PKey1",
		"PKey2",
		"Error",
		"X",
		"Y",
		"Z",
	}
}

func (cpk *CompositePrimaryKey) columnsToPtrs(cols []string, customPtrs map[string]interface{}) ([]interface{}, error) {
	ret := make([]interface{}, 0, len(cols))
	for _, col := range cols {
		if val, ok := customPtrs[col]; ok {
			ret = append(ret, val)
			continue
		}

		switch col {
		case "Id":
			ret = append(ret, &cpk.ID)
		case "PKey1":
			ret = append(ret, &cpk.PKey1)
		case "PKey2":
			ret = append(ret, &cpk.PKey2)
		case "Error":
			ret = append(ret, &cpk.Error)
		case "X":
			ret = append(ret, &cpk.X)
		case "Y":
			ret = append(ret, &cpk.Y)
		case "Z":
			ret = append(ret, &cpk.Z)
		default:
			return nil, fmt.Errorf("unknown column: %s", col)
		}
	}
	return ret, nil
}

func (cpk *CompositePrimaryKey) columnsToValues(cols []string) ([]interface{}, error) {
	ret := make([]interface{}, 0, len(cols))
	for _, col := range cols {
		switch col {
		case "Id":
			ret = append(ret, cpk.ID)
		case "PKey1":
			ret = append(ret, cpk.PKey1)
		case "PKey2":
			ret = append(ret, cpk.PKey2)
		case "Error":
			ret = append(ret, cpk.Error)
		case "X":
			ret = append(ret, cpk.X)
		case "Y":
			ret = append(ret, cpk.Y)
		case "Z":
			ret = append(ret, cpk.Z)
		default:
			return nil, fmt.Errorf("unknown column: %s", col)
		}
	}

	return ret, nil
}

// newCompositePrimaryKey_Decoder returns a decoder which reads a row from *spanner.Row
// into CompositePrimaryKey. The decoder is not goroutine-safe. Don't use it concurrently.
func newCompositePrimaryKey_Decoder(cols []string) func(*spanner.Row) (*CompositePrimaryKey, error) {
	customPtrs := map[string]interface{}{}

	return func(row *spanner.Row) (*CompositePrimaryKey, error) {
		var cpk CompositePrimaryKey
		ptrs, err := cpk.columnsToPtrs(cols, customPtrs)
		if err != nil {
			return nil, err
		}

		if err := row.Columns(ptrs...); err != nil {
			return nil, err
		}

		return &cpk, nil
	}
}

// Insert returns a Mutation to insert a row into a table. If the row already
// exists, the write or transaction fails.
func (cpk *CompositePrimaryKey) Insert(ctx context.Context) *spanner.Mutation {
	return spanner.Insert("CompositePrimaryKeys", CompositePrimaryKeyColumns(), []interface{}{
		cpk.ID, cpk.PKey1, cpk.PKey2, cpk.Error, cpk.X, cpk.Y, cpk.Z,
	})
}

// InsertCompositePrimaryKeyAll returns slice of Mutation to insert rows into a table. If the row already
// exists, the write or transaction fails.
func InsertCompositePrimaryKeyAll(ctx context.Context, rows []*CompositePrimaryKey) []*spanner.Mutation {
	muts := make([]*spanner.Mutation, 0, len(rows))
	for _, r := range rows {
		muts = append(muts, r.Insert(ctx))
	}
	return muts
}

// Update returns a Mutation to update a row in a table. If the row does not
// already exist, the write or transaction fails.
func (cpk *CompositePrimaryKey) Update(ctx context.Context) *spanner.Mutation {
	return spanner.Update("CompositePrimaryKeys", CompositePrimaryKeyColumns(), []interface{}{
		cpk.ID, cpk.PKey1, cpk.PKey2, cpk.Error, cpk.X, cpk.Y, cpk.Z,
	})
}

// UpdateCompositePrimaryKeyAll returns slice of Mutation to update rows in a table. If the row does not
// already exist, the write or transaction fails.
func UpdateCompositePrimaryKeyAll(ctx context.Context, rows []*CompositePrimaryKey) []*spanner.Mutation {
	muts := make([]*spanner.Mutation, 0, len(rows))
	for _, r := range rows {
		muts = append(muts, r.Update(ctx))
	}
	return muts
}

// InsertOrUpdate returns a Mutation to insert a row into a table. If the row
// already exists, it updates it instead. Any column values not explicitly
// written are preserved.
func (cpk *CompositePrimaryKey) InsertOrUpdate(ctx context.Context) *spanner.Mutation {
	return spanner.InsertOrUpdate("CompositePrimaryKeys", CompositePrimaryKeyColumns(), []interface{}{
		cpk.ID, cpk.PKey1, cpk.PKey2, cpk.Error, cpk.X, cpk.Y, cpk.Z,
	})
}

// InsertOrUpdateCompositePrimaryKeyAll returns slice of Mutation to insert rows into a table. If the row
// already exists, it updates it instead. Any column values not explicitly
// written are preserved.
func InsertOrUpdateCompositePrimaryKeyAll(ctx context.Context, rows []*CompositePrimaryKey) []*spanner.Mutation {
	muts := make([]*spanner.Mutation, 0, len(rows))
	for _, r := range rows {
		muts = append(muts, r.InsertOrUpdate(ctx))
	}
	return muts
}

// UpdateColumns returns a Mutation to update specified columns of a row in a table.
func (cpk *CompositePrimaryKey) UpdateColumns(ctx context.Context, cols ...string) (*spanner.Mutation, error) {
	// add primary keys to columns to update by primary keys
	colsWithPKeys := append(cols, CompositePrimaryKeyPrimaryKeys()...)

	values, err := cpk.columnsToValues(colsWithPKeys)
	if err != nil {
		return nil, newErrorWithCode(codes.InvalidArgument, "CompositePrimaryKey.UpdateColumns", "CompositePrimaryKeys", err)
	}

	return spanner.Update("CompositePrimaryKeys", colsWithPKeys, values), nil
}

// UpdateCompositePrimaryKeyColumnsAll returns slice of Mutation to update specified columns of rows in a table.
func UpdateCompositePrimaryKeyColumnsAll(ctx context.Context, rows []*CompositePrimaryKey, cols ...string) ([]*spanner.Mutation, error) {
	// add primary keys to columns to update by primary keys
	colsWithPKeys := append(cols, CompositePrimaryKeyPrimaryKeys()...)

	muts := make([]*spanner.Mutation, 0, len(rows))
	for _, r := range rows {
		values, err := r.columnsToValues(colsWithPKeys)
		if err != nil {
			return nil, newErrorWithCode(codes.InvalidArgument, "CompositePrimaryKey.UpdateColumns", "CompositePrimaryKeys", err)
		}

		muts = append(muts, spanner.Update("CompositePrimaryKeys", colsWithPKeys, values))
	}

	return muts, nil
}

// FindCompositePrimaryKey gets a CompositePrimaryKey by primary key
func FindCompositePrimaryKey(ctx context.Context, db YORODB, pKey1 string, pKey2 int64) (*CompositePrimaryKey, error) {
	key := spanner.Key{pKey1, pKey2}
	row, err := db.ReadRow(ctx, "CompositePrimaryKeys", key, CompositePrimaryKeyColumns())
	if err != nil {
		return nil, newError("FindCompositePrimaryKey", "CompositePrimaryKeys", err)
	}

	decoder := newCompositePrimaryKey_Decoder(CompositePrimaryKeyColumns())
	cpk, err := decoder(row)
	if err != nil {
		return nil, newErrorWithCode(codes.Internal, "FindCompositePrimaryKey", "CompositePrimaryKeys", err)
	}

	return cpk, nil
}

// Delete deletes the CompositePrimaryKey from the database.
func (cpk *CompositePrimaryKey) Delete(ctx context.Context) *spanner.Mutation {
	values, _ := cpk.columnsToValues(CompositePrimaryKeyPrimaryKeys())
	return spanner.Delete("CompositePrimaryKeys", spanner.Key(values))
}

// DeleteCompositePrimaryKeyAll deletes the CompositePrimaryKey rows from the database.
func DeleteCompositePrimaryKeyAll(ctx context.Context, rows []*CompositePrimaryKey) []*spanner.Mutation {
	muts := make([]*spanner.Mutation, 0, len(rows))
	for _, r := range rows {
		muts = append(muts, r.Delete(ctx))
	}
	return muts
}

// FindCompositePrimaryKeysByError retrieves multiple rows from 'CompositePrimaryKeys' as a slice of CompositePrimaryKey.
//
// Generated from index 'CompositePrimaryKeysByError'.
func FindCompositePrimaryKeysByError(ctx context.Context, db YORODB, e int64) ([]*CompositePrimaryKey, error) {
	const sqlstr = `SELECT ` +
		`Id, PKey1, PKey2, Error, X, Y, Z ` +
		`FROM CompositePrimaryKeys@{FORCE_INDEX=CompositePrimaryKeysByError} ` +
		`WHERE Error = @param0`

	stmt := spanner.NewStatement(sqlstr)
	stmt.Params["param0"] = e

	decoder := newCompositePrimaryKey_Decoder(CompositePrimaryKeyColumns())

	// run query
	YOLog(ctx, sqlstr, e)
	iter := db.Query(ctx, stmt)
	defer iter.Stop()

	// load results
	res := []*CompositePrimaryKey{}
	for {
		row, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, newError("FindCompositePrimaryKeysByError", "CompositePrimaryKeys", err)
		}

		cpk, err := decoder(row)
		if err != nil {
			return nil, newErrorWithCode(codes.Internal, "FindCompositePrimaryKeysByError", "CompositePrimaryKeys", err)
		}

		res = append(res, cpk)
	}

	return res, nil
}

// FindCompositePrimaryKeysByXY retrieves multiple rows from 'CompositePrimaryKeys' as a slice of CompositePrimaryKey.
//
// Generated from index 'CompositePrimaryKeysByXY'.
func FindCompositePrimaryKeysByXY(ctx context.Context, db YORODB, x string, y string) ([]*CompositePrimaryKey, error) {
	const sqlstr = `SELECT ` +
		`Id, PKey1, PKey2, Error, X, Y, Z ` +
		`FROM CompositePrimaryKeys@{FORCE_INDEX=CompositePrimaryKeysByXY} ` +
		`WHERE X = @param0 AND Y = @param1`

	stmt := spanner.NewStatement(sqlstr)
	stmt.Params["param0"] = x
	stmt.Params["param1"] = y

	decoder := newCompositePrimaryKey_Decoder(CompositePrimaryKeyColumns())

	// run query
	YOLog(ctx, sqlstr, x, y)
	iter := db.Query(ctx, stmt)
	defer iter.Stop()

	// load results
	res := []*CompositePrimaryKey{}
	for {
		row, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, newError("FindCompositePrimaryKeysByXY", "CompositePrimaryKeys", err)
		}

		cpk, err := decoder(row)
		if err != nil {
			return nil, newErrorWithCode(codes.Internal, "FindCompositePrimaryKeysByXY", "CompositePrimaryKeys", err)
		}

		res = append(res, cpk)
	}

	return res, nil
}
