package server

import (
	"fmt"
	"github.com/ReneKroon/ttlcache/v2"
	"github.com/fossteams/fossteams-backend/internal/errors"
	"github.com/fossteams/fossteams-backend/internal/messages"
	v1 "github.com/fossteams/fossteams-backend/internal/responses/api/v1"
	teams "github.com/fossteams/teams-api"
	models "github.com/fossteams/teams-api/pkg/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Server struct {
	logger *logrus.Logger
	cache  *ttlcache.Cache
	e      *gin.Engine
	teams  *teams.TeamsClient
}

const cacheTTL = 300

func (s *Server) setupRoutes() {
	apiEndpoint := s.e.Group("/api")
	v1Endpoint := apiEndpoint.Group("/v1")

	s.setupApiV1(v1Endpoint)
}

func (s *Server) setupApiV1(endpoint *gin.RouterGroup) {
	endpoint.GET("/conversations", s.v1GetConversations)
	endpoint.GET("/conversations/:id", s.v1GetSingleConversation)
	endpoint.GET("/conversations/:id/profilePicture", s.v1GetConversationProfilePicture)
}

func parseMembers(members []models.ChatMember) []v1.ChatMember {
	var pMembers []v1.ChatMember

	for _, m := range members {
		pMembers = append(pMembers, v1.ChatMember{
			Mri:      m.Mri,
			Role:     string(m.Role),
			TenantId: m.TenantId,
			ObjectId: m.ObjectId,
		})
	}
	return pMembers
}

func parseChannels(channels []models.Channel) []v1.Channel {
	var chans []v1.Channel
	for _, c := range channels {
		chans = append(chans, v1.Channel{
			Id:           c.Id,
			DisplayName:  c.DisplayName,
			Description:  c.Description,
			Creator:      c.Creator,
			ParentTeamId: c.ParentTeamId,
			LastMessage:  processMessage(c.LastMessage),
		})
	}
	return chans
}

func processMessage(message models.Message) v1.ShortMessage {
	return v1.ShortMessage{
		Id:      message.Id,
		Content: message.Content,
		From:    message.From,
	}

}

func New(logger *logrus.Logger) (*Server, error) {
	if logger == nil {
		logger = logrus.New()
	}

	engine := gin.Default()
	teams, err := teams.New()

	if err != nil {
		return nil, fmt.Errorf("unable to initialize Teams client: %v", err)
	}

	cache := ttlcache.NewCache()
	err = cache.SetTTL(cacheTTL * time.Second)
	if err != nil {
		return nil, fmt.Errorf("unable to set ttlcache TTL: %v", err)
	}

	s := Server{
		logger: logger,
		e:      engine,
		teams:  teams,
		cache:  cache,
	}
	s.setupCors()
	s.setupRoutes()

	return &s, nil
}

func (s *Server) Start(addr string) error {
	return s.e.Run(addr)
}

func (s *Server) v1GetConversations(c *gin.Context) {
	const conversationsKey = "conversations"
	if v, err := s.cache.Get(conversationsKey); err != ttlcache.ErrNotFound {
		// cache hit
		c.JSON(http.StatusOK, v)
		return
	}

	// Fetch conversations
	if !s.checkTeams(c) {
		return
	}

	conv, err := s.teams.GetConversations()
	if err != nil {
		s.logger.Errorf("unable to get conversations: %v", err)
		c.JSON(http.StatusInternalServerError, errors.ApiError{
			Message: "unable to fetch conversations",
		})
		return
	}

	// Process conversations
	var chats []v1.Chat
	var teams []v1.Team

	for _, c := range conv.Chats {
		pc := v1.Chat{
			Id:                  c.Id,
			Title:               c.Title,
			LastMessage:         processMessage(c.LastMessage),
			IsOneOnOne:          c.IsOneOnOne,
			Creator:             c.Creator,
			IsRead:              c.IsRead,
			Members:             parseMembers(c.Members),
			IsLastMessageFromMe: c.IsLastMessageFromMe,
		}

		chats = append(chats, pc)
	}

	for _, t := range conv.Teams {
		pt := v1.Team{
			Creator:     t.Creator,
			Id:          t.Id,
			DisplayName: t.DisplayName,
			Channels:    parseChannels(t.Channels),
		}
		teams = append(teams, pt)
	}

	resp := v1.Conversations{
		Chats: chats,
		Teams: teams,
	}

	err = s.cache.Set(conversationsKey, resp)
	if err != nil {
		s.logger.Warnf("unable to set cache entry: %v", err)
	}

	c.JSON(http.StatusOK, resp)
	return
}

func (s *Server) v1GetSingleConversation(c *gin.Context) {
	// Conversation id is :id
	convId := c.Param("id")
	cacheKey := "conversations/" + convId

	if v, err := s.cache.Get(cacheKey); err != ttlcache.ErrNotFound {
		// cache hit
		c.JSON(http.StatusOK, v)
		return
	}

	if !s.checkTeams(c) {
		return
	}
	if convId == "" {
		c.JSON(http.StatusBadRequest, errors.ApiError{Message: "invalid conversation ID"})
		return
	}

	s.teams.Debug(true)
	s.teams.ChatSvc().DebugDisallowUnknownFields(false)

	chatMessages, err := s.teams.GetMessages(&models.Channel{Id: convId})
	if err != nil {
		s.logger.Errorf("unable to get messages for convId=%s: %v", convId, err)
		c.JSON(http.StatusInternalServerError, errors.ApiError{Message: "unable to get messages"})
		return
	}

	pMessages := s.parseMessages(chatMessages)
	resp := v1.ConversationResponse{
		Messages: pMessages,
	}
	err = s.cache.Set(cacheKey, resp)
	if err != nil {
		s.logger.Warnf("unable to set cache entry: %v", err)
	}
	c.JSON(http.StatusOK, resp)
	return
}

func (s *Server) v1GetConversationProfilePicture(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please provide an id"})
		return
	}

	buff, err := s.teams.GetTeamsProfilePicture(id)
	if err != nil {
		s.logger.Errorf("unable to get profile picture: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "unable to get profile picture"})
		return
	}

	c.Status(http.StatusOK)
	_, _ = c.Writer.Write(buff)
}

func (s *Server) parseMessages(msgs []models.ChatMessage) []v1.Message {
	var pMessages []v1.Message
	threads := map[string][]v1.Message{}
	parentMessages := map[string]v1.Message{}

	seenParent := map[string]bool{}

	for _, m := range msgs {
		msg := v1.Message{
			ShortMessage: v1.ShortMessage{
				Id:           m.Id,
				CleanContent: messages.ParseMessageContent(m.Content),
				Content:      m.Content,
				From:         m.From,
			},
			ImDisplayName:       m.ImDisplayName,
			OriginalArrivalTime: m.OriginalArrivalTime,
			ConversationId:      m.ConversationId,
			ParentID:            parseParentId(m.ConversationLink),
			SequenceId:          m.SequenceId,
			MessageType:         string(m.Type),
			Type:                m.MessageType,
			Subject:             m.Properties.Subject,
			Title:               m.Properties.Title,
			Reactions:           parseReactions(m.Properties.Emotions),
		}

		if msg.ParentID == msg.Id {
			parentMessages[msg.Id] = msg
		} else {
			threads[msg.ParentID] = append(threads[msg.ParentID], msg)
		}
	}

	for threadId, threadMessages := range threads {
		if val, ok := parentMessages[threadId]; !ok {
			s.logger.Warnf("thread %s doesn't exist", threadId)
			continue
		} else {
			val.Replies = append(val.Replies, threadMessages...)
			sort.Sort(bySequenceId(val.Replies))
			pMessages = append(pMessages, val)
			seenParent[threadId] = true
		}
	}

	for k, msg := range parentMessages {
		if _, ok := seenParent[k]; !ok {
			pMessages = append(pMessages, msg)
		}
	}

	sort.Sort(bySequenceId(pMessages))

	return pMessages
}

func parseReactions(emotions []models.Emotion) map[string]int {
	reactions := map[string]int{}
	for _, em := range emotions {
		reactions[strings.ToLower(em.Key)] = len(em.Users)
	}
	return reactions
}

func parseParentId(link string) string {
	mUrl, err := url.Parse(link)
	if err != nil {
		return ""
	}

	// get message id
	kv := strings.Split(mUrl.Path, ";")
	if len(kv) != 2 {
		return ""
	}

	messageId := strings.Split(kv[1], "=")
	if messageId[0] != "messageid" {
		return ""
	}
	if len(messageId) != 2 {
		return ""
	}

	return messageId[1]
}

type bySequenceId []v1.Message

func (b bySequenceId) Len() int {
	return len(b)
}

func (b bySequenceId) Less(i, j int) bool {
	return b[i].SequenceId < b[j].SequenceId
}

func (b bySequenceId) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

var _ sort.Interface = bySequenceId{}

func (s *Server) checkTeams(c *gin.Context) bool {
	if s.teams == nil {
		c.JSON(http.StatusInternalServerError, errors.ApiError{
			Message: "teams client not ready",
		})
		return false
	}
	return true
}

func (s *Server) setupCors() {
	s.e.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://127.0.0.1:8080"},
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Origin"},
		MaxAge:       12 * time.Hour,
	}))
}
