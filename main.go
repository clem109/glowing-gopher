package main

import (
	"log"
	"io/ioutil"
	"fmt"
	"net/http"
	"sync"
	// "github.com/gin-gonic/gin"
	"github.com/go-yaml/yaml" 
)  

// Config is a struct to parse the healthcheck yaml.
// for go-yaml to work correctly all types need to be
// exported otherwise it won't work.
type Config struct {
	Push []struct {
		Name			  string 	`yaml:"name"`
		URL 			  string 	`yaml:"url"`
		Expression  string 	`yaml:"expression"`
	} `yaml:"push"`
	Receive []struct {
		Name 				string 	`yaml:"name"`
		Message 		string 	`yaml:"message"`
		Token 			string 	`yaml:"token"`
	} `yaml:"receive"`
}

// Health shows the status of our apps to anyone calling our healtcheck endpoint
type Health struct {
	Endpoints []struct {
		Name 			string `json:"name"`
		Status		string `json:"status"`
		Message 	string `json:"message"`
	}
	Jobs []struct {
		Name 			string `json:"name"`
		Status		string `json:"status"`
		Message 	string `json:"message"`
		Time 			string `json:"last_trigger"` 
	}
}

var err error
var healthCheck []byte

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func unMarshalYaml (c *Config, path string) {
	healthCheck, err = ioutil.ReadFile(path)
	checkError(err)
	err = yaml.Unmarshal([]byte(healthCheck), &c) 
	checkError(err)

}

func checkHealth (h *Health, c *Config) {

	testEndpoints := func (wg *sync.WaitGroup, url string, name string) {
		client := &http.Client{}
		resp, err := client.Get(url)
		checkError(err)
		defer resp.Body.Close()

		fmt.Printf("Name: %s\nURL: %s \nResponse Status: %v\n\n", name, url, resp.Status)
		defer wg.Done()
	}

	var wg sync.WaitGroup
	wg.Add(len(c.Push))

	for i := 0; i < len(c.Push); i++ {
    go testEndpoints(&wg, c.Push[i].URL, c.Push[i].Name)
  }
	wg.Wait()

}

func main() {

	var c Config
	var h Health
	unMarshalYaml(&c, "./config.yml")
	checkHealth(&h, &c)
		
	// r := gin.Default()

	// r.GET("/healthcheck", func(g *gin.Context) {
	// 	g.JSON(200, gin.H{
	// 		"json": "healthcheck",
	// 	})
	// })

	// r.Run(":3333") // listen and serve on 0.0.0.0:3333
}