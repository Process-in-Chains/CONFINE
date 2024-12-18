package attestation

import (
	"app/utils/http"
	"bytes"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/edgelesssys/ego/attestation"
	"github.com/edgelesssys/ego/attestation/tcbstatus"
	"github.com/edgelesssys/ego/eclient"
	"io/ioutil"
	"log"
	"os"
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
	if ValidateCertificate(minerCertBytes) != nil {
		panic(errors.New("Miner's TLS certificate is not valid"))
	}
	//Use the miner's organization TLS certificate to get the report.
	//Create a TLS config that uses the server certificate as root CA so that future connections to the server can be verified.
	cert, _ := x509.ParseCertificate(minerCertBytes)
	minerTLSConfig := &tls.Config{RootCAs: x509.NewCertPool(), ServerName: "localhost"}
	minerTLSConfig.RootCAs.AddCert(cert)
	// Get the report via attested TLS channel
	reportBytes := http.HttpGet(minerTLSConfig, serverAddr+"/report")
	// Verify the report
	if err := VerifyReport(reportBytes, minerCertBytes, expectedMeasurement); err != nil {
		panic(err)
	}
	log.Println("Remote attestation succesfully completed")
	return minerCertBytes, reportBytes
}

/*This function verifies that 1)The report is signed by a functioning INTEL SGX TEE 2)That the tuple  */
func VerifyReport(reportBytes []byte, certBytes []byte, expectedMeasurement []byte) error {
	//1)Report validation via endorser starts here
	// Verify the report validity and extract the report data.
	report, err := eclient.VerifyRemoteReport(reportBytes)
	log.Println("Attestation report successfully verified and decrypyed")
	if err == attestation.ErrTCBLevelInvalid {
		log.Printf("Warning: TCB level is invalid: %v\n%v\n", report.TCBStatus, tcbstatus.Explain(report.TCBStatus))
		//fmt.Println("We'll ignore this issue in this sample. For an app that should run in production, you must decide which of the different TCBStatus values are acceptable for you to continue.")
	} else if err != nil {
		return err
	}
	//TODO HERE WE SHOULD ALSO ADD A FRESHNESS CHECK
	//2)Secure miner verification starts here
	// You can either verify the UniqueID or the tuple (SignerID, ProductID, SecurityVersion, Debug). Here we verify with the UNIQUEID.
	if !bytes.Equal([]byte(hex.EncodeToString(report.UniqueID)), expectedMeasurement) {
		return errors.New("The report measurement do not match the expected one")
	}
	log.Println("Report data verification: measurement inside the report matches the expected one")
	//3)Veify that the report data matches the miner's TLS certificate
	//Hash the TLS certificate
	hash := sha256.Sum256(certBytes)
	//Verify if the TLS certificate hash is the same as the one in the report
	if !bytes.Equal(report.Data[:len(hash)], hash[:]) {
		//fmt.Printf(report.Data)
		return errors.New("Report data does not match the certificate's hash")
	}
	log.Println("Report data verification: TLS certificate inside the report matches the expected one")
	return nil
}

/*Verify that the certificate valid*/
func ValidateCertificate(certBytes []byte) error {
	// Parse the certificate
	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return err
	}
	// Check if the public key is able to verify the certificate's signature
	if err := cert.CheckSignature(cert.SignatureAlgorithm, cert.RawTBSCertificate, cert.Signature); err != nil {
		return errors.New("Certificate signature is not valid")
	}
	//Extract the public key from the certificate
	pubKey := cert.PublicKey.(*rsa.PublicKey)
	//Verifies if the public key of the certificate belongs to an authorized miner organization
	pubBytes, _ := x509.MarshalPKIXPublicKey(pubKey)

	if !isFromExpectedOrganization(base64.StdEncoding.EncodeToString(pubBytes)) {
		return errors.New("The certificate does not belong to an authorized miner organization")
	}
	log.Println("The provided certificate belongs to an authorized miner organization")
	return nil
}

/*Check if the certificate's public key belong to a known organization*/
func isFromExpectedOrganization(pubKey string) bool {
	//TODO Implement a function that verifies if the public key belongs to an authorized miner organization
	// Step 1: Read the JSON file
	file, err := os.Open("mining-data/provision-data/minerList.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return false
	}
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	// Step 2: Parse the JSON content into a slice of strings
	var keys []string
	if err := json.Unmarshal(byteValue, &keys); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return false
	}
	// Step 3: Decode each public key string into a public key object
	for _, keyStr := range keys {
		if keyStr == pubKey {
			return true
		}
	}
	return false
}
