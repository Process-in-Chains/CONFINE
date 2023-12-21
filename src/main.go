package main

import (
	logreception "app/secure-miner/log-consumption"
	"app/secure-miner/log-consumption/miningAlgorithms"
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

/*
The main of the Secure Miner.
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
		fmt.Printf("Command list:------------------------------------------------\n1: Start CONFINE protocol\n2: Mine specific log\n3: Show miner public key\n4: Reset Secure Miner's memory\n")
		var command string
		fmt.Scanln(&command)
		/*This command initiates the mining process via CONFINE protocol.*/
		if command == "1" {
			if TESTMODE {
				test.STOPMONITORING = false
				go test.PrintRamUsage()
				println("TESTMODE - TEST STARTED AT: ", time.Now().UnixMilli())
			}
			logrequest.LogRequest("process-01", port, segmentsize)
			test.WaitUntilStop()

		}
		/*This command initiates the mining process using 'event_log.xes' located in './mining-data/provision-data/process-01'.*/
		if command == "2" {
			log_path := "./mining-data/provision-data/" + "process-01" + "/event_log.xes"
			eventLog := xes.ReadXes(log_path)
			prosessMiningAlgorithms.HeuristicMiner(eventLog.XesToSlices(), "process-01")
		}
		/*This comand reads the public key associated with the Secure Miner*/
		if command == "3" {
			pubKey, err := encryption.ParsePublicKeyToString("./public.pem")
			if err != nil {
				fmt.Println("Error reading public key:", err)
			} else {
				fmt.Println("Public key:\n" + pubKey)
			}

		}
		/*This comand resets the result knwoledge gained from the previous CONFINE elaborations. It resets the knowledge by deleting files and emptying matrices and maps.*/
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
