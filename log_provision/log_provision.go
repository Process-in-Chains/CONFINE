package main

//go build -o logprovision log_provision/log_provision.go && ./logprovision -port 8087 -log event_log_TEST_2.xes
import (
	utilsHTTP "app/utils/attestation"
	utilsAttestation "app/utils/http"
	"context"
	"crypto/rsa"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"os"
	"strconv"
	"strings"
	//"strings"
	"app/utils/encryption"
	"crypto/x509"
	//"encoding/hex"
	"app/utils/xes"
	"flag"
	"fmt"
	"log"
	"net/url"
	//"encoding/json"
	"crypto/rand"
	"io/ioutil"
	"net/http"
)

var MYLOGPATH = "./mining-data/provision-data/process-01/event_log_TEST.xes"
var MYREFERENCE = "http://localhost:"

const PROCESSNAME = "process-01"

func main() {
	serverPort := flag.Int("port", 8081, "server address")
	provisionData := flag.String("log", "event_log_TEST.xes", "event log to provide")
	flag.Parse()
	MYREFERENCE = MYREFERENCE + strconv.Itoa(*serverPort)
	MYLOGPATH = "./mining-data/provision-data/process-01/" + *provisionData
	fmt.Println("Serving event log: ", MYLOGPATH)
	server := NewLogProvider(*serverPort)
	err2 := server.Start()
	if err2 != nil && err2 != http.ErrServerClosed {
		fmt.Printf("Failed to start server: %v\n", err2)
	}

	fmt.Println(server)
}

type LogProvider struct {
	port   int
	server *http.Server
}

func NewLogProvider(port int) *LogProvider {
	handler := http.NewServeMux()
	handler.HandleFunc("/logrequest", handleLogRequest)
	handler.HandleFunc("/tracelistrequest", handleTraceListRequest)
	s := &LogProvider{port: port}
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
	return s
}
func (s *LogProvider) Start() error {
	fmt.Printf("Log Provider listening on port %d...\n", s.port)
	return s.server.ListenAndServe()
}
func (s *LogProvider) Stop(ctx context.Context) error {
	fmt.Println("Shutting down Log Provider...")
	return s.server.Shutdown(ctx)
}
func getSegmentForm(secret string, key string, segmentNumber string, publicKey string, myreference string) url.Values {
	return url.Values{"key": {key}, "secret": {secret}, "segmentNumber": {segmentNumber}, "publicKey": {publicKey}, "myreference": {myreference}}
}
func getHeaderForm(decryptionToken string, bytes int, hashList []string, publicKey string) url.Values {
	// Create a map to hold the values
	data := map[string]interface{}{
		"bytes":         bytes,
		"hashList":      hashList,
		"segmentNumber": len(hashList),
	}
	// Convert the map to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	// Convert the JSON data to a string
	jsonString := string(jsonData)
	fmt.Println(jsonString)
	formData := url.Values{"header": {jsonString}, "decryptionToken": {decryptionToken}, "publicKey": {publicKey}}
	return formData
}
func sendSegments(symKey []byte, certBytes []byte, encryptedKey []byte, serverAddr string, publicKey string, myreference string) {
	files, err := ioutil.ReadDir("./mining-data/provision-data/" + PROCESSNAME + "/" + base64.StdEncoding.EncodeToString(symKey))
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "batch_") && strings.HasSuffix(file.Name(), ".xes") {
			// Estrai il numero di batch dal nome del file
			batchNumberStr := strings.TrimPrefix(file.Name(), "batch_")
			batchNumberStr = strings.TrimSuffix(batchNumberStr, ".xes")
			// Converti il numero di batch da stringa a intero
			batchNumber, err := strconv.Atoi(batchNumberStr)
			if err != nil {
				log.Printf("Errore nella conversione del numero di batch per il file %s: %v\n", file.Name(), err)
				continue
			}
			// Utilizza il numero di batch nell'iterazione corrente
			segmentNumber := batchNumber
			// Esegui le operazioni necessarie con il numero di batch
			fmt.Printf("File: %s, Numero di Batch: %d\n", file.Name(), segmentNumber)
			content, e := encryption.EncryptXES("./mining-data/provision-data/"+PROCESSNAME+"/"+base64.StdEncoding.EncodeToString(symKey)+"/"+file.Name(), symKey)
			if e != nil {
				fmt.Println(e)
				return
			}
			text := string(content)
			// Create a TLS config that uses the server certificate as root
			// CA so that future connections to the server can be verified.
			cert, _ := x509.ParseCertificate(certBytes)
			tlsConfig := &tls.Config{RootCAs: x509.NewCertPool(), ServerName: "localhost"}
			tlsConfig.RootCAs.AddCert(cert)
			utilsAttestation.HttpPOST(tlsConfig, serverAddr+"/secret", getSegmentForm(text, string(encryptedKey), batchNumberStr, publicKey, myreference))
		}
	}
	err = os.RemoveAll("./mining-data/provision-data/" + PROCESSNAME + "/" + base64.StdEncoding.EncodeToString(symKey))
	if err != nil {
		log.Fatal(err)
	}

}
func handleTraceListRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
		r.Body.Close()
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("New trace list request received")
	// On the server side, receive the public key bytes and deserialize it
	pubKeyBytes := r.Form.Get("publicKey")
	deserializedPubKey, err := x509.ParsePKIXPublicKey([]byte(pubKeyBytes))
	if err != nil {
		panic(err)
	}
	objectPublicKey := deserializedPubKey.(*rsa.PublicKey)
	//objectPublicKey := encryption.LoadPublicKeyFromFile("./public.pem")
	//serverAddr := r.Form.Get("logreceiver")
	//certBytes, _ := remoteAttestation(serverAddr)

	// Genera una nuova chiave simmetrica casuale
	symKey := encryption.GenerateRandomDecryptionToken()
	fmt.Println(len(symKey))
	fmt.Println(symKey)
	// Cripta la chiave simmetrica con RSA
	encryptedKey, err := rsa.EncryptPKCS1v15(rand.Reader, objectPublicKey, symKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(encryptedKey))
	traceSizeList, err := xes.GetTraceSize(MYLOGPATH)
	if err != nil {
		log.Fatal(err)
	}
	jsonBytes, err := json.Marshal(traceSizeList)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	//encryptedTraceSizeList, _ := encryption.EncryptDataWithSymetric(jsonBytes, symKey)
	fmt.Println(string(jsonBytes))
	//response := map[string]string{"traceList": string(encryptedTraceSizeList), "encryptedKey": string(symKey)}
	response := map[string]string{"traceList": string(jsonBytes), "encryptedKey": string(encryptedKey)}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Fatal(err)
	}
}
func handleLogRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
		r.Body.Close()
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("New log request received")
	serverAddr := r.Form.Get("logreceiver")
	fmt.Println(serverAddr)
	pubKeyBytes := r.Form.Get("publicKey")
	tracelistString := r.Form.Get("loglist")
	var traceList []string
	err := json.Unmarshal([]byte(tracelistString), &traceList)
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
	}
	batchSizeKB, err := strconv.Atoi(r.Form.Get("segmentSize"))
	if err != nil {
		panic(err)
	}
	deserializedPubKey, err := x509.ParsePKIXPublicKey([]byte(pubKeyBytes))
	if err != nil {
		panic(err)
	}
	_ = deserializedPubKey.(*rsa.PublicKey)
	certBytes, _ := utilsHTTP.RemoteAttestation(serverAddr, []byte(pubKeyBytes))
	// Genera una nuova chiave simmetrica casuale
	symKey := encryption.GenerateRandomDecryptionToken()
	// Cripta la chiave simmetrica con RSA
	encryptedKey, _ := encryption.EncryptSymmetricKey(symKey, "./public.pem")
	targetXes := xes.ReadXes(MYLOGPATH)
	//TODO: FILTER XES HERE BEFORE SENDING(LOOK AT TODOS IN log_request AND log_provision)---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	hashList, err := xes.SplitXESFile2(*targetXes, batchSizeKB, "./mining-data/provision-data/"+PROCESSNAME+"/"+base64.StdEncoding.EncodeToString(symKey))
	if err != nil {
		fmt.Println("Error splitting XES file:", err)
		return
	}
	fileInfo, err := os.Stat(MYLOGPATH)
	if err != nil {
		log.Fatal(err)
	}
	// Conversione in kilobyte
	fileSizeKB := int(fileInfo.Size() / 1024)
	myPublicKey := readPublicKey("./public.pem")
	myPubKeyBytes, err := x509.MarshalPKIXPublicKey(&myPublicKey)
	publicKeyString := base64.StdEncoding.EncodeToString(myPubKeyBytes)
	header := getHeaderForm(string(symKey), fileSizeKB, hashList, publicKeyString)
	tlsConfig := &tls.Config{RootCAs: x509.NewCertPool(), ServerName: "localhost"}
	cert, _ := x509.ParseCertificate(certBytes)
	tlsConfig.RootCAs.AddCert(cert)
	utilsAttestation.HttpPOST(tlsConfig, serverAddr+"/secret", header)
	sendSegments(symKey, certBytes, encryptedKey, serverAddr, publicKeyString, MYREFERENCE)
	fmt.Println("Sent log over attested TLS channel.")
}
func readPublicKey(filePath string) rsa.PublicKey {
	// Read the contents of the PEM file
	pemData, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading PEM file:", err)
	}
	// Decode the PEM data
	block, _ := pem.Decode(pemData)
	if block == nil {
		fmt.Println("Failed to decode PEM block")
	}
	// Parse the DER-encoded public key
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing public key:", err)
	}
	// Assert the type of the public key to RSA
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		fmt.Println("Public key is not an RSA key")
	}
	return *rsaPubKey
}
