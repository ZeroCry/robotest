package specs

import (
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

		var bandwagonLogin = framework.BandwagonConfig{
			Organization: defaults.BandwagonOrganization,
			Username:     defaults.BandwagonUsername,
			Password:     defaults.BandwagonPassword,
			Email:        defaults.BandwagonEmail,
		}

		if ctx.Bandwagon.Organization != "" {
			bandwagonLogin.Organization = ctx.Bandwagon.Organization
		}
		if ctx.Bandwagon.Username != "" {
			bandwagonLogin.Username = ctx.Bandwagon.Username
		}
		if ctx.Bandwagon.Password != "" {
			bandwagonLogin.Password = ctx.Bandwagon.Password
		}
		if ctx.Bandwagon.Email != "" {
			bandwagonLogin.Email = ctx.Bandwagon.Email
		}

		if ctx.Wizard {
			Skip("this test cannot run in wizard mode")
		}

		BeforeEach(func() {
			domainName = ctx.ClusterName
			siteURL = framework.SiteURL()
		})

		shouldHandleNewDeploymentScreen := func() {
			installer := installermodel.Open(f.Page, framework.InstallerURL())

			By("filling out license text field if required")
			installer.FillOutLicenseIfRequired(ctx.License)

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
			enableRemoteAccess := true
			bandwagon.Complete(f.Page,
				domainName,
				&bandwagonLogin,
				enableRemoteAccess)
		}

		shouldNavigateToSite := func() {
			By("opening a site page")
			ui.EnsureUser(f.Page, framework.SiteURL(), login)
			site.Open(f.Page, domainName)
		}

		It("should handle installation", func() {
			ui.EnsureUser(f.Page, framework.InstallerURL(), ctx.Login)
			shouldHandleNewDeploymentScreen()
			shouldHandleRequirementsScreen()
			shouldHandleInProgressScreen()
			shouldHandleBandwagonScreen()
			shouldNavigateToSite()
		})
	})
}
