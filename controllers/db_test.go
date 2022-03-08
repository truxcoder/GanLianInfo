package controllers

import (
	"GanLianInfo/models"
	"GanLianInfo/utils"
	"faker"
	"fmt"
	"go/ast"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"

	"gorm.io/gorm"

	log "github.com/truxcoder/truxlog"
)

type Student struct {
	Personnel
	Code string
}

type Field struct {
	Name string
	Type string
	Tag  string
}

// Schema 表示数据库一张表
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

type Department struct {
	ID          string    `json:"deptId" gorm:"size:50;primaryKey"`
	Name        string    `json:"name" gorm:"size:50;not Null"`
	ShortName   string    `json:"shortName" gorm:"size:50"`
	DeptType    int       `json:"deptType"`
	DataStatus  int       `json:"dataStatus"`
	Code        string    `json:"code"`
	LevelCode   string    `json:"levelCode"`
	BusOrgCode  string    `json:"busOrgCode"`
	BusDeptCode string    `json:"busDeptCode"`
	ParentId    string    `json:"parentId"`
	Sort        int       `json:"sort"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}

func (d *Department) UnmarshalJSON(data []byte) error {
	type TempD Department // 定义与Department字段一致的新类型
	dt := struct {
		CreateTime string `json:"createTime"`
		UpdateTime string `json:"updateTime"`
		*TempD            // 避免直接嵌套Department进入死循环
	}{
		TempD: (*TempD)(d),
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(data, &dt); err != nil {
		return err
	}
	var err error
	d.CreateTime, err = time.Parse(timeFormat, dt.CreateTime)
	if err != nil {
		return err
	}
	d.UpdateTime, err = time.Parse(timeFormat, dt.UpdateTime)
	if err != nil {
		return err
	}
	return nil
}

func TestDepartmentSync(t *testing.T) {
	url := "http://30.29.2.6:8686/unionapi/dept/list/json"
	contentType := "application/x-www-form-urlencoded"
	data := "Authorization=438019355f6940fba3b98316d97fd5f0&foo=bar"
	var departments []Department

	resp, err := http.Post(url, contentType, strings.NewReader(data))
	if err != nil {
		fmt.Printf("post failed, err:%v\n", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("get resp failed, err:%v\n", err)
	}

	list := jsoniter.Get(b, "data").ToString()
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if err = json.Unmarshal([]byte(list), &departments); err != nil {
		log.Error(err)
	}
	for _, v := range departments {
		if v.DataStatus == 0 {
			db.Create(&v)
			//log.Successf("depart: %+v", v)
		}
	}
}

func TestAssociation(t *testing.T) {
	var p []models.Appraisal
	result := db.Debug().Preload("Level", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).Find(&p)

	log.Error(result)

	log.Successf("result: %+v\n", p)

}

func TestInsert(t *testing.T) {
	var p models.Personnel
	faker.New(&p)
	db.Create(&p)
	var person models.Personnel
	log.Info("currentID:", p.ID)
	db.First(&person, "id = ?", p.ID)
	T := reflect.TypeOf(person)
	V := reflect.ValueOf(person)
	for i := 0; i < T.NumField(); i++ {
		fmt.Printf("%s:%+v\n", T.Field(i).Name, V.Field(i).Interface())
	}
}

func TestInsertPersonnel(t *testing.T) {
	//var idList []string

	// ****本地读取json文件方案开始
	//type Department struct {
	//	Id         string `json:"deptId"`
	//	Name       string `json:"name"`
	//	ShortName  string `json:"shortName"`
	//	DeptType   int    `json:"deptType"`
	//	DataStatus int    `json:"dataStatus"`
	//	BusOrgCode string `json:"busOrgCode"`
	//	ParentId   string `json:"parentId"`
	//}
	//bytes, _ := ioutil.ReadFile("departmentData.json")
	//var dList, organList []string
	//var departments []Department
	//err := json.Unmarshal(bytes, &departments)
	//if err != nil {
	//	log.Error(err)
	//}
	//for i := 0; i < len(departments); i++ {
	//	if departments[i].DataStatus == 0 && departments[i].DeptType == 0 {
	//		dList = append(dList, departments[i].Id)
	//	} else if departments[i].DataStatus == 0 && departments[i].DeptType == 1 {
	//		organList = append(organList, departments[i].Id)
	//	}
	//}
	// ****本地读取json文件方案结束

	var organ, department []models.Department
	var result *gorm.DB
	result = db.Select("id", "name", "short_name", "bus_org_code").Where("dept_type = ?", 1).Find(&organ)
	organLength := result.RowsAffected
	if result.Error != nil {
		log.Error(result.Error)
	}
	result = db.Select("id").Where("dept_type = ?", 0).Find(&department)
	if result.Error != nil {
		log.Error(result.Error)
	}
	deptLength := result.RowsAffected

	for i := 0; i < 25; i++ {
		var p models.Personnel
		faker.New(&p)
		p.OrganID = organ[rand.Int63n(organLength)].ID
		p.DepartmentId = department[rand.Int63n(deptLength)].ID
		p.Birthday = utils.GetBirthdayFromIdCode(p.IdCode)
		db.Create(&p)
	}

	//decoder := json.NewDecoder(filePtr)
	//for decoder.More() {
	//	var department Department
	//	err = decoder.Decode(&department)
	//	log.Infof("department:%+v\n", department)
	//}

	//var o []models.Organ
	//result := db.Find(&o)
	//log.Successf("result: %d\n", result.RowsAffected)
	//for i := 0; i < len(o); i++ {
	//	idList = append(idList, o[i].ID)
	//}

	//for i := 0; i < 25; i++ {
	//	var p models.Personnel
	//	faker.New(&p)
	//	p.OrganId = faker.Choice(organList)
	//	p.DepartmentID = faker.Choice(dList)
	//	p.Birthday = utils.GetBirthdayFromIdCode(p.IdCode)
	//	db.Create(&p)
	//}
}

//func TestInsertOrganData(t *testing.T) {
//	nameList := []string{"四川省资阳强制隔离戒毒所", "四川省绵阳强制隔离戒毒所", "四川省眉山强制隔离戒毒所", "四川省女子强制隔离戒毒所"}
//	shortNameList := []string{"资阳所", "绵阳所", "眉山所", "女子所"}
//	stmt := "INSERT INTO ORGANS (id, created_at, updated_at, version, name, short_name, parent, authorized_size, orders) VALUES (?,?,?,?,?,?,?,?,?)"
//	db.Exec(stmt, faker.UUID(), time.Now(), time.Now(), 0, "四川省戒毒管理局", "戒毒局", "0", 1000, 1)
//	var organ models.Organ
//	db.First(&organ)
//	firstId := organ.ID
//	for i := 0; i < len(nameList); i++ {
//		var o models.Organ
//		o.Name = nameList[i]
//		o.ShortName = shortNameList[i]
//		o.Parent = firstId
//		o.AuthorizedSize = faker.Number("200,600")
//		o.Orders = 2 + i
//		db.Create(&o)
//	}
//}

func TestGetBirthday(t *testing.T) {
	for i := 0; i < 10; i++ {
		c := faker.IdCode()
		log.Infof("IdCode:%s----Birthday:%+v\n", c, utils.GetBirthdayFromIdCode(c))
	}
}

func TestConstraint(t *testing.T) {
	sql := "ALTER TABLE disciplines ADD CONSTRAINT fk_DisDisDicts FOREIGN KEY (dict_id) REFERENCES dis_dicts(id) ON UPDATE CASCADE ON DELETE SET NULL;"
	log.Success(sql)
}

func TestFakeModuleData(t *testing.T) {
	stmt := `INSERT INTO "GANLIAN"."modules"("id","created_at","updated_at","version","name","title","paths","rank","component","redirect","icon","parent","orders") VALUES(4376914533814279,'2021-10-14 11:00:05.596000 +08:00','2021-11-03 10:55:50.709434 +08:00',1,'Module','模块管理','module',2,'Module','','el-icon-video-camera-solid',4376914533814290,11);
INSERT INTO "GANLIAN"."modules"("id","created_at","updated_at","version","name","title","paths","rank","component","redirect","icon","parent","orders") VALUES(4376914533814280,'2021-10-13 10:41:30.722000 +08:00','2021-11-11 17:43:42.572501 +08:00',3,'Organization','机构管理','/organ',1,'Layout','/organ/organ','el-icon-school',0,1);
INSERT INTO "GANLIAN"."modules"("id","created_at","updated_at","version","name","title","paths","rank","component","redirect","icon","parent","orders") VALUES(4376914533814281,'2021-10-14 11:00:05.596000 +08:00','2021-11-03 10:55:50.714444 +08:00',1,'Permission','权限管理','permission',2,'Permission','','el-icon-help',4376914533814290,12);
INSERT INTO "GANLIAN"."modules"("id","created_at","updated_at","version","name","title","paths","rank","component","redirect","icon","parent","orders") VALUES(4376914533814282,'2021-10-14 11:00:05.596000 +08:00','2021-11-03 10:55:50.721811 +08:00',2,'Other','其他项目','/other',1,'Layout','/other/fixedasset','el-icon-s-goods',0,20);
INSERT INTO "GANLIAN"."modules"("id","created_at","updated_at","version","name","title","paths","rank","component","redirect","icon","parent","orders") VALUES(4376914533814283,'2021-11-10 08:43:28.158741 +08:00','2021-11-10 08:44:05.862704 +08:00',1,'Personnel','人员信息','personnel',2,'Personnel','','el-icon-user-solid',4376914533814284,25);
INSERT INTO "GANLIAN"."modules"("id","created_at","updated_at","version","name","title","paths","rank","component","redirect","icon","parent","orders") VALUES(4376914533814284,'2021-10-14 11:05:27.227000 +08:00','2021-11-10 08:43:44.491025 +08:00',1,'Person','人员管理','/personnel',1,'Layout','/personnel/personnel','el-icon-user-solid',0,5);
INSERT INTO "GANLIAN"."modules"("id","created_at","updated_at","version","name","title","paths","rank","component","redirect","icon","parent","orders") VALUES(4376914533814285,'2021-10-14 11:00:05.596000 +08:00','2021-11-11 17:45:41.925086 +08:00',2,'Organ','机构管理','organ',2,'Organ','','el-icon-s-management',4376914533814280,2);
INSERT INTO "GANLIAN"."modules"("id","created_at","updated_at","version","name","title","paths","rank","component","redirect","icon","parent","orders") VALUES(4376914533814286,'2021-11-18 17:34:04.249652 +08:00','2021-11-18 17:34:37.828962 +08:00',1,'Appraisal','考核管理','appraisal',2,'Appraisal','','el-icon-document-checked',4376914533814284,29);
INSERT INTO "GANLIAN"."modules"("id","created_at","updated_at","version","name","title","paths","rank","component","redirect","icon","parent","orders") VALUES(4376914533814287,'2021-10-14 11:00:05.596000 +08:00','2021-11-03 10:55:50.691038 +08:00',0,'Position','职数配置','position',2,'Position','','el-icon-phone',4376914533814280,3);
INSERT INTO "GANLIAN"."modules"("id","created_at","updated_at","version","name","title","paths","rank","component","redirect","icon","parent","orders") VALUES(4376914533814288,'2021-11-15 15:40:59.046977 +08:00','2021-11-15 15:41:26.696675 +08:00',1,'Level','职务等级管理','level',2,'Level','','el-icon-s-data',4376914533814284,28);
INSERT INTO "GANLIAN"."modules"("id","created_at","updated_at","version","name","title","paths","rank","component","redirect","icon","parent","orders") VALUES(4376914533814289,'2021-11-16 12:50:28.304749 +08:00','2021-11-16 12:50:28.304749 +08:00',0,'Post','任职管理','post',2,'Post','','el-icon-star-on',4376914533814284,26);
INSERT INTO "GANLIAN"."modules"("id","created_at","updated_at","version","name","title","paths","rank","component","redirect","icon","parent","orders") VALUES(4376914533814290,'2021-10-14 11:00:05.596000 +08:00','2021-11-03 10:55:50.703205 +08:00',0,'System','系统管理','/system',1,'Layout','/system/user','el-icon-s-tools',0,10);
INSERT INTO "GANLIAN"."modules"("id","created_at","updated_at","version","name","title","paths","rank","component","redirect","icon","parent","orders") VALUES(4376914533814291,'2021-10-14 11:00:05.596000 +08:00','2021-11-03 10:55:50.727858 +08:00',4,'User','用户管理','user',2,'User','','el-icon-user-solid',4376914533814282,22);
INSERT INTO "GANLIAN"."modules"("id","created_at","updated_at","version","name","title","paths","rank","component","redirect","icon","parent","orders") VALUES(4376914533814292,'2021-11-15 16:50:29.963478 +08:00','2021-11-15 16:50:29.963478 +08:00',0,'Position','职务管理','position',2,'Position','','el-icon-s-custom',4376914533814284,27);

`
	db.Exec(stmt)
}

func TestMigrate(t *testing.T) {
	var err error
	var modelList = []interface{}{&models.Personnel{}, &models.Department{}, &models.Resume{}, &models.Appraisal{},
		&models.Post{}, &models.Position{}, &models.Level{}, &models.Award{}, &models.Punish{}, &models.Module{}, &models.RoleDict{},
		&models.PermissionDict{}, &models.Discipline{}, &models.DisDict{}}
	err = db.Migrator().DropTable(modelList...)
	if err != nil {
		log.Error(err)
	}
	err = db.AutoMigrate(modelList...)
	if err != nil {
		log.Error(err)
	}
	//db.AutoMigrate( &models.Department{})
	log.Info(db.Migrator().CurrentDatabase())
	//if err != nil {
	//	log.Error(err)
	//}
}
func TestCreateDepartmentTable(t *testing.T) {
	var err error
	err = db.Migrator().DropTable(&models.Department{})
	if err != nil {
		log.Error(err)
	}
	err = db.AutoMigrate(&models.Department{})
	if err != nil {
		log.Error(err)
	}
}

func TestGenPost(t *testing.T) {
	var posts []struct {
		ID     int64     `json:"id"`
		EndDay time.Time `json:"endDay"`
	}
	db.Table("posts").Select("id,end_day").Find(&posts)
	var zero time.Time
	var total int
	for _, v := range posts {
		if v.EndDay.Year() > 2022 {
			log.Successf("id:%d, endDay: %v \n", v.ID, v.EndDay)
			db.Table("posts").Where("id = ?", v.ID).Update("end_day", zero)
			total++
		}
		//db.Table("posts").Updates(v)
	}
	log.Successf("total: %d\n", total)
}

func TestDMText(t *testing.T) {
	var txt models.TextSize

	//var txt map[string]interface{}
	db.Table("text_sizes").First(&txt)
	//db.Raw("select name from text_sizes").Scan(name)
	//n, _ := txt.Name.GetLength()
	//str, _ := txt.Name.ReadString(1, int(n))
	log.Successf("txt:%+v", txt)
}

func TestSearch(t *testing.T) {
	var personnel models.Personnel
	var p struct {
		OrganId string `json:"organID"`
	}
	p.OrganId = "222"
	db.Debug().Model(models.Personnel{}).Where(&p).Where("organ_id in ?", []string{"999", "888"}).Find(&personnel)
}

func TestCreateTable(t *testing.T) {
	var err error
	err = db.Migrator().DropTable(&models.Affair{})
	if err != nil {
		log.Error(err)
	}
	err = db.AutoMigrate(&models.Affair{})
	if err != nil {
		log.Error(err)
	}
}

func TestCreatePostTable(t *testing.T) {
	var err error
	err = db.Migrator().DropTable(&models.PersonReport{})
	if err != nil {
		log.Error(err)
	}
	err = db.AutoMigrate(&models.PersonReport{})
	if err != nil {
		log.Error(err)
	}
}

func TestCreateLevelPositionTable(t *testing.T) {
	var err error
	err = db.Migrator().DropTable(&models.Position{}, &models.Level{})
	if err != nil {
		log.Error(err)
	}
	err = db.AutoMigrate(&models.Position{}, &models.Level{})
	if err != nil {
		log.Error(err)
	}
}

func TestCreateModuleTable(t *testing.T) {
	var err error
	err = db.Migrator().DropTable(&models.Module{})
	if err != nil {
		log.Error(err)
	}
	err = db.AutoMigrate(&models.Module{})
	if err != nil {
		log.Error(err)
	}
}

func TestCreateDictTable(t *testing.T) {
	var err error
	err = db.Migrator().DropTable(&models.Discipline{})
	if err != nil {
		log.Error(err)
	}
	err = db.AutoMigrate(&models.Discipline{})
	if err != nil {
		log.Error(err)
	}
}

func TestReflect(t *testing.T) {
	//type SmallString string
	//var temp SmallString
	//temp = "oookkk"
	//tp := reflect.ValueOf(temp).Type()
	//log.Success(reflect.TypeOf(temp).Name())
	//log.Success(tp.Kind())

	s := Student{
		Personnel: Personnel{Name: "张飞"},
		Code:      "5159373",
	}
	//log.Success(reflect.TypeOf(s).Kind().String() == "struct")
	//log.Success(reflect.Indirect(reflect.ValueOf(s)).Type().String())
	modelType := reflect.Indirect(reflect.ValueOf(&s)).Type()
	schema := &Schema{
		Model:    &s,
		Name:     modelType.Name(), //modelType.Name() 获取到结构体的名称作为表名。
		fieldMap: make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name, //p.Name 即字段名，p.Type 即字段类型，通过 (Dialect).DataTypeOf() 转换为数据库的字段类型
				Type: p.Type.Name(),
			}
			if v, ok := p.Tag.Lookup("truxorm"); ok {
				field.Tag = v //p.Tag 即额外的约束条件
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		} else if p.Anonymous {
			for j := 0; j < p.Type.NumField(); j++ {
				pp := p.Type.Field(j)
				field := &Field{
					Name: pp.Name,
					Type: pp.Type.Name(),
				}
				if v, ok := pp.Tag.Lookup("truxorm"); ok {
					field.Tag = v //p.Tag 即额外的约束条件
				}
				schema.Fields = append(schema.Fields, field)
				schema.FieldNames = append(schema.FieldNames, pp.Name)
				schema.fieldMap[pp.Name] = field
			}
			log.Successf("P.type:%v\n", p.Type.NumField())
		}

	}
	log.Successf("Name:%s\nschemaFieldNames:%+v\nFields:%+v\n", schema.Name, schema.FieldNames, schema.Fields)
	for _, v := range schema.Fields {
		log.Infof("Filed:%+v\n", v)
	}
}

func TestStrings(t *testing.T) {
	aaa := reflect.TypeOf(models.Level{})
	m := reflect.New(aaa).Interface()
	fmt.Printf("%T", m)
	//log.Infof("%t", reflect.New(aaa).Interface())
}

func TestSlice(t *testing.T) {
	list := []string{"32424242423", "23424243242"}
	list2 := []int{2, 3, 4}
	num := 5
	Write(list)
	Write(list2)
	Write(num)
	//result,ok := list.([]interface{})
	//log.Successf("类型：%+v\n", reflect.TypeOf(list).Kind()==reflect.Slice)
	//log.Successf("元素类型：%+v\n", reflect.TypeOf(list2).Elem().String())
}

func Write(v interface{}) {
	switch v.(type) {
	case string:
		s := v.(string)
		log.Infof("%T\n", s)
	case int:
		i := v.(int)
		log.Infof("%T\n", i)
	default:
		log.Infof("%T\n", v)
	}
}

func PrintStruct(m interface{}) {
	T := reflect.Indirect(reflect.ValueOf(m)).Type()
	V := reflect.Indirect(reflect.ValueOf(m))
	for i := 0; i < T.NumField(); i++ {
		fmt.Printf("%s:%+v\n", T.Field(i).Name, V.Field(i).Interface())
	}
}
