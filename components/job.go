package components

import (
	"fmt"
	"net/url"
	"strings"
)

type JenkinsJob struct {
	JenkinsComponent
	Name string
	FullName string
	Url string
	Parent JenkinsFolder
}

func (jenkinsJob JenkinsJob) GenerateURL (urlType string, hostUrl string) *url.URL {

	var reqUrlComp string
	reqParams := &url.Values{}
	jobHierarchyArr := strings.Split(jenkinsJob.FullName, "/")
	jobParents := jobHierarchyArr[: (len(jobHierarchyArr) - 1) ]
	jobName := jobHierarchyArr[len(jobHierarchyArr) - 1]
	for _, parent := range jobParents {
		reqUrlComp += fmt.Sprintf("/job/%v", parent)
	}
	switch urlType {
		case "Create":
			reqUrlComp += "/createItem"
			fmt.Println(jenkinsJob)
			reqParams.Add("name", jobName)
		case "DataXML", "Update":
			reqUrlComp += fmt.Sprintf("/job/%v/config.xml", jobName)
		case "DataJSON":
			reqUrlComp += fmt.Sprintf("/job/%v/api/json", jobName)			
		default:
			fmt.Println("Error: Requestd for unknow url")

	}
	var urlUrl *url.URL
	reqUrlStr := hostUrl + reqUrlComp
	reqUrl, err := urlUrl.Parse(reqUrlStr)
	if err != nil {
		panic(fmt.Sprintf("JenkinsJob generateUrl error %v", err))
	}
	reqUrl.RawQuery = reqParams.Encode()
	return reqUrl
}
