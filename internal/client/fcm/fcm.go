package fcm

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type Client struct {
	client *messaging.Client
}

type PushNotification struct {
	Token    string            `json:"token"`
	Tokens   []string          `json:"tokens,omitempty"`
	Title    string            `json:"title"`
	Body     string            `json:"body"`
	Data     map[string]string `json:"data,omitempty"`
	ImageURL string            `json:"image_url,omitempty"`
}

type BatchPushNotification struct {
	Tokens   []string          `json:"tokens"`
	Title    string            `json:"title"`
	Body     string            `json:"body"`
	Data     map[string]string `json:"data,omitempty"`
	ImageURL string            `json:"image_url,omitempty"`
}

// * NewClient creates a new FCM client
func NewClientFromMessaging(messagingClient *messaging.Client) *Client {
	return &Client{
		client: messagingClient,
	}
}

// * NewClientWithConfig creates a new FCM client with firebase config
func NewClientWithConfig(ctx context.Context, config *firebase.Config, credentialsPath string) (*Client, error) {
	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting messaging client: %v", err)
	}

	return &Client{
		client: client,
	}, nil
}

// * SendToToken sends notification to a single device token
func (c *Client) SendToToken(ctx context.Context, notification *PushNotification) (string, error) {
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: notification.Title,
			Body:  notification.Body,
		},
		Token: notification.Token,
		Data:  notification.Data,
	}

	if notification.ImageURL != "" {
		message.Notification.ImageURL = notification.ImageURL
	}

	response, err := c.client.Send(ctx, message)
	if err != nil {
		return "", fmt.Errorf("error sending message: %v", err)
	}

	return response, nil
}

// * SendToTokens sends notification to multiple device tokens
func (c *Client) SendToTokens(ctx context.Context, notification *BatchPushNotification) (*messaging.BatchResponse, error) {
	if len(notification.Tokens) == 0 {
		return nil, fmt.Errorf("no tokens provided")
	}

	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: notification.Title,
			Body:  notification.Body,
		},
		Tokens: notification.Tokens,
		Data:   notification.Data,
	}

	if notification.ImageURL != "" {
		message.Notification.ImageURL = notification.ImageURL
	}

	response, err := c.client.SendMulticast(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("error sending multicast message: %v", err)
	}

	return response, nil
}

// * SendToTopic sends notification to a topic
func (c *Client) SendToTopic(ctx context.Context, topic string, title, body string, data map[string]string) (string, error) {
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Topic: topic,
		Data:  data,
	}

	response, err := c.client.Send(ctx, message)
	if err != nil {
		return "", fmt.Errorf("error sending message to topic: %v", err)
	}

	return response, nil
}

// * SubscribeToTopic subscribes tokens to a topic
func (c *Client) SubscribeToTopic(ctx context.Context, tokens []string, topic string) (*messaging.TopicManagementResponse, error) {
	response, err := c.client.SubscribeToTopic(ctx, tokens, topic)
	if err != nil {
		return nil, fmt.Errorf("error subscribing to topic: %v", err)
	}

	return response, nil
}

// * UnsubscribeFromTopic unsubscribes tokens from a topic
func (c *Client) UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) (*messaging.TopicManagementResponse, error) {
	response, err := c.client.UnsubscribeFromTopic(ctx, tokens, topic)
	if err != nil {
		return nil, fmt.Errorf("error unsubscribing from topic: %v", err)
	}

	return response, nil
}

// * ValidateToken validates if a token is valid
func (c *Client) ValidateToken(ctx context.Context, token string) error {
	// * Send a dry-run message to validate token
	message := &messaging.Message{
		Token: token,
		Data: map[string]string{
			"validation": "true",
		},
	}

	_, err := c.client.Send(ctx, message)
	return err
}

// * GetInvalidTokens filters out invalid tokens from a list
func (c *Client) GetInvalidTokens(ctx context.Context, response *messaging.BatchResponse) []string {
	var invalidTokens []string

	for i, resp := range response.Responses {
		if !resp.Success {
			// * Check if error is related to invalid token
			if messaging.IsInvalidArgument(resp.Error) ||
				messaging.IsRegistrationTokenNotRegistered(resp.Error) ||
				messaging.IsUnregistered(resp.Error) {
				if i < len(response.Responses) {
					invalidTokens = append(invalidTokens, fmt.Sprintf("token_index_%d", i))
				}
			}
		}
	}

	return invalidTokens
}

// * CreateCustomMessage creates a custom FCM message
func (c *Client) CreateCustomMessage(token string, data map[string]string, androidConfig *messaging.AndroidConfig, apnsConfig *messaging.APNSConfig) *messaging.Message {
	message := &messaging.Message{
		Token: token,
		Data:  data,
	}

	if androidConfig != nil {
		message.Android = androidConfig
	}

	if apnsConfig != nil {
		message.APNS = apnsConfig
	}

	return message
}

// * SendCustomMessage sends custom configured message
func (c *Client) SendCustomMessage(ctx context.Context, message *messaging.Message) (string, error) {
	response, err := c.client.Send(ctx, message)
	if err != nil {
		return "", fmt.Errorf("error sending custom message: %v", err)
	}

	return response, nil
}

// * BuildNotificationData builds data payload for notification
func BuildNotificationData(notificationID, userID, actionURL, notificationType string) map[string]string {
	data := map[string]string{
		"notification_id": notificationID,
		"user_id":         userID,
		"type":            notificationType,
		"click_action":    "FLUTTER_NOTIFICATION_CLICK",
	}

	if actionURL != "" {
		data["action_url"] = actionURL
	}

	return data
}
