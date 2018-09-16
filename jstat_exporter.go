package main

import (
	"flag"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

const (
	namespace = "jstat"
)

var (
	listenAddress = flag.String("web.listen-address", ":9010", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	jstatPath     = flag.String("jstat.path", "/usr/bin/jstat", "jstat path")
	targetPid     = flag.String("target.pid", ":0", "target pid")
)

type Exporter struct {
	jstatPath  string
	targetPid  string
	sv0Cur     prometheus.Gauge
	sv1Cur     prometheus.Gauge
	sv0Used    prometheus.Gauge
	sv1Used    prometheus.Gauge
	edenCur    prometheus.Gauge
	edenUsed   prometheus.Gauge
	oldCur     prometheus.Gauge
	oldUsed    prometheus.Gauge
	metaCur    prometheus.Gauge
	metaUsed   prometheus.Gauge
	classCur   prometheus.Gauge
	classUsed  prometheus.Gauge
	ygcTimes   prometheus.Gauge
	ygcSec    prometheus.Gauge
	fgcTimes   prometheus.Gauge
	fgcSec    prometheus.Gauge
	gcSec     prometheus.Gauge
}

func NewExporter(jstatPath string, targetPid string) *Exporter {
	return &Exporter{
		jstatPath: jstatPath,
		targetPid: targetPid,
		sv0Cur: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "sv0Cur",
			Help:      "sv0Cur",
		}),
		sv1Cur: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "sv1Cur",
			Help:      "sv1Cur",
		}),
		sv0Used: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "sv0Used",
			Help:      "sv0Used",
		}),
		sv1Used: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "sv1Used",
			Help:      "sv1Used",
		}),
		edenCur: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "edenCur",
			Help:      "edenCur",
		}),
		edenUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "edenUsed",
			Help:      "edenUsed",
		}),
		oldCur: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "oldCur",
			Help:      "oldCur",
		}),
		oldUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "oldUsed",
			Help:      "oldUsed",
		}),
		metaCur: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "metaCur",
			Help:      "metaCur",
		}),
		metaUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "metaUsed",
			Help:      "metaUsed",
		}),
		classCur: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "classCur",
			Help:      "classCur",
		}),
		classUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "classUsed",
			Help:      "classUsed",
		}),
		ygcTimes: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "ygcTimes",
			Help:      "ygcTimes",
		}),
		ygcSec: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "ygcSec",
			Help:      "ygcSec",
		}),
		fgcTimes: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "fgcTimes",
			Help:      "fgcTimes",
		}),
		fgcSec: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "fgcSec",
			Help:      "fgcSec",
		}),
		gcSec: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "gcSec",
			Help:      "gcSec",
		}),
	}
}

// Describe implements the prometheus.Collector interface.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.sv0Cur.Describe(ch)
	e.sv1Cur.Describe(ch)
	e.sv0Used.Describe(ch)
	e.sv1Used.Describe(ch)
	e.edenCur.Describe(ch)
	e.edenUsed.Describe(ch)
	e.oldCur.Describe(ch)
	e.oldUsed.Describe(ch)
	e.metaCur.Describe(ch)
	e.metaUsed.Describe(ch)
	e.classCur.Describe(ch)
	e.classUsed.Describe(ch)
	e.ygcTimes.Describe(ch)
	e.ygcSec.Describe(ch)
	e.fgcTimes.Describe(ch)
	e.fgcSec.Describe(ch)
	e.gcSec.Describe(ch)
}

// Collect implements the prometheus.Collector interface.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.JstatGc(ch)
}


func (e *Exporter) JstatGc(ch chan<- prometheus.Metric) {
	out, err := exec.Command(e.jstatPath, "-gc", e.targetPid).Output()
	if err != nil {
		log.Fatal(err)
	}

	for i, line := range strings.Split(string(out), "\n") {
		if i == 1 {
			parts := strings.Fields(line)

			sv0Cur, err := strconv.ParseFloat(parts[0], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.sv0Cur.Set(sv0Cur)
			e.sv0Cur.Collect(ch)
			sv1Cur, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.sv1Cur.Set(sv1Cur)
			e.sv1Cur.Collect(ch)
			sv0Used, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.sv0Used.Set(sv0Used)
			e.sv0Used.Collect(ch)
			sv1Used, err := strconv.ParseFloat(parts[3], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.sv1Used.Set(sv1Used)
			e.sv1Used.Collect(ch)
			edenCur, err := strconv.ParseFloat(parts[4], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.edenCur.Set(edenCur)
			e.edenCur.Collect(ch)
			edenUsed, err := strconv.ParseFloat(parts[5], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.edenUsed.Set(edenUsed)
			e.edenUsed.Collect(ch)
			oldCur, err := strconv.ParseFloat(parts[6], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.oldCur.Set(oldCur)
			e.oldCur.Collect(ch)
			oldUsed, err := strconv.ParseFloat(parts[7], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.oldUsed.Set(oldUsed)
			e.oldUsed.Collect(ch)
			metaCur, err := strconv.ParseFloat(parts[8], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.metaCur.Set(metaCur)
			e.metaCur.Collect(ch)
			metaUsed, err := strconv.ParseFloat(parts[9], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.metaUsed.Set(metaUsed)
			e.metaUsed.Collect(ch)
			classCur, err := strconv.ParseFloat(parts[10], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.classCur.Set(classCur)
			e.classCur.Collect(ch)
			classUsed, err := strconv.ParseFloat(parts[11], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.classUsed.Set(classUsed)
			e.classUsed.Collect(ch)
			ygcTimes, err := strconv.ParseFloat(parts[12], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.ygcTimes.Set(ygcTimes)
			e.ygcTimes.Collect(ch)
			ygcSec, err := strconv.ParseFloat(parts[13], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.ygcSec.Set(ygcSec)
			e.ygcSec.Collect(ch)
			fgcTimes, err := strconv.ParseFloat(parts[14], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.fgcTimes.Set(fgcTimes)
			e.fgcTimes.Collect(ch)
			fgcSec, err := strconv.ParseFloat(parts[15], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.fgcSec.Set(fgcSec)
			e.fgcSec.Collect(ch)
			gcSec, err := strconv.ParseFloat(parts[16], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.gcSec.Set(gcSec)
			e.gcSec.Collect(ch)
		}
	}
}

func main() {
	flag.Parse()

	exporter := NewExporter(*jstatPath, *targetPid)
	prometheus.MustRegister(exporter)

	log.Printf("Starting Server: %s", *listenAddress)
	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head><title>jstat Exporter</title></head>
		<body>
		<h1>jstat Exporter</h1>
		<p><a href="` + *metricsPath + `">Metrics</a></p>
		</body>
		</html>`))
	})
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
