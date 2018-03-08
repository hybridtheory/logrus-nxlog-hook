package nxlog

import (
	"crypto/tls"
	"fmt"
	"net"
	"syscall"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLogrusNxlogHook(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LogrusNxlogHook Suite")
}

const (
	MAX_BUFFER = 2048
)

var (
	tcpListener net.Listener
	tcpResult   = make(chan string)
	udpListener net.PacketConn
	udpResult   = make(chan string)
	tlsListener net.Listener
	tlsResult   = make(chan string)
	udsSocket   = "/tmp/unixdomain"
	udsListener *net.UnixConn
	udsResult   = make(chan string)
)

var _ = BeforeSuite(func() {
	fmt.Println("Opening test servers...")
	tcpListener = createTCPServer("tcp", ":3000")
	udpListener = createUDPServer("udp", ":3001")
	tlsListener = createTLSServer("ssl", ":3002")
	udsListener = createUDSServer("unixgram", udsSocket)
})

var _ = AfterSuite(func() {
	fmt.Println("Closing test servers...")
	tcpListener.Close()
	udpListener.Close()
	tlsListener.Close()
	syscall.Unlink(udsSocket)
})

func handler(received chan string, listener net.Listener) {
	fmt.Printf("Listening to connections to %s...\n", listener.Addr().String())
	for {
		connection, err := listener.Accept()
		go func(r chan string) {
			if err != nil {
				return
			}
			defer connection.Close()

			buffer := make([]byte, MAX_BUFFER)
			_, err = connection.Read(buffer)
			if err != nil {
				return
			}

			message := string(buffer)
			if message != "" {
				r <- message
			}
		}(received)
	}
}

func udpHandler(received chan string, listener net.PacketConn) {
	fmt.Println("Listening to udp connections...")
	for {
		buffer := make([]byte, MAX_BUFFER)
		listener.ReadFrom(buffer)
		message := string(buffer)
		if message != "" {
			received <- message
		}
	}
}

func unixgramHandler(received chan string, listener *net.UnixConn) {
	fmt.Println("Listening to unixgram connections...")
	for {
		buffer := make([]byte, MAX_BUFFER)
		listener.Read(buffer)
		message := string(buffer)
		if message != "" {
			received <- message
		}
	}
}

func createTCPServer(protocol string, endpoint string) net.Listener {
	fmt.Printf("Starting %s server...\n", protocol)
	listener, err := net.Listen(protocol, endpoint)
	Expect(err).To(BeNil())
	go handler(tcpResult, listener)
	return listener
}

func createUDPServer(protocol string, endpoint string) net.PacketConn {
	fmt.Printf("Starting %s server...\n", protocol)
	listener, err := net.ListenPacket(protocol, endpoint)
	Expect(err).To(BeNil())
	go udpHandler(udpResult, listener)
	return listener
}

func createTLSServer(protocol string, endpoint string) net.Listener {
	fmt.Printf("Starting %s server...\n", protocol)
	certificate, _ := tls.LoadX509KeyPair("cert/sample_server.crt", "cert/sample_server.key")
	configuration := &tls.Config{Certificates: []tls.Certificate{certificate}}
	listener, _ := tls.Listen("tcp", endpoint, configuration)
	go handler(tlsResult, listener)
	return listener
}

func createUDSServer(protocol string, endpoint string) *net.UnixConn {
	fmt.Printf("Starting %s server...\n", protocol)
	syscall.Unlink(endpoint)
	address, _ := net.ResolveUnixAddr(protocol, endpoint)
	syscall.Chmod(endpoint, 0755)
	listener, err := net.ListenUnixgram(protocol, address)
	Expect(err).To(BeNil())
	go unixgramHandler(udsResult, listener)
	return listener
}
