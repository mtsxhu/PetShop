package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"math"
	"strconv"
	"webItem/models"
)

type GoodsController struct {
	beego.Controller
}

func (this *GoodsController)ShowIndex()  {
	GetUser(&this.Controller)
	o:=orm.NewOrm()

	// 获取类型数据
	var goodsTypes []models.GoodsType
	o.QueryTable("GoodsType").All(&goodsTypes)
	this.Data["goodsTypes"]=goodsTypes

	// 获取轮播图数据
	var indexGoodsBanner []models.IndexGoodsBanner
	o.QueryTable("IndexGoodsBanner").OrderBy("Index").All(&indexGoodsBanner)
	this.Data["indexGoodsBanner"]=indexGoodsBanner

	// 获取促销商品数据
	var indexPromotionBanner []models.IndexPromotionBanner
	o.QueryTable("IndexPromotionBanner").OrderBy("Index").All(&indexPromotionBanner)
	this.Data["indexPromotionBanner"]=indexPromotionBanner

	// 首页展示商品数据
	goods:=make([]map[string]interface{},len(goodsTypes))
	// 向切片interface中插入类型数据
	for index, value := range goodsTypes {
		// 获取对应类型的首页展示商品
		tmp :=make(map[string]interface{})
		tmp["type"]=value
		goods[index]=tmp
	}
	for _, value := range goods {
		var testGoods []models.IndexTypeGoodsBanner
		var imgGoods []models.IndexTypeGoodsBanner
		// 获取文字
		o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType","GoodsSKU").OrderBy("Index").Filter("GoodsType",value["type"]).Filter("Display_Type",0).All(&testGoods)
		// 获取图片商品数据
		o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType","GoodsSKU").OrderBy("Index").Filter("GoodsType",value["type"]).Filter("Display_Type",1).All(&imgGoods)
		value["textGoods"]=testGoods
		value["imgGoods"]=imgGoods
		fmt.Printf("testGoods:  	%+v\n",testGoods)
		fmt.Printf("imgGoods:	%+v\n",imgGoods)
	}
	this.Data["goods"]=goods
	this.TplName="index.html"
}

// 展示商品详情
func (this *GoodsController)ShowGoodsDetail(){
	// 获取数据
	id,err:=this.GetInt("id")

	// 校验数据
	if err != nil {
		beego.Error("浏览器请求错误")
		this.Redirect("/",302)
	}

	// 处理数据
	o:=orm.NewOrm()
	var goodsSku models.GoodsSKU
	goodsSku.Id=id
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType","Goods").Filter("Id",id).One(&goodsSku)
	// 获取同类数据时间靠前的两天商品
	var goodsNew []models.GoodsSKU
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType",goodsSku.GoodsType).OrderBy("Time").Limit(2,0).All(&goodsNew)
	this.Data["goodsNew"]=goodsNew

	// 添加历史浏览记录,存redis
	userName:=GetUser(&this.Controller)
	if userName != "" {
		var user models.User
		user.Name=userName
		o.Read(&user,"Name")
		conn,err:=redis.Dial("tcp", "192.168.168.106:6379")
		defer conn.Close()
		if err != nil {
			fmt.Println("redis conn err: ",err)
		}
		str,err:=conn.Do("lrem","history_"+strconv.Itoa(user.Id),0,id)
		if err != nil {
			fmt.Println("redis do err: ",err)
		}else {
			fmt.Println("redis do : ",str)
		}
		str,err=conn.Do("lpush", "history_"+strconv.Itoa(user.Id), id)
		if err != nil {
			fmt.Println("redis do err: ",err)
		}else {
			fmt.Println("redis do : ",str)
		}
	}
	// 返回视图
	this.Data["goodsSku"]=goodsSku
	ShowLayout(&this.Controller)
	this.Data["cartCount"]=GetCartCount(&this.Controller)

	this.TplName="detail.html"
}

func ShowLayout(this *beego.Controller){
	// 查询类型
	o:=orm.NewOrm()
	var types []models.GoodsType
	o.QueryTable("GoodsType").All(&types)
	this.Data["types"]=types
	// 获取用户信息
	GetUser(this)
	// 指定layout
	this.Layout="goodsLayout.html"
}

// 展示商品列表页
func (this *GoodsController)ShowGoodsList(){
	id,err:=this.GetInt("typeId")
	if err != nil {
		fmt.Println("get typeId err:",err)
		this.Redirect("/",302)
		return
	}
	ShowLayout(&this.Controller)
	// 获取列表
	o:=orm.NewOrm()
	var goodsNew []models.GoodsSKU
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",id).OrderBy("Time").Limit(2,0).All(&goodsNew)

	var goods []models.GoodsSKU
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",id).All(&goods)
	this.Data["goods"]=goods
	this.Data["goodsNew"]=goodsNew

	// 分页
	count,_:=o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",id).Count()
	pageSize:=3
	pageCount:=math.Ceil(float64(count)/float64(pageSize))
	pageIndex,err:=this.GetInt("pageIndex")
	if err != nil {
		pageIndex=1
	}
	pages:=PageTool(int(pageCount),pageIndex)
	this.Data["pages"]=pages
	this.Data["typeId"]=id
	this.Data["pageIndex"]=pageIndex
	start:=(pageIndex-1)*pageSize
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",id).Limit(pageSize,start).All(&goods)
	this.Data["goods"]=goods

	// 获取上一页页码
	prePage:=pageIndex-1
	if prePage<=1 {
		prePage=1
	}
	this.Data["prePage"]=prePage
	// 获取下一页页码
	nextPage:=pageIndex+1
	if nextPage>int(pageCount) {
		prePage=int(pageCount)
	}
	this.Data["nextPage"]=nextPage

	// 排序
	sort:=this.GetString("sort")
	if sort == "sale" {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",id).OrderBy("Sales").Limit(pageSize,start).All(&goods)
	}else if sort == "price" {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",id).OrderBy("Price").Limit(pageSize,start).All(&goods)
	}else {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",id).Limit(pageSize,start).All(&goods)
	}
	this.Data["goods"]=goods
	this.Data["sort"]=sort
	this.TplName="list.html"
}

func PageTool(pageCount int,pageIndex int)[]int{
	var pages []int
	if pageCount <= 5 {
		pages= make([]int,pageCount)
		for i, _ := range pages {
			pages[i]=i+1
		}
	}else if pageIndex <=3{
		pages=[]int{1,2,3,4,5}
	}else if pageIndex >pageCount{
		pages=[]int{pageCount-4,pageCount-3,pageCount-2,pageCount-1,pageCount}
	}else {
		pages=[]int{pageIndex-2,pageIndex-1,pageIndex,pageIndex+1,pageIndex+2}

	}
	return pages
}

func (this *GoodsController) HandleGoodsSearch() {
	goodsName:=this.GetString("goodsName")
	fmt.Println("goodsName: ",goodsName)
	o:=orm.NewOrm()
	var goods []models.GoodsSKU
	if goodsName == "" {
		o.QueryTable("GoodsSKU").All(&goods)
	}else {
		// 处理数据
		o.QueryTable("GoodsSKU").Filter("Name__icontains",goodsName).All(&goods)
	}
	fmt.Println("========================")
	fmt.Printf("goods:%+v",goods)
	this.Data["goods"]=goods
	ShowLayout(&this.Controller)
	this.TplName="search.html"
}