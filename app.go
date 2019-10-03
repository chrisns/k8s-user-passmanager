package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"

	"github.com/creasty/defaults"
)

func main() {
	var secret string
	if os.Args[1:][0] == "keychain" {
		secret = keychainFetcher(os.Args[1:][1])
	}
	if os.Args[1:][0] == "1password" {
		secret = opgetter(os.Args[1:][1])
	}
	res := &response{}
	err := json.Unmarshal([]byte(secret), &res.Status)
	if err != nil {
		panic(err)
	}
	fmt.Println(formatResponse(res))
}

var defaultOp = func(itemName string) (*opResponse, error) {
	out, err := exec.Command("op", "get", "item", itemName).Output()
	if err != nil {
		return nil, err
	}
	var resp opResponse
	err = json.Unmarshal(out, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func opgetter(itemName string) string {
	resp, err := defaultOp(itemName)
	if err != nil {
		panic(err)
	}
	i := sort.Search(len(resp.Details.Fields), func(i int) bool { return resp.Details.Fields[i].Name == "password" })
	return resp.Details.Fields[i].Value
}

type opResponse struct {
	UUID    string            `json:"uuid"`
	Details opResponseDetails `json:"details"`
}
type opResponseDetails struct {
	Fields []opResponseField `json:"fields"`
	Title  string            `json:"title"`
}

type opResponseField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type responseStatus struct {
	Token                 string `default:"my-bearer-token" json:"token,omitempty"`
	ClientCertificateData string `json:"clientCertificateData,omitempty"`
	ClientKeyData         string `json:"clientKeyData,omitempty"`
}
type response struct {
	APIVersion string         `default:"client.authentication.k8s.io/v1beta1" json:"apiVersion"`
	Kind       string         `default:"ExecCredential" json:"kind"`
	Status     responseStatus `json:"status"`
}

func formatResponse(res *response) string {
	err := defaults.Set(res)
	if err != nil {
		panic(err)
	}
	jsonResponse, _ := json.Marshal(res)
	return string(jsonResponse)
}
