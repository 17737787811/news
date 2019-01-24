package routers

import (
	"newsPublish/controllers"
	"github.com/astaxie/beego"
    "github.com/astaxie/beego/context"
)

func init() {

    beego.InsertFilter("/article/*",beego.BeforeExec,filterFUnc)
    beego.Router("/", &controllers.MainController{})
    beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleRegister")
    beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
    //首页
    beego.Router("/article/index",&controllers.ArticleController{},"get:ShowIndex")
    //添加文章
    beego.Router("/article/addArticle",&controllers.ArticleController{},"get:ShowAdd;post:HandleAdd")
    //查看文章详情
    beego.Router("/article/content",&controllers.ArticleController{},"get:ShowContent")
    //编辑文章
    beego.Router("/article/editArticle",&controllers.ArticleController{},"get:ShowEditArticle;post:HandleEditArticle")
    //删除文章
    beego.Router("/article/deleteArticle",&controllers.ArticleController{},"get:HandleDelete")
    //添加文章类型
    beego.Router("/article/addType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
    //退出登录
    beego.Router("/article/logout",&controllers.UserController{},"get:Logout")
    //删除类型
    beego.Router("/article/deleteType",&controllers.ArticleController{},"get:DeleteType")
    beego.Router("/redis",&controllers.RedisGit{},"get:Showredis")
}

func filterFUnc(ctx*context.Context){
    userName := ctx.Input.Session("userName")
    if userName == nil{
        ctx.Redirect(302,"/login")
        return
    }
}
