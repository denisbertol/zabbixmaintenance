package zabbix

import (
	"strings"
	"encoding/json"
	"fmt"
	"time"
	"strconv"
)

// ZabbixMaintenance - Variavel de dados do Zabbix Maintenance
var ZabbixMaintenance map[string]interface{}

// ZabbixMaintenance - Variavel de dados do Zabbix Maintenance
type Maintenance map[string]interface{}

// Maintenance - Funcao responsavel por executar operacao no objeto maintenance.* do zabbix
func (api *API) MaintenanceVar(method string, data interface{}) (map[string]interface{}, error) {
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

// Maintenance - Funcao responsavel por executar operacao no objeto maintenance.* do zabbix
func (api *API) Maintenance(method string, data interface{}) ([]Maintenance, error) {
	response, err := api.ZabbixRequest("maintenance."+method, data)
	if err != nil {
		return nil, err
	}

	if response.Error.Code != 0 {
		return nil, &response.Error
	}
	res, err := json.Marshal(response.Result)
	var ret []Maintenance
	err = json.Unmarshal(res, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

//CreateMaintenance - Funcao para criacao de manutencao
func CreateMaintenance(api *API, nome string, Hosts []Host, HostGroups []HostGroup, horas int64) (map[string]interface{}, error) {
	params := make(map[string]interface{}, 0)
	timeperiod := make(map[string]interface{}, 0)
	tags := make(map[string]interface{}, 0)

	// Nome da Manutencao com o timestamp para não ter duplicidade
	params["name"] = fmt.Sprintf("%s%s%d", nome, " - ", time.Now().Unix())
	// Horario da ativacao da manutencao
	params["active_since"] = time.Now().Unix()	
	// Horario da finalizacao da ativacao
	horasemseg := (horas * 3600)
	params["active_till"] = time.Now().Unix() + horasemseg
	// Tipo da Manutencao 1 - para nao coletar metricas no periodo
	params["maintenance_type"] = 1

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
	ret, err := api.MaintenanceVar("create", params)
	if err != nil {
		return nil, err
	}

	// Se houve sucesso na chamada
	if len(ret) > 0 {
		return ret, err
	}

	return nil, &Error{0, "", "Erro ao criar Manutencao"}
}

// GetMaintence - Obtem um objeto host pelo nome
func GetMaintenance(api *API, maintenanceid string) ([]Maintenance, error) {

	var Maintenances []Maintenance

	params := make(map[string]interface{}, 0)
	maintenanceidArray := strings.Fields(maintenanceid)

	params["maintenanceids"] = maintenanceidArray
	params["selectTimeperiods"] = "extend"
	params["output"] = "extend"

	Maintenances, err := api.Maintenance("get", params)
	if err != nil {
		return nil, err
	}

	if len(Maintenances) > 0 {
		return Maintenances, nil
	}

	return nil, &Error{0, "", "ERRO: Uma ou mais MAINTENANCE nao foi Encontrada"}
}

//UpdateMaintenance - Funcao para criacao de manutencao
func UpdateMaintenance(api *API, maintenanceid string, horariofinaltp int64) (map[string]interface{}, error) {

	mainantenceresult, err := GetMaintenance(api, maintenanceid)
	if err != nil {
		fmt.Println(err)
	}

	params := make(map[string]interface{}, 0)
	timeperiod := make(map[string]interface{}, 0)

	// Id da Manutencao a ser alterada
	params["maintenanceid"] = maintenanceid
	// Horario da finalizacao da ativacao
	params["active_till"] = horariofinaltp
	
	horainicio, err := strconv.ParseInt(fmt.Sprintf("%v",mainantenceresult[0]["active_since"]), 10, 64)

	periodo := horariofinaltp - horainicio

	var timeperiods [1]map[string]interface{}
	timeperiod["timeperiod_type"] = 0
	timeperiod["start_date"] = mainantenceresult[0]["active_since"]
	timeperiod["period"] = periodo

	timeperiods[0] = timeperiod

	params["timeperiods"] = timeperiods
	ret, err := api.MaintenanceVar("update", params)
	if err != nil {
		return nil, err
	}

	// Se houve sucesso na chamada
	if len(ret) > 0 {
		return ret, err
	}

	return nil, &Error{0, "", "Erro ao criar Manutencao"}
}
