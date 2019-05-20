package main

import (
  "github.com/BurntSushi/toml"
  "fmt"
)

var conf Config

//Read all configurations originating from specified directory
func readConfig(dir string){
  if _, err := toml.DecodeFile(dir, &conf); err != nil {
    fmt.Println(err)
  }
}
