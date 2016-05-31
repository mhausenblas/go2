package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	VERSION            string = "0.1.0"
	GO2_PORT           int    = 6969
	DPID_PARAM         string = "dpid"
	MESOS_DNS_ENDPOINT string = "http://leader.mesos:8123"
)

var (
	mux *http.ServeMux
)

type SRVRecord struct {
	Service string
	Host    string
	IP      string
	Port    string
}

// Lookup takes a process dpid (distributed PID, that is, the ID of a Marathon app)
// and returns an endpoint for a task it serves on in the form (ip, port).
// For example:
//  lookup("/example/app") -> 1.2.3.4, 32787
func lookup(dpid string) (ip string, port string) {
	query := "_" + dpid + "._tcp.marathon.mesos."
	resp, qerr := http.Get(MESOS_DNS_ENDPOINT + "/v1/services/" + query)
	if qerr != nil {
		log.WithFields(log.Fields{"func": "lookup"}).Error("Can't look up ", dpid, " due to ", qerr)
		return
	}
	log.WithFields(log.Fields{"func": "lookup"}).Info("Got response from Mesos-DNS ", resp)
	defer resp.Body.Close()
	body, rerr := ioutil.ReadAll(resp.Body)
	if rerr != nil {
		log.WithFields(log.Fields{"func": "lookup"}).Error("Error reading response from Mesos-DNS due to ", rerr)
	}
	var srvrecords []SRVRecord
	merr := json.Unmarshal(body, &srvrecords)
	if merr != nil {
		log.WithFields(log.Fields{"func": "lookup"}).Error("Error unmarshalling SRV record due to ", merr)
	}
	ip = srvrecords[0].IP
	port = srvrecords[0].Port
	return ip, port
}

func main() {
	mux = http.NewServeMux()
	fmt.Printf("This is the Go2 in version %s listening on port %d\n", VERSION, GO2_PORT)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dpid := r.URL.Query().Get(DPID_PARAM) // extract the /?dpid=$DPID value
		log.WithFields(log.Fields{"handle": "/"}).Info("Got dpid ", dpid)
		ip, port := lookup(dpid)
		fmt.Fprint(w, "http://"+ip+":"+port)
	})
	p := strconv.Itoa(GO2_PORT)
	log.Fatal(http.ListenAndServe(":"+p, mux))
}
