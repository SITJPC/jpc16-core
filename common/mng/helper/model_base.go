package mh

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModelBase struct {
	ID        *primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt *time.Time          `bson:"createdAt,omitempty"`
	UpdatedAt *time.Time          `bson:"updatedAt,omitempty"`
}

func (r *ModelBase) PrepareID(id any) (any, error) {
	if idStr, ok := id.(string); ok {
		return primitive.ObjectIDFromHex(idStr)
	}

	return id, nil
}

// GetID method returns a operand's ID
func (r *ModelBase) GetID() any {
	return r.ID
}

// SetID sets the value of a operand's ID field.
func (r *ModelBase) SetID(id any) {
	objectID := id.(primitive.ObjectID)
	r.ID = &objectID
}

// Creating hook is used here to set the `created_at` field
// value when inserting a new operand into the database.
func (r *ModelBase) Creating() error {
	utc := time.Now().UTC()
	r.CreatedAt = &utc
	return nil
}

// Saving hook is used here to set the `updated_at` field
// value when creating or updating a operand.
func (r *ModelBase) Saving() error {
	utc := time.Now().UTC()
	r.UpdatedAt = &utc
	return nil
}
