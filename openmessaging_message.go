package purecloud

type OpenMessage struct {
	ID              string                `json:"id"`
	Channel         *OpenMessageChannel   `json:"channel"`
	Direction       string                `json:"direction"`
	Type            string                `json:"type"` // Text, Structured, Receipt
	Text            string                `json:"text"`
	Content         []*OpenMessageContent `json:"content,omitempty"`
	RelatedMessages []*OpenMessage        `json:"relatedMessages,omitempty"`
	Metadata        map[string]string     `json:"metadata,omitempty"`
}

type OpenMessageResult struct {
	OpenMessage
	Status         string                `json:"status,omitempty"`
	Reasons        []*StatusReason       `json:"reasons,omitempty"`
	Entity         string                `json:"originatingEntity"`
	IsFinalReceipt bool                  `json:"isFinalReceipt"`
}

type StatusReason struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

/*
{
	"id":"a47b309925689d96a00006c3fd1af102",
	"channel":{
		"id":"edce4efa-4abf-468b-ada7-cd6d35e7bbaf",
		"platform":"Open",
		"type":"Private",
		"to":{"id":"Uec880a5a4d3914c5816475534b62f31d"},
		"from":{"id":"edce4efa-4abf-468b-ada7-cd6d35e7bbaf"},
		"time":"2021-04-18T15:43:44.176Z"
	},
	"type":"Text",
	"content":[{
		"contentType":"Attachment",
		"attachment":{
			"mediaType":"Image",
			"url":"https://api.mypurecloud.jp/api/v2/downloads/pce61ceafa0","mime":"image/jpeg",
			"filename":"coruscant.jpg"
		},
		"reactions":[]
	}],
	"reasons":[],
	"direction":"Outbound",
	"relatedMessages":[]
} 
*/