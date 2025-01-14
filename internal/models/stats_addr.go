// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// StatsAddr stats addr
//
// swagger:model StatsAddr
type StatsAddr struct {

	// balance
	Balance string `json:"Balance,omitempty"`

	// current hash rate
	CurrentHashRate string `json:"CurrentHashRate,omitempty"`

	// current luck
	CurrentLuck string `json:"CurrentLuck,omitempty"`

	// hash rate
	HashRate string `json:"HashRate,omitempty"`

	// immature
	Immature string `json:"Immature,omitempty"`

	// last block found
	LastBlockFound int64 `json:"LastBlockFound,omitempty"`

	// payed
	Payed string `json:"Payed,omitempty"`

	// sum rewards
	SumRewards []*SumRewards `json:"SumRewards"`

	// workers
	Workers []*Worker `json:"Workers"`
}

// Validate validates this stats addr
func (m *StatsAddr) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSumRewards(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateWorkers(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *StatsAddr) validateSumRewards(formats strfmt.Registry) error {
	if swag.IsZero(m.SumRewards) { // not required
		return nil
	}

	for i := 0; i < len(m.SumRewards); i++ {
		if swag.IsZero(m.SumRewards[i]) { // not required
			continue
		}

		if m.SumRewards[i] != nil {
			if err := m.SumRewards[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("SumRewards" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *StatsAddr) validateWorkers(formats strfmt.Registry) error {
	if swag.IsZero(m.Workers) { // not required
		return nil
	}

	for i := 0; i < len(m.Workers); i++ {
		if swag.IsZero(m.Workers[i]) { // not required
			continue
		}

		if m.Workers[i] != nil {
			if err := m.Workers[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("Workers" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this stats addr based on the context it is used
func (m *StatsAddr) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateSumRewards(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateWorkers(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *StatsAddr) contextValidateSumRewards(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.SumRewards); i++ {

		if m.SumRewards[i] != nil {
			if err := m.SumRewards[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("SumRewards" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *StatsAddr) contextValidateWorkers(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Workers); i++ {

		if m.Workers[i] != nil {
			if err := m.Workers[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("Workers" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *StatsAddr) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *StatsAddr) UnmarshalBinary(b []byte) error {
	var res StatsAddr
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
