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
	add()
	get()
	findById()
}

func findById() {
	res, err := http.Get("http://localhost:8080/users/id?id=3")

	if res.StatusCode != http.StatusOK || res.ContentLength == 0 || err != nil {
		log.Printf("error while trying to find by id. Status Code: %v", res.StatusCode)
		return
	}

	var user models.User
	err = json.NewDecoder(res.Body).Decode(&user)

	if err != nil {
		log.Printf("error while trying to decode user data: %v", err)
		return
	}

	prettyJSON, err := json.MarshalIndent(user, "", "  ")

	log.Printf("\nUser found: %v", string(prettyJSON))
}

func add() {
	userRequest := models.AddUserRequest{
		Name:     "awdadwdawd",
		Email:    "awdwaddwadad",
		Password: "asdadasdd",
		Phone:    "qweqweqweq",
		IsAdm:    true,
	}

	jsonData, err := json.Marshal(userRequest)

	if err != nil {
		log.Fatalf("Failed to marshal json: %v", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/users/add", bytes.NewBuffer(jsonData))

	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}

	fmt.Printf("Response: %v\n", resp)
}

func get() {
	resp, err := http.Get("http://localhost:8080/users")
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Expected status 200 OK, got %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	fmt.Println("Response body:")
	fmt.Println(string(body))
}
