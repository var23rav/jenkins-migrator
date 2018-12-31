package api

import (
	"fmt"
	"strings"
	"net/http"
	"net/url"
	"io/ioutil"
	"bytes"
)

var SOURCE_URL string = "http://192.168.0.22:8000"
var DESTINATION_URL string = "http://192.168.0.23:8000"
var IsDev bool = false

type ApiClient struct {
	BaseUrl string
	Token string
	Username string
	Password string
}

func PrepareUrl(url string) string {
	return strings.Trim(url, "/")
}

func (apiClient *ApiClient) updateCsrfToken() {
	
	crumbIssuerUrl := apiClient.BaseUrl + "/crumbIssuer/api/xml"
	var reqUrl *url.URL
	csrfTokenUrl, err := reqUrl.Parse(crumbIssuerUrl)
	if err != nil {
		panic(fmt.Sprintf("Url `%v` parsing error ", crumbIssuerUrl))
	}
	var reqParams = url.Values{}
	reqParams.Add("xpath", `concat(//crumbRequestField,":",//crumb)`)
	csrfTokenUrl.RawQuery = reqParams.Encode()

	reqClient := &http.Client{}
	req, err := http.NewRequest("GET", csrfTokenUrl.String(), nil)
	if err != nil {
		panic("Get Request generation failed")
	}

	req.SetBasicAuth(apiClient.Username, apiClient.Password)
	resp, err := reqClient.Do(req)
	if err != nil {
		panic("Jenkins CSRF Token api request failed")
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
	csrfToken, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("Jenkins CSRF Response reading failed")
	}
	apiClient.Token = string(csrfToken)
}

func (apiClient *ApiClient) prepareAPIHeader(req *http.Request) {

	apiClient.updateCsrfToken()
	csrfTokenArr := strings.Split(apiClient.Token, ":")
	if len(csrfTokenArr) == 2 {
		req.Header.Set(csrfTokenArr[0], csrfTokenArr[1])
	} else {
		fmt.Println("Error: CSRF Token update failed")
	}
	req.SetBasicAuth(apiClient.Username, apiClient.Password)
}

func (apiClient *ApiClient) ApiGetQuery(reqUrl *url.URL) (*http.Response, error) {

	reqClient := &http.Client{}
	req, err := http.NewRequest("GET", reqUrl.String(), nil)
	if err != nil {
		fmt.Println("Error: ApiQuery request generation failed")
	}
	apiClient.prepareAPIHeader(req)
	return reqClient.Do(req)
}

func (apiClient *ApiClient) ApiPostQuery(reqUrl *url.URL, reqData []byte) (*http.Response, error) {

	fmt.Println(">> -- " + reqUrl.String())
	reqClient := &http.Client{}
	req, err := http.NewRequest("POST", reqUrl.String(), bytes.NewBuffer(reqData))
	if err != nil {
		fmt.Println("Error: ApiQuery request generation failed")
	}
	req.Header.Set("Content-Type", "text/xml")
	apiClient.prepareAPIHeader(req)
	return reqClient.Do(req)
}
