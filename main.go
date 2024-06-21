package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "io"
    "log"
    "net/http"
    "os"
    "plugin-portal-scanner/platforms"
    "time"
)

func main() {
    loadEnv()
    authToken := os.Getenv("AUTH_TOKEN")

    resources, err := platforms.SpigotRequests()
    if err != nil {
        log.Fatal(err)
    }

    for _, resource := range resources {
        err := postPluginData(resource.ID, authToken)
        if err != nil {
            log.Printf("Error posting data for resource ID %d: %v\n", resource.ID, err)
        } else {
            fmt.Printf("Successfully posted data for resource ID %d\n", resource.ID)
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

func postPluginData(id int, authToken string) error {
    url := fmt.Sprintf("https://api.pluginportal.link/v1/plugins/spigotmc?id=%d", id)
    req, err := http.NewRequest("POST", url, nil)
    if err != nil {
        return err
    }

    req.Header.Set("Authorization", "Bearer "+authToken)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        bodyBytes, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("error: %s", string(bodyBytes))
    }

    return nil
}
