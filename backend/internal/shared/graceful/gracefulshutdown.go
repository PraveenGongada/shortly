/*
 * Copyright 2025 Praveen Kumar
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package graceful

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/PraveenGongada/shortly/internal/domain/shared/logger"
)

type (
	Operation       func(ctx context.Context) error
	ServerOperation func() error
)

func GracefulShutdown(
	serverOp ServerOperation,
	shutdownTimeout time.Duration,
	operations map[string]Operation,
	log logger.Logger,
) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)

	go func() {
		log.Info(context.Background(), "Starting server...")
		if err := serverOp(); err != nil {
			errChan <- err
		}
	}()

	// Wait for either server error or shutdown signal
	var oscall os.Signal
	select {
	case err := <-errChan:
		log.Error(context.Background(), "Server error", logger.Error(err))
		return
	case oscall = <-signalChan:
		log.Info(context.Background(), "Received system call", logger.Any("signal", oscall))
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	timeAfterExecuted := time.AfterFunc(shutdownTimeout, func() {
		log.Warn(context.Background(), "Force shutdown due to timeout")
		os.Exit(1)
	})
	defer timeAfterExecuted.Stop()

	if len(operations) > 0 {
		log.Info(context.Background(), "Executing shutdown operations...")
		wg := sync.WaitGroup{}
		wg.Add(len(operations))
		for k, op := range operations {
			go func(k string, op Operation) {
				defer wg.Done()
				log.Info(ctx, "Shutting down component", logger.String("component", k))
				if err := op(ctx); err != nil {
					log.Error(ctx, "Error shutting down component",
						logger.String("component", k),
						logger.Error(err))
				}
			}(k, op)
		}
		wg.Wait()
	}

	log.Info(context.Background(), "Graceful shutdown completed")
}
