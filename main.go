package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-yaml/yaml"
)

// Config is a struct to parse the healthcheck yaml.
// for go-yaml to work correctly all types need to be
// exported otherwise it won't work.
type Config struct {
	Push []struct {
		Name string `yaml:"name"`
		URL  string `yaml:"url"`
	} `yaml:"push"`
}

// Health shows the status of our apps to anyone calling our healtcheck endpoint
type Health struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Status string `json:"status"`
	Error  error  `json:"error"`
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

func testEndpoints(wg *sync.WaitGroup, url string, name string, health *[]Health) error {
	defer wg.Done()
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error: Failed to execute the HTTP request. %s\n", err)
		h := Health{name, url, "404", err}
		*health = append(*health, h)
		return err
	}
	defer resp.Body.Close()

	// add to struct
	h := Health{name, url, resp.Status, nil}
	*health = append(*health, h)
	return nil

}

func checkHealth(c *Config) []Health {
	var healthArray []Health
	var wg sync.WaitGroup

	wg.Add(len(c.Push))
	for i := 0; i < len(c.Push); i++ {
		go testEndpoints(&wg, c.Push[i].URL, c.Push[i].Name, &healthArray)
	}
	wg.Wait()
	return healthArray
}

func startServer(port string) {
	r := gin.Default()

	r.GET("/healthcheck", func(g *gin.Context) {
		resp := checkHealth(&c)
		g.JSON(200, resp)
	})

	r.Run(port) // listen and serve on the port provided
}

var c Config
var h Health

func main() {
	unMarshalYaml(&c, "./config.yml")
	port := getEnv("PORT", ":3333")
	startServer(port)
}
