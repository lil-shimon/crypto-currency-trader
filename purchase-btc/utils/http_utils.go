package utils

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

func CreateHttpRequest(method, url string, header, query map[string]string, data []byte) ([]byte, error) {

	/// GET POSTしか受け付けない
	if method != "GET" && method != "POST" {
		return nil, errors.New("method was neither GET or POST")
	}

	/// HTTP requestを作成
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	/// Query Parameterを作成
	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}

	/// エンコード
	req.URL.RawQuery = q.Encode()

	/// Headerの処理
	for key, value := range header {
		req.Header.Add(key, value)
	}

	/// HTTP Client
	httpC := &http.Client{}

	/// requestを投げる
	res, err := httpC.Do(req)
	if err != nil {
		return nil, err
	}

	/// 開放処理
	defer res.Body.Close()

	/// resの中身を取得
	bd, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bd, nil
}
