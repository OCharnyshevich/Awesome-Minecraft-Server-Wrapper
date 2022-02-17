package minecraft

import (
	"encoding/json"
	"fmt"
	"github.com/OCharnyshevich/Awesome-Minecraft-Server-Wrapper/http"
	"log"
)

const VersionManifestURL = "https://launchermeta.mojang.com/mc/game/version_manifest.json"

type VersionManifest struct {
	Latest   Latest    `json:"latest"`
	Versions []Version `json:"versions"`
}

type Latest struct {
	Release  string `json:"release"`
	Snapshot string `json:"snapshot"`
}

type Version struct {
	ID          string `json:"id"` // Version number
	Type        string `json:"type"`
	URL         string `json:"url"`
	Time        string `json:"time"`
	ReleaseTime string `json:"ReleaseTime"`
}

type VersionDetail struct {
	Assets                 string                        `json:"assets"`
	ID                     string                        `json:"id"`
	ComplianceLevel        int                           `json:"complianceLevel"`
	Downloads              map[string]VersionDetailsFile `json:"downloads"`
	JavaVersion            JavaVersion                   `json:"javaVersion"`
	MainClass              string                        `json:"mainClass"`
	MinimumLauncherVersion uint                          `json:"minimumLauncherVersion"`
	ReleaseTime            string                        `json:"releaseTime"`
	Time                   string                        `json:"time"`
	Type                   string                        `json:"type"`

	//arguments  string `json:"arguments"`
	//assetIndex string `json:"assetIndex"`
	//libraries string `json:"libraries"`
	//logging string `json:"logging"`
}

type VersionDetailsFile struct {
	Sha1 string `json:"sha1"`
	Size uint   `json:"size"`
	URL  string `json:"url"`
}

type JavaVersion struct {
	Component    string `json:"component"`
	MajorVersion uint   `json:"majorVersion"`
}

func PreloadManifest(versionManifest *VersionManifest) {
	body := http.DownloadJSON(VersionManifestURL)
	jsonErr := json.Unmarshal(body, versionManifest)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}

func (m VersionManifest) GetDetails(version string) (*VersionDetail, error) {
	versionDetail := &VersionDetail{}
	v, e := m.getVersion(version)

	if e != nil {
		return nil, e
	}

	body := http.DownloadJSON(v.URL)
	jsonErr := json.Unmarshal(body, versionDetail)

	return versionDetail, jsonErr
}

func (m VersionManifest) getVersion(version string) (*Version, error) {
	for _, v := range m.Versions {
		if version == v.ID {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("version '%s' dosn't exist", version)
}
