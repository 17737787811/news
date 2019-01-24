package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"newsPublish/models"
	"math"
	"github.com/gomodule/redigo/redis"
	//"database/sql/driver"
	"bytes"
	"encoding/gob"

)

type ArticleController struct {
	beego.Controller
}


//展示首页内容
func(this*ArticleController)ShowIndex(){

	//登录校验
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login",302)
		return
	}

	//获取所有文章数据
	//获取orm对象
	o := orm.NewOrm()
	//获取所有文章    select * from article;  queryseter
	qs :=o.QueryTable("Article")
	var articles []models.Article
	//qs.All(&articles)
	//实现分页
	//获取总记录数和总页数

	pageSize := 2


	//处理首页末页内容
	pageIndex,err :=this.GetInt("pageIndex")
	if err != nil{
		pageIndex = 1
	}

	//获取数据库部分数据  获取几条数据    从哪里开始
	start := pageSize * (pageIndex - 1)

	//获取所有类型数据并传递给前段展示
	var articleTypes []models.ArticleType
	//存到redis中
	conn,err:=redis.Dial("tcp",":6379")
	if err != nil{
		beego.Error("获取失败",err)
		return
	}
	defer conn.Close()

	data,err:=redis.Bytes(conn.Do("get","articleTypes"))
	if len(data)==0{
		o.QueryTable("ArticleType").All(&articleTypes)
		beego.Info("从MYSQL获取数据")
		var buffer bytes.Buffer
		encode:=gob.NewEncoder(&buffer)
		encode.Encode(&articleTypes)
		conn.Do("set","articleTypes",buffer.Bytes())
	}else{
		//var articleTypes2 []models.ArticleType
		dec:=gob.NewDecoder(bytes.NewReader(data))
		dec.Decode(&articleTypes)
	}
	//编码


//提取出来
	//data,_:=redis.Bytes(conn.Do("get","articleTypes"))

	//beego.Info(articleTypes2)





	this.Data["articleTypes"] = articleTypes


	//下拉框改变的时候，获取不同类型的文章
	//获取数据

	var count int64
	typeName := this.GetString("select")
	if typeName == ""{
		count,_ = qs.RelatedSel("ArticleType").Count()
		qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articles)
	}else {
		count,_ = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()
		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
	}
	//获取到类型数据,根据这个数据获取相应文章
	//默认多表查询是惰性查询


	pageCount := math.Ceil(float64(count) / float64(pageSize))
	this.Data["count"] = count
	this.Data["pageCount"] = pageCount

	//select * from Article where ArticleType.typeName = typeName
	this.Data["pageIndex"] = pageIndex

	this.Data["typeName"] = typeName
	//传递数据给前端
	this.Data["articles"] = articles
	//指定视图
	this.Layout = "layout.html"
	this.TplName = "index.html"
}

//处理下拉框选中的时候，加载不同的文章类型
func(this*ArticleController)HandleSelect(){
	//获取所有类型

	//获取所有文章

	//操作分页



	//指定视图
	this.TplName = "index.html"
}

//展示添加文章页面
func(this*ArticleController)ShowAdd(){
	//获取所有类型,传递给前端显示
	//获取orm对象
	o := orm.NewOrm()
	//获取所有
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)
	//传递给前段
	this.Data["articleTypes"] = articleTypes

	//指定视图
	this.Layout = "layout.html"
	this.TplName = "add.html"
}

//处理添加文章数据
func(this*ArticleController)HandleAdd(){
	//获取数据
	articleName :=this.GetString("articleName")
	content :=this.GetString("content")
	//file
	file,head,err :=this.GetFile("uploadname")

	//获取数据
	if articleName == "" || content == "" || err != nil{
		beego.Error("获取用户添加数据失败",err)
		this.TplName = "add.html"
		return
	}
	defer file.Close()
	//需要判断大小
	if head.Size > 5000000{
		beego.Error("图片太大，我不收")
		this.TplName = "add.html"
		return
	}
	//需要校验格式
	ext :=path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
		beego.Error("文件格式不正确")
		this.TplName = "add.html"
		return
	}

	//防止重名
	//beego.Info("time.now = ",time.Now().Format("2006-01-02 15:04:05"))
	fileName := time.Now().Format("20060102150405")
	//操作数据
	this.SaveToFile("uploadname","./static/img/"+fileName+ext)

	//把数据插入到数据库
	//获取orm对象
	o := orm.NewOrm()
	//获取插入对象
	var article models.Article
	//给插入对象赋值
	article.Title = articleName
	article.Content = content
	article.Img = "/static/img/"+fileName+ext

	//获取类型数据
	typeName := this.GetString("select")
	//获取类型对象
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Read(&articleType,"TypeName")
	//把类型对象插入文章中
	article.ArticleType = &articleType


	//插入
	o.Insert(&article)

	//返回数据
	this.Redirect("/article/index",302)
}

//展示文章详情页
func(this*ArticleController)ShowContent(){
	//获取数据
	id,err := this.GetInt("id")
	//校验数据
	if err != nil{
		beego.Error("获取数据错误",err)
		this.TplName = "index.html"
		return
	}
	//处理数据
	//查询数据库，获取文章信息
	//获取orm对象
	o := orm.NewOrm()
	//获取查询对象
	var article models.Article
	//指定查询条件
	article.Id2 = id
	//查询
	o.Read(&article)

	//多对多查询两种方式
	//o.LoadRelated(&article,"Users")
	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id2",article.Id2).Distinct().All(&users)
	this.Data["users"] = users



	//阅读次数加一
	//更新操作
	article.ReadCount += 1
	o.Update(&article)

	//返回数据
	this.Data["article"] = article
	this.Data["title"] = "商品详情"
	//添加浏览记录  多对多插入
	m2m:=o.QueryM2M(&article,"Users")
	//向这个表里面插入对象指针，获取对象
	var user models.User
	userName := this.GetSession("userName")
	//根据session中的用户名获取当前用户
	user.Name = userName.(string)
	o.Read(&user,"Name")
	//插入对象
	m2m.Add(user)



	//指定视图
	this.Layout = "layout.html"

	//this.LayoutSections = make(map[string]string)
	//this.LayoutSections["jsFile"] = "index.js"
	this.TplName = "content.html"
}

//展示编辑文章页面
func(this*ArticleController)ShowEditArticle(){
	//填充的文章原来的数据
	//获取数据
	id ,err :=this.GetInt("id")
	//数据校验
	if err != nil{
		beego.Error("获取数据错误",err)
		this.TplName = "index.html"
		return
	}
	//处理数据
	//查询
	o := orm.NewOrm()
	//获取查询对象
	var article models.Article
	//指定查询条件
	article.Id2 = id
	//查询
	o.Read(&article)


	//返回数据
	this.Data["article"] = article


	//指定视图
	this.TplName = "update.html"
}

//封装上传文件函数  重复使用相同代码   函数参数   返回值   写接口
func UploadFunc(this*ArticleController,fileName string)string{
	//file
	file,head,err :=this.GetFile(fileName)

	//获取数据
	if  err != nil{
		beego.Error("获取用户添加数据失败",err)
		this.TplName = "add.html"
		return ""
	}
	defer file.Close()
	//需要判断大小
	if head.Size > 5000000{
		beego.Error("图片太大，我不收")
		this.TplName = "add.html"
		return ""
	}
	//需要校验格式
	ext :=path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
		beego.Error("文件格式不正确")
		this.TplName = "add.html"
		return ""
	}

	//防止重名
	//beego.Info("time.now = ",time.Now().Format("2006-01-02 15:04:05"))
	filePath := time.Now().Format("20060102150405")
	//操作数据
	this.SaveToFile(fileName,"./static/img/"+filePath+ext)
	return "/static/img/"+filePath+ext
}

//处理编辑数据
func(this*ArticleController)HandleEditArticle(){
	//获取数据
	id,err :=this.GetInt("id")
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	filePath := UploadFunc(this,"uploadname")
	//校验数据
	if err != nil || articleName == "" || content == "" || filePath == ""{
		beego.Error("获取数据错误")
		this.TplName = "update.html"
		return
	}

	//处理数据
	//更新
	//获取orm对象
	o := orm.NewOrm()
	//获取更新对象
	var article models.Article
	//给更新条件赋值
	article.Id2 = id
	//先read一下，判断要更新的数据
	err = o.Read(&article)
	//更新
	if err != nil{
		beego.Error("更新数据不存在")
		this.TplName = "update.html"
		return
	}
	article.Title = articleName
	article.Content = content
	article.Img = filePath
	o.Update(&article)

	//返回数据
	this.Redirect("/article/index",302)
}

//处理删除文章
func(this*ArticleController)HandleDelete(){
	//获取数据
	id,err :=this.GetInt("id")
	//校验数据
	if err != nil{
		beego.Error("删除请求数据错误")
		this.TplName = "index.html"
		return
	}
	//处理数据
	//删除操作
	o := orm.NewOrm()
	//定义一个删除对象
	var article models.Article
	//指定删除条件
	article.Id2 = id
	//删除
	_,err =o.Delete(&article)
	if err != nil{
		beego.Error("删除失败")
		this.TplName = "index.html"
		return
	}
	//返回数据
	this.Redirect("/article/index",302)
}

//展示添加文章类型页面
func(this*ArticleController)ShowAddType(){
	//获取所有类型
	o := orm.NewOrm()

	qs := o.QueryTable("ArticleType")
	var articleTypes []models.ArticleType
	qs.All(&articleTypes)

	//传递数据给前段
	this.Data["articleTypes"] = articleTypes

	this.TplName = "addType.html"
}

//处理添加类型数据
func(this*ArticleController)HandleAddType(){
	//获取数据
	typeName := this.GetString("typeName")
	//校验数据
	if typeName == ""{
		beego.Error("类型数据不能为空")
		this.TplName = "addType.html"
		return
	}
	//处理数据
	o := orm.NewOrm()
	//获取插入对象
	var articleType models.ArticleType
	//给插入对象赋值
	articleType.TypeName = typeName
	//插入
	o.Insert(&articleType)
	//返回数据
	this.Redirect("/article/addType",302)
}

//删除文章类型
func(this*ArticleController)DeleteType(){
	//获取数据
	id ,err :=this.GetInt("id")
	//校验数据
	if err != nil {
		beego.Error("获取类型数据失败")
		this.TplName = "addType.html"
		return
	}
	//删除数据
	o := orm.NewOrm()
	//获取删除对象
	var articleType models.ArticleType
	//给删除条件赋值
	articleType.Id = id
	//删除
	o.Delete(&articleType)

	this.Redirect("/article/addType",302)

}