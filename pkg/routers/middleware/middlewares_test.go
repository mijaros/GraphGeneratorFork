package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/soch-fit/GraphGenerator/pkg/configuration"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	configuration.SetupTestingEnv()
	configuration.SetCookieInsecure()
}

func TestRecover(t *testing.T) {
	sut := gin.Default()
	sut.Use(Recover())
	sut.GET("/test", func(context *gin.Context) {
		param, e := context.GetQuery("test")
		if !e {
			panic("Expected panic")
		}
		context.JSON(200, gin.H{"param": param})
	})

	tester := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/test?test=ahoj", nil)
	if err != nil {
		t.Errorf("Couldn't start server %v", err.Error())
	}
	sut.ServeHTTP(tester, req)
	assert.Equal(t, http.StatusOK, tester.Code)
	resp, err := json.Marshal(gin.H{"param": "ahoj"})
	assert.JSONEq(t, tester.Body.String(), string(resp))

	tester = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/test?something=else", nil)
	sut.ServeHTTP(tester, req)
	assert.Equal(t, http.StatusInternalServerError, tester.Code)
}

func TestSetUpCookieMiddleware(t *testing.T) {
	sut := gin.Default()
	sut.Use(SetUpCookieMiddleware())
	testIdentifier := ""
	checkChange := false
	testingFunc := func(context *gin.Context) {
		id, ok := context.Get("identifier")
		if !ok {
			t.Errorf("Identifier must be always set")
			return
		}
		if checkChange {
			assert.Equal(t, testIdentifier, id)
		}
		testIdentifier = id.(string)
		context.Status(200)
	}

	sut.GET("/test", testingFunc)
	sut.POST("/test", testingFunc)

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/test", nil)
	assert.Nil(t, err)

	sut.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Set-Cookie"), testIdentifier)

	rec = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/test", nil)
	req.Header.Set("Cookie", SessionCookieName+"="+testIdentifier)
	checkChange = true
	sut.ServeHTTP(rec, req)

	testIdentifier = "And now for something completely different."
	rec = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/test", nil)
	assert.Nil(t, err)

	req.Header.Set("Cookie", SessionCookieName+"="+testIdentifier)
	checkChange = true
	sut.ServeHTTP(rec, req)
	//assert.Equal(t, SessionCookieName+"="+testIdentifier, rec.Header().Get("Set-Cookie"))
}
