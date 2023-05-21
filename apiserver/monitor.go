package apiserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

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

	adminPassword, err := apiconfig.SrvConfig.GetAdminPassword()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := ldaputils.Connect(apiconfig, "cn=admin,cn=config", adminPassword)
	if err != nil {
		log.Fatal(err)
	}

	//cn=Bind,cn=Operations,cn=Monitor
	bind, err := ldaputils.GetOne(conn, "cn=Bind,cn=Operations,cn=Monitor", "monitorOpInitiated", "monitorOpCompleted")
	if err != nil {
		log.Fatalf("Failed to search in monitor. %s", err)
	}

	type message struct {
		Time   int64    `json:"time"`
		Legend []string `json:"legend"`
		Value  []int64  `json:"value"`
	}

	var Message message
	var legends []string
	var values []int64

	now := time.Now().Unix()
	Message.Time = now

	for x, y := range bind {
		//Response.R = append(Response.R, message{x, y})
		legends = append(legends, x)
		yNum, err := strconv.ParseInt(y, 10, 64)
		if err != nil {
			log.Fatalf("Failed to convert to int64: %s. %s", y, err)
		}
		values = append(values, yNum)
	}
	Message.Legend = legends
	Message.Value = values
	responseJSON, _ := json.Marshal(Message)
	fmt.Printf("Message: %s\n", string(responseJSON))

	ctx.Header("Access-Control-Expose-Headers", "Content-Range")
	ctx.Header("Content-Range", "bytes 0-9/*")
	ctx.JSON(http.StatusOK, Message)
	conn.Close()
}
