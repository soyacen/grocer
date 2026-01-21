package pyroscopex

import (
	"log"
	"os"
	"runtime"

	"github.com/grafana/pyroscope-go"
)

func Init() {
	// These 2 lines are only required if you're using mutex or block profiling
	// Read the explanation below for how to set these rates:

	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)
	// 初始化 Pyroscope
	stopProfiler, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: os.Getenv("APP_NAME"),
		ServerAddress:   os.Getenv("PYROSCOPE_SERVER"),
		AuthToken:       os.Getenv("PYROSCOPE_AUTH_TOKEN"),

		// 采样率
		SampleRate: 100, // Hz

		// 标签
		Tags: map[string]string{
			"hostname": os.Getenv("HOSTNAME"),
			"env":      os.Getenv("ENV"),
			"version":  os.Getenv("VERSION"),
			"region":   os.Getenv("REGION"),
		},

		// 分析类型
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},

		// 日志级别
		Logger: pyroscope.StandardLogger,

		// 自动检测配置
		DisableGCRuns: false,
	})

	if err != nil {
		log.Printf("Failed to start profiler: %v", err)
	} else {
		defer stopProfiler()
	}

	// 你的应用代码
	runApplication()
}
