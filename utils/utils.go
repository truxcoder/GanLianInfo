package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

func GetAgeFromIdCode(code string) int {
	reg := regexp.MustCompile(`^[1-9]\d{5}(18|19|20)(\d{2})((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`)
	//reg := regexp.MustCompile(`^[1-9]\d{5}(18|19|20)`)
	params := reg.FindStringSubmatch(code)
	birYear, _ := strconv.Atoi(params[1] + params[2])
	birMonth, _ := strconv.Atoi(params[3])
	age := time.Now().Year() - birYear
	if int(time.Now().Month()) < birMonth {
		age--
	}
	return age
}

func GetAgeFromBirthday(birthday time.Time) int {
	year := birthday.Year()
	month := int(birthday.Month())
	day := birthday.Day()
	age := time.Now().Year() - year
	if int(time.Now().Month()) < month || (int(time.Now().Month()) == month && time.Now().Day() < day) {
		age--
	}
	return age
}

func GetBirthdayFromIdCode(code string) time.Time {
	//code = strings.TrimSpace(code)
	reg := regexp.MustCompile(`^[1-9]\d{5}(18|19|20)(\d{2})((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`)
	//reg := regexp.MustCompile(`^[1-9]\d{5}(18|19|20)`)
	params := reg.FindStringSubmatch(code)
	if len(params) < 7 {
		return time.Time{}
	}
	year, err := strconv.Atoi(params[1] + params[2])
	if err != nil {
		fmt.Println(err)
		return time.Time{}
	}
	month, err := strconv.Atoi(params[3])
	if err != nil {
		fmt.Println(err)
		return time.Time{}
	}
	day, err := strconv.Atoi(params[6])
	if err != nil {
		fmt.Println(err)
		return time.Time{}
	}
	birthday := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return birthday
}

func GetGenderFromIdCode(code string) string {
	if code == "" {
		return ""
	}
	if len(code) != 18 {
		return ""
	}
	s := code[16:17]
	g, err := strconv.Atoi(s)
	if err != nil {
		return ""
	}
	if g%2 == 0 {
		return "女"
	}
	return "男"
}

func StructToMap(in interface{}) (map[string]interface{}, error) {
	// 当前函数只接收struct类型
	inV := reflect.Indirect(reflect.ValueOf(in))
	inT := reflect.Indirect(reflect.ValueOf(in)).Type()
	//v := reflect.ValueOf(in)
	//if v.Kind() == reflect.Ptr { // 结构体指针
	//	v = v.Elem()
	//}
	if inT.Kind() != reflect.Struct {
		return nil, fmt.Errorf("StructToMap函数的参数只能为struct指针; got %+v", inT)
	}

	out := make(map[string]interface{})
	for i := 0; i < inT.NumField(); i++ {
		p := inT.Field(i)
		if !p.Anonymous {
			out[p.Name] = inV.Field(i).Interface()
		} else {
			field := inV.Field(i)
			for j := 0; j < p.Type.NumField(); j++ {
				pp := p.Type.Field(j)
				out[pp.Name] = field.Field(j).Interface()
			}
		}
	}
	return out, nil
}
