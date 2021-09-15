package jobs

import (
	"sync"
	"time"

	"github.com/ca17/teamsedge/service/app"
	"github.com/ca17/teamsedge/service/models"
	"github.com/robfig/cron/v3"
)

var Sched *cron.Cron
var nameMapLock sync.Mutex
var nameMap = make(map[cron.EntryID]string, 0)
var Names = make([]string, 0)

type namedJob struct {
	Id  cron.EntryID
	err error
}

func newNamedJob(id cron.EntryID, err error) *namedJob {
	return &namedJob{Id: id, err: err}
}

func addNamedJob(name string, ljob *namedJob) {
	nameMapLock.Lock()
	defer nameMapLock.Unlock()
	nameMap[ljob.Id] = name
	Names = append(Names, name)
}

func Init() {
	loc, _ := time.LoadLocation(app.Config.System.Location)
	Sched = cron.New(cron.WithLocation(loc))

	// 定时上报任务
	addNamedJob(app.EdgeInformTask, newNamedJob(Sched.AddFunc("@every 60s", func() {
		app.Publish(app.TeamsEdgeInform, models.EdgeInformMessage{
			Eid: app.GetEdgeID(),
		})
	})))

	Sched.Start()
}
