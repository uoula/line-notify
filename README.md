# LINE Notify

Line notify with proxy option

## Installation

```bash
go get -u github.com/uoula/line-notify
```

## Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/uoula/line-notify/notifier"
)

func main() {

    notifier := notifier.Notifier{}

    token := "O6****************0ItHq2WyNQcxFhOS8i2"
    message := "Hello World!"
    
    response, err := notifier.NotifyMessage(token, message)
    if err != nil {
        log.Println(err.Error())
    }
    fmt.Println(response)

}
```

### Notify With Proxy

```go
package main

import (
	"net/url"

	"github.com/uoula/line-notify/notifier"
)

func main() {

    notifier := notifier.Notifier{}
    notifier.Proxy, _ = url.Parse("http://proxy.example.com:port")

    token := "O6****************0ItHq2WyNQcxFhOS8i2"
    message := "Hello World!"

    notifier.NotifyMessage(token, message)
}
```

### Notify With Image URL

```go
token := "O6****************0ItHq2WyNQcxFhOS8i2"
message := "Hello World!"
imageURL := "http://example.com/image.jpg"
imageThumbnail := "http://example.com/thumbnail.jpg" //optional (can be left blank)

notifier.NofityImageURL(token, message, imageURL, imageThumbnail)
```

### Notify With Image File (PNG, JPEG only)

```go
token := "O6****************0ItHq2WyNQcxFhOS8i2"
message := "Hello World!"
imageFile := "/path/to/your/file.jpeg"

notifier.NotifyImageFile(token, message, imageFile)
```

### Notify With Sticker

```go
token := "O6****************0ItHq2WyNQcxFhOS8i2"
message := "Hello World!"
stickerPackageId := 446
stickerId := 1988

notifier.NotifySticker(token, message, stickerPackageId, stickerId)
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.