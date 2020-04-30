package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"flag"

	//"fmt"
	"gocv.io/x/gocv"
	"math/big"

	quic "github.com/lucas-clemente/quic-go"
)

const addr = "localhost:4242"

// We start a server echoing data on the first stream the client opens,
// then connect with a client, send the message, and wait for its receipt.
func main() {
	multipath := flag.Bool("m", false, "multipath")
	err := clientMain(*multipath)
	if err != nil {
		panic(err)
	}
}

func clientMain(multipath bool) error {
	cfgClient := &quic.Config{
		CreatePaths: multipath,
	}
	session, err := quic.DialAddr(addr, &tls.Config{InsecureSkipVerify: true}, cfgClient)
	if err != nil {
		return err
	}

	stream, err := session.OpenStreamSync()
	if err != nil {
		return err
	}

	webcam, _ := gocv.VideoCaptureDevice(0)
	img := gocv.NewMat()

	for {
		webcam.Read(&img)
		bufr := img.ToBytes()
		err = putBlock(stream, bufr)
		if err != nil {
			return err
		}
	}

	return nil
}

func putBlock(stream quic.Stream, bufr []byte) error {
	bs := make([]byte, 4)
	msgLen := len(bufr)
	binary.LittleEndian.PutUint32(bs, uint32(msgLen))
	_, err := stream.Write(bs)
	if err != nil {
		return err
	}

	_, err = stream.Write(bufr)
	if err != nil {
		return err
	}

	return nil
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
