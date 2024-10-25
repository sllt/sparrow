package gen

import (
	"time"
)

var (
	DefaultRequestTimeout = 5

	DefaultCompressionType      CompressionType  = CompressionTypeGZIP
	DefaultCompressionLevel     CompressionLevel = CompressionDefault
	DefaultCompressionThreshold int              = 1024

	DefaultLogFilter = []LogLevel{
		LogLevelTrace,
		LogLevelDebug,
		LogLevelInfo,
		LogLevelWarning,
		LogLevelError,
		LogLevelPanic,
	}

	DefaultKeepAlivePeriod time.Duration = 15 * time.Second
	DefaultTCPBufferSize   int           = 65535
	DefaultPort            uint16        = 11144

	DefaultNetworkFlags = NetworkFlags{
		Enable:                       true,
		EnableRemoteSpawn:            true,
		EnableRemoteApplicationStart: true,
		EnableFragmentation:          false,
		EnableProxyTransit:           false,
		EnableProxyAccept:            true,
		EnableImportantDelivery:      true,
	}

	DefaultNetworkProxyFlags = NetworkProxyFlags{
		Enable:                       true,
		EnableRemoteSpawn:            false,
		EnableRemoteApplicationStart: false,
		EnableEncryption:             false,
		EnableImportantDelivery:      true,
	}

	DefaultLogLevels = []LogLevel{
		LogLevelSystem,
		LogLevelTrace,
		LogLevelDebug,
		LogLevelInfo,
		LogLevelWarning,
		LogLevelError,
		LogLevelPanic,
	}
)
