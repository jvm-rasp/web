package socket

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"server/model"
	"server/repository"
	"server/socket/message"
	"strconv"
	"strings"
)

// 日志topic 或者命令
const DAEMON_TOPIC byte = 0x01
const AGENT_TOPIC byte = 0x02
const MODULE_TOPIC byte = 0x03
const ATTACK_TOPIC byte = 0x04

// LOG_PACKAGE_CONSTANT_LENGTH 消息体的固定长度 149
const LOG_PACKAGE_CONSTANT_LENGTH = 149

var raspHostRepository = repository.NewRaspHostRepository()
var javaProcessRepository = repository.NewJavaProcessInfoRepository()

func NewSockekServer(port int) error {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%v", port))
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	defer tcpListener.Close()
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("get a connection")
		go handle(tcpConn)
	}
	return nil
}

func handle(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	// 长度分割
	scanner.Split(splitLog)
	for scanner.Scan() {
		scannedPack := new(LogPackage)
		scannedPack.Unpack(bytes.NewReader(scanner.Bytes()))
		Handler(scannedPack)
	}
}

// 处理日志
func Handler(p *LogPackage) {
	switch p.Type {
	case DAEMON_TOPIC:
		handleDaemonLog(p)
		break
	default:

	}
}

// splitLog 将字节码流分割为消息体
func splitLog(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// TODO 疑似npe
	if !atEOF && len(data) >= 4 && string(data[0:3]) == string(magicBytes[:]) {
		if len(data) > 13 {
			hostNamelength := int32(0)
			err := binary.Read(bytes.NewReader(data[5:9]), binary.BigEndian, &hostNamelength)
			if err != nil {
				return 0, nil, err
			}
			bodySize := int32(0)
			err = binary.Read(bytes.NewReader(data[9:13]), binary.BigEndian, &bodySize)
			if err != nil {
				return 0, nil, err
			}
			length := int(hostNamelength) + int(bodySize)
			if length+LOG_PACKAGE_CONSTANT_LENGTH <= len(data) {
				return LOG_PACKAGE_CONSTANT_LENGTH + length, data[:LOG_PACKAGE_CONSTANT_LENGTH+length], nil
			}
		}
	}
	return
}

func handleDaemonLog(p *LogPackage) {
	daemonMessage := &message.DaemonMessage{}
	if p != nil {
		err := json.Unmarshal(p.Body, &daemonMessage)
		if err != nil {
			// TODO
			return
		}
	}

	// 不同logid 处理
	switch daemonMessage.Logid {
	case message.DAEMON_STARTUP_LOGID:
		handleStartupLog(daemonMessage)
	case message.HOST_ENV_LOGID:
		handleHostEnvLog(daemonMessage)
	case message.HEART_BEAT_LOGID:
		handleHeartbeatLog(daemonMessage)
	case message.JAVA_PROCESS_STARTUP:
		handleFindJavaProcessLog(daemonMessage)
	case message.JAVA_PROCESS_SHUTDOWN:
		handleRemoveJavaProcessLog(daemonMessage)
	}
}

func handleStartupLog(message *message.DaemonMessage) {
	host := &model.RaspHost{
		HostName:     message.HostName,
		Ip:           message.Ip,
		HeatbeatTime: message.Ts,
	}
	detailMap := make(map[string]string)
	err := json.Unmarshal([]byte(message.Detail), &detailMap)
	if err != nil {
		// TODO
		return
	}
	host.AgentMode = detailMap["agentMode"]
	dbData, err := raspHostRepository.QueryRaspHost(host.HostName)
	if err != nil {
		// todo
		return
	}
	if len(dbData) <= 0 {
		err = raspHostRepository.CreateRaspHost(host)
		if err != nil {
			// todo
			return
		}
	}
	dbData, err = raspHostRepository.QueryRaspHost(host.HostName)
	if err != nil {
		// todo
		return
	}
	fmt.Println("host.AgentMode")
	fmt.Println(host.AgentMode)
}

func handleHostEnvLog(message *message.DaemonMessage) {
	// 主机相关的信息，安装之后不发生变化，相对固定
	detailMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(message.Detail), &detailMap)
	if err != nil {
		// TODO
		return
	}

	installDir := detailMap["installDir"].(string)
	version := detailMap["version"].(string)
	binFileHash := detailMap["binFileHash"].(string)
	osType := detailMap["osType"].(string)

	totalMem := detailMap["totalMem"].(float64)
	cpuCounts := detailMap["cpuCounts"].(float64)
	freeDisk := detailMap["freeDisk"].(float64)

	buildDateTime := detailMap["buildDateTime"].(string)
	buildGitBranch := detailMap["buildGitBranch"].(string)
	buildGitCommit := detailMap["buildGitCommit"].(string)

	dbData, err := raspHostRepository.QueryRaspHost(message.HostName)
	if err != nil {
		// todo
		return
	}

	var host *model.RaspHost
	if len(dbData) == 0 {
		host = &model.RaspHost{
			HostName: message.HostName,
			Ip:       message.Ip,
		}
	} else {
		host = dbData[0]
	}

	host.InstallDir = installDir
	host.Version = version
	host.ExeFileHash = binFileHash
	host.InstallDir = osType
	host.TotalMem = totalMem
	host.CpuCounts = cpuCounts
	host.FreeDisk = freeDisk
	host.BuildDateTime = buildDateTime
	host.BuildGitBranch = buildGitBranch
	host.BuildGitBranch = buildGitCommit
	host.HeatbeatTime = message.Ts

	if len(dbData) == 0 {
		err = raspHostRepository.CreateRaspHost(host)
		if err != nil {
			// todo
			return
		}
	} else {
		err := raspHostRepository.UpdateRaspHost(host)
		if err != nil {
			return
		}
	}
}

func handleHeartbeatLog(message *message.DaemonMessage) {
	dbData, err := raspHostRepository.QueryRaspHost(message.HostName)
	if err != nil {
		// todo
		return
	}

	var host *model.RaspHost
	if len(dbData) == 0 {
		host = &model.RaspHost{
			HostName: message.HostName,
			Ip:       message.Ip,
		}
	} else {
		host = dbData[0]
	}

	host.HeatbeatTime = message.Ts

	if len(dbData) == 0 {
		err = raspHostRepository.CreateRaspHost(host)
		if err != nil {
			// todo
			return
		}
	} else {
		err := raspHostRepository.UpdateRaspHost(host)
		if err != nil {
			return
		}
	}

	list, err := javaProcessRepository.GetAllJavaProcessInfos(message.HostName)
	if err != nil {
		return
	}

	// 删除process信息
	detailMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(message.Detail), &detailMap)
	if err != nil {
		// TODO
		return
	}

	for _, v := range list {
		if len(detailMap) == 0 || detailMap[string(v.Pid)] == nil {
			err := javaProcessRepository.DeleteProcess(v.ID)
			if err != nil {
				return
			}
		}
	}
}

func handleFindJavaProcessLog(message *message.DaemonMessage) {
	detailMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(message.Detail), &detailMap)
	if err != nil {
		panic(err)
		return
	}

	pid := detailMap["javaPid"].(float64)
	fmt.Printf("pid:%v", pid)
	startTime := detailMap["startTime"].(string)
	status := detailMap["injectedStatus"].(string)

	cmdLines := detailMap["cmdLines"].([]interface{})
	var paramSlice []string
	for _, param := range cmdLines {
		switch v := param.(type) {
		case string:
			paramSlice = append(paramSlice, v)
		case int:
			strV := strconv.FormatInt(int64(v), 10)
			paramSlice = append(paramSlice, strV)
		default:
			panic("params type not supported")
		}
	}

	process := &model.JavaProcessInfo{
		Status:      status,
		Pid:         int(pid),
		StartTime:   startTime,
		CmdlineInfo: strings.Join(paramSlice, ","),
		HostName:    message.HostName,
	}
	fmt.Printf("SaveProcessInfo:%v\n", process.Pid)
	err = javaProcessRepository.SaveProcessInfo(process)
	if err != nil {
		panic(err)
	}
}

func handleRemoveJavaProcessLog(message *message.DaemonMessage) {
	pid, err := strconv.ParseInt(message.Detail, 10, 32)
	if err != nil {
		// TODO
		return
	}
	err = javaProcessRepository.DeleteProcessByPid(message.HostName, uint(pid))
	if err != nil {
		// TODO
		return
	}
}
