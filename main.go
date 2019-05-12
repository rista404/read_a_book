package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	. "github.com/logrusorgru/aurora"
)

const folder = "/.config/read_a_book"
const configFileName = "config.json"

func getConfigDir() string {
	// Get the user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Can't find user home dir.")
		os.Exit(2)
	}

	return homeDir + folder
}

func getConfigPath() string {
	dir := getConfigDir()

	return dir + "/" + configFileName
}

type Config struct {
	from uint8
	to   uint8
	book Book
}

type Book struct {
	Name       string `json:"name"`
	Author     string `json:"author"`
	Percentage uint8  `json:"percentage"`
}

func (b Book) String() string {
	author := fmt.Sprintf("by %s", b.Author)
	percent := fmt.Sprintf("[%d%%]", b.Percentage)
	return fmt.Sprintf("%s %s %s %s", Gray(15, "Go read"), b.Name, Gray(15, author), Green(percent))
}

func (b *Book) Save() {
	// Marshal into JSON
	json, err := json.Marshal(&b)
	if err != nil {
		log.Fatal(err)
	}

	// Make the config directory
	_ = os.MkdirAll(getConfigDir(), os.ModePerm)

	// Open the file for creating
	f, err := os.Create(getConfigPath())
	defer f.Close()

	if err != nil {
		fmt.Println("Error creating file: ", err)
		return
	}

	// Write the JSON to the file
	_, err = f.Write(json)
	if err != nil {
		fmt.Println("Error writing to the file: ", err)
		return
	}
}

func (b *Book) ReadFromFile() {
	filePath := getConfigPath()

	// Open the file for reading
	content, err := ioutil.ReadFile(filePath)

	if err != nil {
		fmt.Println("Error reading from file: ", err)
		return
	}

	err = json.Unmarshal(content, &b)

	if err != nil {
		fmt.Println("Error unmarshaling json: ", err)
		return
	}
}

func shouldDisplay() bool {
	now := time.Now()
	morning := 10
	evening := 19

	hour := now.Hour()
	return hour < morning || hour >= evening
}

func main() {
	if len(os.Args) > 1 {
		cmd := os.Args[1]

		// Sets the book in config file
		if cmd == "set" {
			b := Book{"", "", 0}
			scanner := bufio.NewScanner(os.Stdin)

			// Read the book name
			fmt.Print("Book name: ")
			scanner.Scan()
			b.Name = scanner.Text()

			// Read the book author
			fmt.Print("Book author: ")
			scanner.Scan()
			b.Author = scanner.Text()

			// Save in the file
			b.Save()
			fmt.Println("Book saved! Start reading 📚")
			return
		}

		// Change the read percentage of the book
		if n, err := strconv.Atoi(cmd); err == nil {
			if n < 0 || n > 100 {
				fmt.Println("Please provide a number between 0 and 100 to change the read percentage of the book.")
				return
			}

			b := new(Book)
			b.ReadFromFile()
			b.Percentage = uint8(n)
			b.Save()

			fmt.Println("Wohoo! Keep on reading 📖")

			return
		}

		fmt.Println("Unknown command")

		return
	}

	if shouldDisplay() {
		b := new(Book)
		b.ReadFromFile()
		if len(b.Name) == 0 {
			return
		}
		fmt.Print(b)
	}
}
