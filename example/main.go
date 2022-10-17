package main

import (
	"fmt"

	notify "github.com/willdot/gomacosnotify"
)

func main() {
	n, err := notify.New()
	if err != nil {
		panic(err)
	}

	notification := notify.Notification{
		Title:        "Demo notification",
		SubTitle:     "Some information",
		Message:      "Some more detailed information",
		ContentImage: "../RandomImage.png",
		CloseText:    "go away",
		Actions:      []string{"Option 1", "Option 2"},
	}

	_ = notification.SetTimeout(5)

	resp, err := n.Send(notification)
	if err != nil {
		panic(err)
	}

	fmt.Printf("action type: %s\n", resp.ActivationType)
	fmt.Printf("action value: %s\n", resp.ActivationValue)

}
