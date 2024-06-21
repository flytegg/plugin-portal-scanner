package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
)

type Resource struct {
    Name string `json:"name"`
    ID   int    `json:"id"`
}

func spigotRequests() {
    get, err := http.Get("https://api.spiget.org/v2/resources/free?size=5&sort=-downloads&fields=id%2Cname")
    if err != nil {
        log.Fatal(err)
    }
    defer get.Body.Close()

    bodyBytes, err := io.ReadAll(get.Body)
    if err != nil {
        log.Fatal(err)
    }

    var resources []Resource
    err = json.Unmarshal(bodyBytes, &resources)
    if err != nil {
        log.Fatal(err)
    }

    for _, resource := range resources {
        fmt.Printf("Name: %s, ID: %d\n", resource.Name, resource.ID)
    }
}
