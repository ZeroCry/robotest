package opscenter

import (
	"fmt"

	"github.com/gravitational/robotest/e2e/uimodel/defaults"
	"github.com/gravitational/robotest/e2e/uimodel/utils"

	log "github.com/Sirupsen/logrus"
	. "github.com/onsi/gomega"
	web "github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

// OpsCenter is opscenter ui model
type OpsCenter struct {
	page *web.Page
	url  string
}

// Open navigates to opscenter URL and returns ui model
func Open(page *web.Page) OpsCenter {
	log.Infof("trying to open opscenter")
	url := utils.GetOpsCenterURL(page)
	Expect(page.Navigate(url)).To(Succeed())
	Eventually(page.FindByClass("grv-portal"), defaults.AppLoadTimeout).
		Should(BeFound(), "waiting for opscenter to load")

	utils.PauseForComponentJs()
	return OpsCenter{page: page, url: url}
}

// DeleteSite deletes cluster by its name
func (o *OpsCenter) DeleteSite(domainName string, awsAccessKey string, awsSecrenKey string) {
	log.Infof("selecting a site to delete")
	deploymentIndex := getDeploymentIndex(o.page, domainName)
	Expect(deploymentIndex).To(BeNumerically(">=", 0), "expected to find a valid deployment index")
	utils.SetDropdownValue2(o.page, fmt.Sprintf(".grv-portal-sites tr:nth-child(%v)", deploymentIndex+1), "button", "Delete...")

	log.Infof("entering AWS credentials")
	elems := o.page.FindByName("aws_access_key")
	count, _ := elems.Count()
	if count > 0 {
		Expect(elems).To(BeFound(), "expected to find an input field for AWS access key")
		Expect(elems.SendKeys(awsAccessKey)).To(Succeed(), "expected to input AWS access key")

		elems = o.page.FindByName("aws_secret_key")
		Expect(elems).To(BeFound(), "expected to find an input field for AWS secret key")
		Expect(elems.SendKeys(awsSecrenKey)).To(Succeed(), "expected to input AWS secret key")
	}

	log.Infof("confirming cluster name")
	elems = o.page.FindByName("deploymentName")
	Expect(elems).To(BeFound())
	Expect(elems.SendKeys(domainName)).To(Succeed())

	log.Infof("confirming the action")
	Expect(o.page.Find(".grv-dialog .btn-danger").Click()).To(Succeed())
	Eventually(
		func() bool {
			log.Infof("checking if cluster %v disappered from the list", domainName)
			return getDeploymentIndex(o.page, domainName) >= 0
		},
		defaults.OpsCenterDeleteSiteTimeout,
		defaults.OpsCenterDeleteSitePollInterval).ShouldNot(BeTrue(), "cluster should disappear from the cluster list")
}

func getDeploymentIndex(page *web.Page, domainName string) int {
	var deploymentIndex int
	const scriptTemplate = `
            var rows = Array.prototype.slice.call(document.querySelectorAll(".grv-portal-sites .grv-table .grv-portal-sites-tag"));
            return rows.findIndex( (tag) => {
		    return (tag.innerText == "Name:%v");
            });
        `

	script := fmt.Sprintf(scriptTemplate, domainName)
	Expect(page.RunScript(script, nil, &deploymentIndex)).To(Succeed())
	return deploymentIndex
}
