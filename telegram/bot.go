package telegram

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/korjavin/dstgbot/api"
	"github.com/korjavin/dstgbot/cache"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	groupID int64
	cache   *cache.MessageCache
	botName string
}

func NewBot(token string, groupID string, cache *cache.MessageCache) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	gid, err := strconv.ParseInt(groupID, 10, 64)
	if err != nil {
		return nil, err
	}

	// Get bot info
	self, err := botAPI.GetMe()
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:     botAPI,
		groupID: gid,
		cache:   cache,
		botName: self.UserName,
	}, nil
}

func (b *Bot) Start(deepseekClient *api.Client) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.Chat.ID == b.groupID {
			if err := b.handleMessage(update.Message, deepseekClient); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}

	return nil
}

func (b *Bot) handleMessage(msg *tgbotapi.Message, client *api.Client) error {
	// Store message in cache
	replyToID := 0
	if msg.ReplyToMessage != nil {
		replyToID = msg.ReplyToMessage.MessageID
	}

	b.cache.Add(cache.Message{
		ID:        msg.MessageID,
		Text:      msg.Text,
		ReplyToID: replyToID,
		Timestamp: msg.Time(),
	})

	// Check if message is for the bot
	if !b.isMessageForBot(msg) {
		return nil
	}

	// Handle unsupported message types
	if msg.Photo != nil || msg.Voice != nil || msg.Audio != nil {
		reply := tgbotapi.NewMessage(msg.Chat.ID, "Sorry, I only support text messages for now.")
		reply.ReplyToMessageID = msg.MessageID
		_, err := b.api.Send(reply)
		return err
	}

	// Get conversation context
	messages := b.getConversationContext(msg)

	// Get response from DeepSeek
	response, err := client.CreateChatCompletion(context.Background(), messages)
	if err != nil {
		log.Printf("DeepSeek API error: %v", err)
		reply := tgbotapi.NewMessage(msg.Chat.ID, "Sorry, I'm having trouble processing your request. Please try again later.")
		reply.ReplyToMessageID = msg.MessageID
		_, err := b.api.Send(reply)
		return err
	}

	// Send response
	reply := tgbotapi.NewMessage(msg.Chat.ID, response)
	reply.ReplyToMessageID = msg.MessageID
	_, err = b.api.Send(reply)
	return err
}

func (b *Bot) isMessageForBot(msg *tgbotapi.Message) bool {
	// Check if message is a reply to the bot
	if msg.ReplyToMessage != nil && msg.ReplyToMessage.From.UserName == b.botName {
		return true
	}

	// Check if message mentions the bot
	if strings.Contains(strings.ToLower(msg.Text), "@"+strings.ToLower(b.botName)) {
		return true
	}

	return false
}

func (b *Bot) getConversationContext(msg *tgbotapi.Message) []api.Message {
	var messages []api.Message

	// Add system message if available
	if systemMsg := os.Getenv("SYSTEM_MSG"); systemMsg != "" {
		messages = append(messages, api.Message{
			Role:    "system",
			Content: systemMsg,
		})
	}

	// Get message thread from cache
	cachedMessages := b.cache.GetThread(msg.MessageID)
	for _, cachedMsg := range cachedMessages {
		role := "user"
		if cachedMsg.ID == msg.MessageID && msg.From.UserName == b.botName {
			role = "assistant"
		}
		messages = append(messages, api.Message{
			Role:    role,
			Content: cachedMsg.Text,
		})
	}

	return messages
}
