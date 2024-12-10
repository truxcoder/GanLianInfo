package controllers

import (
	"GanLianInfo/models"
	"GanLianInfo/utils"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	log "github.com/truxcoder/truxlog"
)

const (
	timeFormat  = "2006-01-02 15:04:05"
	dayFormat   = "2006-01-02"
	monthFormat = "2006-01"
)

// Personnel 人员同步结构体
type Personnel struct {
	ID           int64        `json:"id,string"`
	UserID       string       `json:"userId"`
	Name         string       `json:"realName"`
	Gender       string       `json:"gender"`
	IdCode       string       `json:"idCode"`
	Birthday     time.Time    `json:"birthday"`
	OrganID      string       `json:"organID"`
	DepartmentID string       `json:"deptId"`
	Username     string       `json:"userName"`
	UserType     int8         `json:"userType"`
	DataStatus   int8         `json:"dataStatus"`
	Phone        string       `json:"phone"`
	Sort         int          `json:"sort"`
	CreateTime   time.Time    `json:"createTime"`
	UpdateTime   time.Time    `json:"updateTime"`
	Same         []PersonSame `json:"same"`
}

type PersonSame struct {
	ID           int64     `json:"id,string"`
	Name         string    `json:"name"`
	Gender       string    `json:"gender"`
	IdCode       string    `json:"idCode"`
	Birthday     time.Time `json:"birthday"`
	OrganID      string    `json:"organID"`
	DepartmentID string    `json:"departmentId"`
	UserType     int8      `json:"userType"`
	Phone        string    `json:"phone"`
}

type PersonnelChanged struct {
	ID           int64 `json:"id,string"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Version      int8
	UserId       string    `json:"userId"`
	Name         string    `json:"realName"`
	Gender       string    `json:"gender"`
	IdCode       string    `json:"idCode"`
	Birthday     time.Time `json:"birthday"`
	OrganID      string    `json:"organID"`
	DepartmentID string    `json:"deptId"`
	UserType     int8      `json:"userType"`
	DataStatus   int8      `json:"dataStatus"`
	Phone        string    `json:"phone"`
	Sort         int       `json:"sort"`
	CreateTime   time.Time `json:"createTime"`
	UpdateTime   time.Time `json:"updateTime"`
}

type Per struct {
	UserID        string    `json:"user_id"`
	Gender        string    `json:"gender"`
	Nation        string    `json:"nation"`
	PoliceCode    string    `json:"police_code"`
	Political     string    `json:"political"`
	JoinPartyDay  time.Time `json:"join_party_day"`
	StartJobDay   time.Time `json:"start_job_day"`
	FullTimeEdu   string    `json:"full_time_edu"`
	FullTimeMajor string    `json:"full_time_major"`
	PartTimeEdu   string    `json:"part_time_edu"`
	BePoliceDay   time.Time `json:"be_police_day"`
}

func (b *PersonnelChanged) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = utils.GenId()
	b.Version = 0
	return
}
func (b *PersonnelChanged) BeforeUpdate(tx *gorm.DB) (err error) {
	b.Version++
	return
}

type PerSlice []Personnel

type Department struct {
	ID          string    `json:"deptId"`
	Name        string    `json:"name"`
	ShortName   string    `json:"shortName"`
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

// DepartmentChanged 多定义这个是因为Department从大数据中心的数据Unmarshal的时候要处理时间。而从前端返回接收的时候时间已经是处理好的了，不需要再处理。
// 不能再利用自定义的Unmarshal接收。所以必须再定义一个字段一样，名称不一样的类型来接收。Personnel也一样。
type DepartmentChanged struct {
	ID          string    `json:"deptId"`
	Name        string    `json:"name"`
	ShortName   string    `json:"shortName"`
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

// UnmarshalJSON 为Personnel类型重写UnmarshalJSON方法
func (d *Personnel) UnmarshalJSON(data []byte) error {
	type TempP Personnel // 定义与Department字段一致的新类型
	dt := struct {
		CreateTime string `json:"createTime"`
		UpdateTime string `json:"updateTime"`
		*TempP            // 避免直接嵌套Department进入死循环
	}{
		TempP: (*TempP)(d),
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

// Len 重写Len方法
func (p PerSlice) Len() int {
	return len(p)
}

// Swap 重写Swap方法
func (p PerSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Less 重写Less方法
func (p PerSlice) Less(i, j int) bool {
	return p[j].Name < p[i].Name
}

func DataSync(c *gin.Context) {
	data := GetPersonnelDataFromInterface()
	var p, added, updated, deleted PerSlice
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, &p)
	if err != nil {
		log.Error(err)
	}
	sort.Sort(p)
	for i := 0; i < len(p); i++ {
		v := p[i]
		var id string
		var personnel models.Personnel
		if v.IdCode == "" {
			continue
		}

		result := db.Model(models.Personnel{}).Where("id_code = ?", v.IdCode).Limit(1).Find(&personnel)
		isFound := result.RowsAffected == 1
		isValid := v.UserType == 1 && v.DataStatus == 0
		isUpdated := v.UpdateTime.After(personnel.UpdateTime)

		if !isFound && !isValid {
			continue
		}
		if isFound && !isUpdated {
			continue
		}
		v.Birthday = utils.GetBirthdayFromIdCode(v.IdCode)
		v.Gender = utils.GetGenderFromIdCode(v.IdCode)

		// 这里在redis里建立一个map, 避免为了查找organID多次查询数据库
		//stmt := `select id from departments where bus_org_code = (select bus_org_code from departments where id = ?) and dept_type = 1`
		//db.Raw(stmt, v.DepartmentID).Scan(&id)
		id = getOrganIdFromDepartmentId(v.DepartmentID)
		v.OrganID = id
		// 过滤掉泸州所和攀枝花所
		if v.OrganID == "c84c0a0ae2e54c5baf8c9d8c86fc9761" || v.OrganID == "6a8f659d05a74ee582c4880083ed606d" {
			continue
		}

		if !isFound && isValid {
			var _add []PersonSame
			db.Table("personnels").Where("name = ?", v.Name).Find(&_add)
			if len(_add) > 0 {
				//added = append(added, _add...)
				v.Same = append(v.Same, _add...)
			}
			added = append(added, v)
		} else if isFound && isValid && isUpdated {
			v.ID = personnel.ID
			updated = append(updated, v)
		} else if isFound && isUpdated {
			// TODO: 这里的删除逻辑需要大数据中心开放身份证验证接口，否则无法实现
			//v.ID = personnel.ID
			//deleted = append(deleted, v)
		}
	}
	// 如果所有数据为空，证明无需添加、删除、更新，则往redis里写入updateTime为当前时间，便于下一步请求大数据中心数据时
	// updateStartTime从当前时间计算
	if len(added) == 0 && len(updated) == 0 && len(deleted) == 0 && rdb != nil {
		//res, _ := rdb.Exists(ctx, "updateTime").Result()
		now := time.Now()
		rdb.Set(ctx, "updateTime", now, time.Hour*2400)
		//if res == 0 {
		//	rdb.HSet(ctx, "personOrganMap", _map)
		//}

	}
	r := gin.H{"code": 20000, "add": &added, "update": &updated, "delete": &deleted}
	c.JSON(200, r)
}

func DataSure(c *gin.Context) {
	var r gin.H
	method := c.Query("method")
	var p []PersonnelChanged
	if method != "delete" {
		if err := c.ShouldBindJSON(&p); err != nil {
			log.Error(err)
			//r = Errors.ServerError
			r = GetError(CodeBind)
			c.JSON(200, r)
			return
		}
	}
	if method == "add" {
		//for _, v := range p {
		//	result = db.Table("personnels").Create(&v)
		//}
		if result := db.Table("personnels").Create(p); result.Error != nil {
			//r = Errors.Insert
			r = GetError(CodeAdd)
			c.JSON(200, r)
			return
		}
		setIdCodeMapToCache()
		r = gin.H{"code": 20000, "message": "添加成功!"}
		c.JSON(200, r)
		return
	}
	if method == "update" {
		for _, v := range p {
			//log.Successf("v:%+v\n", v)
			db.Table("personnels").Omit("birthday").Where("id = ?", v.ID).Updates(&v)
		}
		r = gin.H{"code": 20000, "message": "更新成功!"}
		c.JSON(200, r)
		return
	}
	if method == "delete" {
		var id struct {
			Id []string `json:"id"`
		}
		if err := c.ShouldBindJSON(&id); err != nil {
			log.Error(err)
			//r = Errors.ServerError
			r = GetError(CodeBind)
			c.JSON(200, r)
			return
		}
		db.Where("id_code in ?", &id.Id).Delete(models.Personnel{})
		setIdCodeMapToCache()
		r = gin.H{"code": 20000, "message": "删除成功!"}
		c.JSON(200, r)
		return
	}
}

func DepartmentSync(c *gin.Context) {
	data := GetDepartmentDataFromInterface()
	var departments []Department
	var added, updated, deleted []Department
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, &departments)
	if err != nil {
		log.Error(err)
	}
	for i := 0; i < len(departments); i++ {
		v := departments[i]
		var department models.Department
		result := db.Model(models.Department{}).Where("id = ?", v.ID).Limit(1).Find(&department)
		isFound := result.RowsAffected == 1
		isValid := v.DataStatus == 0
		isUpdated := v.UpdateTime.After(department.UpdateTime)
		if !isFound && !isValid {
			continue
		}
		if isFound && !isUpdated {
			continue
		}

		if !isFound && isValid {
			added = append(added, v)
		} else if isFound && isValid && isUpdated {
			updated = append(updated, v)
		} else if isFound && isUpdated {
			deleted = append(deleted, v)
		}
	}

	r := gin.H{"code": 20000, "add": &added, "update": &updated, "delete": &deleted}
	c.JSON(200, r)
}

func DepartmentSure(c *gin.Context) {
	var r gin.H
	method := c.Query("method")
	var d []DepartmentChanged
	if method != "delete" {
		if err := c.ShouldBindJSON(&d); err != nil {
			log.Error(err)
			//r = Errors.ServerError
			r = GetError(CodeBind)
			c.JSON(200, r)
			return
		}
	}
	if method == "add" {
		if result := db.Table("departments").Create(&d); result.Error != nil {
			//r = Errors.Insert
			r = GetError(CodeAdd)
			c.JSON(200, r)
			return
		}
		setDepartmentMapToCache()
		setDepartmentSliceToCache()
		r = gin.H{"code": 20000, "message": "添加成功!"}
		c.JSON(200, r)
		return
	}
	if method == "update" {
		for _, v := range d {
			//log.Successf("v:%+v\n", v)
			db.Table("departments").Where("id = ?", v.ID).Updates(&v)
		}
		setDepartmentMapToCache()
		setDepartmentSliceToCache()
		r = gin.H{"code": 20000, "message": "更新成功!"}
		c.JSON(200, r)
		return
	}
	if method == "delete" {
		var id struct {
			Id []string `json:"id"`
		}
		if err := c.ShouldBindJSON(&id); err != nil {
			log.Error(err)
			r = GetError(CodeBind)
			c.JSON(200, r)
			return
		}
		db.Where("id in ?", id.Id).Delete(models.Department{})
		setDepartmentMapToCache()
		setDepartmentSliceToCache()
		r = gin.H{"code": 20000, "message": "删除成功!"}
		c.JSON(200, r)
		return
	}
}

func GetPersonnelDataFromInterface() []byte {
	url := "http://30.29.2.6:8686/unionapi/user/list/json"
	contentType := "application/x-www-form-urlencoded"
	baseTime := time.Date(2024, 7, 11, 12, 0, 0, 0, time.Local)
	if rdb != nil {
		res, _ := rdb.Exists(ctx, "updateTime").Result()
		if res > 0 {
			temp, _ := rdb.Get(ctx, "updateTime").Result()
			updateTime, err := time.Parse(time.RFC3339, temp)
			if err != nil {
				log.Error(err)
			} else {
				baseTime = updateTime
				//log.Successf("baseTime:%v\n", baseTime)
			}
		}
	}
	data := "Authorization=438019355f6940fba3b98316d97fd5f0&updateStartTime=" + baseTime.Format(timeFormat)
	resp, err := http.Post(url, contentType, strings.NewReader(data))
	if err != nil {
		fmt.Printf("post failed, err:%v\n", err)
		return nil
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("get resp failed, err:%v\n", err)
		return nil
	}
	list := jsoniter.Get(b, "data").ToString()
	return []byte(list)
}

func GetPersonnelDataFromInterfaceForAccount() []byte {
	url := "http://30.29.2.6:8686/unionapi/user/list/json"
	contentType := "application/x-www-form-urlencoded"
	baseTime := time.Date(2022, 1, 21, 12, 0, 0, 0, time.Local)
	if rdb != nil {
		res, _ := rdb.Exists(ctx, "accountUpdateTime").Result()
		if res > 0 {
			temp, _ := rdb.Get(ctx, "accountUpdateTime").Result()
			updateTime, err := time.Parse(time.RFC3339, temp)
			if err != nil {
				log.Error(err)
			} else {
				baseTime = updateTime
				//log.Successf("GetPersonnelDataFromInterfaceForAccount baseTime:%v\n", baseTime)
			}
		}
	}
	data := "Authorization=438019355f6940fba3b98316d97fd5f0&updateStartTime=" + baseTime.Format(timeFormat)
	resp, err := http.Post(url, contentType, strings.NewReader(data))
	if err != nil {
		fmt.Printf("post failed, err:%v\n", err)
		return nil
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("get resp failed, err:%v\n", err)
		return nil
	}
	list := jsoniter.Get(b, "data").ToString()
	return []byte(list)
}

func GetDepartmentDataFromInterface() []byte {
	url := "http://30.29.2.6:8686/unionapi/dept/list/json"
	contentType := "application/x-www-form-urlencoded"
	data := "Authorization=438019355f6940fba3b98316d97fd5f0&foo=bar"
	resp, err := http.Post(url, contentType, strings.NewReader(data))
	if err != nil {
		fmt.Printf("post failed, err:%v\n", err)
		return nil
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("get resp failed, err:%v\n", err)
		return nil
	}
	list := jsoniter.Get(b, "data").ToString()
	return []byte(list)
}

func DataSyncText(c *gin.Context) {
	id := getOrganIdFromDepartmentId("000166d685844db9ad0cf4aade9f7528")
	log.Successf("id:%s\n", id)
}
