package sqles

import (
	"context"
	"testing"
)

func TestConnect(t *testing.T) {
	db, err := Connect(context.Background(), DefaultDriver, ":memory")
	if err != nil {
		t.Errorf("got err %q, want <nil>", err)
	}
	if db == nil {
		t.Error("Expected database got <nil>")
	}
}
