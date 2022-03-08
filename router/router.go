package router

import (
	"GanLianInfo/controllers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func genDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		if controllers.FixDB() {
			c.Next()
		} else {
			r := controllers.Errors.DatabaseError
			c.JSON(200, r)
		}
	}
}

func Register() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	//router := gin.Default()
	// 新建一个没有任何默认中间件的路由
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(genDB())
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("X-Token")
	router.Use(cors.New(config))
	router.POST("/initdb", controllers.InitDb)
	router.POST("/fake", controllers.Fake)
	router.POST("/add", controllers.Add)
	router.POST("/update", controllers.Update)
	router.POST("/delete", controllers.Delete)
	router.POST("/dashboard", controllers.DashboardData)
	user := router.Group("user")
	{
		user.POST("/login", controllers.Login)
		//user.POST("/logout", controllers.Logout)
		user.POST("/info", controllers.UserInfo)
		user.POST("/organ", controllers.GetPersonOrganId)
		user.POST("/organs", controllers.GetPersonOrgans)
		//user.POST("/userinfo", controllers.GetUserInfo)
		//user.POST("/police", controllers.PoliceInfo)
		//user.POST("/photo", controllers.PolicePhoto)

	}
	module := router.Group("module")
	{
		module.POST("/list", controllers.GetModuleList)
		module.POST("/add", controllers.ModuleAdd)
		module.POST("/update", controllers.ModuleUpdate)
		module.POST("/delete", controllers.ModuleDelete)
		module.POST("/order", controllers.ModuleOrder)
		module.POST("/role", controllers.ModuleRole)
	}
	personnel := router.Group("personnel")
	{
		personnel.POST("/list", controllers.PersonnelList)
		personnel.POST("/detail", controllers.PersonnelDetail)
		personnel.POST("/update", controllers.PersonnelUpdate)
		personnel.POST("/searchName", controllers.SearchPersonnelName)
		personnel.POST("/nameList", controllers.GetPersonnelName)
		personnel.POST("/name_list", controllers.PersonnelNameList)
		personnel.POST("/dict", controllers.EduDictList)
		personnel.POST("/resume", controllers.PersonnelResume)
	}
	//organ := router.Group("organ")
	//{
	//	//organ.GET("list", controllers.GetOrganList)
	//	organ.POST("/add", controllers.OrganAdd)
	//	organ.POST("/update", controllers.OrganUpdate)
	//	organ.POST("/delete", controllers.OrganDelete)
	//}
	department := router.Group("department")
	{
		department.POST("list", controllers.GetDepartmentList)
		department.POST("organ", controllers.GetOrganList)
	}
	training := router.Group("training")
	{
		training.POST("list", controllers.TrainingList)
		training.POST("person_list", controllers.TrainPersonList)
		training.POST("add", controllers.TrainPersonAdd)
		training.POST("delete", controllers.TrainPersonDelete)
		training.POST("detail", controllers.TrainingDetail)
	}
	post := router.Group("post")
	{
		post.POST("list", controllers.PostList)
		post.POST("/detail", controllers.PostDetail)
	}
	position := router.Group("position")
	{
		position.POST("list", controllers.PositionList)
	}
	level := router.Group("level")
	{
		level.POST("list", controllers.LevelList)
	}
	appraisal := router.Group("appraisal")
	{
		appraisal.POST("list", controllers.AppraisalList)
		appraisal.POST("detail", controllers.AppraisalDetail)
	}
	award := router.Group("award")
	{
		award.POST("list", controllers.AwardList)
		award.POST("detail", controllers.AwardDetail)
	}
	punish := router.Group("punish")
	{
		punish.POST("list", controllers.PunishList)
		punish.POST("detail", controllers.PunishDetail)
	}
	discipline := router.Group("discipline")
	{
		discipline.POST("list", controllers.DisciplineList)
		discipline.POST("detail", controllers.DisciplineDetail)
	}
	disDict := router.Group("dis_dict")
	{
		disDict.POST("list", controllers.DisDictList)
	}
	permission := router.Group("permission")
	{
		permission.POST("list", controllers.PermissionList)
		permission.POST("manage", controllers.PermissionManage)
		permission.POST("policy", controllers.GetPolicy)
		permission.POST("check", controllers.PermissionCheck)
		permission.POST("add", controllers.PermissionAdd)
		permission.POST("delete", controllers.PermissionDelete)
	}
	role := router.Group("role")
	{
		role.POST("list", controllers.RoleList)
		role.POST("add", controllers.RoleAdd)
		role.POST("update", controllers.RoleUpdate)
		role.POST("delete", controllers.RoleDelete)
		role.POST("permission", controllers.GetRolePermission)
	}
	roleDict := router.Group("role_dict")
	{
		roleDict.POST("list", controllers.RoleDictList)
		roleDict.POST("add", controllers.RoleDictAdd)
		roleDict.POST("update", controllers.RoleDictUpdate)
		roleDict.POST("delete", controllers.RoleDictDelete)
	}
	data := router.Group("data")
	{
		data.POST("sync", controllers.DataSync)
		data.POST("sure", controllers.DataSure)
	}
	report := router.Group("report")
	{
		report.POST("list", controllers.ReportList)
		report.POST("one", controllers.ReportOne)
		report.POST("detail", controllers.ReportDetail)
		report.POST("steps", controllers.ReportSteps)
		report.POST("add", controllers.ReportAdd)
		report.POST("update", controllers.ReportUpdate)
		report.POST("person_add", controllers.PersonReportAdd)
	}
	entryExit := router.Group("entry_exit")
	{
		entryExit.POST("list", controllers.EntryExitList)
	}
	affair := router.Group("affair")
	{
		affair.POST("list", controllers.AffairList)
		affair.POST("one", controllers.AffairOne)
	}
	return router
}

func Start(r *gin.Engine) {
	//以下是优雅的关机方法
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
