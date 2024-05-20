package export

import (
	"encoding/csv"
	"fmt"
	. "local/BrowserArtifact/src"
	"os"
)

func ExportCSV(path string, artifacts []BrowserArtifact) {
	// Open file
	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()

	// Write CSV data to file
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"ArtifactType",
		"Dest",
		"Src",
		"User",
		"App",
		"Action",
		"Cached",
		"Cookie",
		"Url",
		"UrlDomain",
		"HttpMethod",
		"HttpReferrer",
		"HttpUserAgent",
		"HttpContentType",
		"Duration",
		"Status",
		"BytesIn",
		"BytesOut",
		"Timestamp",
		"TimestampType",
		"Typed",
		"VisitCount",
		"Title",
		"BookmarkTitle",
		"Metadata",
		"Filename",
		"Fieldname",
		"Value",
		"AddonName",
		"AddonType",
		"Active",
		"Visible",
		"Description",
		"FullDescription",
		"Version",
		"SourceURI",
		"HomePageURL",
		"AboutURL",
		"ReviewURL",
		"CreatorName",
		"CreatorURL",
		"AverageRating",
		"RatingCount",
	}
	err = writer.Write(header)

	for _, artifact := range artifacts {
		record := []string{
			artifact.ArtifactType,
			artifact.Dest,
			artifact.Src,
			artifact.User,
			artifact.App,
			artifact.Action,
			fmt.Sprintf("%t", artifact.Cached),
			artifact.Cookie,
			artifact.Url,
			artifact.UrlDomain,
			artifact.HttpMethod,
			artifact.HttpReferrer,
			artifact.HttpUserAgent,
			artifact.HttpContentType,
			fmt.Sprintf("%d", artifact.Duration),
			artifact.Status,
			fmt.Sprintf("%d", artifact.BytesIn),
			fmt.Sprintf("%d", artifact.BytesOut),
			fmt.Sprintf("%d", artifact.Timestamp),
			artifact.TimestampType,
			fmt.Sprintf("%d", artifact.Typed),
			fmt.Sprintf("%d", artifact.VisitCount),
			artifact.Title,
			artifact.BookmarkTitle,
			artifact.Metadata,
			artifact.Filename,
			artifact.Fieldname,
			artifact.Value,
			artifact.AddonName,
			artifact.AddonType,
			fmt.Sprintf("%t", artifact.Active),
			fmt.Sprintf("%t", artifact.Visible),
			artifact.Description,
			artifact.FullDescription,
			artifact.Version,
			artifact.SourceURI,
			artifact.HomePageURL,
			artifact.AboutURL,
			artifact.ReviewURL,
			artifact.CreatorName,
			artifact.CreatorURL,
			fmt.Sprintf("%f", artifact.AverageRating),
			fmt.Sprintf("%d", artifact.RatingCount),
		}
		err := writer.Write(record)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
