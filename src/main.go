package main

import (
	"app/secure-miner/log-elaboration/miningAlgorithms"
	logmanagement "app/secure-miner/log-management"
	logreception "app/secure-miner/log-reception"
	logrequest "app/secure-miner/log-request"
	"app/utils/reset"
	"app/utils/test"
	"app/utils/xes"
	"crypto/x509"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

//ego-go build -buildvcs=false && ego sign enclave.json && sudo ego run ./app -segsize 2000 -port 8094 -test true
///cd Desktop/Davide/linux-sgx/external/dcap_source/QuoteGeneration/pccs/ && node pccs_server.js
/*
Secure Miner main.
It initializes the necessary variables, parses command-line arguments, starts the log receiver, and prompts the user for commands to execute.
*/

var DeclareProcessModelPath = ""

func main() {
	var port string
	var segmentsize int
	var TESTMODESTR string
	flag.StringVar(&port, "port", "", "Port number")
	flag.StringVar(&TESTMODESTR, "test", "false", "Test mode")
	flag.IntVar(&segmentsize, "segsize", 100, "Segement size to be processed in the TEE. Value in KB.")
	flag.Parse()
	TESTMODE, err := strconv.ParseBool(TESTMODESTR)
	if err != nil {
		fmt.Println("Error:", err)
	}
	if port == "" {
		fmt.Println("Missing port number")
		return
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println("Number port error:", err)
		return
	}
	fmt.Println("Setting up CONFINE...")
	logReceiver := logreception.NewLogReceiver(p)
	go logReceiver.Start()
	//TODO Disclose TLS public key of the log receiver somewhere
	//TODO ... option 1) copia a manella
	//TODO ... option 2) scrittura automatica in un file
	for true {
		time.Sleep(2000 * time.Millisecond)
		fmt.Println()
		fmt.Printf("Command list:\n1: CONFINE WITH INCREMENTAL DISCOVERY - Discover process model with the incremental HeuristicsMiner via CONFINE protocol\n2: CONFINE WITH NON-INCREMENTAL DISCOVERY - Discover process model with non incremental HeuristicsMiner via CONFINE protocol\n3: CONFINE WITH INCREMENTAL CONFORMANCE CHECKING - Incremental Declare Conformance checking via CONFINE protocol \n4: CONFINE WITH NON-INCREMENTAL CONFORMANCE CHECKING - Non-incremental Declare Conformance checking via CONFINE protocol \n5: Classic HeuristicsMiner execution using the local event log at '/mining-data/provision-data/process-01/event_log.xes'\n6: Show TLS public key of the secure miner\n")
		fmt.Println()
		var command string
		fmt.Scanln(&command)
		/*This command initiates the CONFINE protocol with the HeuristicMiner (discovery algorithm) */
		if command == "1" {
			if TESTMODE {
				test.STOPMONITORING = false
				test.TEST_MODE = true
				go test.PrintRamUsage()
				println("TESTMODE - TEST STARTED AT: ", time.Now().UnixMilli())
			}
			logReceiver.SetAlgorithm("HeuristicsMiner")
			logrequest.LogRequest("process-01", port, segmentsize, logReceiver.GetTLSCertificate(), logReceiver.GetTLSPrivateKey())
			test.WaitUntilStop()
			time.Sleep(1000 * time.Millisecond)
			fmt.Println("Process discovery with HeuristicsMiner completed! Process model saved at /mining-data/output/heuristicsMiner_output.pnml")
		}
		if command == "2" {
			if TESTMODE {
				test.STOPMONITORING = false
				test.TEST_MODE = true
				go test.PrintRamUsage()
				println("TESTMODE - TEST STARTED AT: ", time.Now().UnixMilli())
			}
			logReceiver.SetAlgorithm("ClassicHeuristics")
			logrequest.LogRequest("process-01", port, segmentsize, logReceiver.GetTLSCertificate(), logReceiver.GetTLSPrivateKey())
			test.WaitUntilStop()
			time.Sleep(1000 * time.Millisecond)
			fmt.Println("Process discovery with HeuristicsMiner completed! Process model saved at /mining-data/output/heuristicsMiner_output.pnml")
		}
		/*This command initiates the CONFINE protocol with DeclareConformance (conformance checking algorithm)*/
		if command == "3" {
			fmt.Println("Insert the path of the declare process model inside the '/mining-data/input/' folder of the project (e.g., '/declareModel.json')")
			var processmodelpath string
			fmt.Scanln(&processmodelpath)
			_, err := ioutil.ReadFile("./mining-data/input/" + processmodelpath)
			if err != nil {
				fmt.Println("Error while opening the declare process model '/mining-data/input/" + processmodelpath)
				fmt.Println(err.Error())
			} else {
				DeclareProcessModelPath = "./mining-data/input/" + processmodelpath
				if TESTMODE {
					test.STOPMONITORING = false
					test.TEST_MODE = true
					go test.PrintRamUsage()
					println("TESTMODE - TEST STARTED AT: ", time.Now().UnixMilli())
				}
				logReceiver.SetAlgorithm("IncrementalDeclareConformance")
				logReceiver.SetProcessModel("./mining-data/input/" + processmodelpath)
				logrequest.LogRequest("process-01", port, segmentsize, logReceiver.GetTLSCertificate(), logReceiver.GetTLSPrivateKey())
				test.WaitUntilStop()
				time.Sleep(1000 * time.Millisecond)
				fmt.Println("Conformance checking with DeclareConformance completed! Conformance results saved at /mining-data/output/declareConformance.json")
			}
		}
		if command == "4" {
			fmt.Println("Insert the path of the declare process model inside the '/mining-data/input/' folder of the project (e.g., '/declareModel.json')")
			var processmodelpath string
			fmt.Scanln(&processmodelpath)
			_, err := ioutil.ReadFile("./mining-data/input/" + processmodelpath)
			if err != nil {
				fmt.Println("Error while opening the declare process model '/mining-data/input/" + processmodelpath)
				fmt.Println(err.Error())
			} else {
				DeclareProcessModelPath = "./mining-data/input/" + processmodelpath
				if TESTMODE {
					test.STOPMONITORING = false
					test.TEST_MODE = true
					go test.PrintRamUsage()
					println("TESTMODE - TEST STARTED AT: ", time.Now().UnixMilli())
				}
				logReceiver.SetAlgorithm("DeclareConformance")
				logReceiver.SetProcessModel("./mining-data/input/" + processmodelpath)
				logrequest.LogRequest("process-01", port, segmentsize, logReceiver.GetTLSCertificate(), logReceiver.GetTLSPrivateKey())
				test.WaitUntilStop()
				time.Sleep(1000 * time.Millisecond)
				fmt.Println("Conformance checking with DeclareConformance completed! Conformance results saved at /mining-data/output/declareConformance.json")
			}
		}
		/*This command initiates the mining process using 'event_log.xes' located in './mining-data/provision-data/process-01'.*/
		if command == "5" {
			log_path := "./mining-data/provision-data/" + "process-01" + "/event_log.xes"
			_, err := ioutil.ReadFile(log_path)
			if err != nil {
				fmt.Println("Error while opening the event log at", log_path)
				fmt.Println(err.Error())
			} else {
				eventLog := xes.ReadXes(log_path)
				prosessMiningAlgorithms.HeuristicMiner(eventLog.XesToSlices(), "process-01")
			}
		}
		/*This comand reads the public key associated with the Secure Miner*/
		if command == "6" {
			tlsCert := logReceiver.GetTLSCertificate()
			cert, _ := x509.ParseCertificate(tlsCert)
			pubBytes, _ := x509.MarshalPKIXPublicKey(cert.PublicKey)
			fmt.Println("TLS public key of the secure miner:\n" + base64.StdEncoding.EncodeToString(pubBytes))
		}

		logreception.FIRSTTS = false
		logreception.FIRSTATT = false
		logmanagement.FIRSTCOMP = false
		err := reset.DeleteAllFilesInSubfolders("mining-data/consumption-data/process-01/")
		if err != nil {
			fmt.Println("Error deleting files:", err)
		}
		reset.ReplaceWithEmptyMatrix("./mining-data/consumption-data/process-01/miningMetadata/dependencyMatrix.json")
		reset.ReplaceWithEmptyMatrix("./mining-data/consumption-data/process-01/miningMetadata/dependencyMatrix2len.json")
		reset.ReplaceWithEmptyMap("./mining-data/consumption-data/process-01/miningMetadata/map.json")
		reset.DeleteEmptySubfolders("mining-data/consumption-data/process-01/")
	}
}
