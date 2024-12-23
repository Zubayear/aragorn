package main

import (
	"context"
	"github.com/Zubayear/aragorn/api/handlers"
	"github.com/Zubayear/aragorn/api/routes"
	"github.com/joho/godotenv"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"time"

	"github.com/Zubayear/aragorn/pkg/sms"
	"github.com/Zubayear/aragorn/pkg/user"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var logger *zap.Logger

// NewProductionLogger creates and returns a new zap.Logger configured for production use with optional verbosity.
// It writes logs to both a rolling log file and stdout, supports caller information, and stack tracing for errors.
func NewProductionLogger(verbose bool) (*zap.Logger, error) {
	// Encoder config for JSON output
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Set log level based on verbosity
	logLevel := zapcore.InfoLevel
	if verbose {
		logLevel = zapcore.DebugLevel
	}

	// Configure lumberjack for file rolling
	logWriter := &lumberjack.Logger{
		Filename:   "logs/app.log", // Log file name
		MaxSize:    1024,           // Maximum size in MB before it gets rolled
		MaxBackups: 3,              // Maximum number of backup files
		MaxAge:     28,             // Maximum number of days to retain old log files
		Compress:   true,           // Compress the rolled files
	}

	// Create a core that writes to the file with rolling and to stdout
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(logWriter),     // Output to lumberjack
			zap.NewAtomicLevelAt(logLevel), // Log level
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(
				encoderConfig,
			), // Console output can also be JSON or you can change it to ConsoleEncoder
			zapcore.AddSync(os.Stdout),     // Output to stdout
			zap.NewAtomicLevelAt(logLevel), // Log level
		),
	)

	// Add options to log caller information
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger, nil
}

// InitializeDBPool establishes a connection pool to a PostgreSQL database using the provided connection string.
// It configures the connection pool settings for maximum/minimum connections, idle time, lifetime, and health checks.
// It terminates the process if the pool creation or connection test fails.
// Returns the created *pgxpool.Pool instance.
func InitializeDBPool(connString string, logger *zap.Logger) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		logger.Fatal("Failed to parse database connection string", zap.Error(err))
	}

	// Connection pool settings
	config.MaxConns = 50                       // Maximum number of connections in the pool
	config.MinConns = 10                       // Minimum number of idle connections
	config.MaxConnIdleTime = 5 * time.Minute   // Maximum idle time for a connection
	config.MaxConnLifetime = 30 * time.Minute  // Maximum lifetime of a connection
	config.HealthCheckPeriod = 1 * time.Minute // Interval for health checks

	// Create the connection pool
	dbPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Fatal("Failed to create database connection pool", zap.Error(err))
	}

	// Test the connection pool
	err = dbPool.Ping(context.Background())
	if err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}

	logger.Info("Database connection pool successfully created")
	return dbPool
}

// InitializeRedisClient initializes a new Redis client using the given connection string and returns the client instance.
func InitializeRedisClient(connString string) *redis.Client {
	//return redis.NewClient(&redis.Options{
	//	Addr:     connString,
	//	Password: "",
	//	DB:       0,
	//})
	return redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName: "mymaster",
		SentinelAddrs: []string{
			"redis-node-0.redis-headless.mobireach-staging.svc.cluster.local:26379",
			"redis-node-1.redis-headless.mobireach-staging.svc.cluster.local:26379",
			"redis-node-2.redis-headless.mobireach-staging.svc.cluster.local:26379",
		},
		Password:    "uaS1eegaih3AeYoh",
		PoolSize:    10,
		PoolTimeout: 10 * time.Second,
	})
}

func main() {
	_ = godotenv.Load()
	logger, _ = NewProductionLogger(false)
	dbConnString := os.Getenv("DB_CONN_STRING")
	pool := InitializeDBPool(dbConnString, logger)
	defer pool.Close()

	redisConnString := os.Getenv("REDIS_CONN_STRING")
	redisClient := InitializeRedisClient(redisConnString)
	defer redisClient.Close()

	userRepo := user.NewUserRepository(pool, redisClient)
	userService := user.NewService(userRepo)
	smsRepo := sms.NewRepository(pool)
	smsService := sms.NewService(smsRepo)
	cpHandler := handlers.NewCpHandler(userRepo)

	r := routes.Router(cpHandler, userService, smsService)

	if err := http.ListenAndServe(":42069", r); err != nil {
		logger.Fatal("Error in ListenAndServe: %s", zap.Error(err))
	}
}
