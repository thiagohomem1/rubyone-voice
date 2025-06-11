package database

import (
	"context"
	"fmt"
	"log"
	"time"
	"rubyone-voice/config"
	"rubyone-voice/models"
	
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	cfg := config.GetConfig()
	
	// Configurar DSN para PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)
	
	// Configurar logger baseado no ambiente
	var logLevel logger.LogLevel
	switch cfg.Env {
	case "production":
		logLevel = logger.Error
	case "staging":
		logLevel = logger.Warn
	default:
		logLevel = logger.Info
	}
	
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}
	
	var err error
	
	// Tentar conectar com retry logic para ambientes containerizados
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err == nil {
			break
		}
		
		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		
		if i < maxRetries-1 {
			backoffDuration := time.Duration(i+1) * time.Second
			if backoffDuration > 10*time.Second {
				backoffDuration = 10 * time.Second
			}
			time.Sleep(backoffDuration)
		}
	}
	
	if err != nil {
		log.Fatal("Failed to connect to database after all retries:", err)
	}
	
	// Configurar pool de conexões para produção
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB:", err)
	}
	
	// Configurações de pool otimizadas para produção
	if cfg.Env == "production" {
		sqlDB.SetMaxIdleConns(25)
		sqlDB.SetMaxOpenConns(200)
		sqlDB.SetConnMaxLifetime(2 * time.Hour)
		sqlDB.SetConnMaxIdleTime(15 * time.Minute)
	} else {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(50)
		sqlDB.SetConnMaxLifetime(time.Hour)
		sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	}
	
	// Verificar se a conexão está funcionando
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	
	log.Printf("Connected to PostgreSQL database successfully (host: %s, db: %s, env: %s)", 
		cfg.DBHost, cfg.DBName, cfg.Env)
	
	// Auto migrate all models
	if err := migrateModels(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	
	log.Println("Database migration completed successfully")
}

func migrateModels() error {
	return DB.AutoMigrate(
		&models.Tenant{},
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.User{},
		&models.UserTenant{},
		&models.UserRole{},
		&models.Call{},
	)
}

func GetDB() *gorm.DB {
	if DB == nil {
		Connect()
	}
	return DB
}

func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func GracefulShutdown(ctx context.Context) error {
	if DB == nil {
		return nil
	}
	
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	
	// Canal para sinalizar que o close foi completado
	done := make(chan error, 1)
	
	go func() {
		done <- sqlDB.Close()
	}()
	
	select {
	case err := <-done:
		if err != nil {
			log.Printf("Error during database close: %v", err)
			return err
		}
		log.Println("Database connection closed gracefully")
		return nil
	case <-ctx.Done():
		log.Println("Database close operation timed out")
		return ctx.Err()
	}
}

func IsHealthy() bool {
	if DB == nil {
		return false
	}
	
	sqlDB, err := DB.DB()
	if err != nil {
		return false
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	return sqlDB.PingContext(ctx) == nil
}

func GetDBStats() map[string]interface{} {
	if DB == nil {
		return map[string]interface{}{
			"status": "disconnected",
		}
	}
	
	sqlDB, err := DB.DB()
	if err != nil {
		return map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		}
	}
	
	stats := sqlDB.Stats()
	return map[string]interface{}{
		"status":           "connected",
		"max_open_conns":   stats.MaxOpenConnections,
		"open_conns":       stats.OpenConnections,
		"in_use":           stats.InUse,
		"idle":             stats.Idle,
		"wait_count":       stats.WaitCount,
		"wait_duration":    stats.WaitDuration.String(),
		"max_idle_closed":  stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
} 