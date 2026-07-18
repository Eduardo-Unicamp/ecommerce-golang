package repository

import (
	"context"
	"database/sql"
	"first-api/internal/model"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGetOrderByID(t *testing.T) {
	dummyID := "019f5936-935b-7e2a-985e-468106c73243"

	successOrderResult := &model.Order{
		ID:         uuid.MustParse(dummyID),
		CustomerID: uuid.MustParse(dummyID),
		Status:     model.PENDING,
		Items: []model.OrderItem{
			{
				ID:           uuid.MustParse(dummyID),
				ProductID:    uuid.MustParse(dummyID),
				SellingPrice: decimal.NewFromFloat(1.0),
				UnitsOrdered: 1,
			},
		},
	}

	testCases := []struct {
		name           string
		ID             string
		customerID     string
		mockfunc       func(mockDB pgxmock.PgxPoolIface)
		expectedResult *model.Order
		expectsError   bool
	}{
		{
			name:       "Erro: Begin nao funciona, deve retornar nil e error",
			ID:         dummyID,
			customerID: dummyID,
			mockfunc: func(mockDB pgxmock.PgxPoolIface) {
				mockDB.ExpectBegin().WillReturnError(sql.ErrConnDone)
			},
			expectedResult: nil,
			expectsError:   true,
		},
		{
			name:       "Sucesso,retorna *model.Order,nil",
			ID:         "019f5936-935b-7e2a-985e-468106c73243",
			customerID: "019f5936-935b-7e2a-985e-468106c73243",
			mockfunc: func(mockDB pgxmock.PgxPoolIface) {
				mockDB.ExpectBegin()

				orderRows := mockDB.NewRows([]string{"id", "status", "customer_id"}).
					AddRow(dummyID, "PENDING", dummyID)

				mockDB.ExpectQuery("^SELECT \\* FROM orders WHERE id=\\$1 AND customer_id=\\$2;$").
					WithArgs(dummyID, dummyID).
					WillReturnRows(orderRows)

				itemRows := mockDB.NewRows([]string{"id", "selling_price", "units", "product_id"}).
					AddRow(dummyID, decimal.NewFromFloat(1.0), 1, dummyID)

				mockDB.ExpectQuery("^SELECT id,selling_price,units,product_id FROM order_items WHERE order_id=\\$1$").
					WithArgs(dummyID).
					WillReturnRows(itemRows)

				mockDB.ExpectCommit()

			},
			expectedResult: successOrderResult,
			expectsError:   false,
		},
	}
	//Arrange
	//Act
	for _, tt := range testCases {
		mockDB, err := pgxmock.NewPool()
		if err != nil {
			log.Fatalf("Erro ao criar mock do banco: %v", err)
		}
		defer mockDB.Close()
		orderRepository := NewOrderRepository(mockDB)

		tt.mockfunc(mockDB)
		result, err := orderRepository.GetOrderByID(context.Background(), tt.ID, tt.customerID)

		//Assert
		if !tt.expectsError {
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result)
		} else {
			assert.Error(t, err)
		}

	}

}
