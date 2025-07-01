package api

import (
	"runtime/debug"

	"github.com/infinite-iroha/touka"
)

type InfoApiStruct struct {
	Version      string   `json:"version"`
	License      string   `json:"license"`
	Author       []string `json:"author"`
	BuildVersion string   `json:"build_version"`
	GoVersion    string   `json:"go_version"`
}

func infoHandle(version string) touka.HandlerFunc {
	return func(c *touka.Context) {
		buildinfo, ok := debug.ReadBuildInfo()
		if !ok {
			c.JSON(500, touka.H{"error": "no build info"})
			return
		}
		c.JSON(200, InfoApiStruct{
			Version: version,
			License: "Mozilla Public License 2.0",
			//Author:       "WJQSERVER",
			Author:       []string{"WJQSERVER"},
			BuildVersion: buildinfo.Main.Version,
			GoVersion:    buildinfo.GoVersion,
		})
	}
}
