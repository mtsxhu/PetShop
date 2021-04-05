package controllers

import (
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
	"github.com/garyburd/redigo/redis"
	"regexp"
	"strconv"
	"webItem/models"
)

type UserController struct {
	beego.Controller
}

func GetUser(this *beego.Controller)  string{
	userName:=this.GetSession("userName")
	if userName == nil {
		this.Data["userName"]=""
	}else {
		this.Data["userName"]=userName.(string)
		return userName.(string)
	}
	return ""
}

func (this *UserController)ShowRegister(){
	this.TplName="register.html"
}

//用户注册
 func (this *UserController)HandleRegister(){
	user_name:=this.GetString("user_name")
	pwd:=this.GetString("pwd")
	cpwd:=this.GetString("cpwd")
	email:=this.GetString("email")
	allow:=this.GetString("allow")
	if user_name==""||pwd == "" || cpwd == "" || email == "" || allow == "" {
		this.Data["errMSG"]="请完善信息"
		this.TplName="register.html"
		return
	}
	//邮箱正则表达式
	reg,_:=regexp.Compile("^[A-Za-z0-9\u4e00-\u9fa5]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")
	res:=reg.FindString(email)
	if res==""{
		this.Data["errMSG"]="邮箱格式不正确"
		this.TplName="register.html"
		return
	}

	//插入数据
	o:=orm.NewOrm()
	user:=models.User{}
	user.Name=user_name
	user.Email=email
	user.PassWord=pwd
	_,err:=o.Insert(&user)
	if err!=nil{
		this.Data["errMSG"]="用户名已存在"
		this.TplName="register.html"
		return
	}
	//邮箱注册业务
	emailConfig:=`{"username":"734774492@qq.com","password":"rbbtirehfumzbdac","host":"smtp.qq.com","port":587}`
	emailConn:=utils.NewEMail(emailConfig)
	//发送者
	emailConn.From="734774492@qq.com"
	emailConn.To=[]string{email}
	//邮件标题
	emailConn.Subject="天天生鲜注册"
	//激活请求地址
	emailConn.Text="192.168.168.106:8080/active?id="+strconv.Itoa(user.Id)
	err=emailConn.Send()
	if err!=nil {
		beego.Info("邮件发送失败",err)
	}
	this.Ctx.WriteString("注册成功，请去邮箱激活页面")
}

//激活处理
func (this *UserController)ActiveUser(){
	//获取数据
	id,err:=this.GetInt("id")

	//校验数据
	if err!=nil {
		this.Data["errMSG"]="要激活的用户不存在"
		this.TplName="register.html"
		return
	}

	//处理数据
	o:=orm.NewOrm()
	var user models.User
	user.Id=id
	err=o.Read(&user)
	if err!=nil {
		this.Data["errMSG"]="要激活的用户不存在"
		this.TplName="register.html"
		return
	}
	user.Active=true
	o.Update(&user)

	//返回视图
	this.Redirect("/login",302)
}

func (this *UserController)ShowLogin(){
	userName:=this.Ctx.GetCookie("userName")
	temp,_:=base64.StdEncoding.DecodeString(userName)
	if string(temp) == "" {
		this.Data["userName"]=""
		this.Data["pwd"]=""
	}else {
		this.Data["userName"]=string(temp)
		this.Data["checked"]="checked"
	}

	this.TplName="login.html"
}

//登陆处理
func (this *UserController)HandleLogin(){
	//获得数据
	userName:=this.GetString("userName")
	pwd:=this.GetString("pwd")
	//校验数据
	if userName==""|| pwd==""{
		this.Data["errMSG"]="登陆信息不完整"
		this.TplName="login.html"
		return
	}
	o:=orm.NewOrm()
	var user models.User
	user.Name=userName
	//根据name字段查询
	err:=o.Read(&user,"Name")
	if err != nil {
		this.Data["errMSG"]="用户名或密码错误"
		this.TplName="login.html"
		return
	}
	if user.PassWord!=pwd {
		this.Data["errMSG"]="用户名或密码错误"
		this.TplName="login.html"
		return
	}
	if user.Active!=true {
		this.Data["errMSG"]="用户未激活"
		this.TplName="login.html"
		return
	}
	beego.Info(pwd)
	rembUser:=this.GetString("rembUser")
	if rembUser == "on" {
		temp:=base64.StdEncoding.EncodeToString([]byte(userName))
		this.Ctx.SetCookie("userName",temp,24*3600*30)
	}else {
		this.Ctx.SetCookie("userName",userName,-1)
	}
	this.SetSession("userName",userName)
	this.Redirect("/",302)
}

//退出登陆
func (this *UserController) Logout()  {
	this.DelSession("username")
	this.Redirect("/login",302)
}

//用户中心页面
func (this *UserController)ShowCenterInfo(){
	userName:=GetUser(&this.Controller)
	this.Data["userName"]=userName

	//查询地址
	o:=orm.NewOrm()
	var addr models.Address
	o.QueryTable("Address").RelatedSel("User").Filter("User__Name",userName).Filter("Isdefualt",true).One(&addr)
	if addr.Id==0{
		this.Data["addr"]=""
	}else {
		this.Data["addr"]=addr
	}

	// 获取历史浏览记录
	conn,err:=redis.Dial("tcp", "192.168.168.106:6379")
	if err != nil {
		fmt.Println("redis conn err: ",err)
	}
	var user models.User
	user.Name=userName
	o.Read(&user,"Name")
	rep,err:=conn.Do("lrange","history_"+strconv.Itoa(user.Id),0,4)
	defer conn.Close()
	var goodsSkus []models.GoodsSKU
	goodsIds,err:=redis.Ints(rep,err)
	for _, value := range goodsIds {
		var goods models.GoodsSKU
		goods.Id=value
		o.Read(&goods)
		goodsSkus=append(goodsSkus,goods)
	}
	this.Data["goodsSkus"]=goodsSkus
	fmt.Println("goodsSkus:",goodsSkus)
	this.Layout="userCenterLayout.html"
	this.TplName="user_center_info.html"
}

//用户中心订单页
func (this *UserController)ShowCenterOrder(){
	GetUser(&this.Controller)
	this.Layout="userCenterLayout.html"
	this.TplName="user_center_order.html"
}

//用户中心地址页
func (this *UserController)ShowCenterSite(){
	userName:=GetUser(&this.Controller)

	o:=orm.NewOrm()
	var addr models.Address
	o.QueryTable("Address").RelatedSel("User").Filter("User__Name",userName).Filter("Isdefualt",true).One(&addr)
	this.Data["addr"]=addr
	//设置视图布局
	this.Layout="userCenterLayout.html"
	//设置页面
	this.TplName="user_center_site.html"
}

//用户中心地址页数据处理
func (this *UserController)HandleCenterSite(){
	//获取数据
	reciver:=this.GetString("reciver")
	address:=this.GetString("address")
	zipCode:=this.GetString("zipCode")
	phone:=this.GetString("phone")
	//校验数据
	if reciver == ""||address == ""||zipCode == ""||phone == "" {
		this.Redirect("/userCenterSite",302)
		return
	}
	//处理数据
	//插入数据
	o:=orm.NewOrm()
	var addrUser models.Address
	addrUser.Isdefualt=true
	err:=o.Read(&addrUser,"Isdefualt")
	if err == nil {
		addrUser.Isdefualt=false
		o.Update(&addrUser)
	}

	var user models.User
	user.Name=GetUser(&this.Controller)
	o.Read(&user,"Name")

	var addrUserNew models.Address
	addrUserNew.Isdefualt=true
	addrUserNew.Addr=address
	addrUserNew.Zip_code=zipCode
	addrUserNew.Phone=phone
	addrUserNew.Receiver=reciver
	addrUserNew.User=&user
	o.Insert(&addrUserNew)
	//返回数据
	this.Redirect("/user/userCenterInfo",302)
}