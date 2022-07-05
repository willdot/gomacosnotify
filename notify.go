package notify

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

//go:embed assets/alerter
var alerter []byte

type Notifier struct {
	alerterLocation string

	Timeout      int
	Title        string
	SubTitle     string
	ContentImage string
	Message      string
}

const (
	tempDirectory = "gomacosnotify"
	binaryName    = "alerter"
)

// New will initialize a new notifier. It will install the alerter binary into a temp location on your machine
func New() (*Notifier, error) {
	err := install()
	if err != nil {
		return nil, err
	}

	return &Notifier{
		alerterLocation: path.Join(os.TempDir(), tempDirectory, binaryName),
	}, nil
}

func install() error {
	p := path.Join(os.TempDir(), tempDirectory)

	// ensure the temp directory already exists
	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		return err
	}

	f := filepath.Join(p, binaryName)

	_, err = os.Stat(f)
	if err == nil {
		// it must exist so don't install again
		return nil
	}

	// if the error is that it doesn't exist, that's expected, so ignore it
	if !os.IsNotExist(err) {
		return err
	}

	err = os.WriteFile(f, alerter, 0644)
	if err != nil {
		return err
	}

	err = os.Chmod(f, 0755)
	if err != nil {
		return err
	}
	return nil
}

// New will initialize a new notifier but allows the location of a prexisting alerter binary to be defined. This will
// not install anything on your machine
func NewWithCustomPath(alerterLocation string) *Notifier {
	return &Notifier{
		alerterLocation: alerterLocation,
	}
}

// Will send a notification
func (n *Notifier) Notify() error {
	args := []string{
		"-closeLabel", "ignore",
		"-json",
	}

	if n.Message == "" {
		return errors.New("message must be set")
	}
	args = append(args, "-message", n.Message)

	if n.Title == "" {
		return errors.New("title must be set")
	}

	args = append(args, "-title", n.Title)

	if n.SubTitle != "" {
		args = append(args, "-subtitle", n.SubTitle)
	}

	if n.Timeout > 0 {
		args = append(args, "-timeout", fmt.Sprintf("%v", n.Timeout))
	}

	if n.ContentImage != "" {
		args = append(args, "-contentImage", n.ContentImage)
	}

	_, err := exec.Command(n.alerterLocation, args...).Output()
	if err != nil {
		return nil
	}
	return nil
}
