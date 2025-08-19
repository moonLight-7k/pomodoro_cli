package config

import (
	"os"
	"strconv"
	"time"

	"pomodoro_cli/internal/errors"
)

type Config struct {
	WorkDuration     time.Duration
	BreakDuration    time.Duration
	ProgressBarWidth int
	MaxSessionTime   time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		WorkDuration:     25 * time.Minute,
		BreakDuration:    5 * time.Minute,
		ProgressBarWidth: 30,
		MaxSessionTime:   12 * time.Hour, // Reasonable maximum
	}
}

func ParseArgs() (*Config, error) {
	config := DefaultConfig()

	if len(os.Args) < 3 || len(os.Args) > 4 {
		return nil, &errors.AppError{
			Code:    errors.ErrInvalidArgs,
			Message: "Usage: pomodoro <work_time> <break_time> [-h]",
			Details: map[string]interface{}{
				"examples": []string{
					"pomodoro 25 5      # 25 minutes work, 5 minutes break",
					"pomodoro 1 1 -h    # 1 hour work, 1 hour break",
				},
			},
		}
	}

	useHours := false
	if len(os.Args) == 4 {
		if os.Args[3] == "-h" {
			useHours = true
		} else {
			return nil, &errors.AppError{
				Code:    errors.ErrInvalidFlag,
				Message: "Unknown flag: " + os.Args[3],
				Details: map[string]interface{}{
					"valid_flags": []string{"-h"},
				},
			}
		}
	}

	workTime, err := parseTimeArg(os.Args[1], "work_time")
	if err != nil {
		return nil, err
	}

	breakTime, err := parseTimeArg(os.Args[2], "break_time")
	if err != nil {
		return nil, err
	}

	var workDuration, breakDuration time.Duration
	if useHours {
		workDuration = time.Duration(workTime) * time.Hour
		breakDuration = time.Duration(breakTime) * time.Hour
	} else {
		workDuration = time.Duration(workTime) * time.Minute
		breakDuration = time.Duration(breakTime) * time.Minute
	}

	if workDuration > config.MaxSessionTime {
		return nil, &errors.AppError{
			Code:    errors.ErrInvalidDuration,
			Message: "Work session too long",
			Details: map[string]interface{}{
				"max_allowed": config.MaxSessionTime.String(),
				"requested":   workDuration.String(),
			},
		}
	}

	if breakDuration > config.MaxSessionTime {
		return nil, &errors.AppError{
			Code:    errors.ErrInvalidDuration,
			Message: "Break session too long",
			Details: map[string]interface{}{
				"max_allowed": config.MaxSessionTime.String(),
				"requested":   breakDuration.String(),
			},
		}
	}

	config.WorkDuration = workDuration
	config.BreakDuration = breakDuration

	return config, nil
}

func parseTimeArg(arg, name string) (int, error) {
	value, err := strconv.Atoi(arg)
	if err != nil {
		return 0, &errors.AppError{
			Code:    errors.ErrInvalidNumber,
			Message: name + " must be a valid integer",
			Details: map[string]interface{}{
				"provided": arg,
			},
		}
	}

	if value <= 0 {
		return 0, &errors.AppError{
			Code:    errors.ErrInvalidNumber,
			Message: name + " must be positive",
			Details: map[string]interface{}{
				"provided": value,
			},
		}
	}

	const maxValue = 999
	if value > maxValue {
		return 0, &errors.AppError{
			Code:    errors.ErrInvalidNumber,
			Message: name + " is too large",
			Details: map[string]interface{}{
				"provided":    value,
				"max_allowed": maxValue,
			},
		}
	}

	return value, nil
}
