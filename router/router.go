package router

import (
	"GanLianInfo/auth"
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
			r := controllers.GetError(controllers.CodeDatabase)
			c.JSON(200, r)
			c.Abort()
		}
	}
}

func Register() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	authMiddleware := auth.JWTAuthMiddleware()
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
	router.POST("/add", authMiddleware, controllers.Add)
	router.POST("/update", authMiddleware, controllers.Update)
	router.POST("/delete", authMiddleware, controllers.Delete)
	router.POST("/dashboard", authMiddleware, controllers.DashboardData)
	router.POST("/upload", authMiddleware, controllers.Upload)
	router.POST("/pre", authMiddleware, controllers.PreEdit)
	router.POST("/pre_batch", authMiddleware, controllers.PreBatchEdit)
	user := router.Group("user")
	{
		user.POST("/login", controllers.Login)
		user.POST("/info", controllers.UserInfo)
	}
	module := router.Group("module", authMiddleware)
	{
		module.POST("/list", controllers.GetModuleList)
		module.POST("/order", controllers.ModuleOrder)
		module.POST("/role", controllers.ModuleRole)
	}
	personnel := router.Group("personnel", authMiddleware)
	{
		personnel.POST("/list", controllers.PersonnelList)
		personnel.POST("/detail", controllers.PersonnelDetail)
		personnel.POST("/update", controllers.PersonnelUpdate)
		//personnel.POST("/delete", controllers.PersonnelDelete)
		personnel.DELETE("/:id", controllers.PersonnelDelete)
		personnel.POST("/export_list", controllers.PersonnelExportList)
		personnel.POST("/base_list", controllers.PersonnelBaseList)
		personnel.POST("/name_list", controllers.PersonnelNameList)
		personnel.POST("/resume", controllers.PersonnelResume)
		personnel.POST("/update_edu", controllers.PersonnelUpdateEdu)
		personnel.POST("/update_id_code", controllers.UpdateIdCode)
		personnel.POST("/update_birthday", controllers.UpdateBirthday)
		personnel.POST("/update_status", controllers.PersonnelUpdateStatus)
		personnel.POST("/organ", controllers.GetPersonOrganId)
		personnel.POST("/organs", controllers.GetPersonOrgans)
	}
	custom := router.Group("custom", authMiddleware)
	{
		custom.POST("list", controllers.CustomList)
	}
	router.POST("/personnel/dict", controllers.EduDictList)
	//organ := router.Group("organ")
	//{
	//	//organ.GET("list", controllers.GetOrganList)
	//	organ.POST("/add", controllers.OrganAdd)
	//	organ.POST("/update", controllers.OrganUpdate)
	//	organ.POST("/delete", controllers.OrganDelete)
	//}
	department := router.Group("department", authMiddleware)
	{
		department.POST("list", controllers.DepartmentList)
		department.POST("organ", controllers.OrganList)
		department.POST("headcount", controllers.HeadcountList)
		department.POST("update", controllers.DepartmentUpdate)
		department.POST("position", controllers.DepartmentPosition)
	}
	training := router.Group("training", authMiddleware)
	{
		training.POST("list", controllers.TrainingList)
		training.POST("person_list", controllers.TrainPersonList)
		training.POST("add", controllers.TrainPersonAdd)
		training.POST("delete", controllers.TrainPersonDelete)
		training.POST("detail", controllers.TrainingDetail)
	}
	post := router.Group("post", authMiddleware)
	{
		post.POST("list", controllers.PostList)
		post.POST("detail", controllers.PostDetail)
		post.POST("add", controllers.PostAdd)
		post.POST("update", controllers.PostUpdate)
	}
	position := router.Group("position", authMiddleware)
	{
		position.POST("list", controllers.PositionList)
		position.POST("check", controllers.PositionCheck)
	}
	appointment := router.Group("appointment", authMiddleware)
	{
		appointment.POST("list", controllers.AppointmentList)
		appointment.POST("table", controllers.AppointmentTableDetail)
	}
	level := router.Group("level", authMiddleware)
	{
		level.POST("list", controllers.LevelList)
	}
	appraisal := router.Group("appraisal", authMiddleware)
	{
		appraisal.POST("list", controllers.AppraisalList)
		appraisal.POST("detail", controllers.AppraisalDetail)
		appraisal.POST("batch", controllers.AppraisalBatch)
		appraisal.POST("pre_batch", controllers.AppraisalPreBatch)
	}
	award := router.Group("award", authMiddleware)
	{
		award.POST("list", controllers.AwardList)
		award.POST("detail", controllers.AwardDetail)
		award.POST("batch", controllers.AwardBatch)
		award.POST("pre_batch", controllers.AwardPreBatch)
	}
	punish := router.Group("punish", authMiddleware)
	{
		punish.POST("list", controllers.PunishList)
		punish.POST("detail", controllers.PunishDetail)
	}
	discipline := router.Group("discipline", authMiddleware)
	{
		discipline.POST("list", controllers.DisciplineList)
		discipline.POST("detail", controllers.DisciplineDetail)
	}
	disDict := router.Group("dis_dict", authMiddleware)
	{
		disDict.POST("list", controllers.DisDictList)
	}
	permission := router.Group("permission", authMiddleware)
	{
		permission.POST("list", controllers.PermissionList)
		permission.POST("manage", controllers.PermissionManage)
		permission.POST("policy", controllers.GetPolicy)
		permission.POST("check", controllers.PermissionCheck)
		permission.POST("act_check", controllers.PermissionActCheck)
		permission.POST("add", controllers.PermissionAdd)
		permission.POST("delete", controllers.PermissionDelete)
	}
	role := router.Group("role", authMiddleware)
	{
		role.POST("list", controllers.RoleList)
		role.POST("add", controllers.RoleAdd)
		role.POST("update", controllers.RoleUpdate)
		role.POST("delete", controllers.RoleDelete)
		role.POST("permission", controllers.GetRolePermission)
	}
	roleDict := router.Group("role_dict", authMiddleware)
	{
		roleDict.POST("list", controllers.RoleDictList)
		roleDict.POST("add", controllers.RoleDictAdd)
		roleDict.POST("update", controllers.RoleDictUpdate)
		roleDict.POST("delete", controllers.RoleDictDelete)
	}
	data := router.Group("data", authMiddleware)
	{
		data.POST("sync", controllers.DataSync)
		data.POST("department_sync", controllers.DepartmentSync)
		data.POST("account_sync", controllers.AccountSync)
		data.POST("sure", controllers.DataSure)
		data.POST("department_sure", controllers.DepartmentSure)
		data.POST("account_sure", controllers.AccountSure)
		data.POST("test", controllers.DataSyncText)
	}
	report := router.Group("report", authMiddleware)
	{
		report.POST("list", controllers.ReportList)
		report.POST("one", controllers.ReportOne)
		report.POST("detail", controllers.ReportDetail)
		report.POST("steps", controllers.ReportSteps)
		report.POST("add", controllers.ReportAdd)
		report.POST("update", controllers.ReportUpdate)
		report.POST("person_add", controllers.PersonReportAdd)
	}
	entryExit := router.Group("entry_exit", authMiddleware)
	{
		entryExit.POST("list", controllers.EntryExitList)
	}
	affair := router.Group("affair", authMiddleware)
	{
		affair.POST("list/:category", controllers.AffairList)
		affair.POST("detail", controllers.AffairDetail)
		affair.POST("one", controllers.AffairOne)
	}
	family := router.Group("family", authMiddleware)
	{
		family.POST("detail", controllers.FamilyDetail)
	}
	account := router.Group("account", authMiddleware)
	{
		account.POST("list", controllers.AccountList)
		account.POST("base_list", controllers.AccountBaseList)
	}
	review := router.Group("review", authMiddleware)
	{
		review.POST("list", controllers.ReviewList)
		review.POST("feedback", controllers.FeedbackList)
		review.POST("pass", controllers.ReviewPass)
	}
	analysis := router.Group("analysis", authMiddleware)
	{
		analysis.POST("police", controllers.AnalysisPoliceTeamData)
		analysis.POST("leader", controllers.AnalysisLeaderTeamData)
	}
	leader := router.Group("leader", authMiddleware)
	{
		leader.POST("list", controllers.LeaderList)
	}
	talent := router.Group("talent", authMiddleware)
	{
		talent.POST("list/:category", controllers.TalentList)
		talent.POST("add", controllers.TalentAdd)
		talent.POST("pick_list", controllers.TalentPickList)
		talent.POST("pick_add", controllers.TalentPickAdd)
		talent.POST("detail_list", controllers.TalentDetailList)
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
