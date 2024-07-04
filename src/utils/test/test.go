package test

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

var STOPMONITORING = false
var TEST_MODE = false

func PrintRamUsage() {
	var ramList = []int{}
	var timestampList = []int{}
	for !STOPMONITORING {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		ramList = append(ramList, int(m.Alloc))
		timestampList = append(timestampList, int(time.Now().UnixMilli()))
		time.Sleep(100 * time.Millisecond)
	}
	// Save RAM usage and timestamp to CSV file
	err := saveToCSV(ramList, timestampList)
	if err != nil {
		fmt.Println("Error saving RAM usage to CSV:", err)
	}
	sum := 0.0
	for _, v := range ramList {
		sum += float64(v)
	}
	avg := sum / float64(len(ramList))
	var firstTimestamp = timestampList[0]
	var lastTimestamp = timestampList[len(timestampList)-1]
	fmt.Printf("TESTMODE - Average RAM Usage in bytes: %.2f\n", avg)
	fmt.Printf("TESTMODE - Test duration in seconds: %d\n", (lastTimestamp-firstTimestamp)/1000)

}
func WaitUntilStop() {
	for !STOPMONITORING {
		time.Sleep(50 * time.Millisecond)
	}
}

func saveToCSV(ramList []int, timestampList []int) error {
	file, err := os.Create("mining-data/output/test" + time.Now().String() + ".csv")
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	// Write header
	err = writer.Write([]string{"Timestamp", "RAM Usage (Bytes)"})
	if err != nil {
		return err
	}
	// Write data rows
	for i := 0; i < len(ramList); i++ {
		err = writer.Write([]string{strconv.Itoa(timestampList[i]), strconv.Itoa(ramList[i])})
		if err != nil {

			return err
		}
	}
	return nil
}
func getSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Add file size if it is not a directory
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}
