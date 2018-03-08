# NXLog Hook for [Logrus](https://github.com/sirupsen/logrus)

[![Build Status](https://travis-ci.org/affectv/logrus-nxlog-hook.svg?branch=master)]
(https://travis-ci.org/affectv/logrus-nxlog-hook)

This logrus hook allows sending logs to your [NXLog](https://nxlog.co/) instance over
any of the protocols the server supports.

## Usage

We are using the standard net package, so in order to use this logrus hook
we will need to provide:

* A protocol where NXLog is listening (`tcp`, `udp`, `ssl` or `unixgram`)
* The endpoint where NXLog is listening ("ip:port" string in case of an ip protocol
or a valid path in case of uds)

The hook must be configured with:

```go
package main

import (
  log "github.com/sirupsen/logrus"
  nxlog "github.com/affectv/logrus-nxlog-hook"
)

func main() {
  hook := nxlog.NewNxlogHook("tcp", "ip:port", nil)
  log.AddHook(hook)
  log.Info("info message to be logged")
}
```

### Filter by level

NXLog allows us to filter for certain levels, but maybe we want not to send them
in the first place. So the hook allows setting the minimum logging level that
will trigger this hook to forward the messages.

```go
func main() {
  hook, _ := nxlog.NewNxlogHook("tcp", "ip:port", nil)
  hook.Level = log.InfoLevel
  log.AddHook(hook)
  log.Info("info message to be logged")
}
```

### Specify formatting

The default logging formatter is JSON (which NXLog is able to decode easily
to append new variables, remove them, etc.).

We can change the formatter as follows:

```go
func main() {
  hook, _ := nxlog.NewNxlogHook("tcp", "ip:port", nil)
  hook.Formatter = &log.TextFormatter{}
  log.AddHook(hook)
  log.Info("info message to be logged")
}
```

## Examples

Some examples of how to use the plugin are the following:

### TCP connections

This plugin allows to connect with NXLog via the `im_tcp` module.

```go
func main() {
  hook, err := nxlog.NewNxlogHook("tcp", "ip:port", nil)
  if err == nil {
    log.AddHook(hook)
  }
  log.Info("tcp info message to be logged")
}
```

### UDP connections

This plugin allows to connect with NXLog via the `im_udp` module.

```go
func main() {
  hook, err := nxlog.NewNxlogHook("udp", "ip:port", nil)
  if err == nil {
    log.AddHook(hook)
  }
  log.Info("udp info message to be logged")
}
```

### SSL connections

This plugin allows to connect with NXLog via the `im_ssl` module.

In order to test the behaviour, we are adding the steps to have NXLog `im_ssl`
up and running with a basic SSL configuration.

First we can create a compatible SSL certificate with:

```bash
openssl req -newkey rsa:2048 -nodes -keyout sample-key.pem -x509 -days 365 -out sample-certificate.pem
```

Then we can configure NXLog to accept connections via SSL with:

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

To be able to send messages to the SSL listener, we must pass the TLS configuration
specifying a client's certificate as well:

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
  hook, err := nxlog.NewNxlogHook("ssl", "ip:port", &tls.Config{
    InsecureSkipVerify: true,
    Certificates: []tls.Certificate{certificate},
  })  
  if err == nil {
    log.AddHook(hook)
  }
  log.Info("ssl info message to be logged")
}
```

### UDS connections

This plugin allows to connect with NXLog via the im_uds module.

```go
func main() {
  hook, err := nxlog.NewNxlogHook("unixgram", "/var/run/nxlog/devlog", nil)
  if err == nil {
    log.AddHook(hook)
  }
  log.Info("uds info message to be logged")
}
```

## Testing

To execute the plugin tests we have the [ginkgo](https://github.com/onsi/ginkgo)
and [gomega](https://github.com/onsi/gomega) dependencies.

To execute the tests:

```go
ginkgo ./...
```
