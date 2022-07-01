package notify

import (
	"testing"
)

func TestNotifyInstalls(t *testing.T) {
	notifier, err := New()
	if err != nil {
		t.Fatal(err)
	}

	err = notifier.Notify()
	if err != nil {
		t.Fatal(err)
	}
}
