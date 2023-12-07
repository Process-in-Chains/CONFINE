package xes

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// XES is the root XML structure of the XES event log
/**
type XES struct {
	XMLName      xml.Name     `xml:"log"`
	Extensions   []Extension  `xml:"extension"`
	Global       []Global     `xml:"global"`
	Classifiers  []Classifier `xml:"classifier"`
	StringAttrib []Attribute  `xml:"string"`
	IntAttrib    []Attribute  `xml:"int"`
	DateAttrib   []Attribute  `xml:"date"`
	Traces       []Trace      `xml:"trace"`
}**/
type XES struct {
	XMLName xml.Name `xml:"log"`
	//Extensions  []Extension  `xml:"extension"`
	//Global      []Global     `xml:"global"`
	//Classifiers []Classifier `xml:"classifier"`
	//StringAttrib []Attribute  `xml:"string"`
	//IntAttrib    []Attribute  `xml:"int"`
	//DateAttrib   []Attribute  `xml:"date"`
	Traces []Trace `xml:"trace"`
}

// Extension represents an extension in the XES log
type Extension struct {
	Name   string `xml:"name,attr"`
	Prefix string `xml:"prefix,attr"`
	URI    string `xml:"uri,attr"`
}

// Global represents the global attributes of the XES log
type Global struct {
	Scope      string      `xml:"scope,attr"`
	StringAttr []Attribute `xml:"string"`
	IntAttr    []Attribute `xml:"int"`
	DateAttr   []Attribute `xml:"date"`
}

// Classifier represents a classifier in the XES log
type Classifier struct {
	Name     string      `xml:"name,attr"`
	Keys     []Attribute `xml:"string"`
	Position int         `xml:"position,attr"`
}

// Trace represents an individual trace in the event log
type Trace struct {
	Name   []Attribute `xml:"string" `
	Events []Event     `xml:"event"`
}

func (trace Trace) GetId() (string, error) {
	for _, trace_attr := range trace.Name {
		if trace_attr.Key == "concept:name" {
			return trace_attr.Value, nil
		}
	}
	return "", errors.New("Case Id not defined")
}
func (trace Trace) getTraceSize() (error, int64) {
	traceXML, err := xml.MarshalIndent(trace, "", "")
	if err != nil {
		return err, 0
	}
	return nil, int64(len(string(traceXML)))
}
func (trace Trace) traceToSlice() []string {
	names := make([]string, len(trace.Events))
	for i, event := range trace.Events {
		for _, attr := range event.Attributes {
			if attr.Key == "concept:name" {
				names[i] = attr.Value
				break
			}
		}
	}
	return names
}
func (trace Trace) TraceToByte() []byte {
	emptyXes := XES{}
	emptyXes.Traces = append(emptyXes.Traces, trace)
	byteTrace, err := xml.Marshal(&emptyXes)
	if err != nil {
		log.Fatal(err)
	}
	return byteTrace
}

// Event represents an event within a trace
type Event struct {
	Attributes []Attribute `xml:"string"`
	Timestamp  Attribute   `xml:"date"`
}

// Function to get the attribute from the event from the event
func (event Event) GetAttributeValue(attributeName string) Attribute {
	var attribute Attribute
	for _, attr := range event.Attributes {
		if attr.Key == attributeName {
			attribute = attr
		}
	}
	return attribute

}

func (event Event) getTimestamp() string {
	return event.Timestamp.Value
}

// Attribute represents an attribute within an element
type Attribute struct {
	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr"`
}

func (xes XES) getTrace(caseId string) *Trace {
	for _, trace := range xes.Traces {
		for _, trace_attr := range trace.Name {
			if (trace_attr.Key == "concept:name") && (trace_attr.Value == caseId) {
				return &trace
			}
		}
	}
	return nil
}
func (xes XES) XesToSlices() [][]string {
	var traces [][]string
	for _, t := range xes.Traces {
		traceSlice := t.traceToSlice()
		traces = append(traces, traceSlice)
		//fmt.Println(traces,traceSlice)
	}
	return traces
}
func (xes XES) getHeaderSize() (error, int64) {
	emptyXes := xes
	emptyXes.Traces = xes.Traces[:0]
	xesXML, err := xml.MarshalIndent(emptyXes, "", "")
	if err != nil {
		return err, 0
	}
	return nil, int64(len(string(xesXML)))

}
func ReadXes(filepath string) *XES {
	// Read the XES event log file
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()
	// Read the file content
	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}
	return ParseXes(content)
}
func ParseXes(content []byte) *XES {
	// Initialize XES struct to hold the parsed data
	xes := XES{}
	// Unmarshal the XML content into the XES struct
	err := xml.Unmarshal(content, &xes)
	if err != nil {
		fmt.Println("Error parsing XML:", err)
		return nil
	}
	return &xes
}
