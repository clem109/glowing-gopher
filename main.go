package main

import (
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
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
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

func testEndpoints(wg *sync.WaitGroup, url string, name string, health *[]Health) []Health {
	client := &http.Client{}
	resp, err := client.Get(url)
	checkError(err)
	defer resp.Body.Close()

	// add to struct
	h := Health{name, resp.Status, "message"}
	*health = append(*health, h)
	defer wg.Done()

	return *health
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
