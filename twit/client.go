package twit

import (

	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ChimeraCoder/anaconda"
)

type Twit struct {
	api *anaconda.TwitterApi
}

type keys struct {
	Key         string `json:"key"`
	Secret      string `json:"secret"`
	Token       string `json:"token"`
	TokenSecret string `json:"tokenSecret"`
}

func NewTwit() *Twit {
	raw, err := ioutil.ReadFile("./keys.json")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	var k keys
	json.Unmarshal(raw, &k)


	anaconda.SetConsumerKey(k.Key)
	anaconda.SetConsumerSecret(k.Secret)
	api := anaconda.NewTwitterApi(k.Token, k.TokenSecret)
	t := Twit{api}
	return &t
}

func (t *Twit) tweet(status string) (int64, error) {
	tweet, err := t.api.PostTweet(status, nil)
	return tweet.Id, err
}
