package main

import (
    "fmt"
    "github.com/go-resty/resty/v2"
    "github.com/joho/godotenv"
    "log"
    "net/http"
    "os"
    "plugin-portal-scanner/platforms"
    "time"
)

func main() {
    loadEnv()

    registerModrinth()

    //for _, resource := range resources {
    //    err := postPluginData(resource.ID, authToken, "hangar")
    //    if err != nil {
    //        log.Printf("Error posting data for resource ID %s: %v\n", resource.ID, err)
    //    } else {
    //        fmt.Printf("Successfully posted data for resource ID %s\n", resource.ID)
    //    }
    //    time.Sleep(100 * time.Millisecond)
    //}

}

func registerHangar() {
    authToken := os.Getenv("AUTH_TOKEN")

    resources, err := platforms.HangarRequests()
    if err != nil {
        log.Fatal(err)
    }

    for _, resource := range resources {
        err := postPluginData(resource.ID, authToken, "hangar")
        if err != nil {
            log.Printf("Error posting data for resource ID %s: %v\n", resource.ID, err)
        } else {
            fmt.Printf("Successfully posted data for resource ID %s\n", resource.ID)
        }
        time.Sleep(100 * time.Millisecond)
    }
}

func registerSpigotMC() {
    authToken := os.Getenv("AUTH_TOKEN")

    resources, err := platforms.SpigotRequests()
    if err != nil {
        log.Fatal(err)
    }

    for _, resource := range resources {
        err := postPluginData(resource.ID, authToken, "spigotmc")
        if err != nil {
            log.Printf("Error posting data for resource ID %s: %v\n", resource.ID, err)
        } else {
            fmt.Printf("Successfully posted data for resource ID %s\n", resource.ID)
        }
        time.Sleep(100 * time.Millisecond)
    }
}

func registerModrinth() {
    authToken := os.Getenv("AUTH_TOKEN")

    resources, err := platforms.ModrinthRequests()
    if err != nil {
        log.Fatal(err)
    }

    for _, resource := range resources {
        err := postPluginData(resource.ID, authToken, "modrinth")
        if err != nil {
            log.Printf("Error posting data for resource ID %s: %v\n", resource.ID, err)
        } else {
            fmt.Printf("Successfully posted data for resource ID %s\n", resource.ID)
        }
        time.Sleep(100 * time.Millisecond)
    }
}

func loadEnv() {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

func postPluginData(id string, authToken string, platformString string) error {
    url := fmt.Sprintf("https://api.pluginportal.link/v1/plugins/%s/%s", platformString, id)

    client := resty.New().
        EnableTrace()

    resp, err := client.R().
        SetHeader("Authorization", "Bearer "+authToken).
        Post(url)

    if err != nil {
        return err
    }

    if resp.StatusCode() != http.StatusOK {
        return fmt.Errorf("error: %s", resp.String())
    }

    return nil
}
