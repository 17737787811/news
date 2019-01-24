package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
)
type RedisGit struct {
	beego.Controller
}
func (this*RedisGit)Showredis(){

	conn,err:=redis.Dial("tcp",":6379")
	if err!=nil{
		beego.Error("连接错误",err)
		return
	}
	resp,err:=conn.Do("mget","kk","bb","aa")
result,_:=redis.Values(resp,err)

    var v1,v2 string
	var v3 int
	redis.Scan(result,&v1,&v2,&v3)
	beego.Info(v1,v2,v3)

}