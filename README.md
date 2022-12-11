# gophercord

A powerful Go library to interact with the Discord API.

## Installation

```bash
go get github.com/davipatricio/gophercord
```

## Usage

```go
package main

import (
	"os"

	"github.com/davipatricio/gophercord/client"
)

func main() {
	bot := client.NewClient("abc")
	bot.Connect()

	// wait for a control c
	<-make(chan os.Signal, 1)
}
```

## Contributing

Feel free to open an issue or a pull request.
Fork the repository, make your changes and open a pull request.

## License

[MIT](https://choosealicense.com/licenses/mit/)
