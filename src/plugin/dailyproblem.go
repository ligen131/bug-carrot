package plugin

import (
	"bug-carrot/controller"
	"bug-carrot/model"
	"bug-carrot/param"
	"bug-carrot/util"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type dailyproblemWords struct {
	DailyProblem string
	Today        string
	Yesterday    string
	Lastday      string
	Day          string
	Setting      string
	Problem      string
	Solution     string
	Query        string
	Div1         string
	Div2         string
	Help1        string
	Help2        string
	HelpQuery    string
	HelpSetting  string
}

type dailyproblem struct {
	Index             param.PluginIndex
	Words             dailyproblemWords
	ProblemAnnounced  string
	SolutionAnnounced string
}

var dailyproblemAdminGroup = []int64{
	444185515,
	706605585,
	285976171,
	786012798,
} // 💩

var dailyproblemGroup = []int64{
	706605585,
	744140247,
} // 💩

func (p *dailyproblem) GetPluginName() string {
	return p.Index.PluginName
}
func (p *dailyproblem) GetPluginAuthor() string {
	return p.Index.PluginAuthor
}
func (p *dailyproblem) CanTime() bool {
	return p.Index.FlagCanTime
}
func (p *dailyproblem) CanMatchedGroup() bool {
	return p.Index.FlagCanMatchedGroup
}
func (p *dailyproblem) CanMatchedPrivate() bool {
	return p.Index.FlagCanMatchedPrivate
}
func (p *dailyproblem) CanListen() bool {
	return p.Index.FlagCanListen
}
func (p *dailyproblem) NeedDatabase() bool {
	return p.Index.FlagUseDatabase
}
func (p *dailyproblem) DoIgnoreRiskControl() bool {
	return p.Index.FlagIgnoreRiskControl
}

// 以上六个函数在注册时被唯一调用，并以此为依据加入相应的 queue
// 无需修改

// IsTime : 是你需要的时间吗？
func (p *dailyproblem) IsTime() bool {
	now := time.Now()
	return ((now.Hour() >= 8 && p.ProblemAnnounced != now.Format("20060102")) ||
		(now.Hour() >= 20 && p.SolutionAnnounced != now.Format("20060102"))) && now.Minute() <= 10
}

// DoTime : 当到了你需要的时间，要做什么呢？
func (p *dailyproblem) DoTime() error {
	now := time.Now()
	if now.Hour() >= 20 && p.SolutionAnnounced != now.Format("20060102") {
		query, err := p.queryDailyproblem(now.Format("20060102"), false, false)
		if err != nil {
			util.ErrorPrint(err, nil, "定时任务出错！无法发送题目公告。")
			return err
		}
		for _, group := range dailyproblemGroup {
			query += `今天的题解来啦！\n`
			util.QQGroupSend(group, query)
		}
		p.SolutionAnnounced = now.Format("20060102")
	} else if now.Hour() >= 8 && p.ProblemAnnounced != now.Format("20060102") && p.SolutionAnnounced != now.Format("20060102") {
		query, err := p.queryDailyproblem(now.Format("20060102"), false, true)
		if err != nil {
			util.ErrorPrint(err, nil, "定时任务出错！无法发送题目公告。")
			return err
		}
		for _, group := range dailyproblemGroup {
			query += `今天的每日一题来啦！\n`
			util.QQGroupSend(group, query)
		}
		p.ProblemAnnounced = now.Format("20060102")
	}
	return nil
}

// IsMatchedGroup : 是你想收到的群 @ 消息吗？
func (p *dailyproblem) IsMatchedGroup(msg param.GroupMessage) bool {
	return msg.ExistWord("r", []string{"每日"}) &&
		msg.ExistWord("n", []string{"一题"})
}

// DoMatchedGroup : 收到了想收到的群 @ 消息，要做什么呢？
func (p *dailyproblem) DoMatchedGroup(msg param.GroupMessage) error {
	sendText, _ := p.doDailyproblemActions(msg.GroupId, msg.UserId, msg.RawMessage)
	group := msg.GroupId
	util.QQGroupSend(group, sendText)
	return nil
}

// IsMatchedPrivate : 是你想收到的私聊消息吗？
func (p *dailyproblem) IsMatchedPrivate(msg param.PrivateMessage) bool {
	return msg.ExistWord("r", []string{"每日"}) &&
		msg.ExistWord("n", []string{"一题"})
}

// DoMatchedPrivate : 收到了想收到的私聊消息，要做什么呢？
// 备注：我们建议大部分功能只对群聊开启，增强 bot 在群聊中的存在感，私聊功能可以提供给管理员
func (p *dailyproblem) DoMatchedPrivate(msg param.PrivateMessage) error {
	sendText, _ := p.doDailyproblemActions(0, msg.UserId, msg.RawMessage)
	user := msg.UserId
	util.QQSend(user, sendText)
	return nil
}

// Listen : 监听到非 @ 的群消息，要做什么呢？
// 备注：我们建议只对极少数的功能采取监听行为
// 除去整活效果较好的特殊场景，我们一般希望 bot 只有在被 @ 到的时候才会对应s发言
func (p *dailyproblem) Listen(msg param.GroupMessage) {
	if msg.ExistWord("r", []string{"每日"}) && msg.ExistWord("n", []string{"一题"}) {
		sendText, _ := p.doDailyproblemActions(msg.GroupId, msg.UserId, msg.RawMessage)
		group := msg.GroupId
		util.QQGroupSend(group, sendText)
	}
}

// Close : 项目要关闭了，要做什么呢？
func (p *dailyproblem) Close() {
}

// DefaultPluginRegister : 创造一个插件实例，并调用 controller.PluginRegister
// 在 main.go 的 pluginRegister 函数中调用来实现注册
func DailyproblemPluginRegister() {
	p := &dailyproblem{
		Index: param.PluginIndex{
			PluginName:            "dailyproblem", // 插件名称
			PluginAuthor:          "ligen131",     // 插件作者
			FlagCanTime:           true,           // 是否能在特殊时间做出行为
			FlagCanMatchedGroup:   true,           // 是否能回应群聊@消息
			FlagCanMatchedPrivate: true,           // 是否能回应私聊消息
			FlagCanListen:         true,           // 是否能监听群消息
			FlagUseDatabase:       true,           // 是否用到了数据库（配置文件中配置不使用数据库的话，用到了数据库的插件会不运行）
			FlagIgnoreRiskControl: false,          // 是否无视风控（为 true 且 RiskControl=true 时将自动无视群聊功能，建议设置为 false）
		},
		Words: dailyproblemWords{
			DailyProblem: "每日一题",
			Today:        "今",
			Yesterday:    "昨",
			Lastday:      "前",
			Day:          "天",
			Setting:      "设置",
			Problem:      "题目",
			Solution:     "题解",
			Query:        "查询",
			Div1:         "div1",
			Div2:         "div2",
			Help1:        "帮助",
			Help2:        "help",
			HelpQuery: `每日一题 帮助列表
[今天/昨天/n*前天]每日一题
每日一题 查询 {date} - date 格式为 yyyymmdd，如 20221001
每日一题 帮助 - 打开此帮助

每天 8:00 公布当天每日一题，20:00 公布题解（自动）
前往 https://vjudge.net/group/hustacm 查看历史每日一题
每日一题投稿：https://docs.qq.com/sheet/DV0t0RGZBV1ZMdHhz
`,
			HelpSetting: `每日一题 管理员设置
[今天/昨天/n*前天]每日一题
每日一题 查询 {date} - date 格式为 yyyymmdd，如 20221001
每日一题 帮助 - 打开此帮助
每日一题 设置 {date} {link1} {link2} {link3} {link4} - 四个{link}分别代表 div1 题目, div1 题解, div2 题目, div2 题解
每日一题 设置 {date} div1/div2 题目/题解 {link} - 更新/设置特定 div 题目/题解链接，注意参数中间有空格，div 和数字之间没有空格，如 每日一题 设置 div1 题解 https://codeforces/...

每天 8:00 公布当天每日一题，20:00 公布题解（自动）
前往 https://vjudge.net/group/hustacm 查看历史每日一题
每日一题投稿：https://docs.qq.com/sheet/DV0t0RGZBV1ZMdHhz
`,
		},
		ProblemAnnounced:  "20220930",
		SolutionAnnounced: "20220930",
	}
	controller.PluginRegister(p)
}

func (p *dailyproblem) queryDailyproblem(date string, isAdmin bool, isAnnounce bool) (string, error) {
	m := model.GetModel()
	defer m.Close()

	dp, err := m.GetDailyproblemByDate(date)
	if (err != nil || dp == param.Dailyproblem{}) {
		util.ErrorPrint(err, date, fmt.Sprintf("查询每日一题出错 date=%s", date))
		return fmt.Sprintf("%s 没有每日一题哦\n", date), err
	}
	req := ""
	now := time.Now()
	year, _ := strconv.Atoi(date[0:4])
	month, _ := strconv.Atoi(date[4:6])
	day, _ := strconv.Atoi(date[6:8])
	if !isAdmin && (year > now.Year() || month > int(now.Month()) || day > now.Day()) {
		date = now.Format("20060102")
		dp, err = m.GetDailyproblemByDate(date)
		if (err != nil || dp == param.Dailyproblem{}) {
			util.ErrorPrint(err, date, fmt.Sprintf("查询每日一题出错 date=%s", date))
			return fmt.Sprintf("%s 没有每日一题哦\n", date), err
		}
	}
	if !isAdmin && date == now.Format("20060102") {
		dpy, err := m.GetDailyproblemByDate(now.Add(-24 * time.Hour).Format("20060102"))
		if now.Hour() < 8 || dp.Div1Problem == "" || dp.Div2Problem == "" {
			req += "今天的题目暂未公布，请稍等一会吧！\n"
			if err == nil && dpy.Div1Problem != "" && dpy.Div2Problem != "" {
				req += fmt.Sprintf("【昨日题目】\nDiv.1: %s\nDiv.2: %s\n", dpy.Div1Problem, dpy.Div2Problem)
			}
		} else {
			req += fmt.Sprintf("【今日题目】\nDiv.1: %s\nDiv.2: %s\n", dp.Div1Problem, dp.Div2Problem)
		}
		if now.Hour() < 20 || dp.Div1Solution == "" || dp.Div2Solution == "" {
			if !isAnnounce {
				req += "今天的题解暂未公布。\n"
			}
			if err == nil && dpy.Div1Solution != "" && dpy.Div2Solution != "" {
				req += fmt.Sprintf("【昨日题解】\nDiv.1: %s\nDiv.2: %s\n", dpy.Div1Solution, dpy.Div2Solution)
			}
		} else {
			req += fmt.Sprintf("【今日题解】\nDiv.1: %s\nDiv.2: %s\n", dp.Div1Solution, dp.Div2Solution)
		}
	} else {
		req += fmt.Sprintf("【%s 每日一题】\nDiv.1: %s\nDiv.2: %s\n【题解】\nDiv.1: %s\nDiv.2: %s\n", date, dp.Div1Problem, dp.Div2Problem, dp.Div1Solution, dp.Div2Solution)
	}
	return req, nil
}

func (p *dailyproblem) updateDailyproblem(date string, div int, isProblem bool, link string) (string, error) {
	m := model.GetModel()
	defer m.Close()

	query, err := p.queryDailyproblem(date, true, false)
	req := "更新前查询结果：\n" + query
	dp, err := m.UpdateDailyproblem(date, div, isProblem, link)
	if (err != nil || dp == param.Dailyproblem{}) {
		util.ErrorPrint(err, nil, fmt.Sprintf("更新每日一题失败！date=%s, div=%d, isProblem=%v, link=%s", date, div, isProblem, link))
		return req + "\n⚠更新每日一题失败！请检查日期是否正确。\n", err
	}
	query, err = p.queryDailyproblem(date, true, false)
	req += "\n更新成功！\n" + query
	return req, nil
}

func (p *dailyproblem) updateDailyproblemAll(date string, link1 string, link2 string, link3 string, link4 string) (string, error) {
	m := model.GetModel()
	defer m.Close()

	query, err := p.queryDailyproblem(date, true, false)
	req := "更新前查询结果：\n" + query
	dp, err := m.UpdateDailyproblem(date, 1, true, link1)
	if (err != nil || dp == param.Dailyproblem{}) {
		util.ErrorPrint(err, nil, fmt.Sprintf("更新每日一题失败！date=%s, link1=%s, link2=%s, link3=%s, link4=%s", date, link1, link2, link3, link4))
		return req + "\n⚠更新每日一题失败！请检查日期是否正确。\n", err
	}
	dp, err = m.UpdateDailyproblem(date, 1, false, link2)
	if (err != nil || dp == param.Dailyproblem{}) {
		util.ErrorPrint(err, nil, fmt.Sprintf("更新每日一题失败！date=%s, link1=%s, link2=%s, link3=%s, link4=%s", date, link1, link2, link3, link4))
		return req + "\n⚠更新每日一题失败！请检查日期是否正确。\n", err
	}
	dp, err = m.UpdateDailyproblem(date, 2, true, link3)
	if (err != nil || dp == param.Dailyproblem{}) {
		util.ErrorPrint(err, nil, fmt.Sprintf("更新每日一题失败！date=%s, link1=%s, link2=%s, link3=%s, link4=%s", date, link1, link2, link3, link4))
		return req + "\n⚠更新每日一题失败！请检查日期是否正确。\n", err
	}
	dp, err = m.UpdateDailyproblem(date, 2, false, link4)
	if (err != nil || dp == param.Dailyproblem{}) {
		util.ErrorPrint(err, nil, fmt.Sprintf("更新每日一题失败！date=%s, link1=%s, link2=%s, link3=%s, link4=%s", date, link1, link2, link3, link4))
		return req + "\n⚠更新每日一题失败！请检查日期是否正确。\n", err
	}

	query, err = p.queryDailyproblem(date, true, false)
	req += "\n更新成功！\n" + query
	return req, nil
}

func (p *dailyproblem) doDailyproblemActions(GroupId int64, UserId int64, message string) (string, error) {
	msg := strings.Split(message, " ")
	size := len(msg)
	if size == 0 {
		return "", fmt.Errorf("无效消息")
	}
	if strings.Index(msg[0], "CQ:at") >= 0 {
		msg = msg[1:]
	}
	ti := time.Now()
	isAdmin := false
	for _, account := range dailyproblemAdminGroup {
		if account == GroupId {
			isAdmin = true
			break
		}
	}

	// [今天/昨天/n*前天]每日一题
	// 每日一题 查询 {date} - date 格式为 yyyymmdd，如 20221001
	// 每日一题 设置 {date} {link1} {link2} {link3} {link4} - 四个{link}分别代表 div1 题目, div1 题解, div2 题目, div2 题解
	// 每日一题 设置 {date} div1/div2 题目/题解 {link} - 更新/设置特定 div 题目/题解链接，注意参数中间有空格，div 和数字之间没有空格，如 每日一题 设置 div1 题解 https://codeforces/...
	if ind := strings.Index(msg[0], p.Words.Day); ind > 0 {
		sub := msg[0][0:ind]
		if sub == p.Words.Yesterday {
			ti = ti.Add(-24 * time.Hour)
		} else {
			ti = ti.Add(time.Duration(-24 * (strings.Count(sub, p.Words.Lastday) + 1) * int(time.Hour)))
		}
	} else if msg[0] == p.Words.DailyProblem {
		if size >= 3 && msg[1] == p.Words.Query {
			return p.queryDailyproblem(msg[2], isAdmin, false)
		} else if size >= 2 && (msg[1] == p.Words.Help1 || msg[1] == p.Words.Help2) {
			if isAdmin {
				return p.Words.HelpSetting, nil
			} else {
				return p.Words.HelpQuery, nil
			}
		} else if size >= 6 && msg[1] == p.Words.Setting && isAdmin {
			if msg[3] == p.Words.Div1 || msg[3] == p.Words.Div2 {
				date := msg[2]
				div := 1
				isProblem := true
				link := msg[5]
				if msg[3] == p.Words.Div2 {
					div = 2
				}
				if msg[4] == p.Words.Solution {
					isProblem = false
				}
				return p.updateDailyproblem(date, div, isProblem, link)
			} else if size >= 7 {
				return p.updateDailyproblemAll(msg[2], msg[3], msg[4], msg[5], msg[6])
			} else {
				return "命令无效，请检查输入\n" + p.Words.HelpSetting, fmt.Errorf("管理命令无效")
			}
		}
	}
	return p.queryDailyproblem(ti.Format("20060102"), isAdmin, false)
}
