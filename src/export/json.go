package export

import (
	"encoding/json"
	"fmt"
	. "local/BrowserArtifact/src"
	"os"
)

func ExportJSON(path string, artifacts []BrowserArtifact) {
	// Open file
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(artifacts, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Write JSON data to file
	_, err = file.Write(jsonData)
	if err != nil {
		return
	}
}

func ExportJSONLine(path string, artifacts []BrowserArtifact) {
	// Open file
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Write JSON data to file
	for _, artifact := range artifacts {
		jsonData, err := json.Marshal(artifact)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = file.Write(jsonData)
		if err != nil {
			return
		}
	}
}
