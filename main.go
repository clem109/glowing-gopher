package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	// "github.com/gin-gonic/gin"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-yaml/yaml"
)

// Config is a struct to parse the healthcheck yaml.
// for go-yaml to work correctly all types need to be
// exported otherwise it won't work.
type Config struct {
	Push []struct {
		Name       string `yaml:"name"`
		URL        string `yaml:"url"`
		Expression string `yaml:"expression"`
	} `yaml:"push"`
	Receive []struct {
		Name    string `yaml:"name"`
		Message string `yaml:"message"`
		Token   string `yaml:"token"`
	} `yaml:"receive"`
}

// Health shows the status of our apps to anyone calling our healtcheck endpoint
type Health struct {
	Endpoints []struct {
		Name    string `json:"name"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	Jobs []struct {
		Name    string `json:"name"`
		Status  string `json:"status"`
		Message string `json:"message"`
		Time    string `json:"last_trigger"`
	}
}

var err error

// UTILITY FUNCTIONS
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var healthCheck []byte

func unMarshalYaml(c *Config, path string) {
	healthCheck, err = ioutil.ReadFile(path)
	checkError(err)
	err = yaml.Unmarshal([]byte(healthCheck), &c)
	checkError(err)

}

func checkHealth(h *Health, c *Config) {
	testEndpoints := func(wg *sync.WaitGroup, url string, name string) {
		client := &http.Client{}
		resp, err := client.Get(url)
		checkError(err)
		defer resp.Body.Close()
		fmt.Printf("Name: %s\nURL: %s \nResponse Status: %v\nTime: %v\n",
			name, url, resp.Status, time.Now())
		defer wg.Done()
	}

	var wg sync.WaitGroup
	wg.Add(len(c.Push))

	for i := 0; i < len(c.Push); i++ {
		go testEndpoints(&wg, c.Push[i].URL, c.Push[i].Name)
	}
	wg.Wait()
}

func startServer(port string) {
	r := gin.Default()

	r.GET("/healthcheck", func(g *gin.Context) {
		checkHealth(&h, &c)
		g.JSON(200, c.Push)
	})

	r.Run(port) // listen and serve on the port provided
}

var c Config
var h Health

func main() {
	unMarshalYaml(&c, "./config.yml")
	checkHealth(&h, &c)
	port := getEnv("PORT", ":3333")
	startServer(port)

}
