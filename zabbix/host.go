package zabbix

import (
	"encoding/json"
	"strings"
)

// Host - Estrutura do Zabbix Host
type Host map[string]interface{}

// Host - Funcao responsavel por executar operacao no objeto host.* do zabbix
func (api *API) Host(method string, data interface{}) ([]Host, error) {
	response, err := api.ZabbixRequest("host."+method, data)
	if err != nil {
		return nil, err
	}

	if response.Error.Code != 0 {
		return nil, &response.Error
	}

	res, err := json.Marshal(response.Result)
	var ret []Host
	err = json.Unmarshal(res, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// GetHosts - Obtem um objeto host pelo nome
func GetHosts(api *API, hosts string) ([]Host, error) {

	var Hosts []Host
	// first, clean/remove the comma
	strLimpa := strings.Replace(hosts, ",", " ", -1)
	// convert 'clened' comma separated string to slice
	hostsArray := strings.Fields(strLimpa)

	params := make(map[string]interface{}, 0)
	filter := make(map[string]interface{}, 0)
	filter["host"] = hostsArray
	params["filter"] = filter
	params["output"] = "extend"
	params["select_groups"] = "extend"
	params["templated_hosts"] = 1
	Hosts, err := api.Host("get", params)
	if err != nil {
		return nil, err
	}

	// Se algum host nao foi encontrado, retorna erro
	if len(Hosts) < len(hostsArray) {
		return nil, &Error{0, "", "ERRO: Um ou mais HOST nao foi Encontrado"}
	}

	if len(Hosts) > 0 {
		return Hosts, nil
	}

	return nil, &Error{0, "", "ERRO: Um ou mais HOST nao foi Encontrado"}
}

// GetAllHosts - Obtem todos os hosts
func GetAllHosts(api *API) ([]Host, error) {

	params := make(map[string]interface{}, 0)
	params["output"] = "extend"
	params["select_groups"] = "extend"
	params["templated_hosts"] = 1
	Hosts, err := api.Host("get", params)
	if err != nil {
		return nil, err
	}

	if len(Hosts) > 0 {
		return Hosts, nil
	}

	return nil, &Error{0, "", "ERRO: Nenhum HOST Encontrado"}
}
