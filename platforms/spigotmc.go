package platforms

import (
    "encoding/json"
    "io"
    "net/http"
)

type Resource struct {
    Name string `json:"name"`
    ID   int    `json:"id"`
}

func SpigotRequests() ([]Resource, error) {
    get, err := http.Get("https://api.spiget.org/v2/resources/free?size=5&sort=-downloads&fields=id%2Cname")
    if err != nil {
        return nil, err
    }
    defer get.Body.Close()

    bodyBytes, err := io.ReadAll(get.Body)
    if err != nil {
        return nil, err
    }

    var resources []Resource
    err = json.Unmarshal(bodyBytes, &resources)
    if err != nil {
        return nil, err
    }

    return resources, nil
}
