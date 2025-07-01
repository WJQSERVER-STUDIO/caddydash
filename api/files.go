package api

import (
	"caddydash/db"

	"github.com/infinite-iroha/touka"
)

func FilesParams(cdb *db.ConfigDB) touka.HandlerFunc {
	return func(c *touka.Context) {
		params, err := cdb.RangeAllParams()
		if err != nil {
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}
		c.JSON(200, params)
	}
}

func FilesTemplates(cdb *db.ConfigDB) touka.HandlerFunc {
	return func(c *touka.Context) {
		templates, err := cdb.GetAllTemplates()
		if err != nil {
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}
		c.JSON(200, templates)
	}
}

func FilesRendered(cdb *db.ConfigDB) touka.HandlerFunc {
	return func(c *touka.Context) {
		rendered, err := cdb.RangeAllReandered()
		if err != nil {
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}
		c.JSON(200, rendered)
	}
}
