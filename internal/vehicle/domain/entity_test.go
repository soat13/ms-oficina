package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVehicle(t *testing.T) {
	t.Run("should create a new vehicle", func(t *testing.T) {
		v, err := NewVehicle("customer-id", "plate", "model", "brand", 2023)
		assert.NoError(t, err)
		assert.NotNil(t, v)
		assert.Equal(t, "customer-id", v.CustomerId)
		assert.Equal(t, "plate", v.Plate)
		assert.Equal(t, "model", v.Model)
		assert.Equal(t, "brand", v.Brand)
		assert.Equal(t, 2023, v.Year)
	})

	t.Run("should return an error when plate is empty", func(t *testing.T) {
		_, err := NewVehicle("customer-id", "", "model", "brand", 2023)
		assert.ErrorIs(t, err, ErrInvalidPlate)
	})

	t.Run("should return an error when model is empty", func(t *testing.T) {
		_, err := NewVehicle("customer-id", "plate", "", "brand", 2023)
		assert.ErrorIs(t, err, ErrInvalidModel)
	})

	t.Run("should return an error when brand is empty", func(t *testing.T) {
		_, err := NewVehicle("customer-id", "plate", "model", "", 2023)
		assert.ErrorIs(t, err, ErrInvalidBrand)
	})

	t.Run("should return an error when year is invalid", func(t *testing.T) {
		_, err := NewVehicle("customer-id", "plate", "model", "brand", 0)
		assert.ErrorIs(t, err, ErrInvalidYear)
	})

	t.Run("should return an error when customer id is empty", func(t *testing.T) {
		_, err := NewVehicle("", "plate", "model", "brand", 2023)
		assert.ErrorIs(t, err, ErrInvalidCustomerId)
	})
}

func TestVehicle_Getters(t *testing.T) {
	v := &Vehicle{
		CustomerId: "customer-123",
		Plate:      "ABC-1234",
		Brand:      "TestBrand",
	}

	t.Run("should get the correct plate", func(t *testing.T) {
		assert.Equal(t, "ABC-1234", v.GetPlate())
	})

	t.Run("should get the correct brand", func(t *testing.T) {
		assert.Equal(t, "TestBrand", v.GetBrand())
	})

	t.Run("should get the correct customer id", func(t *testing.T) {
		assert.Equal(t, "customer-123", v.GetCustomerId())
	})
}
