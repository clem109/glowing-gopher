package main

import (
	"log"
	"fmt"
	"github.com/gin-gonic/gin"
"github.com/go-yaml/yaml" 
)  

type StructA struct {
	A string `yaml:"a"`
}

type StructB struct {
	// Embedded structs are not treated as embedded in YAML by default. To do that,
	// add the ",inline" annotation below
	StructA `yaml:",inline"`
	B       string `yaml:"b"`
}

var data = `
a: a string from struct A
b: a string from struct B
`

func main() {

	var b StructB

	err := yaml.Unmarshal([]byte(data), &b)
	if err != nil {
			log.Fatalf("cannot unmarshal data: %v", err)
	}
	fmt.Println(b.A)
	fmt.Println(b.B)

	
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})



	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"json": "healthcheck",
		})
	})

	r.Run(":3333") // listen and serve on 0.0.0.0:3333
}