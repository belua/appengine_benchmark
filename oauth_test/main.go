// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/gob"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"code.google.com/p/goauth2/oauth"
)

var config = &oauth.Config{
	ClientId:     "", // Set by --clientid or --clientid_file
	ClientSecret: "", // Set by --secret or --secret_file
	Scope:        "https://www.googleapis.com/auth/taskqueue",
	AuthURL:      "https://accounts.google.com/o/oauth2/auth",
	TokenURL:     "https://accounts.google.com/o/oauth2/token",
}

func main() {
	//	config.ClientId = fileContents("clientid.dat")
	//	config.ClientSecret = fileContents("clientsecret.dat")
	//	c := getOAuthClient(config)
	serviceAccountsJSON()
}

func tokenFromFile(file string) (*oauth.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := new(oauth.Token)
	err = gob.NewDecoder(f).Decode(t)
	return t, err
}

func getOAuthClient(config *oauth.Config) *http.Client {
	t := &oauth.Transport{
		Config:    config,
		Transport: &logTransport{http.DefaultTransport},
	}
	err := t.AuthenticateClient()
	if err != nil {
		log.Fatalf("Token exchange error: %v", err)
	}
	return t.Client()
}

func openUrl(url string) {
	try := []string{"xdg-open", "google-chrome", "open"}
	for _, bin := range try {
		err := exec.Command(bin, url).Run()
		if err == nil {
			return
		}
	}
	log.Printf("Error opening URL in browser.")
}

func fileContents(filename string) []byte {
	slurp, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading %q: %v", filename, err)
	}
	return slurp
}
