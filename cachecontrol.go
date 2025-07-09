package main

import (
	"net/http"

	"github.com/fenthope/cachectrl"
	"github.com/infinite-iroha/touka"
)

var rules = make(map[string]http.Header)

func setCacheControl(r *touka.Engine) {
	rules["/v0/api/"] = http.Header{
		"Cache-Control": {"no-cache", "no-store", "must-revalidate"},
		"Pragma":        {"no-cache"},
		"Expires":       {"0"},
	}

	rules["/locales/"] = http.Header{
		"Cache-Control": {"public", "max-age=43200", "must-revalidate"},
	}

	rules["/js/"] = http.Header{
		"Cache-Control": {"public", "max-age=21600", "must-revalidate"},
	}

	rules["/css/"] = http.Header{
		"Cache-Control": {"public", "max-age=43200", "must-revalidate"},
	}

	noMatchHeaders := http.Header{
		"Cache-Control": {"no-cache", "no-store", "must-revalidate"},
	}

	r.Use(cachectrl.NewCacheControl(cachectrl.CacheOptions{
		Rules:          rules,
		NoMatchHeaders: noMatchHeaders,
	}))
}
