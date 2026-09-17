package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/josiastomasnanez/finflow/internal/lab"
)

func TestLabOKHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/lab/ok", nil)
	labOKHandler(c)
	if w.Code != http.StatusOK {
		t.Fatalf("code %d", w.Code)
	}
}

func TestLabErrorHandlerClampsTo5xx(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/lab/error", strings.NewReader(`{"code":404}`))
	c.Request.Header.Set("Content-Type", "application/json")
	labErrorHandler(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("code %d", w.Code)
	}
}

func TestLabMemoryRetainAndRelease(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mem := &lab.Hold{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/lab/memory", strings.NewReader(`{"megabytes":2}`))
	c.Request.Header.Set("Content-Type", "application/json")
	labMemoryRetainHandler(mem)(c)
	if w.Code != http.StatusOK {
		t.Fatalf("retain %d", w.Code)
	}
	if mem.MiB() != 2 {
		t.Fatalf("held %d", mem.MiB())
	}

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest(http.MethodPost, "/api/lab/memory/release", strings.NewReader(`{}`))
	c2.Request.Header.Set("Content-Type", "application/json")
	labMemoryReleaseHandler(mem)(c2)
	if mem.MiB() != 0 {
		t.Fatalf("after release %d", mem.MiB())
	}
}
