package health

import (
	"github.com/Komdosh/golang-microservices/src/api/utils/test_utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConstants(t *testing.T) {
	assert.EqualValues(t, "alive", health)
}

func TestHealth(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	c := test_utils.GetMockedContext(request, response)

	Health(c)

	assert.EqualValues(t, http.StatusOK, response.Code)
	assert.EqualValues(t, "\"alive\"", response.Body.String())
}
