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

// Message represents a single message retrieved from the API.
type Message struct {
	XMLName                  xml.Name `xml:"message"`
	ID                       string   `xml:"id"`
	From                     string   `xml:"from"`
	To                       string   `xml:"to"`
	Subject                  string   `xml:"subject"`
	Date                     string   `xml:"date"`
	Body                     string   `xml:"body"`
	Status                   string   `xml:"status"`
	Attachment               string   `xml:"attachment"`
	Unread                   string   `xml:"unread"`
	Label                    string   `xml:"label"`
	Receivers                []string `xml:"receivers>to"`
	CcReceivers              []string `xml:"ccreceivers>cc"`
	BccReceivers             []string `xml:"bccreceivers>bcc"`
	SenderPicture            string   `xml:"senderPicture"`
	MarkedInLVS              string   `xml:"markedInLVS"`
	FromTeam                 string   `xml:"fromTeam"`
	TotalNrOtherToReciviers  string   `xml:"totalNrOtherToReciviers"`
	TotalnrOtherCcReceivers  string   `xml:"totalnrOtherCcReceivers"`
	TotalnrOtherBccReceivers string   `xml:"totalnrOtherBccReceivers"`
	CanReply                 string   `xml:"canReply"`
	HasReply                 string   `xml:"hasReply"`
	HasForward               string   `xml:"hasForward"`
	SendDate                 string   `xml:"sendDate"`
}

// MessageListResponse represents the structure for the response of a message list.
type MessageListResponse struct {
	XMLName  xml.Name  `xml:"response"`
	Messages []Message `xml:"message"`
}

// ShowMessageResponse represents the structure for the response of a single message retrieval.
type ShowMessageResponse struct {
	XMLName xml.Name `xml:"response"`
	Message Message  `xml:"message"`
}
