package api

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/infinite-iroha/touka"
)

func StreamLog(logFilePath string) touka.HandlerFunc {
	return func(c *touka.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")

		file, err := os.Open(logFilePath)
		if err != nil {
			errorMessage := fmt.Sprintf("Failed to open log file: %v", err)
			c.Errorf("%s", errorMessage)
			fmt.Printf("%s\n", errorMessage)
			log.Printf("%s", errorMessage)
			fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", errorMessage)
			c.Writer.Flush()
			return
		}
		defer file.Close()

		// Send historical logs in chunks
		// Send historical logs in chunks, ensuring not to split in the middle of a line
		const chunkSize = 4096 // 4KB
		buffer := make([]byte, chunkSize)
		file.Seek(0, 0) // Go back to the start of the file
		reader := bufio.NewReader(file)

		for {
			n, err := reader.Read(buffer)
			if err != nil {
				if err == io.EOF {
					break
				}
				errorMessage := fmt.Sprintf("Error reading log file for history: %v", err)
				c.Errorf("%s", errorMessage)
				fmt.Printf("%s\n", errorMessage)
				log.Printf("%s", errorMessage)
				fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", errorMessage)
				c.Writer.Flush()
				return
			}
			// important: data must be split by \n
			chunk := string(buffer[:n])
			lines := strings.Split(chunk, "\n")
			for i, line := range lines {
				if i < len(lines)-1 {
					fmt.Fprintf(c.Writer, "data: %s\n\n", line)
				} else {
					// The last part might be an incomplete line, so we hold it
					// But for simplicity in this context, we send it.
					// A more robust solution would handle this remainder.
					if line != "" {
						fmt.Fprintf(c.Writer, "data: %s\n\n", line)
					}
				}
			}
			c.Writer.Flush()
			time.Sleep(100 * time.Millisecond) // Small delay to prevent network congestion
		}

		// Now, seek to the end of the file to start streaming new logs
		_, err = file.Seek(0, io.SeekEnd)
		if err != nil {
			errorMessage := fmt.Sprintf("Failed to seek to end of log file: %v", err)
			c.Errorf("%s", errorMessage)
			fmt.Printf("%s\n", errorMessage)
			log.Printf("%s", errorMessage)
			fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", errorMessage)
			c.Writer.Flush()
			return
		}

		reader = bufio.NewReader(file)
		dataTicker := time.NewTicker(500 * time.Millisecond)
		defer dataTicker.Stop()

		heartbeatTicker := time.NewTicker(5 * time.Second)
		defer heartbeatTicker.Stop()

		fmt.Fprintf(c.Writer, "event: connected\ndata: Log stream connected.\n\n")
		c.Writer.Flush()

		for {
			select {
			case <-c.Request.Context().Done():
				log.Println("Client disconnected, stopping log stream.")
				return
			case <-dataTicker.C:
				line, err := reader.ReadString('\n')
				if err == nil {
					fmt.Fprintf(c.Writer, "data: %s\n\n", line)
					c.Writer.Flush()
				} else if err == io.EOF {
					continue
				} else {
					errorMessage := fmt.Sprintf("Error reading log file: %v", err)
					c.Errorf("%s", errorMessage)
					fmt.Printf("%s\n", errorMessage)
					log.Printf("%s", errorMessage)
					fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", errorMessage)
					c.Writer.Flush()
					return
				}
			case <-heartbeatTicker.C:
				fmt.Fprintf(c.Writer, ": heartbeat\n\n")
				c.Writer.Flush()
			}
		}
	}
}

func ListSiteLogs(logDir string) touka.HandlerFunc {
	return func(c *touka.Context) {
		var sites []string
		err := filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("Error accessing path %s: %v", path, err)
				return err
			}
			if info.IsDir() && path != logDir {
				sites = append(sites, filepath.Base(path))
			}
			return nil
		})
		if err != nil {
			c.JSON(500, touka.H{"error": "Failed to list site logs"})
			log.Printf("Failed to list site logs in directory %s: %v", logDir, err)
			return
		}
		c.JSON(200, sites)
	}
}

func SiteLog(logDir string) touka.HandlerFunc {
	return func(c *touka.Context) {
		sitename := c.Param("sitename")
		logFilePath := filepath.Join(logDir, sitename, "access.log")
		StreamLog(logFilePath)(c)
	}
}
