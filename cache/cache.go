package cache

import (
	"sync"
	"time"
)

type Message struct {
	ID        int
	Text      string
	ReplyToID int
	Timestamp time.Time
}

type MessageCache struct {
	messages []Message
	mu       sync.RWMutex
	capacity int
}

func NewMessageCache(capacity int) *MessageCache {
	return &MessageCache{
		messages: make([]Message, 0, capacity),
		capacity: capacity,
	}
}

func (c *MessageCache) Add(msg Message) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.messages) >= c.capacity {
		c.messages = c.messages[1:]
	}
	c.messages = append(c.messages, msg)
}

func (c *MessageCache) GetByID(id int) (Message, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, msg := range c.messages {
		if msg.ID == id {
			return msg, true
		}
	}
	return Message{}, false
}

func (c *MessageCache) GetThread(id int) []Message {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var thread []Message
	currentID := id

	for {
		msg, found := c.GetByID(currentID)
		if !found {
			break
		}
		thread = append(thread, msg)
		currentID = msg.ReplyToID
	}

	// Reverse to maintain chronological order
	for i, j := 0, len(thread)-1; i < j; i, j = i+1, j-1 {
		thread[i], thread[j] = thread[j], thread[i]
	}

	return thread
}
