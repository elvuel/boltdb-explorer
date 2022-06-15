package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/elvuel/boltdb-explorer/backend/bolt"
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
)

var (
	port = os.Getenv("PORT")
)

func debug(o interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(o)
}

func init() {
	if port == "" {
		port = "8080"
	}
}

func main() {
	engine := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	engine.Use(cors.New(corsConfig))

	engine.POST("/upload", func(c *gin.Context) {
		c.Request.ParseMultipartForm(32 << 20)
		f, h, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		defer f.Close()

		n := fmt.Sprintf("bolt_%d.tmp.db", time.Now().Unix())
		tmpFile, err := os.OpenFile(n, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer tmpFile.Close()

		_, err = io.Copy(tmpFile, f)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result, err := bolt.Read(n)
		defer os.Remove(n)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": result, "filename": h.Filename})
	})

	engine.POST("/download", func(c *gin.Context) {
		payload := make(map[string]interface{})
		err := json.NewDecoder(c.Request.Body).Decode(&payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer c.Request.Body.Close()

		n := fmt.Sprintf("bolt_%d.tmp.download.db", time.Now().Unix())
		_, err = bolt.WriteMap(n, payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer os.Remove(n)
		c.FileAttachment(n, n)
	})

	engine.Run(":" + port)
}
