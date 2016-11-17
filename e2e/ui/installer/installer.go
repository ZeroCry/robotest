package installer

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gravitational/robotest/e2e/ui/common"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/sclevine/agouti"
	am "github.com/sclevine/agouti/matchers"
)

var defaultTimeout = 20 * time.Second

type Installer struct {
	page *agouti.Page
}

func OpenInstallerWithSite(page *agouti.Page, domainName string) *Installer {
	urlPrefix := fmt.Sprintf("/web/installer/site/%v", domainName)
	r, _ := regexp.Compile("/web/.*")
	url, _ := page.URL()
	url = r.ReplaceAllString(url, urlPrefix)

	return OpenInstaller(page, url)
}

func OpenInstaller(page *agouti.Page, URL string) *Installer {
	Expect(page.Navigate(URL)).To(Succeed())
	Eventually(page.FindByClass("grv-installer"), defaultTimeout).Should(am.BeFound())
	return &Installer{page: page}
}

func (i *Installer) CreateAwsSite(
	domainName string, accessKey string, secretKey string, region string, keyPair string, vpc string) string {
	By("Setting deployment name")
	page := i.page
	specifyDomainName(page, domainName)

	By("Setting provisioner")
	Expect(page.FindByClass("--aws").Click()).To(Succeed())
	Expect(page.FindByName("aws_access_key").Fill(accessKey)).To(Succeed())
	Expect(page.FindByName("aws_secret_key").Fill(secretKey)).To(Succeed())
	Expect(page.FindByClass("grv-installer-btn-new-site").Click()).To(Succeed())
	Eventually(page.FindByClass("grv-installer-aws-region"), defaultTimeout).Should(am.BeFound())

	By("Setting region")
	common.SetDropDownValue(page, "grv-installer-aws-region", region)
	By("Setting key pair")
	common.SetDropDownValue(page, "grv-installer-aws-key-pair", keyPair)
	By("Setting Vpc")
	common.SetDropDownValue(page, "grv-installer-aws-vpc", vpc)

	i.proceedToReqs()

	pageURL, _ := page.URL()
	return pageURL
}

func (i *Installer) CreateOnPremNewSite(domainName string) string {
	page := i.page
	By("Setting deployment name")
	specifyDomainName(page, domainName)

	By("Setting provisioner")
	Eventually(page.FindByClass("fa-check"), defaultTimeout).Should(am.BeFound())
	Expect(page.FindByClass("--metal").Click()).To(Succeed())

	i.proceedToReqs()

	pageURL, _ := page.URL()
	return pageURL
}

func (i *Installer) ProceedToSite() {
	Expect(i.page.Find(".grv-installer-progress-result .btn-primary").Click()).To(Succeed())
}

func (i *Installer) IsCreateSiteStep() bool {
	count, _ := i.page.FindByClass("grv-installer-fqdn").Count()
	return count != 0
}

func (i *Installer) IsInProgressStep() bool {
	count, _ := i.page.FindByClass("grv-installer-progres-indicator").Count()
	return count != 0
}

func (i *Installer) IsRequirementsReviewStep() bool {
	count, _ := i.page.FindByClass("grv-installer-provision-reqs").Count()
	return count != 0
}

func (i *Installer) StartInstallation() {
	btn := i.page.Find(".grv-installer-footer .btn-primary")
	Expect(btn).To(am.BeFound())
	Expect(btn.Click()).To(Succeed())
	Eventually(i.IsInProgressStep, defaultTimeout).Should(BeTrue())
}

func (i *Installer) IsInstallCompleted() bool {
	count, _ := i.page.Find(".grv-installer-progress-result .fa-check").Count()
	return count != 0
}

func (i *Installer) SelectFlavor(index int) {
	cssSelector := fmt.Sprintf(".grv-slider-value-desc:nth-child(%v) span", index)
	el := i.page.First(cssSelector)
	Expect(el).To(am.BeFound())
	Expect(el.Click()).To(Succeed())
}

func specifyDomainName(page *agouti.Page, domainName string) {
	Eventually(page.FindByName("domainName"), defaultTimeout).Should(am.BeFound())
	Expect(page.FindByName("domainName").Fill(domainName)).To(Succeed())
	Eventually(page.FindByClass("fa-check"), defaultTimeout).Should(am.BeFound())
}

func (i *Installer) proceedToReqs() {
	Expect(i.page.FindByClass("grv-installer-btn-new-site").Click()).To(Succeed())
	Eventually(i.page.FindByClass("grv-installer-provision-reqs"), defaultTimeout).Should(am.BeFound())
}
