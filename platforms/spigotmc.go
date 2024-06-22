package platforms

import (
    "encoding/json"
    "io"
    "net/http"
)

func SpigotRequests() ([]Resource, error) {
    get, err := http.Get("https://api.spiget.org/v2/resources/free?size=50&sort=-downloads&fields=id%2Cname")
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
