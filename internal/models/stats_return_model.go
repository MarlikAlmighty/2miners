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

// StatsReturnModel stats return model
//
// swagger:model StatsReturnModel
type StatsReturnModel struct {

	// candidates total
	CandidatesTotal int64 `json:"candidatesTotal,omitempty"`

	// hashrate
	Hashrate float32 `json:"hashrate,omitempty"`

	// immature total
	ImmatureTotal int64 `json:"immatureTotal,omitempty"`

	// luck
	Luck float32 `json:"luck,omitempty"`

	// matured total
	MaturedTotal int64 `json:"maturedTotal,omitempty"`

	// miners total
	MinersTotal int64 `json:"minersTotal,omitempty"`

	// nodes
	Nodes []*NodeModel `json:"nodes"`

	// now
	Now int64 `json:"now,omitempty"`

	// payments total
	PaymentsTotal int64 `json:"paymentsTotal,omitempty"`

	// pool charts
	PoolCharts []*PoolChartsModel `json:"poolCharts"`

	// stats
	Stats *StatsReturnModelStats `json:"stats,omitempty"`
}

// Validate validates this stats return model
func (m *StatsReturnModel) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateNodes(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePoolCharts(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStats(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *StatsReturnModel) validateNodes(formats strfmt.Registry) error {
	if swag.IsZero(m.Nodes) { // not required
		return nil
	}

	for i := 0; i < len(m.Nodes); i++ {
		if swag.IsZero(m.Nodes[i]) { // not required
			continue
		}

		if m.Nodes[i] != nil {
			if err := m.Nodes[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("nodes" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *StatsReturnModel) validatePoolCharts(formats strfmt.Registry) error {
	if swag.IsZero(m.PoolCharts) { // not required
		return nil
	}

	for i := 0; i < len(m.PoolCharts); i++ {
		if swag.IsZero(m.PoolCharts[i]) { // not required
			continue
		}

		if m.PoolCharts[i] != nil {
			if err := m.PoolCharts[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("poolCharts" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *StatsReturnModel) validateStats(formats strfmt.Registry) error {
	if swag.IsZero(m.Stats) { // not required
		return nil
	}

	if m.Stats != nil {
		if err := m.Stats.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("stats")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this stats return model based on the context it is used
func (m *StatsReturnModel) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateNodes(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidatePoolCharts(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateStats(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *StatsReturnModel) contextValidateNodes(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Nodes); i++ {

		if m.Nodes[i] != nil {
			if err := m.Nodes[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("nodes" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *StatsReturnModel) contextValidatePoolCharts(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.PoolCharts); i++ {

		if m.PoolCharts[i] != nil {
			if err := m.PoolCharts[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("poolCharts" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *StatsReturnModel) contextValidateStats(ctx context.Context, formats strfmt.Registry) error {

	if m.Stats != nil {
		if err := m.Stats.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("stats")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *StatsReturnModel) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *StatsReturnModel) UnmarshalBinary(b []byte) error {
	var res StatsReturnModel
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// StatsReturnModelStats stats return model stats
//
// swagger:model StatsReturnModelStats
type StatsReturnModelStats struct {

	// last block found
	LastBlockFound int64 `json:"lastBlockFound,omitempty"`

	// nshares
	Nshares int64 `json:"nshares,omitempty"`

	// round shares
	RoundShares float32 `json:"roundShares,omitempty"`
}

// Validate validates this stats return model stats
func (m *StatsReturnModelStats) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this stats return model stats based on context it is used
func (m *StatsReturnModelStats) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *StatsReturnModelStats) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *StatsReturnModelStats) UnmarshalBinary(b []byte) error {
	var res StatsReturnModelStats
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}