package telegram

import (
	"context"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

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

	log.Printf("Bot username: %s", self.UserName)

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

	log.Println("Bot started listening for messages.")

	for update := range updates {
		if update.Message != nil {
			log.Printf("Received message from chat ID %d (Type: %s): %s", update.Message.Chat.ID, update.Message.Chat.Type, update.Message.Text)

			if update.Message.Chat.ID == b.groupID || update.Message.Chat.Type == "private" {
				go func(msg *tgbotapi.Message) {
					if err := b.handleMessage(msg, deepseekClient); err != nil {
						log.Printf("Error handling message: %v", err)
					}
				}(update.Message)
			} else {
				log.Printf("Ignoring message from chat ID %d (expected %d)", update.Message.Chat.ID, b.groupID)
			}
		}
	}

	return nil
}

func (b *Bot) handleMessage(msg *tgbotapi.Message, client *api.Client) error {
	log.Printf("Handling message ID %d from %s: %s", msg.MessageID, msg.From.UserName, msg.Text)

	// Random delay between 5-15 seconds
	delay := time.Duration(rand.Intn(11)+5) * time.Second
	log.Printf("Sleeping for %v before processing...", delay)
	time.Sleep(delay)

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
	log.Printf("Added message %d to cache. ReplyToID: %d", msg.MessageID, replyToID)

	// Check if message is for the bot
	if !b.isMessageForBot(msg) {
		log.Println("Message is not for the bot, ignoring.")
		return nil
	}
	/*
		// Handle unsupported message types
		if msg.Photo != nil || msg.Voice != nil || msg.Audio != nil {
			reply := tgbotapi.NewMessage(msg.Chat.ID, "Sorry, I only support text messages for now.")
			reply.ReplyToMessageID = msg.MessageID
			sentMsg, err := b.api.Send(reply)
			if err != nil {
				log.Printf("Error sending message: %v", err)
				return err
			}
			// Store bot's message to cache
			b.cache.Add(cache.Message{
				ID:        sentMsg.MessageID,
				Text:      reply.Text,
				ReplyToID: msg.MessageID,
				Timestamp: sentMsg.Time(),
			})
			return nil
		}
	*/

	// Get conversation context
	messages := b.getConversationContext(msg)
	log.Printf("Conversation context: %v", messages)

	// Get response from DeepSeek
	log.Println("Sending request to DeepSeek API...")
	response, err := client.CreateChatCompletion(context.Background(), messages)
	if err != nil {
		log.Printf("DeepSeek API error: %v", err)
		reply := tgbotapi.NewMessage(msg.Chat.ID, "Sorry, I'm having trouble processing your request. Please try again later.")
		reply.ReplyToMessageID = msg.MessageID
		sentMsg, err := b.api.Send(reply)
		if err != nil {
			log.Printf("Error sending message: %v", err)
			return err
		}
		// Store bot's message to cache
		b.cache.Add(cache.Message{
			ID:        sentMsg.MessageID,
			Text:      reply.Text,
			ReplyToID: msg.MessageID,
			Timestamp: sentMsg.Time(),
		})
		return err
	}
	log.Println("Received response from DeepSeek API.")

	// Send response
	reply := tgbotapi.NewMessage(msg.Chat.ID, response)
	reply.ReplyToMessageID = msg.MessageID
	sentMsg, err := b.api.Send(reply)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}
	log.Printf("Sent message to Telegram: %s", response)

	// Store bot's message to cache
	b.cache.Add(cache.Message{
		ID:        sentMsg.MessageID,
		Text:      response,
		ReplyToID: msg.MessageID,
		Timestamp: sentMsg.Time(),
	})
	return err
}

func (b *Bot) isMessageForBot(msg *tgbotapi.Message) bool {
	// Private DM: Always for the bot
	if msg.Chat.Type == "private" {
		log.Println("Message is a private DM.")
		return true
	}

	// Check if message is a reply to the bot
	if msg.ReplyToMessage != nil && msg.ReplyToMessage.From.UserName == b.botName {
		log.Println("Message is a reply to the bot.")
		return true
	}

	text := strings.ToLower(msg.Text)
	botName := strings.ToLower(b.botName)

	// Check if message mentions the bot (with or without @)
	if strings.Contains(text, botName) {
		log.Println("Message mentions the bot name.")
		return true
	}

	// Check if message mentions "Titania"
	if strings.Contains(text, "titania") {
		log.Println("Message mentions Titania.")
		return true
	}

	return false
}

func (b *Bot) getConversationContext(msg *tgbotapi.Message) []api.Message {
	var messages []api.Message
	messageThread := b.getMessageThread(msg.MessageID, 10)

	// Add system message if available
	if systemMsg := os.Getenv("SYSTEM_MSG"); systemMsg != "" {
		messages = append(messages, api.Message{
			Role:    "system",
			Content: systemMsg,
		})
	}

	// Add message thread to context
	for _, cachedMsg := range messageThread {
		role := "user"
		//if cachedMsg.ID == msg.MessageID && msg.From.UserName == b.botName {
		//	role = "assistant"
		//}
		if cachedMsg.ID == msg.MessageID {
			if msg.From.UserName == b.botName {
				role = "assistant"
			}
		}
		messages = append(messages, api.Message{
			Role:    role,
			Content: cachedMsg.Text,
		})
	}

	return messages
}

func (b *Bot) getMessageThread(messageID int, limit int) []cache.Message {
	var thread []cache.Message
	currentID := messageID
	count := 0

	for count < limit {
		msg, found := b.cache.GetByID(currentID)
		if !found {
			break
		}
		thread = append(thread, msg)
		currentID = msg.ReplyToID
		count++
	}

	// Reverse to maintain chronological order
	for i, j := 0, len(thread)-1; i < j; i, j = i+1, j-1 {
		thread[i], thread[j] = thread[j], thread[i]
	}

	log.Printf("Getting thread for message ID %d from cache. Found %d messages.", messageID, len(thread))
	return thread
}
