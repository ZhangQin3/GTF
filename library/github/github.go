package github

import "se"

type Github struct{ se.Page }

func OpenGithub(params ...func(caps map[string]interface{})) (Github, error) {
	url := "http://www.github.com/"
	page, err := se.OpenPage(url, nil, params...)

	return Github(page), err
}

func (p *Github) OpenURL(url string) error {
	page, err := se.OpenPage(url, nil)
	*p = Github(page)
	return err
}

func (p *Github) SignIn() *se.Element {
	return p.Link(se.ByLinkText, "Sign in")
}

func (p *Github) UserName() *se.Element {
	return p.TextBox(se.ById, "login_field")
}

func (p *Github) Password() *se.Element {
	return p.PasswordBox(se.ById, "password")
}

func (p *Github) Signin() *se.Element {
	return p.SubmitBtn(se.ByValue, "Sign in")
}

func (p *Github) NewProject() *se.Element {
	return p.Link(se.ByLinkText, "Start a project")
}

func (p *Github) Profile() *se.Element {
	// return p.Link(se.ByCssSelector, "a.header-nav-link.name.tooltipped.tooltipped-sw.js-menu-target")
	return p.Link(se.ByClassName, "header-nav-link name tooltipped tooltipped-sw js-menu-target")
}

func (p *Github) Logout() *se.Element {
	return p.Form(se.ByClassName, "logout-form")
}

// func OpenWanSetup(wd se.WebDriver) {
// 	page := struct{ Page }{Page{WD: wd}}
// 	k := page.link(LinkText, "WAN Setup").Click()
// 	println("====================", k)
// }

func (p *Github) HostName() *se.Element {
	return p.Table(se.ByClassName, "common_table").Tr(se.ByIndex, 1).Td(se.ByIndex, 1)
}

func (p *Github) WanSetup() *se.Element {
	return p.Link(se.ByLinkText, "Utilities")
}
