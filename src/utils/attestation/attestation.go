package attestation

import (
	"app/utils/http"
	"bytes"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/edgelesssys/ego/attestation"
	"github.com/edgelesssys/ego/attestation/tcbstatus"
	"github.com/edgelesssys/ego/eclient"
)

/*
Se usiamo self signed certificates non abbiamo modo di verificare l'identità del miner. In questo caso, dovremmo quindi inserire un'identity evidence nel report. Questo può essere fatto con una coppia di chiavi pubbliche e private randomiche generate nella TEE. La chiave pubblica può essere inserita nel report e la chiave privata può essere usata per criptare i segmenti del log.
Se usiamo CA signed certificates abbiamo un modo per verificare l'identità del miner. Pero, non possiamo usare la chiave del certificato, poiche generata esternamente alla TEE. Qundi dovremmo crearla a runtime ed includere la sua pubblica nel report.

*/

/*Method invoked to execute the remote attestation of the Secure Miner hosted in a given address*/
func RemoteAttestation(serverAddr string, receivedKey []byte) ([]byte, []byte) {
	// Get server certificate. Skip TLS certificate verification because the certificate is self-signed and we will verify it using the report instead.
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	// Get miner's organization TLS certificate
	minerCertBytes := http.HttpGet(tlsConfig, serverAddr+"/cert")
	// Validate the certificate
	if validateCertificate(minerCertBytes) != nil {
		panic(errors.New("TLS certificate is not valid"))

	}
	//Verify certificate validity
	//Use the miner's organization TLS certificate to get the report.
	//Create a TLS config that uses the server certificate as root CA so that future connections to the server can be verified.
	cert, _ := x509.ParseCertificate(minerCertBytes)
	minerTLSConfig := &tls.Config{RootCAs: x509.NewCertPool(), ServerName: "localhost"}
	minerTLSConfig.RootCAs.AddCert(cert)
	// Get the report via attested TLS channel
	reportBytes := http.HttpGet(minerTLSConfig, serverAddr+"/report")
	// Verify the report
	if err := VerifyReport(reportBytes, minerCertBytes, receivedKey); err != nil {
		//TODO HARDWARE REQUIREMENT HERE. IF THE REMOTE ATTESTATION FAILS, THE SECURE MINER SHOULD BE STOPPED.
		//panic(err)
		fmt.Println("Log receiver not attested")
	}
	// Create a TLS config that uses the server certificate as root CA so that future connections to the server can be verified.
	//cert, _ := x509.ParseCertificate(minerCertBytes)
	//tlsConfig = &tls.Config{RootCAs: x509.NewCertPool(), ServerName: "localhost"}
	//tlsConfig.RootCAs.AddCert(cert)
	fmt.Println("Remote Attestation completed")
	return minerCertBytes, reportBytes
}

/*This function verifies that 1)The report is signed by a functioning INTEL SGX TEE 2)That the tuple  */
func VerifyReport(reportBytes []byte, certBytes []byte, senderKey []byte) error {
	//1)Report validation via endorsere starts here
	// Verify the report validity and extract the report data.
	report, err := eclient.VerifyRemoteReport(reportBytes)
	if err == attestation.ErrTCBLevelInvalid {
		fmt.Printf("Warning: TCB level is invalid: %v\n%v\n", report.TCBStatus, tcbstatus.Explain(report.TCBStatus))
		fmt.Println("We'll ignore this issue in this sample. For an app that should run in production, you must decide which of the different TCBStatus values are acceptable for you to continue.")
	} else if err != nil {
		return err
	}
	//2)Secure miner verification starts here
	// You can either verify the UniqueID or the tuple (SignerID, ProductID, SecurityVersion, Debug). Here we verify the tuple.
	//Check security version
	if report.SecurityVersion < 2 {
		return errors.New("invalid security version")
	}
	//Check product ID
	if binary.LittleEndian.Uint16(report.ProductID) != 1234 {
		return errors.New("invalid product")
	}
	//Check the developer that signed the trusted app
	if !isFromExpectedDeveloper(report, senderKey) {
		return errors.New("the attested signer is not the expected developer")
	}
	//Check if the sender is running in simulation mode
	if !report.Debug {
		return errors.New("invalid debug")
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

/*Check that the signer of the trusted app is the expected developer*/
func isFromExpectedDeveloper(report attestation.Report, receivedKey []byte) bool {
	reportSigner := report.SignerID
	fmt.Println(reportSigner)
	return bytes.Equal(reportSigner, receivedKey)
}

/*Check if the certificate's public key belong to a known organization*/
func isFromExpectedOrganization(pubKey *rsa.PublicKey) bool {
	//TODO Implement a function that verifies if the public key belongs to an authorized miner organization
	return true
}
