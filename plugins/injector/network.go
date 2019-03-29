package injector

import (
	"trump/pkg/middleware"
	"log"
	"time"
	"github.com/docker/libchan/spdy"
	"trump/pkg/inject"
	"fmt"
	"trump/pkg/proxy"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
)

var PLUGIN_PRIORITY = 10
var Plugin = middleware.Plugin{
	Process:  ProcessMsg,
	Priority: PLUGIN_PRIORITY,
	Init:     Init,
	Shutdown: Shutdown,
}
var SERVER_CERT = ""
var SERVER_KEY = ""
var CA_CERT = ""
var SERVER_PORT uint16 = 9998

var listener net.Listener
var con net.Conn

func Init() {
	// Load server cert
	serverCert, err := tls.LoadX509KeyPair(SERVER_CERT, SERVER_KEY)
	if err != nil {
		log.Fatalf("failed to load server certificate: %v", err)
	}
	// Load CA cert
	caCert, err := ioutil.ReadFile(CA_CERT)
	if err != nil {
		log.Fatalf("failed to load ca certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatalf("failed to parse client certificate authority")
	}
	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()
	listener, err = tls.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", SERVER_PORT), tlsConfig)
	if err != nil {
		log.Fatalf("failed to start listening for connections: %v", err)
	}
	go func() {
		for {
			con, err = listener.Accept()
			if err == nil {
				p, err := spdy.NewSpdyStreamProvider(con, true)
				if err == nil {
					t := spdy.NewTransport(p)
					r, err := t.WaitReceiveChannel()
					if err == nil {
						msg := &inject.Data{}
						for {
							err = r.Receive(msg)
							if err != nil {
								fmt.Printf("failed to receive message: %v", err)
								continue
							}
							proxy.InjectMsg(msg, PLUGIN_PRIORITY)
						}
						continue
					}
				}
			}
			log.Printf("failed to accept connection: %v", err)
			time.Sleep(time.Second)
		}
	}()
}

func ProcessMsg(arg interface{}) interface{} {
	return arg
}

func Shutdown() {
	if con != nil {
		con.Close()
	}
	if listener != nil {
		listener.Close()
	}
}
