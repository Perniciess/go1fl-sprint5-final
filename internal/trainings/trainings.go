package trainings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

type Training struct {
	Steps        int
	TrainingType string
	Duration     time.Duration
	personaldata.Personal
}

var (
	ErrParseSteps          = errors.New("ошибка во время парсинга шагов")
	ErrParseDuration       = errors.New("ошибка во время парсинга длительности")
	ErrInvalidFormat       = errors.New("неверный формат данных")
	ErrInvalidTrainingType = errors.New("неизвестный тип тренировки")
)

func (t *Training) Parse(datastring string) (err error) {
	parts := strings.Split(datastring, ",")
	if len(parts) != 3 {
		return ErrInvalidFormat
	}

	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return ErrParseSteps
	}
	if steps <= 0 {
		return ErrParseSteps
	}

	duration, err := time.ParseDuration(strings.TrimSpace(parts[2]))
	if err != nil {
		return ErrParseDuration
	}
	if duration <= 0 {
		return ErrParseDuration
	}

	t.Steps = steps
	t.TrainingType = strings.TrimSpace(parts[1])
	t.Duration = duration
	return nil
}

func (t Training) ActionInfo() (string, error) {
	var err error
	var calories float64

	switch t.TrainingType {
	case "Ходьба":
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	case "Бег":
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	default:
		return "", ErrInvalidTrainingType
	}

	if err != nil {
		return "", err
	}

	distance := spentenergy.Distance(t.Steps, t.Height)
	meanSpeed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		t.TrainingType,
		t.Duration.Hours(),
		distance,
		meanSpeed,
		calories,
	), nil
}
