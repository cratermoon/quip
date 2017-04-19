package twit

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ChimeraCoder/anaconda"

	"github.com/cratermoon/quip/storage"
)

type Twit struct {
	api *anaconda.TwitterApi
	kit *storage.Kit
}

type keys struct {
	Key         string `json:"key"`
	Secret      string `json:"secret"`
	Token       string `json:"token"`
	TokenSecret string `json:"tokenSecret"`
}

func NewTwit() *Twit {
	kit, err := storage.NewKit()
	if err != nil {
		log.Println("Error starting quip service", err)
		return nil
	}
	raw, err := kit.FileObject("keys.json")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	var k keys
	json.Unmarshal(raw, &k)

	anaconda.SetConsumerKey(k.Key)
	anaconda.SetConsumerSecret(k.Secret)
	api := anaconda.NewTwitterApi(k.Token, k.TokenSecret)
	t := Twit{api: api}
	return &t
}

// Tweet posts the given status string to twitter
func (t *Twit) Tweet(status string) (int64, error) {
	tweet, err := t.api.PostTweet(status, nil)
	return tweet.Id, err
}

// Delete removes the status with the given id from the twitter timeline
func (t *Twit) Delete(id int64) (int64, error) {
	tweet, err := t.api.DeleteTweet(id, false)
	return tweet.Id, err
}
