package main

import (
	"flag"
	"fmt"
	. "local/BrowserArtifact/src"
	. "local/BrowserArtifact/src/browsers/chrome"
	. "local/BrowserArtifact/src/browsers/firefox"
	. "local/BrowserArtifact/src/export"
	"os"
	"runtime"
	"sort"
	"time"
)

// Define command line arguments

var OsName string
var browserArg string

var startDateString string
var endDateString string

var outputDirectory string
var fileBaseName string
var outputFormat string

var startDate time.Time
var endDate time.Time

var profile string
var verboseLevel string
var logFile string

func init() {
	// Define command line arguments
	flag.StringVar(&browserArg, "browser", "all", "Browser: chrome, firefox")
	flag.StringVar(&outputDirectory, "output_directory", ".", "Output Directory")
	flag.StringVar(&fileBaseName, "file_base_name", "BrowserArtifacts", "File Base Name")
	flag.StringVar(&outputFormat, "format", "json", "Output Format: json, json_line, csv")

	flag.StringVar(&startDateString, "start_date", "2000-01-01", "Start Date")
	flag.StringVar(&endDateString, "end_date", "now", "End Date")
	flag.StringVar(&logFile, "log_file", "", "Log File")
	flag.StringVar(&verboseLevel, "verbose", "info", "Verbose Level: debug, info, warn, error")

	flag.StringVar(&profile, "profile", "all", "User Profile")

	flag.Parse()
}

func log(level string, source string, message string) {
	Log.Log(level, "main", source, message)
}

func findProfile() []string {
	var foundProfile []string
	switch OsName {
	case "windows":
		//List directory in C:/Users
		dir, err := os.ReadDir("C:\\Users")
		if err != nil {
			return foundProfile
		}
		for _, entry := range dir {
			if entry.IsDir() {
				foundProfile = append(foundProfile, entry.Name())
			}

		}
	case "darwin":
		// Mac
		dir, err := os.ReadDir("/Users")
		if err != nil {
			return foundProfile
		}
		for _, entry := range dir {
			if entry.IsDir() {
				foundProfile = append(foundProfile, entry.Name())
			}
		}

	case "linux":
		// Linux
		dir, err := os.ReadDir("/home")
		if err != nil {
			return foundProfile
		}
		for _, entry := range dir {
			if entry.IsDir() {
				foundProfile = append(foundProfile, entry.Name())
			}
		}

	}

	return foundProfile
}

func findInstalledBrowser() []string {
	var foundBrowser []string

	switch OsName {
	case "windows":
		//Check if Chrome is installed
		if _, err := os.Stat("C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"); err == nil {
			foundBrowser = append(foundBrowser, "chrome")
		}
		//Check if firefox is installed
		if _, err := os.Stat("C:\\Program Files\\Mozilla firefox\\firefox.exe"); err == nil {
			foundBrowser = append(foundBrowser, "firefox")
		}
	case "darwin":
		// Mac
		//Check if Chrome is installed
		if _, err := os.Stat("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"); err == nil {
			foundBrowser = append(foundBrowser, "chrome")
		}
		//Check if firefox is installed
		if _, err := os.Stat("/Applications/firefox.app/Contents/MacOS/firefox"); err == nil {
			foundBrowser = append(foundBrowser, "firefox")
		}
	case "linux":
		// Linux
		//Check if Chrome is installed
		if _, err := os.Stat("/usr/bin/google-chrome"); err == nil {
			foundBrowser = append(foundBrowser, "chrome")
		}
		//Check if firefox is installed
		if _, err := os.Stat("/usr/bin/firefox"); err == nil {
			foundBrowser = append(foundBrowser, "firefox")
		}
	}

	return foundBrowser
}

func debugTimestamp(artifacts []BrowserArtifact) {
	// Debug epoch length
	const epochSec = 1000000000

	// Dictionnary
	var dict = make(map[string]string)

	for _, artifact := range artifacts {
		value, ok := dict[artifact.ArtifactType]
		if !ok {
			if artifact.Timestamp >= epochSec*1000*1000*1000 {
				if value != "nano" {
					fmt.Println("WARN: Inconsistent timestamp format for ", artifact.ArtifactType, " new: nano, old: ", value)
				}
				dict[artifact.ArtifactType] = "nano"
			} else if artifact.Timestamp >= epochSec*1000*1000 {
				if value != "micro" {
					fmt.Println("WARN: Inconsistent timestamp format for ", artifact.ArtifactType, " new: micro, old: ", value)
				}
				dict[artifact.ArtifactType] = "micro"
			} else if artifact.Timestamp >= epochSec*1000 {
				if value != "milli" {
					fmt.Println("WARN: Inconsistent timestamp format for ", artifact.ArtifactType, " new: milli, old: ", value)
				}
				dict[artifact.ArtifactType] = "milli"
			} else {
				if value != "sec" {
					fmt.Println("WARN: Inconsistent timestamp format for ", artifact.ArtifactType, " new: sec, old: ", value)
				}
				dict[artifact.ArtifactType] = "sec"
			}
		}
	}

	fmt.Println("Timestamp format: ", dict)

}

func FilterArtifacts(artifacts []BrowserArtifact, startDate time.Time, endDate time.Time) []BrowserArtifact {
	var filteredArtifacts []BrowserArtifact
	//debugTimestamp(artifacts)

	for _, artifact := range artifacts {
		// Convert ms to s, ms
		timestamp_sec := artifact.Timestamp / 1000000
		timestamp_nano := (artifact.Timestamp % 1000000) * 1000

		timestamp := time.Unix(int64(timestamp_sec), int64(timestamp_nano))
		//fmt.Println(artifact.ArtifactType, ": ", artifact.Timestamp, ":", timestamp.Format("2006-01-02"))

		if timestamp.Unix() >= startDate.Unix() && timestamp.Unix() <= endDate.Unix() {
			filteredArtifacts = append(filteredArtifacts, artifact)
		}
	}

	return filteredArtifacts
}

func argVerify() bool {
	isValid := true

	if verboseLevel != "debug" && verboseLevel != "info" && verboseLevel != "warn" && verboseLevel != "error" {
		fmt.Println("Invalid verbose level: ", verboseLevel)
		isValid = false
	}

	if outputFormat != "json" && outputFormat != "json_line" && outputFormat != "csv" {
		fmt.Println("Invalid output format: ", outputFormat)
		isValid = false
	}

	// Check if dates are in correct format "YYYY-MM-DD"
	if startDateString != "now" {
		_, err := time.Parse("2006-01-02", startDateString)
		if err != nil {
			fmt.Println("Invalid start date format: ", startDateString)
			isValid = false
		}
	}

	if endDateString != "now" {
		_, err := time.Parse("2006-01-02", endDateString)
		if err != nil {
			fmt.Println("Invalid end date format: ", endDateString)
			isValid = false
		}
	}

	return isValid
}

// Main function
func main() {
	if !argVerify() {
		return
	}

	// Set log level
	Log.SetLevel(verboseLevel)

	// Set log file
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("Failed to open log file: ", logFile)
			return
		}
		Log.SetOutput(file)
	} else {
		Log.SetOutput(os.Stdout)
	}

	OsName = runtime.GOOS
	log("info", "main", "OS: "+OsName)

	// Parse Date String
	if endDateString == "now" {
		// Log Debug
		endDate = time.Now()
	} else {
		endDate, _ = time.Parse("2006-01-02", endDateString)
	}

	startDate, _ = time.Parse("2006-01-02", startDateString)

	var profiles []string
	if profile == "all" {
		profiles = findProfile()
		log("info", "main", "Profiles: "+fmt.Sprint(profiles))

	} else {
		profiles = append(profiles, profile)
	}

	var browsers []string
	if browserArg == "all" {
		browsers = findInstalledBrowser()
	} else {
		browsers = append(browsers, browserArg)
	}

	log("info", "main", "Browsers: "+fmt.Sprint(browsers))

	var artifacts []BrowserArtifact

	// Loop through profiles
	for _, profile := range profiles {
		for _, browser := range browsers {
			switch browser {
			case "chrome":
				log("debug", "main", "Processing chrome artifacts for profile: "+profile)
				artifacts = append(artifacts, GetChromeArtifacts(profile, OsName)...)
			case "firefox":
				log("debug", "main", "Processing firefox artifacts for profile: "+profile)
				artifacts = append(artifacts, GetFirefoxArtifacts(profile, OsName)...)
			}
		}
	}

	log("info", "main", "Total Artifacts: "+fmt.Sprint(len(artifacts)))
	filteredArtifacts := FilterArtifacts(artifacts, startDate, endDate)
	log("info", "main", "Filtered Artifacts: "+fmt.Sprint(len(filteredArtifacts)))

	fmt.Println("Sorting artifacts...")
	//Sorting Artifacts by timestamp
	sort.Slice(filteredArtifacts[:], func(i, j int) bool {
		return filteredArtifacts[i].Timestamp < filteredArtifacts[j].Timestamp
	})

	// Join output directory and file base name
	outputFile := outputDirectory + "/" + fileBaseName
	log("info", "main", "Exporting to "+outputFile+" in "+outputFormat+" format")
	// Export artifacts to JSON
	switch outputFormat {
	case "json":
		ExportJSON(outputFile, filteredArtifacts)
	case "json_line":
		ExportJSONLine(outputFile, filteredArtifacts)
	case "csv":
		ExportCSV(outputFile, filteredArtifacts)
	}
	log("info", "main", "Export completed")
}
