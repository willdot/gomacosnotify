package notify

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

const (
	tempDirectory = "gomacosnotify"
	binaryName    = "alerter"
)

var defaultTimeoutSeconds = 10

//go:embed assets/alerter
var alerter []byte

// Notifier is a client that allows notifications to be sent
type Notifier struct {
	alerterLocation string
}

// Notification contains all of the fields and settings that a notification can have
type Notification struct {
	Title        string
	SubTitle     string
	ContentImage string
	Message      string
	CloseText    string
	Actions      []string
	timeout      *int
}

// SetTimeout will set the timeout in seconds. It must be >= 0 and if not set, a default will be applied.
// If 0 is supplied, the notification won't timeout and will remain until the user has closed it or taken
// action.
func (n *Notification) SetTimeout(timeout int) error {
	if timeout < 0 {
		return errors.New("timeout must be greater or equal to 0")
	}
	n.timeout = &timeout
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

// Send will send a notification.
// If a timeout is set, and a user interacts with the notification after the timeout has expired,
// the user interaction will not be captured.
// This function blocks and will not return unless either a configured timeout expires or a user takes "action"
// on the notification.
func (n *Notifier) Send(notification Notification) (Response, error) {
	var resp Response

	if notification.Message == "" {
		return resp, errors.New("message must be set")
	}
	if notification.Title == "" {
		return resp, errors.New("title must be set")
	}

	args := []string{
		"-json", // ensures the response is in JSON format so we can unmarshal it
	}

	args = append(args, "-message", notification.Message)
	args = append(args, "-title", notification.Title)

	if notification.SubTitle != "" {
		args = append(args, "-subtitle", notification.SubTitle)
	}

	if notification.timeout == nil {
		notification.timeout = &defaultTimeoutSeconds
	}
	if *notification.timeout >= 0 {
		args = append(args, "-timeout", fmt.Sprintf("%v", *notification.timeout))
	}

	if notification.ContentImage != "" {
		args = append(args, "-contentImage", notification.ContentImage)
	}

	if notification.CloseText != "" {
		args = append(args, "-closeLabel", notification.CloseText)
	}

	if len(notification.Actions) > 0 {

		actions := strings.Join(notification.Actions, ",")
		args = append(args, "-actions", actions)
	}

	output, err := exec.Command(n.alerterLocation, args...).Output()
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(output, &resp)
	if err != nil {
		return resp, errors.Wrap(err, "error decoding response")
	}

	fmt.Println(resp)

	return resp, nil
}

func install() error {
	tempDirPath := path.Join(os.TempDir(), tempDirectory)

	// ensure the temp directory already exists
	err := os.MkdirAll(tempDirPath, os.ModePerm)
	if err != nil {
		return err
	}

	alerterPath := filepath.Join(tempDirPath, binaryName)

	_, err = os.Stat(alerterPath)
	// it must exist so don't install again
	if err == nil {
		return nil
	}

	// we are expecting a not exists error, so only return if it's a different error
	if !os.IsNotExist(err) {
		return err
	}

	err = os.WriteFile(alerterPath, alerter, 0644)
	if err != nil {
		return err
	}

	// make the binary runnable
	err = os.Chmod(alerterPath, 0755)
	if err != nil {
		return err
	}
	return nil
}
