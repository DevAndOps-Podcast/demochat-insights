package insights

import (
	"context"
	"demochat-insights/internal/repositories/insights"
	"time"
)

type Service struct {
	insightsRepo *insights.Repository
}

func New(insightsRepo *insights.Repository) *Service {
	return &Service{insightsRepo: insightsRepo}
}

type PublishMessageRequest struct {
	UserID  int64  `json:"user_id"`
	Message string `json:"message"`
}

func (s *Service) PublishMessage(ctx context.Context, req PublishMessageRequest) error {
	tx, err := s.insightsRepo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	currentInsights, err := s.insightsRepo.GetInsights(ctx)
	if err != nil {
		return err
	}

	currentInsights.TotalMessages++

	if err := s.insightsRepo.UpdateUserActivity(ctx, tx, req.UserID); err != nil {
		return err
	}

	mostActiveUserID, err := s.insightsRepo.GetMostActiveUserID(ctx)
	if err != nil {
		return err
	}

	currentInsights.MostActiveUserID = mostActiveUserID

	now := time.Now().Unix()
	if currentInsights.FirstMessageTimestamp == nil {
		currentInsights.FirstMessageTimestamp = &now
	}
	currentInsights.LastMessageTimestamp = &now

	if *currentInsights.LastMessageTimestamp != *currentInsights.FirstMessageTimestamp {
		duration := *currentInsights.LastMessageTimestamp - *currentInsights.FirstMessageTimestamp
		currentInsights.AverageMessageRate = float64(currentInsights.TotalMessages) / float64(duration)
	}

	if err := s.insightsRepo.UpdateInsights(ctx, tx, currentInsights); err != nil {
		return err
	}

	return tx.Commit()
}

type Insights struct {
	MostActiveUser     int64   `json:"most_active_user_id"`
	TotalMessages      int64   `json:"total_messages"`
	AverageMessageRate float64 `json:"average_message_rate"`
}

func (s *Service) GetInsights(ctx context.Context) (*Insights, error) {
	i, err := s.insightsRepo.GetInsights(ctx)
	if err != nil {
		return nil, err
	}

	return &Insights{
		MostActiveUser:     i.MostActiveUserID,
		TotalMessages:      i.TotalMessages,
		AverageMessageRate: i.AverageMessageRate,
	}, nil
}
