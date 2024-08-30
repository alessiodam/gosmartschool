package structs

// SearchUsersRequest represents the data structure for searching users by name.
type SearchUsersRequest struct {
	Val          string `url:"val"`
	Type         int    `url:"type"`
	ParentNodeID string `url:"parentNodeId"`
	XML          string `url:"xml"`
}

// MessageListRequest represents the data structure for requesting a list of messages.
type MessageListRequest struct {
	Command string `url:"command"`
}

// ShowMessageRequest represents the data structure for requesting a specific message by ID.
type ShowMessageRequest struct {
	Command string `url:"command"`
}

// DeleteMessageRequest represents the data structure for deleting a message by ID.
type DeleteMessageRequest struct {
	Command string `url:"command"`
}
