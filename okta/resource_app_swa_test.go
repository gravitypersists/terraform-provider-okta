package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/okta/okta-sdk-golang/okta"
)

// Test creation of a simple AWS SWA app. The preconfigured apps are created by name.
func TestAccOktaAppSwaApplicationPreconfig(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestSwaConfigPreconfig(ri)
	updatedConfig := buildTestSwaConfigPreconfigUpdated(ri)
	resourceName := buildResourceFQN(appSwa, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(appSwa, createDoesAppExist(okta.NewSwaApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSwaApplication())),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSwaApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
				),
			},
		},
	})
}

// Test creation of a custom SAML app.
func TestAccOktaAppSwaApplication(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestSwaConfig(ri)
	updatedConfig := buildTestSwaConfigUpdated(ri)
	resourceName := buildResourceFQN(appSwa, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(appSwa, createDoesAppExist(okta.NewSwaApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSwaApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "button_field", "btn-login"),
					resource.TestCheckResourceAttr(resourceName, "password_field", "txtbox-password"),
					resource.TestCheckResourceAttr(resourceName, "username_field", "txtbox-username"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://example.com/login.html"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSwaApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://example.com/login.html"),
					resource.TestCheckResourceAttr(resourceName, "button_field", "btn-login"),
					resource.TestCheckResourceAttr(resourceName, "password_field", "txtbox-password"),
					resource.TestCheckResourceAttr(resourceName, "username_field", "txtbox-username"),
				),
			},
		},
	})
}

// Add and remove groups/users
func TestAccOktaAppSwaApplicationUserGroups(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestSwaGroupsUsers(ri)
	updatedConfig := buildTestSwaRemoveGroupsUsers(ri)
	resourceName := buildResourceFQN(appSwa, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(appSwa, createDoesAppExist(okta.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "1"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
		},
	})
}

func buildTestSwaConfigPreconfig(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  preconfigured_app		= "aws_console"
  label         		= "%s"
}
`, appSwa, name, name)
}

func buildTestSwaConfigPreconfigUpdated(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  preconfigured_app		= "aws_console"
  label         		= "%s"
  status 	   	 		= "INACTIVE"
}
`, appSwa, name, name)
}

func buildTestSwaGroupsUsers(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "okta_group" "group-%d" {
  name = "testAcc-%d"
}
resource "okta_user" "user-%d" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "blah"
  login       = "test-acc-%d@testing.com"
  email       = "test-acc-%d@testing.com"
  status      = "ACTIVE"
}

resource "%s" "%s" {
  preconfigured_app = "aws_console"
  label       	    = "%s"
  users {
    id = "${okta_user.user-%d.id}"
    username = "${okta_user.user-%d.email}"
  }
  groups = ["${okta_group.group-%d.id}"]
}
`, rInt, rInt, rInt, rInt, rInt, appSwa, name, name, rInt, rInt, rInt)
}

func buildTestSwaRemoveGroupsUsers(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "okta_group" "group-%d" {
  name = "testAcc-%d"
}

resource "okta_user" "user-%d" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "blah"
  login       = "test-acc-%d@testing.com"
  email       = "test-acc-%d@testing.com"
  status      = "ACTIVE"
}

resource "%s" "%s" {
  preconfigured_app  = "aws_console"
  label              = "%s"
}
`, rInt, rInt, rInt, rInt, rInt, appSwa, name, name)
}

func buildTestSwaConfig(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label           = "%s"
  button_field	  = "btn-login"
  password_field  = "txtbox-password"
  username_field  = "txtbox-username"
  url		  = "https://example.com/login.html"
}
`, appSwa, name, name)
}

func buildTestSwaConfigUpdated(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label           = "%s"
  status 	  = "INACTIVE"
  button_field	  = "btn-login"
  password_field  = "txtbox-password"
  username_field  = "txtbox-username"
  url		  = "https://example.com/login.html"
}
`, appSwa, name, name)
}
