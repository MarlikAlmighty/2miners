// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// WorkerGroupModel worker group model
//
// swagger:model WorkerGroupModel
type WorkerGroupModel struct {

	// hr
	Hr float32 `json:"hr,omitempty"`

	// hr2
	Hr2 float32 `json:"hr2,omitempty"`

	// last beat
	LastBeat int64 `json:"lastBeat,omitempty"`

	// offline
	Offline *bool `json:"offline,omitempty"`
}

// Validate validates this worker group model
func (m *WorkerGroupModel) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this worker group model based on context it is used
func (m *WorkerGroupModel) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *WorkerGroupModel) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *WorkerGroupModel) UnmarshalBinary(b []byte) error {
	var res WorkerGroupModel
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
