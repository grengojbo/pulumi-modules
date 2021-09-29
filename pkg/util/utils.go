package util

import (
	"github.com/grengojbo/pulumi-modules/interfaces"
	// "github.com/grengojbo/pulumi-modules/automation"
	// "github.com/grengojbo/pulumi-modules/config"
	log "github.com/sirupsen/logrus"
	// "github.com/spf13/cobra"
)

// ListSecurityGroup подготовка cписка групп безопасности
func ListSecurityGroup(sg interfaces.SecurityGroupArgs) (sgList []interfaces.PortSecurityGroupArgs) {
	// sgList := []interfaces.PortSecurityGroupArgs{}
	if sg.Http {
		i := interfaces.PortSecurityGroupArgs{
			Name: "HTTP",
			Allows: []string{"0.0.0.0/0"},
			Port: 80,
		}
		sgList = append(sgList, i)
	}
	if sg.Https {
		i := interfaces.PortSecurityGroupArgs{
			Name: "HTTPS",
			Allows: []string{"0.0.0.0/0"},
			Port: 443,
		}
		sgList = append(sgList, i)
	}
	if len(sg.SshAllows) > 0 {
		for _, v := range sg.SshAllows {
			i := interfaces.PortSecurityGroupArgs{
				Name: "SSH",
				Allows: []string{v},
				Port: 22,
			}
			sgList = append(sgList, i)
		}
	}
	if len(sg.KubeApi) > 0 {
		for _, v := range sg.KubeApi {
			i := interfaces.PortSecurityGroupArgs{
				Name: "KubeApi",
				Allows: []string{v},
				Port: 6443,
			}
			sgList = append(sgList, i)
		}
	}
	
	if len(sg.Vault) > 0 {
		for _, v := range sg.Vault {
			i := interfaces.PortSecurityGroupArgs{
				Name: "Vault",
				Allows: []string{v},
				Port: 8200,
			}
			sgList = append(sgList, i)
		}
	}

	if len(sg.PostgreSqlAllows) > 0 {
		for _, v := range sg.PostgreSqlAllows {
			i := interfaces.PortSecurityGroupArgs{
				Name: "PostgreSql",
				Allows: []string{v},
				Port: 5432,
			}
			sgList = append(sgList, i)
		}
	}

	if len(sg.MySqlAllows) > 0 {
		for _, v := range sg.MySqlAllows {
			i := interfaces.PortSecurityGroupArgs{
				Name: "MySql",
				Allows: []string{v},
				Port: 3306,
			}
			sgList = append(sgList, i)
		}
	}
	if len(sg.MongoAllows) > 0 {
		for _, v := range sg.MongoAllows {
			i := interfaces.PortSecurityGroupArgs{
				Name: "MongoDB",
				Allows: []string{v},
				Port: 27017,
			}
			sgList = append(sgList, i)
		}
	}

	if len(sg.MsSqlAllows) > 0 {
		for _, v := range sg.MsSqlAllows {
			i := interfaces.PortSecurityGroupArgs{
				Name: "MsSql",
				Allows: []string{v},
				Port: 1433,
			}
			sgList = append(sgList, i)
		}
	}

	if len(sg.Custom) > 0 {
		for _, row := range sg.Custom {
			for _, v := range row.Allows {
				i := interfaces.PortSecurityGroupArgs{
					Name: row.Name,
					Allows: []string{v},
					Port: row.Port,
				}
				sgList = append(sgList, i)
			}
		}
	}

	for _, row := range sgList {
		// log.Infof("Name: %s, Port: %d, Allow: %s", row.Name, row.Port, row.Allows[0])
		log.Debugf("Name: %s, Port: %d, Allow: %s", row.Name, row.Port, row.Allows[0])
	}
	return sgList
}