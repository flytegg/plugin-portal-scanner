package platforms

import (
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
	"log"
)

type ModrinthProject struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

type Response struct {
	Hits []ModrinthProject `json:"hits"`
}

func ModrinthRequests() ([]Resource, error) {
	client := resty.New().
		SetHeaders(map[string]string{
			"User-Agent": "Plugin Portal Scanner (https://github.com/flytegg/plugin-portal)",
		})

	var allResources []Resource
	offset := 0
	limit := 100
	requestCount := 0
	startTime := time.Now()

	// Rate limiting: 250 requests per minute = 1 request per 240ms
	rateLimiter := time.NewTicker(240 * time.Millisecond)
	defer rateLimiter.Stop()

	log.Printf("Starting Modrinth plugin scan...")

	for {
		// Wait for rate limiter
		<-rateLimiter.C

		var response Response
		requestCount++

		log.Printf("Fetching Modrinth plugins - Offset: %d, Total found so far: %d", offset, len(allResources))

		_, err := client.R().
			EnableTrace().
			SetResult(&response).
			SetQueryParams(map[string]string{
				"facets": `[["project_type:plugin"],["categories:bukkit","categories:spigot","categories:paper","categories:purpur","categories:folia"]]`,
				"offset": strconv.Itoa(offset),
				"limit":  strconv.Itoa(limit),
			}).
			Get("https://api.modrinth.com/v2/search")

		if err != nil {
			log.Printf("Error fetching Modrinth plugins: %v", err)
			return nil, err
		}

		// If no more hits are returned, break the loop
		if len(response.Hits) == 0 {
			break
		}

		// Convert to Resource type
		for _, hit := range response.Hits {
			allResources = append(allResources, Resource{
				Name: hit.Title,
				ID:   hit.Slug,
			})
		}

		offset += limit
	}

	duration := time.Since(startTime)
	log.Printf("Modrinth scan completed:")
	log.Printf("- Total plugins found: %d", len(allResources))
	log.Printf("- Total requests made: %d", requestCount)
	log.Printf("- Time taken: %v", duration)
	log.Printf("- Average request rate: %.2f requests/minute", float64(requestCount)/(duration.Minutes()))

	return allResources, nil
}
