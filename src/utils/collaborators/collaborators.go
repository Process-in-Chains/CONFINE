package collaborators

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Reference struct {
	PublicKey    string `json:"public_Key"`
	WebReference string `json:"http_reference"`
	MergeKey     string `json:"merge_key"`
}

func GetReference(webReference string) (Reference, error) {
	references, error := GetReferences()
	if error != nil {
		log.Fatalf("Error getting references: %v", error)
	}
	// Access each field of each entry
	for _, ref := range references {
		if ref.WebReference == webReference {
			return ref, nil
		}
	}
	return Reference{}, fmt.Errorf("Reference not found")
}

func GetReferences() ([]Reference, error) {
	// Read the JSON file
	data, err := ioutil.ReadFile("mining-data/collaborators/process-01/references.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	// Create a variable of the struct type
	var references []Reference

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(data, &references)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return nil, err
	}

	return references, nil
}

func GetReferenceByPublicKey(id string) (Reference, error) {
	// Create a variable of the struct type
	//var references []Reference
	references, error := GetReferences()
	if error != nil {
		log.Fatalf("Error getting references: %v", error)
	}
	// Access each field of each entry
	for _, ref := range references {
		if ref.PublicKey == id {
			return ref, nil
		}
	}
	return Reference{}, fmt.Errorf("Reference not found")

}
