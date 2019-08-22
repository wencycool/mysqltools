package vm

import (
	"github.com/jinzhu/gorm"
	"github.com/shirou/gopsutil/disk"
	"time"
)

//usage
type DiskUsageStat struct {
	DeviceId				uint64
	Device 					string
	MountPoint 				string
	Fstype 					string
	Opts 					string
	Total 					uint64
	Free 					uint64
	Used 					uint64
	UsedPercent 			float64
	InodesTotal 			uint64
	InodesUsed 				uint64
	InodesFree 				uint64
	InodesUsedPercent 		float64
	InsertTime 		time.Time		`grom:"column:insert_time    ;		type:datetime;		index"`
}

func (DiskUsageStat)TableName()string  {
	return "vm_diskstat"
}


func getDiskUsageStat(interval time.Duration,ch chan []*DiskUsageStat)  {
	ticker := time.NewTicker(interval)
	f := func() (dstats []*DiskUsageStat) {
		partitions,_ := disk.Partitions(true)
		t_now := time.Now()
		for id,eachP := range partitions {
			dusage,_ := disk.Usage(eachP.Mountpoint)
			dstats = append(dstats, &DiskUsageStat{
				DeviceId:		   uint64(id),
				Device:            eachP.Device,
				MountPoint:        eachP.Mountpoint,
				Fstype:            eachP.Fstype,
				Opts:              eachP.Opts,
				Total:             dusage.Total,
				Free:              dusage.Free,
				Used:              dusage.Used,
				UsedPercent:       dusage.UsedPercent,
				InodesTotal:       dusage.InodesTotal,
				InodesUsed:        dusage.InodesUsed,
				InodesFree:        dusage.InodesFree,
				InodesUsedPercent: dusage.InodesUsedPercent,
				InsertTime:        t_now,
			})
		}
		return dstats
	}
	for {
		select {
		case <-ticker.C:
			//开始写入数据
			ch <- f()
		}
	}
	ticker.Stop()
}

//将采集的数据放入数据库中
func SaveDiskUsageStat(db *gorm.DB,interval time.Duration)  {
	var ch chan []*DiskUsageStat
	ch = make(chan []*DiskUsageStat,10)
	go getDiskUsageStat(interval,ch)
	//对于ch出来的数据插入到数据库中
	if ! db.HasTable(&DiskUsageStat{}) {
		db.CreateTable(&DiskUsageStat{})
	}
	for eachDiskUsageStats := range ch {
		tx := db.Begin()
		for _,eachDiskUsageStat := range eachDiskUsageStats {
			tx.Create(eachDiskUsageStat)
		}
		tx.Commit()
	}
}