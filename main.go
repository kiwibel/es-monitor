package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kiwibel/es-monitor/metric"
	"github.com/kiwibel/es-monitor/slack"
)

var slackChannel = os.Getenv("SLACK_CHANNEL")
var healthEndpoint = "/_cluster/health/"

type status struct {
	Cluster string `json:"cluster_name"`
	Status  string `json:"status"`
}

type cluster struct {
	Name     string
	URL      string
	Password string
}

var clusters = []cluster{
	{"Monitoring",
		fmt.Sprintf("%s%s", os.Getenv("ES_URL_MON"), healthEndpoint),
		os.Getenv("ES_PASSWORD_MON"),
	},

	{"Development",
		fmt.Sprintf("%s%s", os.Getenv("ES_URL_DEV"), healthEndpoint),
		os.Getenv("ES_PASSWORD_DEV"),
	},

	{"Staging",
		fmt.Sprintf("%s%s", os.Getenv("ES_URL_STG"), healthEndpoint),
		os.Getenv("ES_PASSWORD_STG"),
	},

	{"Production",
		fmt.Sprintf("%s%s", os.Getenv("ES_URL_PROD"), healthEndpoint),
		os.Getenv("ES_PASSWORD_PROD"),
	},
}

var emojiMap = map[string]string{
	"green":  ":white_check_mark:",
	"yellow": ":warning:",
	"red":    ":x:",
}

func main() {
	pollClusters()
	// Can wrap it into infinite loop
	// time.Sleep(time.Second * pollingInterval)
}

func initRequest(c cluster) (*http.Request, http.Client) {
	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}
	req, err := http.NewRequest(http.MethodGet, c.URL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("elastic", c.Password)
	req.Header.Set("User-Agent", "elastic-monitor")
	return req, spaceClient

}

func notifySlack(currentCluster, currentStatus string) {
	//	currentStatus = "yellow"
	slack.SendMessageToChannel(slackChannel,
		fmt.Sprintf("The status of the cluster `%s` is %s %s\n---------------\n",
			currentCluster, currentStatus, emojiMap[currentStatus]),
	)
}

func pollClusters() string {

	// Checks for flags if any
	// esUrlFlag = flag.String(name string, value string, usage string)

	for _, cl := range clusters {
		fmt.Printf("Polling cluster state for %s cluster\n", cl.Name)
		statusCode := 1.0 // 1 is green, 0 is non-green

		req, spaceClient := initRequest(cl)
		res, getErr := spaceClient.Do(req)
		if getErr != nil {
			log.Fatal(getErr)
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		body, readErr := ioutil.ReadAll(res.Body)

		//		Enable for troubleshooting
		//		fmt.Println(string(body))

		if readErr != nil {
			log.Fatal(readErr)
		}

		clusterHealth := status{}
		jsonErr := json.Unmarshal(body, &clusterHealth)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		fmt.Printf("Cluster: %s\nStatus: %s\n", cl.Name, clusterHealth.Status)

		if clusterHealth.Status != "green" {
			statusCode = 0.0
			fmt.Printf("Sending alert to Slack channel %s", slackChannel)
			notifySlack(cl.Name, clusterHealth.Status)
		}

		if metric.PutMetric(cl.Name, statusCode) != nil {
			return "Couldn't emit metrics to Cloudwatch"
		}

		defer fmt.Printf("Successfully checked status for %s cluster and sent metrics to Cloudwatch.\n", cl.Name)
	}

	return "Exiting OK"
}
