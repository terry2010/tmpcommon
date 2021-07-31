package Common

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CommonKey struct {
	RedisTaskQueue string
	Version        float64
}

type ServerInfo struct {
	IP           string    `json:"ip"`
	Port         string    `json:"port"`
	HostIP       string    `json:"hostIP"`
	HostPort     string    `json:"hostPort"`   //当在docker内启动， 这个是docker指定的外部端口， 否则和Port值一致
	CreateTime   time.Time `json:"createTime"` //创建时间
	ActiveTime   time.Time `json:"activeTime"`
	TaskLoad     int       `json:"taskLoad"` //任务负载情况， 计算方法为  int( TaskActiveNumber / TaskNumTotal * 100)
	TaskNumTotal int       `json:"taskNumTotal"`
	Status       bool      `json:"status"`
	Err          error     `json:"err"`
}

/*
检查网络的参数是否正常
*/
func (s ServerInfo) CheckServerHealth() (ok bool, err error) {
	url := s.PingURL()
	body, err := HttpGet(time.Duration(2)*time.Second, url)
	if nil == err {
		if strings.Compare(string(body), url) != 0 {

			return false, errors.New("expect:[" + url + "], got:[" + string(body) + "]")
		} else {
			return true, nil
		}
	} else {
		log.Println(err)
		return false, err
	}
}

/**
检查服务是否可用
*/
func (s ServerInfo) CheckServerAvailability() (ok bool, err error) {
	url := s.StatusURL()
	body, err := HttpGet(time.Duration(1)*time.Second, url)
	if nil == err {
		var tmpServerInfo ServerInfo
		err := Json.Unmarshal(body, &tmpServerInfo)
		if nil == err {
			if tmpServerInfo.TaskLoad < 100 {
				return true, nil
			} else if false == tmpServerInfo.Status {
				return false, errors.New("server [" + s.IP + ":" + s.HostPort + "] status is disabled")
			} else {
				return false, errors.New("server [" + s.IP + ":" + s.HostPort + "] is overload,load:" + strconv.Itoa(s.TaskLoad) + ",maxNum:" + strconv.Itoa(s.TaskNumTotal))
			}
		} else {
			return false, err
		}

	} else {
		return false, err
	}
}

func (s ServerInfo) PingURL() string {
	return "http://" + strings.Trim(s.HostIP, "\r") + ":" + strings.Trim(s.HostPort, "\r") + "/ping"
}

func (s ServerInfo) RegisterURL() string {
	return "http://" + s.IP + ":" + s.HostPort + "/election/register"
}

func (s ServerInfo) StatusURL() string {
	return "http://" + s.IP + ":" + s.HostPort + "/status"
}

func (s ServerInfo) SlaveTaskAddURL() string {
	return "http://" + s.IP + ":" + s.HostPort + "/client/task/add"
}

type ProxyInfo struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
	Desc string `json:"desc"`
}

/**
 {
	  "apiID": 1,
	  "uniqueID": "123",
	  "url": ["http://news.mydrivers.com/1/645/645356.htm"],
	  "openType": ["proxy"],
	  "client": ["pcshot"],
	  "beforeTask": ["autoLogin"],
	  "task": ["getHTML", "screenshot", "uploadGuardSAE"],
	  "afterTask": [""],
	  "finishAction": ["standardCallback"]

  }
*/
type TaskPreData struct {
	TaskID       string //算出来的 唯一id，用于判断url是否重复
	ApiID        int      `json:"apiID"`
	UniqueID     string   `json:"uniqueID"`
	URL          string   `json:"url"`
	OpenType     string   `json:"openType"`
	ProxyID      int      `json:"proxyID"`
	Client       string   `json:"client"`
	BeforeTask   []string `json:"beforeTask"`
	Task         []string `json:"task"`
	AfterTask    []string `json:"afterTask"`
	FinishAction []string `json:"finishAction"`
}

//URL : TaskMasterData
var TaskMasterDataList = sync.Map{}

//TaskMasterDataList = array(url=>array( taskID => array(taskApiData)))
//TaskMasterDataList[URL] =  TaskMasterData
type TaskMasterData struct {
	URL        string                    `json:"url"`
	ServerData map[string]TaskServerData `json:"serverData"` //taskID :TaskServerData
	ApiData    map[string][]TaskApiData  `json:"apiData"`    //taskID :TaskApiData
	Err        string                    `json:"err"`
}

type TaskMasterDataIDList struct {
	TaskID   string //算出来的 唯一id，用于判断url和打开方式是否重复
	ApiID    int    `json:"apiID"`
	UniqueID string `json:"uniqueID"`
}

type TaskApiData struct {
	TaskID   string //算出来的 唯一id，用于判断url和打开方式是否重复
	ApiID    int    `json:"apiID"`
	UniqueID string `json:"uniqueID"`

	URL      string `json:"url"`
	Client   string `json:"client"`
	OpenType string `json:"openType"`

	CreateTime time.Time `json:"createTime"`

	TaskMap string `json:"taskMap"`
}

func (t TaskApiData) GetTaskID() string {
	var tmpApiData = TaskApiData{
		URL:      t.URL,
		Client:   t.Client,
		TaskMap:  t.TaskMap,
		OpenType: t.OpenType,
	}
	str, _ := Json.MarshalToString(tmpApiData)
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

type TaskServerData struct {
	TaskID         string       `json:"taskID"`
	URL            string       `json:"url"`
	TimeoutSeconds int          `json:"timeoutSeconds"`
	CreateTime     time.Time    `json:"createTime"`
	TimeoutTime    time.Time    `json:"timeoutTime"`
	ProxyID        string       `json:"proxyID"`
	TaskMap        string       `json:"taskMap"`
	Task           []ClientTask `json:"task"`
	MasterIP       string       `json:"masterIP"`
	MasterPort     string       `json:"masterPort"`
	ServerIP       string       `json:"serverIP"`
	ServerPort     string       `json:"serverPort"`
}

func (t TaskServerData) ServerTaskAddURL() string {
	return "http://" + t.ServerIP + ":" + t.ServerPort + "/server/task/add"
}

type ClientTaskList struct {
	TaskID            string          `json:"taskID"`
	URL               string          `json:"url"`
	StartTime         time.Time       `json:"startTime"`   //开始执行时间
	TimeOutTime       time.Time       `json:"timeOutTime"` //超时时间
	TimeoutSeconds    int             `json:"timeoutSeconds"`
	CreateTime        time.Time       `json:"createTime"`
	ClientSuicideTime time.Time       `json:"clientSuicideTime"` //客户端存活后，自杀时间
	ProxyID           string          `json:"proxyID"`
	TaskMap           string          `json:"taskMap"`
	Task              []ClientTask    `json:"task"`
	History           []ClientHistory `json:"history"`
	CallBackURL       string          `json:"callBackUrl"`
}

type ClientTask struct {
	Opration string            `json:"opration"`
	Param    map[string]string `json:"param"`
}

type ClientHistory struct {
	Time       time.Time
	StatusCode int
	URL        string
}

type ChromeOption struct {
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	ScaleFactor bool   `json:"scaleFactor"`
	Mobile      bool   `json:"mobile"`
	FitWindow   bool   `json:"fitWindow"`
	UserAgent   string `json:"userAgent"`
	Desc        string `json:"desc"`
}

func (t TaskPreData) GetTaskID() string {
	var tmpData TaskPreData
	tmpData = t
	tmpData.TaskID = ""
	tmpData.ApiID = 0
	tmpData.UniqueID = ""
	str, _ := Json.MarshalToString(tmpData)
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

/**
map[HashID][]TaskPreData
*/
type TaskPreNodeData map[string][]TaskPreData

/**值为
map["www.baidu.com"][]"aaaaaaaa"
map["www.baidu.com"][]"bbbbbbbb"
这里的aaaa和bbbb 对应 TaskPreNodeData中的HashID
*/
type TaskPreDataTaskURLMap map[string][]string

type TaskData struct {
	ApiID        int      `json:"apiID"`
	UniqueID     string   `json:"uniqueID"`
	URL          string   `json:"url"`
	OpenType     string   `json:"openType"`
	Client       string   `json:"client"`
	ClientOption string   `json:"option"`
	BeforeTask   []string `json:"beforeTask"`
	Task         []string `json:"task"`
	AfterTask    []string `json:"afterTask"`
	FinishAction []string `json:"finishAction"`

	StartTime time.Time `json:"startTime"`
}

type StandardCallbackData struct {
	TaskID    string          `json:"taskID"`
	ApiID     string          `json:"apiID"`
	UniqueID  string          `json:"uniqueID"`
	URL       string          `json:"url"`
	FinishURL string          `json:"finishURL"`
	Client    string          `json:"client"`
	SaeURL    string          `json:"saeURL"`
	HTML      string          `json:"html"`
	History   []ClientHistory `json:"history"`
	Err       string          `json:"err"`
}

type ResultData struct {
	Code      int    `json:"code"`
	Err       string `json:"err"`
	Msg       string `json:"msg"`
	Data      string `json:"data"`
	Operation string `json:"operation"`
	Pid       int    `json:"pid"`
}

type ServerRetryData struct {
	RetryFormData url.Values
	RetryTime     int
	Err           error
}
