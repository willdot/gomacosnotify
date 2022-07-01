package notify

import (
	_ "embed"
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

	fmt.Println(err)

	// don't create if already exists
	if err != nil {
		if !os.IsNotExist(err) {
			return nil
		}
		return err
	}

	err = os.WriteFile(f, alerter, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = os.Chmod(f, 0755)
	if err != nil {
		fmt.Println(err.Error())
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
		"-timeout", "1",
		"-title", "Bugsnag",
		"-subtitle", "This is a subtitle",
		"-json",
	}

	_, err := exec.Command(n.alerterLocation, args...).Output()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return nil
}
