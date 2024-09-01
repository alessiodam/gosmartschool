package structs

type Message struct {
	ID                string `xml:"id"`
	From              string `xml:"from"`
	Subject           string `xml:"subject"`
	Date              string `xml:"date"`
	Status            string `xml:"status"`
	Attachment        string `xml:"attachment"`
	Unread            string `xml:"unread"`
	Label             string `xml:"label"`
	Deleted           string `xml:"deleted"`
	AllowReply        string `xml:"allowreply"`
	AllowReplyEnabled string `xml:"allowreplyenabled"`
	HasReply          string `xml:"hasreply"`
	HasForward        string `xml:"hasForward"`
	RealBox           string `xml:"realBox"`
	SendDate          string `xml:"sendDate"`
	SenderPicture     string `xml:"fromImage"`
}

type MessageXMLData struct {
	Messages []Message `xml:"messages>message"`
}

type MessageXMLAction struct {
	Data MessageXMLData `xml:"data"`
}

type MessageXMLResponse struct {
	Actions []MessageXMLAction `xml:"actions>action"`
}

type MessageXMLServer struct {
	Response MessageXMLResponse `xml:"response"`
}
