package main

import "encoding/json"
import "fmt"
import "os"

type Configuration struct {
	UnclassifiedWorkDir   string
	TwitterConsumerKey    string
	TwitterConsumerSecret string
	TwitterMaxConcurrent  int
	TwitterSources        []string
}

func main() {
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	var conf = Configuration{}
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("error:", err)
	}
	predator := NewPredator(&conf)
	predator.Run()
}
