package controllers

import (
	"GanLianInfo/models"

	"github.com/gin-gonic/gin"
	log "github.com/truxcoder/truxlog"
)

func InitDb(c *gin.Context) {
	var err error
	err = db.Migrator().DropTable(&models.Personnel{}, &models.Module{}, &models.Department{})
	if err != nil {
		log.Error(err)
	}
	err = db.AutoMigrate(&models.Personnel{}, &models.Module{}, &models.Department{})
	if err != nil {
		log.Error(err)
	}
	//db.AutoMigrate( &models.Department{})
	log.Info(db.Migrator().CurrentDatabase())
	//if err != nil {
	//	log.Error(err)
	//}
	//db := dao.Connect()
	//db.Migrator().DropTable(&models.Personnel{})
	//err := db.AutoMigrate(&models.Personnel{})
	//if err != nil {
	//	fmt.Println("数据迁移错误：",err)
	//	return
	//}
	//r := gin.H{"code": 20000,"data": "数据迁移正常"}
	//c.JSON(200, r)
	//asset := models.Asset{Code:"12345678",Name:"交换机",Brand: "华为",ProductModel: "s5720",BuyTime: "2020-10-22",UseTime: "2019-11-11",ScrapYear: 8,Price: 1520.5,Status:"正常",Position: "三大队车间",Manager: "李晓波"}
	//result := db.Create(&asset)
}

func AddD(c *gin.Context) {
	//engine := dao.Connect()
	//defer engine.Close()
	//s := engine.New()
	//personnel := models.Personnel{
	//	Code:  "3232322-23223232",
	//	Name:  "李晓波",
	//	Phone: "18080565555",
	//}
	//err := personnel.BeforeCreate()
	//if err != nil {
	//	fmt.Printf("BeforeCreate err: %v\n",err)
	//}
	//stmt, err := db.Prepare(`insert into personnel(id,name,code,phone,createdAt,updatedAt,Version) values (?,?,?,?,?,?,?)`)
	//if err != nil {
	//	fmt.Printf("err: %v\n",err)
	//}
	//result,err :=stmt.Exec(personnel.ID,personnel.Name,personnel.Code,personnel.Phone,personnel.CreatedAt,
	//	personnel.UpdatedAt,personnel.Version)
	//if err != nil {
	//	fmt.Printf("err: %v\n",err)
	//}
	//fmt.Printf("result: %v\n",result)
}

func Fake(c *gin.Context) {
	//db := dao.Connect()
	//	//db.MapperFunc(strings.ToUpper)
	//	var personnel []models.Personnel
	//	//stmt, _ := db.Prepare("select * from personnel")
	//	//rows, err := db.Queryx("select * from personnel")
	//	err := db.Select(&personnel, "select * from personnel")
	//	if err != nil {
	//		fmt.Printf("数据库查询错误，err: %v\n",err)
	//	}
	//	defer db.Close()
	//for key,val := range personnel {
	//	fmt.Printf("%v-%+v\n",key,val)
	//}
	//	r := gin.H{"code":20000,"data":personnel}
	//	c.JSON(200, r)
}
