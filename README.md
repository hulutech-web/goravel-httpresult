# http_result

#### 一、安装
```go
go get -u github.com/hulutech-web/http_result

```
##### 1.1 注册服务提供者:config/app.go
```go
import	"github.com/hulutech-web/http_result"

func init() {
"providers": []foundation.ServiceProvider{
	....
	&http_result.ServiceProvider{},
}
}

```
##### 1.2发布资源  
```go
go run . artisan vendor:publish --package=github.com/hulutech-web/http_result

```

#### 二、使用

##### 2.1 使用说明:自定义默认返回
发布资源后，config/http_result.go中的配置文件中有默认的返回状态码和message，可自行修改
```go
config.Add("http_result", map[string]any{
		"Code":    500, //自定义修改默认状态码
		"Message": "获取成功",//自定义修改默认消息
})
```
#### 使用：ctx为goravel控制器中的默认(ctx http.Context) "github.com/goravel/framework/contracts/http"
##### 方式一：
一、成功返回：
```
import	httpfacades "github.com/hulutech-web/http_result"

return httpfacades.NewResult(ctx).Success("", user)
```
二、失败返回：
```go
import	httpfacades "github.com/hulutech-web/http_result"
httpfacades.NewResult(ctx).Error(500, "用户不存在", "no users find")
```
三、表单验证错误：
```go
httpfacades.NewResult(ctx).ValidError("验证失败", errors.All())
```
四、分页查询+加上过滤条件：其中conditionMap:map[string]interface{}{}过滤条件,[]string{"excepts"}...为排除的key值
```go

func (r *MaterialController) Index(ctx http.Context) http.Response {
    var user models.User
    facades.Auth(ctx).User(&user)
    partner_id := cast.ToUint(ctx.Value("partner"))
    materials := []models.Material{}
    
    queries := ctx.Request().Queries()
    conditionMap := map[string]interface{}{
    "partner_id": partner_id,
    }
    result, _ := httpfacades.NewResult(ctx).SearchByParams(queries, conditionMap,[]string{"excepts"}...).ResultPagination(&materials)
    return result
}
```

说明： 配合前端表格渲染，其中固定参数为pageSize,currentPage,sort,order，其他参数将默认采用like模糊查询
##### SearchByParams添加排除模糊搜索的字段，ResultPagination添加关联查询并返回（需在模型中定义关联）
##### example：?name=xxx&pageSize=10&currentPage=1&sort=age&order=desc
##### 解释一：根据表中的列name模糊查询，每页10条，当前第一页，按照age列降序排列
##### 解释二：查询参数：order仅支持asc,desc默认采用asc,sort表示需要查询的列，order参数和sort参数需要同时出现
##### 解释三：前端，可以配合vue组件，https://vxetable.cn，使用效率更高
```go
// Index 用户分页查询，支持搜索，路由参数?name=xxx&pageSize=1&currentPage=1&sort=xxx&order=xxx,等其他任意的查询参数
// @Summary      用户分页查询
// @Description  用户分页查询
// @Tags         用户分页查询
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @Param Authorization header string false "Bearer 用户令牌"
// @Param  name  query  string  false  "name"
// @Param  pageSize  query  string  false  "pageSize"
// @Param  currentPage  query  string  false  "currentPage"
// @Param  sort  query  string  false  "sort"
// @Param  order  query  string  false  "order"
// @Success 200 {string} json {}
// @Router       /api/user [get]
```

##### 解释四：返回结果
```json
{
    "data": [
        {
            "id": 6,
            "created_at": "2024-05-19 11:22:22",
            "updated_at": "2024-05-19 11:22:22",
            "name": "Karolann Waelchi",
            "mobile": "Annalise Koss",
            "password": "eyJpdiI6Im5JbVNXQ2pWV0FOVkxFTG0iLCJ2YWx1ZSI6IkJKYTc5bWt0WWRrUFRPYVJlMW5NcWN0SXFWK29iYVBqIn0=",
            "area": "",
            "contact": "Juana Russel",
            "contact_mobile": "9210102772",
            "address": "95469 New Bypassshire",
            "id_card": "4034197872575788",
            "control_arr": null,
            "pid": 0,
            "parent": null,
            "children": null,
            "deleted_at": null
        }
    ],
    "total": 1,
    "links": {
        "first": "http://localhost:3000/api/user/indexs?pageSize=2&currentPage=1",
        "last": "http://localhost:3000/api/user/indexs?pageSize=2&currentPage=0",
        "prev": "http://localhost:3000/api/user/indexs?pageSize=2&currentPage=0",
        "next": "http://localhost:3000/api/user/indexs?pageSize=2&currentPage=2"
    },
    "meta": {
        "total_page": 1,
        "current_page": 1,
        "per_page": 2,
        "total": 1
    }
}
```
