package rpc

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	consts "github.com/YaleOpenLab/openx/consts"
)

// CARPC contains a list of all ComplyAdvantage related RPCs
var CARPC = map[int][]string{
	1: {"/user/ca/search", "name", "birthyear"},
	2: {"/admin/ca/users/all"},
}

// setupCAHandlers sets up rpc handlers that are involved with integrating ComplyAdvantage into openx
func setupCAHandlers() {
	searchComplyAdvantage()
	getAllCAUsers()
}

// CAResponse defines a struct that ComplyAdvantage returns
type CAResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"string"`
	Content struct {
		Data struct {
			ID         int64  `json:"id"`
			Ref        string `json:"ref"`
			Searcherid int64  `json:"searcher_id"`
			Assigneeid int64  `json:"assignee_id"`
			Filters    struct {
				Birthyear      int64    `json:"birth_year"`
				Countrycodes   []int    `json:"country_codes"`
				Removedeceased int      `json:"remove_deceased"`
				Types          []string `json:"types"`
				Exactmatch     bool     `json:"exact_match"`
				Fuzziness      float64  `json:"fuzziness"`
			}

			Matchstatus   string   `json:"match_status"`
			Risklevel     string   `json:"risk_level"`
			Searchterm    string   `json:"search_term"`
			Submittedterm string   `json:"submitted_term"`
			Clientref     string   `json:"client_ref"`
			Totalhits     int      `json:"total_hits"`
			Updatedat     string   `json:"updated_at"`
			Createdat     string   `json:"created_at"`
			Tags          []string `json:"tags"`
			Limit         int      `json:"limit"`
			Offset        int      `json:"offset"`
			Shareurl      string   `json:"share_url"`
			Hits          []struct {
				Doc struct {
					Aka []struct {
						Name string `json:"name"`
					} `json:"aka"`
					Assets []struct {
						Publicurl string `json:"public_url"`
						Source    string `json:"source"`
						Type      string `json:"type"`
					} `json:"assets"`
					Entitytype string `json:"entity_type"`
					Fields     []struct {
						Name   string `json:"name"`
						Source string `json:"source"`
						Tag    string `json:"tag"`
						Value  string `json:"value"`
					} `json:"fields"`
					ID    string
					Media []struct {
						Date    string `json:"date"`
						Snippet string `json:"snippet"`
						Title   string `json:"title"`
						URL     string `json:"url"`
					} `json:"media"`
					Name    string   `json:"name"`
					Sources []string `json:"sources"`
					Types   []string `json:"types"`
				} `json:"doc"`
				Matchtypes    []string `json:"match_types"`
				Score         float64  `json:"score"`
				Matchstatus   string   `json:"match_status"`
				Iswhitelisted bool     `json:"is_whitelisted"`
			} `json:"hits"`
		} `json:"data"`
	} `json:"content"`
}

// PostRequestCA is a handler that makes it easy to send out POST requests
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

	defer func() {
		if ferr := res.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	x, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("did not read from ioutil", err)
		return dummy, err
	}

	return x, nil
}

// PostAndSendCA is a handler that POSTs data and returns the response
func PostAndSendCA(w http.ResponseWriter, r *http.Request, body string, payload io.Reader) {
	data, err := PostRequestCA(body, payload)
	if erpc.Err(w, err, erpc.StatusBadRequest, "did not receive success response") {
		return
	}
	log.Println(string(data))
	var x kycDepositResponse
	err = json.Unmarshal(data, &x)
	if erpc.Err(w, err, erpc.StatusInternalServerError, "did not unmarshal json") {
		return
	}
	erpc.MarshalSend(w, x)
}

// searchComplyAdvantage searches for a particular entity on ComplyAdvantage's platform
func searchComplyAdvantage() {
	http.HandleFunc(CARPC[1][0], func(w http.ResponseWriter, r *http.Request) {
		_, err := userValidateHelper(w, r, CARPC[1][1:], "GET")
		if err != nil {
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
		// also decide what should be done with the specific user and what message must be displayed on the frontend
		PostAndSendCA(w, r, body, payload)
	})
}

type caAllUserResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Content struct {
		Data []struct {
			ID        int    `json:"id"`
			Email     string `json:"email"`
			Name      string `json:"name"`
			Phone     string `json:"phone"`
			Updatedat string `json:"updated_at"`
			Createdat string `json:"created_at"`
		} `json:"data"`
	} `json:"content"`
}

// getAllCAUsers gets a list of all users searched for using ComplyAdvantage
func getAllCAUsers() {
	http.HandleFunc(CARPC[2][0], func(w http.ResponseWriter, r *http.Request) {
		_, err := userValidateHelper(w, r, CARPC[2][1:], "GET")
		if err != nil {
			return
		}

		body := "https://api.complyadvantage.com/users?api_key=" + consts.KYCAPIKey
		data, err := erpc.GetRequest(body)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		var x caAllUserResponse
		err = json.Unmarshal(data, &x)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}
		erpc.MarshalSend(w, x)
	})
}
