// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Scrape `SHOW BINARY LOGS`

package collector

import (
	"context"
	"database/sql"
    "strconv"
	"strings"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"os"
)

const (
	// Subsystem.
	backup = "backup"
)

// Metric descriptors.
var (
	backupStatusDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, backup, "status"),
		"Status of database backups.",
		[]string{}, nil,
	)
)

// ScrapeBackupStatus collects from `SHOW BINARY LOGS`.
type ScrapeBackupStatus struct{}

// Name of the Scraper. Should be unique.
func (ScrapeBackupStatus) Name() string {
	return "backup_status"
}

// Help describes the role of the Scraper.
func (ScrapeBackupStatus) Help() string {
	return "Collect the status of database backups"
}

// Version of MySQL from which scraper is available.
func (ScrapeBackupStatus) Version() float64 {
	return 5.1
}

// Scrape collects data from database connection and sends it over channel as prometheus metric.
func (ScrapeBackupStatus) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error {
    var status float64
	var delblank string
	filename := "/tmp/mysqlBakStatus.txt"
    content, err := os.ReadFile(filename)
    if err != nil {
		status = 10000
		ch <- prometheus.MustNewConstMetric(
			backupStatusDesc, prometheus.GaugeValue, status,
		)

        return nil
    }
	delblank = strings.Replace(string(content), " ", "", -1)
    delnewline, _ := strconv.ParseFloat(strings.Replace(delblank,"\n", "", -1), 64)
	if delnewline == 1 {
        status = 1
    } else {
		status = 0
	}
	ch <- prometheus.MustNewConstMetric(
		backupStatusDesc, prometheus.GaugeValue, status,
	)

	return nil
}

// check interface
var _ Scraper = ScrapeBackupStatus{}