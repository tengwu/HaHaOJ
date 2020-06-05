package dao

import (
	"errors"
	"time"
)

//单个对象的增删改查
type SingleData interface { //对一个数据的增删改查, 注释掉的方法是应该需要实现的, 未注释的实现了可以使用泛型
	GetRedisKey() string  //获取存到redis的key
	GetID() int64         //获取对象数据库id
	GetTableName() string //获取数据库表名
	GetRedisExpire() time.Duration
	GetRelatedKeysMap() []interface{}
	GetSelf() interface{}
}

//下面是泛型方法

//获取相关redis keys
func GetRelatedKeys(sd SingleData) []string {
	keys := []string{sd.GetRedisKey()}
	if rk := sd.GetRelatedKeysMap(); rk != nil {
		num := len(rk)
		for i := 0; i < num; i += 2 {
			keys = append(keys, rk[i].(string))
		}
	}
	return keys
}

//将自己上传到redis中
func PutToRedis(sd SingleData) error {
	self := sd.GetSelf()
	if err := putObjToRedis(sd.GetRedisKey(), self, sd.GetRedisExpire()); err != nil {
		return err
	}
	if rk := sd.GetRelatedKeysMap(); rk != nil {
		if err := rdb.MSet(ctx, rk...).Err(); err != nil {
			return err
		}
	}
	return nil
}

func IsInRedis(sd SingleData) bool {
	return rdb.Exists(ctx, sd.GetRedisKey()).Val() > 0
}

//如果已经在redis,更新过期时间,否则加入
func PutToRedisIfNotIn(sd SingleData) error {
	key := sd.GetRedisKey()
	if rdb.Exists(ctx, key).Val() > 0 {
		rdb.Expire(ctx, key, sd.GetRedisExpire())
	} else {
		return putObjToRedis(key, sd.GetSelf(), sd.GetRedisExpire())
	}
	return nil
}

func DeleteFromRedis(sd SingleData) error {
	rdb.Del(ctx, GetRelatedKeys(sd)...)
	return nil
}

//创建
func Create(sd SingleData) error {
	self := sd.GetSelf()
	if num, err := engine.InsertOne(self); err != nil || num != 1 {
		return err
	}
	return PutToRedis(sd)
}

//删除
func Delete(sd SingleData) error {
	id := sd.GetID()
	if id == 0 {
		return errors.New("不存在该用户的ID")
	}
	if err := DeleteFromRedis(sd); err != nil {
		return err
	}
	sql := "delete from " + sd.GetTableName() + " where id=?"
	if _, err := engine.Exec(sql, id); err != nil {
		return err
	}
	return nil
}

func GetColsOfSelf(sd SingleData, wants ...string) interface{} {
	self := sd.GetSelf()
	if _, err := engine.ID(sd.GetID()).Cols(wants...).Get(self); err != nil {
		return err
	}
	return self
}
func GetSelfAll(sd SingleData) interface{} {
	key := sd.GetRedisKey()
	self := sd.GetSelf()
	if rdb.Exists(ctx, key).Val() > 0 {
		if err := GetObjFromRedis(key, self); err != nil {
			return nil
		}
	} else {
		if exist, err := engine.ID(sd.GetID()).Get(self); !exist || err != nil {
			return nil
		}
		if err := PutToRedisIfNotIn(sd); err != nil {
			return nil
		}
	}
	return self
}

//获取某一个数据的单字段内容
func OneCol(sd SingleData, want string) *Col {
	id := sd.GetID()
	x := new(Col)
	key := sd.GetRedisKey()
	if rdb.Exists(ctx, key).Val() > 0 {
		x.data = rdb.HGet(ctx, key, want).Val()
	} else {
		sql := "select " + want + " from " + sd.GetTableName() + " where id = ?"
		if _, err := engine.SQL(sql, id).Get(&x.data); err != nil {
			return nil
		}
	}
	return x
}

//获取多列
func Cols(sd SingleData, wants ...string) []Col {
	id := sd.GetID()
	n := len(wants)
	ret := make([]Col, n)
	x := make([]interface{}, n)
	key := sd.GetRedisKey()
	if rdb.Exists(ctx, key).Val() > 0 {
		x = rdb.HMGet(ctx, key, wants...).Val()
	} else {
		sql := ToSqlSelect(wants...) + " from " + sd.GetTableName() + " where id = ?"
		if _, err := engine.SQL(sql, id).Get(&x); err != nil {
			return nil
		}
	}
	for i := 0; i < n; i++ {
		ret[i].data = x[i]
	}
	return ret
}

//map更新某些列, 不要改变作为rediskey的相关列
func UpdateCols(sd SingleData, mp map[string]interface{}) error {
	args := make([]interface{}, 0)
	sql := "update " + sd.GetTableName() + " set "
	first := true
	for k, v := range mp {
		t := typeAnalyzed(v)
		args = append(args, k, t)
		if first {
			sql += k + "=?"
			first = false
		} else {
			sql += "," + k + "=?"
		}
	}
	sql += " where id=?"
	n := len(mp)
	sqlArgs := make([]interface{}, 2+n)
	sqlArgs[0] = sql
	for i := 0; i < n; i++ {
		sqlArgs[i+1] = args[2*i+1]
	}
	sqlArgs[n+1] = sd.GetID()
	if _, err := engine.Exec(sqlArgs...); err != nil {
		return err
	}
	if key := sd.GetRedisKey(); rdb.Exists(ctx, key).Val() > 0 {
		if _, err := rdb.HMSet(ctx, key, args...).Result(); err != nil {
			return err
		}
		rdb.Expire(ctx, key, sd.GetRedisExpire())
	}
	return nil
}

//多条数据的查询
type ManyData interface {
	GetIDs(cols []string, values []interface{}, a ...int) []int64 //获取满足条件的id,条数限制为a[0],起始位置为a[1], a为空时不限制
	//GetColsByIDs(wants ...string) [][]Col
	//GetTableName() string
}

//多态
func Count(md ManyData, cols []string, values []interface{}, a ...int) int {
	return len(md.GetIDs(cols, values, a...))
}
