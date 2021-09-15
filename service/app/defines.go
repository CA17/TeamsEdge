package app

type JsonOptions struct {
	Id    interface{} `json:"id,string"`
	Value interface{} `json:"value"`
}

const (
	EdgeInformTask = "EdgeInformTask"

	// TeamsEdgeBootstrap 启动消息
	TeamsEdgeBootstrap = "teamsedge/bootstrap"
	// TeamsEdgeInform 定时上报信息
	TeamsEdgeInform = "teamsedge/inform"
	// TeamsEdgeHostpingReq 网络 Ping 请求
	TeamsEdgeHostpingReq = "teamsedge/hostping_req"
	// TeamsEdgeHostpingRep 网络 Ping 响应
	TeamsEdgeHostpingRep = "teamsedge/hostping_rep"
)
