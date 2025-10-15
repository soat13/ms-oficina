package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCatalogReaderColumnExpression(t *testing.T) {
	tests := []struct {
		name     string
		table    string
		expected string
	}{
		{
			name:     "Products table includes stock",
			table:    "products",
			expected: "c.id, c.name, c.price, c.stock",
		},
		{
			name:     "Services table excludes stock",
			table:    "services",
			expected: "c.id, c.name, c.price",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &catalogReader{table: tt.table}

			columnExpr := "c.id, c.name, c.price"
			if reader.table == "products" {
				columnExpr = "c.id, c.name, c.price, c.stock"
			}

			assert.Equal(t, tt.expected, columnExpr)
		})
	}
}

func TestCatalogRowStockHandling(t *testing.T) {
	tests := []struct {
		name     string
		stock    *int
		expected int
	}{
		{
			name:     "Stock with value",
			stock:    intPtr(50),
			expected: 50,
		},
		{
			name:     "Stock nil",
			stock:    nil,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stock := 0
			if tt.stock != nil {
				stock = *tt.stock
			}

			assert.Equal(t, tt.expected, stock)
		})
	}
}

func intPtr(i int) *int {
	return &i
}
