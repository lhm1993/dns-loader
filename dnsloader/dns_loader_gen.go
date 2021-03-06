package dnsloader

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"github.com/briandowns/spinner"
	"go.uber.org/ratelimit"
)

//GloablGenerator define  global control object
var GloablGenerator Generator

type myDNSLoaderGenerator struct {
	caller     Caller
	timeout    time.Duration
	qps        uint32
	status     uint32
	duration   time.Duration
	ctx        context.Context
	cancelFunc context.CancelFunc
	callCount  uint64
	workers    int
	result     map[uint8]uint64
}

func (mlg *myDNSLoaderGenerator) init() {
	log.Println("initial common loader...")
	mlg.result = make(map[uint8]uint64)
	log.Printf("initial Process Done QPS[%d]", mlg.qps)
}

func (mlg *myDNSLoaderGenerator) Start() bool {
	log.Println("starting Loader...")
	mlg.ctx, mlg.cancelFunc = context.WithTimeout(context.Background(), mlg.duration)
	mlg.callCount = 0
	currentStatus := mlg.Status()
	if currentStatus != STATUS_STARTING && currentStatus != STATUS_STOPPED {
		return false
	}
	atomic.StoreUint32(&mlg.status, STATUS_STARTING)
	if mlg.qps > 0 {
		interval := time.Duration(1e9 / mlg.qps)
		log.Printf("setting throttle %v", interval)
	}
	atomic.StoreUint32(&mlg.status, STATUS_STARTED)

	log.Println("new goroutine to generating dns packets")
	s := spinner.New(spinner.CharSets[36], 100*time.Millisecond)
	log.Println("new goroutine to receive dns data from server")
	go func() {
		// recive data from connections
		b := make([]byte, 4)
		dnsclient := mlg.caller.(*DNSClient)
		for {
			n, err := dnsclient.Conn.Read(b)
			if err == nil && n > 0 {
				code := b[3] & 0x0f
				mlg.result[code] = mlg.result[code] + 1
			}
		}
	}()

	var limiter ratelimit.Limiter
	if mlg.qps > 0 {
		limiter = ratelimit.New(int(mlg.qps))
	}

	mlg.generatorLoad(limiter, s)
	log.Println("waiting for program shutdown....")
	time.Sleep(5)
	log.Printf("[Result]total packets sum:%d", mlg.CallCount())
	log.Printf("[Result]runing time %v", mlg.duration)
	var counter uint64
	for k, v := range mlg.result {
		counter = v + counter
		log.Printf("[Result]status %s:%d [%.2f%%]", DNSRcodeReverse[k], v, float64(v*100)/float64(mlg.CallCount()))
	}
	restUnknown := mlg.CallCount() - counter
	log.Printf("[Result]status unknown:%d [%.2f%%]", restUnknown, float64(restUnknown*100)/float64(mlg.CallCount()))
	return true
}

func (mlg *myDNSLoaderGenerator) prepareStop(err error) {
	log.Printf("prepare to stop load test [%s]\n", err)
	atomic.StoreUint32(&mlg.status, STATUS_STOPPING)
	log.Println("try to stop channel...")
	atomic.StoreUint32(&mlg.status, STATUS_STOPPED)
	log.Println("stop load test success!")
}

func (mlg *myDNSLoaderGenerator) sendNewRequest() {
	defer func() {
		if p := recover(); p != nil {
			err, _ := interface{}(p).(error)
			log.Println(err)
		}
	}()
	rawRequest := mlg.caller.BuildReq()
	mlg.caller.Call(rawRequest)
}

func (mlg *myDNSLoaderGenerator) generatorLoad(limiter ratelimit.Limiter, spinnerInstance *spinner.Spinner) {
	spinnerInstance.Start()
	if mlg.qps > 0 {
		for {
			select {
			case <-mlg.ctx.Done():
				spinnerInstance.Stop()
				mlg.prepareStop(mlg.ctx.Err())
				return
			default:
			}
			limiter.Take()
			rawRequest := mlg.caller.BuildReq()
			mlg.caller.Call(rawRequest)
			atomic.AddUint64(&mlg.callCount, 1)
		}
	} else {
		for {
			select {
			case <-mlg.ctx.Done():
				spinnerInstance.Stop()
				mlg.prepareStop(mlg.ctx.Err())
				return
			default:
			}
			rawRequest := mlg.caller.BuildReq()
			mlg.caller.Call(rawRequest)
			atomic.AddUint64(&mlg.callCount, 1)
		}
	}

}

func (mlg *myDNSLoaderGenerator) Stop() bool {
	if !atomic.CompareAndSwapUint32(
		&mlg.status, STATUS_STARTED, STATUS_STOPPING) {
		return false
	}
	mlg.cancelFunc()
	for {
		if atomic.LoadUint32(&mlg.status) == STATUS_STOPPED {
			break
		}
		time.Sleep(time.Microsecond)
	}
	return true
}

func (mlg *myDNSLoaderGenerator) Status() uint32 {
	return atomic.LoadUint32(&mlg.status)
}
func (mlg *myDNSLoaderGenerator) CallCount() uint64 {
	return atomic.LoadUint64(&mlg.callCount)
}

// NewDNSLoaderGenerator will return a new instance of generator
// using param from GeneratorParam
func NewDNSLoaderGenerator(param GeneratorParam) (Generator, error) {
	log.Println("New Load Generator")
	if err := param.ValidCheck(); err != nil {
		return nil, err
	}
	mlg := &myDNSLoaderGenerator{
		caller:   param.Caller,
		timeout:  param.Timeout,
		qps:      param.QPS,
		duration: param.Duration,
		status:   STATUS_ORIGINAL,
	}
	mlg.init()
	return mlg, nil
}

// GenTrafficFromConfig function will do traffic generate job
// from configuration
func GenTrafficFromConfig(config *Configuration) {
	dnsclient, err := NewDNSClientWithConfig(config)
	if err != nil {
		log.Panicf("%s", err.Error())
	}
	log.Println("config the dns loader success")
	log.Printf("current configuration for dns loader is server:%s|port:%d\n",
		dnsclient.Config.Server, dnsclient.Config.Port)
	// log.Printf("%+v", config)
	param := GeneratorParam{
		Caller:   dnsclient,
		Timeout:  1000 * time.Millisecond,
		QPS:      uint32(config.QPS),
		Duration: time.Second * time.Duration(config.Duration),
	}
	log.Printf("initialize load %+v", param)
	gen, err := NewDNSLoaderGenerator(param)
	if err != nil {
		log.Panicf("load generator initialization fail :%s", err)
	}
	log.Println("start load generator")
	GloablGenerator = gen
	gen.Start()
}
