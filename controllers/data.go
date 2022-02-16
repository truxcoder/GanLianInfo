package controllers

import (
	"GanLianInfo/models"
	"GanLianInfo/utils"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	log "github.com/truxcoder/truxlog"
)

const (
	timeFormat  = "2006-01-02 15:04:05"
	dayFormat   = "2006-01-02"
	monthFormat = "2006-01"
)

type Personnel struct {
	ID           string    `json:"userId"`
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

type PersonnelChanged struct {
	ID           string    `json:"userId"`
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
	ID            string    `json:"user_id"`
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

func DataSync(c *gin.Context) {
	var id string
	data := GetPersonnelDataFromInterface()
	var p []Personnel
	var added, updated, deleted []Personnel
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, &p)
	if err != nil {
		log.Error(err)
	}
	for i := 0; i < len(p); i++ {
		v := p[i]
		var personnel models.Personnel
		result := db.Model(models.Personnel{}).Where("id = ?", v.ID).Limit(1).Find(&personnel)
		if result.RowsAffected == 0 && (v.UserType == 1 || v.UserType == 2) && v.DataStatus == 0 {
			if v.IdCode != "" {
				v.Birthday = utils.GetBirthdayFromIdCode(v.IdCode)
				v.Gender = utils.GetGenderFromIdCode(v.IdCode)
			}
			stmt := `select id from departments where bus_org_code = (select bus_org_code from departments where id = ?) and dept_type = 1`
			db.Raw(stmt, v.DepartmentID).Scan(&id)
			v.OrganID = id
			added = append(added, v)
		} else if result.RowsAffected == 1 && (v.UserType == 1 || v.UserType == 2) && v.DataStatus == 0 && v.UpdateTime.After(personnel.UpdateTime) {
			if v.IdCode != "" {
				v.Birthday = utils.GetBirthdayFromIdCode(v.IdCode)
				v.Gender = utils.GetGenderFromIdCode(v.IdCode)
			}
			stmt := `select id from departments where bus_org_code = (select bus_org_code from departments where id = ?) and dept_type = 1`
			db.Raw(stmt, v.DepartmentID).Scan(&id)
			v.OrganID = id
			updated = append(updated, v)
		} else if result.RowsAffected == 1 && v.UpdateTime.After(personnel.UpdateTime) {
			if v.IdCode != "" {
				v.Birthday = utils.GetBirthdayFromIdCode(v.IdCode)
				v.Gender = utils.GetGenderFromIdCode(v.IdCode)
			}
			stmt := `select id from departments where bus_org_code = (select bus_org_code from departments where id = ?) and dept_type = 1`
			db.Raw(stmt, v.DepartmentID).Scan(&id)
			v.OrganID = id
			deleted = append(deleted, v)
		}
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
			r = Errors.ServerError
			c.JSON(200, r)
			return
		}
	}
	if method == "add" {
		for _, v := range p {
			db.Table("personnels").Create(&v)
		}
		r = gin.H{"code": 20000, "message": "添加成功!"}
		c.JSON(200, r)
		return
	}
	if method == "update" {
		for _, v := range p {
			db.Table("personnels").Updates(&v)
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
			r = Errors.ServerError
			c.JSON(200, r)
			return
		}
		db.Delete(models.Personnel{}, &id.Id)
		r = gin.H{"code": 20000, "message": "删除成功!"}
		c.JSON(200, r)
		return
	}
}

func GetPersonnelDataFromInterface() []byte {
	url := "http://30.29.2.6:8686/unionapi/user/list/json"
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
