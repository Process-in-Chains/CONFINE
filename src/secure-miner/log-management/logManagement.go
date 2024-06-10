package logmanagement

import (
	logelaboration "app/secure-miner/log-elaboration"
	"app/utils/test"
	"app/utils/xes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

/*Variable for testing */
var FIRSTCOMP = false

/*This function handles a log segment sent by a provisioner to the Secure Miner*/
func HandleSegment(logSegment xes.XES, processName string, publicKey string, myreference string, mergekey string, traceMap map[string]map[string]bool, algorithm string, logElaborator logelaboration.LogElaborator) {
	/*For each trace in the log segment...*/
	for _, trace := range logSegment.Traces {
		//traceId, _ := trace.GetId()
		traceId, error := trace.GetAttributeValue(mergekey)
		if error != nil {
			log.Fatal(error)
			return
		}
		/*If the traceId exist and the provisioner has not delivered yet...*/
		if _, ok := traceMap[traceId]; ok && !traceMap[traceId][myreference] {
			/*If the algorithm it's the HeuristicsMiner, we apply an incremental approach on the traces. Each completed trace is computed individually and increment the global result.*/
			//TODO: THIS CHECK SHOULD BE REMOVED IN FUTURE VERSIONS FOR BETTER GENERALIZATION.
			if algorithm == "HeuristicsMiner" {
				/*Mark the trace of the provisioner as arrived*/
				traceMap[traceId][myreference] = true
				/*If each provisioner has delivered the trace*/
				if isTraceCompleted(traceMap[traceId]) {
					/*Get the references to the partial traces in memory*/
					files, err := os.ReadDir("/mining-data/consumption-data/" + processName + "/trace_" + traceId)
					if err != nil {
						log.Fatal(err)
					}
					mergedTrace := xes.Trace{}
					/*Merge each partial trace in memory*/
					for nfile, traceFile := range files {
						if nfile == 0 {
							mergedTrace, _ = xes.MergeTraces(trace, xes.ReadXes("mining-data/consumption-data/" + processName + "/trace_" + traceId + "/" + traceFile.Name()).Traces[0])
						} else {
							currentTrace := xes.ReadXes("mining-data/consumption-data/" + processName + "/trace_" + traceId + "/" + traceFile.Name())
							mergedTrace, _ = xes.MergeTraces(mergedTrace, currentTrace.Traces[0])
						}
					}
					/*Format the merged trace in an empty log*/
					traceArray := []xes.Trace{}
					traceArray = append(traceArray, mergedTrace)
					logWithMergedTrace := xes.XES{Traces: traceArray}
					if !FIRSTCOMP {
						fmt.Println("TESTMODE - FIRST COMPUTATION AT:", time.Now().UnixMilli())
						FIRSTCOMP = true
					}
					/*Use the log elaborator and apply the HeuristicsMiner algorithm on the merged trace*/
					logElaborator.ApplyAlgorithm(algorithm, processName, logWithMergedTrace)
				} else /*If some provisioner has not already sent its trace...*/ {
					/*Store the trace in memory*/
					storeTrace(processName, traceId, trace, url.PathEscape(fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))))
				}
			}
			/*If the algorithm is DeclareConformance, we adopt a lazy computation approach. Therefore, we compute just one time, when the whole log is completed.*/
			//TODO THIS CHECK SHOULD BE REMOVED FOR BETTER GENERALIZATION (AS PER HEURISTICSMINER)
			if algorithm == "DeclareConformance" {
				/*Set the trace from*/
				traceMap[traceId][myreference] = true
				mergedTrace := xes.Trace{}
				//mergedTrace, err := xes.MergeTraces(trace, xes.ReadXes("mining-data/consumption-data/" + processName + "/trace_" + traceId + "/trace_"+traceId+"_merged.xes"))
				if _, err := os.Stat("mining-data/consumption-data/" + processName + "/trace_" + traceId + "/trace_" + traceId + "_merged.xes"); os.IsNotExist(err) {
					mergedTrace = trace
				} else {
					mergedTrace, _ = xes.MergeTraces(trace, xes.ReadXes("mining-data/consumption-data/" + processName + "/trace_" + traceId + "/trace_" + traceId + "_merged.xes").Traces[0])
				}
				storeTrace(processName, traceId, mergedTrace, traceId+"_merged")
			}
		} else {
			fmt.Printf("'%s' is not a key in the map or sender not expected\n", traceId)
		}
		/*Store the updated tracemap*/
		jsonData, err := json.MarshalIndent(traceMap, "", "  ")
		if err != nil {
			fmt.Println("Error converting JSON:", err)
			return
		}
		err = ioutil.WriteFile("./mining-data/consumption-data/"+processName+"/traceMap.json", jsonData, 0644)
		if err != nil {
			fmt.Println("Errore nella scrittura del file JSON:", err)
			return
		}
		/*If all the provisioners have delivered their traces....*/
		if allTracesCollected(traceMap) {
			/*In this case, we can apply the DeclareConformance algorithm*/
			//TODO this check should be removed for better generalization
			if algorithm == "DeclareConformance" {
				/*Collect all the merged traces*/
				finalLog := collectMergedTraces(processName)
				if !FIRSTCOMP {
					fmt.Println("TESTMODE - FIRST COMPUTATION AT:", time.Now().UnixMilli())
					FIRSTCOMP = true
				}
				/*Apply the DeclareConformance algorithm on the merged traces*/
				logElaborator.ApplyAlgorithm(algorithm, processName, finalLog)
			}
			/*Print for testing*/
			fmt.Println("TESTMODE - TEST ENDED AT: ", time.Now().UnixMilli())
			test.STOPMONITORING = true
		}
	}
}

/*Check if all the provisioners have delivered their piece of trace*/
func isTraceCompleted(myMap map[string]bool) bool {
	for _, value := range myMap {
		if !value {
			return false
		}
	}
	return true
}

/*Check if all the traces*/
func allTracesCollected(traceMap map[string]map[string]bool) bool {
	for _, value := range traceMap {
		if !isTraceCompleted(value) {
			return false
		}
	}
	return true
}

/*Read all the merged trace in memory and put it in an event log*/
func collectMergedTraces(processName string) xes.XES {
	// Get the list of traces folders in consumption-data
	tracesDir := "mining-data/consumption-data/" + processName
	traceFolders, err := ioutil.ReadDir(tracesDir)
	if err != nil {
		log.Fatal("Error reading trace folders:", err)
	}
	// Iterate over each trace folder
	finalLog := xes.XES{}
	for _, traceFolder := range traceFolders {
		if traceFolder.IsDir() && traceFolder.Name() != "miningMetadata" {
			traceFolderPath := filepath.Join(tracesDir, traceFolder.Name())
			// Get the list of trace files in the current trace folder
			traceFiles, err := ioutil.ReadDir(traceFolderPath)
			if err != nil {
				log.Fatal("Error reading trace files:", err)
			}
			// Iterate over each trace file
			for _, traceFile := range traceFiles {
				if !traceFile.IsDir() {
					traceFilePath := filepath.Join(traceFolderPath, traceFile.Name())
					// Read the trace file
					currentTrace := xes.ReadXes(traceFilePath).Traces[0]
					finalLog.Traces = append(finalLog.Traces, currentTrace)
				}
			}
		}
	}
	return finalLog

}

/*Store the trace in memory*/
func storeTrace(processName string, traceId string, trace xes.Trace, fileName string) {
	filename := "mining-data/consumption-data/" + processName + "/trace_" + traceId + "/trace_" + fileName + ".xes"
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

//func storeMergedTrace(processName string, traceId string, trace xes.Trace) {
//	filename := "mining-data/consumption-data/" + processName + "/trace_" + traceId + "/trace_" + traceId + "_merged.xes"
//	byteTrace := trace.TraceToByte()
//	err := os.MkdirAll(filepath.Dir(filename), os.ModePerm)
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = os.WriteFile(filename, byteTrace, 0644)
//	if err != nil {
//		log.Fatal(err)
//	}
//	//}
//}


