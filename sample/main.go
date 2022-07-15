package main

import(
        "log"
        "net/http"
)

var servicePort = ":8080"

func main() {
    log.SetPrefix("[LOG] ")
    log.SetFlags(3)

    log.Println("SampleService started successfully on port ", servicePort)

    fs := http.FileServer(http.Dir("./"))

    log.Fatal(http.ListenAndServe(servicePort, fs))
}
