#go-clementine#

[![goci](https://goci.herokuapp.com/project/image/github.com/brunoqc/go-clementine)](http://goci.me/project/github.com/brunoqc/go-clementine)

Control [Clementine](http://www.clementine-player.org/) using the [remote control feature](https://code.google.com/p/clementine-player/wiki/RemoteControl) and [goprotobuf](https://code.google.com/p/goprotobuf/).

For now, there's only 3 functions:
- SimplePlay()
- SimplePause()
- SimpleStop()

Those functions connect to Clementine, send the command and disconnect.

I'll support persistent connection later. If you need a feature you can open [a new issue](https://github.com/brunoqc/go-clementine/issues/new).

##Sample code##

```go
package main

import "github.com/brunoqc/go-clementine"

func main() {
	clementine := clementine.Clementine{
		Host:     "127.0.0.1",
		Port:     5500,
		AuthCode: 28615,
	}
	errPause := clementine.SimpleStop()
	if errPause != nil {
		panic(errPause)
	}
}
```
