package rpc

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	consts "github.com/YaleOpenLab/openx/consts"
)

func setupCAHandlers() {
	searchComplyAdvantage()
	getAllCAUsers()
}

type CAResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"string"`
	Content struct {
		Data struct {
			Id          int64  `json:"id"`
			Ref         string `json:"ref"`
			Searcher_id int64  `json:"searcher_id"`
			Assignee_id int64  `json:"assignee_id"`
			Filters     struct {
				birth_year      int64    `json:"birth_year"`
				Country_codes   []int    `json:"country_codes"`
				Remove_deceased int      `json:"remove_deceased"`
				Types           []string `json:"types"`
				Exact_match     bool     `json:"exact_match"`
				Fuzziness       float64  `json:"fuzziness"`
			}

			Match_status   string   `json:"match_status"`
			Risk_level     string   `json:"risk_level"`
			Search_term    string   `json:"search_term"`
			Submitted_term string   `json:"submitted_term"`
			Client_ref     string   `json:"client_ref"`
			Total_hits     int      `json:"total_hits"`
			Updated_at     string   `json:"updated_at"`
			Created_at     string   `json:"created_at"`
			Tags           []string `json:"tags"`
			Limit          int      `json:"limit"`
			Offset         int      `json:"offset"`
			Share_url      string   `json:"share_url"`
			Hits           []struct {
				Doc struct {
					Aka []struct {
						Name string `json:"name"`
					} `json:"aka"`
					Assets []struct {
						Public_url string `json:"public_url"`
						Source     string `json:"source"`
						Type       string `json:"type"`
					} `json:"assets"`
					Entity_type string `json:entity_type`
					Fields      []struct {
						Name   string `json:"name"`
						Source string `json:"source"`
						Tag    string `json:"tag"`
						Value  string `json:"value"`
					} `json:"fields"`
					Id    string
					Media []struct {
						Date    string `json:"date"`
						Snippet string `json:"snippet"`
						Title   string `json:"title"`
						Url     string `json:"url"`
					} `json:"media"`
					Name    string   `json:"name"`
					Sources []string `json:"sources"`
					Types   []string `json:"types"`
				} `json:"doc"`
				Match_types    []string `json:"match_types"`
				Score          float64  `json:"score"`
				Match_status   string   `json:"match_status"`
				Is_whitelisted bool     `json:"is_whitelisted"`
			} `json:"hits"`
		} `json:"data"`
	} `json:"content"`
}

// PostRequest is a handler that makes it easy to send out POST requests
func PostRequestCA(body string, payload io.Reader) ([]byte, error) {
	// the body must be the param that you usually pass to curl's -d option
	var dummy []byte
	req, err := http.NewRequest("POST", body, payload)
	if err != nil {
		log.Println("did not create new POST request", err)
		return dummy, err
	}

	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("did not make request", err)
		return dummy, err
	}

	defer res.Body.Close()
	x, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("did not read from ioutil", err)
		return dummy, err
	}

	return x, nil
}

// PostAndSend is a handler that POSTs data and returns the response
func PostAndSendCA(w http.ResponseWriter, r *http.Request, body string, payload io.Reader) {
	data, err := PostRequestCA(body, payload)
	if err != nil {
		log.Println("did not receive success response", err)
		responseHandler(w, StatusBadRequest)
		return
	}
	log.Println(string(data))
	var x kycDepositResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println("did not unmarshal json", err)
		responseHandler(w, StatusInternalServerError)
		return
	}
	MarshalSend(w, x)
}

func searchComplyAdvantage() {
	http.HandleFunc("/user/ca/search", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		_, err := UserValidateHelper(w, r)
		if err != nil {
			responseHandler(w, StatusUnauthorized)
			return
		}

		if r.URL.Query()["name"] == nil || r.URL.Query()["birthyear"] == nil {
			responseHandler(w, StatusBadRequest)
			return
		}

		name := r.URL.Query()["name"][0]
		birthyear := r.URL.Query()["birthyear"][0]
		body := "https://api.complyadvantage.com/searches?api_key=" + consts.KYCAPIKey
		data := `{
  "search_term": "` + name + `",
  "client_ref": "testnet+ElChapo",
  "fuzziness": 0.8,
  "filters": {
      "birth_year": "` + birthyear + `"` + `
  },
	"limit": 5,
  "share_url": 1
}`
		payload := bytes.NewBuffer([]byte(data))
		// TODO: analyze the response and check whether the user is clear or not. If not,
		// also decide what should be done with the specific user and what message must be dipslayed on the frontend
		PostAndSendCA(w, r, body, payload)
	})
}

type caAllUserResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Content struct {
		Data []struct {
			Id         int    `json:"id"`
			Email      string `json:"email"`
			Name       string `json:"name"`
			Phone      string `json:"phone"`
			Updated_at string `json:"updated_at"`
			Created_at string `json:"created_at"`
		} `json:"data"`
	} `json:"content"`
}

func getAllCAUsers() {
	http.HandleFunc("/admin/ca/users/all", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		_, err := UserValidateHelper(w, r)
		if err != nil {
			responseHandler(w, StatusUnauthorized)
			return
		}

		body := "https://api.complyadvantage.com/users?api_key=" + consts.KYCAPIKey
		data, err := GetRequest(body)
		if err != nil {
			responseHandler(w, StatusInternalServerError)
			return
		}

		var x caAllUserResponse
		err = x.UnmarshalJSON(data)
		if err != nil {
			responseHandler(w, StatusInternalServerError)
			return
		}
		MarshalSend(w, x)
	})
}
