package napnap

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextRemoteIpAddress(t *testing.T) {
	nap := New()
	nap.ForwardRemoteIpAddress = true
	c := NewContext(nap, nil, nil)

	c.Request, _ = http.NewRequest("POST", "/", nil)

	c.Request.Header.Set("X-Real-IP", " 10.10.10.10  ")
	c.Request.Header.Set("X-Forwarded-For", "  20.20.20.20, 30.30.30.30")
	c.Request.RemoteAddr = "  40.40.40.40:42123 "

	assert.Equal(t, "10.10.10.10", c.RemoteIpAddress())

	c.Request.Header.Del("X-Real-IP")
	assert.Equal(t, "20.20.20.20", c.RemoteIpAddress())

	c.Request.Header.Set("X-Forwarded-For", "30.30.30.30  ")
	assert.Equal(t, "30.30.30.30", c.RemoteIpAddress())

	c.Request.Header.Del("X-Forwarded-For")
	assert.Equal(t, "40.40.40.40", c.RemoteIpAddress())
}

func TestContextContentType(t *testing.T) {
	nap := New()
	c := NewContext(nap, nil, nil)

	c.Request, _ = http.NewRequest("POST", "/", nil)
	c.Request.Header.Set("Content-Type", "application/json; charset=utf-8")

	assert.Equal(t, c.ContentType(), "application/json")
}

func TestContextSetCookie(t *testing.T) {
	w := httptest.NewRecorder()
	nap := New()
	c := NewContext(nap, nil, w)

	c.SetCookie("user", "jason", 1, "/", "localhost", true, true)
	assert.Equal(t, "user=jason; Path=/; Domain=localhost; Max-Age=1; HttpOnly; Secure", c.Writer.Header().Get("Set-Cookie"))
}

func TestContextGetCookie(t *testing.T) {
	nap := New()
	c := NewContext(nap, nil, nil)

	c.Request, _ = http.NewRequest("GET", "/get", nil)
	c.Request.Header.Set("Cookie", "user=jason")
	cookie, _ := c.Cookie("user")
	assert.Equal(t, "jason", cookie)
}
