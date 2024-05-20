package firefox

type loginJSON struct {
	NextID int `json:"nextId"`
	Logins []struct {
		ID                  int         `json:"id"`
		Hostname            string      `json:"hostname"`
		HTTPRealm           interface{} `json:"httpRealm"`
		FormSubmitURL       string      `json:"formSubmitURL"`
		UsernameField       string      `json:"usernameField"`
		PasswordField       string      `json:"passwordField"`
		EncryptedUsername   string      `json:"encryptedUsername"`
		EncryptedPassword   string      `json:"encryptedPassword"`
		GUID                string      `json:"guid"`
		EncType             int         `json:"encType"`
		TimeCreated         int64       `json:"timeCreated"`
		TimeLastUsed        int64       `json:"timeLastUsed"`
		TimePasswordChanged int64       `json:"timePasswordChanged"`
		TimesUsed           int         `json:"timesUsed"`
	} `json:"logins"`
	PotentiallyVulnerablePasswords   []interface{} `json:"potentiallyVulnerablePasswords"`
	DismissedBreachAlertsByLoginGUID struct {
	} `json:"dismissedBreachAlertsByLoginGUID"`
	Version int `json:"version"`
}

type addonJSON struct {
	Schema int `json:"schema"`
	Addons []struct {
		ID    string `json:"id"`
		Icons struct {
			Num32  string `json:"32"`
			Num64  string `json:"64"`
			Num128 string `json:"128"`
		} `json:"icons"`
		Name            string `json:"name"`
		Version         string `json:"version"`
		SourceURI       string `json:"sourceURI"`
		HomepageURL     string `json:"homepageURL"`
		SupportURL      string `json:"supportURL"`
		AmoListingURL   string `json:"amoListingURL"`
		Description     string `json:"description"`
		FullDescription string `json:"fullDescription"`
		WeeklyDownloads int    `json:"weeklyDownloads"`
		Type            string `json:"type"`
		Creator         struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"creator"`
		Developers  []interface{} `json:"developers"`
		Screenshots []struct {
			URL             string `json:"url"`
			Width           int    `json:"width"`
			Height          int    `json:"height"`
			ThumbnailURL    string `json:"thumbnailURL"`
			ThumbnailWidth  int    `json:"thumbnailWidth"`
			ThumbnailHeight int    `json:"thumbnailHeight"`
			Caption         string `json:"caption"`
		} `json:"screenshots"`
		ContributionURL string  `json:"contributionURL"`
		AverageRating   float64 `json:"averageRating"`
		ReviewCount     int     `json:"reviewCount"`
		ReviewURL       string  `json:"reviewURL"`
		UpdateDate      int64   `json:"updateDate"`
	} `json:"addons"`
}

type extensionJSON struct {
	SchemaVersion int `json:"schemaVersion"`
	Addons        []struct {
		ID                  string      `json:"id"`
		SyncGUID            string      `json:"syncGUID"`
		Version             string      `json:"version"`
		Type                string      `json:"type"`
		Loader              interface{} `json:"loader"`
		UpdateURL           interface{} `json:"updateURL"`
		InstallOrigins      interface{} `json:"installOrigins"`
		ManifestVersion     int         `json:"manifestVersion"`
		OptionsURL          interface{} `json:"optionsURL"`
		OptionsType         interface{} `json:"optionsType"`
		OptionsBrowserStyle bool        `json:"optionsBrowserStyle"`
		AboutURL            interface{} `json:"aboutURL"`
		DefaultLocale       struct {
			Name         string      `json:"name"`
			Description  string      `json:"description"`
			Creator      string      `json:"creator"`
			HomepageURL  string      `json:"homepageURL"`
			Developers   interface{} `json:"developers"`
			Translators  interface{} `json:"translators"`
			Contributors interface{} `json:"contributors"`
		} `json:"defaultLocale"`
		Visible                bool        `json:"visible"`
		Active                 bool        `json:"active"`
		UserDisabled           bool        `json:"userDisabled"`
		AppDisabled            bool        `json:"appDisabled"`
		EmbedderDisabled       bool        `json:"embedderDisabled"`
		InstallDate            int64       `json:"installDate"`
		UpdateDate             int64       `json:"updateDate,omitempty"`
		ApplyBackgroundUpdates int         `json:"applyBackgroundUpdates"`
		Path                   string      `json:"path"`
		Skinnable              bool        `json:"skinnable"`
		SourceURI              string      `json:"sourceURI"`
		ReleaseNotesURI        interface{} `json:"releaseNotesURI"`
		SoftDisabled           bool        `json:"softDisabled"`
		ForeignInstall         bool        `json:"foreignInstall"`
		StrictCompatibility    bool        `json:"strictCompatibility"`
		Locales                []struct {
			Name         string      `json:"name"`
			Description  string      `json:"description"`
			Creator      string      `json:"creator"`
			HomepageURL  string      `json:"homepageURL"`
			Developers   interface{} `json:"developers"`
			Translators  interface{} `json:"translators"`
			Contributors interface{} `json:"contributors"`
			Locales      []string    `json:"locales"`
		} `json:"locales"`
		TargetApplications []struct {
			ID         string `json:"id"`
			MinVersion string `json:"minVersion"`
			MaxVersion string `json:"maxVersion"`
		} `json:"targetApplications"`
		TargetPlatforms []interface{} `json:"targetPlatforms"`
		SignedState     int           `json:"signedState,omitempty"`
		SignedDate      int64         `json:"signedDate"`
		Seen            bool          `json:"seen"`
		Dependencies    []interface{} `json:"dependencies"`
		Incognito       string        `json:"incognito,omitempty"`
		UserPermissions struct {
			Permissions []string `json:"permissions"`
			Origins     []string `json:"origins"`
		} `json:"userPermissions"`
		OptionalPermissions struct {
			Permissions []interface{} `json:"permissions"`
			Origins     []interface{} `json:"origins"`
		} `json:"optionalPermissions"`
		Icons struct {
			Num16  string `json:"16"`
			Num48  string `json:"48"`
			Num128 string `json:"128"`
		} `json:"icons"`
		IconURL              interface{} `json:"iconURL"`
		BlocklistState       int         `json:"blocklistState"`
		BlocklistURL         interface{} `json:"blocklistURL"`
		StartupData          interface{} `json:"startupData"`
		Hidden               bool        `json:"hidden"`
		InstallTelemetryInfo struct {
			Source    string `json:"source"`
			SourceURL string `json:"sourceURL"`
			Method    string `json:"method"`
		} `json:"installTelemetryInfo"`
		RecommendationState interface{} `json:"recommendationState"`
		RootURI             string      `json:"rootURI"`
		Location            string      `json:"location"`
	} `json:"addons"`
}
