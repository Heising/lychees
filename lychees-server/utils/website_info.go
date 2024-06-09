package utils

import (
	"compress/gzip"
	"context"
	"github.com/google/brotli/go/cbrotli"
	"golang.org/x/net/html"
	"io"
	"lychees-server/logs"
	"lychees-server/models"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

type websiteInfo struct {
	icon    string
	title   string
	clearAt int64
}

var websiteInfoMap = make(map[string]*websiteInfo)
var websiteInfoLock sync.RWMutex
var client = http.Client{}

// 定义全局常量用于存储固定的请求头信息
const (
	AcceptEncoding = "gzip,br"
	AcceptLanguage = "zh-CN,zh;q=0.9"
	UserAgent      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
	Accept         = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"
	Referer        = "https://redli.cn/"
)


func init() {

	go func() {
		// 每30分钟清除缓存的信息
		for range time.Tick(30 * time.Minute) {
			websiteInfoLock.Lock()
			logs.Logger.Infof("缓存网站信息的长度是：%d", len(websiteInfoMap))
			for index := range websiteInfoMap {
				if websiteInfoMap[index].clearAt < time.Now().Unix() {
					delete(websiteInfoMap, index)
				}

			}
			websiteInfoLock.Unlock()
		}
	}()

}

// 设置网站信息缓存
func setInfoCache(newBookMarkData *models.Item) {
	websiteInfoLock.Lock()
	defer websiteInfoLock.Unlock()
	websiteInfoMap[newBookMarkData.URL] = &websiteInfo{
		icon:    newBookMarkData.Icon,
		title:   newBookMarkData.Title,
		clearAt: time.Now().Add(time.Hour).Unix(),
	}
}

// 拿网站信息缓存
func getInfoCache(url string) *websiteInfo {
	websiteInfoLock.RLock()
	defer websiteInfoLock.RUnlock()
	logs.Logger.Infof("有人读缓存%s", url)
	if _, ok := websiteInfoMap[url]; ok {
		//if websiteInfoMap[url].clearAt < time.Now().Unix() {
		//	return nil
		//}
		return websiteInfoMap[url]
	}
	logs.Logger.Infof("没有缓存到%s", url)
	return nil
}
func newRequestTimeout(urlString string) (*http.Request, context.CancelFunc) {
	// 创建一个带有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	// 包装超时上下文
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		logs.Logger.Info(err)
		return nil, cancel
	}

	// 将上下文与请求关联
	req = req.WithContext(ctx)

	req.Header.Add("accept-encoding", AcceptEncoding)
	req.Header.Add("accept-language", AcceptLanguage)
	req.Header.Add("user-agent", UserAgent)
	req.Header.Add("accept", Accept)
	req.Header.Add("referer", Referer)
	return req, cancel
}

// 检测图标
func checkIcon(newBookMarkData *models.Item) {
	if strings.HasPrefix(newBookMarkData.Icon, "https://") {

	} else if strings.HasPrefix(newBookMarkData.Icon, "http://") {
		newBookMarkData.URL = strings.Replace(newBookMarkData.URL, "http://", "https://", 1)
	} else if strings.HasPrefix(newBookMarkData.Icon, "//") {
		newBookMarkData.Icon = "https:" + newBookMarkData.Icon
	} else if strings.HasPrefix(newBookMarkData.Icon, "/") {
		newBookMarkData.Icon = extractDomain(newBookMarkData.URL) + newBookMarkData.Icon
		logs.Logger.Info(newBookMarkData.Icon)
		logs.Logger.Info("字符串以斜杠开头")
	} else {
		logs.Logger.Info("字符串不以斜杠开头")
		newBookMarkData.Icon = extractDomain(newBookMarkData.URL) + "/" + newBookMarkData.Icon
	}

	requestIcon, cancelIcon := newRequestTimeout(newBookMarkData.Icon)
	if requestIcon == nil {
		logs.Logger.Info("newRequestTimeout失败")
		return
	}
	defer cancelIcon()
	responseIcon, err := client.Do(requestIcon)
	if err != nil {
		logs.Logger.Info("client Do失败")
		return
	}
	// 检查对方的状态码
	//logs.Logger.Infof("图标响应码是:%d", responseIcon.StatusCode)

	if responseIcon.StatusCode == http.StatusForbidden || responseIcon.StatusCode == http.StatusNotFound {
		logs.Logger.Info("有防盗链，或者404，默认添加/favicon.ico")
		newBookMarkData.Icon = extractDomain(newBookMarkData.URL) + "/favicon.ico"
	}
}
func GetFaviconTitle(newBookMarkData *models.Item) {
	//先处理链接
	var urlString string
	if strings.HasPrefix(newBookMarkData.URL, "https://") || strings.HasPrefix(newBookMarkData.URL, "http://") {
		urlString = newBookMarkData.URL
	} else if strings.HasPrefix(newBookMarkData.URL, "//") {
		newBookMarkData.URL = "https:" + newBookMarkData.URL
		urlString = newBookMarkData.URL
	} else if strings.HasPrefix(newBookMarkData.URL, "/") {
		newBookMarkData.URL = "https:/" + newBookMarkData.URL
		urlString = newBookMarkData.URL
	} else {
		newBookMarkData.URL = "https://" + newBookMarkData.URL
		urlString = newBookMarkData.URL
	}
	//从缓存拿信息
	cache := getInfoCache(newBookMarkData.URL)
	if cache != nil {
		if !newBookMarkData.IsSvg && newBookMarkData.Icon == "" {
			newBookMarkData.Icon = cache.icon

		}

		if newBookMarkData.Title == "" {
			newBookMarkData.Title = cache.title
		}

		return
	}
	request, cancel := newRequestTimeout(urlString)
	if request == nil {
		return
	}
	defer cancel()
	resp, err := client.Do(request)
	if err != nil {
		// 处理超时错误
		if err, ok := err.(net.Error); ok && err.Timeout() {

			logs.Logger.Info("请求超时")

			if newBookMarkData.Title == "" {
				newBookMarkData.Title = getDomain(newBookMarkData.URL)
			}

			if newBookMarkData.Icon == "" {
				logs.Logger.Info("未找到匹配的图标，默认添加/favicon.ico")
				newBookMarkData.Icon = extractDomain(newBookMarkData.URL) + "/favicon.ico"
			}
			return

		}
		//不知道啥错误，直接打印
		logs.Logger.Info(err)
		return
	}
	// 获取并检测响应的状态码
	statusCode := resp.StatusCode
	logs.Logger.Infof("响应状态码:%d", statusCode)
	// 检测响应头中的"Location"字段
	location := resp.Header.Get("Location")
	if location != "" {
		logs.Logger.Info("重定向URL:", location)
	} else {
		logs.Logger.Info("未检测到重定向URL")
	}
	// 检查响应的协议
	if resp.Request.URL.Scheme == "https" {
		//fmt.Println("该网站支持 HTTPS")
		if strings.HasPrefix(newBookMarkData.URL, "http://") {
			newBookMarkData.URL = strings.Replace(newBookMarkData.URL, "http://", "https://", 1)
		}
	}
	//if err != nil {
	//	logs.Logger.Info(err)
	//	return
	//
	//}
	defer resp.Body.Close()
	var titleIcon *websiteInfo
	//var respBody []byte
	// 检查响应是否使用Brotli编码
	if resp.Header.Get("Content-Encoding") == "br" {
		// 创建Brotli读取器
		reader := cbrotli.NewReader(resp.Body)

		//respBody, err = io.ReadAll(reader)
		//if err != nil {
		//	logs.Logger.Info("error decoding br response", err)
		//	return
		//}
		defer reader.Close()

		titleIcon = extractTitleIcon(reader)

		// 打印解码后的响应
		//logs.Logger.Info("Decoded cbrotli response:\n", string(respBody))
	} else if resp.Header.Get("Content-Encoding") == "gzip" {

		// 读取HTML内容
		reader, err := gzip.NewReader(resp.Body)
		//respBody, err = io.ReadAll(reader)

		if err != nil {
			logs.Logger.Info("error decoding gzip response", err)
			return
		}
		defer reader.Close()

		titleIcon = extractTitleIcon(reader)

		// 打印解码后的响应
		//logs.Logger.Info("Decoded gzip response:\n", string(respBody))
	} else {
		//respBody, err = io.ReadAll(resp.Body)
		//if err != nil {
		//	logs.Logger.Info("error response", err)
		//}
		titleIcon = extractTitleIcon(resp.Body)

		// 打印解码后的响应
		//logs.Logger.Info("Decoded response:\n", string(respBody))
	}
	if titleIcon == nil {
		if newBookMarkData.Title == "" {
			newBookMarkData.Title = getDomain(newBookMarkData.URL)
		}

		if newBookMarkData.Icon == "" {
			logs.Logger.Info("未找到匹配的图标，默认添加/favicon.ico")
			newBookMarkData.Icon = extractDomain(newBookMarkData.URL) + "/favicon.ico"
		}
		return

	}
	if !newBookMarkData.IsSvg && newBookMarkData.Icon == "" {
		if titleIcon.icon != "" {
			newBookMarkData.Icon = titleIcon.icon
			checkIcon(newBookMarkData)
		} else {
			logs.Logger.Info("未找到匹配的图标，默认添加/favicon.ico")
			newBookMarkData.Icon = extractDomain(newBookMarkData.URL) + "/favicon.ico"
		}
	}

	if newBookMarkData.Title == "" {
		//拿标题
		logs.Logger.Info("拿标题")

		if titleIcon.title != "" {
			logs.Logger.Info("匹配到标题")

			newBookMarkData.Title = titleIcon.title
		} else {
			logs.Logger.Info("未找到匹配的标题")
			newBookMarkData.Title = getDomain(newBookMarkData.URL)

		}
	}

	setInfoCache(newBookMarkData)

}

// 提取域名端口和返还https协议
func extractDomain(urlString string) string {
	// 去除字符串两端的空格
	urlString = strings.TrimSpace(urlString)

	// 使用url.Parse解析URL
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		logs.Logger.Info(err)
	}

	// 提取域名部分
	//domainParts := strings.Split(parsedURL.Hostname(), ".")
	//domain := domainParts[len(domainParts)-2] + "." + domainParts[len(domainParts)-1]
	// Extract the protocol and domain
	//protocol := parsedURL.Scheme
	return "https://" + parsedURL.Host
}

// 只提取域名
func getDomain(urlString string) string {
	// 去除字符串两端的空格
	urlString = strings.TrimSpace(urlString)

	// 使用url.Parse解析URL
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		logs.Logger.Info(err)
	}
	domainParts := strings.Split(parsedURL.Hostname(), ".")
	domain := domainParts[len(domainParts)-2]
	return domain
}

func extractTitleIcon(htmlReader io.Reader) *websiteInfo {

	var titleIcon websiteInfo

	// 解析HTML
	root, err := html.Parse(htmlReader)
	if err != nil {
		logs.Logger.Infof("HTML解析错误:", err)
		return nil
	}
	// 变量用于存储最大属性值
	maxSize := 0
	// 处理每个节点
	var processNode func(*html.Node)
	processNode = func(n *html.Node) {
		// fmt.Println("n.Type :", n.Type)

		if (n.Type == html.ElementNode && (n.Data == "html" || n.Data == "head")) || n.Type == html.DocumentNode {
			// 递归处理子节点
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				processNode(c)

			}
			return

		}
		if n.Type == html.ElementNode && n.Data == "title" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				logs.Logger.Infof("标题是:%s", c.Data)
				titleIcon.title = c.Data
			}
		}

		// 捕获icon
		if n.Type == html.ElementNode && n.Data == "link" {
			var capture = false

			for _, attr := range n.Attr {
				if attr.Key == "rel" && (attr.Val == "icon" || attr.Val == "shortcut icon") {

					if attr.Val == "icon" || attr.Val == "shortcut icon" {


						capture = true

					}

				}

			}
			if capture {

				var num int

				for _, attr := range n.Attr {
					if attr.Key == "sizes" {


						// 查找第一个非数字字符的位置
						index := strings.IndexFunc(attr.Val, func(r rune) bool {
							return !unicode.IsDigit(r)
						})

						// 提取数字部分
						var numStr string
						if index != -1 {
							numStr = attr.Val[:index]
						} else {
							numStr = attr.Val
						}

						// 将提取的数字部分转换为整数
						num, err = strconv.Atoi(numStr)
						if err != nil {
							logs.Logger.Error(err)
							return
						}

					}

				}

				// 比较属性值并保留最大值
				if num >= maxSize {
					maxSize = num
					for _, attr := range n.Attr {
						if attr.Key == "href" {
							titleIcon.icon = attr.Val

						}

					}

				}

			}
		}

	}

	// 从根节点开始处理
	processNode(root)

	logs.Logger.Infof("提取的数据:%v", titleIcon)
	return &titleIcon
}
