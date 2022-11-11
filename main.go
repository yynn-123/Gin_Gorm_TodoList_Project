package main

// 技术选型:导入gorm、gin框架
// 运行之前记得go mod tidy进行查缺补漏
import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)

// Todo 待办列表字段信息
type Todo struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

func main() {
	// 连接数据库
	db, err := gorm.Open("mysql", "root:123456@(127.0.0.1)/todolist?charset=utf8mb4&parseTime=True&loc=Local")
	// 连接数据库报错处理
	if err != nil {
		panic(err)
	}
	// defer实现程序运行后关闭数据库连接
	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	// 数据迁移
	db.AutoMigrate(&Todo{})
	// 定义默认路由引擎
	r := gin.Default()
	// 静态文件指向
	r.Static("/static", "static")
	// 渲染模板
	r.LoadHTMLGlob("templates/**")
	// 代办列表首页展示
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// 待办列表的功能:增、删、改、查
	// 根据前端请求,判断出所有功能具有部分相同的路径信息,所以此处采用路由组实现
	v1Group := r.Group("/v1")
	{
		// 新增
		v1Group.POST("/todo", func(c *gin.Context) {
			// 定义todo变量,类型为Todo
			var todo Todo
			// 前端发送json请求,后端采用BindJSON进行数据接收
			// 在结构体中已经定义了Todo类型,并且通过反序列化来接收前端传过来的信息
			// errorBind进行错误处理
			errorBind := c.BindJSON(&todo)
			if errorBind != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"message": "绑定数据失败",
					"data":    err.Error(),
				})
				return
			}
			// 数据库进行数据写入,以及进行写入的错误处理
			errorCreate := db.Create(&todo).Error
			if errorCreate != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"message": "数据写入失败",
					"data":    err.Error(),
				})
				return
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"message": "SUC0000",
					"data":    todo,
				})
			}
		})
		// 修改
		v1Group.PUT("/todo/:id", func(c *gin.Context) {
			// 定义updateObject变量,类型为Todo
			var updateObject Todo
			// 取到前端传过来的id,并判断id是否存在
			id, boolGet := c.Params.Get("id")
			if !boolGet {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"message": "无效的修改",
					"data":    err.Error(),
				})
				return
			} else {
				// 根据id查询数据库中的对应记录,并将查到的副本传给&updateObject
				db.Where("id = ?", id).First(&updateObject)
				// errorBind进行错误处理
				errorBind := c.BindJSON(&updateObject)
				if errorBind != nil {
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"message": "绑定数据失败",
						"data":    err.Error(),
					})
					return
				} else {
					// 进行数据库对应记录修改,然后返回一个errorSave
					errorSave := db.Save(&updateObject).Error
					if errorSave != nil {
						c.JSON(http.StatusOK, gin.H{
							"code":    200,
							"message": "数据修改失败",
							"data":    err.Error(),
						})
						return
					} else {
						c.JSON(http.StatusOK, gin.H{
							"code":    200,
							"message": "SUC0000",
							"data":    updateObject,
						})
					}
				}
			}
		})
		// 删除:相关信息同修改
		v1Group.DELETE("/todo/:id", func(c *gin.Context) {
			id, boolDelete := c.Params.Get("id")
			if !boolDelete {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"message": "无效的删除",
					"data":    err.Error(),
				})
				return
			} else {
				errorDelete := db.Where("id = ?", id).Delete(Todo{}).Error
				if errorDelete != nil {
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"message": "删除失败",
						"data":    err.Error(),
					})
					return
				}
			}
		})
		// 查询
		v1Group.GET("/todo", func(c *gin.Context) {
			// 定义一个todoList对象，类型为Todo切片
			var todolist []Todo
			// 查询数据库中全部记录信息
			errorFind := db.Find(&todolist).Error
			if errorFind != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"message": "查询数据库失败",
					"data":    err.Error(),
				})
				return
			} else {
				c.JSON(http.StatusOK, todolist)
			}
		})
	}
	errRun := r.Run(":8091")
	if errRun != nil {
		return
	}
}
