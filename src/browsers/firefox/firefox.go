package firefox

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pierrec/lz4"
	"io"
	. "local/BrowserArtifact/src"
	"os"
)

const chunkSize = 256 * 1024

func log(level string, source string, message string) {
	Log.Log(level, "firefox", source, message)
}

func getBasePath(profile string, osName string) (string, string) {
	switch osName {
	case "windows":
		return fmt.Sprintf("C:\\Users\\%s\\AppData\\Roaming\\Mozilla\\firefox\\Profiles\\", profile), fmt.Sprintf("C:\\Users\\%s\\AppData\\Local\\Mozilla\\firefox\\Profiles\\", profile)
	case "darwin":
		return fmt.Sprintf("/Users/%s/Library/Application Support/firefox/Profiles/", profile), ""
	case "linux":
		return fmt.Sprintf("/home/%s/.mozilla/firefox/", profile), ""
	default:
		return "", ""
	}
}

func getFirefoxProfile(basePath string) []string {
	var profiles []string
	dir, err := os.ReadDir(basePath)
	if err != nil {
		return profiles
	}
	for _, entry := range dir {
		if entry.IsDir() {
			profiles = append(profiles, entry.Name())
		}
	}
	return profiles
}

func GetFirefoxArtifacts(profile string, osName string) []BrowserArtifact {
	var artifacts []BrowserArtifact
	basePath, localBasePath := getBasePath(profile, osName)
	firefoxProfiles := getFirefoxProfile(basePath)

	for _, firefoxProfile := range firefoxProfiles {
		artifacts = append(artifacts, processHistory(basePath+firefoxProfile+"/places.sqlite")...)
		artifacts = append(artifacts, processDownloads(basePath+firefoxProfile+"/places.sqlite")...)
		artifacts = append(artifacts, processBookmarks(basePath+firefoxProfile+"/places.sqlite")...)
		artifacts = append(artifacts, processFormHistory(basePath+firefoxProfile+"/formhistory.sqlite")...)
		artifacts = append(artifacts, processCookies(basePath+firefoxProfile+"/cookies.sqlite")...)

		if localBasePath != "" {
			artifacts = append(artifacts, processCache(localBasePath+firefoxProfile+"/cache2/")...)
		} else {
			artifacts = append(artifacts, processCache(basePath+firefoxProfile+"/cache2/")...)
		}

		artifacts = append(artifacts, processFavicons(basePath+firefoxProfile+"/favicons.sqlite")...)
		artifacts = append(artifacts, processLogins(basePath+firefoxProfile+"/logins.json")...)
		artifacts = append(artifacts, processAddons(basePath+firefoxProfile+"/addons.json")...)
		artifacts = append(artifacts, processExtensions(basePath+firefoxProfile+"/extensions.json")...)
		artifacts = append(artifacts, processBookmarksBackup(basePath+firefoxProfile+"/bookmarkbackups")...)
	}

	for i, artifact := range artifacts {
		artifact.User = profile
		artifact.App = "firefox"
		artifacts[i] = artifact
	}

	return artifacts
}

func processHistory(path string) []BrowserArtifact {
	if CheckPath(path, false) == false {
		log("error", "history", "File not found: "+path)
		return nil
	}
	var places []BrowserArtifact

	// Open the database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log("error", "history", "Error opening database: "+err.Error())
		return nil
	}
	defer db.Close()

	query := `SELECT history.visit_date, history.visit_type, place.url, ifnull(referrer.url,"") as referrer, ifnull(place.title,"") as title, place.visit_count, place.typed
FROM moz_historyvisits as history
JOIN moz_places as place ON place.id = history.place_id
LEFT JOIN moz_historyvisits as referrer_history ON referrer_history.id = history.from_visit
LEFT JOIN moz_places as referrer ON referrer.id = referrer_history.place_id;`

	rows, err := db.Query(query)
	if err != nil {
		log("error", "history", "Error querying database: "+err.Error())
		return nil
	}

	type rowStruct struct {
		visit_date  int
		visit_type  int
		url         string
		referrer    string
		title       string
		visit_count int
		typed       int
	}

	for rows.Next() {
		row := rowStruct{}

		err = rows.Scan(&row.visit_date, &row.visit_type, &row.url, &row.referrer, &row.title, &row.visit_count, &row.typed)

		if err != nil {
			log("error", "history", "Error scanning row: "+err.Error())
			continue
		}

		artifact := BrowserArtifact{}
		artifact.ArtifactType = "history"
		artifact.Url = row.url
		artifact.Title = row.title
		artifact.VisitCount = row.visit_count
		artifact.Typed = row.typed
		artifact.TimestampType = "visit_date"
		artifact.Timestamp = row.visit_date
		artifact.HttpReferrer = row.referrer
		places = append(places, artifact)
	}
	log("info", "history", fmt.Sprintf("Found %d history entries in %s", len(places), path))
	return places
}

func processBookmarks(path string) []BrowserArtifact {
	if CheckPath(path, false) == false {
		log("error", "bookmarks", "File not found: "+path)
		return nil
	}
	var places []BrowserArtifact

	// Open the database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log("error", "bookmarks", "Error opening database: "+err.Error())
		return nil
	}
	defer db.Close()

	query := "SELECT bookmark.title as bookmark_title, bookmark.dateAdded, bookmark.lastModified, ifnull(place.url,\"\") as url, ifnull(place.title,\"\") as title, ifnull(place.visit_count,0)\nfrom moz_bookmarks as bookmark\nLEFT JOIN moz_places as place ON bookmark.fk = place.id;"
	type rowStruct struct {
		bookmark_title string
		dateAdded      int
		lastModified   int
		url            string
		title          string
		visit_count    int
	}

	rows, err := db.Query(query)
	if err != nil {
		log("error", "bookmarks", "Error querying database: "+err.Error())
		return nil
	}
	for rows.Next() {
		row := rowStruct{}
		err = rows.Scan(&row.bookmark_title, &row.dateAdded, &row.lastModified, &row.url, &row.title, &row.visit_count)
		if err != nil {
			log("error", "bookmarks", "Error scanning row: "+err.Error())
			continue
		}

		artifact := BrowserArtifact{}
		artifact.ArtifactType = "bookmark"
		artifact.Title = row.title
		artifact.Url = row.url
		artifact.VisitCount = row.visit_count
		artifact.BookmarkTitle = row.bookmark_title
		artifact.TimestampType = "dateAdded"
		artifact.Timestamp = row.dateAdded
		places = append(places, artifact)

		artifact = BrowserArtifact{}
		artifact.ArtifactType = "bookmark"
		artifact.Title = row.title
		artifact.Url = row.url
		artifact.VisitCount = row.visit_count
		artifact.BookmarkTitle = row.bookmark_title
		artifact.TimestampType = "lastModified"
		artifact.Timestamp = row.lastModified
		places = append(places, artifact)

	}

	log("info", "bookmarks", fmt.Sprintf("Found %d bookmarks in %s", len(places), path))
	return places
}

func processDownloads(path string) []BrowserArtifact {
	if CheckPath(path, false) == false {
		log("error", "downloads", "File not found: "+path)
		return nil
	}

	// Open the database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log("error", "downloads", "Error opening database: "+err.Error())
		return nil
	}
	defer db.Close()

	query := "SELECT t1.place_id, t1.content as file, t2.content as metadata, t1.flags, t1.expiration, t1.dateAdded, t1. lastModified, place.url, ifnull(place.title,\"\"), ifnull(place.description,\"\"), place.visit_count\nFROM moz_annos as t1\nLEFT JOIN moz_annos as t2 ON t1.place_id = t2.place_id\nLEFT JOIN moz_places as place ON t1.place_id = place.id\nWHERE t1.anno_attribute_id = 1;"
	type rowStruct struct {
		place_id     int
		file         string
		metadata     string
		flags        int
		expiration   int
		dateAdded    int
		lastModified int
		url          string
		title        string
		description  string
		visit_count  int
	}
	var downloads []BrowserArtifact

	rows, err := db.Query(query)
	if err != nil {
		log("error", "downloads", "Error querying database: "+err.Error())
		return nil
	}
	for rows.Next() {
		row := rowStruct{}
		err = rows.Scan(&row.place_id, &row.file, &row.metadata, &row.flags, &row.expiration, &row.dateAdded, &row.lastModified, &row.url, &row.title, &row.description, &row.visit_count)
		if err != nil {
			log("error", "downloads", "Error scanning row: "+err.Error())
			continue
		}

		artifact := BrowserArtifact{}
		artifact.ArtifactType = "download"
		artifact.Url = row.url
		artifact.Title = row.title
		artifact.TimestampType = "dateAdded"
		artifact.Timestamp = row.dateAdded
		artifact.VisitCount = row.visit_count
		artifact.Metadata = row.metadata
		artifact.Filename = row.file
		downloads = append(downloads, artifact)

		artifact = BrowserArtifact{}
		artifact.ArtifactType = "download"
		artifact.Url = row.url
		artifact.Title = row.title
		artifact.TimestampType = "lastModified"
		artifact.Timestamp = row.lastModified
		artifact.VisitCount = row.visit_count
		artifact.Metadata = row.metadata
		artifact.Filename = row.file
		downloads = append(downloads, artifact)
	}

	log("info", "downloads", fmt.Sprintf("Found %d downloads in %s", len(downloads), path))
	return downloads
}

func processFormHistory(path string) []BrowserArtifact {
	if CheckPath(path, false) == false {
		log("error", "formhistory", "File not found: "+path)
		return nil
	}
	var formHistory []BrowserArtifact

	// Open the database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log("error", "formhistory", "Error opening database: "+err.Error())
		return nil
	}
	defer db.Close()

	query := "SELECT fieldname, value, firstUsed, lastUsed, source FROM moz_formhistory as formhistory\nJOIN moz_history_to_sources as history2source ON formhistory.id = history2source.history_id\nJOIN moz_sources as source ON history2source.source_id = source.id;"

	type rowStruct struct {
		fieldname string
		value     string
		firstUsed int
		lastUsed  int
		source    string
	}

	rows, err := db.Query(query)
	if err != nil {
		log("error", "formhistory", "Error querying database: "+err.Error())
		return nil
	}

	for rows.Next() {
		row := rowStruct{}
		err = rows.Scan(&row.fieldname, &row.value, &row.firstUsed, &row.lastUsed, &row.source)
		if err != nil {
			log("error", "formhistory", "Error scanning row: "+err.Error())
			continue
		}

		artifact := BrowserArtifact{}
		artifact.ArtifactType = "formhistory"
		artifact.Url = row.source
		artifact.Fieldname = row.fieldname
		artifact.Value = row.value
		artifact.TimestampType = "firstUsed"
		artifact.Timestamp = row.firstUsed
		formHistory = append(formHistory, artifact)

		artifact = BrowserArtifact{}
		artifact.ArtifactType = "formhistory"
		artifact.Url = row.source
		artifact.Fieldname = row.fieldname
		artifact.Value = row.value
		artifact.TimestampType = "lastUsed"
		artifact.Timestamp = row.lastUsed
		formHistory = append(formHistory, artifact)

	}

	log("info", "formhistory", fmt.Sprintf("Found %d form history entries in %s", len(formHistory), path))
	return formHistory
}

func processCookies(path string) []BrowserArtifact {
	if CheckPath(path, false) == false {
		log("error", "cookies", "File not found: "+path)
		return nil
	}
	var cookies []BrowserArtifact

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log("error", "cookies", "Error opening database: "+err.Error())
		return nil
	}
	defer db.Close()

	query := "SELECT host, name, value, path, expiry, lastAccessed, creationTime FROM moz_cookies;"
	type rowStruct struct {
		host         string
		name         string
		value        string
		path         string
		expiry       int
		lastAccessed int
		creationTime int
	}

	rows, err := db.Query(query)
	if err != nil {
		log("error", "cookies", "Error querying database: "+err.Error())
		return nil
	}
	for rows.Next() {
		row := rowStruct{}
		err = rows.Scan(&row.host, &row.name, &row.value, &row.path, &row.expiry, &row.lastAccessed, &row.creationTime)
		if err != nil {
			log("error", "cookies", "Error scanning row: "+err.Error())
			continue
		}

		artifact := BrowserArtifact{}
		artifact.ArtifactType = "cookie"
		artifact.Url = row.host
		artifact.Cookie = row.name + "=" + row.value
		artifact.TimestampType = "creationTime"
		artifact.Timestamp = row.creationTime
		cookies = append(cookies, artifact)

		artifact = BrowserArtifact{}
		artifact.ArtifactType = "cookie"
		artifact.Url = row.host
		artifact.Cookie = row.name + "=" + row.value
		artifact.TimestampType = "lastAccessed"
		artifact.Timestamp = row.lastAccessed
		cookies = append(cookies, artifact)
	}

	log("info", "cookies", fmt.Sprintf("Found %d cookies in %s", len(cookies), path))
	return cookies
}

/*
Credits to https://github.com/JamesHabben/FirefoxCache2/blob/master/firefox-cache2-file-parser.py
This function is based on the above code
Thanks to James Habben
*/
func parseCacheFile(filename string) (error, []BrowserArtifact) {
	artifacts := make([]BrowserArtifact, 0)

	file, err := os.Open(filename)
	if err != nil {
		return err, nil
	}
	defer file.Close()

	_, err = file.Seek(-4, io.SeekEnd)
	if err != nil {
		return err, nil
	}

	var metaStart uint32
	err = binary.Read(file, binary.BigEndian, &metaStart)
	if err != nil {
		return err, nil
	}

	numHashChunks := metaStart / chunkSize
	if metaStart%chunkSize != 0 {
		numHashChunks++
	}

	_, err = file.Seek(int64(metaStart+4+numHashChunks*2), io.SeekStart)
	if err != nil {
		return err, nil
	}

	var version, fetchCount, lastFetchInt, lastModInt, frecency, expireInt, keySize uint32

	err = binary.Read(file, binary.BigEndian, &version)
	if err != nil {
		return err, nil
	}

	err = binary.Read(file, binary.BigEndian, &fetchCount)
	if err != nil {
		return err, nil
	}

	err = binary.Read(file, binary.BigEndian, &lastFetchInt)
	if err != nil {
		return err, nil
	}

	err = binary.Read(file, binary.BigEndian, &lastModInt)
	if err != nil {
		return err, nil
	}

	err = binary.Read(file, binary.BigEndian, &frecency)
	if err != nil {
		return err, nil
	}

	err = binary.Read(file, binary.BigEndian, &expireInt)
	if err != nil {
		return err, nil
	}

	err = binary.Read(file, binary.BigEndian, &keySize)
	if err != nil {
		return err, nil
	}

	key := make([]byte, keySize)
	_, err = file.Read(key)
	if err != nil {
		return err, nil
	}

	lastFetchArtifact := BrowserArtifact{}
	lastFetchArtifact.ArtifactType = "cache"
	lastFetchArtifact.Url = string(key)
	// Convert sec to microsec
	lastFetchArtifact.Timestamp = int(lastFetchInt) * 1000 * 1000
	lastFetchArtifact.VisitCount = int(fetchCount)

	if lastFetchInt == lastModInt {
		lastFetchArtifact.TimestampType = "lastFetch (lastMod)"
	} else {
		lastModArtifact := BrowserArtifact{}
		lastModArtifact.ArtifactType = "cache"
		lastModArtifact.Url = string(key)
		lastModArtifact.TimestampType = "lastMod"
		lastModArtifact.Timestamp = int(lastModInt)
		lastModArtifact.VisitCount = int(fetchCount)
		artifacts = append(artifacts, lastModArtifact)
	}
	artifacts = append(artifacts, lastFetchArtifact)

	return nil, artifacts
}

func processCache(path string) []BrowserArtifact {
	if CheckPath(path, true) == false {
		log("error", "cache", "Directory not found: "+path)
		return nil
	}
	var cache []BrowserArtifact

	// List all files in the cache directory
	dir, err := os.ReadDir(path + "entries")
	if err != nil {
		log("error", "cache", "Error reading cache directory: "+err.Error())
		return nil
	}

	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}

		err, artifacts := parseCacheFile(path + "entries/" + entry.Name())
		if err != nil {
			log("error", "cache", "Error parsing cache file: "+err.Error())
			continue
		}
		cache = append(cache, artifacts...)
	}

	log("info", "cache", fmt.Sprintf("Found %d cache entries in %s", len(cache), path))
	return cache
}

func processFavicons(path string) []BrowserArtifact {
	if CheckPath(path, false) == false {
		log("error", "favicons", "File not found: "+path)
		return nil
	}
	var favicons []BrowserArtifact

	// Open the database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log("error", "favicons", "Error opening database: "+err.Error())
		return nil
	}
	defer db.Close()

	query := "SELECT icon_url, expire_ms FROM moz_icons;"
	type rowStruct struct {
		icon_url string
		expires  int
	}

	rows, err := db.Query(query)
	if err != nil {
		log("error", "favicons", "Error querying database: "+err.Error())
		return nil
	}

	for rows.Next() {
		row := rowStruct{}
		err = rows.Scan(&row.icon_url, &row.expires)
		if err != nil {
			log("error", "favicons", "Error scanning row: "+err.Error())
			continue
		}

		artifact := BrowserArtifact{}
		artifact.ArtifactType = "favicon"
		artifact.Url = row.icon_url
		artifact.TimestampType = "expires -7 days (favicons)"
		artifact.Timestamp = row.expires - 7*24*60*60*1000
		// Convert milliseconds to microseconds
		artifact.Timestamp = artifact.Timestamp * 1000
		favicons = append(favicons, artifact)
	}

	log("info", "favicons", fmt.Sprintf("Found %d favicons in %s", len(favicons), path))
	return favicons
}

func processLogins(path string) []BrowserArtifact {
	if CheckPath(path, false) == false {
		log("error", "logins", "File not found: "+path)
		return nil
	}
	var logins []BrowserArtifact

	// Open JSON file
	file, err := os.Open(path)
	if err != nil {
		log("error", "logins", "Error opening file: "+err.Error())
		return nil
	}
	defer file.Close()

	// Parse JSON file
	jsonParser := json.NewDecoder(file)
	var loginJSON loginJSON
	err = jsonParser.Decode(&loginJSON)
	if err != nil {
		log("error", "logins", "Error parsing JSON: "+err.Error())
		return nil
	}

	for _, login := range loginJSON.Logins {
		artifact := BrowserArtifact{}
		artifact.Url = login.Hostname
		artifact.ArtifactType = "login"
		artifact.TimestampType = "dateCreated"
		artifact.Timestamp = int(login.TimeCreated)
		// Convert milliseconds to microseconds
		artifact.Timestamp = artifact.Timestamp * 1000
		logins = append(logins, artifact)

		artifact = BrowserArtifact{}
		artifact.Url = login.Hostname
		artifact.ArtifactType = "login"
		artifact.TimestampType = "dateLastUsed"
		artifact.Timestamp = int(login.TimeLastUsed)
		// Convert milliseconds to microseconds
		artifact.Timestamp = artifact.Timestamp * 1000
		logins = append(logins, artifact)

		artifact = BrowserArtifact{}
		artifact.ArtifactType = "login"
		artifact.Url = login.Hostname
		artifact.TimestampType = "datePasswordChanged"
		artifact.Timestamp = int(login.TimePasswordChanged)
		// Convert milliseconds to microseconds
		artifact.Timestamp = artifact.Timestamp * 1000
		logins = append(logins, artifact)
	}

	log("info", "logins", fmt.Sprintf("Found %d logins in %s", len(logins), path))
	return logins
}

func processAddons(path string) []BrowserArtifact {
	if CheckPath(path, false) == false {
		log("error", "addons", "File not found: "+path)
		return nil
	}
	var addons []BrowserArtifact

	// Open JSON file
	file, err := os.Open(path)
	if err != nil {
		log("error", "addons", "Error opening file: "+err.Error())
		return nil
	}
	defer file.Close()

	// Parse JSON file
	jsonParser := json.NewDecoder(file)
	var addonJSON addonJSON
	err = jsonParser.Decode(&addonJSON)
	if err != nil {
		log("error", "addons", "Error parsing JSON: "+err.Error())
		return nil
	}

	for _, addon := range addonJSON.Addons {
		artifact := BrowserArtifact{}
		artifact.ArtifactType = "addon"
		artifact.AddonName = addon.Name
		artifact.AddonType = addon.Type
		artifact.Description = addon.Description
		artifact.FullDescription = addon.FullDescription
		artifact.Version = addon.Version
		artifact.SourceURI = addon.SourceURI
		artifact.HomePageURL = addon.HomepageURL
		artifact.ReviewURL = addon.AmoListingURL
		artifact.CreatorName = addon.Creator.Name
		artifact.CreatorURL = addon.Creator.URL
		artifact.AverageRating = float32(addon.AverageRating)
		artifact.RatingCount = addon.ReviewCount
		artifact.TimestampType = "updateDate"
		artifact.Timestamp = int(addon.UpdateDate)
		// Convert milliseconds to microseconds
		artifact.Timestamp = artifact.Timestamp * 1000
		addons = append(addons, artifact)
	}

	log("info", "addons", fmt.Sprintf("Found %d addons in %s", len(addons), path))
	return addons
}

func processExtensions(path string) []BrowserArtifact {
	if CheckPath(path, false) == false {
		log("error", "extensions", "File not found: "+path)
		return nil
	}
	var extensions []BrowserArtifact

	// Open JSON file
	file, err := os.Open(path)
	if err != nil {
		log("error", "extensions", "Error opening file: "+err.Error())
		return nil
	}
	defer file.Close()

	// Parse JSON file
	jsonParser := json.NewDecoder(file)
	var extensionJSON extensionJSON
	err = jsonParser.Decode(&extensionJSON)
	if err != nil {
		log("error", "extensions", "Error parsing JSON: "+err.Error())
		return nil
	}

	for _, addon := range extensionJSON.Addons {
		artifact := BrowserArtifact{}
		artifact.ArtifactType = "extension"
		artifact.AddonName = addon.DefaultLocale.Name
		artifact.Version = addon.Version
		artifact.AddonType = addon.Type
		if tmp, ok := addon.AboutURL.(string); ok {
			artifact.AboutURL = tmp
		}
		artifact.Description = addon.DefaultLocale.Description
		artifact.CreatorName = addon.DefaultLocale.Creator
		artifact.HomePageURL = addon.DefaultLocale.HomepageURL
		artifact.Active = addon.Active
		artifact.Visible = addon.Visible
		artifact.SourceURI = addon.SourceURI
		artifact.Url = addon.RootURI
		artifact.TimestampType = "updateDate"
		// Convert milliseconds to nanoseconds
		artifact.Timestamp = int(addon.UpdateDate) * 1000
		extensions = append(extensions, artifact)

		artifact = BrowserArtifact{}
		artifact.ArtifactType = "extension"
		artifact.AddonName = addon.DefaultLocale.Name
		artifact.Version = addon.Version
		artifact.AddonType = addon.Type
		if tmp, ok := addon.AboutURL.(string); ok {
			artifact.AboutURL = tmp
		}
		artifact.Description = addon.DefaultLocale.Description
		artifact.CreatorName = addon.DefaultLocale.Creator
		artifact.HomePageURL = addon.DefaultLocale.HomepageURL
		artifact.Active = addon.Active
		artifact.Visible = addon.Visible
		artifact.SourceURI = addon.SourceURI
		artifact.Url = addon.RootURI
		artifact.TimestampType = "installDate"
		artifact.Timestamp = int(addon.InstallDate)
	}

	log("info", "extensions", fmt.Sprintf("Found %d extensions in %s", len(extensions), path))
	return extensions
}

func parseBookmarkBackupFile(path string) (error, []BrowserArtifact) {

	artifacts := []BrowserArtifact{}

	file, err := os.Open(path)
	if err != nil {
		return err, nil
	}
	defer file.Close()

	// Read the LZ4 header (8 bytes) to verify the format
	header := make([]byte, 8)
	_, err = io.ReadFull(file, header)
	if err != nil {
		return err, nil
	}

	// Verify the header format (expect "mozLz40\x00")
	expectedHeader := []byte("mozLz40\x00")
	if !bytes.Equal(header[:8], expectedHeader) {
		return err, nil
	}

	// Read the uncompressed size (4 bytes, little-endian)
	var size uint32
	err = binary.Read(file, binary.LittleEndian, &size)
	if err != nil {
		return err, nil
	}

	// Create a buffer to store the LZ4 block
	lz4Block := make([]byte, size)
	_, err = io.ReadFull(file, lz4Block)
	if err != nil && err != io.ErrUnexpectedEOF {
		return err, nil
	}

	// Decompress the LZ4 block
	decompressedData := make([]byte, size)
	_, err = lz4.UncompressBlock(lz4Block, decompressedData)
	if err != nil {
		return err, nil
	}

	return nil, artifacts
}

func processBookmarksBackup(path string) []BrowserArtifact {
	log("debug", "bookmarks", "Not implemented yet")
	// Check if the profile has bookmarkbackups
	if CheckPath(path, true) == false {
		log("error", "bookmarks", "Directory not found: "+path)
		return nil
	}

	var bookmarks []BrowserArtifact
	/*
		// List all files in the bookmarkbackups directory
		dir, err := os.ReadDir(path)
		if err != nil {
			return nil
		}

		// Loop through all files in the bookmarkbackups directory
		for _, entry := range dir {
			if entry.IsDir() {
				continue
			}

			// Files are in JSONLZ4 format
			// Open the file
			err, artifacts := parseBookmarkBackupFile(path + "/" + entry.Name())
			if err != nil {
				continue
			}
			bookmarks = append(bookmarks, artifacts...)

		}
	*/
	return bookmarks
}
