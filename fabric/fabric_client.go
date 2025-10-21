package fabric

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Connect establishes a gateway connection to Hyperledger Fabric and returns the network & gateway instances.
func Connect() (*client.Network, *client.Gateway, error) { //Fabric channel and gateway connection to Fabric Peers
	// Paths to crypto material and Configs like we set from the CLI
	certPath := "fabric/crypto/cert.pem"
	keyPath := "fabric/crypto/priv_sk"
	tlsCertPath := "fabric/ca.crt"
	peerEndpoint := "localhost:7051"
	mspID := "Org1MSP"
	channelName := "mychannel"
	// Load the identity and private key
	id, err := newIdentity(certPath, mspID) //Loads the certificate and create a new Fabric Identity Object
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create identity: %v", err)
	}

	sign, err := newSign(keyPath) //Loads the private key and creates a signing function for transactions
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create signer: %v", err)
	}

	// Load peer TLS certificate
	cert, err := loadCertificate(tlsCertPath) //Loads peer TLS certificate
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load peer TLS certificate: %v", err)
	}
	certPool := x509.NewCertPool() //Creates a Certificate pool
	certPool.AddCert(cert)         //Adds the loaded TLS Certificate to the Certificate Pool
	// Create gRPC connection
	transportCreds := credentials.NewTLS(&tls.Config{ //Creating TLS credentials using the Certificate Pool
		RootCAs:            certPool,
		InsecureSkipVerify: false,
	})
	grpcConn, err := grpc.NewClient(peerEndpoint, grpc.WithTransportCredentials(transportCreds)) //Opens a new GRPC connection to the peer using TLS
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create gRPC connection: %v", err)
	}
	// Create Gateway connection
	gw, err := client.Connect( //Creates a Gate
		id,
		client.WithSign(sign),
		client.WithClientConnection(grpcConn),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect gateway: %v", err)
	}
	network := gw.GetNetwork(channelName)
	return network, gw, nil
}

// ----------------------------------------------------------------------
// Helper functions
// ----------------------------------------------------------------------

// newIdentity loads an X.509 certificate and creates a Fabric identity
func newIdentity(certPath string, mspID string) (*identity.X509Identity, error) {
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %v", err)
	}
	cert, err := identity.CertificateFromPEM(certPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %v", err)
	}
	id, err := identity.NewX509Identity(mspID, cert)
	if err != nil {
		return nil, fmt.Errorf("failed to create X509 identity: %v", err)
	}
	return id, nil
}

// newSign creates a signing function from a PEM-encoded private key
func newSign(keyPath string) (identity.Sign, error) {
	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %v", err)
	}

	privateKey, err := identity.PrivateKeyFromPEM(keyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %v", err)
	}
	return sign, nil
}

// loadCertificate reads a PEM-encoded certificate from disk
func loadCertificate(path string) (*x509.Certificate, error) {
	certPEM, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read cert file: %v", err)
	}

	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block from cert")
	}
	return x509.ParseCertificate(block.Bytes)
}
