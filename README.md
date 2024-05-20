# Browser Artefacts

"A tool to control them all."

## Usage

```
Usage of BrowserArtifact.exe:
  -browser string
        Browser: chrome, firefox [default: all] (default "all")
  -end_date string
        End Date [default: now] (default "now")
  -start_date string
        Start Date [default: 2000-01-01] (default "2000-01-01")
  -file_base_name string
        File Base Name [default: ''] (default "BrowserArtifacts")
  -format string
        Output Format: json, json_line, csv [default: json] (default "json")
  -log_file string
        Log File [default: '']
  -output_directory string
        Output Directory [default: .] (default ".")
  -profile string
        Profile [default: all] (default "all")
  -verbose string
        Verbose Level: debug, info, warn, error [default: info] (default "info")
```


## Handled Artefacts

### Firefox

- [x] History (SQLite): 
  - `C:\Users\XXX\AppData\Roaming\Mozilla\Firefox\Profiles\XXX\places.sqlite`
- [x] Downloads (SQLite):
  - `C:\Users\XXX\AppData\Roaming\Mozilla\Firefox\Profiles\XXX\downloads.sqlite`
- [x] Cookies (SQLite): 
  - `C:\Users\XXX\AppData\Roaming\Mozilla\Firefox\Profiles\XXX\cookies.sqlite`
- [x] Cache (Miscellaneous):
  - `C:\Users\XXX\AppData\Roaming\Mozilla\Firefox\Profiles\XXX\cache2\entries`
  - `C:\Users\XXX\AppData\Roaming\Mozilla\Firefox\Profiles\XXX\cache2\startupCache`
- [x] Form History (SQLite):
  - `C:\Users\XXX\AppData\Roaming\Mozilla\Firefox\Profiles\XXX\formhistory.sqlite`
- [x] Addons & Extensions (SQLite):
  - `C:\Users\XXX\AppData\Roaming\Mozilla\Firefox\Profiles\XXX\extensions.sqlite`
  - `C:\Users\XXX\AppData\Roaming\Mozilla\Firefox\Profiles\XXX\addons.sqlite`
- [x] Favicons (SQLite):
  - `C:\Users\XXX\AppData\Roaming\Mozilla\Firefox\Profiles\XXX\favicons.sqlite`
- [x] Logins (JSON):
  - `C:\Users\XXX\AppData\Roaming\Mozilla\Firefox\Profiles\XXX\logins.json`
- [] Passords (KDBX?):
  - `C:\Users\XXX\AppData\Roaming\Mozilla\Firefox\Profiles\XXX\key4.db`
  - `C:\Users\XXX\AppData\Roaming\Mozilla\Firefox\Profiles\XXX\key3.db`
- [] Session Data:
  - [] Current Session:
    - `C:\Users\XXX\AppData\Roaming\Mozilla\Firefox\Profiles\XXX\sessionstore.json`
  - [] Last Session:
    - `C:\Users\XXX\AppData\Roaming\Mozilla\Firefox\Profiles\XXX\sessionstore-backups\recovery.jsonlz4`
- [] Thumbnails (Folder):
  - `C:\Users\XXX\AppData\Roaming\Mozilla\Firefox\Profiles\XXX\thumbnails`
- [] Bookmarks backup (jsonlz4):
  - `C:\Users\XXX\AppData\Roaming\Mozilla\Firefox\Profiles\XXX\bookmarkbackups`

### Chrome

- [x] History (SQLite): 
  - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\Default\History`
  - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\ChromeDefaultData\History`
- [x] Cookies (SQLite): 
  - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\Default\Cookies`
  - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\ChromeDefaultData\Cookies`
- [] Cache (Miscellaneous):
  - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\Default\Cache`
  - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\ChromeDefaultData\Cache`
- [x] Bookmarks (JSON):
  - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\Default\Bookmarks`
  - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\ChromeDefaultData\Bookmarks`
- [x] Form History (SQLite):
  - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\Default\Web Data`
  - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\ChromeDefaultData\Web Data`
- [x] Favicons (SQLite):
  - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\Default\Favicons`
  - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\ChromeDefaultData\Favicons`
- [x] Login Data (SQLite):
  - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\Default\Login Data`
  - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\ChromeDefaultData\Login Data`
- [x] Extensions & Addons:
  - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\Default\Extensions`
  - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\ChromeDefaultData\Extensions`
- [] Session Data:
  - [] Current Session:
    - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\Default\Current Session`
    - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\ChromeDefaultData\Current Session`
    - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\Default\Current Tabs`
    - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\ChromeDefaultData\Current Tabs`
  - [] Last Session:
    - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\Default\Last Session`
    - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\ChromeDefaultData\Last Session`
    - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\Default\Last Tabs`
    - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\ChromeDefaultData\Last Tabs`
- [] Thumbnails (SQLite):
  - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\Default\Top Sites`
  - `C:\Users\XXX\AppData\Local\Google\Chrome\User Data\ChromeDefaultData\Thumbnails`



## Tested on

- [x] Windows
- [] Debian 
- [] MacOSs