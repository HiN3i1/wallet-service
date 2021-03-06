// Copyright (c) 2018-2021 The CYBAVO developers
// All Rights Reserved.
// NOTICE: All information contained herein is, and remains
// the property of CYBAVO and its suppliers,
// if any. The intellectual and technical concepts contained
// herein are proprietary to CYBAVO
// Dissemination of this information or reproduction of this materia
// is strictly forbidden unless prior written permission is obtained
// from CYBAVO.

package cybavo

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var baseURL = "https://sofatest.sandbox.cybavo.com"

func buildChecksum(params []string, secret string, time int64, r string) string {
	params = append(params, fmt.Sprintf("t=%d", time))
	params = append(params, fmt.Sprintf("r=%s", r))
	sort.Strings(params)
	params = append(params, fmt.Sprintf("secret=%s", secret))
	return fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(params, "&"))))
}

func MakeRequest(apiCode string, apiSecret string, method string, api string, params []string, ginBody io.ReadCloser) ([]byte, error) {

	var err error

	postBody, err := ioutil.ReadAll(ginBody)

	if err != nil {
		return nil, errors.New("invalid parameters")
	}

	if method == "" || api == "" {
		return nil, errors.New("invalid parameters")
	}

	client := &http.Client{}
	r := RandomString(8)
	if r == "" {
		return nil, errors.New("can't generate random byte string")
	}
	t := time.Now().Unix()
	url := fmt.Sprintf("%s%s?t=%d&r=%s", baseURL, api, t, r)
	if len(params) > 0 {
		url += fmt.Sprintf("&%s", strings.Join(params, "&"))
	}

	var req *http.Request
	if len(postBody) == 0 {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, bytes.NewReader(postBody))
		params = append(params, string(postBody))
	}
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-CODE", apiCode)
	req.Header.Set("X-CHECKSUM", buildChecksum(params, apiSecret, t, r))
	if postBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	log.Debug("Request URL:", url)
	log.Debug("\tX-CHECKSUM:\t", req.Header.Get("X-CHECKSUM"))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		result := &ErrorCodeResponse{}
		_ = json.Unmarshal(body, result)
		msg := fmt.Sprintf("%s, Error: %s", res.Status, result.String())
		return body, errors.New(msg)
	}
	return body, nil
}
