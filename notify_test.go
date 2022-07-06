package notify

import (
	"os"
	"path"
	"testing"
)

func TestNewInstalls(t *testing.T) {
	t.Cleanup(func() {
		cleanTempDir(t)
	})

	// ensure we are using a clean directory setup
	cleanTempDir(t)

	_, err := New()
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(path.Join(os.TempDir(), tempDirectory))
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewInstallAlreadyExists(t *testing.T) {
	t.Cleanup(func() {
		cleanTempDir(t)
	})

	// ensure we are using a clean directory setup
	cleanTempDir(t)

	// run install once
	_, err := New()
	if err != nil {
		t.Fatal(err)
	}

	// run install again. Shouldn't get an error since it should already exist

	_, err = New()
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(path.Join(os.TempDir(), tempDirectory))
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewWithoutInstallation(t *testing.T) {
	t.Cleanup(func() {
		cleanTempDir(t)
	})

	// ensure we are using a clean directory setup
	cleanTempDir(t)

	NewWithCustomPath("some path")

	_, err := os.Stat(path.Join(os.TempDir(), tempDirectory))
	if err == nil {
		t.Fatal("Expected an error due to the file not being found but didn't get one indicating that the file was installed")
	}
}

func TestNotifyNoTitleSet(t *testing.T) {
	t.Cleanup(func() {
		cleanTempDir(t)
	})

	n, err := New()
	if err != nil {
		t.Fatal(err)
	}

	n.Message = "hello"
	_, err = n.Notify()
	if err == nil {
		t.Fatal("expected an error due to no title being set")
	}
}

func TestNotifyNoMessageSet(t *testing.T) {
	t.Cleanup(func() {
		cleanTempDir(t)
	})

	n, err := New()
	if err != nil {
		t.Fatal(err)
	}

	n.Title = "hello"
	_, err = n.Notify()
	if err == nil {
		t.Fatal("expected an error due to no message being set")
	}
}

func cleanTempDir(t *testing.T) {
	installDir := path.Join(os.TempDir(), tempDirectory)
	if err := os.RemoveAll(installDir); err != nil {
		t.Fatal(err)
	}
}
