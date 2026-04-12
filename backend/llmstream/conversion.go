package llmstream

func NewConversation(systemPrompt string) *ConversationManager {
	return &ConversationManager{
		History:      []Message{{Role: "system", Content: systemPrompt}},
		SystemPrompt: systemPrompt,
	}
}

// AddUser appends a user turn
func (c *ConversationManager) AddUser(content string) {
	c.History = append(c.History, Message{Role: "user", Content: content})
}

// AddAssistant appends the collected streaming response
func (c *ConversationManager) AddAssistant(content string) {
	c.History = append(c.History, Message{Role: "assistant", Content: content})
}
