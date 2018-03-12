package main

import (
	"log"
	"io/ioutil"
	"fmt"
	// "github.com/gin-gonic/gin"
	"github.com/go-yaml/yaml" 
)  

// Config is a struct to parse the healthcheck yaml.
// for go-yaml to work correctly all types need to be
// exported otherwise it won't work.
type Config struct {
	Push []struct {
		Name			 string 	`yaml:"name"`
		URL 			 string 	`yaml:"url"`
		Expression string `yaml:"expression"`
	} `yaml:"push"`
	Receive []struct {
		Name 		string `yaml:"name"`
		Message string `yaml:"message"`
		Token 	string `yaml:"token"`
	} `yaml:"receive"`
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

func main() {

	var c Config
	unMarshalYaml(&c, "./config.yml")
	fmt.Println(c)
	


	// r := gin.Default()

	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })



	// r.GET("/healthcheck", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"json": "healthcheck",
	// 	})
	// })

	// r.Run(":3333") // listen and serve on 0.0.0.0:3333
}