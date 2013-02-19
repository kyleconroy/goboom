package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Config struct {
	Backend string
}

type Client struct {
}

type List struct {
	Name  string
	Items []Item
}

type Item struct {
	Name  string
	Value string
}

//
type Storage interface {
	Lists() []List
	GetList(name string) List
	Populate() error
}

type jsonStore struct {
	Lists []List
}

type JsonStorage struct {
	Filename string
	db       jsonStore
}

func load(config Config) (Storage, error) {
	storage := JsonStorage{Filename: "temp.json"}
	err := storage.Populate()
	return &storage, err
}

// The list of Lists in your JSON data, sorted by number of items descending.
func (js *JsonStorage) Lists() []List {
	return js.db.Lists
}

// The list of Lists in your JSON data, sorted by number of items descending.
func (js *JsonStorage) GetList(name string) List {
	return List{}
}

// The list of Lists in your JSON data, sorted by number of items descending.
func (js *JsonStorage) Populate() error {
	filebyte, err := ioutil.ReadFile(js.Filename)
	if err != nil {
		return err
	}
	var db jsonStore
	err = json.Unmarshal(filebyte, &db)

	if err != nil {
		return err
	}

	fmt.Println(db)

	js.db = db

	return nil
}

func main() {
	storage, err := load(Config{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(storage.Lists()))
}
