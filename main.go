// https://tour.golang.org/concurrency/1

package main

import (
	"app/log_consumption"
	prosessMiningAlgorithms "app/log_consumption/miningAlgorithms"
	"app/log_request"
	"app/utils/reset"
	"app/utils/test"
	"app/utils/xes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

func main() {
	var port string
	var segmentsize int
	var TESTMODE bool
	flag.StringVar(&port, "port", "", "Port number")
	flag.IntVar(&segmentsize, "segsize", 100, "Segement size to be processed in the TEE. Value in KB.")
	flag.BoolVar(&TESTMODE, "test", false, "Test mode")
	flag.Parse()
	if port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			fmt.Println("Number port error:", err)
			return
		}
		server := logreception.NewLogReceiver(p)
		go server.Start()
	} else {
		fmt.Println("Missing port number")
		return
	}
	for true {
		fmt.Printf("Command list:----------------------------------------------------------\n1: Start mining (Inter-organizational)\n2: Start mining (Only owned traces)\n3: Show public key\n4: Reset knowledge\n")
		var command string
		fmt.Scanln(&command)
		if command == "1" {
			if TESTMODE {
				test.STOPMONITORING = false
				go test.PrintRamUsage()
			}
			println("TESTMODE - TEST STARTED AT: ", time.Now().UnixMilli())
			logrequest.LogRequest("process-01", port, segmentsize)
			test.WaitUntilStop()

		}
		if command == "2" {
			log_path := "./mining-data/provision-data/" + "process-01" + "/event_log.xes"
			eventLog := xes.ReadXes(log_path)
			prosessMiningAlgorithms.HeuristicMiner(eventLog.XesToSlices(), "process-01")
		}
		if command == "3" {

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
			a := *rsaPubKey
			pubKeyBytes, err := x509.MarshalPKIXPublicKey(&a)
			if err != nil {
				panic(err)
			}
			fmt.Println("Public key:\n" + base64.StdEncoding.EncodeToString(pubKeyBytes))
		}
		if command == "4" {
			err := reset.DeleteAllFilesInSubfolders("mining-data/consumption-data/process-01/")
			if err != nil {
				fmt.Println("Error deleting files:", err)
			}
			reset.ReplaceWithEmptyMatrix("./mining-data/consumption-data/process-01/miningMetadata/dependencyMatrix.json")
			reset.ReplaceWithEmptyMatrix("./mining-data/consumption-data/process-01/miningMetadata/dependencyMatrix2len.json")
			reset.ReplaceWithEmptyMap("./mining-data/consumption-data/process-01/miningMetadata/map.json")
			reset.DeleteEmptySubfolders("mining-data/consumption-data/process-01/")
			err = reset.DeleteTraceFolders()
			if err != nil {
				fmt.Println("Error deleting trace folders:", err)
			}
		}
		logreception.FIRSTTS = false
		logreception.FIRSTATT = false
		logreception.FIRSTCOMP = false
	}
}
