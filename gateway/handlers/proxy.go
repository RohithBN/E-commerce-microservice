package handlers

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

var serviceURLs = map[string]string{
	"auth":     "http://localhost:8081",
	"products": "http://localhost:8082",
	"cart":     "http://localhost:8083",
	"orders":   "http://localhost:8084",
}

func ProxyHandler(service, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get service URL
		serviceURL, exists := serviceURLs[service]
		if !exists {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Service not found"})
			return
		}

		// Create the target URL
		target, err := url.Parse(fmt.Sprintf("%s%s", serviceURL, path))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse target URL"})
			return
		}

		// Create reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(target)

		// Update the request URL
		c.Request.URL.Host = target.Host
		c.Request.URL.Scheme = target.Scheme
		c.Request.Header.Set("X-Forwarded-Host", c.Request.Header.Get("Host"))

		// Forward user information if authenticated
		if user, exists := c.Get("user"); exists {
			c.Request.Header.Set("X-User", fmt.Sprintf("%v", user))
			c.Request.Header.Set("X-User-ID", fmt.Sprintf("%v", user.(map[string]interface{})["id"]))

		}

		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)

			modifiedPath := path

			// used for replacing :productId with the actual product ID 
			for _, param := range c.Params {
				modifiedPath = strings.Replace(modifiedPath, ":"+param.Key, param.Value, 1)
			}

			req.URL.Path = modifiedPath
		}

		// Serve the request
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// Helper function to get service URL
func GetServiceURL(service string) (string, error) {
	if url, exists := serviceURLs[service]; exists {
		return url, nil
	}
	return "", fmt.Errorf("service %s not found", service)
}
