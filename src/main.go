package main

import (
	"app/secure-miner/log-elaboration/miningAlgorithms"
	logreception "app/secure-miner/log-reception"
	logrequest "app/secure-miner/log-request"
	"app/utils/encryption"
	"app/utils/reset"
	"app/utils/test"
	"app/utils/xes"
	"flag"
	"fmt"
	"strconv"
	"time"
)

//ego-go build -buildvcs=false && ego sign enclave.json && OE_SIMULATION=1 ego run ./app -segsize 2000 -port 8080 -test true
/*
Secure Miner main.
It initializes the necessary variables, parses command-line arguments, starts the log receiver, and prompts the user for commands to execute.
*/
func main() {
	var port string
	var segmentsize int
	var TESTMODE bool
	flag.StringVar(&port, "port", "", "Port number")
	flag.IntVar(&segmentsize, "segsize", 100, "Segement size to be processed in the TEE. Value in KB.")
	flag.BoolVar(&TESTMODE, "test", false, "Test mode")
	flag.Parse()
	if port == "" {
		fmt.Println("Missing port number")
		return
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println("Number port error:", err)
		return
	}
	logReceiver := logreception.NewLogReceiver(p)
	go logReceiver.Start()
	for true {
		fmt.Printf("Command list:------------------------------------------------\n1: CONFINE DISCOVERY-Discover process model with HeuristicsMiner via CONFINE protocol\n2:CONFINE-CONFORMANCE CHEKCING: Conformance checking of the JSON declare model declare model at '/mining-data/input/declareModel.json' via CONFINE protocol \n3: Apply HeuristicsMiner using the local event log at '/mining-data/provision-data/process-01/event_log.xes'\n4: Show miner public key\n5: Reset Secure Miner's memory\n")
		var command string
		fmt.Scanln(&command)
		/*This command initiates the mining process via CONFINE protocol.*/
		if command == "1" {
			if TESTMODE {
				test.STOPMONITORING = false
				go test.PrintRamUsage()
				println("TESTMODE - TEST STARTED AT: ", time.Now().UnixMilli())
			}
			logReceiver.SetAlgorithm("HeuristicsMiner")
			logrequest.LogRequest("process-01", port, segmentsize)
			test.WaitUntilStop()

		}
		if command == "2" {
			if TESTMODE {
				test.STOPMONITORING = false
				go test.PrintRamUsage()
				println("TESTMODE - TEST STARTED AT: ", time.Now().UnixMilli())
			}
			logReceiver.SetAlgorithm("DeclareConformance")
			logrequest.LogRequest("process-01", port, segmentsize)
			test.WaitUntilStop()

		}
		/*This command initiates the mining process using 'event_log.xes' located in './mining-data/provision-data/process-01'.*/
		if command == "3" {
			log_path := "./mining-data/provision-data/" + "process-01" + "/event_log.xes"
			eventLog := xes.ReadXes(log_path)
			prosessMiningAlgorithms.HeuristicMiner(eventLog.XesToSlices(), "process-01")
		}
		/*This comand reads the public key associated with the Secure Miner*/
		if command == "4" {
			pubKey, err := encryption.ParsePublicKeyToString("./public.pem")
			if err != nil {
				fmt.Println("Error reading public key:", err)
			} else {
				fmt.Println("Public key:\n" + pubKey)
			}

		}
		/*This comand resets the result knwoledge gained from the previous CONFINE elaborations. It resets the knowledge by deleting files and emptying matrices and maps.*/
		if command == "5" {
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
