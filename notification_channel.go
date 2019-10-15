package purecloud

// NotificationChannel  defines a Notification Channel
type NotificationChannel struct {
	ID           string  `json:"channelId"`
	WebsocketURL string  `json:"wssUrl"`
	SelfURI      string  `json:"selfURI"`
	Version      uint32  `json:"version"`
	Client       *Client `json:"-"`
}

type Subscriber interface {
	ActivityID() string
}

// CreateNotificationChannel creates a new channel for notifications
func (client *Client) CreateNotificationChannel() (*NotificationChannel, error) {
	channel := &NotificationChannel{}
	if err := client.Post("/notifications/channels", struct{}{}, &channel); err != nil {
		return nil, err
	}
	channel.Client = client
	return channel, nil
}

func (channel *NotificationChannel) Subscribe(subscriber Subscriber) error {
	return channel.Client.Post(
		"/notifications/channels/"+channel.ID+"/subscriptions",
		struct {
			ID string `json:"id"`
		}{ID: subscriber.ActivityID()},
		nil,
	)
}
