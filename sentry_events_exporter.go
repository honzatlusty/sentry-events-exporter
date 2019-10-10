package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
        "os"
	"os/exec"
)

type OptsHandler struct {
	Token string
	URL   string
	Hours string
}

func (oh *OptsHandler) handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "#Sentry Event Exporter for Prometheus V.1.0 : Path : %s \n \n", r.URL.Path[1:])
	out, err := exec.Command("/bin/bash", "sentry_events.sh", oh.Token, oh.URL, oh.Hours).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", out)
}

func usage() {
	fmt.Print("Example usage: ./exporter -token <sentry api token> -url <sentry url> -hours <number of the past hours to search for events from>\n" +
		"<token> - token for Sentry API                   - Mandatory\n" +
		"<url>   - Sentry url                             - Optional - defaults to \"http://localhost\"\n" +
		"<hours> - how many hours to look back for events - Optional - defaults to 1\n")
}

func main() {
	token := flag.String("token", "", "token for Sentry API")
	url := flag.String("url", "http://localhost", "Sentry url")
	hoursAgo := flag.String("hours", "1", "how many hours to look back for events")
	flag.Parse()
	if *token == "" {
                usage()
                os.Exit(2)
        }
	myFilesHandler := &OptsHandler{Token: *token, URL: *url, Hours: *hoursAgo}
	http.HandleFunc("/metrics", myFilesHandler.handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

/*Based on https://medium.com/@jomzsg/prometheus-wildlfy-exporter-using-shell-script-and-golang-cf7ec4179c0b*/
