package page

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	"github.com/Stuhub-io/logger"
	commonutils "github.com/Stuhub-io/utils"
	"github.com/Stuhub-io/utils/activityutils"
	"github.com/Stuhub-io/utils/pageutils"
	"github.com/Stuhub-io/utils/timeutils"
	"github.com/Stuhub-io/utils/userutils"
)

type Service struct {
	cfg    config.Config
	logger logger.Logger
	mailer ports.Mailer
	repo   *ports.Repository
}

type NewServiceParams struct {
	config.Config
	logger.Logger
	ports.Mailer
	*ports.Repository
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		cfg:    params.Config,
		logger: params.Logger,
		mailer: params.Mailer,
		repo:   params.Repository,
	}
}

func (s *Service) GetPagesByOrgPkID(
	query domain.PageListQuery,
	curUser *domain.User,
) (d []domain.Page, e *domain.Error) {

	parentPagePkID := query.ParentPagePkID

	if parentPagePkID != nil {
		parentPage, err := s.repo.Page.GetByID(
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

		permissions := s.repo.Page.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
			Page:     *parentPage,
			User:     curUser,
			PageRole: pageRole,
		})

		if !permissions.CanView {
			return nil, domain.ErrPermissionDenied
		}
	}

	d, e = s.repo.Page.List(context.Background(), query, curUser)
	return d, e
}

func (s *Service) UpdatePageByPkID(
	pagePkID int64,
	updateInput domain.PageUpdateInput,
	user *domain.User,
) (d *domain.Page, e *domain.Error) {

	page, err := s.repo.Page.GetByID(
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
	permissions := s.repo.Page.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *page,
		User:     user,
		PageRole: pageRole,
	})

	if !permissions.CanEdit {
		return nil, domain.ErrPermissionDenied
	}

	if page.ViewType != domain.PageViewTypeFolder {
		go s.repo.PageAccessLog.Upsert(
			context.Background(),
			pagePkID,
			user.PkID,
			domain.PageEdit,
		)
	}

	d, e = s.repo.Page.Update(context.Background(), pagePkID, updateInput)

	return d, e
}

func (s *Service) RenamePageByPkID(pagePkID int64, input domain.RenamePageInput, user *domain.User) (d *domain.Page, e *domain.Error) {
	page, err := s.repo.Page.GetByID(context.Background(), "", nil, domain.PageDetailOptions{}, nil)
	if err != nil {
		return nil, err
	}
	pageRole := s.GetPageRolesByUser(context.Background(), pagePkID, user)
	permissions := s.repo.Page.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *page,
		User:     user,
		PageRole: pageRole,
	})

	if !permissions.CanEdit {
		return nil, domain.ErrPermissionDenied
	}

	d, e = s.repo.Page.Update(context.Background(), pagePkID, domain.PageUpdateInput{
		Name: &input.Name,
	})

	go func() {
		relatedPagePkIDs := []int64{pagePkID}
		snapshot := commonutils.ToJsonStr(activityutils.UserRenamePageMetaV2{
			Page:    page,
			NewName: input.Name,
		})
		_, err := s.repo.ActivityV2.Create(context.Background(), domain.ActivityV2Input{
			ActionCode:       domain.ActionUserRenamePage,
			UserPkID:         user.PkID,
			Snapshot:         snapshot,
			RelatedPagePkIDs: relatedPagePkIDs,
		})

		if err != nil {
			e := fmt.Errorf(err.Message)
			s.logger.Error(e, "[RenamePageByPkID] - activityV2Repository.Create error")
			return
		}
	}()

	return d, e
}

func (s *Service) GetPageDetailByIdOrPkID(
	pageID string,
	PkID *int64,
	curUser *domain.User,
) (d *domain.Page, e *domain.Error) {

	var userPkID *int64 = nil
	if curUser != nil {
		userPkID = &curUser.PkID
	}

	d, e = s.repo.Page.GetByID(
		context.Background(),
		pageID,
		PkID,
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
	permission := s.repo.Page.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
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
		go s.repo.PageAccessLog.Upsert(
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
		parentPage, e = s.repo.Page.GetByID(context.Background(), "", &parentPagePkID, domain.PageDetailOptions{}, nil)
		if e != nil {
			return d, e
		}
		parentPageCurRole := s.GetPageRolesByUser(context.Background(), parentPage.PkID, curUser)
		parentPagePermission := s.repo.Page.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
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

func (s *Service) GetPageUnsafe(
	pagePkID int64,
	userPkID *int64,
) (d *domain.Page, e *domain.Error) {
	// Ignore Permission Check
	d, e = s.repo.Page.GetByID(context.Background(), "", &pagePkID, domain.PageDetailOptions{
		Document: true,
		Asset:    true,
		Author:   true,
	}, userPkID)

	return d, e
}

func (s *Service) ArchivedPageByPkID(
	pagePkID int64,
	curUser *domain.User,
) (d *domain.Page, e *domain.Error) {

	page, err := s.repo.Page.GetByID(
		context.Background(),
		"",
		&pagePkID,
		domain.PageDetailOptions{
			Asset: true, // Need for display file type UI in activity
		},
		nil,
	)
	if err != nil {
		return nil, err
	}

	curRole := s.GetPageRolesByUser(context.Background(), pagePkID, curUser)
	permissions := s.repo.Page.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *page,
		User:     curUser,
		PageRole: curRole,
	})

	if !permissions.CanDelete {
		return nil, domain.ErrPermissionDenied
	}

	d, e = s.repo.Page.Archive(context.Background(), pagePkID)

	if e == nil {
		go func() {
			// Archive Activity
			var originalParentPage *domain.Page = nil
			if page.ParentPagePkID != nil {
				p, pErr := s.repo.Page.GetByID(context.Background(), "", page.ParentPagePkID, domain.PageDetailOptions{}, nil)
				if pErr != nil {
					e := fmt.Errorf(pErr.Message)
					s.logger.Error(e, "[ArchivedPageByPkID] - pageRepository.GetByID error")
					return
				}
				originalParentPage = p
			}

			relatedPagePkIDs := []int64{pagePkID}
			if originalParentPage != nil {
				relatedPagePkIDs = append(relatedPagePkIDs, originalParentPage.PkID)
			}

			snapshot := commonutils.ToJsonStr(activityutils.UserArchivePageMetaV2{
				Page:       page,
				ParentPage: originalParentPage,
			})

			_, er := s.repo.ActivityV2.Create(context.Background(), domain.ActivityV2Input{
				ActionCode:       domain.ActionUserArchivePage,
				UserPkID:         curUser.PkID,
				RelatedPagePkIDs: relatedPagePkIDs,
				Snapshot:         snapshot,
			})

			if er != nil {
				e := fmt.Errorf(er.Message)
				s.logger.Error(e, "[ArchivedPageByPkID] - activityV2Repository.Create error")
			}
		}()
	}

	return d, e
}

func (s *Service) MovePageByPkID(
	pagePkID int64,
	moveInput domain.PageMoveInput,
	curUser *domain.User,
) (d *domain.Page, e *domain.Error) {

	// Check Permission
	p, err := s.repo.Page.GetByID(
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
	permissions := s.repo.Page.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *p,
		User:     curUser,
		PageRole: curRole,
	})

	if !permissions.CanMove {
		return nil, domain.ErrPermissionDenied
	}

	d, e = s.repo.Page.Move(context.Background(), pagePkID, moveInput.ParentPagePkID)

	if e != nil {
		return d, e
	}

	// Log Activity
	go func() {
		// Include related page
		relatedPagePkIDs := []int64{pagePkID}

		srcParentPagePkID := p.ParentPagePkID
		if srcParentPagePkID != nil {
			relatedPagePkIDs = append(relatedPagePkIDs, *srcParentPagePkID)
		}

		parentPage, err := s.repo.Page.GetByID(context.Background(), "", moveInput.ParentPagePkID, domain.PageDetailOptions{}, nil)
		if err == nil && parentPage != nil {
			relatedPagePkIDs = append(relatedPagePkIDs, parentPage.PkID)
		}

		snapshot := commonutils.ToJsonStr(activityutils.UserMovePageMetaV2{
			Page:          p,
			DesParentPage: parentPage,
		})

		_, aErr := s.repo.ActivityV2.Create(context.Background(), domain.ActivityV2Input{
			ActionCode:       domain.ActionUserMovePage,
			UserPkID:         curUser.PkID,
			RelatedPagePkIDs: relatedPagePkIDs,
			Snapshot:         snapshot,
		})
		if aErr != nil {
			e := fmt.Errorf(aErr.Message)
			s.logger.Error(e, "[MovePageByPkID] - activityV2Repository.Create error")
		}
	}()

	return d, e
}

func (s *Service) CreatePublicPageToken(
	pageID string,
) (d *domain.PagePublicToken, e *domain.Error) {
	page, err := s.repo.Page.GetByID(
		context.Background(),
		pageID,
		nil,
		domain.PageDetailOptions{},
		nil,
	)
	if err != nil {
		return nil, domain.ErrDatabaseQuery
	}
	d, e = s.repo.Page.CreatePublicToken(context.Background(), page.PkID)
	return d, e
}

func (s *Service) ArchiveAllPublicPageToken(pageID string) (e *domain.Error) {
	page, err := s.repo.Page.GetByID(
		context.Background(),
		pageID,
		nil,
		domain.PageDetailOptions{},
		nil,
	)
	if err != nil {
		return domain.ErrDatabaseQuery
	}
	e = s.repo.Page.ArchiveAllPublicToken(context.Background(), page.PkID)
	return e
}

func (s *Service) UpdateGeneralAccess(
	pagePkID int64,
	updateInput domain.PageGeneralAccessUpdateInput,
	curUser *domain.User,
) (*domain.Page, *domain.Error) {

	page, err := s.repo.Page.GetByID(
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
	permission := s.repo.Page.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *page,
		User:     curUser,
		PageRole: curRole,
	})

	if !permission.CanShare {
		return nil, domain.ErrPermissionDenied
	}

	page, err = s.repo.Page.UpdateGeneralAccess(context.Background(), pagePkID, updateInput)
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
		parent, err := s.repo.Page.GetByID(
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
		permission := s.repo.Page.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
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
	page, err := s.repo.Page.CreateDocumentPage(context.Background(), pageInput)

	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}

	// Log Activity
	go func() {
		commonutils.RetryFunc(3, func() error {
			createdPage, e := s.repo.Page.GetByID(context.Background(), "", &page.PkID, domain.PageDetailOptions{
				Author: true,
			}, nil)
			if e != nil {
				return fmt.Errorf(e.Message)
			}
			relatedPagePkIDs := []int64{createdPage.PkID}
			if parentPage != nil {
				relatedPagePkIDs = append(relatedPagePkIDs, parentPage.PkID)
			}

			var err *domain.Error
			var snapshot string
			var activityCode domain.ActionCode

			if pageInput.ViewType == domain.PageViewTypeDoc {
				snapshot = commonutils.ToJsonStr(activityutils.UserCreateDocumentMeta{
					ParentPage: parentPage,
					ChildPage:  createdPage,
				})
				activityCode = domain.ActionUserCreateDocument
			} else {
				activityCode = domain.ActionUserCreateFolder
				pageRoles, pErr := s.repo.Page.GetPageRoles(context.Background(), createdPage.PkID)
				if pErr != nil {
					e := fmt.Errorf(pErr.Message)
					s.logger.Error(e, "[CreateDocumentPage] - pageRepository.GetPageRoles error")
					return e
				}
				snapshot = commonutils.ToJsonStr(activityutils.UserCreateFolderMeta{
					ParentPage: parentPage,
					ChildPage:  *createdPage,
					PageRoles:  pageRoles,
				})
			}
			_, err = s.repo.ActivityV2.Create(context.Background(), domain.ActivityV2Input{
				ActionCode:       activityCode,
				UserPkID:         curUser.PkID,
				Snapshot:         snapshot,
				RelatedPagePkIDs: relatedPagePkIDs,
			})

			if err != nil {
				e := fmt.Errorf(err.Message)
				s.logger.Error(e, "[CreateDocumentPage] - activityV2Repository.Create error")
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

	page, err := s.repo.Page.GetByID(
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
	permissions := s.repo.Page.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *page,
		User:     curUser,
		PageRole: curRole,
	})

	if !permissions.CanEdit {
		return nil, domain.ErrPermissionDenied
	}

	if page.ViewType != domain.PageViewTypeFolder {
		go s.repo.PageAccessLog.Upsert(
			context.Background(),
			page.PkID,
			curUser.PkID,
			domain.PageEdit,
		)
	}

	// Activity Log
	d, e = s.repo.Page.UpdateContent(context.Background(), pagePkID, content)

	return d, e
}

func (s *Service) ValidateDocumentPublicToken(token string) (d bool, e *domain.Error) {
	// d, e = s.repo.Page.ValidatePublicToken(context.Background(), token)
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
		parent, err := s.repo.Page.GetByID(
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
		permission := s.repo.Page.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
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

	page, err := s.repo.Page.CreateAsset(context.Background(), assetInput)
	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}
	// Log Activity
	go func() {
		// Join upload activity if latest upload activity is less than 5 minutes
		RelatedPagePkIDs := []int64{}
		if parentPage != nil {
			RelatedPagePkIDs = append(RelatedPagePkIDs, parentPage.PkID)
		}

		recentActivity, err := s.repo.ActivityV2.One(context.Background(), domain.ActivityV2ListQuery{
			ActionCodes:      []domain.ActionCode{domain.ActionUserUploadedAssets},
			UserPkIDs:        []int64{curUser.PkID},
			RelatedPagePkIDs: RelatedPagePkIDs,
		})

		if recentActivity != nil && err == nil {
			assetMeta := activityutils.UserUploadedAssetsMeta{}
			if err := json.Unmarshal([]byte(recentActivity.Snapshot), &assetMeta); err == nil {
				now := time.Now()

				recentActivityCreatedAt := timeutils.ParseTime(recentActivity.CreatedAt)

				if ((assetMeta.ParentPage == nil && parentPagePkID == nil) ||
					(assetMeta.ParentPage != nil && parentPagePkID != nil &&
						assetMeta.ParentPage.PkID == *parentPagePkID)) &&
					now.Sub(*recentActivityCreatedAt) < time.Minute*5 {
					// Join assets to existed activity
					assetMeta.Assets = append(assetMeta.Assets, *page)
					// Update activity

					_, e := s.repo.ActivityV2.Update(context.Background(), recentActivity.PkID, domain.ActivityV2Input{
						ActionCode:       domain.ActionUserUploadedAssets,
						UserPkID:         curUser.PkID,
						Snapshot:         commonutils.ToJsonStr(assetMeta),
						RelatedPagePkIDs: []int64{page.PkID}, // Add new related page
					})
					if e != nil {
						e := fmt.Errorf(e.Message)
						s.logger.Error(e, "[CreateAssetPage] - activityV2Repository.Update error")
						return
					}

					return
				}
			}
		}

		dataSnapshot := commonutils.ToJsonStr(activityutils.UserUploadedAssetsMeta{
			ParentPage: parentPage,
			Assets:     []domain.Page{*page},
		})

		relatedPagePkIDs := []int64{page.PkID}
		if parentPage != nil {
			relatedPagePkIDs = append(relatedPagePkIDs, parentPage.PkID)
		}

		_, err = s.repo.ActivityV2.Create(context.Background(), domain.ActivityV2Input{
			ActionCode:       domain.ActionUserUploadedAssets,
			UserPkID:         curUser.PkID,
			Snapshot:         dataSnapshot,
			RelatedPagePkIDs: relatedPagePkIDs,
		})
		if err != nil {
			e := fmt.Errorf(err.Message)
			s.logger.Error(e, "[CreateAssetPage] - activityV2Repository.Create error")
			return
		}
	}()

	go s.repo.PageAccessLog.Upsert(
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

	existingPage, err := s.repo.Page.GetByID(
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
	permissions := s.repo.Page.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *existingPage,
		User:     curUser,
		PageRole: curRole,
	})

	if !permissions.CanShare {
		return nil, nil, domain.ErrPermissionDenied
	}

	exisingPageRoleUser, _ := s.repo.Page.GetPageRoleByEmail(
		context.Background(),
		input.PagePkID,
		input.Email,
	)

	if exisingPageRoleUser != nil {
		return nil, nil, domain.ErrExisitingPageRoleUser
	}

	pageRoleUser, err := s.repo.Page.CreatePageRole(context.Background(), input)
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

	page, err := s.repo.Page.GetByID(
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
	permissions := s.repo.Page.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *page,
		User:     curUser,
		PageRole: curRole,
	})

	if !permissions.CanShare {
		return nil, domain.ErrPermissionDenied
	}

	pageRoleUsers, err := s.repo.Page.GetPageRoles(
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
	exisingPage, err := s.repo.Page.GetByID(
		context.Background(),
		"",
		&input.PagePkID,
		domain.PageDetailOptions{},
		nil,
	)

	curRole := s.GetPageRolesByUser(context.Background(), input.PagePkID, curUser)
	permissions := s.repo.Page.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
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

	exisingPageRoleUser, _ := s.repo.Page.GetPageRoleByEmail(
		context.Background(),
		input.PagePkID,
		input.Email,
	)
	if exisingPageRoleUser == nil {
		return domain.ErrNotFound
	}

	return s.repo.Page.UpdatePageRole(context.Background(), input)
}

func (s *Service) DeletePageRoleUser(
	input domain.PageRoleDeleteInput,
	curUser *domain.User,
) *domain.Error {
	existingPage, err := s.repo.Page.GetByID(
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
	permissions := s.repo.Page.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *existingPage,
		User:     curUser,
		PageRole: curRole,
	})

	if !permissions.CanShare {
		return domain.ErrPermissionDenied
	}

	exisingPageRoleUser, _ := s.repo.Page.GetPageRoleByEmail(
		context.Background(),
		input.PagePkID,
		input.Email,
	)
	if exisingPageRoleUser == nil {
		return domain.ErrNotFound
	}

	return s.repo.Page.DeletePageRole(context.Background(), input)
}

func (s *Service) GetPageRolesByUser(ctx context.Context, pagePkID int64, user *domain.User) *domain.PageRole {
	if user == nil {
		return nil
	}
	role, err := s.repo.Page.GetPageRoleByEmail(ctx, pagePkID, user.Email)
	if err != nil {
		return nil
	}
	return &role.Role
}

func (s Service) RequestPagePermission(pageID string, email string) *domain.Error {

	page, pErr := s.repo.Page.GetByID(context.Background(), pageID, nil, domain.PageDetailOptions{}, nil)
	if pErr != nil {
		return pErr
	}

	_, err := s.repo.Page.CreatePageAccessRequest(context.Background(), domain.PageRoleRequestCreateInput{
		PagePkID: page.PkID,
		Email:    email,
	})

	return err
}

func (s Service) ListRequestPagePermissions(pagePkID int64) ([]domain.PageRoleRequestLog, *domain.Error) {
	return s.repo.Page.ListPageAccessRequestByPagePkID(context.Background(), domain.PageRoleRequestLogQuery{
		PagePkIDs: []int64{pagePkID},
		Status:    []domain.PageRoleRequestLogStatus{domain.PRSLPending},
	})
}

func (s Service) RejectPagePermissions(pagePkID int64, emails []string) *domain.Error {
	existingPage, err := s.repo.Page.GetByID(
		context.Background(),
		"",
		&pagePkID,
		domain.PageDetailOptions{},
		nil,
	)
	if err != nil {
		return err
	}

	err = s.repo.Page.UpdatePageAccessRequestStatus(context.Background(), domain.PageRoleRequestLogQuery{
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

	err = s.repo.Page.UpdatePageAccessRequestStatus(context.Background(), domain.PageRoleRequestLogQuery{
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
	page, pErr := s.repo.Page.GetByID(
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
	permission := s.repo.Page.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *page,
		User:     curUser,
		PageRole: curRole,
	})

	if !permission.CanView {
		return domain.ErrPermissionDenied
	}

	_, err := s.repo.Page.StarPage(context.Background(), input)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) RemovePageFromStarred(input domain.StarPageInput, curUser *domain.User) *domain.Error {
	// Handler Permissions
	page, pErr := s.repo.Page.GetByID(
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
	permission := s.repo.Page.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *page,
		User:     curUser,
		PageRole: curRole,
	})

	if !permission.CanView {
		return domain.ErrPermissionDenied
	}
	err := s.repo.Page.UnstarPage(context.Background(), input)
	if err != nil {
		return err
	}
	return nil
}
