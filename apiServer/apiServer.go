package apiServer

import (
	"git_test/apiServer/heartbeat"
	"git_test/apiServer/locate"
	"git_test/apiServer/objects"
	"git_test/apiServer/temp"
	"git_test/apiServer/versions"
	"log"
	"net/http"
	"os"
)

func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
