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

var mockUserService *domainMock.MockIUserService

func initialUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Initialize APM tracer
	tracer := apm.DefaultTracer()
	router.InitRouter(gin.Default(), nil, nil, tracer)

	ctrl := gomock.NewController(t)
	mockUserService = domainMock.NewMockIUserService(ctrl)
}

func Test_CreateUser(t *testing.T) {
	initialUserHandler(t)

	// variables
	user := &mysqlModel.User{
		Name:    "user1",
		Balance: decimal.NewFromFloat(0),
	}
	reqBody := userHdl.CreateUserReq{Name: user.Name}
	reqBodyBytes, _ := json.Marshal(reqBody)

	// mock
	mockUserService.EXPECT().
		CreateUser(gomock.Any(), gomock.Eq(user)).
		Return(nil)

	// request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/user", bytes.NewReader(reqBodyBytes))

	// handler
	hdl := userHdl.NewUserHandler(mockUserService)
	hdl.CreateUser()(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "user created")
}
