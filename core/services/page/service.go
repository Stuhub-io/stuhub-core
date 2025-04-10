package page

import (
	"context"
	"fmt"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	"github.com/Stuhub-io/logger"
	commonutils "github.com/Stuhub-io/utils"
	"github.com/Stuhub-io/utils/activityutils"
	"github.com/Stuhub-io/utils/pageutils"
	"github.com/Stuhub-io/utils/userutils"
)

type Service struct {
	cfg                     config.Config
	logger                  logger.Logger
	pageRepository          ports.PageRepository
	pageAccessLogRepository ports.PageAccessLogRepository
	orgRepository           ports.OrganizationRepository
	activityRepository      ports.ActivityRepository
	mailer                  ports.Mailer
}

type NewServiceParams struct {
	config.Config
	logger.Logger
	ports.PageRepository
	ports.PageAccessLogRepository
	ports.OrganizationRepository
	ports.ActivityRepository
	ports.Mailer
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		cfg:                     params.Config,
		logger:                  params.Logger,
		pageRepository:          params.PageRepository,
		pageAccessLogRepository: params.PageAccessLogRepository,
		mailer:                  params.Mailer,
		orgRepository:           params.OrganizationRepository,
		activityRepository:      params.ActivityRepository,
	}
}

func (s *Service) GetPagesByOrgPkID(
	query domain.PageListQuery,
	curUser *domain.User,
) (d []domain.Page, e *domain.Error) {

	parentPagePkID := query.ParentPagePkID

	if parentPagePkID != nil {
		parentPage, err := s.pageRepository.GetByID(
			context.Background(),
			"",
			parentPagePkID,
			domain.PageDetailOptions{},
			nil,
		)
		if err != nil {
			return nil, err
		}

		pageRole := s.GetPageRolesByUser(context.Background(), *parentPagePkID, curUser)

		permissions := s.pageRepository.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
			Page:     *parentPage,
			User:     curUser,
			PageRole: pageRole,
		})

		if !permissions.CanView {
			return nil, domain.ErrPermissionDenied
		}
	}

	d, e = s.pageRepository.List(context.Background(), query, curUser)
	return d, e
}

func (s *Service) UpdatePageByPkID(
	pagePkID int64,
	updateInput domain.PageUpdateInput,
	user *domain.User,
) (d *domain.Page, e *domain.Error) {

	page, err := s.pageRepository.GetByID(
		context.Background(),
		"",
		&pagePkID,
		domain.PageDetailOptions{},
		nil,
	)
	if err != nil {
		return nil, err
	}

	pageRole := s.GetPageRolesByUser(context.Background(), pagePkID, user)
	permissions := s.pageRepository.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *page,
		User:     user,
		PageRole: pageRole,
	})

	if !permissions.CanEdit {
		return nil, domain.ErrPermissionDenied
	}

	if page.ViewType != domain.PageViewTypeFolder {
		go s.pageAccessLogRepository.Upsert(
			context.Background(),
			pagePkID,
			user.PkID,
			domain.PageEdit,
		)
	}

	d, e = s.pageRepository.Update(context.Background(), pagePkID, updateInput)

	// Log Activity
	// FIXME: Move rename to separate API
	// go func() {
	// 	commonutils.RetryFunc(3, func() error {
	// 		metadata := commonutils.ToJsonStr(activityutils.UserUpdatePageInfoMeta{
	// 			OldPageName:  page.Name,
	// 			OldPageCover: page.CoverImage,
	// 			OldViewType:  page.ViewType.String(),
	// 		})

	// 		_, err := s.activityRepository.Create(context.Background(), domain.ActivityInput{
	// 			ActionCode: domain.ActionUserUpdatePageInfo,
	// 			PagePkID:   &page.PkID,
	// 			ActorPkID:  user.PkID,
	// 			MetaData:   &metadata,
	// 		})
	// 		if err != nil {
	// 			e := fmt.Errorf(err.Message)
	// 			s.logger.Error(e, "[Activity]: Failed to log activity for update page info")
	// 			return e
	// 		}
	// 		return nil
	// 	})
	// }()

	return d, e
}

func (s *Service) GetPageDetailByID(
	pageID string,
	publicTokenID string,
	curUser *domain.User,
) (d *domain.Page, e *domain.Error) {

	var pagePkID *int64

	// FIXME: remove public token features
	if pageID == "" {
		token, err := s.pageRepository.GetPublicTokenByID(context.Background(), publicTokenID)
		if token.ArchivedAt != "" {
			return nil, domain.NewErr("Public page is expired", domain.ResourceInvalidOrExpiredCode)
		}
		if err != nil {
			return nil, domain.ErrDatabaseQuery
		}
		pagePkID = &token.PagePkID
	}

	var userPkID *int64 = nil
	if curUser != nil {
		userPkID = &curUser.PkID
	}

	d, e = s.pageRepository.GetByID(
		context.Background(),
		pageID,
		pagePkID,
		domain.PageDetailOptions{
			Asset:    true,
			Document: true,
			Author:   true,
		},
		userPkID,
	)

	if e != nil {
		return d, e
	}

	curRole := s.GetPageRolesByUser(context.Background(), d.PkID, curUser)
	permission := s.pageRepository.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *d,
		User:     curUser,
		PageRole: curRole,
	})

	// Assign Current User Permission
	d.Permissions = &permission
	if !permission.CanView {
		return nil, domain.ErrPermissionDenied
	}

	if d.ViewType != domain.PageViewTypeFolder && curUser != nil {
		go s.pageAccessLogRepository.Upsert(
			context.Background(),
			d.PkID,
			curUser.PkID,
			domain.PageOpen,
		)
	}

	// Include Parent Page Detail
	var parentPage *domain.Page
	parentPkIDs := pageutils.PagePathToPkIDs(d.Path)

	if len(parentPkIDs) > 0 {
		parentPagePkID := parentPkIDs[len(parentPkIDs)-1]
		parentPage, e = s.pageRepository.GetByID(context.Background(), "", &parentPagePkID, domain.PageDetailOptions{}, nil)
		if e != nil {
			return d, e
		}
		parentPageCurRole := s.GetPageRolesByUser(context.Background(), parentPage.PkID, curUser)
		parentPagePermission := s.pageRepository.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
			Page:     *parentPage,
			User:     curUser,
			PageRole: parentPageCurRole,
		})
		// Assign Parent Page Permission
		parentPage.Permissions = &parentPagePermission
	}

	d.ParentPage = parentPage

	return d, e
}

func (s *Service) ArchivedPageByPkID(
	pagePkID int64,
	curUser *domain.User,
) (d *domain.Page, e *domain.Error) {
	// Recursive archive all children
	page, err := s.pageRepository.GetByID(
		context.Background(),
		"",
		&pagePkID,
		domain.PageDetailOptions{},
		nil,
	)
	if err != nil {
		return nil, err
	}

	curRole := s.GetPageRolesByUser(context.Background(), pagePkID, curUser)
	permissions := s.pageRepository.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *page,
		User:     curUser,
		PageRole: curRole,
	})

	if !permissions.CanDelete {
		return nil, domain.ErrPermissionDenied
	}

	d, e = s.pageRepository.Archive(context.Background(), pagePkID)

	var pP *domain.Page = nil
	if page.ParentPagePkID != nil {
		p, pErr := s.pageRepository.GetByID(context.Background(), "", page.ParentPagePkID, domain.PageDetailOptions{}, nil)
		if pErr != nil {
			return d, e
		}
		pP = p
	}

	go func() {
		commonutils.RetryFunc(3, func() error {
			pPName := ""
			if pP != nil {
				pPName = pP.Name
			}
			metadata := commonutils.ToJsonStr(activityutils.UserRemovePageMeta{
				OldParentPagePkID: page.ParentPagePkID,
				OldParentPageName: &pPName,
			})

			_, err := s.activityRepository.Create(context.Background(), domain.ActivityInput{
				ActionCode: domain.ActionUserRemovePage,
				PagePkID:   &page.PkID,
				OrgPkID:    &page.OrganizationPkID,
				ActorPkID:  curUser.PkID,
				MetaData:   &metadata,
			})

			if err != nil {
				e := fmt.Errorf(err.Message)
				s.logger.Error(e, "[Activity]: Failed to log activity for remove page")
				return e
			}
			return nil
		})
	}()

	return d, e
}

func (s *Service) MovePageByPkID(
	pagePkID int64,
	moveInput domain.PageMoveInput,
	curUser *domain.User,
) (d *domain.Page, e *domain.Error) {

	// Check Permission
	p, err := s.pageRepository.GetByID(
		context.Background(),
		"",
		&pagePkID,
		domain.PageDetailOptions{},
		nil,
	)
	if err != nil {
		return nil, err
	}

	curRole := s.GetPageRolesByUser(context.Background(), pagePkID, curUser)
	permissions := s.pageRepository.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *p,
		User:     curUser,
		PageRole: curRole,
	})

	if !permissions.CanMove {
		return nil, domain.ErrPermissionDenied
	}

	d, e = s.pageRepository.Move(context.Background(), pagePkID, moveInput.ParentPagePkID)

	parentPage, err := s.pageRepository.GetByID(context.Background(), "", moveInput.ParentPagePkID, domain.PageDetailOptions{}, nil)
	if err != nil {
		return d, e
	}
	// Log Activity
	go func() {
		commonutils.RetryFunc(3, func() error {
			var pName *string = nil
			if parentPage != nil {
				pName = &parentPage.Name
			}
			var oldPName *string = nil
			oldPName = &p.Name

			metadata := commonutils.ToJsonStr(activityutils.UserMovePageMeta{
				OldParentPagePkID: p.ParentPagePkID,
				NewParentPagePkID: d.ParentPagePkID,
				OldParentPageName: oldPName,
				NewParentPageName: pName,
			})
			_, err := s.activityRepository.Create(context.Background(), domain.ActivityInput{
				ActionCode: domain.ActionUserMovePage,
				PagePkID:   &d.PkID,
				OrgPkID:    &d.OrganizationPkID,
				ActorPkID:  curUser.PkID,
				MetaData:   &metadata,
			})

			if err != nil {
				e := fmt.Errorf(err.Message)
				s.logger.Error(e, "[Activity]: Failed to log activity for move page")
				return e
			}
			return nil
		})
	}()

	return d, e
}

func (s *Service) CreatePublicPageToken(
	pageID string,
) (d *domain.PagePublicToken, e *domain.Error) {
	page, err := s.pageRepository.GetByID(
		context.Background(),
		pageID,
		nil,
		domain.PageDetailOptions{},
		nil,
	)
	if err != nil {
		return nil, domain.ErrDatabaseQuery
	}
	d, e = s.pageRepository.CreatePublicToken(context.Background(), page.PkID)
	return d, e
}

func (s *Service) ArchiveAllPublicPageToken(pageID string) (e *domain.Error) {
	page, err := s.pageRepository.GetByID(
		context.Background(),
		pageID,
		nil,
		domain.PageDetailOptions{},
		nil,
	)
	if err != nil {
		return domain.ErrDatabaseQuery
	}
	e = s.pageRepository.ArchiveAllPublicToken(context.Background(), page.PkID)
	return e
}

func (s *Service) UpdateGeneralAccess(
	pagePkID int64,
	updateInput domain.PageGeneralAccessUpdateInput,
	curUser *domain.User,
) (*domain.Page, *domain.Error) {

	page, err := s.pageRepository.GetByID(
		context.Background(),
		"",
		&pagePkID,
		domain.PageDetailOptions{},
		nil,
	)

	if err != nil {
		return nil, err
	}

	curRole := s.GetPageRolesByUser(context.Background(), pagePkID, curUser)
	permission := s.pageRepository.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *page,
		User:     curUser,
		PageRole: curRole,
	})

	if !permission.CanShare {
		return nil, domain.ErrPermissionDenied
	}

	page, err = s.pageRepository.UpdateGeneralAccess(context.Background(), pagePkID, updateInput)
	if err != nil {
		return nil, err
	}

	return page, nil
}

// Document Controller.
func (s *Service) CreateDocumentPage(
	pageInput domain.DocumentPageInput,
	curUser *domain.User,
) (*domain.Page, *domain.Error) {
	parentPagePkID := pageInput.ParentPagePkID
	var parentPage *domain.Page

	// Check Parent Page Permission
	if parentPagePkID != nil {
		parent, err := s.pageRepository.GetByID(
			context.Background(),
			"",
			parentPagePkID,
			domain.PageDetailOptions{},
			nil,
		)
		if err != nil {
			return nil, err
		}
		parentPage = parent

		curRole := s.GetPageRolesByUser(context.Background(), *pageInput.ParentPagePkID, curUser)
		permission := s.pageRepository.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
			Page:     *parentPage,
			User:     curUser,
			PageRole: curRole,
		})

		if !permission.CanEdit {
			return nil, domain.ErrPermissionDenied
		}
	}

	if curUser == nil {
		return nil, domain.ErrPermissionDenied
	}

	// FIXME: Check if user is a member of the organization
	page, err := s.pageRepository.CreateDocumentPage(context.Background(), pageInput)

	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}
	// Log Activity
	go func() {
		commonutils.RetryFunc(3, func() error {
			var pName *string = nil
			if parentPage != nil {
				pName = &parentPage.Name
			}

			metadata := commonutils.ToJsonStr(activityutils.UserCreatePageMeta{
				ParentPagePkID: pageInput.ParentPagePkID,
				ParentPageName: pName,
				NewPageName:    page.Name,
				NewPagePkID:    page.PkID,
				NewPageID:      page.ID,
			})

			_, err := s.activityRepository.Create(context.Background(), domain.ActivityInput{
				ActionCode: domain.ActionUserCreatePage,
				PagePkID:   &page.PkID,
				OrgPkID:    &page.OrganizationPkID,
				ActorPkID:  curUser.PkID,
				MetaData:   &metadata,
			})
			if err != nil {
				e := fmt.Errorf(err.Message)
				s.logger.Error(e, "Failed to log activity")
				return e
			}
			return nil
		})
	}()

	return page, nil
}

func (s *Service) UpdateDocumentContentByPkID(
	pagePkID int64,
	content domain.DocumentInput,
	curUser *domain.User,
) (d *domain.Page, e *domain.Error) {

	page, err := s.pageRepository.GetByID(
		context.Background(),
		"",
		&pagePkID,
		domain.PageDetailOptions{},
		nil,
	)

	if err != nil {
		return nil, err
	}

	curRole := s.GetPageRolesByUser(context.Background(), pagePkID, curUser)
	permissions := s.pageRepository.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *page,
		User:     curUser,
		PageRole: curRole,
	})

	if !permissions.CanEdit {
		return nil, domain.ErrPermissionDenied
	}

	if page.ViewType != domain.PageViewTypeFolder {
		go s.pageAccessLogRepository.Upsert(
			context.Background(),
			page.PkID,
			curUser.PkID,
			domain.PageEdit,
		)
	}

	// Activity Log
	d, e = s.pageRepository.UpdateContent(context.Background(), pagePkID, content)
	go func() {
		commonutils.RetryFunc(3, func() error {
			metadata := commonutils.ToJsonStr(activityutils.UserVisitePageMeta{})

			_, err := s.activityRepository.Create(context.Background(), domain.ActivityInput{
				ActionCode: domain.ActionUserCreatePage,
				PagePkID:   &page.PkID,
				ActorPkID:  curUser.PkID,
				MetaData:   &metadata,
			})
			if err != nil {
				e := fmt.Errorf(err.Message)
				s.logger.Error(e, "Failed to log activity")
				return e
			}
			return nil
		})
	}()

	return d, e
}

// Generate Public Token.
func (s *Service) GenerateDocumentPublicToken(pagePkID int64) (d string, e *domain.Error) {
	// d, e = s.pageRepository.GeneratePublicToken(context.Background(), pagePkID)
	return d, e
}

func (s *Service) ValidateDocumentPublicToken(token string) (d bool, e *domain.Error) {
	// d, e = s.pageRepository.ValidatePublicToken(context.Background(), token)
	return d, e
}

// Asset Controller.
func (s *Service) CreateAssetPage(
	assetInput domain.AssetPageInput,
	curUser *domain.User,
) (*domain.Page, *domain.Error) {

	parentPagePkID := assetInput.ParentPagePkID
	var parentPage *domain.Page
	if parentPagePkID != nil {
		parent, err := s.pageRepository.GetByID(
			context.Background(),
			"",
			parentPagePkID,
			domain.PageDetailOptions{},
			nil,
		)
		parentPage = parent
		if err != nil {
			return nil, err
		}

		curRole := s.GetPageRolesByUser(context.Background(), *parentPagePkID, curUser)
		permission := s.pageRepository.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
			Page:     *parentPage,
			User:     curUser,
			PageRole: curRole,
		})

		if !permission.CanEdit {
			return nil, domain.ErrPermissionDenied
		}
	}

	//FIXME: GetOrgMembers Error
	// members, err := s.orgRepository.GetOrgMembers(context.Background(), assetInput.OrganizationPkID)
	// if err != nil {
	// 	return nil, err
	// }
	// isOrgMember := sliceutils.Find(members, func(member domain.OrganizationMember) bool {
	// 	return member.OrganizationPkID == assetInput.OrganizationPkID
	// }) != nil

	// if !isOrgMember {
	// 	return nil, domain.ErrPermissionDenied
	// }

	page, err := s.pageRepository.CreateAsset(context.Background(), assetInput)
	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}
	// Log Activity
	go func() {
		commonutils.RetryFunc(3, func() error {
			var pName *string = nil
			if parentPage != nil {
				pName = &parentPage.Name
			}

			metadata := commonutils.ToJsonStr(activityutils.UserCreatePageMeta{
				ParentPagePkID: assetInput.ParentPagePkID,
				ParentPageName: pName,
				NewPageName:    page.Name,
				NewPagePkID:    page.PkID,
				NewPageID:      page.ID,
			})

			_, err := s.activityRepository.Create(context.Background(), domain.ActivityInput{
				ActionCode: domain.ActionUserCreatePage,
				PagePkID:   &page.PkID,
				OrgPkID:    &page.OrganizationPkID,
				ActorPkID:  curUser.PkID,
				MetaData:   &metadata,
			})
			if err != nil {
				e := fmt.Errorf(err.Message)
				s.logger.Error(e, "Failed to log activity")
				return e
			}
			return nil
		})
	}()

	go s.pageAccessLogRepository.Upsert(
		context.Background(),
		page.PkID,
		curUser.PkID,
		domain.PageUpload,
	)

	return page, nil
}

// Page Role

func (s *Service) AddPageRoleUser(
	input domain.PageRoleCreateInput,
	curUser *domain.User,
) (*domain.PageRoleUser, *domain.Page, *domain.Error) {
	if curUser == nil {
		return nil, nil, domain.ErrPermissionDenied
	}

	existingPage, err := s.pageRepository.GetByID(
		context.Background(),
		"",
		&input.PagePkID,
		domain.PageDetailOptions{
			Organization: true,
		},
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	curRole := s.GetPageRolesByUser(context.Background(), input.PagePkID, curUser)
	permissions := s.pageRepository.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *existingPage,
		User:     curUser,
		PageRole: curRole,
	})

	if !permissions.CanShare {
		return nil, nil, domain.ErrPermissionDenied
	}

	exisingPageRoleUser, _ := s.pageRepository.GetPageRoleByEmail(
		context.Background(),
		input.PagePkID,
		input.Email,
	)

	if exisingPageRoleUser != nil {
		return nil, nil, domain.ErrExisitingPageRoleUser
	}

	pageRoleUser, err := s.pageRepository.CreatePageRole(context.Background(), input)
	if err != nil {
		return nil, nil, domain.ErrDatabaseMutation
	}

	err = s.mailer.SendMailCustomTemplate(ports.SendSendGridMailCustomTemplatePayload{
		ToName: userutils.GetUserFullName(
			pageRoleUser.User.FirstName,
			pageRoleUser.User.LastName,
		),
		ToAddress:        pageRoleUser.User.Email,
		TemplateHTMLName: "share_people",
		Data: map[string]string{
			"sender": userutils.GetUserFullName(
				curUser.FirstName,
				curUser.LastName,
			),
			"url": fmt.Sprintf("%s/%s/%s", s.cfg.RemoteBaseURL, existingPage.Organization.Slug, existingPage.ID),
		},
		Subject: "Share with you",
	})
	if err != nil {
		return nil, nil, err
	}

	return pageRoleUser, existingPage, nil
}

func (s *Service) GetPageRoleUsers(
	input domain.PageRoleGetAllInput,
	curUser *domain.User,
) ([]domain.PageRoleUser, *domain.Error) {

	pagePkID := input.PagePkID

	page, err := s.pageRepository.GetByID(
		context.Background(),
		"",
		&pagePkID,
		domain.PageDetailOptions{},
		nil,
	)
	if err != nil {
		return nil, err
	}

	curRole := s.GetPageRolesByUser(context.Background(), input.PagePkID, curUser)
	permissions := s.pageRepository.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *page,
		User:     curUser,
		PageRole: curRole,
	})

	if !permissions.CanShare {
		return nil, domain.ErrPermissionDenied
	}

	pageRoleUsers, err := s.pageRepository.GetPageRoles(
		context.Background(),
		input.PagePkID,
	)

	if err != nil {
		return nil, err
	}

	return pageRoleUsers, nil
}

func (s *Service) UpdatePageRoleUser(
	input domain.PageRoleUpdateInput,
	curUser *domain.User,
) *domain.Error {
	exisingPage, err := s.pageRepository.GetByID(
		context.Background(),
		"",
		&input.PagePkID,
		domain.PageDetailOptions{},
		nil,
	)

	curRole := s.GetPageRolesByUser(context.Background(), input.PagePkID, curUser)
	permissions := s.pageRepository.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *exisingPage,
		User:     curUser,
		PageRole: curRole,
	})

	if !permissions.CanShare {
		return domain.ErrPermissionDenied
	}

	if err != nil {
		return err
	}

	exisingPageRoleUser, _ := s.pageRepository.GetPageRoleByEmail(
		context.Background(),
		input.PagePkID,
		input.Email,
	)
	if exisingPageRoleUser == nil {
		return domain.ErrNotFound
	}

	return s.pageRepository.UpdatePageRole(context.Background(), input)
}

func (s *Service) DeletePageRoleUser(
	input domain.PageRoleDeleteInput,
	curUser *domain.User,
) *domain.Error {
	existingPage, err := s.pageRepository.GetByID(
		context.Background(),
		"",
		&input.PagePkID,
		domain.PageDetailOptions{},
		nil,
	)
	if err != nil {
		return err
	}

	curRole := s.GetPageRolesByUser(context.Background(), input.PagePkID, curUser)
	permissions := s.pageRepository.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *existingPage,
		User:     curUser,
		PageRole: curRole,
	})

	if !permissions.CanShare {
		return domain.ErrPermissionDenied
	}

	exisingPageRoleUser, _ := s.pageRepository.GetPageRoleByEmail(
		context.Background(),
		input.PagePkID,
		input.Email,
	)
	if exisingPageRoleUser == nil {
		return domain.ErrNotFound
	}

	return s.pageRepository.DeletePageRole(context.Background(), input)
}

func (s *Service) GetPageRolesByUser(ctx context.Context, pagePkID int64, user *domain.User) *domain.PageRole {
	if user == nil {
		return nil
	}
	role, err := s.pageRepository.GetPageRoleByEmail(ctx, pagePkID, user.Email)
	if err != nil {
		return nil
	}
	return &role.Role
}

func (s Service) RequestPagePermission(pageID string, email string) *domain.Error {

	page, pErr := s.pageRepository.GetByID(context.Background(), pageID, nil, domain.PageDetailOptions{}, nil)
	if pErr != nil {
		return pErr
	}

	_, err := s.pageRepository.CreatePageAccessRequest(context.Background(), domain.PageRoleRequestCreateInput{
		PagePkID: page.PkID,
		Email:    email,
	})

	return err
}

func (s Service) ListRequestPagePermissions(pagePkID int64) ([]domain.PageRoleRequestLog, *domain.Error) {
	return s.pageRepository.ListPageAccessRequestByPagePkID(context.Background(), domain.PageRoleRequestLogQuery{
		PagePkIDs: []int64{pagePkID},
		Status:    []domain.PageRoleRequestLogStatus{domain.PRSLPending},
	})
}

func (s Service) RejectPagePermissions(pagePkID int64, emails []string) *domain.Error {
	existingPage, err := s.pageRepository.GetByID(
		context.Background(),
		"",
		&pagePkID,
		domain.PageDetailOptions{},
		nil,
	)
	if err != nil {
		return err
	}

	err = s.pageRepository.UpdatePageAccessRequestStatus(context.Background(), domain.PageRoleRequestLogQuery{
		PagePkIDs: []int64{pagePkID},
		Emails:    emails,
	}, domain.PRSLRejected)

	if err != nil {
		return err
	}

	for _, email := range emails {
		go func(email string) {
			err = s.mailer.SendMailCustomTemplate(ports.SendSendGridMailCustomTemplatePayload{
				ToName:           email,
				ToAddress:        email,
				TemplateHTMLName: "share_request_rejected",
				Data: map[string]string{
					"page": existingPage.Name,
				},
				Subject: "Access request reply",
			})
			if err != nil {
				s.logger.Info(err.Message)
			}
		}(email)
	}

	return nil
}

func (s Service) AcceptRequestPagePermission(input domain.PageRoleCreateInput, curUser *domain.User) *domain.Error {
	_, pageDetails, err := s.AddPageRoleUser(input, curUser)
	if err != nil {
		return err
	}

	err = s.pageRepository.UpdatePageAccessRequestStatus(context.Background(), domain.PageRoleRequestLogQuery{
		PagePkIDs: []int64{input.PagePkID},
		Emails:    []string{input.Email},
	}, domain.PRSLApproved)

	if err != nil {
		return err
	}

	err = s.mailer.SendMailCustomTemplate(ports.SendSendGridMailCustomTemplatePayload{
		ToName: userutils.GetUserFullName(
			curUser.FirstName,
			curUser.LastName,
		),
		ToAddress:        input.Email,
		TemplateHTMLName: "share_request_accepted",
		Data: map[string]string{
			"page": pageDetails.Name,
			"sender": userutils.GetUserFullName(
				curUser.FirstName,
				curUser.LastName,
			),
			"url": fmt.Sprintf("%s/%s/%s", s.cfg.RemoteBaseURL, pageDetails.Organization.Slug, pageDetails.ID),
		},
		Subject: "Access request reply",
	})
	if err != nil {
		return err
	}

	return nil
}

func (s Service) AddPageToStarred(input domain.StarPageInput, curUser *domain.User) *domain.Error {
	// Handler Permissions
	page, pErr := s.pageRepository.GetByID(
		context.Background(),
		"",
		&input.PagePkID,
		domain.PageDetailOptions{
			Author: true,
		},
		nil,
	)
	if pErr != nil {
		return pErr
	}

	curRole := s.GetPageRolesByUser(context.Background(), input.PagePkID, curUser)
	permission := s.pageRepository.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *page,
		User:     curUser,
		PageRole: curRole,
	})

	if !permission.CanView {
		return domain.ErrPermissionDenied
	}

	_, err := s.pageRepository.StarPage(context.Background(), input)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) RemovePageFromStarred(input domain.StarPageInput, curUser *domain.User) *domain.Error {
	// Handler Permissions
	page, pErr := s.pageRepository.GetByID(
		context.Background(),
		"",
		&input.PagePkID,
		domain.PageDetailOptions{
			Author: true,
		},
		nil,
	)
	if pErr != nil {
		return pErr
	}

	curRole := s.GetPageRolesByUser(context.Background(), input.PagePkID, curUser)
	permission := s.pageRepository.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *page,
		User:     curUser,
		PageRole: curRole,
	})

	if !permission.CanView {
		return domain.ErrPermissionDenied
	}
	err := s.pageRepository.UnstarPage(context.Background(), input)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) CreateUserActivity(input domain.ActivityInput, curUser *domain.User) *domain.Error {
	// FIXME: Check Permissions

	_, err := s.activityRepository.Create(context.Background(), input)
	if err != nil {
		return err
	}
	return nil
}
