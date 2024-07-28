package organization

import (
	"context"
	"fmt"
	"sync"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
)

type Service struct {
	cfg            config.Config
	orgRepository  ports.OrganizationRepository
	userRepository ports.UserRepository
	hasher         ports.Hasher
	mailer         ports.Mailer
}

type NewServiceParams struct {
	config.Config
	ports.OrganizationRepository
	ports.UserRepository
	ports.Hasher
	ports.Mailer
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		cfg:            params.Config,
		orgRepository:  params.OrganizationRepository,
		userRepository: params.UserRepository,
		hasher:         params.Hasher,
		mailer:         params.Mailer,
	}
}

func (s *Service) CreateOrganization(dto CreateOrganizationDto) (*CreateOrganizationResponse, *domain.Error) {
	existingOrg, err := s.orgRepository.GetOwnerOrgByName(context.Background(), dto.OwnerPkID, dto.Name)
	if err != nil && err.Error != domain.NotFoundErr {
		return nil, err
	}

	if existingOrg != nil {
		return nil, domain.ErrExistOwnerOrg(dto.Name)
	}

	org, err := s.orgRepository.CreateOrg(context.Background(), dto.OwnerPkID, dto.Name, dto.Description, dto.Avatar)
	if err != nil {
		return nil, err
	}

	return &CreateOrganizationResponse{
		Org: org,
	}, nil
}

func (s *Service) GetOrganizationDetailBySlug(slug string) (*domain.Organization, *domain.Error) {
	org, err := s.orgRepository.GetOrgBySlug(context.Background(), slug)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func (s *Service) GetJoinedOrgs(userPkID int64) ([]*domain.Organization, *domain.Error) {
	orgs, err := s.orgRepository.GetOrgsByUserPkID(context.Background(), userPkID)
	if err != nil {
		return nil, err
	}
	return orgs, nil
}

func (s *Service) InviteMemberByEmails(dto InviteMemberByEmailsDto) (*InviteMemberByEmailsResponse, *domain.Error) {
	_, err := s.orgRepository.GetOwnerOrgByPkId(context.Background(), dto.OwnerPkId, dto.OrgInfo.PkId)
	if err != nil {
		return nil, err
	}

	var sentEmails []string
	var failedEmails []string

	var wg sync.WaitGroup

	for _, info := range dto.InviteInfos {
		wg.Add(1)
		go func(info EmailInviteInfo) {
			defer wg.Done()

			existingMember, _ := s.orgRepository.GetOrgMemberByEmail(context.Background(), dto.OrgInfo.PkId, info.Email)
			if existingMember != nil {
				return
			}

			err := s.mailer.SendMail(ports.SendSendGridMailPayload{
				FromName:   dto.OwnerFullName + " via Stuhub",
				ToName:     "",
				ToAddress:  info.Email,
				TemplateId: s.cfg.SendgridSetPasswordTemplateId,
				Data: map[string]string{
					"url":        "",
					"owner_name": dto.OwnerFullName,
					"org_name":   dto.OrgInfo.Name,
					"org_avatar": dto.OrgInfo.Avatar,
				},
				Subject: "Authenticate your email",
			})
			if err != nil {
				failedEmails = append(failedEmails, info.Email)
				fmt.Printf("Err sending invitation for email: %s", info.Email)
				return
			}

			var memberUserPkID *int64
			memberUser, err := s.userRepository.GetUserByEmail(context.Background(), info.Email)
			if err != nil && err.Error == domain.NotFoundErr {
				salt := s.hasher.GenerateSalt()
				newUser, err := s.userRepository.GetOrCreateUserByEmail(context.Background(), info.Email, salt)
				if err != nil {
					fmt.Printf("Failed to create new user: %s", info.Email)
					return
				}

				memberUserPkID = &newUser.PkID
			}

			if memberUser != nil {
				memberUserPkID = &memberUser.PkID
			}

			_, err = s.orgRepository.AddMemberToOrg(context.Background(), dto.OrgInfo.PkId, memberUserPkID, info.Role)
			if err != nil {
				fmt.Printf("Err add member to org: %s", info.Email)
				return
			}

			sentEmails = append(sentEmails, info.Email)
		}(info)
	}

	wg.Wait()

	return &InviteMemberByEmailsResponse{
		SentEmails:   sentEmails,
		FailedEmails: failedEmails,
	}, nil
}

func (s *Service) MakeValidateInvitationURL(token string) string {
	baseUrl := s.cfg.RemoteBaseURL

	return baseUrl + "?token=" + token
}
