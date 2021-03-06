// Code generated by entc, DO NOT EDIT.

package item

import (
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the item type in the database.
	Label = "item"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCustomerID holds the string denoting the customer_id field in the database.
	FieldCustomerID = "customer_id"
	// FieldTitle holds the string denoting the title field in the database.
	FieldTitle = "title"
	// FieldPrice holds the string denoting the price field in the database.
	FieldPrice = "price"
	// Table holds the table name of the item in the database.
	Table = "items"
)

// Columns holds all SQL columns for item fields.
var Columns = []string{
	FieldID,
	FieldCustomerID,
	FieldTitle,
	FieldPrice,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// PriceValidator is a validator for the "price" field. It is called by the builders before save.
	PriceValidator func(int64) error
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)
