package model

import "time"

type AmfEventnotifylist struct {
	Subscription struct {
		Eventlist []struct {
			Type          string `json:"type"`
			Immediateflag bool   `json:"immediateFlag"`
		} `json:"eventlist"`

		EventNotifyUri      string `json:"eventNotifyUri"`
		NotifyCorrelationId string `json:"notifyCorrelationId"`
		NfId                string `json:"nfId"`
		AnyUE               bool   `json:"anyUE"`
	} `json:"subscription"`
	SubscriptionId string `json:"subscriptionId"`
	Reportlist     []struct {
		Type  string `json:"type"`
		State struct {
			Active bool `json:"active"`
		} `json:"state"`
		TimeStamp time.Time `json:"timeStamp"`
		AnyUE     bool      `json:"anyUe"`
		Supi      string    `json:"supi"`
	} `json:"reportList"`
}

type RegistrationAccept struct {
	Date   time.Time `json:"date"`
	Amf    Amf
	Ue     Ue
	PlmnId PlmnId
}

type Amf struct {
	Id     string `json:"id"`
	Locale string `json:"locale"`
}

type Ue struct {
	Suci string `json:"suci"`
	Supi string `json:""`
}

type PlmnId struct {
	Mcc string `json:"mcc"`
	Mnc string `json:"mnc"`
}

/* CONFIG */
type Config struct {
	Port     int
	MongoURI string
	DBName   string
}

type CollectionInfo struct {
	DocumentName    string `json:"Name"`
	NumberOfRecords int64
}

type Article struct {
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}
