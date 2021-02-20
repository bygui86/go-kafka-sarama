package config

import (
	"github.com/bygui86/go-kafka/tracing/kafka-producer/logging"
	"github.com/bygui86/go-kafka/tracing/kafka-producer/utils"
)

const (
	enableMonitoringEnvVar = "ENABLE_MONITORING" // bool
	enableTracingEnvVar    = "ENABLE_TRACING"    // bool
	tracingTechEnvVar      = "TRACING_TECH"      //  available values: jaeger, zipkin

	enableMonitoringDefault = true
	enableTracingDefault    = true
	tracingTechDefault      = TracingTechJaeger
)

func LoadConfig() *Config {
	logging.Log.Info("Load global configurations")

	tracingTech := utils.GetStringEnv(tracingTechEnvVar, tracingTechDefault)
	if tracingTech != TracingTechJaeger && tracingTech != TracingTechZipkin {
		logging.SugaredLog.Warnf("Tracing technology %s not supported, fallback to %s",
			tracingTech, TracingTechJaeger)
		tracingTech = TracingTechJaeger
	}

	return &Config{
		enableMonitoring: utils.GetBoolEnv(enableMonitoringEnvVar, enableMonitoringDefault),
		enableTracing:    utils.GetBoolEnv(enableTracingEnvVar, enableTracingDefault),
		tracingTech:      tracingTech,
	}
}
