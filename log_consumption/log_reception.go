package logreception

//ego-go build -buildvcs=false && ego sign enclave.json && OE_SIMULATION=1 ego run ./app -segsize 2000 -port 8080 -test true
//ego sign enclave.json
//ego run ./app
import (
	prosessMiningAlgorithms "app/log_consumption/miningAlgorithms"
	"app/utils/encryption"
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
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
	//"log"
	"github.com/edgelesssys/ego/enclave"
	"log"
)

var DEBUG = false
var FIRSTTS = false
var FIRSTATT = false
var FIRSTCOMP = false
var (
	mutex sync.Mutex
)

type LogReceiver struct {
	port   int
	server *http.Server
}

func NewLogReceiver(port int) *LogReceiver {
	cert, priv := createCertificate()
	hash := sha256.Sum256(cert)
	report, err := enclave.GetRemoteReport(hash[:])
	if err != nil {
		fmt.Println(err)
	}
	handler := http.NewServeMux()
	handler.HandleFunc("/cert", func(w http.ResponseWriter, r *http.Request) {
		if DEBUG {
			fmt.Println("New certificate request received")
		}
		w.Write(cert)
	})
	handler.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		if !FIRSTATT {
			fmt.Println("TESTMODE - FIRST ATTESTATION AT:", time.Now().UnixMilli())
			FIRSTATT = true
		}
		if DEBUG {
			fmt.Println("New report request received")
		}
		w.Write(report)
	})
	handler.HandleFunc("/secret", secretLogHandler)
	tlsCfg := tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{cert},
				PrivateKey:  priv,
			},
		},
	}
	s := &LogReceiver{port: port}
	s.server = &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		Handler:   handler,
		TLSConfig: &tlsCfg,
	}

	return s
}
func (s *LogReceiver) Start() error {
	fmt.Printf("Log Receiver listening on port %d...\n", s.port)
	return s.server.ListenAndServe()
}
func (s *LogReceiver) Stop(ctx context.Context) error {
	fmt.Println("Shutting down Log Receiver...")
	return s.server.Shutdown(ctx)
}

/*
*

	func waitForInterrupt(ctx context.Context) {
	    c := make(chan os.Signal, 1)
	    signal.Notify(c, os.Interrupt)
	    select {
	    case <-c:
	        fmt.Println("\nReceived interrupt signal")
	    case <-ctx.Done():
	    }
	}

*
*/
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

// {'messageType':'0,1,2', 'secret':{eventLog, index}, 'key':decryptionToken/"", 'header':{bytes:_,hashList:_,segmentNumber:},}
func secretLogHandler(w http.ResponseWriter, r *http.Request) {
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
	if r.Form.Has("header") {
		if DEBUG {
			//bytes, segmentNumber, hashList := readHeader(headerString)
			fmt.Println("New header received: ", r.Form.Get("header"))
		}
	} else {
		if !FIRSTTS {
			println("TESTMODE - FIRST SEGMENT RECEIVED AT:", time.Now().UnixMilli())
			FIRSTTS = true
		}
		eventLog := r.Form.Get("secret")
		_ = r.Form.Get("segmentNumber")
		senderPublicKey := r.Form.Get("publicKey")
		senderReference := r.Form.Get("myreference")
		privateKeyPath := "./private.pem"
		privateKey, err := encryption.LoadPrivateKeyFromFile(privateKeyPath)
		if err != nil {
			log.Fatal("Error loading private key:", err)
		}
		key := r.Form.Get("key")
		// Decrypt the symmetric key using RSA
		symKey, err := encryption.DecryptSymmetricKey([]byte(key), privateKey)
		if err != nil {
			log.Fatal("Error decrypting symmetric key:", err)
		}
		// Decrypt the XES data using AES
		decryptedData, err := encryption.DecryptXES([]byte(eventLog), symKey)
		if err != nil {
			log.Fatal("Error decrypting XES data:", err)
		}
		//err = createFolderIfNotExists("data/enclave/" + base64.StdEncoding.EncodeToString(symKey))
		//if err != nil {
		//	return
		//}
		//logFilenameDecrypted := "data/enclave/" + base64.StdEncoding.EncodeToString(symKey) + "/" + segmentNumber + ".xes"
		//HERE YOU HAVE TO PARSE decryptedData to turn it into a XES object. Than
		//you can extract traces and save them or marge them
		xesSegment := xes.ParseXes(decryptedData)
		//Read the traces from the traceMap file
		mutex.Lock()
		writtenTraceMap, err := ioutil.ReadFile("mining-data/consumption-data/process-01/traceMap.json")
		readTraceMap := make(map[string]map[string]bool)
		err = json.Unmarshal(writtenTraceMap, &readTraceMap)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}
		handleTrace(*xesSegment, "process-01", senderPublicKey, senderReference, readTraceMap)
		mutex.Unlock()
		//err = ioutil.WriteFile(logFilenameDecrypted, decryptedData, 0644)
		//if err != nil {
		//	fmt.Println("Error writing to file:", err)
		//	response, _ := prepareResponse(false)
		//	w.Write([]byte(response))
		//	return
		//}
		response, _ := prepareResponse(true)
		w.Write([]byte(response))
		if DEBUG {
			fmt.Println("New event log received and stored in the trusted memory zone")
		}
	}
}
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
	// Verifica se la cartella esiste
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// La cartella non esiste, quindi la crea
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return err
		}
		log.Printf("La cartella %s è stata creata.", folderPath)
	}
	return nil
}
func handleTrace(eventLog xes.XES, processName string, publicKey string, myreference string, traceMap map[string]map[string]bool) {
	for _, trace := range eventLog.Traces {
		traceId, _ := trace.GetId()
		if _, ok := traceMap[traceId]; ok && !traceMap[traceId][myreference] {
			traceMap[traceId][myreference] = true
			if allValuesTrue(traceMap[traceId]) {
				if DEBUG {
					fmt.Println("It's time to mine Trace: ", traceId)
				}
				//TODO: FIX HERE TO ADD CASE IF MINER DOESN'T HAVE ANY TRACE(LOOK AT TODOS IN log_request AND log_provision)-------------------------------------------------------------------------------------------------------------
				my := xes.ReadXes("/mining-data/provision-data/process-01/trace_" + traceId + "/trace_" + traceId + ".xes")
				mergedTrace, _ := xes.MergeTraces(trace, my.Traces[0])
				//iterate over files in the data/enclave/process_01 folder
				files, err := os.ReadDir("/mining-data/consumption-data/" + processName + "/trace_" + traceId)
				if err != nil {
					log.Fatal(err)
				}
				for _, traceFile := range files {
					currentTrace := xes.ReadXes("mining-data/consumption-data/" + processName + "/trace_" + traceId + "/" + traceFile.Name())
					if DEBUG {
						fmt.Println("Merging: ", "mining-data/consumption-data/"+processName+"/trace_"+traceId+"/"+traceFile.Name())
					}
					mergedTrace, _ = xes.MergeTraces(mergedTrace, currentTrace.Traces[0])
				}
				if DEBUG {
					for _, ev := range mergedTrace.Events {
						fmt.Println(ev.GetAttributeValue("concept:name").Value, ev.Timestamp.Value)
					}
				}
				tr := []xes.Trace{}
				tr = append(tr, mergedTrace)
				x := xes.XES{Traces: tr}
				//MINE(mergedTrace)------------------------------HEURISTC------------------------------------------------------
				if !FIRSTCOMP {
					fmt.Println("TESTMODE - FIRST COMPUTATION AT:", time.Now().UnixMilli())
					FIRSTCOMP = true
				}

				prosessMiningAlgorithms.HeuristicMiner(x.XesToSlices(), processName)

			} else {
				if DEBUG {
					fmt.Println("Write in memory Trace: ", traceId)
				}
				filename := "mining-data/consumption-data/" + processName + "/trace_" + traceId + "/trace_" + url.PathEscape(fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))) + ".xes"
				//byteTrace := trace.traceToByte()
				byteTrace := trace.TraceToByte()
				err := os.MkdirAll(filepath.Dir(filename), os.ModePerm)
				if err != nil {
					log.Fatal(err)
				}

				err = os.WriteFile(filename, byteTrace, 0644)
				if err != nil {
					log.Fatal(err)
				}
				//}
			}

		} else {
			fmt.Printf("'%s' is not a key in the map or sender not expected\n", traceId)
		}

	}

	jsonData, err := json.MarshalIndent(traceMap, "", "  ")
	if err != nil {
		fmt.Println("Error converting JSON:", err)
		return
	}
	// WRITE Json in file
	err = ioutil.WriteFile("./mining-data/consumption-data/"+processName+"/traceMap.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Errore nella scrittura del file JSON:", err)
		return
	}
	if allTracesCompleted(traceMap) {
		fmt.Println("TESTMODE - TEST ENDED AT: ", time.Now().UnixMilli())
		test.STOPMONITORING = true
	}

}
func allValuesTrue(myMap map[string]bool) bool {
	for _, value := range myMap {
		if !value {
			return false
		}
	}
	return true
}

func allTracesCompleted(traceMap map[string]map[string]bool) bool {
	for _, value := range traceMap {
		if !allValuesTrue(value) {
			return false
		}
	}
	return true
}
