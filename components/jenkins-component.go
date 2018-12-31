package components

import (
	"fmt"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"jenkins-migrator/api"
)

func Me() {
	fmt.Println("var23--Me")
}

type JenkinsComponent struct {}

func (component JenkinsComponent) Get(jobUrl string, apiClient *api.ApiClient) {

	fmt.Println("-----------------------------------")
	fmt.Println(jobUrl)
	jobUrl = api.PrepareUrl(jobUrl) + "/api/json"
	fmt.Println(jobUrl)
	var reqUrl *url.URL
	reqJobUrl, err := reqUrl.Parse(jobUrl)
	if err != nil {
		panic(fmt.Sprintf("URL parsing failed %v", err))
	}
	resp, err := apiClient.ApiGetQuery(reqJobUrl)
	if err != nil {
		panic(fmt.Sprintf("URL api query error for `%v`", jobUrl))
	}
	defer resp.Body.Close()
	jobsDataRespString, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Respone reading failed for `%v`", jobUrl))
	}
	// fmt.Println(string(jobsDataRespString))
	var jenkinsFolder JenkinsFolder
	err = json.Unmarshal(jobsDataRespString, &jenkinsFolder)
	if err != nil {
		panic(fmt.Sprintf("JSON Response unmarshal failed with err `%v`", err))
	}
	if err == nil {
	
		destApiClient := api.ApiClient{BaseUrl: api.DESTINATION_URL, Username: "admin", Password: "admin"}
		sourceApiClient := api.ApiClient{BaseUrl: api.SOURCE_URL, Username: "admin", Password: "admin"}
		if len(jenkinsFolder.Jobs) == 0 {
			
			var jenkinsJob JenkinsJob
			err = json.Unmarshal(jobsDataRespString, &jenkinsJob)
			if err != nil {
				panic(fmt.Sprintf("JSON Response unmarshal failed with err `%v`", err))
			}
			MigrateJobs(jenkinsJob, sourceApiClient, destApiClient)
		} else {

			// fmt.Println(jenkinsFolder.Jobs)
			MigrateJobs(jenkinsFolder, sourceApiClient, destApiClient)
			// for index, jenkinChildJob := range jenkinsFolder.Jobs {
			for _, jenkinChildJob := range jenkinsFolder.Jobs {
				// if index > 0 {
				// 	break
				// }
				jenkinChildJobUrl := jenkinChildJob.Url
				// fmt.Println(jenkinChildJobUrl)
				component.Get(jenkinChildJobUrl, apiClient)
			}
		}
	}
	// fmt.Println(string(jobsDataRespString))
	// return jenkinsFolder, nil
}

func GetJobData(component IComponent, apiClient api.ApiClient) []byte {
	reqJobDataUrl := component.GenerateURL("DataXML", apiClient.BaseUrl)
	return component.getJobData(reqJobDataUrl, apiClient)
}
func (component JenkinsComponent) getJobData(reqJobDataUrl *url.URL, apiClient api.ApiClient) []byte {

	respJobData, err := apiClient.ApiGetQuery(reqJobDataUrl)
	if err != nil {
		panic(fmt.Sprintf("Migration api failed `%v`", err))
	}
	defer respJobData.Body.Close()
	respJobDataText, err := ioutil.ReadAll(respJobData.Body)
	if err != nil {
		panic(fmt.Sprintf("Job Data Fetch api failed `%v`", err))
	}
	fmt.Println(respJobData.StatusCode)
	// fmt.Println(string(respJobDataText))
	return respJobDataText
}

func CreateJob(component IComponent, apiClient api.ApiClient, jobData []byte) bool {
	reqJobCreateUrl := component.GenerateURL("Create", apiClient.BaseUrl)
	return component.createJob(reqJobCreateUrl, apiClient, jobData)
}
func (component JenkinsComponent) createJob(reqJobCreateUrl *url.URL, apiClient api.ApiClient, jobData []byte) bool {

	// Create new Job
	fmt.Println(reqJobCreateUrl)
	respJobCreate, err := apiClient.ApiPostQuery(reqJobCreateUrl, jobData)
	respJobCreateText, err := ioutil.ReadAll(respJobCreate.Body)
	defer respJobCreate.Body.Close()
	if err != nil {
		panic(fmt.Sprintf("New Job Create api failed `%v`", err))
	}
	// fmt.Println(respJobCreate.Status)
	fmt.Println(respJobCreate.StatusCode)
	if api.IsDev {
		fmt.Println(string(respJobCreateText))
	}
	return respJobCreate.StatusCode == 200
}

func UpdateJob(component IComponent, apiClient api.ApiClient, jobData []byte) bool {
	reqJobUpdateUrl := component.GenerateURL("Update", apiClient.BaseUrl)
	return component.updateJob(reqJobUpdateUrl, apiClient, jobData)
}
func (component JenkinsComponent) updateJob(reqJobUpdateUrl *url.URL, apiClient api.ApiClient, jobData []byte) bool {

	// Update existing Job
	fmt.Println(reqJobUpdateUrl)
	respJobCreate, err := apiClient.ApiPostQuery(reqJobUpdateUrl, jobData)
	respJobCreateText, err := ioutil.ReadAll(respJobCreate.Body)
	defer respJobCreate.Body.Close()
	if err != nil {
		panic(fmt.Sprintf("New Job Create api failed `%v`", err))
	}
	// fmt.Println(respJobCreate.Status)
	fmt.Println(respJobCreate.StatusCode)
	if api.IsDev {
		fmt.Println(string(respJobCreateText))
	}
	return respJobCreate.StatusCode == 200
}

func MigrateJobs(component IComponent, sourceApiClient api.ApiClient, destApiClient api.ApiClient) {
	// component.migrateJobs(sourceApiClient, destApiClient)
	respJobData := GetJobData(component, sourceApiClient)
	isNewJobCreated := CreateJob(component, destApiClient, respJobData)
	fmt.Println("New Job created : " + strconv.FormatBool(isNewJobCreated))
	if !isNewJobCreated {
		isJobupdated := UpdateJob(component, destApiClient, respJobData)
		fmt.Println("Existing Job Updated: " + strconv.FormatBool(isJobupdated))
	}
}
func (component JenkinsComponent) migrateJobs(sourceApiClient api.ApiClient, destApiClient api.ApiClient) {

	// respJobData := component.IgetJobData(sourceApiClient)
	// isNewJobCreated := component.IcreateJob(destApiClient, respJobData)
	// fmt.Println("New Job created : " + strconv.FormatBool(isNewJobCreated))
	// if !isNewJobCreated {
	// 	isJobupdated := component.IupdateJob(destApiClient, respJobData)
	// 	fmt.Println("Existing Job Updated: " + strconv.FormatBool(isJobupdated))
	// }

		// respJobData := GetJobData(component, sourceApiClient)
		// isNewJobCreated := CreateJob(component, destApiClient, respJobData)
		// fmt.Println("New Job created : " + strconv.FormatBool(isNewJobCreated))
		// if !isNewJobCreated {
		// 	isJobupdated := UpdateJob(component, destApiClient, respJobData)
		// 	fmt.Println("Existing Job Updated: " + strconv.FormatBool(isJobupdated))
		// }
}