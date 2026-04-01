package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/oficina-utils/pkg/valueobjects/plate"
)

func TestNewVehicle(t *testing.T) {
	validCustomerID := uuid.New()
	validPlate, _ := plate.New("ABC1234")
	validBrand := "Toyota"
	validModel := "Corolla"
	validYear := 2020
	currentYear := time.Now().Year()

	tests := []struct {
		name       string
		id         uuid.UUID
		customerID uuid.UUID
		plateVO    plate.Plate
		brand      string
		model      string
		year       int
		wantErr    error
	}{
		{
			name:       "valid vehicle",
			id:         uuid.Nil,
			customerID: validCustomerID,
			plateVO:    validPlate,
			brand:      validBrand,
			model:      validModel,
			year:       validYear,
			wantErr:    nil,
		},
		{
			name:       "valid vehicle with specific ID",
			id:         uuid.New(),
			customerID: validCustomerID,
			plateVO:    validPlate,
			brand:      validBrand,
			model:      validModel,
			year:       validYear,
			wantErr:    nil,
		},
		{
			name:       "valid vehicle with current year",
			id:         uuid.Nil,
			customerID: validCustomerID,
			plateVO:    validPlate,
			brand:      validBrand,
			model:      validModel,
			year:       currentYear,
			wantErr:    nil,
		},
		{
			name:       "valid vehicle with next year",
			id:         uuid.Nil,
			customerID: validCustomerID,
			plateVO:    validPlate,
			brand:      validBrand,
			model:      validModel,
			year:       currentYear + 1,
			wantErr:    nil,
		},
		{
			name:       "empty brand",
			id:         uuid.Nil,
			customerID: validCustomerID,
			plateVO:    validPlate,
			brand:      "",
			model:      validModel,
			year:       validYear,
			wantErr:    ErrInvalidVehicleBrand,
		},
		{
			name:       "whitespace only brand",
			id:         uuid.Nil,
			customerID: validCustomerID,
			plateVO:    validPlate,
			brand:      "   ",
			model:      validModel,
			year:       validYear,
			wantErr:    ErrInvalidVehicleBrand,
		},
		{
			name:       "empty model",
			id:         uuid.Nil,
			customerID: validCustomerID,
			plateVO:    validPlate,
			brand:      validBrand,
			model:      "",
			year:       validYear,
			wantErr:    ErrInvalidVehicleModel,
		},
		{
			name:       "whitespace only model",
			id:         uuid.Nil,
			customerID: validCustomerID,
			plateVO:    validPlate,
			brand:      validBrand,
			model:      "   ",
			year:       validYear,
			wantErr:    ErrInvalidVehicleModel,
		},
		{
			name:       "invalid year - too old",
			id:         uuid.Nil,
			customerID: validCustomerID,
			plateVO:    validPlate,
			brand:      validBrand,
			model:      validModel,
			year:       1899,
			wantErr:    ErrInvalidVehicleYear,
		},
		{
			name:       "invalid year - too far in future",
			id:         uuid.Nil,
			customerID: validCustomerID,
			plateVO:    validPlate,
			brand:      validBrand,
			model:      validModel,
			year:       currentYear + 2,
			wantErr:    ErrInvalidVehicleYear,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewVehicle(tt.id, tt.customerID, tt.plateVO, tt.brand, tt.model, tt.year)
			if err != tt.wantErr {
				t.Errorf("NewVehicle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				if got.ID == uuid.Nil {
					t.Error("NewVehicle() ID should not be Nil")
				}
				if got.CustomerID != tt.customerID {
					t.Errorf("NewVehicle() CustomerID = %v, want %v", got.CustomerID, tt.customerID)
				}
				if got.Plate != tt.plateVO {
					t.Errorf("NewVehicle() Plate = %v, want %v", got.Plate, tt.plateVO)
				}
				if got.Brand != tt.brand {
					t.Errorf("NewVehicle() Brand = %v, want %v", got.Brand, tt.brand)
				}
				if got.Model != tt.model {
					t.Errorf("NewVehicle() Model = %v, want %v", got.Model, tt.model)
				}
				if got.Year != tt.year {
					t.Errorf("NewVehicle() Year = %v, want %v", got.Year, tt.year)
				}
				if got.CreatedAt.IsZero() {
					t.Error("NewVehicle() CreatedAt should not be zero")
				}
				if got.UpdatedAt.IsZero() {
					t.Error("NewVehicle() UpdatedAt should not be zero")
				}
			}
		})
	}
}

func TestVehicleChangePlate(t *testing.T) {
	customerID := uuid.New()
	oldPlate, _ := plate.New("ABC1234")
	newPlate, _ := plate.New("XYZ5678")

	vehicle, _ := NewVehicle(uuid.Nil, customerID, oldPlate, "Toyota", "Corolla", 2020)
	originalUpdatedAt := vehicle.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	err := vehicle.ChangePlate(newPlate)

	if err != nil {
		t.Errorf("ChangePlate() error = %v, want nil", err)
	}
	if vehicle.Plate != newPlate {
		t.Errorf("ChangePlate() Plate = %v, want %v", vehicle.Plate, newPlate)
	}
	if !vehicle.UpdatedAt.After(originalUpdatedAt) {
		t.Error("ChangePlate() should update UpdatedAt")
	}
}

func TestVehicleChangeBrand(t *testing.T) {
	customerID := uuid.New()
	plateVO, _ := plate.New("ABC1234")

	vehicle, _ := NewVehicle(uuid.Nil, customerID, plateVO, "Toyota", "Corolla", 2020)
	originalUpdatedAt := vehicle.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	tests := []struct {
		name     string
		newBrand string
		wantErr  error
	}{
		{
			name:     "valid brand change",
			newBrand: "Honda",
			wantErr:  nil,
		},
		{
			name:     "empty brand",
			newBrand: "",
			wantErr:  ErrInvalidVehicleBrand,
		},
		{
			name:     "whitespace only brand",
			newBrand: "   ",
			wantErr:  ErrInvalidVehicleBrand,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := vehicle.ChangeBrand(tt.newBrand)
			if err != tt.wantErr {
				t.Errorf("ChangeBrand() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if vehicle.Brand != tt.newBrand {
					t.Errorf("ChangeBrand() Brand = %v, want %v", vehicle.Brand, tt.newBrand)
				}
				if !vehicle.UpdatedAt.After(originalUpdatedAt) {
					t.Error("ChangeBrand() should update UpdatedAt")
				}
			}
		})
	}
}

func TestVehicleChangeModel(t *testing.T) {
	customerID := uuid.New()
	plateVO, _ := plate.New("ABC1234")

	vehicle, _ := NewVehicle(uuid.Nil, customerID, plateVO, "Toyota", "Corolla", 2020)
	originalUpdatedAt := vehicle.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	tests := []struct {
		name     string
		newModel string
		wantErr  error
	}{
		{
			name:     "valid model change",
			newModel: "Camry",
			wantErr:  nil,
		},
		{
			name:     "empty model",
			newModel: "",
			wantErr:  ErrInvalidVehicleModel,
		},
		{
			name:     "whitespace only model",
			newModel: "   ",
			wantErr:  ErrInvalidVehicleModel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := vehicle.ChangeModel(tt.newModel)
			if err != tt.wantErr {
				t.Errorf("ChangeModel() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if vehicle.Model != tt.newModel {
					t.Errorf("ChangeModel() Model = %v, want %v", vehicle.Model, tt.newModel)
				}
				if !vehicle.UpdatedAt.After(originalUpdatedAt) {
					t.Error("ChangeModel() should update UpdatedAt")
				}
			}
		})
	}
}

func TestVehicleChangeYear(t *testing.T) {
	customerID := uuid.New()
	plateVO, _ := plate.New("ABC1234")
	currentYear := time.Now().Year()

	vehicle, _ := NewVehicle(uuid.Nil, customerID, plateVO, "Toyota", "Corolla", 2020)
	originalUpdatedAt := vehicle.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	tests := []struct {
		name    string
		newYear int
		wantErr error
	}{
		{
			name:    "valid year change",
			newYear: 2021,
			wantErr: nil,
		},
		{
			name:    "current year",
			newYear: currentYear,
			wantErr: nil,
		},
		{
			name:    "next year",
			newYear: currentYear + 1,
			wantErr: nil,
		},
		{
			name:    "invalid year - too old",
			newYear: 1899,
			wantErr: ErrInvalidVehicleYear,
		},
		{
			name:    "invalid year - too far in future",
			newYear: currentYear + 2,
			wantErr: ErrInvalidVehicleYear,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := vehicle.ChangeYear(tt.newYear)
			if err != tt.wantErr {
				t.Errorf("ChangeYear() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if vehicle.Year != tt.newYear {
					t.Errorf("ChangeYear() Year = %v, want %v", vehicle.Year, tt.newYear)
				}
				if !vehicle.UpdatedAt.After(originalUpdatedAt) {
					t.Error("ChangeYear() should update UpdatedAt")
				}
			}
		})
	}
}
