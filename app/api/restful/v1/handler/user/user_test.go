package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	router "banking/app/api"
	userHdl "banking/app/api/restful/v1/handler/user"
	domainMock "banking/domain/mock"
	mysqlModel "banking/model/mysql"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.elastic.co/apm/v2"
)

func initialUserHandler(t *testing.T) (*gin.Context, *httptest.ResponseRecorder, *domainMock.MockIUserService, *domainMock.MockIAPIKeyService) {
	gin.SetMode(gin.TestMode)

	// Initialize APM tracer
	tracer := apm.DefaultTracer()
	router.InitRouter(gin.Default(), nil, nil, nil, tracer)

	ctrl := gomock.NewController(t)
	mockUserService := domainMock.NewMockIUserService(ctrl)
	mockAPIKeyService := domainMock.NewMockIAPIKeyService(ctrl)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	t.Cleanup(func() {
		ctrl.Finish()
	})

	return c, w, mockUserService, mockAPIKeyService
}

func Test_CreateUser(t *testing.T) {
	c, w, mockUserService, mockAPIKeyService := initialUserHandler(t)

	// variables
	user := &mysqlModel.User{
		Name:     "user1",
		Email:    "user1@yopmail.com",
		Password: "password",
		Balance:  decimal.NewFromFloat(0),
	}
	reqBody := userHdl.CreateUserReq{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	// mock
	mockUserService.EXPECT().
		CreateUser(gomock.Any(), gomock.Eq(user)).
		Return(nil)

	// request
	c.Request = httptest.NewRequest("POST", "/api/v1/user", bytes.NewReader(reqBodyBytes))

	// handler
	hdl := userHdl.NewUserHandler(mockUserService, mockAPIKeyService)
	hdl.CreateUser()(c)

	// Check status code
	assert.Equal(t, http.StatusCreated, w.Code, w.Body.String())

	// Check response body
	var actualResponse userHdl.CreateUserResp
	err = json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, actualResponse.Data.Name)
	assert.Equal(t, user.Balance, actualResponse.Data.Balance)
}
