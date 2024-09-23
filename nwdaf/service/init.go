package service

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/ciromacedo/nwdaf/analyticsinfo"
	"github.com/ciromacedo/nwdaf/datacollection"
	"github.com/ciromacedo/nwdaf/eventssubscription"
	mongoDBLibLogger "github.com/free5gc/MongoDBLibrary/logger"
	openApiLogger "github.com/free5gc/openapi/logger"
	pathUtilLogger "github.com/free5gc/path_util/logger"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/ciromacedo/nwdaf/consumer"
	nwdaf_context "github.com/ciromacedo/nwdaf/context"
	"github.com/ciromacedo/nwdaf/factory"
	"github.com/ciromacedo/nwdaf/logger"
	"github.com/ciromacedo/nwdaf/util"
	"github.com/free5gc/MongoDBLibrary"
	"github.com/free5gc/http2_util"
	"github.com/free5gc/logger_util"
	"github.com/free5gc/openapi/models"
)

type timerFunc func()

type NWDAF struct{}

type (
	Config struct {
		nwdafcfg string
	}
)

var config Config

//var nwdafcontext nwdaf_context.NWDAFContext

var nwdafCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "nwdafcfg",
		Usage: "config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}
func settimer() {
	timer2 := time.NewTimer(2 * time.Second)

	<-timer2.C
	fmt.Println("second timer started")
}

func timerStart(callbackFunc timerFunc, value int) {
	timer1 := time.NewTimer(time.Duration(value) * time.Second)
	<-timer1.C
	fmt.Println("timer has started")
	callbackFunc()
}

var nfid string

func TimerCallback_() {
	//var nfid string
	self := nwdaf_context.NWDAF_Self()
	util.InitNwdafContext(self)
	//profile := consumer.BuildNFInstance(self)
	//nfid = profile.NfInstanceId
	patchItem := []models.PatchItem{
		{
			Op:     "replace",
			Path:   "/nfStatus",
			From:   "NWDAf",
			Value:  "REGISTERED",
			Scheme: models.UriScheme(self.UriScheme),
		},
	}

	fmt.Println("callback called...")
	DoneAsync(TimerCallback_)
	consumer.SendNFPeriodicHeartbeat(self.NrfUri, nfid, patchItem)
	fmt.Println("it is running")
}
func DoneAsync(callbackFunc timerFunc) {
	//r := make(chan int)
	fmt.Println("Warming up ...")
	go func() {
		fmt.Println("Done123 ...")
		time.Sleep(10 * time.Second)
		//r <- 1
		callbackFunc()
		fmt.Println("Done ...")
	}()
	//return r
}

func (*NWDAF) GetCliCmd() (flags []cli.Flag) {
	return nwdafCLi
}

func (nwdaf *NWDAF) Initialize(c *cli.Context) error {
	config = Config{
		nwdafcfg: c.String("nwdafcfg"),
	}

	if config.nwdafcfg != "" {
		if err := factory.InitConfigFactory(config.nwdafcfg); err != nil {
			return err
		}
	} else {
		//DefaultAmfConfigPath := path_util.Free5gcPath("./config/nwdafcfg.yaml")

		if err := factory.InitConfigFactory(util.DefaultNwdafConfigPath); err != nil {
			return err
		}

		//if err := factory.InitConfigFactory(DefaultAmfConfigPath); err != nil {
		//	return err
		//}
	}

	nwdaf.setLogLevel()

	if err := factory.CheckConfigVersion(); err != nil {
		return err
	}

	return nil
}

func (nwdaf *NWDAF) setLogLevel() {
	if factory.NwdafConfig.Logger == nil {
		initLog.Warnln("NWDAF config without log level setting!!!")
		return
	}

	if factory.NwdafConfig.Logger.UDR != nil {
		if factory.NwdafConfig.Logger.UDR.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.NwdafConfig.Logger.UDR.DebugLevel); err != nil {
				initLog.Warnf("UDR Log level [%s] is invalid, set to [info] level",
					factory.NwdafConfig.Logger.UDR.DebugLevel)
				logger.SetLogLevel(logrus.InfoLevel)
			} else {
				initLog.Infof("UDR Log level is set to [%s] level", level)
				logger.SetLogLevel(level)
			}
		} else {
			initLog.Infoln("UDR Log level not set. Default set to [info] level")
			logger.SetLogLevel(logrus.InfoLevel)
		}
		logger.SetReportCaller(factory.NwdafConfig.Logger.UDR.ReportCaller)
	}

	if factory.NwdafConfig.Logger.PathUtil != nil {
		if factory.NwdafConfig.Logger.PathUtil.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.NwdafConfig.Logger.PathUtil.DebugLevel); err != nil {
				pathUtilLogger.PathLog.Warnf("PathUtil Log level [%s] is invalid, set to [info] level",
					factory.NwdafConfig.Logger.PathUtil.DebugLevel)
				pathUtilLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				pathUtilLogger.SetLogLevel(level)
			}
		} else {
			pathUtilLogger.PathLog.Warnln("PathUtil Log level not set. Default set to [info] level")
			pathUtilLogger.SetLogLevel(logrus.InfoLevel)
		}
		pathUtilLogger.SetReportCaller(factory.NwdafConfig.Logger.PathUtil.ReportCaller)
	}

	if factory.NwdafConfig.Logger.OpenApi != nil {
		if factory.NwdafConfig.Logger.OpenApi.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.NwdafConfig.Logger.OpenApi.DebugLevel); err != nil {
				openApiLogger.OpenApiLog.Warnf("OpenAPI Log level [%s] is invalid, set to [info] level",
					factory.NwdafConfig.Logger.OpenApi.DebugLevel)
				openApiLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				openApiLogger.SetLogLevel(level)
			}
		} else {
			openApiLogger.OpenApiLog.Warnln("OpenAPI Log level not set. Default set to [info] level")
			openApiLogger.SetLogLevel(logrus.InfoLevel)
		}
		openApiLogger.SetReportCaller(factory.NwdafConfig.Logger.OpenApi.ReportCaller)
	}

	if factory.NwdafConfig.Logger.MongoDBLibrary != nil {
		if factory.NwdafConfig.Logger.MongoDBLibrary.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.NwdafConfig.Logger.MongoDBLibrary.DebugLevel); err != nil {
				mongoDBLibLogger.MongoDBLog.Warnf("MongoDBLibrary Log level [%s] is invalid, set to [info] level",
					factory.NwdafConfig.Logger.MongoDBLibrary.DebugLevel)
				mongoDBLibLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				mongoDBLibLogger.SetLogLevel(level)
			}
		} else {
			mongoDBLibLogger.MongoDBLog.Warnln("MongoDBLibrary Log level not set. Default set to [info] level")
			mongoDBLibLogger.SetLogLevel(logrus.InfoLevel)
		}
		mongoDBLibLogger.SetReportCaller(factory.NwdafConfig.Logger.MongoDBLibrary.ReportCaller)
	}
}

func (nwdaf *NWDAF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range nwdaf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (nwdaf *NWDAF) Start() {
	// get config file info
	config := factory.NwdafConfig
	mongodb := config.Configuration.Mongodb

	initLog.Infof("NWDAF Config Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)

	// Connect to MongoDB
	MongoDBLibrary.SetMongoDB(mongodb.Name, mongodb.Url)

	initLog.Infoln("Server started")

	router := logger_util.NewGinWithLogrus(logger.GinLog)
	//router1 := logger_util.NewGinWithLogrus(logger.GinLog)

	// Order is important for the same route pattern.
	//datarepository.AddService(router)
	analyticsinfo.AddService(router)
	eventssubscription.AddService(router)
	datacollection.AddService(router)

	// analyticsinfo.AddService(router1)
	// eventssubscription.AddService(router1)
	// datacollection.AddService(router1)

	nwdafLogPath := util.NwdafLogPath
	nwdafPemPath := util.NwdafPemPath
	nwdafKeyPath := util.NwdafKeyPath

	self := nwdaf_context.NWDAF_Self()
	util.InitNwdafContext(self)
	//nwdafcontext = self

	addr := fmt.Sprintf("%s:%d", self.BindingIPv4, self.SBIPort)
	initLog.Errorf("***********[%s:%d:%s]", self.BindingIPv4, self.SBIPort, addr)
	// addr1 := fmt.Sprintf("%s:%d", "127.0.0.59", 29599)
	// initLog.Errorf("***********[%s:%d:%s]", self.BindingIPv4, self.SBIPort, addr1)

	profile := consumer.BuildNFInstance(self)
	nfid = profile.NfInstanceId
	var newNrfUri string
	var err error
	newNrfUri, self.NfId, err = consumer.SendRegisterNFInstance(self.NrfUri, profile.NfInstanceId, profile)
	if err == nil {
		//		self.NrfUri = newNrfUri
		newNrfUri = newNrfUri
		//start timer
		//timerStart(settimer, 4)
		fmt.Println("Let's start ...")
		//val := DoneAsync()
		DoneAsync(TimerCallback_)
		//fmt.Println("Done is running ...")

		//fmt.Println(<- val)
	} else {
		initLog.Errorf("Send Register NFInstance Error[%s]", err.Error())
	}

	fmt.Println("running ...")
	server, err := http2_util.NewServer(addr, nwdafLogPath, router)
	fmt.Println("second running ...")
	if server == nil {
		initLog.Errorf("Initialize HTTP server failed: %+v", err)
		return
	}

	if err != nil {
		initLog.Warnf("Initialize HTTP server: %+v", err)
	}

	/*server1, err1 := http2_util.NewServer(addr1, nwdafLogPath, router1)
	if server1 == nil {
		initLog.Errorf("Initialize HTTP server failed: %+v", err)
		return
	}

	if err1 != nil {
		initLog.Warnf("Initialize HTTP server: %+v", err)
	}*/

	/* init subscriber data collect */
	datacollection.InitEventExposureSubscriber(self)

	serverScheme := factory.NwdafConfig.Configuration.Sbi.Scheme
	if serverScheme == "http" {
		err = server.ListenAndServe()
		// consumer.SendNFPeriodicHeartbeat(self.NrfUri, profile.NfInstanceId, patchItem)
		// time.Sleep(10 * time.Second)
		// continue

	} else if serverScheme == "https" {
		err = server.ListenAndServeTLS(nwdafPemPath, nwdafKeyPath)
	}

	if err != nil {
		initLog.Fatalf("HTTP server setup failed: %+v", err)
	}

	/*serverScheme1 := factory.NwdafConfig.Configuration.Sbi.Scheme
	if serverScheme1 == "http" {
		err = server1.ListenAndServe()
	} else if serverScheme1 == "https" {
		err = server1.ListenAndServeTLS(nwdafPemPath, nwdafKeyPath)
	}*/

	if err != nil {
		initLog.Fatalf("HTTP server setup failed: %+v", err)
	}

}

func (nwdaf *NWDAF) Exec(c *cli.Context) error {

	//NWDAF.Initialize(cfgPath, c)

	initLog.Traceln("args:", c.String("nwdafcfg"))
	args := nwdaf.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./nwdaf", args...)

	nwdaf.Initialize(c)

	var stdout io.ReadCloser
	if readCloser, err := command.StdoutPipe(); err != nil {
		initLog.Fatalln(err)
	} else {
		stdout = readCloser
	}
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	var stderr io.ReadCloser
	if readCloser, err := command.StderrPipe(); err != nil {
		initLog.Fatalln(err)
	} else {
		stderr = readCloser
	}
	go func() {
		in := bufio.NewScanner(stderr)
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	var err error
	go func() {
		if errormessage := command.Start(); err != nil {
			fmt.Println("command.Start Fails!")
			err = errormessage
		}
		wg.Done()
	}()

	wg.Wait()
	return err
}
