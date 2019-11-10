package webssh

import (
	"net/http"

	"github.com/it234/gowebssh/internal/app/manageweb/controllers/common"
	models "github.com/it234/gowebssh/internal/pkg/models/common"
	"github.com/it234/gowebssh/internal/pkg/models/webssh"
	"github.com/it234/gowebssh/pkg/convert"

	"github.com/gin-gonic/gin"
)

type ServerInfo struct{}

// 分页数据
func (ServerInfo) List(c *gin.Context) {
	page := common.GetPageIndex(c)
	limit := common.GetPageLimit(c)
	sort := common.GetPageSort(c)
	key := common.GetPageKey(c)
	var whereOrder []models.PageWhereOrder
	order := "ID DESC"
	if len(sort) >= 2 {
		orderType := sort[0:1]
		order = sort[1:len(sort)]
		if orderType == "+" {
			order += " ASC"
		} else {
			order += " DESC"
		}
	}
	whereOrder = append(whereOrder, models.PageWhereOrder{Order: order})
	if key != "" {
		v := "%" + key + "%"
		var arr []interface{}
		arr = append(arr, v)
		whereOrder = append(whereOrder, models.PageWhereOrder{Where: "alias_name like", Value: arr})
	}
	var total uint64
	list:= []webssh.ServerInfo{}
	err := models.GetPage(&webssh.ServerInfo{}, &webssh.ServerInfo{}, &list, page, limit, &total, whereOrder...)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	common.ResSuccessPage(c, total, &list)
}

// 详情
func (ServerInfo) Detail(c *gin.Context) {
	id := common.GetQueryToUint64(c, "id")
	var model webssh.ServerInfo
	where := webssh.ServerInfo{}
	where.ID = id
	_, err := models.First(&where, &model)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	model.Password = ""
	common.ResSuccess(c, &model)
}

// 更新
func (ServerInfo) Update(c *gin.Context) {
	model := webssh.ServerInfo{}
	err := c.Bind(&model)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	where := webssh.ServerInfo{}
	where.ID = model.ID
	modelOld := webssh.ServerInfo{}
	_, err = models.First(&where, &modelOld)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	if model.Password=="" {
		model.Password = modelOld.Password
	}
	err = models.Save(&model)
	if err != nil {
		common.ResFail(c, "操作失败")
		return
	}
	common.ResSuccessMsg(c)
}

//新增
func (ServerInfo) Create(c *gin.Context) {
	model := webssh.ServerInfo{}
	err := c.Bind(&model)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	err = models.Create(&model)
	if err != nil {
		common.ResFail(c, "操作失败")
		return
	}
	common.ResSuccess(c, gin.H{"id": model.ID})
}

// 删除数据
func (ServerInfo) Delete(c *gin.Context) {
	var ids []uint64
	err := c.Bind(&ids)
	if err != nil || len(ids) == 0 {
		common.ResErrSrv(c, err)
		return
	}
  _,err = models.DeleteByIDS(webssh.ServerInfo{},ids)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	common.ResSuccessMsg(c)
}

func (ServerInfo) Xterm(c *gin.Context) {
	sid,_:=c.GetQuery("sid")
	c.HTML(http.StatusOK, "xterm.html",gin.H{
		"sid": sid,
	})
}

func (ServerInfo) Ws(c *gin.Context) {
	WebSocketHandler(c.Writer,c.Request,checkUserToken,getServerInfo)
}

func checkUserToken(token string) bool {
	return true
}

func getServerInfo(sid string) (m SshLoginModel) {
	var model webssh.ServerInfo
	where := webssh.ServerInfo{}
	where.ID = convert.ToUint64(sid)
	_, err := models.First(&where, &model)
	if err != nil {
		return
	}
	m.PtyCols=100
	m.PtyRows=100
	m.Addr=model.ServerAddress
	m.UserName=model.UserName
	m.Pwd=model.Password
	return 
}

