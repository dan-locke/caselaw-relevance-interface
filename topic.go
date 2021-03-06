package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"fmt"
)

type Topic struct {

	Collection string `json:"collection"`

	Id int `json:"id"`

	Topic string `json:"topic"`

	CaseId int `json:"case_id"`

	CaseTitle string `json:"case_title"`

	DateFiled string `json:"date_filed"`

	Html string `json:"html"`

	PlainText string `json:"plain_text"`

	Extracts []extract`json:"extracts"`
}

type extract struct {

	CitedId []int `json:"cited_id"`

	IdFound []bool `json:"id_manually_found"`

	CitedCase []string `json:"cited_case"`

	CaseExtract []string `json:"case_extract"`

	CitingSentence string `json:"citing_sentence"`

	CitingParagraph string `json:"citing_paragraph"`

	RelevantKeywords []string `json:"relevant_keywords"`

	Query []string `json:"query"`

	EsQuery []map[string]interface{} `json:"es_query"`
	
}

type TopicIndex []struct {

	Topic string

	Name string

	Assessed int

}

func saveTopics(fileName string, t map[string]Topic) error {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)

	err := enc.Encode(t)
	if err != nil {
		log.Panic(err)
		return err
	}

	return ioutil.WriteFile(fileName, buff.Bytes(), 0664)
}

func loadFromDatFile(fileName string) (*map[string]Topic, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	topics := make(map[string]Topic)
	dec := gob.NewDecoder(bytes.NewReader(data))
	err = dec.Decode(&topics)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	return &topics, nil
}

func loadfromFolder(dir, path string) (*map[string]Topic, error) {
	err := os.Chdir(path)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	fileInfo, err := ioutil.ReadDir(".")
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	topics := make(map[string]Topic)
	for _, i := range fileInfo {
		if strings.Contains(i.Name(), ".json") {
			data, err := ioutil.ReadFile(i.Name())
			if err != nil {
				log.Panic(err)
				return nil, err
			}
			t := Topic{}
			err = json.Unmarshal(data, &t)
			if err != nil {
				log.Panic(err)
				return nil, err
			}
			topics[strconv.Itoa(t.Id)] = t
		}
	}

	err = os.Chdir(dir)
	return &topics, err
}

func loadTopics(dir, path, fileName string, load, update bool) (*map[string]Topic, error) {
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
	if load {
		m, err := loadFromDatFile(fileName)
		if os.IsNotExist(err) {
			update = true
		} else {
			return m, err
		}
	}

	if path == "" {
		update = true
	}

	if update {
		m, err := loadfromFolder(dir, path)
		if err != nil {
			return nil, err
		}
		err = saveTopics(fileName, *m)
		return m, err
	}
	return nil, nil
}

// TODO need to add here upper number to be assessed....
func (i *Instance) getTopicList(user int64) (TopicIndex, error) {
	l := TopicIndex{}

	assessed, err := dbGetNumberAssessedPerTopic(i.db, user)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	for k, v := range i.topics {
		l = append(l, struct{Topic string; Name string; Assessed int}{k, v.CaseTitle, assessed[k]})
	}
	return l, nil
}

func (i *Instance) getTopic(topic string) Topic {
	return i.topics[topic]
}

func (i *Instance) updateTopic(topic string, t Topic) error {
	i.topics[topic] = t
	// lazy approach, easier to resave whole file
	return saveTopics(i.config.Topics.DataFileName, i.topics)
}



func loadTopicIdFile(path string) (map[string][]string, error) {
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	topics := map[string][]string{}

	lines := bytes.Split(buff, []byte("\n"))
	for i := range lines {
		if len(lines[i]) == 0 {
			break
		}
		parts := bytes.Fields(lines[i])
		// cnt := 0
		if len(parts) != 2 {
			return nil, fmt.Errorf("Incorrect number of items in line: %q\n", lines[i])
		}
		//id, err := strconv.Atoi(string(parts[0]))
		if err != nil {
			return nil, err
		}

		//topics[id] = append(topics[id], string(parts[1]))
		topics[string(parts[0])] = append(topics[string(parts[0])], string(parts[1]))
	}

	return topics, nil
}

func loadQrel(path string) (map[string][]string, error) {
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	topics := map[string][]string{}

	lines := bytes.Split(buff, []byte("\n"))
	for i := range lines {
		parts := bytes.Fields(lines[i])
		// cnt := 0
		if len(parts) != 4 {
			return nil, fmt.Errorf("Incorrect number of items in line: %q\n", lines[i])
		}
		//id, err := strconv.Atoi(string(parts[0]))
		if err != nil {
			return nil, err
		}

		//topics[id] = append(topics[id], string(parts[1]))
		topics[string(parts[0])] = append(topics[string(parts[0])], string(parts[2]))
	}

	return topics, nil
}
