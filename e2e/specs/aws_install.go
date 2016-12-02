package specs

import (
	"fmt"

	"github.com/gravitational/robotest/e2e/framework"
	"github.com/gravitational/robotest/e2e/model/ui"
	"github.com/gravitational/robotest/e2e/model/ui/defaults"
	installermodel "github.com/gravitational/robotest/e2e/model/ui/installer"
	"github.com/gravitational/robotest/e2e/model/ui/site"
	bandwagon "github.com/gravitational/robotest/e2e/specs/asserts/bandwagon"
	validation "github.com/gravitational/robotest/e2e/specs/asserts/installer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func VerifyAWSInstall(f *framework.T) {

	framework.RoboDescribe("AWS Installation", func() {
		ctx := framework.TestContext
		var domainName string
		var siteURL string
		var login = framework.Login{
			Username: defaults.BandwagonUsername,
			Password: defaults.BandwagonPassword,
		}

		BeforeEach(func() {
			domainName = ctx.ClusterName
			siteURL = framework.SiteURL()
		})

		shouldProvideLicense := func() {
			installer := installermodel.Open(f.Page, framework.InstallerURL())
			By("filling out license text field if required")
			installer.FillOutLicenseIfRequired(ctx.License)
		}

		shouldHandleNewDeploymentScreen := func() {
			installer := installermodel.Open(f.Page, framework.InstallerURL())

			Eventually(installer.IsCreateSiteStep, defaults.FindTimeout).Should(
				BeTrue(),
				"should navigate to installer screen")

			installer.CreateAWSSite(domainName, ctx.AWS)
		}

		shouldHandleRequirementsScreen := func() {
			By("entering domain name")
			installer := installermodel.OpenWithSite(f.Page, domainName)
			Expect(installer.IsRequirementsReviewStep()).To(
				BeTrue(),
				"should be on requirement step")

			By("selecting a flavor")
			installer.SelectFlavorByLabel(ctx.FlavorLabel)

			profiles := installermodel.FindAWSProfiles(f.Page)

			Expect(len(profiles)).To(
				Equal(1),
				"should verify required node number")

			profiles[0].SetInstanceType(ctx.AWS.InstanceType)

			By("starting an installation")
			installer.StartInstallation()
		}

		shouldHandleInProgressScreen := func() {
			validation.WaitForComplete(f.Page, domainName)
		}

		shouldHandleBandwagonScreen := func() {
			enableRemoteAccess := ctx.ForceRemoteAccess || !ctx.Wizard
			// useLocalEndpoint := ctx.ForceLocalEndpoint || ctx.Wizard
			endpoints := bandwagon.Complete(f.Page,
				domainName,
				login,
				enableRemoteAccess)

			By("using local application endpoint")
			serviceLogin := &framework.ServiceLogin{Username: login.Username, Password: login.Password}
			siteEntryURL := endpoints[0]
			// TODO: for terraform, use public installer address
			// terraform nodes are provisioned only with a single private network interface
			if ctx.Provisioner == "terraform" {
				siteEntryURL = fmt.Sprintf("https://%v:%v", framework.InstallerNode().Addr(), defaults.GravityHTTPPort)
			}
			framework.UpdateSiteEntry(siteEntryURL, login, serviceLogin)
		}

		shouldNavigateToSite := func() {
			By("opening a site page")
			ui.EnsureUser(f.Page, framework.SiteURL(), login)
			site.Open(f.Page, domainName)
		}

		It("should handle installation", func() {
			ui.EnsureUser(f.Page, framework.InstallerURL(), ctx.Login)
			shouldProvideLicense()
			shouldHandleNewDeploymentScreen()
			shouldHandleRequirementsScreen()
			shouldHandleInProgressScreen()
			shouldHandleBandwagonScreen()
			shouldNavigateToSite()
		})
	})
}
