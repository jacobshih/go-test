package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
)

type KnowledgeBase struct {
	kbmap map[string]interface{}
}

var kb KnowledgeBase

func jsonTest() {
	const jsonStream = `
	[
		{"Name": "Ed", "Text": "Knock knock."},
		{"Name": "Sam", "Text": "Who's there?"},
		{"Name": "Ed", "Text": "Go fmt."},
		{"Name": "Sam", "Text": "Go fmt who?"},
		{"Name": "Ed", "Text": "Go fmt yourself!"}
	]
`
	type Message struct {
		Name, Text string
	}
	dec := json.NewDecoder(strings.NewReader(jsonStream))

	// read open bracket
	t, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%T: %v\n", t, t)

	// while the array contains values
	for dec.More() {
		var m Message
		// decode an array value (Message)
		err := dec.Decode(&m)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v: %v\n", m.Name, m.Text)
	}

	// read closing bracket
	t, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%T: %v\n", t, t)

	fmt.Printf("%s\n", "jsonTest() end.")
}

func stripPrefix(text string, prefix string) string {
	if strings.HasPrefix(text, prefix) {
		text = strings.TrimPrefix(text, prefix)
		text = strings.TrimLeft(text, " ")
	}
	return text
}

func stringCompareTest() {
	var msg1 = "#echo hello world"
	var msg2 = "#hello hello nothing"
	var id = "0123456789"
	var msg, text string
	const prefixEcho = "#echo"
	const prefixHello = "#hello"
	text = msg1
	if strings.HasPrefix(text, prefixEcho) {
		msg = stripPrefix(text, prefixEcho)
		msg = id + " : " + msg + " (jxlbot echo test)"
		fmt.Printf("process '%v' : %v\n", msg1, msg)
	}
	text = msg2
	if strings.HasPrefix(text, prefixHello) {
		msg = stripPrefix(text, prefixHello)
		msg = id + " : " + msg + " (jxlbot hello test)"
		fmt.Printf("process '%v' : %v\n", msg2, msg)
	}
}

func traceType(o interface{}) {
	fmt.Println("reflect.TypeOf(o): ", reflect.TypeOf(o))
}

func (kb *KnowledgeBase) init() {
	kb.kbmap = make(map[string]interface{})
}

func (kb *KnowledgeBase) learn(jsonString string) {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonString), &m)
	if err != nil {
		panic(err)
	}
	kb.merge(m, kb.kbmap)
}

func (kb *KnowledgeBase) merge(src, dst map[string]interface{}) {
	for k, v := range src {
		_, ok := dst[k].(map[string]interface{})
		if ok {
			kb.merge(src[k].(map[string]interface{}), dst[k].(map[string]interface{}))
		} else {
			_, ok := dst[k].([]interface{})
			if ok {
				dst[k] = append(dst[k].([]interface{}), v.([]interface{})...)
			} else {
				dst[k] = v
			}
		}
	}
}

func (kb *KnowledgeBase) find(data map[string]interface{}, s ...string) interface{} {
	if val, ok := data[s[0]]; ok {
		if reflect.TypeOf(val).Kind() == reflect.String {
			return val
		} else if reflect.TypeOf(val).Kind() == reflect.Slice {
			// FIXME: find in slice
			return val
		} else {
			nextLevel := s[1:]
			if len(nextLevel) > 0 {
				return kb.find(val.(map[string]interface{}), nextLevel...)
			}
			return val
		}
	}
	return nil
}

func (kb *KnowledgeBase) dump() {
	kb.toJSON(kb.kbmap)
}

func (kb *KnowledgeBase) toJSON(m map[string]interface{}) {
	str, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(str))
}

const jsonStream1 = `{
	"group01": {
		"1": {
			"name": "name01",
			"email": "name01@abc.com"
		},
		"2": {
			"name": "name02",
			"email": "name02@abc.com"
		},
		"event": [
			{
			
			}
		]
	}
}
`
const jsonStream2 = `{
	"group01": {
		"1": {
			"nickname": "nickname01"
		},
		"3": {
			"name": "name03",
			"email": "name03@abc.com"
		},
		"event": [
			{
			"20180425":
				{
				"date": "20180425",
				"place": "where"
				}
			}
		]
	}
}
`
const jsonStream3 = `{
	"group01": {
		"1": {
			"description": {
				"company": {
					"name": "company01",
					"address": "address01",
					"telphone": "01234567891"
				}
			}
		},
		"3": {
			"name": "name003",
			"description": {
				"company": {
					"name": "company03",
					"address": "address03",
					"telphone": "01234567893"
				}
			}
		},
		"event": [
			{
			"20180426":
				{
				"date": "20180426",
				"place": "resturant0426"
				}
			}
		]
	}
}
`
const jsonStream4 = `{
	"group01": {
		"4": {
			"description": {
				"company": {
					"name": "company04",
					"address": "address04",
					"telphone": "01234567894"
				}
			}
		},
		"event": [
			{
			"20180501":
				{
				"date": "20180501",
				"place": "space"
				}
			}
		]
	}
}
`
const anotherJson = `{
	"another01": {
		"event": [
			{
			"20180428":
				{
				"date": "20180428",
				"who": ["you","me"],
				"place": "ponorogo"
				}
			}
		]
	}
}
`
const jsonKeyValueOnly = `{
	"mykey": "myvalue"
}
`

func main() {
	//	jsonTest()
	//	stringCompareTest()

	// knowledge base init
	kb := KnowledgeBase{}
	kb.init()

	// learn
	kb.learn(jsonStream1)
	kb.learn(jsonStream2)
	kb.learn(jsonStream3)
	kb.learn(jsonStream4)
	kb.learn(anotherJson)
	kb.learn(jsonKeyValueOnly)
	kb.dump()

	// find
	group1name1 := kb.find(kb.kbmap, "group01", "1", "name")
	fmt.Println("==========================")
	fmt.Println("group1name1: ", group1name1)

	event1 := kb.find(kb.kbmap, "group01", "event", "20180428")
	fmt.Println("==========================")
	fmt.Println("event1: ")
	for _, v := range event1.([]interface{}) {
		kb.toJSON(v.(map[string]interface{}))
	}
}
