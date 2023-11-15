# Ding Go SDK

> [!WARNING]
> This SDK is deprecated in favor of the [Ding API Go SDK](https://github.com/ding-live/ding-golang).

The Ding Go SDK allows Go backends to access the Ding API programmatically.

## Requirements

Go 1.16 or higher

## Installation

Make sure your project is using Go Modules (it will have a `go.mod` file in its
root if it already is):

```sh
go mod init
```

Then, reference ding-go in a Go program with `import`:

```go
import (
	ding "github.com/ding-live/ding-go"
)
```

Run any of the normal `go` commands (`build`/`install`/`test`). The Go
toolchain will resolve and fetch the ding-go module automatically.

Alternatively, you can also explicitly `go get` the package into a project:

```bash
go get -u github.com/ding-live/ding-go
```

## Documentation

Below are a few simple examples:

### Create a new client

```go
cfg := &ding.Config{
	CustomerUUID:      "f0b399ce-eead-4781-bab5-f63240e81a52",
	APIKey:            "87ydfsa987fdyas9h8f7y29ne87fyqds98af",
	MaxNetworkRetries: ding.Int(4),
}

c, err := ding.NewClient(cfg)
```

### Send a message

This will send an OTP code to a client using the best route available

```go
a, err := client.Authenticate(ding.AuthenticateOptions{
	PhoneNumber: "+33xxxxxxxxxx",
	IP:          ding.String("192.168.0.1"),
	DeviceType:  &ding.DeviceTypeIOS,
	AppVersion:  ding.String("1.2.0"),
	CallbackURL: ding.String("https://example.com/callback"),
})
```

### Check a code

When the user enters the code into your app, check whether it is valid

```go
a, err := client.Check("f0b399ce-eead-4781-bab5-f63240e81a52", "3588")
```

### Retry an authentication

```go
a, err := client.Retry("5071dbf5-78d0-497a-b844-c1231808c3e9")
```

### Use a custom HTTP client

If you want more control on the HTTP requests that are performed by the client,
you can use a custom HTTP client.

```go
c, err := ding.NewClient(ding.Config{
	CustomerUUID:      "f0b399ce-eead-4781-bab5-f63240e81a52",
	APIKey:            "87ydfsa987fdyas9h8f7y29ne87fyqds98af",
	CustomHTTPClient:  &http.Client{Timeout: 10 * time.Second},
})
```

### Configure logging

By default, the library logs error messages only (which are sent to stderr). Configure default logging using the LeveledLogger field:

```go
config := &ding.Config{
    LeveledLogger: &ding.Logger{
        Level: ding.LevelInfo,
    },
}
```

It's possible to use non-Ding leveled loggers as well. Ding expects loggers to comply to the following interface:

```go
type LeveledLogger interface {
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
}
```

Some loggers like [Logrus][logrus] and Zap's [SugaredLogger][zapsugaredlogger]
support this interface out-of-the-box so it's possible to set
`DefaultLeveledLogger` to a `*logrus.Logger` or `*zap.SugaredLogger` directly.
For others it may be necessary to write a thin shim layer to support them.

[logrus]: https://github.com/sirupsen/logrus/
[zapsugaredlogger]: https://godoc.org/go.uber.org/zap#SugaredLogger
