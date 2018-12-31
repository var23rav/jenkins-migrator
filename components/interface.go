package components

import (
	"net/url"
	"jenkins-migrator/api"
)

type IComponent interface {
	GenerateURL(string, string) *url.URL
	Get(string, *api.ApiClient)
	getJobData(*url.URL, api.ApiClient) []byte
	createJob(*url.URL, api.ApiClient, []byte) bool
	updateJob(*url.URL, api.ApiClient, []byte) bool
	migrateJobs(api.ApiClient, api.ApiClient)
}