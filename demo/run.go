package main

import (
	jClient "jenkins-migrator/api"
	jComponent "jenkins-migrator/components"
)

func main() {
	sourceClient := &jClient.ApiClient{BaseUrl: jClient.SOURCE_URL, Username: "admin", Password: "admin"}
	// jobUrl := source.BaseUrl + "/api/json"
	jobUrl := sourceClient.BaseUrl
	var jc jComponent.JenkinsComponent
	jc.Get(jobUrl, sourceClient)
	// migrator.Me()
}