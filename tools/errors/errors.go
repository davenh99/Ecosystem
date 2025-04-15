package errors

import (
	"apps/ecosystem/tools"
	"fmt"
	"net/http"
)

func PermissionDenied(w http.ResponseWriter) {
	tools.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}