package vm

import (
	"github.com/jinzhu/gorm"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"strings"
	"time"
)

//当前主机cpu配置信息
type CPUInfo struct {
	CPU        int32    	`gorm:"column:cpu;type:int"`
	VendorID   string   	`gorm:"column:vendorid;type:varchar(100)"`
	Family     string  	 	`gorm:"column:family;type:varchar(100)"`
	Model      string   	`gorm:"column:model;type:varchar(100)"`
	Stepping   int32    	`gorm:"column:stepping;type:int"`
	PhysicalID string   	`gorm:"column:physicalid;type:varchar(30)"`
	CoreID     string   	`gorm:"column:coreid;type:varchar(30)"`
	Cores      int32    	`gorm:"column:cores; type:int"`
	ModelName  string   	`gorm:"column:modelname;type:varchar(100)"`
	Mhz        float64  	`gorm:"column:mhz;type:float"`
	CacheSize  int32    	`gorm:"column:cachesize;type:int"`
	Flags      string 		`gorm:"column:flags;type:varchar(1000)"`
	Microcode  string  		`gorm:"column:microcode;type:varchar(100)"`
	InsertTime 		time.Time		`grom:"column:insert_time    ;		type:datetime;		index"`
}

func (CPUInfo)TableName() string  {
	return "vm_cpuinfo"
}

//获取当前操作系统所有CPU信息
func getCPUInfos() (cpuinfo []*CPUInfo,err error)  {
	stats,err := cpu.Info()
	if err != nil {
		return nil,err
	}
	for _,eachInfo := range stats {
		cpuinfo = append(cpuinfo, &CPUInfo{
			eachInfo.CPU,
			eachInfo.VendorID,
			eachInfo.Family,
			eachInfo.Model,
			eachInfo.Stepping,
			eachInfo.PhysicalID,
			eachInfo.CoreID,
			eachInfo.Cores,
			eachInfo.ModelName,
			eachInfo.Mhz,
			eachInfo.CacheSize,
			strings.Join(eachInfo.Flags,","),
			eachInfo.Microcode,
			time.Now(),
		})
	}
	return cpuinfo,nil
}

//将cpu信息放入到gorm中

func SaveCPUInfo(db *gorm.DB) (error) {
	cpuinfostats,err := getCPUInfos()
	if err != nil {
		return err
	}
	//检查数据库中是否存在该表
	if ! db.HasTable(&CPUInfo{}) {
		db.CreateTable(&CPUInfo{})
	}
	//将数据插入到表中
	tx := db.Begin()
	for _,eachInfo := range cpuinfostats {
		tx.Create(eachInfo)
	}
	tx.Commit()
	return nil
}


//记录CPU的使用率情况,只记录所有CPU资源的平均值
type CPUStat struct {
	TotalUsed 	float64		`gorm:"column:totalused; 	type:float"`	//100% - Idle
	User 		float64		`gorm:"column:user;		 	type:float"` 	//用户使用cpu百分比
	System 		float64		`gorm:"column:system;	 	type:float"`	//系统使用cpu百分比
	IoWait 		float64		`gorm:"column:iowait;	 	type:float"`	//IO等待使用cpu百分比
	Idle 		float64		`gorm:"column:idle;		 	type:float"`	//空闲cpu百分比
	Load1 		float64		`gorm:"column:load1;		type:float"`	//一分钟内平均load数
	Load5 		float64		`gorm:"column:load5;		type:float"`
	Load15		float64		`gorm:"column:load15;		type:float"`
	InsertTime 		time.Time		`grom:"column:insert_time    ;		type:datetime;		index"`
}

func (CPUStat)TableName()string  {
	return "vm_cpustat"
}

//采集cpu使用率信息,t为采集周期,一直采集,采集放入管道中
func getCPUStat(interval time.Duration,ch chan *CPUStat) {
	timestats,_ := cpu.Times(false)
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			//开始写入数据
			timestats_tmp,_ := cpu.Times(false)
			loadstat_tmp, _ := load.Avg()
			userTime := timestats_tmp[0].User - timestats[0].User
			systemTime := timestats_tmp[0].System - timestats[0].System
			iowaitTime := timestats_tmp[0].Iowait - timestats[0].Iowait
			idleTime := timestats_tmp[0].Idle - timestats[0].Idle
			totalTime := userTime + systemTime + iowaitTime + idleTime
			totalUsedTime:= (totalTime - idleTime)
			ch <- &CPUStat{
				TotalUsed: totalUsedTime/totalTime,
				User:      userTime/totalTime,
				System:    systemTime/totalTime,
				IoWait:    iowaitTime/totalTime,
				Idle:      idleTime/totalTime,
				Load1:     loadstat_tmp.Load1,
				Load5:     loadstat_tmp.Load5,
				Load15:    loadstat_tmp.Load15,
				InsertTime: time.Now(),
			}
			timestats = timestats_tmp
		}
	}
	ticker.Stop()
}

//将采集的数据放入数据库中
func SaveCPUStat(db *gorm.DB,interval time.Duration)  {
	var ch chan *CPUStat
	ch = make(chan *CPUStat,10)
	go getCPUStat(interval,ch)
	//对于ch出来的数据插入到数据库中
	if ! db.HasTable(&CPUStat{}) {
		db.CreateTable(&CPUStat{})
	}
	for eachCPUStat := range ch {
		db.Create(eachCPUStat)
	}
}