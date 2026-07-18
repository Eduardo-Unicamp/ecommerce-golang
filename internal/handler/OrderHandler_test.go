package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"first-api/internal/middleware"
	"first-api/internal/mocks"
	"first-api/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetOrderByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockOrderUseCase := mocks.NewMockOrderUseCase(ctrl)
	orderHandler := NewOrderHandler(mockOrderUseCase)

	validUUID := uuid.MustParse("019f5936-935b-7e2a-985e-468106c73243")

	testOrder := model.Order{
		ID:         validUUID,
		CustomerID: validUUID,
		Items: []model.OrderItem{
			{
				ID:           validUUID,
				ProductID:    validUUID,
				SellingPrice: decimal.NewFromInt(1),
				UnitsOrdered: 1,
			},
		},
		Status: model.PENDING,
	}

	//Arrange
	testCases := []struct {
		name           string
		id             string
		setupMocks     func()
		expectedStatus int
		expectedResult *model.Order
		expectedError  error
	}{
		{
			name: "SUCCESS Id válido, retorna order",
			id:   validUUID.String(),
			setupMocks: func() {
				mockOrderUseCase.EXPECT().GetOrderByID(gomock.Any(), "019f5936-935b-7e2a-985e-468106c73243", "019f5936-935b-7e2a-985e-468106c73243").
					Return(&testOrder, nil).Times(1)
			},
			expectedResult: &testOrder,
			expectedStatus: 200,
			expectedError:  nil,
		},
		{
			name: "Id em formato válido, mas order inexistente",
			id:   "019f5936-935b-7e2a-985e-468106c73243", //not in the database
			setupMocks: func() {
				mockOrderUseCase.EXPECT().GetOrderByID(gomock.Any(), "019f5936-935b-7e2a-985e-468106c73243", "019f5936-935b-7e2a-985e-468106c73243").Return(nil, sql.ErrNoRows).Times(1)
			},
			expectedResult: nil,
			expectedStatus: 404,
			expectedError:  model.ErrOrderNotFound,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			tt.setupMocks()

			r := httptest.NewRequest(http.MethodGet, "/order/{order_id}", nil)

			ctxChi := chi.NewRouteContext()
			ctxChi.URLParams.Add("order_id", tt.id)

			ctx := context.WithValue(r.Context(), chi.RouteCtxKey, ctxChi)
			ctx = context.WithValue(ctx, middleware.UserIDKey, tt.id)

			r = r.WithContext(ctx)

			w := httptest.NewRecorder()

			//Act
			orderHandler.GetOrderByID(w, r)

			//Assert
			if tt.expectedError == nil {
				resultOrder := &model.Order{}
				err := json.NewDecoder(w.Body).Decode(resultOrder)

				assert.Equal(t, tt.expectedResult, resultOrder, "")
				assert.Equal(t, tt.expectedStatus, w.Code)
				assert.NoError(t, err, "Nao deve dar erro")
			} else {
				//se da erro vai só escrever o erro direto no w, entao verifica que teve erro garantindo que nao foi criado o header http correto
				assert.NotContains(t, w.Header().Get("Content-Type"), "application/json")
				assert.Equal(t, tt.expectedStatus, w.Code)
			}

		})
	}

}
