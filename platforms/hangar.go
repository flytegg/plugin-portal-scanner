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
	var allProjects []Project
	offset := 0
	
	for {
		projects, err := fetchResources(offset)
		if err != nil {
			return nil, err
		}
		
		// If no more projects are returned, break the loop
		if len(projects) == 0 {
			break
		}
		
		allProjects = append(allProjects, projects...)
		offset += 25 // Increment by the limit value used in fetchResources
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
