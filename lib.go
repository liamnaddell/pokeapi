package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Pokemon struct {
	Weight int    `json:"weight"`
	Name   string `json:"name"`
	Id     int    `json:"id"`
	Types  []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

var home = "/home/liam"

var pokeclient = http.Client{
	Timeout: time.Second * 10, // Maximum of 2 secs
}

func checkerr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func cachePokemon(name string, body []byte) {
	file, err := os.Create(home + "/.pokeapi/" + name)
	checkerr(err)
	_, err2 := file.Write(body)
	checkerr(err2)
}

//check for file
func cff(name string) bool {
	_, err := os.Stat(home + "/.pokeapi/" + name)
	if err != nil {
		return false
	}
	return true

}

func getLink(link string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "pokeapi-getter-go")

	res, getErr := pokeclient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	return res
}
func getPokemon(name string) (*Pokemon, error) {
	url := "http://pokeapi.co/api/v2/pokemon/" + name
	var body []byte
	var ff = cff(name)
	if ff {
		file, err := ioutil.ReadFile(home + "/.pokeapi/" + name)
		body = file
		checkerr(err)
	} else {
		res := getLink(url)
		var readErr error
		body, readErr = ioutil.ReadAll(res.Body)
		checkerr(readErr)
	}

	pokemon := Pokemon{}
	jsonErr := json.Unmarshal(body, &pokemon)
	checkerr(jsonErr)
	if pokemon.Weight == 0 && pokemon.Name == "" {
		return new(Pokemon), errors.New(fmt.Sprintf("%s is not a pokemon", name))
	}
	if !ff {
		cachePokemon(name, body)
	}
	return &pokemon, nil
}

func StartGetPokemon(name string) (*Pokemon, error) {
	pokemon, err := getPokemon(name)
	return pokemon, err
}
