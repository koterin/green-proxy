package main

import(
        "log"
        "net/http"
)

func main() {
    log.SetPrefix("[LOG] ")
    log.SetFlags(3)

    log.Println("supernet started successfully on port 8080")

    fs := http.FileServer(http.Dir("./"))

    log.Fatal(http.ListenAndServe(":8080", fs))
}
