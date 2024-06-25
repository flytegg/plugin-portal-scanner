package platforms

import (
	"github.com/go-resty/resty/v2"
	"strconv"
)

type SpigotResource struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func SpigotRequests() ([]Resource, error) {
	client := resty.New()

	var spigotResources []SpigotResource

	_, err := client.R().
		EnableTrace().
		SetDebug(true).
		SetResult(&spigotResources).
		SetQueryParams(map[string]string{
			"size":   "500",
			"sort":   "-downloads",
			"fields": "id,name",
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
