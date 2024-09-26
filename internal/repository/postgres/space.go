package postgres

import (
	"context"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/internal/repository/model"
)

type SpaceRepository struct {
	store *store.DBStore
	cfg   config.Config
}

type NewSpaceRepositoryParams struct {
	Cfg   config.Config
	Store *store.DBStore
}

func NewSpaceRepository(params NewSpaceRepositoryParams) ports.SpaceRepository {
	return &SpaceRepository{
		store: params.Store,
		cfg:   params.Cfg,
	}
}

func (r *SpaceRepository) CreateSpace(ctx context.Context, orgPkID int64, ownerPkID int64, isPrivate bool, name, description string) (*domain.Space, *domain.Error) {
	var newSpace model.Space
	var newSpaceMember model.SpaceMember

	// -- Start Tx
	tx, doneTx := r.store.NewTransaction()
	newSpace = model.Space{
		OrgPkid:     orgPkID,
		Name:        name,
		Description: description,
		IsPrivate:   isPrivate,
	}
	err := tx.DB().Create(&newSpace).Error
	if err != nil {
		return nil, doneTx(err)
	}
	newSpaceMember = model.SpaceMember{
		SpacePkid: newSpace.Pkid,
		UserPkid:  ownerPkID,
		Role:      domain.Owner.String(),
	}
	err = tx.DB().Create(&newSpaceMember).Error
	if err != nil {
		return nil, doneTx(err)
	}
	commitTxErr := doneTx(nil)
	if commitTxErr != nil {
		return nil, commitTxErr
	}
	// -- End Tx
	return &domain.Space{
		PkID:        newSpace.Pkid,
		ID:          newSpace.ID,
		Name:        newSpace.Name,
		OrgPkID:     newSpace.OrgPkid,
		Description: newSpace.Description,
		IsPrivate:   newSpace.IsPrivate,
		CreatedAt:   newSpace.CreatedAt.String(),
		UpdatedAt:   newSpace.UpdatedAt.String(),
	}, nil
}

func (r *SpaceRepository) GetSpacesByOrgPkID(ctx context.Context, orgPkID int64) ([]domain.Space, *domain.Error) {
	type SpaceMemberWithUser struct {
		model.SpaceMember
		User model.User `gorm:"foreignKey:user_pkid" json:"user"`
	}

	type SpaceWithMembers struct {
		model.Space
		Members []SpaceMemberWithUser `gorm:"foreignKey:space_pkid" json:"members"`
	}
	var spaces []SpaceWithMembers
	err := r.store.DB().Preload("Members").Preload("Members.User").Where("org_pkid = ?", orgPkID).Find(&spaces).Error
	if err != nil {
		return nil, domain.ErrDatabaseQuery
	}
	domainSpaces := make([]domain.Space, 0, len(spaces))
	for _, space := range spaces {

		domainSpaceMembers := make([]domain.SpaceMember, 0, len(space.Members))
		for _, member := range space.Members {
			domainSpaceMember := domain.SpaceMember{
				PkID:      member.Pkid,
				SpacePkID: member.SpacePkid,
				UserPkID:  member.UserPkid,
				Role:      member.Role,
				CreatedAt: member.CreatedAt.String(),
				UpdatedAt: member.UpdatedAt.String(),
				User:      mapUserModelToDomain(member.User),
			}
			domainSpaceMembers = append(domainSpaceMembers, domainSpaceMember)
		}

		domainSpace := domain.Space{
			PkID:        space.Pkid,
			ID:          space.ID,
			Name:        space.Name,
			OrgPkID:     space.OrgPkid,
			Description: space.Description,
			IsPrivate:   space.IsPrivate,
			CreatedAt:   space.CreatedAt.String(),
			UpdatedAt:   space.UpdatedAt.String(),
			Members:     domainSpaceMembers,
		}
		domainSpaces = append(domainSpaces, domainSpace)
	}
	return domainSpaces, nil
}
