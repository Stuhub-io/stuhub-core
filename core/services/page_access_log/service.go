package pageAccessLog

import (
	"context"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
)

type Service struct {
	pageAccessLogRepository ports.PageAccessLogRepository
}

type NewServiceParams struct {
	ports.PageAccessLogRepository
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		pageAccessLogRepository: params.PageAccessLogRepository,
	}
}

func (s *Service) GetLogsByUser(
	query domain.OffsetBasedPagination,
	userPkID int64,
) ([]domain.PageAccessLog, *domain.Error) {
	logs, err := s.pageAccessLogRepository.GetByUserPKID(context.Background(), query, userPkID)

	//TODO: check permision of parent pages

	return logs, err
}
