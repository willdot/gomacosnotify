# About

GoMacOSNotify is a Go library for sending desktop applications on macOS.

To do this it installs a binary called [alerter](https://github.com/vjeantet/alerter) into a temp directory on your machine. However if you wish to install the alerter binary into a location yourself, you can configure it to use your location instead of it installing it. The location it installs to will be something similar to:
`/var/folders/kf/tkk3036976b2x37h6sld0cnc0000gn/T/gomacosnotify`

# Usage
## Creating the notification client
### Using the installation method
``` go
notifier, err := notify.New()
if err != nil {
    panic(err)
}
```
### Supplying custom installation path
``` go
notifier := notify.NewWithCustomPath("some/path/to/alerter")
```

## Sending a notification

``` go
notification := notify.Notification{
    Title:        "Demo notification",
    SubTitle:     "Some information",
    Message:      "Some more detailed information",
}

resp, err := notifier.Send(notification)
if err != nil {
    panic(err)
}

// check the response here
```

The response from the `send` function returns the type of action the user took, and the value of the button they pressed (if any).

`ActivationType` contains the type of action (timeout, closed, contentsClicked etc)
`ActivationValue` contains the value of the button that was pressed, if applicable. (See the Close Button section for more details)

NOTE: If the user actions the notification AFTER a timeout has passed, the action will not be captured as the `Send` function would have returned before the users action was captured. Be sure you use a suitable timeout.


## Configuration
There are options that can be configured.
### Timeout
A timeout (in seconds) can be configured. Some things to note:
* If not configured a default of 10 seconds will be used. 
* A timeout of 0 seconds will mean the notification will not disappear until the user actions it.
* When calling the `Send` function, it will block and not return until either the timeout times out, or the user actions the notification.

``` go
err := notification.SetTimeout(5)
if err != nil {
    panic(err)
}
```

### Close

MacOS notifications support user actions. Currently only closing is supported by this library. To do this, set the `CloseText` field of the notification to be the text you wish to be displayed to the user. You can then use the response from the `Send` function, to see if the user pressed the close button.

The following example will create a notification that has a close button with the text `ignore`.

``` go
notification := notify.Notification{
    Title:        "Demo notification",
    SubTitle:     "Some information",
    Message:      "Some more detailed information",
    CloseText:    "ignore",
}

resp, err := notifier.Send(notification)
if err != nil {
    panic(err)
}

fmt.Printf("action type: %s\n", resp.ActivationType)
fmt.Printf("action value: %s\n", resp.ActivationValue)

// prints
// action type: close
// action value: ignore
```

### Images
You can include a small thumbnail image inside the notification by passing in the path to a `png` file.

``` go
notification := notify.Notification{
    Title:        "Demo notification",
    SubTitle:     "Some information",
    Message:      "Some more detailed information",
    ContentImage: "../RandomImage.png",
}

resp, err := n.Send(notification)
if err != nil {
    panic(err)
}
```