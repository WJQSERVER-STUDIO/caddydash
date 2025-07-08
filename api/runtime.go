package api

import (
	"runtime"

	"github.com/infinite-iroha/touka"
)

func runtimeInfo() touka.HandlerFunc {
	return func(c *touka.Context) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		c.JSON(200, touka.H{
			"alloc":       m.Alloc,
			"total_alloc": m.TotalAlloc,
			"heap_alloc":  m.HeapAlloc,
			"next_gc":     m.NextGC,

			"sys":         m.Sys,
			"num_gc":      m.NumGC,
			"go_routines": runtime.NumGoroutine(),
		})
	}
}
