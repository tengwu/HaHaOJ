package model

type PrivilegeType = uint64 //最多细分64个权限

const (
	CAN_CREATE = uint64(1) << iota
	CAN_DELETE
	CAN_UPDATE
	CAN_SCAN_
)
