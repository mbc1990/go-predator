package main

import "encoding/json"
import "fmt"
import "os"

type Configuration struct {
	TwitterConsumerKey    string
	TwitterConsumerSecret string
	TwitterMaxConcurrent  int
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
