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
	get()
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
		Name:     "test",
		Email:    "test@gmail.com",
		Password: "01234567890",
		Phone:    "4799909090990",
		IsAdm:    true,
	}

	log.Printf("Adding user: %v", userRequest)

	jsonData, _ := json.Marshal(userRequest)

	//if err != nil {
	//	log.Fatalf("Failed to marshal json: %v", err)
	//}

	req, _ := http.NewRequest("POST", "http://localhost:8080/users/add", bytes.NewBuffer(jsonData))

	//if err != nil {
	//	log.Fatalf("Failed to create request: %v", err)
	//}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, _ := client.Do(req)

	//if err != nil {
	//	log.Fatalf("Failed to make request: %v", err)
	//}

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

	//if resp.StatusCode != http.StatusOK {
	//	log.Fatalf("Expected status 200 OK, got %v", resp.StatusCode)
	//}

	log.Printf("Response: %v\n", resp)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	fmt.Println("Response body:")
	fmt.Println(string(body))
}

func delete() {
	req, err := http.NewRequest("DELETE", "http://localhost:8080/users/delete?id=1", nil)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	fmt.Println("Response body:")
	fmt.Println(string(body))
}
