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
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Usage: ./main <absolute path to configuration file>")
		return
	}
	file, _ := os.Open(args[0])
	decoder := json.NewDecoder(file)
	var conf = Configuration{}
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("error:", err)
	}
	predator := NewPredator(&conf)
	predator.Run()
}
