package main

import (
    "log"
    "os"
    "strconv"
    "sync"
    "time"

    "digitalcalc/internal/agent"
    "go.uber.org/zap"
)

func main() {
    logger, err := zap.NewProduction()
    if err != nil {
        log.Fatal("Failed to initialize logger", err)
    }
    defer logger.Sync()

    computingPower, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
    if err != nil || computingPower <= 0 {
        computingPower = 2 
        logger.Warn("COMPUTING_POWER not set or invalid, using default value", zap.Int("default", computingPower))
    }

    var wg sync.WaitGroup
    for i := 0; i < computingPower; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            logger.Info("Agent worker started", zap.Int("id", id))
            for {
                agent.Work(logger)
                time.Sleep(1 * time.Second) 
            }
        }(i)
    }
    wg.Wait()
}
