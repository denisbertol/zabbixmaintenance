package zabbix

import (
	"encoding/json"
	"fmt"
	"time"
)

// ZabbixMaintenance - Variavel de dados do Zabbix Maintenance
var ZabbixMaintenance map[string]interface{}

// Maintenance - Funcao responsavel por executar operacao no objeto maintenance.* do zabbix
func (api *API) Maintenance(method string, data interface{}) (map[string]interface{}, error) {
	response, err := api.ZabbixRequest("maintenance."+method, data)
	if err != nil {
		return nil, err
	}

	if response.Error.Code != 0 {
		return nil, &response.Error
	}

	res, err := json.Marshal(response.Result)
	err = json.Unmarshal(res, &ZabbixMaintenance)
	if err != nil {
		return nil, err
	}

	return ZabbixMaintenance, nil
}

//CreateMaintenance - Funcao para criacao de manutencao
func CreateMaintenance(api *API, nome string, Hosts []Host, HostGroups []HostGroup, horas int64) (map[string]interface{}, error) {
	params := make(map[string]interface{}, 0)
	timeperiod := make(map[string]interface{}, 0)
	tags := make(map[string]interface{}, 0)

	params["name"] = fmt.Sprintf("%s%s%d", nome, " - ", time.Now().Unix())
	params["active_since"] = time.Now().Unix()

	horasemseg := (horas * 3600)

	params["active_till"] = time.Now().Unix() + horasemseg
	params["tags_evaltype"] = 0

	if Hosts != nil {
		// Criar Array com os hostsids
		var hostsids []string
		for i := 0; i < len(Hosts); i++ {
			var Host Host
			Host = Hosts[i]
			hostsids = append(hostsids, fmt.Sprintf("%v", Host["hostid"]))
		}

		params["hostids"] = hostsids
	} else if HostGroups != nil {
		// Criar Array com os hostsids
		var hostgroupsids []string
		for i := 0; i < len(HostGroups); i++ {
			var HostGroup HostGroup
			HostGroup = HostGroups[i]
			hostgroupsids = append(hostgroupsids, fmt.Sprintf("%v", HostGroup["groupid"]))
		}

		params["groupids"] = hostgroupsids
	}
	params["tags"] = tags

	var timeperiods [1]map[string]interface{}
	timeperiod["timeperiod_type"] = 0
	timeperiod["start_date"] = time.Now().Unix()
	timeperiod["period"] = horasemseg

	timeperiods[0] = timeperiod

	params["timeperiods"] = timeperiods
	ret, err := api.Maintenance("create", params)
	if err != nil {
		return nil, err
	}

	// Se houve sucesso na chamada
	if len(ret) > 0 {
		return ret, err
	}

	return nil, &Error{0, "", "Erro ao criar Manutencao"}
}
