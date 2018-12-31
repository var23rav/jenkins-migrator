package components

import (
	"fmt"
	"net/url"
	"strings"
)

type JenkinsFolder struct {
	JenkinsComponent
	Name string
	FullName string
	Url string
	Jobs []JenkinsJob
}

func (jenkinsFolder JenkinsFolder) GenerateURL (urlType string, hostUrl string) *url.URL {

	var reqUrlComp string
	reqParams := &url.Values{}
	jobHierarchyArr := strings.Split(jenkinsFolder.FullName, "/")
	jobParents := jobHierarchyArr[: (len(jobHierarchyArr) - 1) ]
	jobName := jobHierarchyArr[len(jobHierarchyArr) - 1]
	for _, parent := range jobParents {
		reqUrlComp += fmt.Sprintf("/job/%v", parent)
	}
	switch urlType {
		case "Create":
			reqUrlComp += "/createItem"
			fmt.Println(jenkinsFolder)
			reqParams.Add("name", jobName)
			reqParams.Add("mode", "com.cloudbees.hudson.plugins.folder.Folder")
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
		panic(fmt.Sprintf("JenkinsFolder generateUrl error %v", err))
	}
	reqUrl.RawQuery = reqParams.Encode()
	return reqUrl
}
