package apiserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aescanero/openldap-node/config"
	"github.com/aescanero/openldap-node/ldaputils"
	"github.com/gin-gonic/gin"
)

/*currentconnections: current established connections
totalconnections: total established connections
dncache: total DN in cache
entrycache: total entries in cache
idlcache: total IDL in cache
totaloperations: total operations
totalabandon: total ABANDON operation
totaladd: total ADD operations
totalbind: total BIND operations
totalcompare: total COMPARE operations
totaldelete: total DELETE operations
totalextended: total EXTENDED operations
totalmodify: total MODIFY operations
totalmodrdn: total MODRDN operations
totalsearch: total SEARCH operations
totalunbind: total UNBIND operations*/

func monitor(ctx *gin.Context, apiconfig config.Config) {

	conn, err := ldaputils.Connect(apiconfig)
	if err != nil {
		log.Fatal(err)
	}

	//cn=Bind,cn=Operations,cn=Monitor
	bind, err := ldaputils.GetOne(conn, "cn=Bind,cn=Operations,cn=Monitor", "monitorOpInitiated", "monitorOpCompleted")
	if err != nil {
		log.Fatalf("Failed to search in monitor. %s", err)
	}

	message, _ := json.Marshal(bind)
	fmt.Println(string(message))

	ctx.JSON(http.StatusOK, gin.H{"message": message})
}
