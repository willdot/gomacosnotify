package notify

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	tempDirectory = "gomacosnotify"
	binaryName    = "alerter"
)

var defaultTimeoutSeconds = 1

//go:embed assets/alerter
var alerter []byte

// Notifier is a client that allows notifications to be sent
type Notifier struct {
	alerterLocation string
}

// Notification contains all of the fields and settings that a notification can have
type Notification struct {
	Timeout      *int
	Title        string
	SubTitle     string
	ContentImage string
	Message      string
	CloseText    string
}

// SetTimeout will set the timeout in seconds. It must be > 0
func (n *Notification) SetTimeout(timeout int) error {
	if timeout < 0 {
		return errors.New("timeout must be greater than 0")
	}
	n.Timeout = &timeout
	return nil
}

// Response defines the response from a notification, including the actions that the user clicked on
type Response struct {
	ActivationType  string `json:"activationType"`
	ActivationValue string `json:"activationValue"`
}

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

// New will initialize a new notifier but allows the location of a prexisting alerter binary to be defined. This will
// not install anything on your machine
func NewWithCustomPath(alerterLocation string) *Notifier {
	return &Notifier{
		alerterLocation: alerterLocation,
	}
}

// Will send a notification
func (n *Notifier) Notify(notification Notification) (Response, error) {
	var resp Response

	args := []string{
		"-json",
	}

	if notification.Message == "" {
		return resp, errors.New("message must be set")
	}
	args = append(args, "-message", notification.Message)

	if notification.Title == "" {
		return resp, errors.New("title must be set")
	}

	args = append(args, "-title", notification.Title)

	if notification.SubTitle != "" {
		args = append(args, "-subtitle", notification.SubTitle)
	}

	if notification.Timeout == nil {
		notification.Timeout = &defaultTimeoutSeconds
	}
	if *notification.Timeout > 0 {
		args = append(args, "-timeout", fmt.Sprintf("%v", *notification.Timeout))
	}

	if notification.ContentImage != "" {
		args = append(args, "-contentImage", notification.ContentImage)
	}

	if notification.CloseText != "" {
		args = append(args, "-closeLabel", notification.CloseText)
	}

	output, err := exec.Command(n.alerterLocation, args...).Output()
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(output, &resp)
	if err != nil {
		return resp, errors.Wrap(err, "error decoding response")
	}

	return resp, nil
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
