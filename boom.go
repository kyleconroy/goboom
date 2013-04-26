package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"os/user"
	"path"
	"runtime"
)

type Config struct {
	Backend string
}

type Store map[string]map[string]string

type Backend interface {
	Save(Store) error
	Fetch() (Store, error)
}

// In-memory store used for testing puposes
type InMemoryBackend struct {
	Db Store
}

func (b *InMemoryBackend) Fetch() (Store, error) {
	return b.Db, nil
}

func (b *InMemoryBackend) Save(store Store) error {
	b.Db = store
	return nil
}

type JsonBackend struct {
	Filename string
}

func (jb *JsonBackend) Fetch() (Store, error) {
	filebyte, err := ioutil.ReadFile(jb.Filename)

	if err != nil {
		return Store{}, nil
	}

	var db Store

	err = json.Unmarshal(filebyte, &db)

	if err != nil {
		return Store{}, err
	}

	return db, nil
}

func (jb *JsonBackend) Save(store Store) error {
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
	usr, err := user.Current()
	if err != nil {
		return &JsonBackend{}, err
	}
	backend := JsonBackend{Filename: path.Join(usr.HomeDir, ".boom.json")}
	return &backend, nil
}

func copyToClipboard(key string, contents string) error {
	cmd := exec.Command("xclip", "-selection", "clipboard")

	if runtime.GOOS == "darwin" {
		cmd = exec.Command("pbcopy", contents)
	} else if runtime.GOOS == "windows" {
		cmd = exec.Command("clip", contents)
	}

	w, err := cmd.StdinPipe()

	if err != nil {
		return err
	}

	_, err = w.Write([]byte(contents))

	if err != nil {
		return err
	}

	w.Close()

	fmt.Printf("Boom! We just copied %s to your clipboard.\n", key)
	return cmd.Run()
}

type Runner struct {
	storage Store
	backend Backend
}

func (c *Runner) All() error {
	for name, values := range c.storage {
		fmt.Println("  " + name)
		for key, value := range values {
			fmt.Printf("    %s: %s\n", key, value)
		}
	}
	return nil
}

func (c *Runner) Overview() error {
	for name, values := range c.storage {
		fmt.Printf("  %v (%d)\n", name, len(values))
	}

	if len(c.storage) > 0 {
		return nil
	}

	fmt.Println("You don't have anything yet! To start out, create a new list:")
	fmt.Println("  $ boom <list-name>")
	fmt.Println("And then add something to your list!")
	fmt.Println("  $ boom <list-name> <item-name> <item-value>")
	fmt.Println("You can then grab your new item:")
	fmt.Println("  $ boom <item-name>")
	return nil
}

func (c *Runner) Execute() error {
	flag.Parse()

	backend, err := load(Config{})

	if err != nil {
		return err
	}

	store, err := backend.Fetch()

	if err != nil {
		return err
	}

	c.backend = backend
	c.storage = store

	command := flag.Arg(0)
	major := flag.Arg(1)
	minor := flag.Arg(2)

	if command == "" {
		return c.Overview()
	}

	return c.Delegate(command, major, minor)
}

func (c *Runner) Help() error {
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
	_, err := fmt.Println(help)
	return err
}

func (c *Runner) Edit() error {
	return nil
}

func (c *Runner) Switch(backend string) error {
	return nil
}

func (c *Runner) ShowStorage() error {
	return nil
}

func (c *Runner) ShowVersion() error {
	_, err := fmt.Println("0.0.1")
	return err
}

func (c *Runner) DetailList(name string) error {
	values, ok := c.storage[name]

	if !ok {
		return errors.New("Unknown list " + name)
	}

	for key, value := range values {
		fmt.Printf("    %s: %s\n", key, value)
	}

	return nil
}

func (c *Runner) DeleteItem(name string, key string) error {
	fmt.Println("Boom! " + key + " is gone forever")

	if values, ok := c.storage[name]; ok {
		if _, ok = values[key]; ok {
			delete(values, name)
			return c.Save()
		}
	}

	return nil
}

func (c *Runner) DeleteList(name string) error {
	fmt.Printf("You sure you want to delete everything in %s? (y/n):\n", name)

	var answer string
	fmt.Scanf("%s", &answer)

	if answer == "yes" || answer == "y" {
		delete(c.storage, name)
		fmt.Printf("Boom! Deleted all your %s.\n", name)
		return c.Save()
	}

	return nil
}

func (c *Runner) ListExists(name string) bool {
	_, ok := c.storage[name]
	return ok
}

func (c *Runner) AddList(name string) error {
	if _, ok := c.storage[name]; ok {
		return nil
	}

	fmt.Printf("Boom! Created a new list called %s.\n", name)
	c.storage[name] = make(map[string]string)
	return c.Save()
}

func (c *Runner) AddItem(name string, key string, value string) error {
	if values, ok := c.storage[name]; !ok {
		fmt.Printf("Boom! Created a new list called %s.\n", name)
		values = make(map[string]string)
		c.storage[name] = values
	}

	fmt.Printf("Boom! %s in %s is %s. Got it\n", key, name, value)

	c.storage[name][key] = value
	c.Save()

	return nil
}

func (c *Runner) ItemExists(key string) bool {
	for _, values := range c.storage {
		if _, ok := values[key]; ok {
			return true
		}
	}
	return false
}

func (c *Runner) SearchItem(key string) error {
	for _, values := range c.storage {
		if value, ok := values[key]; ok {
			return copyToClipboard(key, value)
		}
	}
	return errors.New("Couldn't find key: " + key)
}

func (c *Runner) SearchListItem(name string, key string) error {
	values, ok := c.storage[name]

	if !ok {
		return errors.New("Unknown list " + name)
	}

	value, ok := values[key]

	if !ok {
		return errors.New("Unknown key " + key)
	}

	return copyToClipboard(key, value)
}

func (c *Runner) Save() error {
	return c.backend.Save(c.storage)
}

func (c *Runner) EchoListItem(listName string, itemName string) error {
	values, ok := c.storage[listName]
	value, ok := values[itemName]

	if !ok {
		return errors.New("Unknown key " + itemName)
	}

	fmt.Println(value)
	return nil
}

func (c *Runner) EchoItem(itemName string) error {
	for _, values := range c.storage {
		if value, ok := values[itemName]; ok {
			fmt.Println(value)
			return nil
		}
	}
	return errors.New("Couldn't find key: " + itemName)
}

func (c *Runner) Delegate(command string, major string, minor string) error {
	switch {
	case command == "all":
		return c.All()
	case command == "edit":
		return c.Edit()
	case command == "switch":
		return c.Switch(major)
	case command == "storage":
		return c.ShowStorage()
	case command == "version":
		return c.ShowVersion()
	case command == "help":
		return c.Help()
	case command == "echo":
		switch {
		case c.ListExists(major):
			return c.EchoListItem(major, minor)
		default:
			return c.EchoItem(major)
		}
	case c.ListExists(command):
		switch {
		case major == "delete":
			return c.DeleteList(command)
		case minor == "delete":
			return c.DeleteItem(command, major)
		case len(minor) > 0:
			return c.AddItem(command, major, minor)
		case len(major) > 0:
			return c.SearchListItem(command, major)
		default:
			return c.DetailList(command)
		}
	case c.ItemExists(command):
		return c.SearchItem(command)
	default:
		return c.AddList(command)
	}

	return errors.New("Command not recognized")
}

func main() {
	runner := Runner{}
	err := runner.Execute()

	if err != nil {
		log.Fatal(err)
	}
}
