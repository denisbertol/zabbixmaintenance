package zabbix

import (
	"encoding/json"
)

// Item - Estrutura do Zabbix Item
type Item map[string]interface{}

// Item - Funcao responsavel por executar operacao no objeto item.* do zabbix
func (api *API) Item(method string, data interface{}) ([]Item, error) {
	response, err := api.ZabbixRequest("item."+method, data)
	if err != nil {
		return nil, err
	}

	if response.Error.Code != 0 {
		return nil, &response.Error
	}

	res, err := json.Marshal(response.Result)
	var ret []Item
	err = json.Unmarshal(res, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

//UpdateItem - Funcao para atualizacao de Item
func UpdateItem(api *API, maintenanceid string, horariofinaltp int64) (map[string]interface{}, error) {
/*
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
*/
	return nil, &Error{0, "", "Erro ao criar Manutencao"}
}

