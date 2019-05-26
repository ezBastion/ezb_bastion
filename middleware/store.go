// This file is part of ezBastion.

//     ezBastion is free software: you can redistribute it and/or modify
//     it under the terms of the GNU Affero General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.

//     ezBastion is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU Affero General Public License for more details.

//     You should have received a copy of the GNU Affero General Public License
//     along with ezBastion.  If not, see <https://www.gnu.org/licenses/>.

package middleware

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ezbastion/ezb_srv/cache"
	"github.com/ezbastion/ezb_srv/model"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Store(storage cache.Storage, conf *model.Configuration) gin.HandlerFunc {
	return func(c *gin.Context) {
		

		tr, _ := c.Get("trace")
		trace := tr.(model.EzbLogs)
		logg := log.WithFields(log.Fields{
			"middleware": "store",
			"xtrack":     trace.Xtrack,
		})
		routeType, _ := c.MustGet("routeType").(string)
		if routeType == "worker" {
			usr, _ := c.Get("account")
			account := usr.(model.EzbAccounts)
			viewApi, err := model.GetViewApi(storage, conf, account.Name, trace.Xtrack)
			if err != nil {
				logg.Error(err)
				c.AbortWithError(http.StatusInternalServerError, errors.New("#S0001"))
				return
			}
			if len(viewApi) == 0 {
				logg.Error("NO API FOUND FOR ", account.Name)
				c.AbortWithError(http.StatusForbidden, errors.New("#S0002"))
				return
			}
			c.Set("ViewApi", viewApi)

			apiPath, err := model.GetApiPath(storage, conf)
			if err != nil {
				logg.Error(err)
				c.AbortWithError(http.StatusInternalServerError, errors.New("#S0003"))
				return
			}
			if len(apiPath) == 0 {
				logg.Error("NO API FOUND")
				c.AbortWithError(http.StatusForbidden, errors.New("#S0004"))
				return
			}
			c.Set("apiPath", apiPath)
			
		}
		if routeType == "internal" {
			w, _ := c.MustGet("wksid").(string)
			wksid, _ := strconv.Atoi(w)
			workers, err := model.GetWorkers(storage, conf)
			var worker model.EzbWorkers
			if err != nil {
				logg.Error(err)
				c.AbortWithError(http.StatusInternalServerError, errors.New("#S0005"))
				return
			}
			for _, w := range workers {
				if w.ID == wksid {
					worker = w
					break
				}
			}
			if worker.ID == 0 {
				logg.Error("NO WORKER FOUND")
				c.AbortWithError(http.StatusForbidden, errors.New("#S0006"))
				return
			}
			trace.Worker = worker.Name
			trace.Controller = "internal"
			c.Set("trace", trace)
			c.Set("worker", worker)
		}
		c.Next()
	}
}


