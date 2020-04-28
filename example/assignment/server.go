package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	//"fmt"
	"gocv.io/x/gocv"
	"io"
	"math/big"

	quic "github.com/lucas-clemente/quic-go"
)

const addr = "localhost:4242"

//const message = "foobar"

// We start a server echoing data on the first stream the client opens,
// then connect with a client, send the message, and wait for its receipt.
func main() {
	err := echoServer()
	if err != nil {
		panic(err)
	}
}

// Start a server that echos all data on the first stream opened by the client
func echoServer() error {
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		return err
	}
	sess, err := listener.Accept()
	if err != nil {
		return err
	}
	stream, err := sess.AcceptStream()
	if err != nil {
		panic(err)
	}
	// Echo through the loggingWriter
	//_, err = io.Copy(loggingWriter{stream}, stream)
	window := gocv.NewWindow("Hello")
	for {
		byteImg, err := getBolck(stream)
		image2, err := gocv.NewMatFromBytes(480, 640, 16, byteImg)

		if err != nil {
			panic(err)
		}
		window.IMShow(image2)
		window.WaitKey(1)
	}
	return err
}

func getBolck(stream quic.Stream) (bufr []byte, e error) {
	bs := make([]byte, 4)
	_, err := io.ReadFull(stream, bs)
	if err != nil {
		return nil, err
	}

	msgLen := binary.LittleEndian.Uint32(bs)

	ls := make([]byte, msgLen)
	_, err = io.ReadFull(stream, ls)

	return ls, nil
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}
