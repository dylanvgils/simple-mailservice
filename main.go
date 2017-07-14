package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type response struct {
	Writer  http.ResponseWriter `json:"-"`
	Code    int                 `json:"code,int"`
	Message string              `json:"message,string"`
}

func (resp *response) sendResponse() {
	resp.Writer.Header().Set("Content-Type", "application/json")

	switch resp.Code {
	case 400:
		resp.Writer.WriteHeader(http.StatusBadRequest)
	default:
		resp.Writer.WriteHeader(http.StatusOK)
	}

	json, _ := json.Marshal(resp)
	resp.Writer.Write(json)
}

func newResponse(w http.ResponseWriter) *response {
	return &response{
		Writer: w,
		Code:    200,
		Message: "Success",
	}
}

func readBody(req *http.Request) []byte {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil
	}

	return body
}

func handleRequest(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	resp := newResponse(w)
	m := newMail(req.RemoteAddr)

	err := json.Unmarshal(readBody(req), &m)
	if err != nil {
		resp.Code = 400
		resp.Message = "Error while parsing json"
	}

	m.sendMail()
	resp.sendResponse()
}

func main() {
	router := httprouter.New()
	router.POST("/", handleRequest)
	log.Fatalln(http.ListenAndServe(":3000", router))
}
