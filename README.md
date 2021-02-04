# Uber Zap logger wrapper with log rotation

## Installation

```
go get -u go.kuoruan.net/log
```

## Quick Start

Use the `global` logger.

```go
package main

import "go.kuoruan.net/log"

func main() {
    log.SetOptions(log.Development())

    log.Debug("this is debug log")
}
```

Or create new:

```go
package main

import "go.kuoruan.net/log"

func main() {
    logger := log.New(
        log.AddCaller(), 
        log.WithLogDirs("log"), 
        log.WithLogToStdout(false), 
        log.WithRotationConfig(log.RotationConfig{
                MaxSize: 500, // MB
                MaxAge: 3, // days
                MaxBackups: 7,
                LocalTime: true,
                Compress: true,
        }),
    )
    
    logger.Info("This is info log")
}
```
