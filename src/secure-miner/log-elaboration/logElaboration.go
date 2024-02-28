package logelaboration

import (
	"app/utils/xes"
	"fmt"
)
import "app/secure-miner/log-elaboration/miningAlgorithms"

type LogElaborator struct {
	algorithms []string
}

func NewLogElaborator() *LogElaborator {
	logElaborator := &LogElaborator{algorithms: []string{"HeuristicsMiner", "DeclareConformance"}}
	return logElaborator
}

func (logElaborator *LogElaborator) ApplyAlgorithm(algorithm string, processName string, eventLog xes.XES) {
	//Check if the algorithm is unsupported
	if logElaborator.isAlgorithmSupported(algorithm) != nil {
		fmt.Println("Algorithm not supported")
		return
	}
	if algorithm == "HeuristicsMiner" {
		prosessMiningAlgorithms.HeuristicMiner(eventLog.XesToSlices(), processName)
	}
	if algorithm == "DeclareConformance" {
		prosessMiningAlgorithms.DeclareConformance(eventLog, "mining-data/input/declareModel.json")

	}
}

// Internal function of the log elaborator that verifies if a string algorithm is supported. If not, it returns an error.
func (logElaborator *LogElaborator) isAlgorithmSupported(algorithm string) error {
	//Check if the algorithm is unsupported
	for _, alg := range logElaborator.algorithms {
		if alg == algorithm {
			return nil
		}
	}
	return fmt.Errorf("Algorithm not supported")
}
