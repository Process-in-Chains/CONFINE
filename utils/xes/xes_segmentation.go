package xes

import (
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"
)

const (
	KB = 1024
)

// SplitXESFile splits the XES file into smaller batches based on the given batch size in MB
//
//	func SplitXESFile(filePath string, batchSizeKB int) error {
//		// Open the original XES file
//		file, err := os.Open(filePath)
//		if err != nil {
//			return err
//		}
//		defer file.Close()
//		// Get the file info
//		fileInfo, err := file.Stat()
//		if err != nil {
//			return err
//		}
//		fileSize := fileInfo.Size()
//		fileSizeKB := int(fileSize / KB)
//		numBatches := fileSizeKB/batchSizeKB + 1
//		// Create the output directory for batch files
//		outputDir := strings.TrimSuffix(filePath, ".xes")
//		err = os.Mkdir(outputDir, 0755)
//		if err != nil {
//			return err
//		}
//		// Create a scanner to read the XES file
//		scanner := bufio.NewScanner(file)
//		scanner.Split(bufio.ScanLines)
//		// Initialize the batch counter and file writer
//		batchCounter := 1
//		batchSizeBytes := batchSizeKB * KB
//		batchPath := fmt.Sprintf("%s/batch_%d.xes", outputDir, batchCounter)
//		batchFile, err := os.Create(batchPath)
//		if err != nil {
//			return err
//		}
//		defer batchFile.Close()
//		batchWriter := bufio.NewWriter(batchFile)
//		// Iterate through the XES file and split into batches
//		for scanner.Scan() {
//			line := scanner.Text()
//			// Write the line to the current batch file
//			_, err := batchWriter.WriteString(line + "\n")
//			if err != nil {
//				return err
//			}
//			// Check if the current batch file exceeds the batch size
//			fileInfo, err = batchFile.Stat()
//			if err != nil {
//				return err
//			}
//			if fileInfo.Size() >= int64(batchSizeBytes) {
//				batchWriter.Flush()
//				batchFile.Close()
//				// Increment the batch counter and create a new batch file
//				batchCounter++
//				batchPath = fmt.Sprintf("%s/batch_%d.xes", outputDir, batchCounter)
//				batchFile, err = os.Create(batchPath)
//				if err != nil {
//					return err
//				}
//				batchWriter = bufio.NewWriter(batchFile)
//			}
//		}
//		// Flush and close the last batch file
//		batchWriter.Flush()
//		batchFile.Close()
//		fmt.Printf("Split into %d batches\n", numBatches)
//		return nil
//	}

func SplitXESFile2(targetXes XES, segmentSizeKB int, outputDir string) ([]string, error) {
	emptyXes := targetXes
	var emptyTraces []Trace
	emptyXes.Traces = emptyTraces
	//emptyXesBytes, _ := xml.MarshalIndent(emptyXes, "", "\t")
	emptyXesBytes, _ := xml.MarshalIndent(emptyXes, "", "")
	//emptyXesBytes, _ := xml.Marshal(emptyXes)
	emptyXesString := string(emptyXesBytes)
	segmentSize := segmentSizeKB * 1024
	currentSize := len(emptyXesString)
	segmentNum := 1
	os.Mkdir(outputDir, 0755)
	//fileName := fmt.Sprintf("%s/batch_%d.xes", outputDir, segmentNum)
	xes := targetXes
	var currentTraces []Trace
	var hashList []string
	for _, trace := range xes.Traces {
		traceXML, err := xml.MarshalIndent(trace, "", "")
		if err != nil {
			return nil, err
		}
		traceSize := len(string(traceXML))
		if traceSize > segmentSize {
			fmt.Println("trace exceding the limit of KBs", traceSize/1024)
			continue
		}
		if currentSize+traceSize <= segmentSize {
			currentTraces = append(currentTraces, trace)
		} else {
			hashLog, err := writeSegment(segmentNum, currentTraces, emptyXes, outputDir)
			if err != nil {
				return nil, err
			}
			hashList = append(hashList, string(hashLog))
			currentTraces = nil
			currentSize = len(emptyXesString)
			currentTraces = append(currentTraces, trace)
			segmentNum++
		}
		currentSize += traceSize
	}
	if len(currentTraces) != 0 {
		hashLog, err := writeSegment(segmentNum, currentTraces, emptyXes, outputDir)
		if err != nil {
			return nil, err
		}
		hashList = append(hashList, string(hashLog))
	}
	return hashList, nil
}
func GetTraceSizeList(filePath string) ([][]string, error) {
	eventLog := ReadXes(filePath)
	err, headerSize := eventLog.getHeaderSize()
	if err != nil {
		return nil, err
	}
	headerList := []string{"h", strconv.FormatInt(headerSize, 10)}
	traceSizeList := [][]string{}
	traceSizeList = append(traceSizeList, headerList)
	for _, trace := range eventLog.Traces {
		id, err := trace.GetId()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		err, traceSize := trace.getTraceSize()
		if err != nil {
			return nil, err
		}
		stringTraceSize := strconv.FormatInt(traceSize, 10)
		traceAndSize := []string{id, stringTraceSize}
		traceSizeList = append(traceSizeList, traceAndSize)
	}
	return traceSizeList, nil
}
func writeSegment(batchNum int, traces []Trace, emptyXes XES, outputDir string) ([]byte, error) {
	emptyXes.Traces = traces
	stringXES, _ := xml.MarshalIndent(emptyXes, "", "")
	hashSHA256 := sha256.Sum256(stringXES)
	fileName := fmt.Sprintf("%s/batch_%d.xes", outputDir, batchNum)
	return hashSHA256[:], ioutil.WriteFile(fileName, stringXES, 0644)
}
func GetTraceSize(filePath string) (map[string]string, error) {
	eventLog := ReadXes(filePath)
	err, headerSize := eventLog.getHeaderSize()
	if err != nil {
		return nil, err
	}
	traceSizeMap := map[string]string{"h": strconv.FormatInt(headerSize, 10)}
	for _, trace := range eventLog.Traces {
		id, err := trace.GetId()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		err, traceSize := trace.getTraceSize()
		if err != nil {
			return nil, err
		}
		stringTraceSize := strconv.FormatInt(traceSize, 10)
		traceSizeMap[id] = stringTraceSize
	}
	return traceSizeMap, nil
}

func MergeTraces(trace1 Trace, trace2 Trace) (Trace, error) {
	var eventList []Event
	var eventList2 []Event
	eventList = append(eventList, trace1.Events...)
	eventList2 = append(eventList2, trace2.Events...)
	eventList = append(eventList, eventList2...)
	sort.Slice(eventList, func(i, j int) bool {
		time1, _ := time.Parse("2006-01-02T15:04:05Z", eventList[i].Timestamp.Value)
		time2, _ := time.Parse("2006-01-02T15:04:05Z", eventList[j].Timestamp.Value)
		return time1.Before(time2)
	})
	newTrace := Trace{Name: trace1.Name, Events: eventList}
	return newTrace, nil

}
