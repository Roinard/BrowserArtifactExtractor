package chrome

type bookmark struct {
	Checksum string `json:"checksum"`
	Roots    struct {
		BookmarkBar struct {
			Children []struct {
				Children []struct {
					DateAdded    string `json:"date_added"`
					DateLastUsed string `json:"date_last_used"`
					GUID         string `json:"guid"`
					ID           string `json:"id"`
					Name         string `json:"name"`
					Type         string `json:"type"`
					URL          string `json:"url,omitempty"`
					Children     []struct {
						DateAdded    string `json:"date_added"`
						DateLastUsed string `json:"date_last_used"`
						GUID         string `json:"guid"`
						ID           string `json:"id"`
						Name         string `json:"name"`
						Type         string `json:"type"`
						URL          string `json:"url"`
					} `json:"children,omitempty"`
					DateModified string `json:"date_modified,omitempty"`
					MetaInfo     struct {
						PowerBookmarkMeta string `json:"power_bookmark_meta"`
					} `json:"meta_info,omitempty"`
				} `json:"children,omitempty"`
				DateAdded    string `json:"date_added"`
				DateLastUsed string `json:"date_last_used"`
				DateModified string `json:"date_modified,omitempty"`
				GUID         string `json:"guid"`
				ID           string `json:"id"`
				Name         string `json:"name"`
				Type         string `json:"type"`
				URL          string `json:"url,omitempty"`
				MetaInfo     struct {
					LastVisitedDesktop string `json:"last_visited_desktop"`
				} `json:"meta_info,omitempty"`
			} `json:"children"`
			DateAdded    string `json:"date_added"`
			DateLastUsed string `json:"date_last_used"`
			DateModified string `json:"date_modified"`
			GUID         string `json:"guid"`
			ID           string `json:"id"`
			Name         string `json:"name"`
			Type         string `json:"type"`
		} `json:"bookmark_bar"`
		Other struct {
			Children     []any  `json:"children"`
			DateAdded    string `json:"date_added"`
			DateLastUsed string `json:"date_last_used"`
			DateModified string `json:"date_modified"`
			GUID         string `json:"guid"`
			ID           string `json:"id"`
			Name         string `json:"name"`
			Type         string `json:"type"`
		} `json:"other"`
		Synced struct {
			Children []struct {
				Children []struct {
					DateAdded    string `json:"date_added"`
					DateLastUsed string `json:"date_last_used"`
					GUID         string `json:"guid"`
					ID           string `json:"id"`
					Name         string `json:"name"`
					Type         string `json:"type"`
					URL          string `json:"url"`
				} `json:"children,omitempty"`
				DateAdded    string `json:"date_added"`
				DateLastUsed string `json:"date_last_used"`
				DateModified string `json:"date_modified,omitempty"`
				GUID         string `json:"guid"`
				ID           string `json:"id"`
				Name         string `json:"name"`
				Type         string `json:"type"`
				URL          string `json:"url,omitempty"`
				MetaInfo     struct {
					LastVisited string `json:"last_visited"`
				} `json:"meta_info,omitempty"`
				MetaInfo0 struct {
					LastVisitedDesktop string `json:"last_visited_desktop"`
				} `json:"meta_info,omitempty"`
				MetaInfo1 struct {
					LastVisited        string `json:"last_visited"`
					LastVisitedDesktop string `json:"last_visited_desktop"`
				} `json:"meta_info,omitempty"`
			} `json:"children"`
			DateAdded    string `json:"date_added"`
			DateLastUsed string `json:"date_last_used"`
			DateModified string `json:"date_modified"`
			GUID         string `json:"guid"`
			ID           string `json:"id"`
			Name         string `json:"name"`
			Type         string `json:"type"`
		} `json:"synced"`
	} `json:"roots"`
	SyncMetadata string `json:"sync_metadata"`
	Version      int    `json:"version"`
}

type extension struct {
	Author     string `json:"author"`
	Background struct {
		Page string `json:"page"`
	} `json:"background"`
	BrowserAction struct {
		DefaultIcon struct {
			Num16 string `json:"16"`
			Num32 string `json:"32"`
			Num64 string `json:"64"`
		} `json:"default_icon"`
		DefaultPopup string `json:"default_popup"`
		DefaultTitle string `json:"default_title"`
	} `json:"browser_action"`
	Commands struct {
		LaunchElementPicker struct {
			Description string `json:"description"`
		} `json:"launch-element-picker"`
		LaunchElementZapper struct {
			Description string `json:"description"`
		} `json:"launch-element-zapper"`
		LaunchLogger struct {
			Description string `json:"description"`
		} `json:"launch-logger"`
		OpenDashboard struct {
			Description string `json:"description"`
		} `json:"open-dashboard"`
		RelaxBlockingMode struct {
			Description string `json:"description"`
		} `json:"relax-blocking-mode"`
		ToggleCosmeticFiltering struct {
			Description string `json:"description"`
		} `json:"toggle-cosmetic-filtering"`
		ToggleJavascript struct {
			Description string `json:"description"`
		} `json:"toggle-javascript"`
	} `json:"commands"`
	ContentScripts []struct {
		AllFrames       bool     `json:"all_frames"`
		Js              []string `json:"js"`
		MatchAboutBlank bool     `json:"match_about_blank,omitempty"`
		Matches         []string `json:"matches"`
		RunAt           string   `json:"run_at"`
	} `json:"content_scripts"`
	ContentSecurityPolicy   string `json:"content_security_policy"`
	DefaultLocale           string `json:"default_locale"`
	Description             string `json:"description"`
	DifferentialFingerprint string `json:"differential_fingerprint"`
	Icons                   struct {
		Num16  string `json:"16"`
		Num32  string `json:"32"`
		Num64  string `json:"64"`
		Num128 string `json:"128"`
	} `json:"icons"`
	Incognito            string `json:"incognito"`
	Key                  string `json:"key"`
	ManifestVersion      int    `json:"manifest_version"`
	MinimumChromeVersion string `json:"minimum_chrome_version"`
	Name                 string `json:"name"`
	OptionsUI            struct {
		OpenInTab bool   `json:"open_in_tab"`
		Page      string `json:"page"`
	} `json:"options_ui"`
	Permissions []string `json:"permissions"`
	ShortName   string   `json:"short_name"`
	Storage     struct {
		ManagedSchema string `json:"managed_schema"`
	} `json:"storage"`
	UpdateURL              string   `json:"update_url"`
	Version                string   `json:"version"`
	WebAccessibleResources []string `json:"web_accessible_resources"`
}
