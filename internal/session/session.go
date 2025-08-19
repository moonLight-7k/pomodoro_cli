package session

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pomodoro_cli/internal/config"
	"pomodoro_cli/internal/errors"
	"pomodoro_cli/internal/terminal"
)

// SessionType represents the type of session
type SessionType int

const (
	WorkSession SessionType = iota
	BreakSession
)

func (st SessionType) String() string {
	switch st {
	case WorkSession:
		return "Work"
	case BreakSession:
		return "Break"
	default:
		return "Unknown"
	}
}

type Session struct {
	Type      SessionType
	Duration  time.Duration
	StartTime time.Time
	EndTime   time.Time
	Completed bool
	Cancelled bool
}

type SessionManager struct {
	config   *config.Config
	terminal *terminal.Terminal
	logger   *errors.Logger
	sessions []Session
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewSessionManager(cfg *config.Config, term *terminal.Terminal, logger *errors.Logger) *SessionManager {
	ctx, cancel := context.WithCancel(context.Background())

	sm := &SessionManager{
		config:   cfg,
		terminal: term,
		logger:   logger,
		sessions: make([]Session, 0),
		ctx:      ctx,
		cancel:   cancel,
	}

	sm.setupSignalHandling()

	return sm
}

func (sm *SessionManager) setupSignalHandling() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		sm.logger.LogInfo("Shutdown signal received", map[string]interface{}{
			"sessions_completed": len(sm.getCompletedSessions()),
			"total_sessions":     len(sm.sessions),
		})

		fmt.Println("\n Exiting Pomodoro.")
		sm.cancel()
		os.Exit(0)
	}()
}

func (sm *SessionManager) RunSession(sessionType SessionType) error {
	var duration time.Duration
	var label string

	switch sessionType {
	case WorkSession:
		duration = sm.config.WorkDuration
		label = "Work"
	case BreakSession:
		duration = sm.config.BreakDuration
		label = "Break"
	}

	session := Session{
		Type:      sessionType,
		Duration:  duration,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(duration),
	}

	sm.logger.LogInfo("Session started", map[string]interface{}{
		"type":     sessionType.String(),
		"duration": duration.String(),
	})

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sm.ctx.Done():
			session.Cancelled = true
			sm.sessions = append(sm.sessions, session)
			return fmt.Errorf("session cancelled")

		case now := <-ticker.C:
			remaining := time.Until(session.EndTime)

			if remaining <= 0 {
				session.Completed = true
				sm.sessions = append(sm.sessions, session)

				if err := sm.terminal.DisplayCompletion(label); err != nil {
					sm.logger.LogError(err)
				}

				sm.logger.LogInfo("Session completed", map[string]interface{}{
					"type":     sessionType.String(),
					"duration": duration.String(),
				})

				time.Sleep(2 * time.Second)
				return nil
			}

			elapsed := now.Sub(session.StartTime)
			progress := elapsed.Seconds() / duration.Seconds()

			sessionInfo := terminal.SessionInfo{
				Label:            label,
				Elapsed:          elapsed,
				Progress:         progress,
				ProgressBarWidth: sm.config.ProgressBarWidth,
			}

			if err := sm.terminal.DisplaySession(sessionInfo); err != nil {
				sm.logger.LogError(errors.NewAppError(
					errors.ErrTerminalNotSupported,
					"Failed to display session",
					map[string]interface{}{
						"error": err.Error(),
					},
				))
			}
		}
	}
}

func (sm *SessionManager) RunPomodoroCycle() error {
	cycle := 0

	for {
		select {
		case <-sm.ctx.Done():
			return fmt.Errorf("pomodoro cycle cancelled")
		default:
		}

		cycle++
		sm.logger.LogInfo("Starting pomodoro cycle", map[string]interface{}{
			"cycle": cycle,
		})

		if err := sm.RunSession(WorkSession); err != nil {
			return fmt.Errorf("work session failed: %w", err)
		}

		select {
		case <-sm.ctx.Done():
			return fmt.Errorf("pomodoro cycle cancelled")
		default:
		}

		if err := sm.RunSession(BreakSession); err != nil {
			return fmt.Errorf("break session failed: %w", err)
		}
	}
}

func (sm *SessionManager) GetStats() map[string]interface{} {
	completed := sm.getCompletedSessions()
	workSessions := 0
	breakSessions := 0
	totalWorkTime := time.Duration(0)
	totalBreakTime := time.Duration(0)

	for _, session := range completed {
		if session.Type == WorkSession {
			workSessions++
			totalWorkTime += session.Duration
		} else {
			breakSessions++
			totalBreakTime += session.Duration
		}
	}

	return map[string]interface{}{
		"total_sessions":     len(sm.sessions),
		"completed_sessions": len(completed),
		"work_sessions":      workSessions,
		"break_sessions":     breakSessions,
		"total_work_time":    totalWorkTime.String(),
		"total_break_time":   totalBreakTime.String(),
	}
}

func (sm *SessionManager) getCompletedSessions() []Session {
	var completed []Session
	for _, session := range sm.sessions {
		if session.Completed {
			completed = append(completed, session)
		}
	}
	return completed
}

func (sm *SessionManager) Close() error {
	sm.cancel()

	stats := sm.GetStats()
	sm.logger.LogInfo("Session manager closing", stats)

	return nil
}
