package organization

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	"github.com/Stuhub-io/utils/userutils"
)

type Service struct {
	cfg             config.Config
	orgRepository   ports.OrganizationRepository
	userRepository  ports.UserRepository
	tokenMaker      ports.TokenMaker
	hasher          ports.Hasher
	mailer          ports.Mailer
	remoteRoute     ports.RemoteRoute
	spaceRepository ports.SpaceRepository
}

type NewServiceParams struct {
	config.Config
	ports.OrganizationRepository
	ports.UserRepository
	ports.TokenMaker
	ports.Hasher
	ports.Mailer
	ports.RemoteRoute
	ports.SpaceRepository
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		cfg:             params.Config,
		orgRepository:   params.OrganizationRepository,
		userRepository:  params.UserRepository,
		tokenMaker:      params.TokenMaker,
		hasher:          params.Hasher,
		mailer:          params.Mailer,
		remoteRoute:     params.RemoteRoute,
		spaceRepository: params.SpaceRepository,
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
	_, err = s.spaceRepository.CreateSpace(context.Background(), org.PkID, dto.OwnerPkID, true, "Privte Space", "")
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
	org, err := s.orgRepository.GetOwnerOrgByPkID(context.Background(), dto.Owner.PkID, dto.OrgInfo.PkID)
	if err != nil {
		return nil, err
	}

	fmt.Print("\n\n", dto.Owner.PkID, org.OwnerID, "\n\n")

	if dto.Owner.PkID != org.OwnerID {
		return nil, domain.ErrUnauthorized
	}

	ownerFullName := userutils.GetUserFullName(dto.Owner.FirstName, dto.Owner.LastName)

	var sentEmails []string
	var failedEmails []string

	var wg sync.WaitGroup

	for _, info := range dto.InviteInfos {
		wg.Add(1)
		go func(info EmailInviteInfo) {
			defer wg.Done()

			existingMember, _ := s.orgRepository.GetOrgMemberByEmail(context.Background(), dto.OrgInfo.PkID, info.Email)
			if existingMember != nil {
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

			_, err = s.orgRepository.AddMemberToOrg(context.Background(), dto.OrgInfo.PkID, memberUserPkID, info.Role)
			if err != nil {
				fmt.Printf("Err add member to org: %s", info.Email)
				return
			}

			token, errToken := s.tokenMaker.CreateOrgInviteToken(*memberUserPkID, org.PkID, domain.OrgInvitationVerificationTokenDuration)
			if errToken != nil {
				fmt.Printf("Err create token url for: %s", info.Email)
				return
			}

			err = s.mailer.SendMail(ports.SendSendGridMailPayload{
				FromName:   ownerFullName + " via Stuhub",
				ToName:     "",
				ToAddress:  info.Email,
				TemplateId: s.cfg.SendgridOrgInvitationTemplateId,
				Data: map[string]string{
					"url":        s.MakeValidateInvitationURL(token, dto.OrgInfo.Slug),
					"owner_name": ownerFullName,
					"org_name":   dto.OrgInfo.Name,
					"org_avatar": dto.OrgInfo.Avatar,
				},
				Subject: "Accept organization invitation",
			})
			if err != nil {
				failedEmails = append(failedEmails, info.Email)
				fmt.Printf("Err sending invitation for email: %s", info.Email)
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

func (s *Service) ValidateOrgInviteToken(dto ValidateOrgInviteTokenDto) (*domain.OrganizationMember, *domain.Error) {
	payload, derr := s.tokenMaker.DecodeOrgInviteToken(dto.Token)
	if derr != nil {
		return nil, domain.ErrTokenExpired
	}

	if payload.UserPkID != dto.CurrentUser.PkID {
		return nil, domain.ErrUnauthorized
	}

	activatedMember, err := s.ActivateMember(ActivateMemberDto{
		MemberPkID: payload.UserPkID,
		OrgPkID:    payload.OrgPkID,
	})
	if err != nil {
		return nil, err
	}

	return activatedMember, nil
}

func (s *Service) ActivateMember(dto ActivateMemberDto) (*domain.OrganizationMember, *domain.Error) {
	member, err := s.orgRepository.GetOrgMemberByUserPkID(context.Background(), dto.OrgPkID, dto.MemberPkID)
	if err != nil {
		return nil, err
	}

	if member.ActivatedAt != "" {
		return member, nil
	}

	updatedMember, err := s.orgRepository.SetOrgMemberActivatedAt(context.Background(), dto.MemberPkID, time.Now())
	if err != nil {
		return nil, err
	}

	return updatedMember, nil
}

func (s *Service) MakeValidateInvitationURL(token, slug string) string {
	baseUrl := s.cfg.RemoteBaseURL + s.remoteRoute.ValidateOrgInvitation(slug)

	return baseUrl + "?token=" + token
}
