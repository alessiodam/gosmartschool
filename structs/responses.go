package structs

import "encoding/xml"

// User represents a single user retrieved from the API.
type User struct {
	XMLName    xml.Name `xml:"user"`
	UserID     string   `xml:"userID"`
	Text       string   `xml:"text"`
	Value      string   `xml:"value"`
	Selectable string   `xml:"selectable"`
	SsID       string   `xml:"ssID"`
	ClassName  string   `xml:"classname"`
	SchoolName string   `xml:"schoolname"`
	Picture    string   `xml:"picture"`
}

// SearchUsersResponse represents the structure for the response of a user search.
type SearchUsersResponse struct {
	XMLName xml.Name `xml:"results"`
	Users   []User   `xml:"user"`
}

type Course struct {
	Icon       string `json:"icon"`
	ID         int    `json:"id"`
	Name       string `json:"name"`
	PlatformID int    `json:"platformId"`
	Teacher    string `json:"teacher"`
	URL        string `json:"url"`
	Descr      string `json:"descr"`
}
