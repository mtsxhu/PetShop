package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"webItem/models"
)

type CartController struct {
	beego.Controller
}

func (this *CartController) HandleAddCart() {
	skuid,err:=this.GetInt("skuid")
	resp:=make(map[string]interface{})
	defer this.ServeJSON()
	if err != nil {
		resp["code"]=1
		resp["msg"]="传递的数据不正确"
		this.Data["json"]=resp
		fmt.Println("HandleAddCart get skuid err:",err)
		return
	}
	count,err:=this.GetInt("count")
	if err != nil {
		resp["code"]=1
		resp["msg"]="传递的数据不正确"
		this.Data["json"]=resp
		fmt.Println("HandleAddCart get skuid err:",err)
		return
	}
	userName:=GetUser(&this.Controller)
	if userName=="" {
		resp["code"]=2
		resp["msg"]="未登录"
		this.Data["json"]=resp
		fmt.Println("未登录",err)
		return
	}
	o:=orm.NewOrm()
	var user models.User
	user.Name=userName
	o.Read(&user,"Name")
	conn,err:=redis.Dial("tcp", "192.168.168.106:6379")
	defer conn.Close()
	if err != nil {
		fmt.Println("redis conn err: ",err)
		return
	}
	preCount,err:=redis.Int(conn.Do("hget","cart_"+strconv.Itoa(user.Id),skuid))
	conn.Do("hset","cart_"+strconv.Itoa(user.Id),skuid,count+preCount)
	ret,err:=conn.Do("hlen","cart_"+strconv.Itoa(user.Id))
	cartCount,_:=redis.Int(ret,err)
	resp["code"]=5
	resp["msg"]="ok"
	resp["cartCount"]=cartCount
	this.Data["json"]=resp
}
func GetCartCount(c *beego.Controller)int{
	// 从redis中读取
	userName:=GetUser(c)
	if userName=="" {
		return 0
	}
	o:=orm.NewOrm()
	var user models.User
	user.Name=userName
	o.Read(&user,"Name")
	conn,err:=redis.Dial("tcp", "192.168.168.106:6379")
	defer conn.Close()
	if err != nil {
		return 0
	}
	ret,err:=conn.Do("hlen","cart_"+strconv.Itoa(user.Id))
	cartCount,_:=redis.Int(ret,err)
	return cartCount
}

func (this *CartController) ShowCart()  {
	conn,err:=redis.Dial("tcp", "192.168.168.106:6379")
	defer conn.Close()
	if err != nil {
		fmt.Println("redis conn err: ",err)
		return
	}
	userName:=GetUser(&this.Controller)
	var user models.User
	user.Name=userName
	o:=orm.NewOrm()
	o.Read(&user,"Name")
	fmt.Printf("%+v\n",user)
	goodsMap,err:=redis.IntMap(conn.Do("hgetall","cart_"+strconv.Itoa(user.Id)))
	goods:=make([]map[string]interface{},len(goodsMap))
	i:=0
	totalPrice:=0
	totalCount:=0
	for index, value := range goodsMap {
		skuid,_:=strconv.Atoi(index)
		var goodsSku models.GoodsSKU
		goodsSku.Id=skuid
		o.Read(&goodsSku)
		tmp:=make(map[string]interface{})
		tmp["goodsSku"]=goodsSku
		tmp["count"]=value
		tmp["addPrice"]=goodsSku.Price*value
		goods[i]=tmp
		i+=1
		totalPrice+=goodsSku.Price*value
		totalCount+=value
	}
	this.Data["totalPrice"]=totalPrice
	this.Data["totalCount"]=totalCount
	this.Data["goods"]=goods
	fmt.Printf("%+v\n",goods)
	this.TplName="cart.html"
}

func (this *CartController) UpdateCart()  {
	skuid,err1:=this.GetInt("skuid")
	count,err2:=this.GetInt("count")
	fmt.Println(skuid,count)
	resp:=make(map[string]interface{})
	defer this.ServeJSON()
	if err1 != nil || err2 != nil {
		resp["code"]=1
		resp["msg"]="请求数据不正确"
		this.Data["json"]=resp
		return
	}

	userName:=GetUser(&this.Controller)
	var user models.User
	user.Name=userName
	o:=orm.NewOrm()
	o.Read(&user,"Name")
	conn,err:=redis.Dial("tcp", "192.168.168.106:6379")
	defer conn.Close()
	if err != nil {
		fmt.Println("redis conn err: ",err)
		return
	}
	conn.Do("hset","cart_"+strconv.Itoa(user.Id),skuid,count)
	resp["code"]=5
	resp["msg"]="ok"
	this.Data["json"]=resp
}












