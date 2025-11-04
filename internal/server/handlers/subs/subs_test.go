package subs

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"log/slog"

	"github.com/P3rCh1/subs-aggregator/internal/config"
	"github.com/P3rCh1/subs-aggregator/internal/models"
	"github.com/P3rCh1/subs-aggregator/internal/storage/postgres"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Create(ctx context.Context, sub *models.Subscription) error {
	args := m.Called(ctx, sub)
	return args.Error(0)
}

func (m *MockDB) Read(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *MockDB) Update(ctx context.Context, sub *models.Subscription) error {
	args := m.Called(ctx, sub)
	return args.Error(0)
}

func (m *MockDB) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDB) List(ctx context.Context, userID uuid.UUID) ([]models.Subscription, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Subscription), args.Error(1)
}

func (m *MockDB) Summary(ctx context.Context, req *models.SumRequest) (int, error) {
	args := m.Called(ctx, req)
	return args.Int(0), args.Error(1)
}

func (m *MockDB) Close() error {
	args := m.Called()
	return args.Error(0)
}

func defaultSub() models.Subscription {
	return models.Subscription{
		ServiceName: "Netflix",
		Price:       1000,
		UserID:      uuid.New(),
		StartDate: models.MonthDate{
			Time:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
	}
}

func defaultSumRequest() models.SumRequest {
	return models.SumRequest{
		UserID: uuid.New(),
		StartDate: models.MonthDate{
			Time:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
		EndDate: models.MonthDate{
			Time:  time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
	}
}

func setup() (*MockDB, *echo.Echo) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{}

	mockDB := &MockDB{}

	api := NewServerAPI(logger, cfg, mockDB)

	e := echo.New()
	e.POST("/subs", api.Create)
	e.GET("/subs/:id", api.Read)
	e.PUT("/subs/:id", api.Update)
	e.DELETE("/subs/:id", api.Delete)
	e.GET("/subs/list/:id", api.List)
	e.POST("/subs/summary", api.Summary)

	return mockDB, e
}

func TestCreate_Success(t *testing.T) {
	mockDB, e := setup()

	sub := defaultSub()

	mockDB.On("Create", mock.Anything, mock.AnythingOfType("*models.Subscription")).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*models.Subscription)
			arg.ID = uuid.New()
		})

	body, _ := json.Marshal(sub)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/subs", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var response models.Subscription
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NotEqual(t, uuid.Nil, response.ID)

	mockDB.AssertExpectations(t)
}

func TestCreate_InvalidJSON(t *testing.T) {
	_, e := setup()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/subs", bytes.NewReader([]byte("{invalid json")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRead_Success(t *testing.T) {
	mockDB, e := setup()

	id := uuid.New()
	expectedSub := defaultSub()
	expectedSub.ID = id

	mockDB.On("Read", mock.Anything, id).Return(&expectedSub, nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/subs/"+id.String(), nil)
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.Subscription
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, expectedSub, response)

	mockDB.AssertExpectations(t)
}

func TestRead_NotFound(t *testing.T) {
	mockDB, e := setup()

	id := uuid.New()

	mockDB.On("Read", mock.Anything, id).Return(nil, postgres.ErrNotFound)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/subs/"+id.String(), nil)
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	mockDB.AssertExpectations(t)
}

func TestRead_InvalidID(t *testing.T) {
	_, e := setup()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/subs/invalid-uuid", nil)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdate_Success(t *testing.T) {
	mockDB, e := setup()

	id := uuid.New()
	sub := defaultSub()
	sub.ServiceName = "Netflix Updated"

	mockDB.On("Update", mock.Anything, mock.AnythingOfType("*models.Subscription")).
		Return(nil)

	body, _ := json.Marshal(sub)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/subs/"+id.String(), bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockDB.AssertExpectations(t)
}

func TestUpdate_NotFound(t *testing.T) {
	mockDB, e := setup()

	id := uuid.New()

	mockDB.On("Update", mock.Anything, mock.AnythingOfType("*models.Subscription")).
		Return(postgres.ErrNotFound)

	body, _ := json.Marshal(defaultSub())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/subs/"+id.String(), bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockDB.AssertExpectations(t)
}

func TestDelete_Success(t *testing.T) {
	mockDB, e := setup()

	id := uuid.New()

	mockDB.On("Delete", mock.Anything, id).Return(nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/subs/"+id.String(), nil)
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockDB.AssertExpectations(t)
}

func TestList_Success(t *testing.T) {
	mockDB, e := setup()

	userID := uuid.New()

	expectedSubs := []models.Subscription{defaultSub(), defaultSub()}
	expectedSubs[0].ServiceName = "Netflix"
	expectedSubs[1].ServiceName = "Spotify"

	mockDB.On("List", mock.Anything, userID).Return(expectedSubs, nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/subs/list/"+userID.String(), nil)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response []models.Subscription
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Len(t, response, 2)
	assert.Equal(t, expectedSubs[0], response[0])
	assert.Equal(t, expectedSubs[1], response[1])

	mockDB.AssertExpectations(t)
}

func TestSummary_Success(t *testing.T) {
	mockDB, e := setup()

	req := defaultSumRequest()
	expectedSum := 18000

	mockDB.On("Summary", mock.Anything, &req).Return(expectedSum, nil)

	body, _ := json.Marshal(req)
	rec := httptest.NewRecorder()
	reqHttp := httptest.NewRequest(http.MethodPost, "/subs/summary", bytes.NewReader(body))
	reqHttp.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	e.ServeHTTP(rec, reqHttp)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]int
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, expectedSum, response["summary"])

	mockDB.AssertExpectations(t)
}

func TestSummary_InvalidDates(t *testing.T) {
	_, e := setup()

	invalidReq := models.SumRequest{
		StartDate: models.MonthDate{
			Time:  time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
		EndDate: models.MonthDate{
			Time:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
	}

	body, _ := json.Marshal(invalidReq)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/subs/summary", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestInternalError(t *testing.T) {
	mockDB, e := setup()

	sub := defaultSub()

	mockDB.On("Create", mock.Anything, mock.AnythingOfType("*models.Subscription")).
		Return(assert.AnError)

	body, _ := json.Marshal(sub)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/subs", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockDB.AssertExpectations(t)
}
