package dao

import (
	"HappyOnlineJudge/common"
	"strconv"
	"time"
)

type UserDao struct {
	ID       int64
	Username string
	User     *User
}

func (ud *UserDao) GetRedisExpire() time.Duration {
	return USER_REDIS_EXPIRE
}
func (ud *UserDao) GetTableName() string {
	return "user"
}
func (ud *UserDao) GetSelf() interface{} {
	if ud.User == nil {
		ud.User = &User{}
	}
	return ud.User
}
func (ud *UserDao) GetRelatedKeysMap() []interface{} {
	return []interface{}{ud.GetNameKey(), ud.GetID()}
}
func (ud *UserDao) GetName() string {
	if ud.Username == "" {
		if ud.User != nil && ud.User.Username != "" {
			ud.Username = ud.User.Username
		} else {
			ud.Username = OneCol(ud, "username").ToString()
		}
	}
	return ud.Username
}
func (ud *UserDao) GetNameKey() string {
	return USER_REDIS_PREFIX + ud.GetName()
}
func (ud *UserDao) GetRedisKey() string { //必须有id
	return USER_REDIS_PREFIX + strconv.FormatInt(ud.GetID(), 10)
}
func (ud *UserDao) GetID() int64 {
	if ud.ID == 0 {
		if ud.User != nil && ud.User.ID != 0 {
			ud.ID = ud.User.ID
		} else if ud.Username != "" || (ud.User != nil && ud.User.Username != "") {
			key := ud.GetNameKey()
			if rdb.Exists(ctx, key).Val() > 0 {
				id := rdb.Get(ctx, key).Val()
				ud.ID = common.StrToInt64(id)
			} else {
				x := new(Col)
				if ok, err := engine.SQL("select id from user where username = ?", ud.Username).Get(&x.data); err == nil && ok {
					ud.ID = x.ToInt64()
				}
			}
		}
	}
	return ud.ID
}

func (ud *UserDao) Update(mp map[string]interface{}) error {
	if err := UpdateCols(ud, mp); err != nil {
		return err
	}
	if newname, ok := mp["username"]; ok {
		rdb.Del(ctx, ud.GetNameKey())
		ud.Username = newname.(string)
		rdb.Set(ctx, ud.GetNameKey(), ud.GetID(), USER_REDIS_EXPIRE)
	}
	return nil
}

type UsersData struct {
	IDs   []int64
	Datas [][]Col
}

func GetTableName() string {
	return "user"
}
func (us *UsersData) GetIDs(cols []string, values []interface{}, a ...int) []int64 { //len(a)=0或2
	if len(a) == 0 {
		engine.Table("user").Where(ToSqlConditions(cols), values...).Cols("id").Find(&us.IDs)
	} else {
		engine.Table("user").Where(ToSqlConditions(cols), values...).Cols("id").Limit(a[0], a[1]).Find(&us.IDs)
	}
	return us.IDs
}
func (us *UsersData) GetColsByIDs(wants []string) [][]Col {
	for _, id := range us.IDs {
		us.Datas = append(us.Datas, Cols(&UserDao{ID: id}, wants...))
	}
	return us.Datas
}
