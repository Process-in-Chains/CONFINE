package main

//CGO_CFLAGS=-I/opt/ego/include CGO_LDFLAGS=-L/opt/ego/lib go build -o logprovision provisioner/log-provision/log_provision.go && ./logprovision -port 8087 -log testing_logs/motivating/hospital.xes -mergekey concept:name -measurement `ego uniqueid app` -skipattestation false
//CGO_CFLAGS=-I/opt/ego/include CGO_LDFLAGS=-L/opt/ego/lib go build -o logprovision provisioner/log-provision/log_provision.go && ./logprovision -port 8088 -log testing_logs/motivating/pharma.xes -mergekey concept:name -measurement `ego uniqueid app` -skipattestation false
//CGO_CFLAGS=-I/opt/ego/include CGO_LDFLAGS=-L/opt/ego/lib go build -o logprovision provisioner/log-provision/log_provision.go && ./logprovision -port 8089 -log testing_logs/motivating/specialized.xes -mergekey concept:name -measurement `ego uniqueid app` -skipattestation false
import (
	utilsAttestation "app/utils/attestation"
	"app/utils/encryption"
	utilsHTTP "app/utils/http"
	"app/utils/xes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var MYLOGPATH = "./mining-data/provision-data/process-01/event_log_TEST.xes"
var MYREFERENCE = "http://localhost:"
var MYMERGEKEY = "concept:name"
var SKIPATTESTATION bool
var EXPECTEDMEASUREMENT = ""

const PROCESSNAME = "process-01"

/*Main function to run the log provider*/
func main() {
	serverPort := flag.Int("port", 8094, "server address")
	provisionData := flag.String("log", "event_log_TEST.xes", "event log to provide")
	mrgkey := flag.String("mergekey", "concept:name", "merge key to be used when merging traces")
	skipattestation := flag.String("skipattestation", "", "set to false if your secure miner is running in simulation")
	measurement := flag.String("measurement", "", "expected measurement of the miner")
	flag.Parse()
	SKIPATTESTATION, _ = strconv.ParseBool(*skipattestation)
	EXPECTEDMEASUREMENT = *measurement
	MYMERGEKEY = *mrgkey
	MYREFERENCE = MYREFERENCE + strconv.Itoa(*serverPort)
	MYLOGPATH = "./mining-data/provision-data/process-01/" + *provisionData
	log.Println("Serving event log: ", MYLOGPATH)
	server := NewLogProvider(*serverPort)
	err2 := server.Start()
	if err2 != nil && err2 != http.ErrServerClosed {
		fmt.Printf("Failed to start server: %v\n", err2)
	}
	fmt.Println(server)
}

/*Log provider object*/
type LogProvider struct {
	port   int
	server *http.Server
}

// Build a new log provider
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

// Start the log provider
func (s *LogProvider) Start() error {
	log.Printf("Log Server listening on port %d...\n", s.port)
	return s.server.ListenAndServe()
}

// Stop the log provider
func (s *LogProvider) Stop(ctx context.Context) error {
	fmt.Println("Shutting down Log Provider...")
	return s.server.Shutdown(ctx)
}

// Function to get the form for the segment
func getSegmentForm(secret string, segmentNumber string, myreference string) url.Values {
	return url.Values{"secret": {secret}, "segmentNumber": {segmentNumber}, "myreference": {myreference}}
}

// Function to get the form for the header
func getHeaderForm(bytes int, hashList []string) url.Values {
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
	formData := url.Values{"header": {jsonString}}
	return formData
}

/**Function that sends the segments to the secure miner**/
func sendSegments(symKey []byte, certBytes []byte, serverAddr string, myreference string) {
	files, err := ioutil.ReadDir("./mining-data/provision-data/" + PROCESSNAME + "/" + base64.StdEncoding.EncodeToString(symKey))
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "batch_") && strings.HasSuffix(file.Name(), ".xes") {
			batchNumberStr := strings.TrimPrefix(file.Name(), "batch_")
			batchNumberStr = strings.TrimSuffix(batchNumberStr, ".xes")
			batchNumber, err := strconv.Atoi(batchNumberStr)
			if err != nil {
				log.Printf("Errore nella conversione del numero di batch per il file %s: %v\n", file.Name(), err)
				continue
			}
			segmentNumber := batchNumber
			content, e := ioutil.ReadFile("./mining-data/provision-data/" + PROCESSNAME + "/" + base64.StdEncoding.EncodeToString(symKey) + "/" + file.Name())
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
			utilsHTTP.HttpPOST(tlsConfig, serverAddr+"/secret", getSegmentForm(text, batchNumberStr, myreference))
			log.Printf("Segment %d sent to the Secure Miner at %s\n", segmentNumber, serverAddr)
		}
	}
	err = os.RemoveAll("./mining-data/provision-data/" + PROCESSNAME + "/" + base64.StdEncoding.EncodeToString(symKey))
	if err != nil {
		log.Fatal(err)
	}

}

// Function to handle the trace list requests from the miner
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
	fmt.Println()
	log.Println("New trace list request received")
	tlsCertBytes := []byte(r.Form.Get("tlsCert"))
	utilsAttestation.ValidateCertificate(tlsCertBytes)
	cert, err := x509.ParseCertificate(tlsCertBytes)
	if err != nil {
		panic(err)
	}
	objectPublicKey := cert.PublicKey.(*rsa.PublicKey)
	symKey := encryption.GenerateRandomDecryptionToken()
	encryptedKey, err := rsa.EncryptPKCS1v15(rand.Reader, objectPublicKey, symKey)
	encryptedKeyBase64 := base64.StdEncoding.EncodeToString(encryptedKey)
	if err != nil {
		log.Fatal(err)
	}
	traceSizeList, err := xes.GetTraceSize(MYLOGPATH, MYMERGEKEY)
	if err != nil {
		log.Fatal(err)
	}
	//TODO This should be encrypted using the symmetric key
	jsonBytes, err := json.Marshal(traceSizeList)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	encryptedTraceList, err := encryption.EncryptDataWithSymetric(jsonBytes, symKey)
	if err != nil {
		fmt.Println("Error encrypting the trace list")
	}
	encryptedTraceListBase64 := base64.StdEncoding.EncodeToString(encryptedTraceList)
	response := map[string]string{"traceList": encryptedTraceListBase64, "encryptedKey": encryptedKeyBase64}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Trace list successfully sent")
}

// Function to handle the log requests from the miner
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
	serverAddr := r.Form.Get("logreceiver")
	log.Println("New event log request received for Secure Miner's log receiver at " + serverAddr)
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
	certBytes := []byte("")
	if !SKIPATTESTATION {
		log.Println("Starting the remote attestation of the miner")
		certBytes, _ = utilsAttestation.RemoteAttestation(serverAddr, []byte(EXPECTEDMEASUREMENT))

	} else {
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		certBytes = utilsHTTP.HttpGet(tlsConfig, serverAddr+"/cert")
	}
	//TODO: Here you should remove sym key and the encryption process. No Encryption is needed. Now we have TLS
	symKey := encryption.GenerateRandomDecryptionToken()
	targetXes := xes.ReadXes(MYLOGPATH)
	//TODO: FILTER XES HERE BEFORE SENDING(LOOK AT TODOS IN log-request AND log-provision)---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	hashList, err := xes.SplitXESFile2(*targetXes, batchSizeKB, "./mining-data/provision-data/"+PROCESSNAME+"/"+base64.StdEncoding.EncodeToString(symKey))
	if err != nil {
		fmt.Println("Error segmenting the XES file:", err)
		return
	}
	fileInfo, err := os.Stat(MYLOGPATH)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Log segmentation completed with segment size " + strconv.Itoa(batchSizeKB) + "KB" + ". Number of segments: " + strconv.Itoa(len(hashList)))
	fileSizeKB := int(fileInfo.Size() / 1024)
	header := getHeaderForm(fileSizeKB, hashList)
	tlsConfig := &tls.Config{RootCAs: x509.NewCertPool(), ServerName: "localhost"}
	cert, _ := x509.ParseCertificate(certBytes)
	tlsConfig.RootCAs.AddCert(cert)
	utilsHTTP.HttpPOST(tlsConfig, serverAddr+"/secret", header)
	sendSegments(symKey, certBytes, serverAddr, MYREFERENCE)
	log.Println("Data transmission completed")
}
