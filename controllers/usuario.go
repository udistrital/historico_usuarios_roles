package controllers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/usuario_rol_crud/helpers"
	"github.com/udistrital/usuario_rol_crud/models"
	"github.com/udistrital/usuario_rol_crud/services"

	"github.com/astaxie/beego"
)

// UsuarioController operations for Usuario
type UsuarioController struct {
	beego.Controller
}

// URLMapping ...
func (c *UsuarioController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Usuario
// @Param	body		body 	models.Usuario	true		"body for Usuario content"
// @Success 201 {int} models.Usuario
// @Failure 403 body is empty
// @router / [post]
func (c *UsuarioController) Post() {
	defer helpers.ErrorController(c.Controller, "UsuarioController")

	var v models.Usuario
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		logs.Error(err)
		c.Data["Message"] = "Error servicio Post: la petición contiene un parámetro incorrecto o no existe ningún registro"
		c.Abort("400")
		return
	}

	if _, err := services.AddUsuario(&v); err != nil {
		logs.Error(err)
		c.Data["Message"] = "Error servicio Post: la petición contiene un parámetro incorrecto o no existe ningún registro"
		c.Abort("400")
		return
	}

	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Registro exitoso", "Data": v}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Usuario by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Usuario
// @Failure 403 :id is empty
// @router /:id [get]
func (c *UsuarioController) GetOne() {
	//defer helpers.ErrorController(c.Controller, "UsuarioController")

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := services.GetUsuarioById(id)
	if err != nil {
		logs.Error(err)
		c.Data["Message"] = "Error en el servicio GetOne: la solicitud contiene un parámetro incorrecto o no existe ningún registro."
		c.Abort("404")
	} else {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Petición exitosa", "Data": v}
	}
	c.ServeJSON()
}

// GetPeriodosByDocumento ...
// @Title Get One Periodo por Documento
// @Description get Usuario by documento and periodos
// @Param	documento		path 	string	true		"The key for staticblock"
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Usuario
// @Failure 403 :id is empty
// @router /:documento/periodos [get]
func (c *UsuarioController) GetPeriodosByDocumento() {
	documento := c.Ctx.Input.Param(":documento")
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
	var result interface{}
	var err error
	var count int64
	if sistemaId, ok := query["sistema_informacion"]; ok {
		// Llamar al servicio de SistemaInformacion
		result, count, err = services.PeriodosPorSistema(&documento, sistemaId, query, fields, sortby, order, offset, limit)
	} else {
		// Llamar al servicio estándar
		result, count, err = services.GetPeriodosPorDocumento(documento, query, fields, sortby, order, offset, limit)
	}
	//periodosUsuario, err := services.GetPeriodosPorDocumento(documento, query, fields, sortby, order, offset, limit)
	if err != nil {
		logs.Error(err)
		c.Data["Message"] = "Error en la consulta de los periodos."
		c.Abort("404")
	} else {
		c.Data["json"] = map[string]interface{}{"Count": count, "Success": true, "Status": 200, "Message": "Petición exitosa", "Data": result}
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Usuario
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Usuario
// @Failure 403
// @router / [get]
func (c *UsuarioController) GetAll() {
	defer helpers.ErrorController(c.Controller, "UsuarioController")

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

	l, err := services.GetAllUsuario(query, fields, sortby, order, offset, limit)
	if err != nil {
		logs.Error(err)
		c.Data["Message"] = "Error servicio GetAll: la solicitud contiene un parámetro incorrecto o no existe ningún registro."
		c.Abort("404")
	} else {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Petición exitosa", "Data": l}
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Usuario
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Usuario	true		"body for Usuario content"
// @Success 200 {object} models.Usuario
// @Failure 403 :id is not int
// @router /:id [put]
func (c *UsuarioController) Put() {
	defer helpers.ErrorController(c.Controller, "UsuarioController")

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Usuario{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		if err := services.UpdateUsuarioById(&v); err == nil {
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Modificacion exitosa", "Data": v}
		} else {
			logs.Error(err)
			c.Data["Message"] = "Error servicio Put: la solicitud contiene un tipo de datos incorrecto o un parámetro no válido"
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		c.Data["Message"] = "Error servicio Put: la solicitud contiene un tipo de datos incorrecto o un parámetro no válido"
		c.Abort("400")
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Usuario
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *UsuarioController) Delete() {
	defer helpers.ErrorController(c.Controller, "UsuarioController")

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := services.DeleteUsuario(id); err == nil {
		d := map[string]interface{}{"Id": id}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Eliminacion exitosa", "Data": d}
	} else {
		logs.Error(err)
		c.Data["Message"] = "Error servicio Delete: la solicitud contiene un tipo de datos incorrecto o un parámetro no válido"
		c.Abort("400")
	}
	c.ServeJSON()
}
