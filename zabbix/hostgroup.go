package zabbix

import (
	"encoding/json"
	"strings"
)

// HostGroup - Estrutura do Zabbix Host
type HostGroup map[string]interface{}

// HostGroup - Funcao responsavel por executar operacao no objeto hostgroup.* do zabbix
func (api *API) HostGroup(method string, data interface{}) ([]HostGroup, error) {
	response, err := api.ZabbixRequest("hostgroup."+method, data)
	if err != nil {
		return nil, err
	}

	if response.Error.Code != 0 {
		return nil, &response.Error
	}

	res, err := json.Marshal(response.Result)
	var ret []HostGroup
	err = json.Unmarshal(res, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// GetHostGroups - Obtem um objeto host pelo nome
func GetHostGroups(api *API, hostgroups string) ([]HostGroup, error) {

	var HostGroups []HostGroup
	// first, clean/remove the comma
	strLimpa := strings.Replace(hostgroups, ",", " ", -1)
	// convert 'clened' comma separated string to slice
	hostGroupsArray := strings.Fields(strLimpa)

	params := make(map[string]interface{}, 0)
	filter := make(map[string]interface{}, 0)
	filter["name"] = hostGroupsArray
	params["filter"] = filter
	params["output"] = "extend"
	HostGroups, err := api.HostGroup("get", params)
	if err != nil {
		return nil, err
	}

	// Se algum host nao foi encontrado, retorna erro
	if len(HostGroups) < len(hostGroupsArray) {
		return nil, &Error{0, "", "ERRO: Um ou mais HOSTGROUP nao foi Encontrado"}
	}

	if len(HostGroups) > 0 {
		return HostGroups, nil
	}

	return nil, &Error{0, "", "ERRO: Um ou mais HOSTGROUP nao foi Encontrado"}
}

// GetAllHostGroups - Obtem um objeto host pelo nome
func GetAllHostGroups(api *API) ([]HostGroup, error) {

	params := make(map[string]interface{}, 0)
	params["output"] = "extend"
	HostGroups, err := api.HostGroup("get", params)
	if err != nil {
		return nil, err
	}

	if len(HostGroups) > 0 {
		return HostGroups, nil
	}

	return nil, &Error{0, "", "ERRO: Um ou mais HOSTGROUP nao foi Encontrado"}
}
