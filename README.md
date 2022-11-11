# log

This is an opinionated, "batteries included" log package that caters to how I think logging
should happen. The functionalities that are being included are:

- rudimentary configuration
- structured logging
- providing PII in log fields with the ability to apply hashing or other means to it or remove them altogether

For further info, have a look at the examples section.

**Warning: When using the plain logging functions of the library instead
of instantiating a logger, you have no configuration options and no way
to handle PII in any special or different way.**

# Base parameters

- Log levels:
  - Debug
  - Info
  - Warn
  - Error
  - Panic
  - Fatal
- Timestamp format: RFC 3339
- Log format: JSON
- Output destinations: stdout, stderr
- Caller info: included
- Stacktrace: only enabled for warn and above
- Key names:
  - Application name key: "app" (set by user)
  - Version key: "version" (set by user)
  - Message key: "msg"
  - Level key: "lvl"
  - Time key: "ts"
  - Name key: "name"
  - Caller key: "caller"
  - Function key: "func"
  - Stacktrace key: "stacktrace"
- Available modes for dealing with PII:
  - none (leaves fields as is)
  - hash (hashes the value with SHA256)
  - mask (uses a custom mask function to mask values -- mask function needs to be provided by the user, when choosing this mode -- log.MaskFunc)
  - remove (removes the whole field from logs)

# Examples

## Instantiate the most basic logger

```go
package main

import "github.com/Rapix-x/log"

func main() {
	// There will be no "app" or "version" field on the logger/the logs
	// The implied log level is "info"
	// The implied PII mode is "none", so PII will be logged as is
	logger := log.MustNewLogger(log.Configuration{})
    defer logger.Sync()
	
    logger.Info("log something")
	// output: {"lvl":"info","ts":"1970-01-01T04:02:00+01:00","caller":"main/main.go:12","func":"main.main","msg":"log something"}
}
```

## Instantiate a production logger

```go
package main

import "github.com/Rapix-x/log"

func main() {
	logger, err := log.NewLogger(log.Configuration{
        ApplicationName: "example-app",
        Version:         "1.0.0",
        MinimumLogLevel: log.WarnLevel,
        PIIMode:         log.PIIModeRemove,
    }) 
    if err != nil {
        log.Fatalf("error occurred while instantiating new logger: %v", err)
    }
    defer logger.Sync()

    logger.Warn("Log something")
	// output: {"lvl":"warn","ts":"1970-01-01T04:02:00+01:00","caller":"main/main.go:19","func":"main.main","msg":"Log something","app":"example-app",
	// "version":"1.0.0","stacktrace":"main.main\n\t/log/main/main.go:17\nruntime.main\n\t/opt/homebrew/Cellar/go/1.19.3/libexec/src/runtime/proc.go:250"}
}
```

## Configuration

As seen in the examples above, the actual configuration of a logger only has four properties to configure.

```go
package main

import "github.com/Rapix-x/log"

func main() {
    conf := log.Configuration{
        ApplicationName: "example-app",
        Version:         "1.0.0",
        MinimumLogLevel: log.WarnLevel,
        PIIMode:         log.PIIModeRemove,
    }
}
```

## PII Mode and beyond

This package provides the capability to handle PII in any logs. All you have to do is to attach the PII
as a field using log statements with fields and wrap it with the provided functions. Below are three examples
that show the vanilla PII functionality, how to provide a custom function when selecting the PIIModeMask
and lastly some PII handling with a custom function just for one data set.

### Vanilla PII Handling

```go
package main

import "github.com/Rapix-x/log"

func main() {
	logger := log.MustNewLogger(log.Configuration{
        PIIMode: log.PIIModeHash,
    })
    defer logger.Sync()
	
    logger.Infow("Log PII fields", log.PII("username", "abc@example.com"))
	// output: {"lvl":"info","ts":"1970-01-01T04:02:00+01:00","caller":"main/main.go:14","func":"main.main","msg":"Log PII fields","username":"9eceb13483d7f187ec014fd6d4854d1420cfc634328af85f51d0323ba8622e21"}
}
```

### Custom Function for "mask" PII Mode

```go
package main

import "github.com/Rapix-x/log"

func main() {
  log.MaskFunc = maskIt
  logger := log.MustNewLogger(log.Configuration{
    PIIMode: log.PIIModeMask,
  })
  defer logger.Sync()

  logger.Infow("Log PII fields", log.PII("usernam", "abc@example.com"))
  // output: {"lvl":"info","ts":"1970-01-01T04:02:00+01:00","caller":"main/main.go:12","func":"main.main","msg":"Log PII fields","username":"gotcha value, hehe"}
}

func maskIt(key, value string) log.ResolvedPIIField {
  field := log.ResolvedPIIField{}
  field.Key = key + "e" // let's fix our typo here, shall we?
  field.Value = "gotcha value, hehe"

  return field
}
```

### Custom Function for Single PII Field

```go
package main

import "github.com/Rapix-x/log"

func main() {
  logger := log.MustNewLogger(log.Configuration{
    PIIMode: log.PIIModeHash,
  })
  defer logger.Sync()

  logger.Infow("Log PII fields", log.CustomPII("username", "abc@example.com", singleFieldMask))
  // output: {"lvl":"info","ts":"1970-01-01T04:02:00+01:00","caller":"main/main.go:11","func":"main.main","msg":"Log PII fields","username":"let's assume this is a hash *coughs in hex*"}
}

func singleFieldMask(mode log.PIIMode, key, value string) log.ResolvedPIIField {
  switch mode {
  case log.PIIModeHash:
    return log.ResolvedPIIField{Key: key, Value: "let's assume this is a hash *coughs in hex*"}
  case log.PIIModeMask:
    return log.ResolvedPIIField{Key: key, Value: "this one is masked"}
  case log.PIIModeRemove:
    return log.ResolvedPIIField{}
  default:
    return log.ResolvedPIIField{Key: key, Value: value}
  }
}
```
