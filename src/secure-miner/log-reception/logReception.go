package logreception

//go build -o logprovision provisioner/log-provision/log_provision.go && ./logprovision -port 8087 -log healthcare_newkeys/specialised_clinic_newkeys.xes -mergekey hospitalCaseId
import (
	"app/secure-miner/log-elaboration"
	logmanagement "app/secure-miner/log-management"
	"app/utils/collaborators"
	"app/utils/test"
	"app/utils/xes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"sync"
	"time"
	//"log"
	"github.com/edgelesssys/ego/enclave"
	"log"
)

/*Variable for debugging prints*/
var DEBUG = false

/*Variable for tests*/
var FIRSTTS = false
var FIRSTATT = false
var (
	mutex sync.Mutex
)

/*The log receiver object activated in the main*/
type LogReceiver struct {
	port             int
	server           *http.Server
	algorithm        string
	declareModelPath string
	tlsCertificate   []byte
	privateKey       crypto.PrivateKey
	logElaborator    *logelaboration.LogElaborator
}

/*Constructor function of the log receiver*/
func NewLogReceiver(port int) *LogReceiver {
	/*Generate a certificate and a private key for TLS with the provisioner.*/
	cert, priv := createCertificate()
	hash := sha256.Sum256(cert)
	s := &LogReceiver{port: port, tlsCertificate: cert, privateKey: priv}
	/*Generate the report (i.e., the attestation evidence) signed by the hardware's TEE. The report contains the hashed TLS certificate	*/
	//TODO REPORT GENERATION AND SIGNING SHOULD BE MOVED IN /report REQUETS. ADD NONCE IN THE PROTOCOL TO AVOID REPLAY ATTACKS.
	report, err := enclave.GetRemoteReport(hash[:])
	if err != nil {
		fmt.Println(err)
	}
	/*Define the handlers for the /cert, /report and /secret HTTP functions*/
	handler := http.NewServeMux()
	handler.HandleFunc("/cert", func(w http.ResponseWriter, r *http.Request) {
		if !FIRSTATT {
			if test.TEST_MODE {
				fmt.Println("TESTMODE - FIRST ATTESTATION AT:", time.Now().UnixMilli())
			}
			FIRSTATT = true
		}
		if DEBUG {
			fmt.Println("New certificate request received")
		}
		w.Write(cert)
	})
	handler.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		if DEBUG {
			fmt.Println("New report request received")
		}
		w.Write(report)
	})
	handler.HandleFunc("/secret", func(w http.ResponseWriter, r *http.Request) {
		secretLogHandler(w, r, s)
	})
	/*Use the certificate for the TLS configuration*/
	tlsCfg := tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{cert},
				PrivateKey:  priv,
			},
		},
	}
	/*Set up the HTTP server of the log receiver*/
	s.server = &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		Handler:   handler,
		TLSConfig: &tlsCfg,
	}
	s.logElaborator = logelaboration.NewLogElaborator()
	return s
}

/*Function to change the Algorithm of the LogReceiver*/
func (s *LogReceiver) SetAlgorithm(algorithm string) {
	s.algorithm = algorithm
}
func (s *LogReceiver) SetProcessModel(declareModelPath string) {
	s.declareModelPath = declareModelPath
}

func (s *LogReceiver) GetTLSCertificate() []byte {
	return s.tlsCertificate
}

func (s *LogReceiver) GetTLSPrivateKey() crypto.PrivateKey {
	return s.privateKey
}

/*Function to start the LogReceiver's server*/
func (s *LogReceiver) Start() error {
	fmt.Printf("Log Receiver listening on port %d...\n", s.port)
	return s.server.ListenAndServeTLS("", "")
}

/*Function to stop the LogReceiver's server*/
func (s *LogReceiver) Stop(ctx context.Context) error {
	fmt.Println("Shutting down Log Receiver...")
	return s.server.Shutdown(ctx)
}

/*Function to create the TLS certificate using the private key*/
func createCertificate() ([]byte, crypto.PrivateKey) {
	template := &x509.Certificate{
		SerialNumber: &big.Int{},
		Subject:      pkix.Name{CommonName: "localhost"},
		NotAfter:     time.Now().Add(time.Hour),
		DNSNames:     []string{"localhost"},
	}
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	cert, _ := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	return cert, priv
}

/*Function that handles the log segments sent by providers*/
func secretLogHandler(w http.ResponseWriter, r *http.Request, logReceiver *LogReceiver) {
	/*The LogReceiver only handles POST requests*/
	if r.Method == "POST" {
		//r.Body = http.MaxBytesReader(w, r.Body, 50000)
		maxRequestSize := int64(70 * 1024 * 1024) // 10 MB in bytes
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)
		err := r.ParseForm()
		if err != nil {
			println(err.Error())
			println("HERE WE HAVE A WRONG FORM", r.Form.Encode())
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			response, _ := prepareResponse(false)
			w.Write([]byte(response))
			return
		}
		r.Body.Close()
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		response, _ := prepareResponse(false)
		w.Write([]byte(response))
		return
	}
	/*If it's an header message...*/
	if r.Form.Has("header") {
		if DEBUG {
			//bytes, segmentNumber, hashList := readHeader(headerString)
			fmt.Println("New header received: ", r.Form.Get("header"))
		}
	} else { /*If the message it's a log segment...*/
		/*Check for the tests*/
		if !FIRSTTS {
			if test.TEST_MODE {
				println("TESTMODE - FIRST SEGMENT RECEIVED AT:", time.Now().UnixMilli())
			}
			FIRSTTS = true
		}
		/*Parse the attributes of the POST form*/
		eventLog := r.Form.Get("secret")
		_ = r.Form.Get("segmentNumber")
		senderReference := r.Form.Get("myreference")
		/*Parse the XES*/
		xesSegment := xes.ParseXes([]byte(eventLog))
		mutex.Lock()
		/*Read the json containing the map of the traces received */
		writtenTraceMap, err := ioutil.ReadFile("mining-data/consumption-data/process-01/traceMap.json")
		readTraceMap := make(map[string]map[string]bool)
		/*Parse the JSON*/
		err = json.Unmarshal(writtenTraceMap, &readTraceMap)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}
		/*Read the json containing the references of the provisioners*/
		reference, error := collaborators.GetReference(senderReference)
		if error != nil {
			log.Fatalf("Error getting references: %v", error)
		}
		/*Get the algorithm set in the main*/
		algorithm := logReceiver.algorithm
		declareModelPath := logReceiver.declareModelPath
		/*Call the trace handler*/
		logmanagement.HandleSegment(*xesSegment, "process-01", senderReference, reference.MergeKey, readTraceMap, algorithm, declareModelPath, *logReceiver.logElaborator)
		mutex.Unlock()
		/*Send response to the provisioner*/
		response, _ := prepareResponse(true)
		w.Write([]byte(response))
		if DEBUG {
			fmt.Println("New event log received and stored in the trusted memory zone")
		}
	}
}

/*Response message to the provisioner*/
func prepareResponse(success bool) (string, bool) {
	message := ""
	if success {
		message = "Log received successfully"
	} else {
		message = "Log reception failed"
	}
	response := Response{
		Success: success,
		Msg:     message}
	jsonData, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error:", err)
		return "", false
	}
	jsonString := string(jsonData)
	return jsonString, true
}

/*Function that reads the header of the log segments*/
func readHeader(header string) (int, int, []string) {
	// Declare a map to hold the decoded data
	var data map[string]interface{}
	// Convert the JSON string to the map
	err := json.Unmarshal([]byte(header), &data)
	if err != nil {
		log.Fatal(err)
	}
	// Access the values from the map
	bytes := int(data["bytes"].(float64))
	segmentNumber := int(data["segmentNumber"].(float64))
	hashList := data["hashList"].([]interface{})
	hashStrings := make([]string, len(hashList))
	for i, v := range hashList {
		hashStrings[i] = v.(string)
	}
	return bytes, segmentNumber, hashStrings
}

type Response struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

func createFolderIfNotExists(folderPath string) error {
	// Verify if the folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// If the folder does not extist, create it
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return err
		}
		log.Printf("La cartella %s è stata creata.", folderPath)
	}
	return nil
}
