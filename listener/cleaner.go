package listener

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	cron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type cleaner struct {
	m       *sync.Map
	gauge   *prometheus.GaugeVec
	minutes time.Duration
}

func newCleaner(gauge *prometheus.GaugeVec, minutes int) cleaner {
	c := cleaner{
		m:       &sync.Map{},
		gauge:   gauge,
		minutes: time.Duration(minutes),
	}

	croner := cron.New()
	_, _ = croner.AddFunc(fmt.Sprintf("@every %dm", minutes), func() {
		c.cleanup()
	})
	croner.Start()
	return c
}

func (c cleaner) add(labels prometheus.Labels) {
	bytes, err := json.Marshal(labels)
	if err != nil {
		logrus.Errorf("marshal labels error: %v", err)
		return
	}
	c.m.Store(string(bytes), time.Now())
}

func (c cleaner) cleanup() {
	beforeTime := time.Now().Add(-time.Minute * c.minutes)
	tobeDeletes := make([]interface{}, 0)
	c.m.Range(func(k, v interface{}) bool {
		var labels prometheus.Labels
		err := json.Unmarshal([]byte(k.(string)), &labels)
		if err != nil {
			logrus.Errorf("unmarshal labels error: %v", err)
			return true
		}
		updatedAt := v.(time.Time)
		if updatedAt.Before(beforeTime) {
			c.gauge.Delete(labels)
			tobeDeletes = append(tobeDeletes, k)
		}
		return true
	})

	for _, tobeDelete := range tobeDeletes {
		c.m.Delete(tobeDelete)
	}
}
