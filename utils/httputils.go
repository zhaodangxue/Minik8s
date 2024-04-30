package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func ReadUnmarshal(rc io.ReadCloser, v interface{}) error {
	value, err := io.ReadAll(rc)
	defer rc.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(value, v)
	if err != nil {
		return err
	}
	return nil
}
func PostWithJson(url string, v interface{}) (*http.Response, error) {
	client := http.Client{}
	value, _ := json.Marshal(v)
	req, err := http.NewRequest("POST", url, bytes.NewReader(value))
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}
func PostWithString(url string, value string) (*http.Response, error) {
	client := http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader([]byte(value)))
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}
func PostWithForm(url string, form map[string][]string) (*http.Response, error) {
	client := http.Client{}
	value, _ := json.Marshal(form)
	reader := bytes.NewReader(value)
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}
func GetUnmarshal(url string, v interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	value, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	return json.Unmarshal(value, v)
}
func PutWithJson(url string, v interface{}) (string, error) {
	client := http.Client{}
	value, _ := json.Marshal(v)
	req, err := http.NewRequest("PUT", url, bytes.NewReader(value))
	if err != nil {
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	value, err = io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return "", err
	}
	return string(value), nil

}
func Delete(url string) (string, error) {
	client := http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	value, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return "", err
	}
	return string(value), nil
}
