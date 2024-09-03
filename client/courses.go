package client

import (
	"encoding/json"
	"gosmartschool/structs"
	"io"
	"net/http"
	"net/url"
)

func (client *SmartSchoolClient) GetCourses() ([]structs.Course, error) {
	data := url.Values{}

	request, body, err := client.sendRequest("POST", "/Topnav/getCourseConfig", data.Encode(), nil)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			client.ApiLogger.Error(err)
		}
	}(request.Body)

	if request.StatusCode == http.StatusOK {
		client.ApiLogger.Info("Course config received")
	}

	var jsonBody map[string]interface{}
	err = json.Unmarshal([]byte(body), &jsonBody)
	if err != nil {
		return nil, err
	}

	var courses []structs.Course
	for _, course := range jsonBody["own"].([]interface{}) {
		courseMap := course.(map[string]interface{})
		courses = append(courses, structs.Course{
			ID:         int(courseMap["id"].(float64)),
			PlatformID: int(courseMap["platformId"].(float64)),
			Name:       courseMap["name"].(string),
			Descr:      courseMap["descr"].(string),
			Icon:       courseMap["icon"].(string),
			URL:        courseMap["url"].(string),
			Teacher:    courseMap["teacher"].(string),
		})
	}

	return courses, nil
}
