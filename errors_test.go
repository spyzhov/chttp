package chttp

import "testing"

func TestError_Error(t *testing.T) {
	var err *Error
	if err.Error() != "" {
		t.Errorf("broken error")
	}
}
