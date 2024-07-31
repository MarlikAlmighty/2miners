// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// MaturedModel matured model
//
// swagger:model MaturedModel
type MaturedModel struct {

	// difficulty
	Difficulty float32 `json:"difficulty,omitempty"`

	// finder
	Finder string `json:"finder,omitempty"`

	// hash
	Hash string `json:"hash,omitempty"`

	// height
	Height int64 `json:"height,omitempty"`

	// orphan
	Orphan *bool `json:"orphan,omitempty"`

	// reward
	Reward int64 `json:"reward,omitempty"`

	// shares
	Shares float32 `json:"shares,omitempty"`

	// timestamp
	Timestamp int64 `json:"timestamp,omitempty"`

	// uncle
	Uncle *bool `json:"uncle,omitempty"`

	// uncle height
	UncleHeight int64 `json:"uncleHeight,omitempty"`
}

// Validate validates this matured model
func (m *MaturedModel) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this matured model based on context it is used
func (m *MaturedModel) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *MaturedModel) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MaturedModel) UnmarshalBinary(b []byte) error {
	var res MaturedModel
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
