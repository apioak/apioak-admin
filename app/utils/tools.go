package utils

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	sequenceBits = uint64(12)
	maxSequence = int64(-1) ^ (int64(-1) << sequenceBits)
	timeLeft = uint8(22)   // 时间戳向左偏移量
	dataLeft = uint8(17)  // 数据中心ID向左偏移量
	workLeft = uint8(12)  // 节点ID向左偏移量
	twepoch = int64(1577808000000) // 2020-01-01 00:00:00 初始时间 常量时间戳(毫秒)

	IdTypeUser = "u"
	IdTypeService = "svc"
	IdTypeRoute = "rt"
	IdTypePlugin = "pu"
	IdTypeCertificate = "cert"
	IdTypeNode = "nd"
)

type Worker struct {
	mu sync.Mutex
	LastStamp int64 // 记录上一次ID的时间戳
	WorkerID int64 // 该节点的ID 分布式情况下,需通过外部配置文件或其他方式为每台机器分配独立的id
	DataCenterID int64 // // 该节点的 数据中心ID
	Sequence int64 // 当前毫秒已经生成的ID序列号(从0 开始累加) 1毫秒内最多生成4096个ID
}

func newWorker() *Worker  {
	return &Worker{
		WorkerID: 1,
		LastStamp: 0,
		Sequence: 0,
	}
}

// 雪花算法生成ID
func snowFlake() (int64, error) {
	worker := newWorker()

	worker.mu.Lock()
	defer worker.mu.Unlock()

	timeStamp := time.Now().UnixNano() / 1e6
	if timeStamp < worker.LastStamp {
		return 0, fmt.Errorf("Time goes back")
	}

	if worker.LastStamp == timeStamp {
		worker.Sequence = (worker.Sequence + 1) & maxSequence
		if worker.Sequence == 0 {
			for timeStamp <= worker.LastStamp {
				timeStamp = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		worker.Sequence = 0
	}

	worker.LastStamp = timeStamp
	id := ((timeStamp - twepoch) << timeLeft) | (worker.DataCenterID << dataLeft)  | (worker.WorkerID << workLeft) | worker.Sequence

	return id, nil
}

// 自动生成ID
func IdGenerate(idType string) (string, error) {

	var id string

	// 获取ID
	snowFlakeId, err := snowFlake()

	// 出错则直接返回错误信息
	if  err != nil {
		return "", err
	} else {
		switch strings.ToLower(idType) {
		case IdTypeUser:
			id = IdTypeUser + "-" + strconv.Itoa(int(snowFlakeId))
		case IdTypeService:
			id = IdTypeService + "-" + strconv.Itoa(int(snowFlakeId))
		case IdTypeRoute:
			id = IdTypeRoute + "-" + strconv.Itoa(int(snowFlakeId))
		case IdTypePlugin:
			id = IdTypePlugin + "-" + strconv.Itoa(int(snowFlakeId))
		case IdTypeCertificate:
			id = IdTypeCertificate + "-" + strconv.Itoa(int(snowFlakeId))
		case IdTypeNode:
			id = IdTypeNode + "-" + strconv.Itoa(int(snowFlakeId))
		default:
			return "", fmt.Errorf("id type error")
		}
	}

	return id, nil
}