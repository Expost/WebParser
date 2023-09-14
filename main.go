package main

import (
	"WebParser/rewrite"
	"WebParser/sanitizer"
	"WebParser/scraper"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.GET("/parser", func(c *gin.Context) {
		url := c.Query("url")
		content, scraperErr := scraper.Fetch(
			url,
			"",
			"",
			"",
			true,
			false,
		)

		if scraperErr != nil {
			c.JSON(200, gin.H{
				"error":    true,
				"messages": scraperErr,
			})
			return
		}

		content = rewrite.Rewriter(url, content, "")

		content = sanitizer.Sanitize(url, content)
		c.JSON(200, gin.H{
			"content": content,
			"url":     url,
		})

	})
	r.Run(":3002")

}
