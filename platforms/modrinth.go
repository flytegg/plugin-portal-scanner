package platforms

import (
	"github.com/go-resty/resty/v2"
	"strconv"
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

	for {
		var response Response

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

	return allResources, nil
}
