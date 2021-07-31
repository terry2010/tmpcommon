package Common

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	//if err != nil {
	//	return "", err
	//}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}

	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}

	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0 : i+1]), nil
}

func GetServerIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println(err)
		return ""
	}
	for _, value := range addrs {
		if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.IsGlobalUnicast() && !ipnet.IP.IsInterfaceLocalMulticast() {
			if ipnet.IP.To4() != nil {
				//log.Println(ipnet.IP.String())
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func Page404(c *gin.Context) {
	//返回404状态码
	c.JSON(http.StatusNotFound, gin.H{
		"code": http.StatusNotFound,
		"msg":  "404, page not exists!",
		"data": "",
	})
}

/**
取得所有已注册slave的列表
*/
func GetActiveSlaveList() (slaveList []ServerInfo, err error) {
	list, err := RedisClient.HGetAll(Config.GetString("key.RedisSlaveServerKey")).Result()

	if nil == err {
		for _, v := range list {

			var slaveInfo ServerInfo
			err := Json.UnmarshalFromString(v, &slaveInfo)

			if nil == err {
				if true == slaveInfo.Status && !(slaveInfo.HostIP == Config.GetString("server.hostIP") && slaveInfo.HostPort == Config.GetString("server.hostPort")) {

					ok, err := slaveInfo.CheckServerHealth()
					if true == ok {
						slaveList = append(slaveList, slaveInfo)
					} else {
						log.Println("GetActiveSlaveList:err:", err)

					}
				}
			} else {
				//跳过异常信息
				log.Println("GetActiveSlaveList:err2:", err)

			}

		}
	} else {
		log.Println("GetActiveSlaveList:err3:", err)
		return slaveList, err
	}


	return slaveList, nil
}

func FastAtoi(num string) int {
	ret , _ := strconv.Atoi(num)
	return ret
}

