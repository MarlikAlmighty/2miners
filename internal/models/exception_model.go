// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ExceptionModel exception model
//
// swagger:model ExceptionModel
type ExceptionModel struct {

	// code
	Code int64 `json:"code,omitempty"`

	// description
	Description string `json:"description,omitempty"`

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this exception model
func (m *ExceptionModel) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this exception model based on context it is used
func (m *ExceptionModel) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ExceptionModel) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ExceptionModel) UnmarshalBinary(b []byte) error {
	var res ExceptionModel
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
