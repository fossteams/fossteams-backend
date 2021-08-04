package v1

import (
	"github.com/fossteams/teams-api/pkg"
)

type Conversations struct {
	Chats []Chat `json:"chats"`
	Teams []Team `json:"teams"`
}

type Chat struct {
	Id                  string       `json:"id"`
	Title               string       `json:"title"`
	LastMessage         ShortMessage `json:"lastMessage"`
	IsOneOnOne          bool         `json:"isOneOnOne"`
	Creator             string       `json:"creator"`
	IsRead              bool         `json:"isRead"`
	IsLastMessageFromMe bool         `json:"isLastMessageFromMe"`
	Members             []ChatMember `json:"members"`
}

type Team struct {
	Creator     string    `json:"creator"`
	Id          string    `json:"id"`
	DisplayName string    `json:"displayName"`
	Channels    []Channel `json:"channels"'`
}

type Channel struct {
	Id           string       `json:"id"`
	DisplayName  string       `json:"displayName"`
	LastMessage  ShortMessage `json:"lastMessage"`
	Description  string       `json:"description"`
	Creator      string       `json:"creator"`
	ParentTeamId string       `json:"parentTeamId"`
}

type ChatMember struct {
	Mri      string `json:"mri"`
	Role     string `json:"role"`
	TenantId string `json:"tenantId"`
	ObjectId string `json:"objectId"`
}

type ShortMessage struct {
	Id      string `json:"id"`
	Content string `json:"content"`
	From    string `json:"from"`
}

type Message struct {
	ShortMessage
	ImDisplayName       string          `json:"imDisplayName"`
	OriginalArrivalTime api.RFC3339Time `json:"originalArrivalTime"`
	ConversationId      string          `json:"conversationId"`
	ParentID            string          `json:"parentID"`
	SequenceId          int64           `json:"sequenceId"`
	MessageType         string          `json:"messageType"`
	Type                string          `json:"type"`
	Subject             string          `json:"subject,omitempty"`
	Title               string          `json:"title,omitempty"`
	Reactions           map[string]int  `json:"reactions,omitempty"`
	Replies             []Message       `json:"replies,omitempty"`
}

type ConversationResponse struct {
	Messages []Message `json:"messages"`
}
