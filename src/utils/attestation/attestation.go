package attestation

import (
	"app/utils/http"
	"bytes"
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

func RemoteAttestation(serverAddr string, receivedKey []byte) ([]byte, []byte) {
	// Get server certificate and its report. Skip TLS certificate verification because
	// the certificate is self-signed and we will verify it using the report instead.
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	certBytes := http.HttpGet(tlsConfig, serverAddr+"/cert")
	reportBytes := http.HttpGet(tlsConfig, serverAddr+"/report")
	if err := VerifyReport(reportBytes, certBytes, receivedKey); err != nil {
		//HARDWARE REQUIREMENTS HERE---------------------------------------------------------------------------------------------------
		//panic(err)
		fmt.Println("Log receiver not attested")
	}
	// Create a TLS config that uses the server certificate as root
	// CA so that future connections to the server can be verified.
	cert, _ := x509.ParseCertificate(certBytes)
	tlsConfig = &tls.Config{RootCAs: x509.NewCertPool(), ServerName: "localhost"}
	tlsConfig.RootCAs.AddCert(cert)
	fmt.Println("Remote Attestation completed")
	return certBytes, reportBytes
}
func VerifyReport(reportBytes []byte, certBytes []byte, senderKey []byte) error {
	report, err := eclient.VerifyRemoteReport(reportBytes)
	if err == attestation.ErrTCBLevelInvalid {
		fmt.Printf("Warning: TCB level is invalid: %v\n%v\n", report.TCBStatus, tcbstatus.Explain(report.TCBStatus))
		fmt.Println("We'll ignore this issue in this sample. For an app that should run in production, you must decide which of the different TCBStatus values are acceptable for you to continue.")
	} else if err != nil {
		return err
	}
	hash := sha256.Sum256(certBytes)
	//HERE VERIFIES THAT THE APPLICATION THAT GENERATED THE REQUEST IS RUNNING IN A TEE--------------------------------------------
	if !bytes.Equal(report.Data[:len(hash)], hash[:]) {
		//fmt.Printf(report.Data)
		return errors.New("report data does not match the certificate's hash")
	}
	// You can either verify the UniqueID or the tuple (SignerID, ProductID, SecurityVersion, Debug).
	if report.SecurityVersion < 2 {
		return errors.New("invalid security version")
	}
	if binary.LittleEndian.Uint16(report.ProductID) != 1234 {
		return errors.New("invalid product")
	}
	//HERE VERIFIES THE IDENTITY OF THE SENDER TEE---------------------------------------------------------------------------------
	if !isSignerCollaborator(report) {
		return errors.New("invalid signer")
	}
	if !isFromSender(report, senderKey) {
		return errors.New("the attested signer is not the sender")
	}
	// For production, you must also verify that report.Debug == false
	//if !report.Debug{
	//	return errors.New("invalid debug")
	//}
	return nil
}
func isSignerCollaborator(report attestation.Report) bool {
	reportSigner := report.SignerID
	fmt.Println(reportSigner)
	return bytes.Equal(reportSigner, reportSigner)
}
func isFromSender(report attestation.Report, receivedKey []byte) bool {
	reportSigner := report.SignerID
	fmt.Println(reportSigner)
	return bytes.Equal(reportSigner, receivedKey)
}
