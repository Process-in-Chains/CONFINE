package logrequest

import (
	"app/utils/collaborators"
	"app/utils/encryption"
	"app/utils/test"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
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
func LogRequest(processName string, receiverPort string, segmentsize int, TLScert []byte, minerPrivateKey crypto.PrivateKey) {
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
	cert, _ := x509.ParseCertificate(TLScert)
	organizationPublicKeyBytes, _ := x509.MarshalPKIXPublicKey(cert.PublicKey)
	//Turn the public key into bytes
	//organizationPublicKeyBytes, err := x509.MarshalPKIXPublicKey(&organizationPublicKey)
	if err != nil {
		log.Fatal(err)
	}
	//Iterate over collaborators and get their tracelists
	if err != nil {
		log.Fatalf("failed to sign message: %v", err)
	}
	for _, item := range references {
		httpReference := item.WebReference
		//TODO HERE THE RESPONSE SHOULD BE ENCRYPTED WITH THE MINER'S ORGANIZATION PUBLIC KEY. THE CLIENT SHOULD DECRYPT THE RESPONSE WITH ITS ORGANIZATION'S PRIVATE KEY.
		response := httpPOST(tlsConfig, httpReference+"/tracelistrequest", string(organizationPublicKeyBytes), "https://localhost:"+receiverPort, segmentsize, "", string(TLScert))
		if err != nil {
			log.Fatal(err)
		}
		var wrapperResponse map[string]string
		err = json.Unmarshal(response, &wrapperResponse)
		if err != nil {
			fmt.Println(err)
		}
		//Extract the sym key from the http response
		encryptedSymKeyBase64 := wrapperResponse["encryptedKey"]
		encryptedKey, err := base64.StdEncoding.DecodeString(encryptedSymKeyBase64)
		if err != nil {
			log.Fatalf("failed to decode base64 encoded encrypted key: %v", err)
		}
		//Decrypt the symetric key with the private key of the miner TLS certificate
		decryptedKey, err := rsa.DecryptPKCS1v15(rand.Reader, minerPrivateKey.(*rsa.PrivateKey), encryptedKey)
		if err != nil {
			log.Fatalf("failed to decrypt key: %v", err)
		}
		//Extract the encrypted trace list from the http resposne
		encryptedTraceListBase64 := wrapperResponse["traceList"]
		encryptedTraceList, err := base64.StdEncoding.DecodeString(encryptedTraceListBase64)
		if err != nil {
			log.Fatalf("failed to decode base64 encoded encrypted trace list: %v", err)
		}
		//Decrypt the trace list with the decrypted symetric key
		decryptedTraceList, err := encryption.DecryptXES(encryptedTraceList, decryptedKey)
		if err != nil {
			log.Fatalf("failed to decrypt trace list with sym key %v", err)
		}
		var responseJson map[string]string
		//err = json.Unmarshal([]byte(wrapperResponse["traceList"]), &responseJson)
		err = json.Unmarshal(decryptedTraceList, &responseJson)
		if err != nil {
			log.Fatal(err)
		}
		for trId, _ := range responseJson {
			if _, ok := readTraceMap[trId]; ok {
				readTraceMap[trId][httpReference] = false
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
		go httpPOST(tlsConfig, httpReference+"/logrequest", string(organizationPublicKeyBytes), "https://localhost:"+receiverPort, segmentsize, string(traceListByte), string(TLScert))
	}
}
func httpPOST(tlsConfig *tls.Config, posturl string, publicKey string, logreceiver string, segmentSize int, loglist string, TLScert string) []byte {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	// Prepare the form data
	formData := url.Values{"logreceiver": {logreceiver}, "publicKey": {publicKey}, "segmentSize": {strconv.Itoa(segmentSize)}, "loglist": {loglist}, "tlsCert": {TLScert}}
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
