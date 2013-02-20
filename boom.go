package main

import (
	"encoding/json"
	"flag"
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

type Store struct {
	Lists []List
	Dirty bool
}

// The list of Lists in your JSON data, sorted by number of items descending.
func (s *Store) Items() []Item {
	items := []Item{}

	for _, list := range s.Lists {
		for _, item := range list.Items {
			items = append(items, item)
		}
	}

	return items
}

func (s *Store) CreateList(name string) {
	s.Lists = append(s.Lists, List{Name: name})
	s.Dirty = true
}

func (s *Store) FindList(name string) (List, bool) {
	for _, list := range s.Lists {
		if list.Name == name {
			return list, true
		}
	}

	return List{}, false
}

func (s *Store) FindItem(name string) (Item, bool) {
	for _, list := range s.Lists {
		for _, item := range list.Items {
			if item.Name == name {
				return item, true
			}
		}
	}

	return Item{}, false
}

func (s *Store) PrintAll() {
	for _, list := range s.Lists {
		for _, item := range list.Items {
			fmt.Println(item.Name)
		}
	}
}

func (s *Store) PrintSummary() {
	for _, list := range s.Lists {
		fmt.Printf("  %v (%d)\n", list.Name, len(list.Items))
	}
}

type Backend interface {
	Save(store Store) error
	Fetch() (Store, error)
}

type JsonBackend struct {
	Filename string
}

func (jb *JsonBackend) Fetch() (Store, error) {
	filebyte, err := ioutil.ReadFile(jb.Filename)

	if err != nil {
		return Store{}, err
	}

	var db Store

	err = json.Unmarshal(filebyte, &db)

	if err != nil {
		return Store{}, err
	}

	return db, nil
}

func (jb *JsonBackend) Save(store Store) error {
	if !store.Dirty {
		return nil
	}

	bytes, err := json.Marshal(store)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(jb.Filename, bytes, 0666)

	if err != nil {
		return err
	}

	return nil
}

func load(config Config) (Backend, error) {
	backend := JsonBackend{Filename: "temp.json"}
	return &backend, nil
}

func showHelpMessage() {
	help := `
	- boom: help ---------------------------------------------------

	boom                          display high-level overview
	boom all                      show all items in all lists
	boom edit                     edit the boom JSON file in $EDITOR
	boom help                     this help text
	boom storage                  shows which storage backend you're using
	boom switch <storage>         switches to a different storage backend

	boom <list>                   create a new list
	boom <list>                   show items for a list
	boom <list> delete            deletes a list

	boom <list> <name> <value>    create a new list item
	boom <name>                   copy item's value to clipboard
	boom <list> <name>            copy item's value to clipboard
	boom open <name>              open item's url in browser
	boom open <list> <name>       open all item's url in browser for a list
	boom random                   open a random item's url in browser
	boom random <list>            open a random item's url for a list in browser
	boom echo <name>              echo the item's value without copying
	boom echo <list> <name>       echo the item's value without copying
	boom <list> <name> delete     deletes an item

	all other documentation is located at:
	https://github.com/holman/boom
	`
	fmt.Println(help)
}

func showEmptyMessage() {
	fmt.Println("You don't have anything yet! To start out, create a new list:")
	fmt.Println("  $ boom <list-name>")
	fmt.Println("And then add something to your list!")
	fmt.Println("  $ boom <list-name> <item-name> <item-value>")
	fmt.Println("You can then grab your new item:")
	fmt.Println("  $ boom <item-name>")
}

type Command struct {
}

func match(commands ...string) bool {
	matched := true
	for i, cmd := range commands {
		if len(flag.Arg(i)) == 0 {
			matched = false
		} else {
			matched = matched && (cmd == flag.Arg(i) || cmd[0] == '<')
		}
	}
	return matched
}



func main() {
	flag.Parse()

	backend, err := load(Config{})

	if err != nil {
		log.Fatal(err)
	}

	store, err := backend.Fetch()

	if err != nil {
		log.Fatal(err)
	}

	if command == "all" {
		return c.all()
	}

	switch {
	case match("all"):
		store.PrintAll()
	case match("help"):
		showHelpMessage()
	case match("switch", "<storage>"):
		//switch
	case match("<list>", "delete"):
		store.DeleteList(flag.Arg(0))

		store.DeleteList(flag.Arg(0))
	case match("<list>", "<name>", "<value>"):
		//store.List(flag.Arg(0))
	case match("<list>"):
		name := flag.Arg(0)
		item, found := store.FindItem(name)

		if found {
			fmt.Println(item)
			break
		}

		list, found := store.FindList(name)

		if found {
			fmt.Println(list)
			break
		}

		store.CreateList(name)
		fmt.Printf("Boom! Created a new list called %s\n", name)
	case len(store.Lists) == 0:
		showEmptyMessage()
	default:
		store.PrintSummary()
	}

	err = backend.Save(store)

	if err != nil {
		log.Fatal(err)
	}
}
