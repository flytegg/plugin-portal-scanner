package platforms

import (
	"github.com/go-resty/resty/v2"
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

	var response Response

	_, err := client.R().
		EnableTrace().
		SetResult(&response).
		SetQueryParams(map[string]string{
			"facets": `[["project_type:plugin"],["categories:bukkit","categories:spigot","categories:paper","categories:purpur","categories:folia"]]`,
			"offset": "0",
			"limit":  "100",
		}).
		Get("https://api.modrinth.com/v2/search")

	if err != nil {
		return nil, err
	}

	// Convert to Resource type
	resources := make([]Resource, len(response.Hits))
	for i, hit := range response.Hits {
		resources[i] = Resource{
			Name: hit.Title,
			ID:   hit.Slug,
		}
	}

	return resources, nil
}
