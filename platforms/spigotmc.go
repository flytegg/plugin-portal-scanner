package platforms

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
)

type SpigotResource struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func SpigotRequests() ([]Resource, error) {
	var allResources []Resource

	// Set to /1 for now, prob will change.
	for i := 0; i < 55000/55000; i += 1 {
		fmt.Println("fetching")
		resources, err := fetchSpigotResources(i)
		if err != nil {
			fmt.Printf("Error fetching resources: %v\n", err)
			return nil, err
		}
		fmt.Printf("Done fetching")

		fmt.Println("Adding resources with offset", i)
		allResources = append(allResources, resources...)
	}

	return allResources, nil
}

func fetchSpigotResources(page int) ([]Resource, error) {
	client := resty.New()

	var spigotResources []SpigotResource

	_, err := client.R().
		EnableTrace().
		SetResult(&spigotResources).
		SetQueryParams(map[string]string{
			"size": "55000",
			//"sort":   "-downloads", // Commented due to timeout issues
			"fields": "name",
			"page":   strconv.Itoa(page),
		}).
		Get("https://api.spiget.org/v2/resources/free")
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, len(spigotResources))
	for i, hit := range spigotResources {
		resources[i] = Resource{
			Name: hit.Name,
			ID:   strconv.Itoa(hit.ID),
		}
	}

	return resources, nil
}
