// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// SumRewards sum rewards
//
// swagger:model SumRewards
type SumRewards struct {

	// name
	Name string `json:"Name,omitempty"`

	// reward
	Reward string `json:"Reward,omitempty"`
}

// Validate validates this sum rewards
func (m *SumRewards) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this sum rewards based on context it is used
func (m *SumRewards) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *SumRewards) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SumRewards) UnmarshalBinary(b []byte) error {
	var res SumRewards
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}