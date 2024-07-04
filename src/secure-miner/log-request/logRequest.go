package logrequest

import (
	"app/utils/collaborators"
	"app/utils/test"
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
	"strconv"
	"time"
)

/*This function contains the logic of the Secure Miner's Log Requester*/
func LogRequest(processName string, receiverPort string, segmentsize int) {
	if test.TEST_MODE {
		println("TESTMODE - INITIALIZATION STARTED AT:", time.Now().UnixMilli())
	}
	//Initalize and write the trace map
	globalTracesMap := make(map[string]map[string]bool)
	_ = os.MkdirAll("./mining-data/consumption-data/"+processName+"/miningMetadata", os.ModePerm)
	_ = os.WriteFile("./mining-data/consumption-data/"+processName+"/miningMetadata/map.json", []byte("{}"), 0644)
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
	//Read collaborator refeferences
	references, err := collaborators.GetReferences()
	if err != nil {
		log.Fatalf("Error getting references: %v", err)
	}
	writtendata, err := ioutil.ReadFile("./mining-data/collaborators/" + processName + "/references.json")
	//var references []map[string]interface{}
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
	//TODO: HERE WE ARE NOT CHECKING THE VALIDITY OF THE PROVISIONER SERVICE. WE SHOULD DO IT OR USE SOME CERTIFICATE.
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	//Get the organization public key to be identified by the collaborators
	//MOVE THE PUBLIC KEY FILES HERE---------------------------------------------------------------------------------------------------------------------------------------------------
	organizationPublicKey := readPublicKey()
	//Turn the public key into bytes
	organizationPublicKeyBytes, err := x509.MarshalPKIXPublicKey(&organizationPublicKey)
	if err != nil {
		log.Fatal(err)
	}
	//Iterate over collaborators and get their tracelists
	for _, item := range references {
		httpReference := item.WebReference
		//TODO HERE THE RESPONSE SHOULD BE ENCRYPTED WITH THE MINER'S ORGANIZATION PUBLIC KEY. THE CLIENT SHOULD DECRYPT THE RESPONSE WITH ITS ORGANIZATION'S PRIVATE KEY.
		response := httpPOST(tlsConfig, httpReference+"/tracelistrequest", string(organizationPublicKeyBytes), "https://localhost:"+receiverPort, segmentsize, "")
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
		for trId, _ := range responseJson {
			if _, ok := readTraceMap[trId]; ok {
				readTraceMap[trId][httpReference] = false
				//TODO: IF OK=FALSE (LOOK ABOVE) THEN THE TRACE IS NOT OWNED BY THE MINER SO A NEW ENTRY SHOULD BE ADDED IN THE TRACEMAP WITH THE ID OF THE TRACE. YOU SHOULD DO ALSO ALL THE STUFF IN _X_
			} else {
				/*h is the header, it's not a case reference*/
				if trId != "h" {
					readTraceMap[trId] = make(map[string]bool)
					readTraceMap[trId][httpReference] = false
					os.MkdirAll("./mining-data/consumption-data/"+processName+"/trace_"+trId, os.ModePerm)
				}
			}
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
		httpReference := item.WebReference
		go httpPOST(tlsConfig, httpReference+"/logrequest", string(organizationPublicKeyBytes), "https://localhost:"+receiverPort, segmentsize, string(traceListByte))
	}
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
