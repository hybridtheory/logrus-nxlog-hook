package nxlog

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Nxlog writer connector", func() {

	Context("with TCP connections", func() {

		It("connects successfully to a listening tcp endpoint", func() {
			_, err := NewWriter("tcp", "127.0.0.1:3000", nil)
			Expect(err).To(BeNil())
		})

		It("shows an error if connecting to an invalid TCP endpoint", func() {
			_, err := NewWriter("tcp", "127.0.0.1:0", nil)
			Expect(err).ToNot(BeNil())
		})

		It("sends properly the data via TCP", func() {
			data := "test input tcp"
			writer, _ := NewWriter("tcp", "127.0.0.1:3000", nil)
			writer.Write([]byte(data), true)
			Eventually(tcpResult).Should(Receive(&data))
		})
	})

	Context("with UDP connections", func() {

		It("connects successfully to a listening UDP endpoint", func() {
			_, err := NewWriter("udp", "127.0.0.1:3001", nil)
			Expect(err).To(BeNil())
		})

		It("sends properly the data via UDP", func() {
			data := "test input udp"
			writer, _ := NewWriter("udp", "127.0.0.1:3001", nil)
			writer.Write([]byte(data), true)
			Eventually(udpResult).Should(Receive(&data))
		})
	})

	Context("with SSL connections", func() {

		It("connects successfully to a listening SSL endpoint", func() {
			_, err := NewWriter("ssl", "127.0.0.1:3002", nil)
			Expect(err).To(BeNil())
		})

		It("shows an error if connecting to an invalid SSL endpoint", func() {
			_, err := NewWriter("ssl", "127.0.0.1:0", nil)
			Expect(err).ToNot(BeNil())
		})

		It("sends properly the data via SSL", func() {
			data := "test input tls"
			writer, _ := NewWriter("ssl", "127.0.0.1:3002", nil)
			writer.Write([]byte(data), true)
			Eventually(tlsResult).Should(Receive(&data))
		})
	})

	Context("with UDS connections", func() {

		It("connects successfully to a listening UDS endpoint", func() {
			_, err := NewWriter("unixgram", udsSocket, nil)
			Expect(err).To(BeNil())
		})

		It("shows an error if connecting to an invalid UDS endpoint", func() {
			_, err := NewWriter("unixgram", "/tmp/fake", nil)
			Expect(err).ToNot(BeNil())
		})

		It("sends properly the data via UDS", func() {
			data := "test input uds"
			writer, _ := NewWriter("unixgram", udsSocket, nil)
			writer.Write([]byte(data), true)
			Eventually(udsResult).Should(Receive(&data))
		})
	})
})
