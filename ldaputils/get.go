package ldaputils

import (
	"log"

	"github.com/go-ldap/ldap"
)

func GetOne(conn *ldap.Conn, baseDN string, atributes ...string) (map[string]string, error) {

	filterDN := "*"

	result, err := conn.Search(ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filterDN,
		atributes,
		nil,
	))

	values := make(map[string]string)

	for _, entry := range result.Entries {
		for _, attr := range entry.Attributes {
			values[attr.Name] = attr.Values[0]
			log.Printf("Values: %s", values[attr.Name])
		}
	}

	if err != nil {
		return nil, err
	}
	return values, nil
}
