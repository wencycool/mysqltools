package vm

import (
	"github.com/jinzhu/gorm"
	"github.com/shirou/gopsutil/mem"
	"time"
)

//处理物理内存和虚拟内存的过程

type SwapMemoryStat struct {
	*mem.SwapMemoryStat
	InsertTime 		time.Time		`grom:"column:insert_time    ;		type:datetime;		index"`
}

func (SwapMemoryStat) TableName()string {
	return "vm_swapstat"
}


func getSwapMemoryStat(interval time.Duration,ch chan *SwapMemoryStat)  {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			//开始写入数据
			swapMemoryStat,_ := mem.SwapMemory()
			ch <- &SwapMemoryStat{swapMemoryStat,time.Now()}
		}
	}
	ticker.Stop()
}

//将采集的数据放入数据库中
func SaveSwapMemoryStat(db *gorm.DB,interval time.Duration)  {
	var ch chan *SwapMemoryStat
	ch = make(chan *SwapMemoryStat,10)
	go getSwapMemoryStat(interval,ch)
	//对于ch出来的数据插入到数据库中
	if ! db.HasTable(&SwapMemoryStat{}) {
		db.CreateTable(&SwapMemoryStat{})
	}
	for eachSwapMemoryStat := range ch {
		db.Create(eachSwapMemoryStat)
	}
}


type VirtualMemoryStat struct {
	*mem.VirtualMemoryStat
	InsertTime 		time.Time		`grom:"column:insert_time    ;		type:datetime;		index"`
}

func (VirtualMemoryStat) TableName()string {
	return "vm_memstat"
}

func getVirtualMemoryStat(interval time.Duration,ch chan *VirtualMemoryStat)  {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			//开始写入数据
			virMemStat,_ := mem.VirtualMemory()
			ch <- &VirtualMemoryStat{virMemStat,time.Now()}
		}
	}
	ticker.Stop()
}

//将采集的数据放入数据库中
func SaveVirtualMemoryStat(db *gorm.DB,interval time.Duration)  {
	var ch chan *VirtualMemoryStat
	ch = make(chan *VirtualMemoryStat,10)
	go getVirtualMemoryStat(interval,ch)
	//对于ch出来的数据插入到数据库中
	if ! db.HasTable(&VirtualMemoryStat{}) {
		db.CreateTable(&VirtualMemoryStat{})
	}
	for eachVirtualMemoryStat := range ch {
		db.Create(eachVirtualMemoryStat)
	}
}