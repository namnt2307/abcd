package module

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

func GetAPI(url string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Sentry_log(err)
		return nil, err
	}

	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		Sentry_log(err)
		return nil, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body, nil
}

func PostAPIJson(url string, data interface{}) ([]byte, error) {
	dataByte, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataByte))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	return body, nil
}

func PostAPI(url string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		Sentry_log(err)
		return nil, err
	}

	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		Sentry_log(err)
		return nil, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body, nil
}

func DeleteAPI(url string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		Sentry_log(err)
		return nil, err
	}

	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		Sentry_log(err)
		return nil, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body, nil
}

func PostAPIJsonTonken(url string, data interface{}, token string) ([]byte, error) {
	dataByte, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataByte))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		Sentry_log(err)
		return nil, err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	// log.Println(string(body))
	return body, nil
}
