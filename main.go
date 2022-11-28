package main

import (
    "database/sql"
    "fmt"
    "os"
    "strconv"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/rds/rdsutils"
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

    opts := session.Options{Config: aws.Config{
            CredentialsChainVerboseErrors: aws.Bool(true),
            Region:                        aws.String(region),
            MaxRetries:                    aws.Int(3),
        }}
    sess := session.Must(session.NewSessionWithOptions(opts))

    authenticationToken, err := rdsutils.BuildAuthToken(dbEndpoint, region, dbUser, sess.Config.Credentials)
    if err != nil {
        panic("failed to create authentication token: " + err.Error())
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
