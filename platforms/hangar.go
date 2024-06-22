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
    get, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer get.Body.Close()

    bodyBytes, err := io.ReadAll(get.Body)
    if err != nil {
        return nil, err
    }

    var response ApiResponse
    err = json.Unmarshal(bodyBytes, &response)
    if err != nil {
        return nil, err
    }

    return response.Result, nil
}

func HangarRequests() ([]Resource, error) {
    var allProjects []Project

    // Fetch the first set of resources
    projects, err := fetchResources(0)
    if err != nil {
        return nil, err
    }
    allProjects = append(allProjects, projects...)

    // Fetch the second set of resources
    projects, err = fetchResources(25)
    if err != nil {
        return nil, err
    }
    allProjects = append(allProjects, projects...)

    // Convert to Resource type
    resources := make([]Resource, len(allProjects))
    for i, project := range allProjects {
        resources[i] = Resource{
            Name: project.Namespace.Owner,
            ID:   project.Namespace.Slug,
        }
    }

    return resources, nil
}
