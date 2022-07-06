package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

const permission = 0644

type Arguments map[string]string

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func Perform(args Arguments, writer io.Writer) error {
	err := validateArgs(args)
	if err != nil {
		return err
	}

	var fileName = args["fileName"]
	var operation = args["operation"]

	switch operation {
	case "add":
		var newUser User

		var item = args["item"]

		err = json.Unmarshal([]byte(item), &newUser)
		if err != nil {
			return err
		}

		users, err := getUserList(fileName)
		if err != nil {
			return err
		}

		for _, item := range users {
			if item.Id == newUser.Id {
				writer.Write([]byte(fmt.Sprintf("Item with id %v already exists", newUser.Id)))
			}
		}

		users = append(users, newUser)

		err = saveUserList(fileName, users)
		if err != nil {
			return err
		}
	case "list":
		var users []User
		var data []byte

		data, err = readFile(fileName)
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, &users)
		if err != nil {
			return err
		}

		if len(users) > 0 {
			writer.Write(data)
		}
	case "findById":
		var user User
		var data []byte
		var userWasFound = false

		var id = args["id"]

		users, err := getUserList(fileName)
		if err != nil {
			return err
		}

		for _, item := range users {
			if item.Id == id {
				userWasFound = true
				user = item
				break
			}
		}

		if userWasFound {
			data, err = json.Marshal(user)
			if err != nil {
				return err
			}
			writer.Write(data)
		} else {
			writer.Write([]byte(""))
		}
	case "remove":
		var userToRemoveIndex = -1
		var userWasFound = false

		var id = args["id"]

		users, err := getUserList(fileName)
		if err != nil {
			return err
		}

		for index, item := range users {
			if item.Id == id {
				userToRemoveIndex = index
				userWasFound = true
				break
			}
		}

		if userWasFound {
			users = append(users[:userToRemoveIndex], users[userToRemoveIndex+1:]...)

			err = saveUserList(fileName, users)
			if err != nil {
				return err
			}
		} else {
			writer.Write([]byte(fmt.Sprintf("Item with id %v not found", id)))
		}
	}

	return nil
}

func parseArgs() Arguments {
	id := flag.String("id", "", "user id to be found in the users list")
	item := flag.String("item", "", "valid json object with the id, email and age fields")
	operation := flag.String("operation", "", "accepts these types of operation: add list findById remove")
	fileName := flag.String("fileName", "", "file name to store users list in the json format")

	flag.Parse()

	var result = Arguments{
		"id":        *id,
		"item":      *item,
		"operation": *operation,
		"fileName":  *fileName,
	}

	return result
}

func validateArgs(args Arguments) error {
	var fileName = args["fileName"]
	if fileName == "" {
		return errors.New("-fileName flag has to be specified")
	}

	var operation = args["operation"]
	switch operation {
	case "":
		return errors.New("-operation flag has to be specified")
	case "add":
		var item = args["item"]
		if item == "" {
			return errors.New("-item flag has to be specified")
		}
	case "list":
		return nil
	case "findById", "remove":
		var id = args["id"]
		if id == "" {
			return errors.New("-id flag has to be specified")
		}
	default:
		return fmt.Errorf("Operation %s not allowed!", operation)
	}

	return nil
}

func readFile(fileName string) ([]byte, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, permission)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func writeFile(fileName string, data []byte) error {
	return os.WriteFile(fileName, data, permission)
}

func getUserList(fileName string) ([]User, error) {
	var users []User
	var data []byte

	data, err := readFile(fileName)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return users, nil
	}

	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func saveUserList(fileName string, users []User) error {
	var data []byte

	data, err := json.Marshal(users)
	if err != nil {
		return err
	}

	err = writeFile(fileName, data)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
