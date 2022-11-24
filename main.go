package main

import (
    "context"
    "database/sql"
    "fmt"
    "os"
    "strconv"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/feature/rds/auth"
    _ "github.com/lib/pq"
)

func getenv(key, fallback string) string {
    value := os.Getenv(key)
    if len(value) == 0 {
        return fallback
    }
    return value
}

func main() {
    var dbName string = getenv("DB", "empty")
    var dbUser string = getenv("USER", "empty")
    var dbHost string = getenv("HOST", "empty")
    var port string = getenv("PORT", "5432")
    var dbPort,_ = strconv.Atoi(port)
    var region string = getenv("REGION", "eu-central-1")
    var dbEndpoint string = fmt.Sprintf("%s:%d", dbHost, dbPort)

    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        panic("configuration error: " + err.Error())
    }

    authenticationToken, err := auth.BuildAuthToken(
        context.TODO(), dbEndpoint, region, dbUser, cfg.Credentials)
    if err != nil {
        //panic("failed to create authentication token: " + err.Error())
        authenticationToken = ""
    }

    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
        dbHost, dbPort, dbUser, authenticationToken, dbName,
    )

    fmt.Println("DEBUG: ", dsn)

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        panic(err)
    } else {
        fmt.Println("Connected")
    }

    err = db.Ping()
    if err != nil {
        panic(err)
    } else {
        fmt.Println("Ping response")
    }
}
