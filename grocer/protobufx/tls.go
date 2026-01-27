package protobufx

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

// AsConfig converts TLSOptions to a tls.Config pointer
// It loads the certificate and key pair, sets up the CA certificates if provided,
// and configures the server name for SNI (Server Name Indication)
func (x *TLSOptions) AsConfig() *tls.Config {
	if x == nil {
		return nil
	}
	// Load the X509 key pair from the specified certificate and key files
	cert, err := tls.LoadX509KeyPair(x.GetCertFile().GetValue(), x.GetKeyFile().GetValue())
	if err != nil {
		// Panic if unable to load the certificate/key pair
		panic(err)
	}
	// Create base TLS configuration with the loaded certificate
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}

	// Handle CA certificate if provided
	if x.GetCaFile() != nil {
		// Read the CA certificate file
		caCert, err := os.ReadFile(x.GetCaFile().GetValue())
		if err != nil {
			// Panic if unable to read the CA certificate file
			panic(err)
		}
		// Create a new certificate pool and append the CA certificate
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		// Set the root CAs for the TLS configuration
		tlsConfig.RootCAs = caCertPool
	}

	// Set the server name if provided for SNI (Server Name Indication)
	if x.GetServerName() != nil {
		tlsConfig.ServerName = x.GetServerName().GetValue()
	}
	return tlsConfig
}