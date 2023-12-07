package logrequest

import (
	"app/utils/xes"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"time"
)

// package logreception
const STOPWRITING = true

// go build -o logrequest log_request/log_request.go && ./logrequest -segsize 2000 -port 8080 -test true
func LogRequest(processName string, receiverPort string, segmentsize int) {
	println("TESTMODE - INITIALIZATION STARTED AT:", time.Now().UnixMilli())
	//Initalize and write the trace map
	//TODO:FIX HERE TO MAKE THE MINER NOT OWNER OF TRACES3
	log_path := "./mining-data/provision-data/" + processName + "/event_log.xes"
	//TODO:FIX HERE TO MAKE THE MINER NOT OWNER OF TRACES
	eventLog := xes.ReadXes(log_path)
	globalTracesMap := make(map[string]map[string]bool)
	if STOPWRITING {
		_ = os.MkdirAll("./mining-data/consumption-data/"+processName+"/miningMetadata", os.ModePerm)
		_ = os.WriteFile("./mining-data/consumption-data/"+processName+"/miningMetadata/map.json", []byte("{}"), 0644)

		//TODO:FIX HERE TO MAKE THE MINER NOT OWNER OF TRACES. THIS FOR SHOULD BE SKIPPED IN THAT CASE
		for _, trace := range eventLog.Traces {
			traceId, _ := trace.GetId()
			traceMap := map[string]bool{
				"0": true,
			}
			globalTracesMap[traceId] = traceMap
			byteTrace := trace.TraceToByte()
			//TODO: THIS THING HERE SHOULD BE DONE ALSO WHEN NOT KNOWN TRACE ARRIVE FROM PROVISIONERS _X_
			_ = os.MkdirAll("./mining-data/provision-data/"+processName+"/trace_"+traceId, os.ModePerm)
			_ = os.MkdirAll("./mining-data/consumption-data/"+processName+"/trace_"+traceId, os.ModePerm)
			_ = os.WriteFile("./mining-data/provision-data/"+processName+"/trace_"+traceId+"/trace_"+traceId+".xes", byteTrace, 0644)
		}

		jsonData, err := json.MarshalIndent(globalTracesMap, "", "  ")
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
	}

	//Read collaborator refeferences
	writtendata, err := ioutil.ReadFile("./mining-data/collaborators/" + processName + "/references.json")
	var references []map[string]interface{}
	err = json.Unmarshal(writtendata, &references)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	//Read tracemap
	writtenTraceMap, err := ioutil.ReadFile("./mining-data/consumption-data/" + processName + "/traceMap.json")
	readTraceMap := make(map[string]map[string]bool)
	err = json.Unmarshal(writtenTraceMap, &readTraceMap)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	//Prepare header exchange
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	publicKey := readPublicKey()
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		log.Fatal(err)
	}
	//Iterate over collaborators and get their tracelists
	for _, item := range references {

		httpReference, _ := item["http_reference"].(string)
		//publicKey, _ := item["public_key"].(string)
		response := httpPOST(tlsConfig, httpReference+"/tracelistrequest", string(pubKeyBytes), "http://localhost:"+receiverPort, segmentsize, "")
		//privateKey, err := loadPrivateKeyFromFile("./private.pem")
		if err != nil {
			log.Fatal(err)
		}
		var wrapperResponse map[string]string
		err = json.Unmarshal(response, &wrapperResponse)
		if err != nil {
			fmt.Println(err)
		}
		var responseJson map[string]string
		err = json.Unmarshal([]byte(wrapperResponse["traceList"]), &responseJson)
		if err != nil {
			log.Fatal(err)
		}
		//TODO: FIX BELOW IF YOU WANT TO INCLUDE TRACES NOT OWNED BY THE MINER (LOOK AT TODOS IN log_reception AND log_provision )----------------------------------------------------------------------------------------------------------------------------------
		for trId, _ := range responseJson {
			if _, ok := readTraceMap[trId]; ok {
				readTraceMap[trId][httpReference] = !STOPWRITING
			}
			//TODO: IF OK=FALSE (LOOK ABOVE) THEN THE TRACE IS NOT OWNED BY THE MINER SO A NEW ENTRY SHOULD BE ADDED IN THE TRACEMAP WITH THE ID OF THE TRACE. YOU SHOULD DO ALSO ALL THE STUFF IN _X_
		}
		jsonData, err := json.MarshalIndent(readTraceMap, "", "  ")
		if err != nil {
			fmt.Println("Error converting JSON:", err)
			return
		}
		// Save the updated tracemap
		err = ioutil.WriteFile("./mining-data/consumption-data/"+processName+"/traceMap.json", jsonData, 0644)
		if err != nil {
			fmt.Println("Errore nella scrittura del file JSON:", err)
			return
		}

	}
	//TODO: TWO POSSIBLE FIX BELOW: 1)REQUEST ALL THE TRACES IN THE TRACE MAP, THEN PROVISIONERS IF DONT HAVE THE TRACE, JUST IGNORE THE REQUESTED ID. 2) DO A CUSTOM TRACELIST FOR EACH LOG_REQUEST
	traceList := make([]string, 0, len(readTraceMap))
	for key := range readTraceMap {
		traceList = append(traceList, key)
	}
	traceListByte, err := json.Marshal(traceList)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	for _, item := range references {
		httpReference, _ := item["http_reference"].(string)
		go httpPOST(tlsConfig, httpReference+"/logrequest", string(pubKeyBytes), "http://localhost:"+receiverPort, segmentsize, string(traceListByte))
	}
}
func httpGet(tlsConfig *tls.Config, url string) []byte {
	client := http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func httpPOST(tlsConfig *tls.Config, posturl string, publicKey string, logreceiver string, segmentSize int, loglist string) []byte {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	// Prepare the form data
	formData := url.Values{"logreceiver": {logreceiver}, "publicKey": {publicKey}, "segmentSize": {strconv.Itoa(segmentSize)}, "loglist": {loglist}}
	// Send the POST request
	response, err := client.PostForm(posturl, formData)
	if err != nil {
		fmt.Println("POST request failed:", err)
		return nil
	}
	defer response.Body.Close()
	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return nil
	}
	// Print the response
	//fmt.Println("Response:", string(body))
	defer response.Body.Close()
	return body
}

func readPublicKey() rsa.PublicKey {
	// Read the contents of the PEM file
	pemData, err := ioutil.ReadFile("./public.pem")
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

// Function to convert RSA public key to string representation
func rsaPubKeyToString(pubKey *rsa.PublicKey) string {
	return fmt.Sprintf("RSA Public Key\nN: %s\nE: %d", pubKey.N.String(), pubKey.E)
}
func sortByFirstValue(arr [][]string) [][]string {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i][0] < arr[j][0]
	})
	return arr
}
func loadPrivateKeyFromFile(keyPath string) (*rsa.PrivateKey, error) {
	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("invalid private key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
