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

	n.Message = "YOYOYOYOYOYOYOYO"
	n.Title = "hello"

	n.SubTitle = "world"
	n.Timeout = 5
	n.ContentImage = "../RandomImage.png"

	n.CloseText = "CLOSE ME"

	resp, err := n.Notify()
	if err != nil {
		panic(err)
	}

	fmt.Printf("action: %s\n", resp.ActivationValue)
	fmt.Printf("action: %s\n", resp.ActivationType)

}
