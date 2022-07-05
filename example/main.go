package main

import notify "github.com/willdot/gomacosnotify"

func main() {
	n, err := notify.New()
	if err != nil {
		panic(err)
	}

	n.Message = "YOYOYOYOYOYOYOYO"
	n.Title = "hello"

	n.SubTitle = "world"
	n.Timeout = 2
	n.ContentImage = "../RandomImage.png"

	err = n.Notify()
	if err != nil {
		panic(err)
	}
}
