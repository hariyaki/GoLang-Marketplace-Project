package handlers

import (
	"fmt"
	"net/http"
)

// HealthHandler implements /healthz

// Health godoc
// @Summary  Health check
// @Tags     system
// @Success  200  {string}  string  "ok"
// @Router   /healthz [get]
type HealthHandler struct{}

func (HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}
