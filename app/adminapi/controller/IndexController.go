package controller

import (
	"context"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql" // import your used driver
	"github.com/go-xorm/builder"
	"lazygo/app/adminapi/common"
	"lazygo/app/configurerpc/configureproto"
	"lazygo/common/rpcclient"
	"lazygo/core/apix"
	"lazygo/core/logx"
)

type (
	IndexController struct {
		apix.Controller
	}
	Auth struct {
		Id       int64  `db:"id" json:"u_id"`
		Mark     string `db:"username" json:"user_name"`
		Password string `db:"mark" json:"mobile_phone"`
	}
	helloRequest struct {
		Name string
		Age  int64
	}
)

func (c *IndexController) Hello() {
	cli, err := rpcclient.Next()
	if err != nil {
		c.W.Write([]byte("request error"))
	}
	resp, err := cli.Rpc(context.Background(), &configureproto.RpcRequest{
		Username: "dev",
		Password: "zyqQgFRULhxcAJs8",
		Name:     "crmuser-rpc",
		Env:      "dev",
		Network:  "inside"})
	if err != nil {
		c.W.Write([]byte("request error"))
	}
	c.W.Write([]byte(resp.Source))
}

func (c *IndexController) Mysql() {
	var auth Auth
	eq := builder.Eq{"username": "pro_xjy"}
	fields := []string{"id", "username", "mark"}
	err := builder.MySQL().
		Select(fields...).
		From("connect_auth").
		Where(eq).
		Limit(1, 0).
		QueryRow(&auth, common.Conn)
	if err != nil {
		logx.Error(err)
	}
	c.W.Write([]byte(auth.Mark))
}

func (c *IndexController) Redis() {
	// err := common.BizRedis.Set("tes","678")
	// 	if err != nil {
	// 	logx.Error(err)
	// }
	tes, err := common.BizRedis.Get("tes")
	if err != nil {
		logx.Error(err)
	}
	c.W.Write([]byte(tes))
}

func (c *IndexController) Hi() {
	var request helloRequest
	json.Unmarshal(c.RequestBody, &request)
	logx.Error(request)
	logx.Error(c.GetString("Name"))
	logx.Error(c.GetInt64("Age"))
	c.W.Write([]byte("hello"))
}
