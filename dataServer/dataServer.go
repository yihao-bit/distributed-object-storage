package dataServer

import (
	"git_test/dataServer/heartbeat"
	"git_test/dataServer/locate"
	"git_test/dataServer/objects"
	"git_test/dataServer/temp"
	"log"
	"net/http"
	"os"
)

func main() {
	locate.CollectObjects()
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
