package fabric

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func ConnectDyn(mspID string, peerEndpoint string, channelName string, certPem string, keyPem string) (*client.Network, *client.Gateway, error) { //Fabric channel and gateway connection to Fabric Peers
	// Paths to crypto material and Configs like we set from the CLI
	// certPath := "fabric/crypto/cert.pem"
	// keyPath := "fabric/crypto/priv_sk"
	tlsCertPath := "fabric/ca-cert.pem"
	peerEndpoint = "localhost:7051"
	// mspID := "Org1MSP"
	// channelName := "mychannel"
	// Load the identity and private key
	id, err := newIdentitynew(certPem, mspID) //Loads the certificate and create a new Fabric Identity Object
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create identity: %v", err)
	}

	sign, err := newSignnew(keyPem) //Loads the private key and creates a signing function for transactions
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create signer: %v", err)
	}

	// Load peer TLS certificate
	cert, err := loadCertificatenew(tlsCertPath) //Loads peer TLS certificate
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
func newIdentitynew(certPEMpassed string, mspID string) (*identity.X509Identity, error) {
	// fmt.Println("newIdentity cert passed", certPEMpassed)
	certPEM := Concatcerts(certPEMpassed)
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
func newSignnew(keyPEMpassed string) (identity.Sign, error) {
	// fmt.Println("keyPEMpassed", keyPEMpassed)
	keyPEM := Concatcerts(keyPEMpassed)
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
func loadCertificatenew(certPEMpassed string) (*x509.Certificate, error) {
	certPEM, err := os.ReadFile(certPEMpassed)
	if err != nil {
		return nil, fmt.Errorf("failed to read cert file: %v", err)
	}
	// certPEM := Concatcerts(certPEMpassed)
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block from cert")
	}
	return x509.ParseCertificate(block.Bytes)
}

func Concatcerts(cert string) []byte {
	certNew := fmt.Sprintf("-----BEGIN PRIVATE KEY-----\n%s\n-----END PRIVATE KEY-----\n", cert)
	return []byte(certNew)
}

func DecryptAESGCM(key, ct []byte) ([]byte, error) {
	if len(ct) < 12 {
		return nil, errors.New("ciphertext too short")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aesgcm.NonceSize()
	if len(ct) < nonceSize {
		return nil, errors.New("ciphertext too short for nonce")
	}
	nonce, ciphertext := ct[:nonceSize], ct[nonceSize:]
	pt, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return pt, nil
}
