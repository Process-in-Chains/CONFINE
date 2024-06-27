package attestation

import (
	"app/utils/http"
	"bytes"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/edgelesssys/ego/attestation"
	"github.com/edgelesssys/ego/attestation/tcbstatus"
	"github.com/edgelesssys/ego/eclient"
)

/*NOTE for remote attestation: if in the /etc/sgx_default_qcnl.conf file the field 'collateral_service' field should be commented
otherwise, the verifyReport method would directly communicate with the collateral_service without passing through your PCCS
comment 'collateral_service' when your pccs is set_up and specify its url in the /etc/sgx_default_qcnl.conf file */

/*Method invoked to execute the remote attestation of the Secure Miner hosted in a given address*/
func RemoteAttestation(serverAddr string, expectedMeasurement []byte) ([]byte, []byte) {
	// Get miner certificate. Skip TLS certificate verification because the certificate is self-signed and we will verify it using the report instead.
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	// Get miner's organization TLS certificate
	minerCertBytes := http.HttpGet(tlsConfig, serverAddr+"/cert")
	// Validate the certificate
	if validateCertificate(minerCertBytes) != nil {
		panic(errors.New("TLS certificate is not valid"))

	}
	//Use the miner's organization TLS certificate to get the report.
	//Create a TLS config that uses the server certificate as root CA so that future connections to the server can be verified.
	cert, _ := x509.ParseCertificate(minerCertBytes)
	minerTLSConfig := &tls.Config{RootCAs: x509.NewCertPool(), ServerName: "localhost"}
	minerTLSConfig.RootCAs.AddCert(cert)
	// Get the report via attested TLS channel
	reportBytes := http.HttpGet(minerTLSConfig, serverAddr+"/report")
	// Verify the report
	fmt.Println("Asking Report verification")
	if err := VerifyReport(reportBytes, minerCertBytes, expectedMeasurement); err != nil {
		panic(err)
		fmt.Println("Log receiver not attested")
	}
	fmt.Println("Report verification failed")
	// Create a TLS config that uses the server certificate as root CA so that future connections to the server can be verified.
	//cert, _ := x509.ParseCertificate(minerCertBytes)
	//tlsConfig = &tls.Config{RootCAs: x509.NewCertPool(), ServerName: "localhost"}
	//tlsConfig.RootCAs.AddCert(cert)
	fmt.Println("Remote Attestation completed")
	return minerCertBytes, reportBytes
}

/*This function verifies that 1)The report is signed by a functioning INTEL SGX TEE 2)That the tuple  */
func VerifyReport(reportBytes []byte, certBytes []byte, expectedMeasurement []byte) error {
	//1)Report validation via endorser starts here
	// Verify the report validity and extract the report data.
	report, err := eclient.VerifyRemoteReport(reportBytes)
	fmt.Println("report returned")
	if err == attestation.ErrTCBLevelInvalid {
		fmt.Printf("Warning: TCB level is invalid: %v\n%v\n", report.TCBStatus, tcbstatus.Explain(report.TCBStatus))
		fmt.Println("We'll ignore this issue in this sample. For an app that should run in production, you must decide which of the different TCBStatus values are acceptable for you to continue.")
	} else if err != nil {
		return err
	}
	//TODO HERE WE SHOULD ALSO ADD A FRESHNESS CHECK
	//2)Secure miner verification starts here
	// You can either verify the UniqueID or the tuple (SignerID, ProductID, SecurityVersion, Debug). Here we verify with the UNIQUEID.
	if !bytes.Equal([]byte(hex.EncodeToString(report.UniqueID)), expectedMeasurement) {
		return errors.New("the report measurement do not match the expected one")
	}
	//3)Veify that the report data matches the miner's TLS certificate
	//Hash the TLS certificate
	hash := sha256.Sum256(certBytes)
	//Verify if the TLS certificate hash is the same as the one in the report
	if !bytes.Equal(report.Data[:len(hash)], hash[:]) {
		//fmt.Printf(report.Data)
		return errors.New("report data does not match the certificate's hash")
	}
	return nil
}

/*Verify that the certificate valid*/
func validateCertificate(certBytes []byte) error {
	// Parse the certificate
	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return err
	}
	// Check if the public key is able to verify the certificate's signature
	if err := cert.CheckSignature(cert.SignatureAlgorithm, cert.RawTBSCertificate, cert.Signature); err != nil {
		return errors.New("certificate signature is not valid")
	}
	//Extract the public key from the certificate
	pubKey := cert.PublicKey.(*rsa.PublicKey)
	//Verifies if the public key of the certificate belongs to an authorized miner organization
	if !isFromExpectedOrganization(pubKey) {
		return errors.New("the certificate does not belong to an authorized miner organization")
	}
	return nil
}

/*Check if the certificate's public key belong to a known organization*/
func isFromExpectedOrganization(pubKey *rsa.PublicKey) bool {
	//TODO Implement a function that verifies if the public key belongs to an authorized miner organization
	return true
}
