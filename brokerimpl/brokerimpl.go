package brokerimpl

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/broersa/ttnrouter/broker"
)

type (
	BrokerImpl struct {
	}
)

// New Implemented Factory
func New() broker.Broker {
	return &BrokerImpl{}
}

func (brokerimpl *BrokerImpl) FindBrokerOnAppEUI(appeui []byte, brokers []string) (string, error) {
	proc := func(ep string, c chan<- string) {
		resp, err := http.Get(ep + "/HasApplication/test") //string(appeui))
		if err != nil {
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		if string(body) == "OK" {
			c <- ep
		}
	}
	c := make(chan string)
	for _, value := range brokers {
		go proc(value, c)
	}
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()
	select {
	case endpoint := <-c:
		return endpoint, nil
	case <-timeout:
	}
	return "", errors.New("Timeout on broker search")
}

func (brokerimpl *BrokerImpl) FindBrokerOnDevAddr(devaddr []byte, brokers []string) (string, error) {
	return "http://localhost:4443/Message", nil
}

func (brokerimpl *BrokerImpl) ForwardMessage(endpoint string, message *broker.Message) (*broker.ResponseMessage, error) {
	j, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(endpoint+"/Message", "application/json", bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var responsemessage broker.ResponseMessage
	err = json.Unmarshal(body, &responsemessage)
	if err != nil {
		return nil, err
	}
	return &responsemessage, nil
}
