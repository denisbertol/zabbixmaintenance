package main

import (
	"bufio"
	"fmt"
	"manutencao/zabbix"
	"os"
	"regexp"
	"strconv"
)

func main() {
	// Cria uma nova API
	api, err := zabbix.NewAPI("http://192.168.15.18/api_jsonrpc.php", "admin", "zabbix")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Efetua o Login
	_, err = api.Login()
	if err != nil {
		fmt.Println(err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	nome := ""
	for nome == "" {
		// Solicita o nome da manutencao na console
		fmt.Print("\nNome da manutencao: ")
		scanner.Scan()
		nome = scanner.Text()
	}

	opcaohost := ""
	for opcaohost != "1" && opcaohost != "2" {
		// Solicita o(s) nome(s) do(s) host(s), ou (s) nome(s) do(s) hostgroup(s)
		fmt.Println("\nAtrelar a manutencao por: ")
		fmt.Println("(1) Nomes de Hosts")
		fmt.Println("(2) Nomes de Host Groups")
		scanner.Scan()
		opcaohost = scanner.Text()
	}

	// Solicita o nome dos hosts ou dos hostgroups
	zabbixHosts := ""
	if opcaohost == "1" {
		for zabbixHosts == "" {
			fmt.Print("\nEntre com o(s) Nome(s) do(s) Host(s) separado por virgula. <ENTER> para listar todos: ")
			scanner.Scan()
			zabbixHosts = scanner.Text()
			if zabbixHosts == "" {
				// Obtem todos os hosts
				hostsresult, err := zabbix.GetAllHosts(api)
				if err != nil {
					fmt.Println(err)
					return
				}
				for i := 0; i < len(hostsresult); i++ {
					var host zabbix.Host
					host = hostsresult[i]
					fmt.Println(fmt.Sprintf("%v", host["host"]))
				}
			}
		}
	} else if opcaohost == "2" {
		for zabbixHosts == "" {
			fmt.Print("\nEntre com o(s) Nome(s) do(s) HostGroups(s) separado por virgula. <ENTER> para listar todos: ")
			scanner.Scan()
			zabbixHosts = scanner.Text()
			if zabbixHosts == "" {
				// Obtem todos os hostsgroups
				hostgroupsresult, err := zabbix.GetAllHostGroups(api)
				if err != nil {
					fmt.Println(err)
					return
				}
				for i := 0; i < len(hostgroupsresult); i++ {
					var hostGroup zabbix.HostGroup
					hostGroup = hostgroupsresult[i]
					fmt.Println(fmt.Sprintf("%v", hostGroup["name"]))
				}
			}
		}
	}

	horas := ""
	horasvalidas := false
	for horasvalidas == false {
		// Solicita o periodo em horas da manutencao a partir do horario atual
		fmt.Print("\nHoras em Manutencao: ")
		scanner.Scan()
		horas = scanner.Text()
		horasvalidas, _ = regexp.MatchString("[0-9][0-9]", horas)
	}

	// Converte o hostid para Int64
	horasint, err := strconv.ParseInt(horas, 10, 64)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Recupera o id dos hosts ou host groups
	if opcaohost == "1" {
		// Obtem os hosts pelo(s) nome(s)
		hostresult, err := zabbix.GetHosts(api, zabbixHosts)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Cria Manutencao pelo host
		maintenceresult, err := zabbix.CreateMaintenance(api, nome, hostresult, nil, horasint)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("\nManutencao Criada com Sucesso. Id da Mantencao: ", maintenceresult["maintenanceids"])
		}
	} else if opcaohost == "2" {
		// Obtem os hosts pelo(s) nome(s)
		hostgroupresult, err := zabbix.GetHostGroups(api, zabbixHosts)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Cria Manutencao pelo hostgroup
		maintenceresult, err := zabbix.CreateMaintenance(api, nome, nil, hostgroupresult, horasint)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("\nManutencao Criada com Sucesso. Id da Mantencao: ", maintenceresult["maintenanceids"])
		}
	}

	// Efetua o Logout
	_, err = api.Logout()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("\nOperacao Concluida. Pressione qualquer tecla para Sair")
	var input string
	fmt.Scanln(&input)

}
