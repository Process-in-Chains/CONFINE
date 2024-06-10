package prosessMiningAlgorithms

import (
	"app/utils/xes"
	"encoding/json"
	"io/ioutil"
	"strings"
)

// go build -o decconf declare-conformance/declareConformance.go && ./decconf
type Deviation struct {
	Type string
	Act0 string
	Act1 string
}
type Tuple struct {
	Value0 string
	Value1 string
}
/*String labels of the Declare contraints*/
var RESPONDED_EXISTENCE = "responded_existence"
var EXISTENCE = "existence"
var ABSENCE = "absence"
var EXACTLY_ONE = "exactly_one"
var INIT = "init"
var SUCCESSION = "succession"
var COEXISTENCE = "coexistence"
var RESPONSE = "response"
var PRECEDENCE = "precedence"
var NONCOEXISTENCE = "noncoexistence"
var ALTRESPONSE = "altresponse"
var CHAINRESPONSE = "chainresponse"
var ALTPRECEDENCE = "altprecedence"
var CHAINPRECEDENCE = "chainprecedence"
var ALTSUCCESSION = "altsuccession"
var CHAINSUCCESSION = "chainsuccession"

/*Function of the Conformance Checking algorithm with Declare models*/
func DeclareConformance(eventLog xes.XES, modelPath string) {
	processModelByets, err := ioutil.ReadFile(modelPath)
	jsonModelStructure := make(map[string]map[string]map[string]int)
	err = json.Unmarshal(processModelByets, &jsonModelStructure)
	if err != nil {
		print(err.Error())
	}
	result := applyAlgorithm(eventLog.XesToSlices(), jsonModelStructure, nil)
	//Write result in a json file
	ginoBytes, _ := json.Marshal(result)
	err = ioutil.WriteFile("mining-data/output/declareConformance.json", ginoBytes, 0644)

}
/*Print the results of the conformance checking*/
func printConformanceResult(result []map[string]interface{}) {
	for _, v := range result {
		println("------------------------------------------------------------")
		for conformanceOutput, conformanceOutputValue := range v {
			if conformanceOutput == "no_constr_total" {
				println(conformanceOutput, conformanceOutputValue.(int))
			}
			if conformanceOutput == "no_dev_total" {
				println(conformanceOutput, conformanceOutputValue.(int))
			}
			if conformanceOutput == "dev_fitness" {
				println(conformanceOutput, conformanceOutputValue.(float64))
			}
			if conformanceOutput == "is_fit" {
				println(conformanceOutput, conformanceOutputValue.(bool))
			}
			if conformanceOutput == "deviations" {
				devList := conformanceOutputValue.([]Deviation)
				for _, dev := range devList {
					println(dev.Type, dev.Act0, dev.Act1)
				}
			}
		}
	}
}

/*Execute all the checks of the algorithm*/
func applyAlgorithm(projectedLog [][]string, model map[string]map[string]map[string]int, parameters map[interface{}]interface{}) []map[string]interface{} {
	if parameters == nil {
		parameters = make(map[interface{}]interface{})
	}
	confCases := make([]map[string]interface{}, 0)
	totalNumConstraints := 0
	for _, v := range model {
		totalNumConstraints += len(v)
	}
	for _, trace := range projectedLog {
		actIdxs := make(map[string][]int)
		for i, act := range trace {
			actIdxs[act] = append(actIdxs[act], i)
		}
		ret := make(map[string]interface{})
		ret["no_constr_total"] = totalNumConstraints
		ret["deviations"] = make([]Deviation, 0)
		checkExistence(trace, model, ret, parameters)
		checkRespondedExistence(trace, model, ret)
		checkAbsence(trace, model, ret, parameters)
		checkExactlyOne(trace, model, ret, parameters)
		checkInit(trace, model, ret, parameters)
		checkCoexistence(trace, model, ret, parameters)
		checkNonCoexistence(trace, model, ret, parameters)
		checkResponse(trace, model, ret, actIdxs, parameters)
		checkPrecedence(trace, model, ret, actIdxs, parameters)
		checkSuccession(trace, model, ret, actIdxs, parameters)
		checkAltResponse(trace, model, ret, actIdxs, parameters)
		checkChainResponse(trace, model, ret, actIdxs, parameters)
		checkAltPrecedence(trace, model, ret, actIdxs, parameters)
		checkChainPrecedence(trace, model, ret, actIdxs, parameters)
		checkAltSuccession(trace, model, ret, actIdxs, parameters)
		checkChainSuccession(trace, model, ret, actIdxs, parameters)
		ret["no_dev_total"] = len(ret["deviations"].([]Deviation))
		ret["dev_fitness"] = 1 - float64(ret["no_dev_total"].(int))/float64(ret["no_constr_total"].(int))
		ret["is_fit"] = ret["no_dev_total"].(int) == 0
		confCases = append(confCases, ret)
	}
	return confCases
}
/*Existence constraint check*/
func checkExistence(trace []string, model map[string]map[string]map[string]int, traceDict map[string]interface{}, parameters map[interface{}]interface{}) {
	if existence, ok := model[EXISTENCE]; ok {
		for act := range existence {
			if !contains(trace, act) {
				deviations := traceDict["deviations"].([]Deviation)
				deviations = append(deviations, Deviation{Type: "EXISTENCE", Act0: act})
				traceDict["deviations"] = deviations
			}
		}
	}
}
/*Not existence constraints check*/
func checkAbsence(trace []string, model map[string]map[string]map[string]int, traceDict map[string]interface{}, parameters map[interface{}]interface{}) {
	if _, ok := model[ABSENCE]; ok {
		for act := range model[ABSENCE] {
			if contains(trace, act) {
				deviations := traceDict["deviations"].([]Deviation)
				deviations = append(deviations, Deviation{Type: "NONEXISTENCE", Act0: act})
				traceDict["deviations"] = deviations
			}
		}
	}
}
/*Check exactly one constraint check*/
func checkExactlyOne(trace []string, model map[string]map[string]map[string]int, traceDict map[string]interface{}, parameters map[interface{}]interface{}) {
	if val, ok := model[EXACTLY_ONE]; ok {
		traceCounter := make(map[string]int)
		for _, act := range trace {
			traceCounter[act]++
		}
		for act := range val {
			if traceCounter[act] != 1 {
				deviations := traceDict["deviations"].([]Deviation)
				deviations = append(deviations, Deviation{Type: "EXACTLY_ONE", Act0: act})
				traceDict["deviations"] = deviations
			}
		}
	}
}
/*Responded existence constrait check*/
func checkRespondedExistence(trace []string, model map[string]map[string]map[string]int, traceDict map[string]interface{}) {
	if val, ok := model[RESPONDED_EXISTENCE]; ok {
		for stringActCouple := range val {
			actCouple := parseString(stringActCouple)
			if contains(trace, actCouple.Value0) && !contains(trace, actCouple.Value1) {
				deviations := traceDict["deviations"].([]Deviation)
				deviations = append(deviations, Deviation{Type: "RESPONDED_EXISTENCE", Act0: actCouple.Value0, Act1: actCouple.Value1})
				traceDict["deviations"] = deviations
			}
		}
	}
}
/*Init constraint check*/
func checkInit(trace []string, model map[string]map[string]map[string]int, traceDict map[string]interface{}, parameters map[interface{}]interface{}) {
	if val, ok := model[INIT]; ok {
		for act := range val {
			if len(trace) == 0 || trace[0] != act {
				deviations := traceDict["deviations"].([]Deviation)
				deviations = append(deviations, Deviation{Type: "INIT", Act0: act})
				traceDict["deviations"] = deviations
			}
		}
	}
}
func checkCoexistence(trace []string, model map[string]map[string]map[string]int, traceDict map[string]interface{}, parameters map[interface{}]interface{}) {
	if val, ok := model[COEXISTENCE]; ok {
		for stringActCouple := range val {
			actCouple := parseString(stringActCouple)
			if (contains(trace, actCouple.Value0) && !contains(trace, actCouple.Value1)) || (contains(trace, actCouple.Value1) && !contains(trace, actCouple.Value0)) {
				deviations := traceDict["deviations"].([]Deviation)
				deviations = append(deviations, Deviation{Type: "COEXISTENCE", Act0: actCouple.Value0, Act1: actCouple.Value1})
				traceDict["deviations"] = deviations
			}
		}
	}
}

func checkNonCoexistence(trace []string, model map[string]map[string]map[string]int, traceDict map[string]interface{}, parameters map[interface{}]interface{}) {
	if val, ok := model[NONCOEXISTENCE]; ok {
		for stringActCouple := range val {
			actCouple := parseString(stringActCouple)
			if contains(trace, actCouple.Value0) && contains(trace, actCouple.Value1) {
				deviations := traceDict["deviations"].([]Deviation)
				deviations = append(deviations, Deviation{Type: "NONCOEXISTENCE", Act0: actCouple.Value0, Act1: actCouple.Value1})
				traceDict["deviations"] = deviations
			}
		}
	}
}

func checkResponse(trace []string, model map[string]map[string]map[string]int, traceDict map[string]interface{}, actIdxs map[string][]int, parameters map[interface{}]interface{}) {
	if val, ok := model[RESPONSE]; ok {
		for stringActCouple := range val {
			actCouple := parseString(stringActCouple)
			if contains(trace, actCouple.Value0) {
				if !contains(trace, actCouple.Value1) || max(actIdxs[actCouple.Value0]) > max(actIdxs[actCouple.Value1]) {
					deviations := traceDict["deviations"].([]Deviation)
					deviations = append(deviations, Deviation{Type: "RESPONSE", Act0: actCouple.Value0, Act1: actCouple.Value1})
					traceDict["deviations"] = deviations
				}
			}
		}
	}
}

func checkPrecedence(trace []string, model map[string]map[string]map[string]int, traceDict map[string]interface{}, actIdxs map[string][]int, parameters map[interface{}]interface{}) {
	if val, ok := model[PRECEDENCE]; ok {
		for stringActCouple := range val {
			actCouple := parseString(stringActCouple)
			if contains(trace, actCouple.Value1) {
				if !contains(trace, actCouple.Value0) || min(actIdxs[actCouple.Value0]) > min(actIdxs[actCouple.Value1]) {
					deviations := traceDict["deviations"].([]Deviation)
					deviations = append(deviations, Deviation{Type: "PRECEDENCE", Act0: actCouple.Value0, Act1: actCouple.Value1})
					traceDict["deviations"] = deviations
				}
			}
		}
	}
}

func checkSuccession(trace []string, model map[string]map[string]map[string]int, traceDict map[string]interface{}, actIdxs map[string][]int, parameters map[interface{}]interface{}) {
	if val, ok := model[SUCCESSION]; ok {
		for stringActCouple := range val {
			actCouple := parseString(stringActCouple)
			if (!contains(trace, actCouple.Value0) || !contains(trace, actCouple.Value1)) || min(actIdxs[actCouple.Value0]) > min(actIdxs[actCouple.Value1]) || max(actIdxs[actCouple.Value0]) > max(actIdxs[actCouple.Value1]) {
				deviations := traceDict["deviations"].([]Deviation)
				deviations = append(deviations, Deviation{Type: "SUCCESSION", Act0: actCouple.Value0, Act1: actCouple.Value1})
				traceDict["deviations"] = deviations
			}
		}
	}
}

func checkAltResponse(trace []string, model map[string]map[string]map[string]int, traceDict map[string]interface{}, actIdxs map[string][]int, parameters map[interface{}]interface{}) {
	if val, ok := model[ALTRESPONSE]; ok {
		for stringActCouple := range val {
			actCouple := parseString(stringActCouple)
			specIdxs := make([][2]interface{}, 0)
			if contains(trace, actCouple.Value0) {
				for _, i := range actIdxs[actCouple.Value0] {
					specIdxs = append(specIdxs, [2]interface{}{actCouple.Value0, i})
				}
			}
			if contains(trace, actCouple.Value1) {
				for _, i := range actIdxs[actCouple.Value1] {
					specIdxs = append(specIdxs, [2]interface{}{actCouple.Value1, i})
				}
			}
			specIdxs = sortSpecIdxs(specIdxs)
			for len(specIdxs) > 0 {
				if specIdxs[0][0] != actCouple.Value0 {
					specIdxs = specIdxs[1:]
				} else {
					break
				}
			}
			isOk := true
			for i := 0; i < len(specIdxs); i++ {
				if i%2 == 0 && (specIdxs[i][0] != actCouple.Value0 || i == len(specIdxs)-1 || specIdxs[i+1][0] != actCouple.Value1) {
					isOk = false
					break
				}
			}
			if !isOk {
				deviations := traceDict["deviations"].([]Deviation)
				deviations = append(deviations, Deviation{"ALTRESPONSE", actCouple.Value0, actCouple.Value1})
				traceDict["deviations"] = deviations
			}
		}
	}
}

func checkChainResponse(trace []string, model map[string]map[string]map[string]int, traceDict map[string]interface{}, actIdxs map[string][]int, parameters map[interface{}]interface{}) {
	if val, ok := model[CHAINRESPONSE]; ok {
		for stringActCouple := range val {
			actCouple := parseString(stringActCouple)
			specIdxs := make([][2]interface{}, 0)
			if contains(trace, actCouple.Value0) {
				for _, i := range actIdxs[actCouple.Value0] {
					specIdxs = append(specIdxs, [2]interface{}{actCouple.Value0, i})
				}
			}
			if contains(trace, actCouple.Value1) {
				for _, i := range actIdxs[actCouple.Value1] {
					specIdxs = append(specIdxs, [2]interface{}{actCouple.Value1, i})
				}
			}
			specIdxs = sortSpecIdxs(specIdxs)
			for len(specIdxs) > 0 {
				if specIdxs[0][0] != actCouple.Value0 {
					specIdxs = specIdxs[1:]
				} else {
					break
				}
			}
			isOk := true
			for i := 0; i < len(specIdxs); i++ {
				if i%2 == 0 && (specIdxs[i][0] != actCouple.Value0 || i == len(specIdxs)-1 || specIdxs[i+1][0] != actCouple.Value1 || specIdxs[i+1][1] != specIdxs[i][1].(int)+1) {
					isOk = false
					break
				}
			}
			if !isOk {
				deviations := traceDict["deviations"].([]Deviation)
				deviations = append(deviations, Deviation{"CHAINRESPONSE", actCouple.Value0, actCouple.Value1})
				traceDict["deviations"] = deviations
			}
		}
	}
}

func checkAltPrecedence(trace []string, model map[string]map[string]map[string]int, traceDict map[string]interface{}, actIdxs map[string][]int, parameters map[interface{}]interface{}) {
	if val, ok := model[ALTPRECEDENCE]; ok {
		for stringActCouple := range val {
			actCouple := parseString(stringActCouple)
			specIdxs := make([][2]interface{}, 0)
			if contains(trace, actCouple.Value0) {
				for _, i := range actIdxs[actCouple.Value0] {
					specIdxs = append(specIdxs, [2]interface{}{actCouple.Value0, i})
				}
			}
			if contains(trace, actCouple.Value1) {
				for _, i := range actIdxs[actCouple.Value1] {
					specIdxs = append(specIdxs, [2]interface{}{actCouple.Value1, i})
				}
			}
			specIdxs = sortSpecIdxs(specIdxs)
			for len(specIdxs) > 1 {
				if specIdxs[1][0] != actCouple.Value1 {
					specIdxs = specIdxs[1:]
				} else {
					break
				}
			}
			isOk := true
			for i := 0; i < len(specIdxs); i++ {
				if i%2 == 0 && (specIdxs[i][0] != actCouple.Value0 || i == len(specIdxs)-1 || specIdxs[i+1][0] != actCouple.Value1) {
					isOk = false
					break
				}
			}
			if !isOk {
				deviations := traceDict["deviations"].([]Deviation)
				deviations = append(deviations, Deviation{"ALTPRECEDENCE", actCouple.Value0, actCouple.Value1})
				traceDict["deviations"] = deviations
			}
		}
	}
}

func checkChainPrecedence(trace []string, model map[string]map[string]map[string]int, traceDict map[string]interface{}, actIdxs map[string][]int, parameters map[interface{}]interface{}) {
	if val, ok := model[CHAINPRECEDENCE]; ok {
		for stringActCouple := range val {
			actCouple := parseString(stringActCouple)
			specIdxs := make([][2]interface{}, 0)
			if contains(trace, actCouple.Value0) {
				for _, i := range actIdxs[actCouple.Value0] {
					specIdxs = append(specIdxs, [2]interface{}{actCouple.Value0, i})
				}
			}
			if contains(trace, actCouple.Value1) {
				for _, i := range actIdxs[actCouple.Value1] {
					specIdxs = append(specIdxs, [2]interface{}{actCouple.Value1, i})
				}
			}
			specIdxs = sortSpecIdxs(specIdxs)
			for len(specIdxs) > 1 {
				if specIdxs[1][0] != actCouple.Value1 {
					specIdxs = specIdxs[1:]
				} else {
					break
				}
			}
			isOk := true
			for i := 0; i < len(specIdxs); i++ {
				if i%2 == 0 && (specIdxs[i][0] != actCouple.Value0 || i == len(specIdxs)-1 || specIdxs[i+1][0] != actCouple.Value1 || specIdxs[i+1][1] != specIdxs[i][1].(int)+1) {
					isOk = false
					break
				}
			}
			if !isOk {
				deviations := traceDict["deviations"].([]Deviation)
				deviations = append(deviations, Deviation{"CHAINPRECEDENCE", actCouple.Value0, actCouple.Value1})
				traceDict["deviations"] = deviations
			}
		}
	}
}

func checkAltSuccession(trace []string, model map[string]map[string]map[string]int, traceDict map[string]interface{}, actIdxs map[string][]int, parameters map[interface{}]interface{}) {
	if val, ok := model[ALTSUCCESSION]; ok {
		for stringActCouple := range val {
			actCouple := parseString(stringActCouple)
			specIdxs := make([][2]interface{}, 0)
			if contains(trace, actCouple.Value0) {
				for _, i := range actIdxs[actCouple.Value0] {
					specIdxs = append(specIdxs, [2]interface{}{actCouple.Value0, i})
				}
			}
			if contains(trace, actCouple.Value1) {
				for _, i := range actIdxs[actCouple.Value1] {
					specIdxs = append(specIdxs, [2]interface{}{actCouple.Value1, i})
				}
			}
			specIdxs = sortSpecIdxs(specIdxs)
			isOk := true
			for i := 0; i < len(specIdxs); i++ {
				if i%2 == 0 && (specIdxs[i][0] != actCouple.Value0 || i == len(specIdxs)-1 || specIdxs[i+1][0] != actCouple.Value1) {
					isOk = false
					break
				}
			}
			if !isOk {
				deviations := traceDict["deviations"].([]Deviation)
				deviations = append(deviations, Deviation{"ALTSUCCESSION", actCouple.Value0, actCouple.Value1})
				traceDict["deviations"] = deviations
			}
		}
	}
}

func checkChainSuccession(trace []string, model map[string]map[string]map[string]int, traceDict map[string]interface{}, actIdxs map[string][]int, parameters map[interface{}]interface{}) {
	if val, ok := model[CHAINSUCCESSION]; ok {
		for stringActCouple := range val {
			actCouple := parseString(stringActCouple)
			specIdxs := make([][2]interface{}, 0)
			if contains(trace, actCouple.Value0) {
				for _, i := range actIdxs[actCouple.Value0] {
					specIdxs = append(specIdxs, [2]interface{}{actCouple.Value0, i})
				}
			}
			if contains(trace, actCouple.Value1) {
				for _, i := range actIdxs[actCouple.Value1] {
					specIdxs = append(specIdxs, [2]interface{}{actCouple.Value1, i})
				}
			}
			specIdxs = sortSpecIdxs(specIdxs)
			isOk := true
			for i := 0; i < len(specIdxs); i++ {
				if i%2 == 0 && (specIdxs[i][0] != actCouple.Value0 || i == len(specIdxs)-1 || specIdxs[i+1][0] != actCouple.Value1 || specIdxs[i+1][1] != specIdxs[i][1].(int)+1) {
					isOk = false
					break
				}
			}
			if !isOk {
				deviations := traceDict["deviations"].([]Deviation)
				deviations = append(deviations, Deviation{"CHAINSUCCESSION", actCouple.Value0, actCouple.Value1})
				traceDict["deviations"] = deviations
			}
		}
	}
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
func parseString(input string) Tuple {
	parts := strings.Split(input, "<<AND>>")
	if len(parts) != 2 {
		// If there are not exactly 2 parts separated by '<<AND>>',
		// return an empty tuple with empty strings
		return Tuple{}
	}
	return Tuple{
		Value0: strings.TrimSpace(parts[0]),
		Value1: strings.TrimSpace(parts[1]),
	}
}

func max(slice []int) int {
	if len(slice) == 0 {
		panic("empty slice")
	}

	maxVal := slice[0]
	for _, num := range slice {
		if num > maxVal {
			maxVal = num
		}
	}
	return maxVal
}

func min(slice []int) int {
	if len(slice) == 0 {
		panic("empty slice")
	}

	minVal := slice[0]
	for _, num := range slice {
		if num < minVal {
			minVal = num
		}
	}
	return minVal
}

func sortSpecIdxs(specIdxs [][2]interface{}) [][2]interface{} {
	for i := 0; i < len(specIdxs)-1; i++ {
		for j := 0; j < len(specIdxs)-i-1; j++ {
			if specIdxs[j][1].(int) > specIdxs[j+1][1].(int) || (specIdxs[j][1].(int) == specIdxs[j+1][1].(int) && specIdxs[j][0].(string) > specIdxs[j+1][0].(string)) {
				specIdxs[j], specIdxs[j+1] = specIdxs[j+1], specIdxs[j]
			}
		}
	}
	return specIdxs
}
