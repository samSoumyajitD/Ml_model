package datacollection

import (
	"net/http"

	"github.com/ciromacedo/nwdaf/model"
	"github.com/ciromacedo/nwdaf/util"
	"github.com/free5gc/openapi"
	"github.com/gin-gonic/gin"
)

func HTTPAmfRegistrationAccept(c *gin.Context) {
	var amfeventnotifylist model.AmfEventnotifylist
	requestBody, err := c.GetRawData()
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadGateway)
		c.Writer.Write([]byte("Internal Error"))
		return
	}

	err = openapi.Deserialize(&amfeventnotifylist, requestBody, "application/json")
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadGateway)
		c.Writer.Write([]byte("Json Parser Error"))
		return
	}

	//amfeventnotifylist.Date = time.Now()
	/* registrar na base */
	util.AddRegistrationAccept(&amfeventnotifylist)
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte("Ok"))
}
