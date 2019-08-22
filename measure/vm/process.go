package vm

import (
	"github.com/jinzhu/gorm"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"os/user"
	"sort"
	"strconv"
	"time"
)

//打印进程信息,只存放最消耗CPU的TOP10的进程信息
type ProcStat struct {
	Pid int32
	Background bool
	CPUPercent float64  //按照此指标排序只存放最消耗CPU的TOP10进程信息,此处为当前进程占用的CPU时间的比例，而非进程从启动到现在cpu时间占比
	CmdLine string
	CreateTime time.Time
	Cwd string
	User string
	process.MemoryInfoStat
	TotalCPUTimes cpu.TimesStat 							`gorm:"-"` //当前操作系统CPU总共时间用于计算CPU占比
	ProcCPUTimes  cpu.TimesStat								`gorm:"-"` //当前进程一共使用时间
	IOCounterStat process.IOCountersStat  	`gorm:"-"`
	NetIO []net.IOCountersStat     			`gorm:"-"`//记录集合在一起的网卡流量信息
	ReadBytes_S uint64  	//每秒读取字节数
	WriteBytes_S uint64 	//每秒写入字节数
	InterfaceName string 	//interface name
	BytesSent_S uint64   	//每秒发送的字节数
	BytesRecv_S uint64  	//每秒接受的字节数
	InsertTime 		time.Time		`grom:"column:insert_time    ;		type:datetime;		index"`
}

func (ProcStat)TableName()string  {
	return "vm_procstat_top10"
}

type PidSort struct {
	Pid int32
	SortKey float64
}
type PidSortList []*PidSort

func (p PidSortList)Len() int  {
	return len(p)
}
func (p PidSortList)Less(i,j int) bool  {
	return p[i].SortKey > p[j].SortKey
}
func (p PidSortList)Swap(i,j int)  {
	p[i],p[j] = p[j],p[i]
}




func geProcStat(interval time.Duration,ch chan []*ProcStat) {
	//对于本次或者上次没出现的进程不进行IO和net指标采集
	getTopPids := func() ([]int32){
		procs,_ := process.Processes()
		var sortList = make([]*PidSort,0)
		for _,p := range procs {
			var (
				pid int32
				sortkey float64
				err error
			)
			pid = p.Pid
			if sortkey,err = p.CPUPercent();err != nil {
				sortkey = 0
			}
			sortList = append(sortList, &PidSort{
				Pid:     pid,
				SortKey: sortkey,
			})
		}
		//对sortList进行排序
		sort.Sort(PidSortList(sortList))
		//保留前10
		var pidList []int32
		var totalPids = len(sortList)
		var topPids = 10
		if totalPids < topPids {
			topPids = totalPids
		}
		for _,eachSortPid := range sortList[:topPids] {
			pidList = append(pidList, eachSortPid.Pid)
		}
		return pidList
	}
	//获取proc相关的快照信息,包括进程信息、进程的IO信息、进程的网络信息,并且与上次获取的ProcStat信息进行对比
	getTopProc := func(pids []int32,lastProcStatList []*ProcStat,insertTime time.Time) (ProcStatList []*ProcStat) {
		for _,eachPid := range pids {
			p,err := process.NewProcess(eachPid)
			if err != nil {
				//说明进程不存在,取下一个
				continue
			}
			var self = new (ProcStat)
			timeStat,_ := cpu.Times(false)
			procCPUTimes,_ := p.Times()
			ioconter,err := p.IOCounters()
			if err != nil {
				ioconter = &process.IOCountersStat{}
			}
			netIOCounterStatList,err := p.NetIOCounters(false)
			if err != nil {
				netIOCounterStatList = []net.IOCountersStat{net.IOCountersStat{}}
			}
			memInfoStat,err := p.MemoryInfo()
			if err != nil {
				memInfoStat = &process.MemoryInfoStat{}
			}
			createTime,_ := p.CreateTime() //从什么时间开始，毫秒，需要换算成当前时间
			self.Pid = p.Pid
			self.Background,_ = p.Background()
			self.CmdLine,_ = p.Cmdline()
			self.CreateTime = time.Unix(0,createTime * 1000000)
			self.Cwd,_ = p.Cwd()
			uids,_ := p.Uids()
			if user,err := user.LookupId(strconv.Itoa(int(uids[0])));err != nil {
				self.User = ""
			}else {
				self.User = user.Name
			}
			self.TotalCPUTimes = timeStat[0]
			self.ProcCPUTimes = *procCPUTimes
			self.MemoryInfoStat = *memInfoStat
			self.IOCounterStat = *ioconter
			self.NetIO = netIOCounterStatList
			self.InsertTime = insertTime
			//检查在lastProcStatList中是否存在相同的PID信息，如果存在求每秒IO、net相关指标信息
			for _,eachLastProcStat := range lastProcStatList {
				if self.Pid == eachLastProcStat.Pid {
					//计算CPU使用率 总的cpu时间totalCpuTime = user + nice + system + idle + iowait + irq + softirq + stealstolen  +  guest
					vm_cpus_all_time := self.TotalCPUTimes.User + self.TotalCPUTimes.Nice + self.TotalCPUTimes.System + self.TotalCPUTimes.Idle + self.TotalCPUTimes.Iowait + self.TotalCPUTimes.Softirq + self.TotalCPUTimes.Steal + self.TotalCPUTimes.Guest
					proc_cpu_all_time := self.ProcCPUTimes.User + self.ProcCPUTimes.Nice + self.ProcCPUTimes.System + self.ProcCPUTimes.Idle + self.ProcCPUTimes.Iowait + self.ProcCPUTimes.Softirq + self.ProcCPUTimes.Steal + self.ProcCPUTimes.Guest
					last_vm_cpus_all_time := eachLastProcStat.TotalCPUTimes.User + eachLastProcStat.TotalCPUTimes.Nice + eachLastProcStat.TotalCPUTimes.System + eachLastProcStat.TotalCPUTimes.Idle + eachLastProcStat.TotalCPUTimes.Iowait + eachLastProcStat.TotalCPUTimes.Softirq + eachLastProcStat.TotalCPUTimes.Steal + eachLastProcStat.TotalCPUTimes.Guest
					last_proc_cpu_all_time := eachLastProcStat.ProcCPUTimes.User + eachLastProcStat.ProcCPUTimes.Nice + eachLastProcStat.ProcCPUTimes.System + eachLastProcStat.ProcCPUTimes.Idle + eachLastProcStat.ProcCPUTimes.Iowait + eachLastProcStat.ProcCPUTimes.Softirq + eachLastProcStat.ProcCPUTimes.Steal + eachLastProcStat.ProcCPUTimes.Guest
					delta_total_cpu_time := (vm_cpus_all_time-last_vm_cpus_all_time)
					if delta_total_cpu_time == 0{
						self.CPUPercent = 0
					}else {
						self.CPUPercent = (proc_cpu_all_time - last_proc_cpu_all_time)/delta_total_cpu_time
					}

					delta_time := uint64(self.InsertTime.Sub(eachLastProcStat.InsertTime).Seconds())
					if delta_time == 0 {
						delta_time = 1
					}
					self.ReadBytes_S = (self.IOCounterStat.ReadBytes - eachLastProcStat.IOCounterStat.ReadBytes)/delta_time
					self.WriteBytes_S = (self.IOCounterStat.WriteBytes - eachLastProcStat.IOCounterStat.WriteBytes)/delta_time
					//处理网卡流量信息
					if len(self.NetIO) == 0 || len(eachLastProcStat.NetIO) == 0 {
						//不存在网卡流量信息，
					}else {
						//因为网卡已经聚合，只处理第一条即可
						netIO := self.NetIO[0]
						lastNetIO := eachLastProcStat.NetIO[0]
						self.InterfaceName = netIO.Name
						self.BytesSent_S = (netIO.BytesSent - lastNetIO.BytesSent)/delta_time
						self.BytesRecv_S = (netIO.BytesRecv - lastNetIO.BytesRecv)/delta_time
					}
				}
			}
			ProcStatList = append(ProcStatList, self)
		}
		return ProcStatList
	}
	ticker := time.NewTicker(interval)
	var lastProcStat []*ProcStat
	for {
		select {
		case <-ticker.C:
			//开始写入数据
			procStat := getTopProc(getTopPids(),lastProcStat,time.Now())
			ch <- procStat
			lastProcStat = procStat
		}
	}
	ticker.Stop()
}

//将采集的数据放入数据库中
func SaveProcStat(db *gorm.DB,interval time.Duration)  {
	var ch chan []*ProcStat
	ch = make(chan []*ProcStat,10)
	go geProcStat(interval,ch)
	//对于ch出来的数据插入到数据库中
	if ! db.HasTable(&ProcStat{}) {
		db.CreateTable(&ProcStat{})
	}
	for eachProcStat := range ch {
		//针对每一行数据都插入到数据库中
		tx := db.Begin()
		for _,eachP := range eachProcStat {
			tx.Create(eachP)
		}
		tx.Commit()
	}
}