package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"plugin-portal-scanner/platforms"
	"strconv"
	"sync"
	"time"
)

func main() {
	loadEnv()

	var wg sync.WaitGroup
	wg.Add(1) // Number of goroutines

	// go func() {
	// 	defer wg.Done()
	// 	registerSpigotMC()
	// }()

	go func() {
	   defer wg.Done()
	   registerModrinth()
	}()
	
	// go func() {
	//    defer wg.Done()
	//    registerHangar()
	// }()

	wg.Wait() // Wait for all goroutines to finish
	log.Println("All registrations completed")
}

func registerHangar() {
	authToken := os.Getenv("AUTH_TOKEN")

	resources, err := platforms.HangarRequests()
	if err != nil {
		log.Fatal(err)
	}

	for _, resource := range resources {
		err := postPluginDataWithRetry(resource.ID, authToken, "hangar")
		if err != nil {
			log.Printf("Error posting data for resource ID %s: %v\n", resource.ID, err)
		} else {
			fmt.Printf("Successfully posted data for resource ID %s\n", resource.ID)
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func registerSpigotMC() {
	authToken := os.Getenv("AUTH_TOKEN")

	resources, err := platforms.SpigotRequests()
	if err != nil {
		log.Fatal(err)
	}

	filteredResources := make([]platforms.Resource, 0)
	for _, resource := range resources {
		id, err := strconv.Atoi(resource.ID)
		if err != nil {
			log.Printf("Error converting resource ID %s to int: %v\n", resource.ID, err)
			continue
		}
		if id >= 0 { // Temp set to 0 in case we want to change
			filteredResources = append(filteredResources, resource)
		}
	}

	for _, resource := range filteredResources {
		go func(resource platforms.Resource) {
			err := postPluginDataWithRetry(resource.ID, authToken, "spigotmc")
			if err != nil {
				log.Printf("Error posting data for resource ID %s: %v\n", resource.ID, err)
			} else {
				fmt.Printf("Successfully posted data for resource ID %s\n", resource.ID)
			}
		}(resource)
		time.Sleep(400 * time.Millisecond)
	}
}

func registerModrinth() {
	authToken := os.Getenv("AUTH_TOKEN")

	resources, err := platforms.ModrinthRequests()
	if err != nil {
		log.Fatal(err)
	}

	total := len(resources)
	successful := 0
	failed := 0
	startTime := time.Now()

	log.Printf("Starting to register %d Modrinth plugins...", total)

	for i, resource := range resources {
		err := postPluginDataWithRetry(resource.ID, authToken, "modrinth")
		if err != nil {
			log.Printf("Error posting data for resource ID %s: %v\n", resource.ID, err)
			failed++
		} else {
			successful++
			if i%50 == 0 { // Progress update every 50 plugins
				log.Printf("Progress: %d/%d Modrinth plugins processed (%.1f%%)", i+1, total, float64(i+1)/float64(total)*100)
			}
		}
		time.Sleep(250 * time.Millisecond)
	}

	duration := time.Since(startTime)
	log.Printf("\nModrinth registration completed:")
	log.Printf("- Total plugins processed: %d", total)
	log.Printf("- Successful registrations: %d", successful)
	log.Printf("- Failed registrations: %d", failed)
	log.Printf("- Success rate: %.1f%%", float64(successful)/float64(total)*100)
	log.Printf("- Time taken: %v", duration)
	log.Printf("- Average processing rate: %.2f plugins/minute", float64(total)/(duration.Minutes()))
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func postPluginDataWithRetry(id string, authToken string, platformString string) error {
	url := fmt.Sprintf("https://api.pluginportal.link/v1/plugins/%s/%s", platformString, id)
	client := resty.New().EnableTrace()

	// First attempt
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+authToken).
		Post(url)
	if err == nil && resp.StatusCode() == http.StatusOK {
		return nil
	}

	// Log first failure
	log.Printf("First attempt failed for resource ID %s: %v\n", id, err)

	// Second attempt
	resp, err = client.R().
		SetHeader("Authorization", "Bearer "+authToken).
		Post(url)
	if err == nil && resp.StatusCode() == http.StatusOK {
		return nil
	}

	// Return error if the second attempt fails
	if err != nil {
		return fmt.Errorf("second attempt failed for resource ID %s: %v", id, err)
	}
	return fmt.Errorf("second attempt failed for resource ID %s: %s", id, resp.String())
}
