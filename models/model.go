package models

import (
	"github.com/astaxie/beego/orm"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

//用户表
type User struct {
	Id 			int
	Name 		string`orm:"size(20)"`				//用户名
	PassWord	string`orm:"size(20)"`				//登陆密码
	Email		string`orm:"size(50)"`				//邮箱
	Active		bool`orm:"defualt(false)"`			//是否激活
	Power		int`orm:"defualt(0)"`				//权限
	Address		[]*Address`orm:"reverse(many)"`
	OrderInfo	[]*OrderInfo`orm:"reverse(many)"`
}

//地址表
type Address struct {
	Id 			int
	Receiver 	string`orm:"size(20)"`
	Addr 		string`orm:"size(50)"`				//用户名
	Zip_code	string`orm:"size(20)"`				//登陆密码
	Phone		string`orm:"size(20)"`
	Isdefualt	bool`orm:"defualt(false)"`
	User 		*User`orm:"rel(fk)"`
	OrderInfo	[]*OrderInfo`orm:"reverse(many)"`
}

//商品表
type Goods struct {
	Id 			int
	Name 		string`orm:"size(20)"`				//用户名
	Detail		string`orm:"size(200)"`
	GoodsSKU	[]*GoodsSKU`orm:"reverse(many)"`
}

//商品类型表
type GoodsType struct {
	Id 			int
	Name 		string
	Logo 		string
	Image		string
	GoodsSKU	[]*GoodsSKU`orm:"reverse(many)"`
	IndexTypeGoodsBanner	[]*IndexTypeGoodsBanner`orm:"reverse(many)"`
}

//商品库存
type GoodsSKU struct {
	Id 				int
	Goods			*Goods`orm:"rel(fk)"`
	GoodsType		*GoodsType`orm:"rel(fk)"`
	Name 			string
	Desc			string
	Price			int
	Unite			string
	Image			string
	Stock			int`orm:"defualt(1)"`
	Sales			int`orm:"defualt(0)"`
	Status			int`orm:"defualt(1)"`
	Time  			time.Time`orm:"auto_now_add"`
	GoodsImage		[]*GoodsImage`orm:"reverse(many)"`
}

//商品图片
type GoodsImage struct {
	Id 				int
	Image			string
	GoodsSKU		*GoodsSKU`orm:"rel(fk)"`
}

//首页轮播商品表
type IndexGoodsBanner struct {
	Id 				int
	GoodsSKU		*GoodsSKU`orm:"rel(fk)"`
	Image 			string
	Index 			int`orm:"defualt(0)"`
}
//首页分类商品列表
type IndexTypeGoodsBanner struct {
	Id				int
	GoodsType		*GoodsType`orm:"rel(fk)"`
	GoodsSKU		*GoodsSKU`orm:"rel(fk)"`
	Display_Type 	int`orm:"defualt(1)"` // 0文字，1图片
	Index 			int`orm:"defualt(0)"`
}

//首页促销商品列表
type IndexPromotionBanner struct {
	Id				int
	Name 			string`orm:"size(20)"`
	Url 			string`orm:"size(50)"`
	Image 			string
	Index 			int`orm:"defualt(0)"`
}

//订单表
type OrderInfo struct {
	Id 				int
	OrderId 		string`orm:"unique"`
	User 			*User`orm:"rel(fk)"`
	Address			*Address`orm:"rel(fk)"`
	Pay_Method		int
	Total_Price		int`orm:"defualt(1)"`
	Transit_Price	int
	Order_status	int`orm:"defualt(1)"`
	Trade_No		string`orm:"defualt('')"`
	Time  			time.Time`orm:"auto_now_add"`
	OrderGoods		[]*OrderGoods`orm:"reverse(many)"`
}

//订单商品表
type OrderGoods struct {
	Id 				int
	OrderInfo 		*OrderInfo`orm:"rel(fk)"`
	GoodsSKU		*GoodsSKU`orm:"rel(fk)"`
	Count 			int`orm:"defualt(1)"`
	Price 			int
	Comment 		string`orm:"defualt('')"`
}

func init() {
	orm.RegisterDataBase("default","mysql",
		"root:123456@tcp(127.0.0.1:3306)/dailyfresh?charset=utf8")
	orm.RegisterModel(new(User),new(OrderGoods),new(OrderInfo),new(IndexPromotionBanner),
		new(IndexTypeGoodsBanner), new(IndexGoodsBanner),
		new(GoodsImage),new(GoodsSKU),new(GoodsType),new(Goods),new(Address),)
	orm.RunSyncdb("default",false,true)
}
















