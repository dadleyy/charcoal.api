package routes

import "fmt"
import "strconv"
import "github.com/albrow/forms"
import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/errors"

const emailDomainRestriction = "restricted_email_domains"

func FindSystemEmailDomains(runtime *net.RequestRuntime) error {
	blueprint := runtime.Blueprint()
	var domains []models.SystemEmailDomain

	total, err := blueprint.Apply(&domains)

	if err != nil {
		runtime.Debugf("ERR_BAD_ROLE_LOOKUP: %s", err.Error())
		return runtime.AddError(fmt.Errorf(errors.ErrFailedQuery))
	}

	for _, domains := range domains {
		runtime.AddResult(domains)
	}

	runtime.SetMeta("total", total)

	return nil
}

func CreateSystemEmailDomain(runtime *net.RequestRuntime) error {
	body, err := forms.Parse(runtime.Request)

	if err != nil {
		return runtime.AddError(err)
	}

	validator := body.Validator()
	validator.Require("domain")
	validator.MinLength("domain", 2)

	// if the validator picked up errors, add them to the request
	// runtime and then return
	if validator.HasErrors() == true {
		for _, m := range validator.Messages() {
			runtime.AddError(fmt.Errorf(m))
		}

		return nil
	}

	domain := models.SystemEmailDomain{Domain: body.Get("domain")}
	cursor := runtime.Model(&domain)
	existing := 0

	if err := cursor.Where("domain = ?", domain.Domain).Count(&existing).Error; err != nil || existing >= 1 {
		return runtime.AddError(fmt.Errorf(errors.ErrDuplicateEntry))
	}

	if err := cursor.Create(&domain).Error; err != nil {
		return runtime.AddError(err)
	}

	runtime.AddResult(domain)
	return nil
}

func DestroySystemEmailDomain(runtime *net.RequestRuntime) error {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.AddError(fmt.Errorf("BAD_DOMAIN_ID"))
	}

	var domain models.SystemEmailDomain

	if err := runtime.First(&domain, id).Error; err != nil {
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	if err := runtime.Delete(&domain).Error; err != nil {
		return runtime.AddError(err)
	}

	return nil
}

func UpdateSystem(runtime *net.RequestRuntime) error {
	var settings models.SystemSettings

	if err := runtime.FirstOrCreate(&settings, models.SystemSettings{}).Error; err != nil {
		return runtime.AddError(err)
	}

	body, err := forms.Parse(runtime.Request)

	if err != nil {
		return runtime.AddError(err)
	}

	updates := make(map[string]interface{})
	restriction := body.KeyExists(emailDomainRestriction)

	if value := body.Get(emailDomainRestriction); restriction && (value != "true" && value != "false") {
		return runtime.AddError(fmt.Errorf("INVALID_EMAIL_RESTRICTION"))
	}

	if restriction {
		value, _ := strconv.ParseBool(body.Get(emailDomainRestriction))
		updates["restricted_email_domains"] = value
	}

	if len(updates) == 0 {
		runtime.AddResult(settings)
		return nil
	}

	cursor := runtime.Model(&settings).Where("id = ?", settings.ID)

	if err := cursor.Updates(updates).Error; err != nil {
		return runtime.AddError(err)
	}

	runtime.AddResult(settings)
	return nil
}

func PrintSystem(runtime *net.RequestRuntime) error {
	var settings models.SystemSettings

	if admin := runtime.IsAdmin(); admin != true {
		runtime.Debugf("non-admin access of system route: %d", runtime.User.ID)
		runtime.AddResult("OK")
		return nil
	}

	if err := runtime.First(&settings).Error; err != nil {
		return runtime.AddError(err)
	}

	runtime.AddResult(settings)
	return nil
}
