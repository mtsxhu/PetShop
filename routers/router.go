package routers

import (
	"webItem/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

var filterFunc= func(ctx* context.Context) {
	userName:=ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302,"/login")
		return
	}
}
func init() {
	beego.InsertFilter("/user/*",beego.BeforeExec,filterFunc)
    //用户注册
	beego.Router("/register", &controllers.UserController{},
	"get:ShowRegister;post:HandleRegister")
    //用户激活
    beego.Router("/active",&controllers.UserController{},"get:ActiveUser")
    //用户登录
    beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
	//退出登陆
	beego.Router("/user/logout",&controllers.UserController{},"get:Logout")
	//商品首页
	beego.Router("/",&controllers.GoodsController{},"get:ShowIndex")
	//用户中心页
	beego.Router("/user/userCenterInfo",&controllers.UserController{},"get:ShowCenterInfo")
	//用户中心订单页
	beego.Router("/user/userCenterOrder",&controllers.UserController{},"get:ShowCenterOrder")
	//用户收货地址页
	beego.Router("/user/userCenterSite",&controllers.UserController{},"get:ShowCenterSite;post:HandleCenterSite")
	// 商品详情展示
	beego.Router("/goodsDetail",&controllers.GoodsController{},"get:ShowGoodsDetail")
	beego.Router("/goodsList",&controllers.GoodsController{},"get:ShowGoodsList")
	beego.Router("/goodsSearch",&controllers.GoodsController{},"post:HandleGoodsSearch")
	beego.Router("/user/cart",&controllers.CartController{},"get:ShowCart")
	beego.Router("/user/UpdateCart",&controllers.CartController{},"post:UpdateCart")
}

