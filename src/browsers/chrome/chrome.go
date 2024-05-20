package chrome

import (
	"database/sql"
	"encoding/json"
	"fmt"
	. "local/BrowserArtifact/src"
	"os"
	"strconv"
	"time"
)

func log(level string, source string, message string) {
	Log.Log(level, "chrome", source, message)
}

func getBasePath(profile string, osName string) string {
	switch osName {
	case "windows":
		return fmt.Sprintf("C:\\Users\\%s\\AppData\\Local\\Google\\Chrome\\User Data\\Default\\", profile)
	case "darwin":
		return fmt.Sprintf("/Users/%s/Library/Application Support/Google/Chrome/")
	case "linux":
		return fmt.Sprintf("/home/%s/.config/google-chrome/")
	default:
		return ""
	}
}

func GetChromeArtifacts(profile string, osName string) []BrowserArtifact {
	artifacts := []BrowserArtifact{}

	basePath := getBasePath(profile, osName)

	artifacts = append(artifacts, processHistory(basePath+"History")...)
	artifacts = append(artifacts, processDownloads(basePath+"History")...)
	artifacts = append(artifacts, processBookmarks(basePath+"Bookmarks")...)
	artifacts = append(artifacts, processCookies(basePath+"Network/Cookies")...)
	artifacts = append(artifacts, processFormHistory(basePath+"Web Data")...)
	artifacts = append(artifacts, processLoginData(basePath+"Login Data")...)
	artifacts = append(artifacts, processExtensions(basePath+"Extensions")...)
	artifacts = append(artifacts, processFavicons(basePath+"Favicons")...)
	//artifacts = append(artifacts, processSession(basePath+"Session")...)
	//artifacts = append(artifacts, processThumbnail(basePath+"Thumbnail")...)
	artifacts = append(artifacts, processCache(basePath+"Cache")...)

	for i, artifact := range artifacts {
		artifact.User = profile
		artifact.App = "chrome"
		artifacts[i] = artifact
	}

	return artifacts

}

func processHistory(path string) []BrowserArtifact {
	// Check if file exists
	if !CheckPath(path, false) {
		log("error", "history", "File not found : "+path)
		return []BrowserArtifact{}
	}

	artifacts := []BrowserArtifact{}

	// Open the database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log("error", "history", "Error opening database: "+err.Error())
	}
	defer db.Close()

	//This timestamp format is used in web browsers such as Apple Safari (WebKit), Google Chrome and Opera (Chromium/Blink).
	//It's a 64-bit value for microseconds since Jan 1, 1601 00:00 UTC. One microsecond is one-millionth of a second.
	//11644474161000000 is the number of microseconds between Jan 1, 1601 and Jan 1, 1970.
	query := "SELECT history.visit_time - 11644474161000000, history.visit_duration, url.url, url.title, url.visit_count, url.typed_count, ifnull(referrer_url.url,\"\") as referrer, ifnull(referrer_url.title,\"\") as referrer_title  FROM visits as history\nLEFT JOIN urls as url ON url.id = history.url\nLEFT JOIN visits as referrer_history ON referrer_history.id = history.opener_visit OR referrer_history.id = history.from_visit\nLEFT JOIN urls as referrer_url ON referrer_history.url = referrer_url.id;"

	rows, err := db.Query(query)
	if err != nil {
		log("error", "history", "Error querying database: "+err.Error())
		return nil
	}

	type rowStruct struct {
		visit_time     int
		visit_duration int
		url            string
		title          string
		visit_count    int
		typed_count    int
		referrer       string
		referrer_title string
	}

	for rows.Next() {
		var row rowStruct
		err = rows.Scan(&row.visit_time, &row.visit_duration, &row.url, &row.title, &row.visit_count, &row.typed_count, &row.referrer, &row.referrer_title)
		if err != nil {
			log("error", "history", "Error scanning row: "+err.Error())
			return nil
		}

		artifact := BrowserArtifact{}
		artifact.ArtifactType = "chrome_history"
		artifact.Timestamp = row.visit_time
		artifact.TimestampType = "visit_date"
		artifact.Title = row.title
		artifact.Url = row.url
		artifact.HttpReferrer = row.referrer
		artifact.VisitCount = row.visit_count
		artifact.Typed = row.typed_count
		artifact.Duration = row.visit_duration
		artifacts = append(artifacts, artifact)
	}

	log("info", "history", fmt.Sprintf("Found %d history entries in %s", len(artifacts), path))
	return artifacts
}

func processDownloads(path string) []BrowserArtifact {
	if !CheckPath(path, false) {
		log("error", "downloads", "File not found : "+path)
		return nil
	}
	artifacts := []BrowserArtifact{}

	// Open the database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log("error", "downloads", "Error opening database: "+err.Error())
		return nil
	}
	defer db.Close()

	query := "SELECT start_time - 11644474161000000 as start_time, target_path,  received_bytes, total_bytes, end_time - 11644474161000000 as end_time, tab_url, tab_referrer_url, mime_type FROM downloads;"

	rows, err := db.Query(query)
	if err != nil {
		log("error", "downloads", "Error querying database: "+err.Error())
		return nil
	}

	type rowStruct struct {
		start_time       int
		target_path      string
		received_bytes   int
		total_bytes      int
		end_time         int
		tab_url          string
		tab_referrer_url string
		mime_type        string
	}

	for rows.Next() {
		var row rowStruct
		err = rows.Scan(&row.start_time, &row.target_path, &row.received_bytes, &row.total_bytes, &row.end_time, &row.tab_url, &row.tab_referrer_url, &row.mime_type)
		if err != nil {
			log("error", "downloads", "Error scanning row: "+err.Error())
			return nil
		}

		artifact := BrowserArtifact{}
		artifact.ArtifactType = "download"
		artifact.Timestamp = row.start_time
		artifact.Duration = row.end_time - row.start_time
		artifact.TimestampType = "dateAdded"
		artifact.Url = row.tab_url
		artifact.HttpReferrer = row.tab_referrer_url
		artifact.Filename = row.target_path
		artifact.MimeType = row.mime_type
		artifact.BytesIn = row.received_bytes
		artifact.BytesOut = row.total_bytes - row.received_bytes
		artifacts = append(artifacts, artifact)
	}

	return artifacts
}

func processBookmarks(path string) []BrowserArtifact {
	// Check if file exists
	if !CheckPath(path, false) {
		log("error", "bookmarks", "File not found : "+path)
		return nil
	}

	artifacts := []BrowserArtifact{}

	// Open JSON file
	file, err := os.Open(path)
	if err != nil {
		log("error", "bookmarks", "Error opening file: "+err.Error())
		return nil
	}
	defer file.Close()

	// Parse JSON file
	decoder := json.NewDecoder(file)
	var data bookmark
	err = decoder.Decode(&data)
	if err != nil {
		log("error", "bookmarks", "Error decoding JSON: "+err.Error())
		return nil
	}

	// Process bookmarks
	for _, bookmark := range data.Roots.BookmarkBar.Children {
		if bookmark.Type != "url" {
			continue
		}

		artifact := BrowserArtifact{}
		artifact.ArtifactType = "bookmark"
		dateAdded, _ := strconv.Atoi(bookmark.DateAdded)
		artifact.Timestamp = dateAdded - 11644474161000000
		artifact.TimestampType = "dateAdded"
		artifact.BookmarkTitle = bookmark.Name
		artifact.Url = bookmark.URL
		artifacts = append(artifacts, artifact)

		artifact = BrowserArtifact{}
		artifact.ArtifactType = "bookmark"
		dateAdded, _ = strconv.Atoi(bookmark.DateLastUsed)
		artifact.Timestamp = dateAdded - 11644474161000000
		artifact.TimestampType = "dateLastUsed"
		artifact.BookmarkTitle = bookmark.Name
		artifact.Url = bookmark.URL
		artifacts = append(artifacts, artifact)

	}

	return artifacts
}

func processCookies(path string) []BrowserArtifact {
	artifacts := []BrowserArtifact{}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log("error", "cookies", "Error opening database: "+err.Error())
		return nil
	}
	defer db.Close()

	query := "SELECT creation_utc - 11644474161000000, last_access_utc - 11644474161000000, last_update_utc - 11644474161000000, host_key, source_port, name, value FROM cookies;"
	type rowStruct struct {
		creationTime int
		lastAccessed int
		lastUpdate   int
		host         string
		sourcePort   int
		name         string
		value        string
	}

	rows, err := db.Query(query)
	if err != nil {
		log("error", "cookies", "Error querying database: "+err.Error())
		return nil
	}
	for rows.Next() {
		row := rowStruct{}
		err = rows.Scan(&row.creationTime, &row.lastAccessed, &row.lastUpdate, &row.host, &row.sourcePort, &row.name, &row.value)
		if err != nil {
			log("error", "cookies", "Error scanning row: "+err.Error())
			return nil
		}

		artifact := BrowserArtifact{}
		artifact.ArtifactType = "cookie"
		artifact.Url = row.host
		artifact.Cookie = row.name + "=" + row.value
		artifact.TimestampType = "creationTime"
		artifact.Timestamp = row.creationTime
		artifact.DestPort = row.sourcePort
		artifacts = append(artifacts, artifact)

		artifact = BrowserArtifact{}
		artifact.ArtifactType = "cookie"
		artifact.Url = row.host
		artifact.Cookie = row.name + "=" + row.value
		artifact.TimestampType = "lastAccessed"
		artifact.Timestamp = row.lastAccessed
		artifact.DestPort = row.sourcePort
		artifacts = append(artifacts, artifact)

		artifact = BrowserArtifact{}
		artifact.ArtifactType = "cookie"
		artifact.Url = row.host
		artifact.Cookie = row.name + "=" + row.value
		artifact.TimestampType = "lastUpdate"
		artifact.Timestamp = row.lastUpdate
		artifact.DestPort = row.sourcePort
		artifacts = append(artifacts, artifact)
	}

	log("info", "cookies", fmt.Sprintf("Found %d cookies in %s", len(artifacts), path))
	return artifacts
}

func processFormHistory(path string) []BrowserArtifact {
	log("debug", "formhistory", "Not implemented yet")

	artifacts := []BrowserArtifact{}

	return artifacts
}

func processLoginData(path string) []BrowserArtifact {
	// Check if file exists
	if !CheckPath(path, false) {
		log("error", "login", "File not found : "+path)
		return nil
	}

	artifacts := []BrowserArtifact{}

	// Open the database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log("error", "login", "Error opening database: "+err.Error())
		return nil
	}
	defer db.Close()

	query := "SELECT date_created - 11644474161000000, date_last_used- 11644474161000000, date_password_modified - 11644474161000000, origin_url, username_value, times_used FROM logins;"
	rows, err := db.Query(query)

	type rowStruct struct {
		dateCreated         int
		dateLastUsed        int
		datePasswordChanged int
		originUrl           string
		usernameValue       string
		timesUsed           int
	}

	for rows.Next() {
		var row rowStruct
		err = rows.Scan(&row.dateCreated, &row.dateLastUsed, &row.datePasswordChanged, &row.originUrl, &row.usernameValue, &row.timesUsed)
		if err != nil {
			log("error", "login", "Error scanning row: "+err.Error())
			return nil
		}

		artifact := BrowserArtifact{}
		artifact.ArtifactType = "login"
		artifact.Timestamp = row.dateCreated
		artifact.TimestampType = "dateCreated"
		artifact.Url = row.originUrl
		artifact.Fieldname = "username"
		artifact.Value = row.usernameValue
		artifacts = append(artifacts, artifact)

		artifact = BrowserArtifact{}
		artifact.ArtifactType = "login"
		artifact.Timestamp = row.dateLastUsed
		artifact.TimestampType = "lastUsed"
		artifact.Url = row.originUrl
		artifact.Fieldname = "username"
		artifact.Value = row.usernameValue
		artifacts = append(artifacts, artifact)

		artifact = BrowserArtifact{}
		artifact.ArtifactType = "login"
		artifact.Timestamp = row.datePasswordChanged
		artifact.TimestampType = "passwordChanged"
		artifact.Url = row.originUrl
		artifact.Fieldname = "username"
		artifact.Value = row.usernameValue
		artifacts = append(artifacts, artifact)

	}
	log("info", "login", fmt.Sprintf("Found %d logins in %s", len(artifacts), path))
	return artifacts
}

func processExtensions(path string) []BrowserArtifact {
	// Check if directory exists
	if !CheckPath(path, true) {
		log("error", "extensions", "Directory not found : "+path)
		return nil
	}
	artifacts := []BrowserArtifact{}

	// Find manifest.json files in Extensions/<name>/<version>/manifest.json
	// List all directories in Extensions
	dirs, err := os.ReadDir(path)
	if err != nil {
		log("error", "extensions", "Error reading directory: "+err.Error())
		return nil
	}

	for _, dir := range dirs {
		// Check if directory is a directory
		if !dir.IsDir() {
			continue
		}
		// List all directories in Extensions/<name>
		versions, err := os.ReadDir(path + "/" + dir.Name())
		if err != nil {
			log("error", "extensions", "Error reading directory: "+err.Error())
			return nil
		}

		for _, version := range versions {
			// Check if directory is a directory
			if !version.IsDir() {
				continue
			}
			manifestPath := path + "/" + dir.Name() + "/" + version.Name() + "/manifest.json"
			// Check if manifest.json exists

			if !CheckPath(manifestPath, false) {
				continue
			}

			// Open manifest.json
			file, err := os.Open(manifestPath)
			if err != nil {
				log("error", "extensions", "Error opening file: "+err.Error())
				continue
			}
			defer file.Close()

			// Parse JSON file
			decoder := json.NewDecoder(file)
			var data extension
			err = decoder.Decode(&data)
			if err != nil {
				log("error", "extensions", "Error decoding JSON: "+err.Error())
				continue
			}

			// Process extension
			extension := BrowserArtifact{}
			extension.ArtifactType = "extension"
			extension.AddonName = data.Name
			extension.TimestampType = "now"
			extension.Timestamp = int(time.Now().UnixMicro())
			extension.Version = data.Version
			extension.SourceURI = data.UpdateURL
			extension.CreatorName = data.Author
			extension.Description = data.Description
			artifacts = append(artifacts, extension)
		}
	}

	log("info", "extensions", fmt.Sprintf("Found %d extensions in %s", len(artifacts), path))

	return artifacts
}

func processFavicons(path string) []BrowserArtifact {
	// Check if file exists
	if !CheckPath(path, false) {
		log("error", "favicons", "File not found : "+path)
		return nil
	}

	artifacts := []BrowserArtifact{}

	// Open the database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log("error", "favicons", "Error opening database: "+err.Error())
		return nil
	}
	defer db.Close()

	query := "SELECT favicon.url as url, map.page_url as referrer, last_updated, last_requested\nFROM favicons as favicon\nLEFT JOIN favicon_bitmaps as bitmap ON favicon.id=bitmap.icon_id\nLEFT JOIN icon_mapping as map ON favicon.id=map.icon_id\n"

	rows, err := db.Query(query)
	if err != nil {
		log("error", "favicons", "Error querying database: "+err.Error())
		return nil
	}

	type rowStruct struct {
		iconUrl       string
		url           string
		lastUpdated   int
		lastRequested int
	}

	for rows.Next() {
		var row rowStruct
		err = rows.Scan(&row.iconUrl, &row.url, &row.lastUpdated, &row.lastRequested)
		if err != nil {
			log("error", "favicons", "Error scanning row: "+err.Error())
			return nil
		}

		// From observation, the last_updated and last_requested timestamps are never both set, sometimes even both are 0

		artifact := BrowserArtifact{}
		artifact.ArtifactType = "favicon"
		artifact.HttpReferrer = row.url
		artifact.Url = row.iconUrl

		if row.lastUpdated != 0 {
			artifact.Timestamp = row.lastUpdated
			artifact.TimestampType = "lastUpdated"
		} else if row.lastRequested != 0 {
			artifact.Timestamp = row.lastRequested
			artifact.TimestampType = "lastRequested"
		} else {
			artifact.Timestamp = int(time.Now().UnixMicro())
			artifact.TimestampType = "unknown"
		}

	}

	return artifacts
}

func processCache(path string) []BrowserArtifact {
	log("debug", "cache", "Not implemented yet")
	artifacts := []BrowserArtifact{}

	return artifacts
}
