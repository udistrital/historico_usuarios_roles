package controllers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/beego/beego/logs"
	"github.com/udistrital/usuario_rol_crud/models"
	"github.com/udistrital/utils_oas/time_bogota"

	"github.com/astaxie/beego"
)

// SistemaInformacionController operations for SistemaInformacion
type SistemaInformacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *SistemaInformacionController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create SistemaInformacion
// @Param	body		body 	models.SistemaInformacion	true		"body for SistemaInformacion content"
// @Success 201 {int} models.SistemaInformacion
// @Failure 403 body is empty
// @router / [post]
func (c *SistemaInformacionController) Post() {
	var v models.SistemaInformacion

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		v.Activo = true
		v.FechaCreacion = time_bogota.TiempoBogotaFormato()
		v.FechaModificacion = time_bogota.TiempoBogotaFormato()
		if _, err := models.AddSistemaInformacion(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "registration successful", "Data": v}
		} else {
			logs.Error(err)
			c.Data["Message"] = "Error service Post: the reques contain an incorrect parameter or no record exists"
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		c.Data["Message"] = "Error service Post: the reques contain an incorrect parameter or no record exists"
		c.Abort("400")
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get SistemaInformacion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.SistemaInformacion
// @Failure 403 :id is empty
// @router /:id [get]
func (c *SistemaInformacionController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetSistemaInformacionById(id)
	if err != nil {
		logs.Error(err)
		c.Data["Message"] = "Error service GetOne: the reques contain an incorrect parameter or no record exists"
		c.Abort("404")
	} else {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "request successful", "Data": v}
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get SistemaInformacion
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.SistemaInformacion
// @Failure 403
// @router / [get]
func (c *SistemaInformacionController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllSistemaInformacion(query, fields, sortby, order, offset, limit)
	if err != nil {
		logs.Error(err)
		c.Data["Message"] = "Error service GetAll: the reques contain an incorrect parameter or no record exists"
		c.Abort("404")
	} else {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "request successful", "Data": l}

	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the SistemaInformacion
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.SistemaInformacion	true		"body for SistemaInformacion content"
// @Success 200 {object} models.SistemaInformacion
// @Failure 403 :id is not int
// @router /:id [put]
func (c *SistemaInformacionController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)

	v := models.SistemaInformacion{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		//se recupera sistema existente para mantener fecha de creacion
		sistema, err := models.GetUsuarioById(id)
		if err != nil {
			logs.Error(err)
			c.Data["Message"] = "Error service Put: the reques contain an incorrect data type or an invalid parameter"
			c.Abort("400")
			return
		}
		v.Activo = true
		v.FechaCreacion = time_bogota.TiempoCorreccionFormato(sistema.FechaCreacion)
		v.FechaModificacion = time_bogota.TiempoBogotaFormato()
		if err := models.UpdateSistemaInformacionById(&v); err == nil {
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "update successful", "Data": v}
		} else {
			logs.Error(err)
			c.Data["Message"] = "Error service Put: the reques contain an incorrect data type or an invalid parameter"
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		c.Data["Message"] = "Error service Put: the reques contain an incorrect data type or an invalid parameter"
		c.Abort("400")
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the SistemaInformacion
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *SistemaInformacionController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteSistemaInformacion(id); err == nil {
		d := map[string]interface{}{"Id": id}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "delete successful", "Data": d}
	} else {
		logs.Error(err)
		c.Data["Message"] = "Error service Delete: the reques contain an incorrect data type or an invalid parameter"
		c.Abort("400")
	}
	c.ServeJSON()
}
