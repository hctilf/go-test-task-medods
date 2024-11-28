package server

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/hctilf/go-test-task-medods/internal/app"
	"github.com/hctilf/go-test-task-medods/internal/controller/http/routes"
)

type Server struct {
	Srv    *fiber.App
	notify chan error
}

func NewServer(app *app.Application) *Server {
	serverFiber := fiber.New(
		fiber.Config{
			ServerHeader:      "Snippetbox",
			ReadTimeout:       app.Config.ReadTimeout,
			WriteTimeout:      app.Config.WriteTimeout,
			IdleTimeout:       app.Config.IdleTimeout,
			AppName:           "snippetbox",
			EnablePrintRoutes: true,
		},
	)

	err := routes.SetRoutes(app, serverFiber)
	if err != nil {
		app.Logger.Fatal(err)
	}

	Server := &Server{
		Srv:    serverFiber,
		notify: make(chan error, 1),
	}

	Server.Start(app)

	return Server
}

func (s *Server) Start(app *app.Application) {
	ln := createListener(app)
	if ln == nil {
		s.notify <- fmt.Errorf("failed to create listener")

		return
	}
	go func() {
		if err := s.Srv.Listener(*ln); err != nil {
			s.notify <- err
		}
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.Srv.ShutdownWithContext(ctx)
}

func createListener(app *app.Application) *net.Listener {
	certFile := "cert.pem"
	keyFile := "key.pem"

	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		app.Logger.Errorf("Self-signed certificate not found, generating...")
		if err := GenerateSelfSignedCert(certFile, keyFile); err != nil {
			return nil
		}
		app.Logger.Errorf("Self-signed certificate generated successfully")
		app.Logger.Errorf("You will need to accept the self-signed certificate in your browser")
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil
	}

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
	ln, err := tls.Listen("tcp", app.Config.Address, tlsConfig)
	if err != nil {
		return nil
	}

	return &ln
}

func GenerateSelfSignedCert(certFile string, keyFile string) error {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 180),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	certOut, err := os.Create(certFile)
	if err != nil {
		return err
	}
	defer certOut.Close()

	_ = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	keyOut, err := os.Create(keyFile)
	if err != nil {
		return err
	}
	defer keyOut.Close()

	_ = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return nil
}
