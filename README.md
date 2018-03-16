# [NXLog](https://logrus.co/) Hook for [Logrus](https://github.com/sirupsen/logrus) <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>

[![GoDoc](https://godoc.org/github.com/nathany/looper?status.svg)](https://godoc.org/github.com/affectv/logrus-nxlog-hook)
[![Build Status](https://travis-ci.org/affectv/logrus-nxlog-hook.svg?branch=master)](https://travis-ci.org/affectv/logrus-nxlog-hook)

This [Logrus](https://github.com/sirupsen/logrus) hook allows sending logs to a [NXLog](https://logrus.co/) server.

## Configuration

The hook must be configured with:

* The protocol NXLog is using (`tcp`, `udp`, `ssl` or `unixgram`)
* The endpoint where NXLog is listening to:
    * As an `"ip:port"` or `"hostname:port"` string, in case of an IP protocol.
    * or as a valid `path`, in case of UDS.

## Usage

After importing the `logrus-nxlog-hook` package, you can create a new hook with:

`NewNxlogHook([PROTOCOL], [ENDPOINT], [OPTIONAL_TLS_CONFIGURATION])`

For example, for a TCP setup:

```go
package main

import (
  log "github.com/sirupsen/logrus"
  nxlog "github.com/affectv/logrus-nxlog-hook"
)

func main() {
  nxlogHook, err := nxlog.NewNxlogHook("tcp", "127.0.0.1:514", nil)
  if err == nil {
    log.AddHook(nxlogHook)
  }
  log.Info("Hello World")
}
```

See the [Examples](#Examples) section for other setups.

### Filtering by Level

This hook allows setting the minimum logging level which
will trigger the forwarding of the messages.

You can specify it with `Level`

e.g.

```go
func main() {
  nxlogHook, _ := nxlog.NewNxlogHook("tcp", "127.0.0.1:514", nil)
  nxlogHook.Level = log.InfoLevel
  log.AddHook(nxlogHook)
  log.Info("Info message to be logged")
}
```

### Specifying the Format

The default logging formatter is JSON, but you can change it as follows, using `Formatter`

e.g.

```go
func main() {
  nxlogHook, _ := nxlog.NewNxlogHook("tcp", "127.0.0.1:514", nil)
  nxlogHook.Formatter = &log.TextFormatter{}
  log.AddHook(nxlogHook)
  log.Info("Info message to be logged")
}
```

## <a name="Examples"></a>Examples

Some examples of how to use the plugin are the following:

### TCP Connection

Specify `tcp` and a `"host:port"` string to `NewNxlogHook`.

e.g.

```go
func main() {
  nxlogHook, _ := nxlog.NewNxlogHook("tcp", "127.0.0.1:514", nil)
  log.AddHook(nxlogHook)
  log.Info("TCP info message to be logged")
}
```

### UDP Connection

Specify `udp` and a `"host:port"` string to `NewNxlogHook`.

```go
func main() {
  nxlogHook, _ := nxlog.NewNxlogHook("udp", "127.0.0.1:514", nil)
  log.AddHook(nxlogHook)
  log.Info("UDP info message to be logged")
}
```

### SSL Connection

Please remember you should have a compatible SSL certificate available to use. If you don't have it, you can create one with:

```bash
openssl req -newkey rsa:2048 -nodes -keyout sample-key.pem -x509 -days 365 -out sample-certificate.pem
```

You must ensure NXLog is configured to accept Connections via SSL as well. The SSL certificate files should be included in the following section of the `nxlog.conf` file:

```xml
<Input in_ssl>
    Module im_ssl
    Host 0.0.0.0
    Port 516
    CertFile	   %CERTSERVER%/sample-certificate.pem
    CertKeyFile	 %CERTSERVER%/sample-key.pem
    AllowUntrusted True
</Input>
```

Now you are ready to connect via SSL. Specify `ssl` and a `"host:port"` string to `NewNxlogHook`.

```go
import (
	"crypto/tls"
	nxlog "github.com/affectv/logrus-nxlog-hook"
	"github.com/sirupsen/logrus"
)

func main() {
  certificate, _ := tls.LoadX509KeyPair(
		"%CERTCLIENT%/sample-certificate.pem",
		"%CERTCLIENT%/sample-key.pem",
	)
  nxlogHook, _ := nxlog.NewNxlogHook("ssl", "127.0.0.1:514", &tls.Config{
    InsecureSkipVerify: true,
    Certificates: []tls.Certificate{certificate},
  })
  log.AddHook(nxlogHook)
  log.Info("SSL info message to be logged")
}
```

### UDS Connection

Specify `unixgram` and `"path"` string to `NewNxlogHook`.

```go
func main() {
  nxlogHook, _ := nxlog.NewNxlogHook("unixgram", "/var/run/nxlog/devlog", nil)
  log.AddHook(nxlogHook)
  log.Info("UDS info message to be logged")
}
```

## Plugin Development and Testing

This plugin was developed using [ginkgo](https://github.com/onsi/ginkgo)
and [gomega](https://github.com/onsi/gomega) test frameworks.

To execute the tests:

```go
ginkgo ./...
```
