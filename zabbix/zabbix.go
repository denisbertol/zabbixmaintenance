package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// JSONRPCResponse - Estrutura de resposta Json do Zabbix
type JSONRPCResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	Error   Error       `json:"error"`
	Result  interface{} `json:"result"`
	ID      int         `json:"id"`
}

// JSONRPCRequest - Estrutura de request Json do Zabbix
type JSONRPCRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Auth    string      `json:"auth,omitempty"`
	ID      int         `json:"id"`
}

// Error - Estrutura de Erro do Zabbix
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (z *Error) Error() string {
	return z.Data
}

// API - Estrutura de dados da API
type API struct {
	url    string
	user   string
	passwd string
	ID     int
	auth   string
	Client *http.Client
}

// NewAPI - Funcao para criacao do objeto API
func NewAPI(server, user, passwd string) (*API, error) {
	return &API{server, user, passwd, 0, "", &http.Client{}}, nil
}

// GetAuth - Funcao para obter o metodo Auth da API apos o login
func (api *API) GetAuth() string {
	return api.auth
}

//ZabbixRequest - Trata as conexoes com o servidor
func (api *API) ZabbixRequest(method string, data interface{}) (JSONRPCResponse, error) {
	// Setup our JSONRPC Request data
	ID := api.ID
	api.ID = api.ID + 1
	jsonobj := JSONRPCRequest{"2.0", method, data, api.auth, ID}

	encoded, err := json.Marshal(jsonobj)

	if err != nil {
		return JSONRPCResponse{}, err
	}

	// Setup our HTTP request
	request, err := http.NewRequest("POST", api.url, bytes.NewBuffer(encoded))
	if err != nil {
		return JSONRPCResponse{}, err
	}
	request.Header.Add("Content-Type", "application/json-rpc")

	// Execute the request
	response, err := api.Client.Do(request)
	if err != nil {
		return JSONRPCResponse{}, err
	}

	var result JSONRPCResponse
	var buf bytes.Buffer

	_, err = io.Copy(&buf, response.Body)
	if err != nil {
		return JSONRPCResponse{}, err
	}

	json.Unmarshal(buf.Bytes(), &result)

	response.Body.Close()

	return result, nil
}

// Login - Funcao responsavel por efetuar o Login
func (api *API) Login() (bool, error) {
	params := make(map[string]string, 0)
	params["user"] = api.user
	params["password"] = api.passwd

	response, err := api.ZabbixRequest("user.login", params)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return false, err
	}

	if response.Error.Code != 0 {
		return false, &response.Error
	}

	api.auth = response.Result.(string)
	return true, nil
}

// Logout - Funcao responsavel por efetuar o logout
func (api *API) Logout() (bool, error) {
	emptyparams := make(map[string]string, 0)
	response, err := api.ZabbixRequest("user.logout", emptyparams)
	if err != nil {
		return false, err
	}

	if response.Error.Code != 0 {
		return false, &response.Error
	}

	return true, nil
}

// Version - Funcao responsavel por obter a versao da API
func (api *API) Version() (string, error) {
	response, err := api.ZabbixRequest("APIInfo.version", make(map[string]string, 0))
	if err != nil {
		return "", err
	}

	if response.Error.Code != 0 {
		return "", &response.Error
	}

	return response.Result.(string), nil
}
