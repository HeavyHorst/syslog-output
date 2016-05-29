```go
package main

import (
	"time"

	"github.com/HeavyHorst/syslog-output"
	"github.com/micro/go-platform/log"
)

func main() {
	out := syslog.NewOutput("tcp", log.OutputName("127.0.0.1:514"))
	l := log.NewLog(log.WithOutput(out), log.WithFields(log.Fields{
		"app":         "test",
		"environment": "dev",
	}))

	l.Info("info")
	l.Error("error")
	l.Fatal("fatal")
}
```
