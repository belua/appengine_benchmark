package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func serviceAccountsJSON() {
	// Your credentials should be obtained from the Google
	// Developer Console (https://console.developers.google.com).
	// Navigate to your project, then see the "Credentials" page
	// under "APIs & Auth".
	// To create a service account client, click "Create new Client ID",
	// select "Service Account", and click "Create Client ID". A JSON
	// key file will then be downloaded to your computer.
	opts, err := oauth2.New(
		google.ServiceAccountJSONKey("clientsecret.json"),
		oauth2.Scope("https://www.googleapis.com/auth/cloud-taskqueue"),
	)
	if err != nil {
		log.Fatal(err)
	}
	// client := http.Client{Transport: &logTransport{opts.NewTransport()}}
	client := http.Client{Transport: opts.NewTransport()}
	// Initiate an http.Client. The following GET request will be
	// authorized and authenticated on the behalf of
	// your service account.
	// resp, err := client.Get("https://www.googleapis.com/taskqueue/v1beta2/projects/fiery-diorama-771/taskqueues/perf-pull-queue/tasks/lease/")
	resp, err := client.PostForm("https://www.googleapis.com/taskqueue/v1beta2/projects/fiery-diorama-771/taskqueues/perf-pull-queue/tasks/lease/", url.Values{"leaseSecs": {"10"}, "numTasks": {"100"}, "project": {"fiery-diorama-771"}, "taskqueue": {"perf-pull-queue"}})
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	println("Body: " + string(b))
}
