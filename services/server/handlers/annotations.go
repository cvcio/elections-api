package handlers

import (
	"github.com/cvcio/elections-api/models/annotation"
	"github.com/cvcio/elections-api/pkg/config"
	"github.com/cvcio/elections-api/pkg/db"
	"github.com/gin-gonic/gin"
)

// Annotations Controller
type Annotations struct {
	cfg *config.Config
	db  *db.DB
}

// Create New Annotation
func (ctrl *Annotations) Create(c *gin.Context) {
	var a annotation.Annotation
	if err := c.Bind(&a); err != nil {
		ResponseError(c, 406, err.Error())
		return
	}
	res, err := annotation.Create(ctrl.db, &a)
	if err != nil {
		ResponseError(c, 401, err.Error())
		return
	}
	ResponseJSON(c, res.IDStr)
}
