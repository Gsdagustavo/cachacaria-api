package main

import (
	"bytes"
	"cachacariaapi/internal/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// TODO: add CLI tools for debugging
func main() {
	//add()
	//findById()
	//delete()
	//update()
	//get()

	register()
	//login()
}

func register() {
	userRequest := models.UserRequest{
		Email:    "test@gmail.com",
		Password: "01234567890",
		Phone:    "4799909090990",
		IsAdm:    true,
	}

	log.Printf("Login: %v", userRequest)

	jsonData, _ := json.Marshal(userRequest)

	req, _ := http.NewRequest("POST", "http://localhost:8080/auth/register", bytes.NewBuffer(jsonData))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, _ := client.Do(req)

	bd, _ := ioutil.ReadAll(resp.Body)

	log.Printf("Status code: %v", resp.StatusCode)
	log.Printf("Response body: %v", string(bd))
}

func login() {
	userRequest := models.UserRequest{
		Email:    "test@gmail.com",
		Password: "01234567890",
		Phone:    "4799909090990",
		IsAdm:    true,
	}

	log.Printf("Login: %v", userRequest)

	jsonData, _ := json.Marshal(userRequest)

	req, _ := http.NewRequest("POST", "http://localhost:8080/auth/login", bytes.NewBuffer(jsonData))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, _ := client.Do(req)

	bd, _ := ioutil.ReadAll(resp.Body)

	log.Printf("Status code: %v", resp.StatusCode)
	log.Printf("Response body: %v", string(bd))
}

func findById() {
	res, err := http.Get("http://localhost:8080/users/id?id=5")

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("user: %v", user)
}

func add() {
	userRequest := models.UserRequest{
		Email:    "test@gmail.com",
		Password: "01234567890",
		Phone:    "4799909090990",
		IsAdm:    true,
	}

	log.Printf("Adding user: %v", userRequest)

	jsonData, _ := json.Marshal(userRequest)

	req, _ := http.NewRequest("POST", "http://localhost:8080/users/add", bytes.NewBuffer(jsonData))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, _ := client.Do(req)

	bd, _ := ioutil.ReadAll(resp.Body)

	log.Printf("Status code: %v", resp.StatusCode)
	log.Printf("Response body: %v", string(bd))
}

func get() {
	resp, err := http.Get("http://localhost:8080/users")
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("Response: %v\n", resp)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	fmt.Println("Response body:")
	fmt.Println(string(body))
}

func delete() {
	req, _ := http.NewRequest("DELETE", "http://localhost:8080/users/delete?id=2", nil)

	client := &http.Client{}

	resp, _ := client.Do(req)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("Response body:")
	fmt.Println(string(body))
}

func update() {
	userRequest := models.UserRequest{
		Email: "@gmail.com",
		//Password: "eqweqewqewqewqeq",
		//Phone:    "479990dadw9090990",
		//IsAdm:    true,
	}

	jsonData, _ := json.Marshal(userRequest)

	req, _ := http.NewRequest("PUT", "http://localhost:8080/users/update?id=7", bytes.NewBuffer(jsonData))

	log.Printf("Updating user: %v", userRequest)
	log.Printf("Request: %v", req)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, _ := client.Do(req)

	bd, _ := ioutil.ReadAll(resp.Body)

	log.Printf("Status code: %v", resp.StatusCode)
	log.Printf("Response body: %v", string(bd))
}
