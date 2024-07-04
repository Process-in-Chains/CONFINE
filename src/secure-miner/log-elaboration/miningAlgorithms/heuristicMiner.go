package prosessMiningAlgorithms

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

// var occurrencyThreshold = 1
var occurrencyThreshold = 10
var dependencyThreshold = 0.3
var andThreshold = 0.1
var longDepencency2Threshold = 0.99 //Usually it's high (almost 1)
var longDistance = 2
var DEBUG = false
var depLongDist = false

type dependencyMeasure struct {
	dep        float32
	occurrence int
}

type dependencyGraphAdjacencies struct {
	nodeName string
	info     dependencyMeasure
}

type causalMatrixLists struct {
	input  []string
	output []string
}

type causalMatrix struct {
	input  *SetOfSets
	output *SetOfSets
}

/*
Event Map
*/
func getKeyFromValue(eventMap map[string]int, x int) string {
	for i, v := range eventMap {
		if v == x {
			return i
		}
	}
	return ""
}

func getMaxValue(m map[string]int) int {
	//Set to -1 for start the index from 0
	max := -1
	for _, v := range m {
		if v > max {
			max = v
		}
	}
	return max
}

func addEventToMap(eventMap map[string]int, e string) {
	_, ok := eventMap[e]
	// If the key exists
	if ok == false {
		eventMap[e] = getMaxValue(eventMap) + 1
	}
	//sortMap(eventMap)
}

func sortMap(eventMap map[string]int) {
	var keys []string
	for i := range eventMap {
		keys = append(keys, i)
	}
	sort.Strings(keys)

	for i, v := range keys {
		eventMap[v] = i
	}
}

/*
Dependecy Matrix
*/

func printDependencyMatrix(dependencyMatrix [][]dependencyMeasure, eventMap map[string]int, mode string) {
	for i := 0; i < len(dependencyMatrix); i++ {
		switch mode {
		case "full":
			fmt.Printf("\n%v) %v ", i, getKeyFromValue(eventMap, i))
			for j := 0; j < len(dependencyMatrix[i]); j++ {
				fmt.Printf(" %v:%v", j, dependencyMatrix[i][j])
			}
		case "dep":
			fmt.Printf("\n%v) %v ", i, getKeyFromValue(eventMap, i))
			for k := 0; k <= getMaxValue(eventMap); k++ {
				fmt.Printf("  %.3f  ", dependencyMatrix[i][k].dep)
			}
		case "occur":
			fmt.Printf("\n%v    ", getKeyFromValue(eventMap, i))
			for k := 0; k <= getMaxValue(eventMap); k++ {
				fmt.Printf("  %v  ", dependencyMatrix[i][k].occurrence)
			}
		}
	}
}

func createDependencyMatrix(n int) [][]dependencyMeasure {
	m := make([][]dependencyMeasure, n)
	for i := range m {
		m[i] = make([]dependencyMeasure, n)
	}
	return m
}

func calculateDependencyMatrix(dependencyMatrix [][]dependencyMeasure, event []string, eventMap map[string]int, n int) {
	for i := range event {
		if i+n < len(event) {
			dependencyMatrix[eventMap[event[i]]][eventMap[event[i+n]]].occurrence++
			if event[i] != event[i+n] {
				dependencyMatrix[eventMap[event[i]]][eventMap[event[i+n]]].dep = float32(dependencyMatrix[eventMap[event[i]]][eventMap[event[i+n]]].occurrence-dependencyMatrix[eventMap[event[i+n]]][eventMap[event[i]]].occurrence) / float32(dependencyMatrix[eventMap[event[i]]][eventMap[event[i+n]]].occurrence+dependencyMatrix[eventMap[event[i+n]]][eventMap[event[i]]].occurrence+1)
				dependencyMatrix[eventMap[event[i+n]]][eventMap[event[i]]].dep = float32(dependencyMatrix[eventMap[event[i+n]]][eventMap[event[i]]].occurrence-dependencyMatrix[eventMap[event[i]]][eventMap[event[i+n]]].occurrence) / float32(dependencyMatrix[eventMap[event[i+n]]][eventMap[event[i]]].occurrence+dependencyMatrix[eventMap[event[i]]][eventMap[event[i+n]]].occurrence+1)
			} else {
				dependencyMatrix[eventMap[event[i]]][eventMap[event[i+n]]].dep = float32(dependencyMatrix[eventMap[event[i]]][eventMap[event[i+n]]].occurrence) / float32(dependencyMatrix[eventMap[event[i]]][eventMap[event[i+n]]].occurrence+1)
			}
		}
	}
}

func calculateReverseDependencyMatrix(reverseDependencyMatrix [][]dependencyMeasure, event []string, eventMap map[string]int) {
	for i := len(event) - 1; i >= 0; i-- {
		if i-1 >= 0 {
			reverseDependencyMatrix[eventMap[event[i]]][eventMap[event[i-1]]].occurrence++
			if event[i] != event[i-1] {
				reverseDependencyMatrix[eventMap[event[i]]][eventMap[event[i-1]]].dep = float32(reverseDependencyMatrix[eventMap[event[i]]][eventMap[event[i-1]]].occurrence-reverseDependencyMatrix[eventMap[event[i-1]]][eventMap[event[i]]].occurrence) / float32(reverseDependencyMatrix[eventMap[event[i]]][eventMap[event[i-1]]].occurrence+reverseDependencyMatrix[eventMap[event[i-1]]][eventMap[event[i]]].occurrence+1)
				reverseDependencyMatrix[eventMap[event[i-1]]][eventMap[event[i]]].dep = float32(reverseDependencyMatrix[eventMap[event[i-1]]][eventMap[event[i]]].occurrence-reverseDependencyMatrix[eventMap[event[i]]][eventMap[event[i-1]]].occurrence) / float32(reverseDependencyMatrix[eventMap[event[i-1]]][eventMap[event[i]]].occurrence+reverseDependencyMatrix[eventMap[event[i]]][eventMap[event[i-1]]].occurrence+1)
			} else {
				reverseDependencyMatrix[eventMap[event[i]]][eventMap[event[i-1]]].dep = float32(reverseDependencyMatrix[eventMap[event[i]]][eventMap[event[i-1]]].occurrence) / float32(reverseDependencyMatrix[eventMap[event[i]]][eventMap[event[i-1]]].occurrence+1)
			}
		}
	}
}

func calculateReverseMatrix(dependencyMatrix [][]dependencyMeasure, eventMap map[string]int) [][]dependencyMeasure {
	reverse := createDependencyMatrix(len(dependencyMatrix))
	for r := range dependencyMatrix {
		for c := range dependencyMatrix[r] {
			reverse[r][c] = dependencyMatrix[c][r]
			reverse[c][r] = dependencyMatrix[r][c]
		}
	}
	return reverse
}

func expandDependencyMatrix(dependencyMatrix [][]dependencyMeasure, d1 int, eventMap map[string]int) [][]dependencyMeasure {
	expandedMatrix := createDependencyMatrix(len(eventMap))
	for i := 0; i < d1; i++ {
		for j := 0; j < d1; j++ {
			expandedMatrix[i][j] = dependencyMatrix[i][j]
		}
	}
	return expandedMatrix
}

// TODO: replace with a thing that have minor computational complexity
func allOccurrency(dependencyMatrix [][]dependencyMeasure, eventMap map[string]int) map[string]int {
	occur := make(map[string]int)
	for r := range dependencyMatrix {
		s := 0
		for c := range dependencyMatrix[r] {
			s += dependencyMatrix[r][c].occurrence
		}
		if s != 0 {
			occur[getKeyFromValue(eventMap, r)] = s
		} else {
			ss := 0
			for x := range dependencyMatrix {
				ss += dependencyMatrix[x][r].occurrence
			}
			occur[getKeyFromValue(eventMap, r)] = ss
		}
	}
	//fmt.Println(occur)
	return occur
}

/*
Dependecy Graph
*/
func calculateDependencyGraph(dependencyMatrix [][]dependencyMeasure, eventMap map[string]int) map[string][]dependencyGraphAdjacencies {
	AdjacencyList := make(map[string][]dependencyGraphAdjacencies)
	for r := range dependencyMatrix {
		for c := range dependencyMatrix[r] {
			if dependencyMatrix[r][c].dep >= float32(dependencyThreshold) && dependencyMatrix[r][c].occurrence >= occurrencyThreshold {
				node := dependencyGraphAdjacencies{nodeName: getKeyFromValue(eventMap, c), info: dependencyMeasure{dep: dependencyMatrix[r][c].dep, occurrence: dependencyMatrix[r][c].occurrence}}
				AdjacencyList[getKeyFromValue(eventMap, r)] = append(AdjacencyList[getKeyFromValue(eventMap, r)], node)
			}
		}
	}
	return AdjacencyList
}

/*
Causal Matrix
*/
func printCasualMatrix(cm map[string]causalMatrix) {
	if DEBUG {
		fmt.Println("(A) Activity (I) Input (O) output")
		fmt.Println("\nA \t I \t O")
		fmt.Println("--------------------")
	}
	for el := range cm {
		fmt.Printf("\n\n%v > ", el)
		for _, v := range cm[el].input.Data {

			//fmt.Printf("%v ", v)
			fmt.Printf("%v ", v)
		}
		fmt.Printf("  |||  ")
		for _, v := range cm[el].output.Data {
			fmt.Printf("%v ", v)
		}
		fmt.Printf("\n")
	}
}

func createCasualMatrix(dependencyGraph map[string][]dependencyGraphAdjacencies, eventMap map[string]int) map[string]causalMatrixLists {
	casualMatrix := make(map[string]causalMatrixLists)
	e := calculateEvent(dependencyGraph)
	//fmt.Println(e)
	slash := []string{"/"}
	for r := range e {
		casualMatrix[r] = causalMatrixLists{input: slash, output: slash}
	}
	//fmt.Println("\n\n\n CAUSALLLLLLL ", casualMatrix, "\n\n\n")
	return casualMatrix
}

func calculateEvent(dependencyGraph map[string][]dependencyGraphAdjacencies) map[string]string {
	ev := make(map[string]string)
	for el := range dependencyGraph {
		ev[el] = ""
		for e := range dependencyGraph[el] {
			ev[dependencyGraph[el][e].nodeName] = ""
		}
	}
	return ev
}

func inputOutputCasualMatrix(eventMap map[string]int, dependencyGraph map[string][]dependencyGraphAdjacencies) map[string]causalMatrixLists {
	casualMatrix := createCasualMatrix(dependencyGraph, eventMap)
	slash := []string{"/"}

	for el := range dependencyGraph {
		for e := range dependencyGraph[el] {
			reflect.DeepEqual(casualMatrix[el].output, slash)
			//if casualMatrix[el].output != slash {
			if !reflect.DeepEqual(casualMatrix[el].output, slash) {
				casualMatrix[el] = causalMatrixLists{input: casualMatrix[el].input, output: append(casualMatrix[el].output, dependencyGraph[el][e].nodeName)}
			} else {
				casualMatrix[el] = causalMatrixLists{input: casualMatrix[el].input, output: []string{dependencyGraph[el][e].nodeName}}
			}
			//if casualMatrix[dependencyGraph[el][e].nodeName].input != "/" {
			if !reflect.DeepEqual(casualMatrix[dependencyGraph[el][e].nodeName].input, slash) {
				casualMatrix[dependencyGraph[el][e].nodeName] = causalMatrixLists{input: append(casualMatrix[dependencyGraph[el][e].nodeName].input, el), output: casualMatrix[dependencyGraph[el][e].nodeName].output}
			} else {
				casualMatrix[dependencyGraph[el][e].nodeName] = causalMatrixLists{input: []string{el}, output: casualMatrix[dependencyGraph[el][e].nodeName].output}
			}
		}
	}
	return casualMatrix
}

func calculateCasualMatrixValue(dependencyMatrix [][]dependencyMeasure, eventMap map[string]int, dep []string, el string) *SetOfSets {
	setOfSets := NewSetOfSets()
	for _, d := range dep {
		set1 := NewSet()
		set1.Add(d)
		for j := 0; j < len(dep); j++ {
			if d != dep[j] {
				calc := float32(dependencyMatrix[eventMap[d]][eventMap[dep[j]]].occurrence+dependencyMatrix[eventMap[dep[j]]][eventMap[d]].occurrence) / float32(dependencyMatrix[eventMap[el]][eventMap[d]].occurrence+dependencyMatrix[eventMap[el]][eventMap[dep[j]]].occurrence+1)
				if calc <= float32(andThreshold) {
					set1.Add(dep[j])
				}
			}
		}
		setOfSets.Add(set1)
	}
	return setOfSets
}

func calculateCasualMatrix(dependencyMatrix [][]dependencyMeasure, reverseDependencyMatrix [][]dependencyMeasure, eventMap map[string]int, dependencyGraph map[string][]dependencyGraphAdjacencies) map[string]causalMatrix {
	casualMatrixIOLists := inputOutputCasualMatrix(eventMap, dependencyGraph)
	casualMatrix := make(map[string]causalMatrix)
	for el := range casualMatrixIOLists {
		if len(casualMatrixIOLists[el].output) > 1 {
			//dep := strings.Split(casualMatrixIOLists[el].output, "")
			//casualMatrix[el] = causalMatrix{input: casualMatrix[el].input, output: calculateCasualMatrixValue(dependencyMatrix, eventMap, dep, el)}
			casualMatrix[el] = causalMatrix{input: casualMatrix[el].input, output: calculateCasualMatrixValue(dependencyMatrix, eventMap, casualMatrixIOLists[el].output, el)}
		} else {
			s := NewSet()
			//s.Add(casualMatrixIOLists[el].output)
			s.Add(casualMatrixIOLists[el].output[0])
			ss := NewSetOfSets()
			ss.Add(s)
			casualMatrix[el] = causalMatrix{input: casualMatrix[el].input, output: ss}
		}
		if len(casualMatrixIOLists[el].input) > 1 {
			//dep := strings.Split(casualMatrixIOLists[el].input, "")
			//casualMatrix[el] = causalMatrix{input: calculateCasualMatrixValue(reverseDependencyMatrix, eventMap, dep, el), output: casualMatrix[el].output}
			casualMatrix[el] = causalMatrix{input: calculateCasualMatrixValue(reverseDependencyMatrix, eventMap, casualMatrixIOLists[el].input, el), output: casualMatrix[el].output}
		} else {
			s := NewSet()
			//s.Add(casualMatrixIOLists[el].input)
			s.Add(casualMatrixIOLists[el].input[0])
			ss := NewSetOfSets()
			ss.Add(s)
			casualMatrix[el] = causalMatrix{input: ss, output: casualMatrix[el].output}
		}
	}
	return casualMatrix
}

/*
Map <-> JSON
*/
// SaveMapToJSON saves a Go map to a JSON file
func SaveMapToJSON(data map[string]int, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// LoadMapFromJSON loads a Go map from a JSON file
func LoadMapFromJSON(filename string) (map[string]int, error) {
	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var data map[string]int
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

/*
Dependecy matrix <-> Json
*/
//MarshalJSON marshals a dependencyMeasure object to JSON.
func (dep dependencyMeasure) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Dep        string `json:"dep"`
		Occurrence string `json:"occurrence"`
	}{
		Dep:        fmt.Sprintf("%f", dep.dep),
		Occurrence: strconv.Itoa(dep.occurrence),
	})
}

// UnmarshalJSON unmarshals JSON data into a dependencyMeasure object.
func (dep *dependencyMeasure) UnmarshalJSON(data []byte) error {
	temp := struct {
		Dep        interface{} `json:"dep"`
		Occurrence interface{} `json:"occurrence"`
	}{}
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}

	switch depValue := temp.Dep.(type) {
	case float64:
		dep.dep = float32(depValue)
	case string:
		parsedDep, err := strconv.ParseFloat(depValue, 32)
		if err != nil {
			return err
		}
		dep.dep = float32(parsedDep)
	default:
		return errors.New("unexpected type for dep field")
	}

	switch occurrenceValue := temp.Occurrence.(type) {
	case float64:
		dep.occurrence = int(occurrenceValue)
	case string:
		parsedOccurrence, err := strconv.Atoi(occurrenceValue)
		if err != nil {
			return err
		}
		dep.occurrence = parsedOccurrence
	default:
		return errors.New("unexpected type for occurrence field")
	}

	return nil
}

// saveToJSON saves the data to a JSON file.
func saveToJSON(data [][]dependencyMeasure, filename string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// readFromJSON reads the data from a JSON file.
func readFromJSON(filename string) ([][]dependencyMeasure, error) {
	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var data [][]dependencyMeasure
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func TracesMapper(eventMap map[string]int, events [][]string, filePath string, processName string) [][]string {

	//filePath := "Log/L_prova4.txt"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v", err)
		return nil
	}

	// Split content by newline characters
	lines := strings.Split(string(content), "\r\n")

	// Print each line
	for _, line := range lines {
		line = strings.TrimSpace(line)
		el := strings.Split(string(line), " ")
		for _, e := range el[2:] {
			addEventToMap(eventMap, e)
		}
		num, _ := strconv.Atoi(el[0][:len(el[0])-1])
		//fmt.Println(num, " ", el[2:])
		for i := 0; i < num; i++ {
			events = append(events, el[2:])
		}
	}
	// Save the map to JSON
	er := SaveMapToJSON(eventMap, "./mining-data/consumption-data/"+processName+"/miningMetadata/map.json")
	if er != nil {
		fmt.Println("Error saving to JSON:", er)
		return nil
	}
	return events
}

func causalMatrixtoPnml(cm map[string]causalMatrix, processName string) (bool, PNML) {
	// Create an initial Petri net
	pnml := CreateInitialPetriNet()

	for el := range cm {
		idEl := strings.ReplaceAll(el, " ", "_")
		//fmt.Printf("\n \n \n --------- %v --------- \n", el)
		if !TransitionExists(pnml, el) {
			id := strings.ReplaceAll(el, " ", "_")
			AddTransition(&pnml, id, el)
		}
		//Input
		for _, v := range cm[el].input.Data {
			//idEl := strings.ReplaceAll(el, " ", "_")
			if v.String() == "/" {
				AddPlace(&pnml, "P0", "P0")
				AddArc(&pnml, NextArcID(pnml, "A"), "P0", idEl)
			}
			if len(v.data) > 1 {
				np := NextPlaceID(pnml, "P")
				AddPlace(&pnml, np, np)
				AddArc(&pnml, NextArcID(pnml, "A"), np, idEl)
				for k, _ := range v.data {
					idk := strings.ReplaceAll(k, " ", "_")
					if !TransitionExists(pnml, k) {
						AddTransition(&pnml, idk, k)
					}
					AddArc(&pnml, NextArcID(pnml, "A"), idk, np)
				}
			}
		}
		//output
		for _, v := range cm[el].output.Data {
			//idEl := strings.ReplaceAll(el, " ", "_")
			//fmt.Printf("%v ", v)
			if len(v.data) == 1 && v.String() == "/" {
				np := NextPlaceID(pnml, "P")
				AddPlace(&pnml, np, np)
				AddArc(&pnml, NextArcID(pnml, "A"), idEl, np)
			}
			//forse >1 e ==1 si può unire ed è simile dal momento che ==1 farebbe il ciclo una sola volta
			if len(v.data) > 1 {
				np := NextPlaceID(pnml, "P")
				AddPlace(&pnml, np, np)
				AddArc(&pnml, NextArcID(pnml, "A"), idEl, np)
				for k, _ := range v.data {
					idk := strings.ReplaceAll(k, " ", "_")
					if !TransitionExists(pnml, k) {
						AddTransition(&pnml, idk, k)
					}
					AddArc(&pnml, NextArcID(pnml, "A"), np, idk)
				}
			}
			if len(v.data) == 1 && v.String() != "/" {
				s := NewSet()
				s.Add(el)
				if cm[v.String()].input.Contains(s) {
					id := strings.ReplaceAll(v.String(), " ", "_")
					if !TransitionExists(pnml, v.String()) {
						AddTransition(&pnml, id, v.String())
					}
					np := NextPlaceID(pnml, "P")
					AddPlace(&pnml, np, np)
					AddArc(&pnml, NextArcID(pnml, "A"), idEl, np)
					AddArc(&pnml, NextArcID(pnml, "A"), np, id)
				}
			}
		}
		//fmt.Printf("\n")
	}
	// Generate PNML XML data
	xmlData, err := xml.MarshalIndent(pnml, "", "    ")
	if err != nil {
		fmt.Println("Error generating PNML XML:", err)
		return false, pnml
	}
	// Create and write PNML file
	file, err := os.Create("./mining-data/output/" + "heuristicsMiner_output.pnml")
	//file, err := os.Create("output/" + processName + "_" + strconv.Itoa(int(timestamp)) + ".pnml")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return false, pnml
	}
	defer file.Close()

	file.Write(xmlData)
	return true, pnml
}

/*
LONG DEP2
*/
func searchNewDependecyLongDistance2(dependencyMatrix2 [][]dependencyMeasure, eventMap map[string]int, occ map[string]int, pnml *PNML, processName string) (int, bool) {
	n := 0
	if DEBUG {
		printDependencyMatrix(dependencyMatrix2, eventMap, "full")
	}
	for r := range dependencyMatrix2 {
		strR := getKeyFromValue(eventMap, r)
		for c := range dependencyMatrix2[r] {
			strC := getKeyFromValue(eventMap, c)
			//fmt.Printf("\n r-> %v (%v)", r, getKeyFromValue(eventMap, r))
			//fmt.Printf("\t c-> %v (%v)", c, getKeyFromValue(eventMap, c))
			//fmt.Printf("\t > %v ", dependencyMatrix2[r][c])
			//fmt.Printf("\t - %v ", float32(dependencyMatrix2[r][c].occurrence)/float32(occ[strR]+1))
			//fmt.Printf("\t > %v ", float32(dependencyMatrix2[r][c].occurrence)/float32(occ[strR]+1) > float32(longDepencency2Threshold))
			//fmt.Printf("\t > %v \n", !FindAPathWithExclusion(FindAllPaths(pnml), strR, strC))
			if float32(dependencyMatrix2[r][c].occurrence)/float32(occ[strR]+1) > float32(longDepencency2Threshold) && !FindAPathWithExclusion(FindAllPaths(pnml), strR, strC) {

				n += 1
				idR := strings.ReplaceAll(strR, " ", "_")
				idC := strings.ReplaceAll(strC, " ", "_")
				if !TransitionExists(*pnml, strR) {
					AddTransition(pnml, idR, strR)
				}
				if !TransitionExists(*pnml, strC) {
					AddTransition(pnml, idC, strC)
				}
				np := NextPlaceID(*pnml, "P")
				//fmt.Println("created new place -> ", np)
				AddPlace(pnml, np, np)
				AddArc(pnml, NextArcID(*pnml, "A"), idR, np)
				AddArc(pnml, NextArcID(*pnml, "A"), np, idC)
			}
		}
	}

	// Generate PNML XML data
	xmlData, err := xml.MarshalIndent(pnml, "", "    ")
	if err != nil {
		fmt.Println("Error generating PNML XML:", err)
		return n, false
	}

	// Create and write PNML file
	//file, err := os.Create("output/" + processName + "_" + strconv.Itoa(int(timestamp)) + ".pnml")
	file, err := os.Create("./mining-data/output/" + "heuristicsMiner_output.pnml")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return n, false
	}
	defer file.Close()

	file.Write(xmlData)
	return n, true
}

func scanEvent(events [][]string, eventMap map[string]int, processName string) {
	for _, ele := range events {
		for _, el := range ele {
			addEventToMap(eventMap, el)
		}
	}
	//er := SaveMapToJSON(eventMap, "/data/miningMetadata/"+processName+"/map.json")
	er := SaveMapToJSON(eventMap, "./mining-data/consumption-data/"+processName+"/miningMetadata/map.json")
	if er != nil {
		fmt.Println("Error saving to JSON:", er)
	}
}

func HeuristicMiner(ev [][]string, processName string) {
	//fmt.Println("-------------------------------\n", ev)
	var d1, d2 int
	var dependencyMatrix [][]dependencyMeasure
	eventMap, e := LoadMapFromJSON("./mining-data/consumption-data/" + processName + "/miningMetadata/map.json")
	d1 = len(eventMap)

	if e != nil {
		fmt.Println("Error loading EventMap from JSON:", e)
		return
	}

	if d1 == 0 {
		eventMap = make(map[string]int)
		scanEvent(ev, eventMap, processName)
		dependencyMatrix = createDependencyMatrix(getMaxValue(eventMap) + 1)
	} else {
		scanEvent(ev, eventMap, processName)
		d2 = len(eventMap)

		// Read the data from the JSON file
		dm, err := readFromJSON("./mining-data/consumption-data/" + processName + "/miningMetadata/dependencyMatrix.json")
		dependencyMatrix = dm
		if err != nil {
			log.Fatal(err)
		}

		if d1 != d2 {
			dependencyMatrix = expandDependencyMatrix(dependencyMatrix, d1, eventMap)
		}
	}
	for _, e := range ev {
		calculateDependencyMatrix(dependencyMatrix, e, eventMap, 1)
	}
	// Save the dependencyMatrix to a JSON file
	e = saveToJSON(dependencyMatrix, "./mining-data/consumption-data/"+processName+"/miningMetadata/dependencyMatrix.json")
	if e != nil {
		log.Fatal(e)
	}

	//fmt.Println("\n\nDEPENDENCY MATRIX\n")
	//printDependencyMatrix(dependencyMatrix, eventMap, "full")

	reverseDependencyMatrix := calculateReverseMatrix(dependencyMatrix, eventMap)
	dependencyGraph := calculateDependencyGraph(dependencyMatrix, eventMap)

	//fmt.Println("\n\nDEPENDENCY GRAPH\n", dependencyGraph)
	causalM := calculateCasualMatrix(dependencyMatrix, reverseDependencyMatrix, eventMap, dependencyGraph)

	//fmt.Println("\n\nCAUSAL MATRIX\n")
	//printCasualMatrix(causalM)
	//printCasualMatrix(causalM)
	//fmt.Println("-------------------")
	//printDependencyMatrix(dependencyMatrix, eventMap, "full")
	pnmlCreation, pnml := causalMatrixtoPnml(causalM, processName)

	if pnmlCreation && depLongDist {
		var dependencyMatrix2 [][]dependencyMeasure
		var occ map[string]int
		if d1 == 0 {
			dependencyMatrix2 = createDependencyMatrix(getMaxValue(eventMap) + 1)
		} else {
			scanEvent(ev, eventMap, processName)
			d2 = len(eventMap)

			dpm2, e := readFromJSON("./mining-data/consumption-data/" + processName + "/miningMetadata/dependencyMatrix" + strconv.Itoa(longDistance) + "len" + ".json")
			if e != nil {
				fmt.Println("Error loading Dependency matrix long ", longDistance, " from JSON:", e)
				return
			}
			dependencyMatrix2 = dpm2
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			if d1 != d2 {
				dependencyMatrix2 = expandDependencyMatrix(dependencyMatrix2, d1, eventMap)
			}
		}
		for _, e := range ev {
			calculateDependencyMatrix(dependencyMatrix2, e, eventMap, longDistance)
		}
		e = saveToJSON(dependencyMatrix2, "./mining-data/consumption-data/"+processName+"/miningMetadata/dependencyMatrix"+strconv.Itoa(longDistance)+"len"+".json")
		if e != nil {
			log.Fatal(e)
		}
		occ = allOccurrency(dependencyMatrix, eventMap)
		//printDependencyMatrix(dependencyMatrix2, eventMap, "full")
		n, success := searchNewDependecyLongDistance2(dependencyMatrix2, eventMap, occ, &pnml, processName)
		if DEBUG {
			if success {
				fmt.Println("\n\nN place added with distance=2: ", n)
			}
		}
	}
}
