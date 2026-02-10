package admin

import (
	"bambamload/constant"
	"bambamload/enum"
	"bambamload/logger"
	"bambamload/models"
	"bambamload/utils"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

func (sa *ServiceAdmin) ApproveOrRejectSupplierKyc(req models.ApproveOrRejectSupplierRequest, user *models.User) error {

	updateMap := make(map[string]interface{})
	switch req.Action {
	case constant.Approve:
		updateMap[fmt.Sprintf("%s_approved", req.DocumentKey)] = true
		updateMap[fmt.Sprintf("%s_comment", req.DocumentKey)] = req.Comment
	case constant.Reject:
		updateMap[fmt.Sprintf("%s_approved", req.DocumentKey)] = false
		updateMap[fmt.Sprintf("%s_comment", req.DocumentKey)] = req.Comment
	}

	err := sa.PostgresRepository.UpdateUser(req.SupplierID, constant.ID, updateMap)
	if err != nil {
		logger.Logger.Errorf("[ApproveOrRejectSupplierKyc]failed to update user supplier: %v", err)
		return errors.New("unable to approve or reject supplier kyc,please try again later")
	}
	//make audit trail

	return nil
}

func (sa *ServiceAdmin) ApproveOrRejectSupplier(req models.ApproveOrRejectSupplierRequest, user *models.User) error {
	updateMap := make(map[string]interface{})
	switch req.Action {
	case constant.Approve:
		updateMap["kyc_status"] = constant.Approved
		updateMap["status"] = constant.Approved
		updateMap["is_active"] = true
	case constant.Reject:
		updateMap["kyc_status"] = constant.Rejected
		updateMap["status"] = constant.Rejected
		updateMap["supplier_reject_reason"] = req.Comment
	}

	err := sa.PostgresRepository.UpdateUser(req.SupplierID, constant.ID, updateMap)
	if err != nil {
		logger.Logger.Errorf("failed to update user supplier: %v", err)
		return errors.New("unable to approve or reject supplier,please try again later")
	}
	return nil
}

func (sa *ServiceAdmin) ChangeSupplierCommissionRate(req models.EditSupplierCommissionRate, user *models.User) error {
	updateMap := make(map[string]interface{})
	updateMap["commission_rate"] = req.Rate
	err := sa.PostgresRepository.UpdateUser(req.SupplierID, constant.ID, updateMap)
	if err != nil {
		logger.Logger.Errorf("failed to update user supplier: %v", err)
		return errors.New("unable to change supplier commission rate")
	}

	return nil
}

func (sa *ServiceAdmin) InviteSupplier(req models.InviteSupplier, user *models.User) (string, error) {

	sEmail := strings.ToLower(strings.TrimSpace(req.Email))
	phoneNumber := utils.StandardiseMSISDN(strings.TrimSpace(req.PhoneNumber))

	if sa.PostgresRepository.UserExists(sEmail, constant.Email) {
		return "Supplier with this email already exists", errors.New("email already exists")
	}

	if sa.PostgresRepository.UserExists(phoneNumber, constant.PhoneNumber) {
		return "Supplier with this phone number already exists", errors.New("phone number already exists")
	}

	ref := utils.GenerateReference("")
	supplier := &models.User{
		Name:         req.ContactPerson,
		Email:        sEmail,
		PhoneNumber:  phoneNumber,
		BusinessName: req.BusinessName,
		Status:       constant.Invited,
		Role:         enum.Supplier,
		Reference:    ref,
	}
	err := sa.PostgresRepository.CreateUser(supplier)
	if err != nil {
		logger.Logger.Errorf("Failed to create supplier: %v", err)
		return "unable to invite supplier, please try again later", err
	}

	url := fmt.Sprintf("%s/auth/onboarding/setup?reference=%s", os.Getenv("FRONTEND_URL"), ref)
	body := utils.BuildSupplierInviteEmail(req.BusinessName, url, req.InvitationMessage)

	err = sa.EmailService.Send(sEmail, "Invitation to Join BamBamLoad", body)
	if err != nil {
		logger.Logger.Errorf("[InviteSupplier]Failed to send email: %v", err)
	}

	return "success", nil
}

func (sa *ServiceAdmin) ResendSupplierInviteEmail(reference string) (string, error) {
	user, err := sa.PostgresRepository.GetUser(reference, constant.Reference)
	if err != nil {
		logger.Logger.Errorf("[ResendSupplierInviteEmail]Failed to get user: %v", err)
		return "unable to send invite mail, please try again later", err
	}
	if user.Role != enum.Supplier {
		return "cannot send mail to non suppliers", errors.New("supplier not allowed")
	}
	if user.Status != constant.Invited {
		return "supplier is not invited", errors.New("supplier is not invited")
	}

	url := fmt.Sprintf("%s/auth/onboarding/setup?reference=%s", os.Getenv("FRONTEND_URL"), user.Reference)
	body := utils.BuildSupplierInviteEmail(user.BusinessName, url, "")

	err = sa.EmailService.Send(user.Email, "Invitation to Join BamBamLoad", body)
	if err != nil {
		logger.Logger.Errorf("[InviteSupplier]Failed to send email: %v", err)
	}

	return "success", nil
}

func (sa *ServiceAdmin) SupplierDashboardCards() (any, error) {
	cards, err := sa.PostgresRepository.GetSupplierCards()
	if err != nil {
		logger.Logger.Errorf("[SupplierDashboardCards]Failed to get supplier cards: %v", err)
		return nil, errors.New("unable to get supplier cards")
	}
	return cards, nil
}

func (sa *ServiceAdmin) GetSuppliers(pm *models.PaginationMetadata, status, searchText string) ([]models.User, *models.PaginationMetadata, error) {

	suppliers, paginationMetaData, err := sa.PostgresRepository.GetSuppliers(pm, status, searchText)
	if err != nil {
		logger.Logger.Errorf("[GetSuppliers]Failed to get suppliers: %v", err)
		return nil, pm, err
	}

	return suppliers, paginationMetaData, nil
}

func (sa *ServiceAdmin) GetSupplier(id string) (any, error) {

	user, err := sa.PostgresRepository.GetUser(id, constant.ID)
	if err != nil {
		logger.Logger.Errorf("[GetSupplier]Failed to get user: %v", err)
		return nil, errors.New("unable to get supplier")
	}
	if user.Role != enum.Supplier {
		return nil, errors.New("user not supplier")
	}
	return user, nil
}

func (sa *ServiceAdmin) ApproveOrRejectSupplierProduct(req models.ApproveOrRejectSupplierProduct, user *models.User) error {

	updateMap := make(map[string]interface{})
	if req.Action == constant.Approve {
		updateMap["approval_status"] = constant.Approved
		updateMap["status"] = constant.Active
		updateMap["approved_by"] = user.Name
		updateMap["date_approved"] = time.Now().UTC()
	} else {
		updateMap["approval_status"] = constant.Rejected
		updateMap["rejected_by"] = user.Name
		updateMap["date_rejected"] = time.Now().UTC()
		updateMap["reject_reason"] = req.Comment
	}

	err := sa.PostgresRepository.UpdateProduct(req.ProductID, constant.ID, updateMap)
	if err != nil {
		logger.Logger.Errorf("[ApproveOrRejectSupplierProduct]Failed to update product: %v", err)
		return errors.New("unable to perform action, please try again later")
	}

	return nil
}

func (sa *ServiceAdmin) GetProductStats() (any, error) {

	stats, err := sa.PostgresRepository.GetProductStats()
	if err != nil {
		logger.Logger.Errorf("[GetProductStats]Failed to get product stats: %v", err)
		return nil, err
	}
	return stats, nil
}
