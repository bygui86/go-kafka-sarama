package main

import (
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/openzipkin/zipkin-go/reporter"

	"github.com/bygui86/go-kafka/tracing/kafka-consumer/commons"
	"github.com/bygui86/go-kafka/tracing/kafka-consumer/config"
	"github.com/bygui86/go-kafka/tracing/kafka-consumer/consumer"
	"github.com/bygui86/go-kafka/tracing/kafka-consumer/logging"
	"github.com/bygui86/go-kafka/tracing/kafka-consumer/monitoring"
	"github.com/bygui86/go-kafka/tracing/kafka-consumer/tracing"
)

const (
	zipkinHost = "localhost"
	zipkinPort = 9411
)

var (
	monitoringServer *monitoring.Server
	jaegerCloser     io.Closer
	zipkinReporter   reporter.Reporter
	kafkaConsumer    *consumer.KafkaConsumer
)

func main() {
	initLogging()

	logging.SugaredLog.Infof("Start %s", commons.ServiceName)

	cfg := loadConfig()

	if cfg.GetEnableMonitoring() {
		monitoringServer = startMonitoringServer()
	}

	if cfg.GetEnableTracing() {
		switch cfg.GetTracingTech() {
		case config.TracingTechJaeger:
			jaegerCloser = initJaegerTracer()
		case config.TracingTechZipkin:
			zipkinReporter = initZipkinTracer()
		}
	}

	kafkaConsumer = startConsumer()

	logging.SugaredLog.Infof("%s up and running", commons.ServiceName)

	startSysCallChannel()

	shutdownAndWait(1)
}

func initLogging() {
	err := logging.InitGlobalLogger()
	if err != nil {
		logging.SugaredLog.Errorf("Logging setup failed: %s", err.Error())
		os.Exit(501)
	}
}

func loadConfig() *config.Config {
	logging.Log.Debug("Load configurations")
	return config.LoadConfig()
}

func startMonitoringServer() *monitoring.Server {
	logging.Log.Debug("Start monitoring")
	server := monitoring.New()
	logging.Log.Debug("Monitoring server successfully created")

	server.Start()
	logging.Log.Debug("Monitoring successfully started")

	return server
}

func initJaegerTracer() io.Closer {
	logging.Log.Debug("Init Jaeger tracer")
	closer, err := tracing.InitTracer()
	if err != nil {
		logging.SugaredLog.Errorf("Jaeger tracer setup failed: %s", err.Error())
		os.Exit(501)
	}
	return closer
}

func initZipkinTracer() reporter.Reporter {
	logging.Log.Debug("Init Zipkin tracer")
	zReporter, err := tracing.InitTestingZipkin(commons.ServiceName, zipkinHost, zipkinPort)
	if err != nil {
		logging.SugaredLog.Errorf("Zipkin tracer setup failed: %s", err.Error())
		os.Exit(501)
	}
	return zReporter
}

func startConsumer() *consumer.KafkaConsumer {
	logging.Log.Debug("Start consumer")
	kConsumer, newErr := consumer.New(commons.ServiceName)
	if newErr != nil {
		logging.SugaredLog.Errorf("Consumer setup failed: %s", newErr.Error())
		os.Exit(501)
	}
	logging.Log.Debug("Consumer successfully created")

	startErr := kConsumer.Start()
	if startErr != nil {
		logging.SugaredLog.Errorf("Consumer start failed: %s", startErr.Error())
		os.Exit(502)
	}
	logging.Log.Debug("Consumer successfully started")

	return kConsumer
}

func startSysCallChannel() {
	syscallCh := make(chan os.Signal)
	signal.Notify(syscallCh, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-syscallCh
}

func shutdownAndWait(timeout int) {
	logging.SugaredLog.Warnf("Termination signal received! Timeout %d", timeout)

	if kafkaConsumer != nil {
		kafkaConsumer.Shutdown(timeout)
	}

	if jaegerCloser != nil {
		err := jaegerCloser.Close()
		if err != nil {
			logging.SugaredLog.Errorf("Jaeger tracer closure failed: %s", err.Error())
		}
	}

	if zipkinReporter != nil {
		err := zipkinReporter.Close()
		if err != nil {
			logging.SugaredLog.Errorf("Zipkin tracer closure failed: %s", err.Error())
		}
	}

	if monitoringServer != nil {
		monitoringServer.Shutdown(timeout)
	}

	time.Sleep(time.Duration(timeout+1) * time.Second)
}
