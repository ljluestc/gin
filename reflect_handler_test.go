// Copyright 2026 The Gin Authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type reflectHandlerService struct {
	called bool
}

func (s *reflectHandlerService) Ping(c *Context) {
	s.called = true
	c.Status(http.StatusNoContent)
}

func TestPointerReceiverReflectedMethodAsHandlerFunc(t *testing.T) {
	router := New()
	service := &reflectHandlerService{}

	method := reflect.ValueOf(service).MethodByName("Ping")
	assert.True(t, method.IsValid(), "expected reflected method to be found on pointer receiver")
	assert.Equal(t, 1, method.Type().NumIn())
	assert.Equal(t, reflect.TypeOf(&Context{}), method.Type().In(0))

	handler := HandlerFunc(func(c *Context) {
		method.Call([]reflect.Value{reflect.ValueOf(c)})
	})

	router.GET("/ping", handler)
	w := PerformRequest(router, http.MethodGet, "/ping")

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.True(t, service.called, "expected reflected handler to update receiver state")
}
