// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// WorkerGroupModel2 worker group model2
//
// swagger:model WorkerGroupModel2
type WorkerGroupModel2 struct {

	// bits
	Bits int64 `json:"bits,omitempty"`

	// hr
	Hr float64 `json:"hr,omitempty"`

	// hr2
	Hr2 float64 `json:"hr2,omitempty"`

	// last beat
	LastBeat int64 `json:"lastBeat,omitempty"`

	// offline
	Offline *bool `json:"offline,omitempty"`
}

// Validate validates this worker group model2
func (m *WorkerGroupModel2) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this worker group model2 based on context it is used
func (m *WorkerGroupModel2) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *WorkerGroupModel2) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *WorkerGroupModel2) UnmarshalBinary(b []byte) error {
	var res WorkerGroupModel2
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}