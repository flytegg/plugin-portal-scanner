package platforms

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type Namespace struct {
	Owner string `json:"owner"`
	Slug  string `json:"slug"`
}

type Project struct {
	Namespace Namespace `json:"namespace"`
}

type ApiResponse struct {
	Result []Project `json:"result"`
}

func fetchResources(offset int) ([]Project, error) {
	url := "https://hangar.papermc.io/api/v1/projects?limit=25&offset=" + strconv.Itoa(offset) + "&sort=-downloads"
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(bodyBytes, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse.Result, nil
}

func HangarRequests() ([]Resource, error) {
	offsets := []int{0, 50, 75, 100, 125, 150}
	var allProjects []Project

	for _, offset := range offsets {
		projects, err := fetchResources(offset)
		if err != nil {
			return nil, err
		}
		allProjects = append(allProjects, projects...)
	}

	resources := make([]Resource, len(allProjects))
	for i, project := range allProjects {
		resources[i] = Resource{
			Name: project.Namespace.Owner,
			ID:   project.Namespace.Slug,
		}
	}

	return resources, nil
}
