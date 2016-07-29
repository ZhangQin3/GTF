package github

import "se"

type Github struct{ se.Page }

func OpenGithub() Github {
	url := "http://www.github.com/"
	page := se.OpenPage(url, nil)

	return Github(page)
}

func (p *Github) SignIn() *se.Element {
	return p.Link(se.ByCssSelector, "a.btn.site-header-actions-btn.mr-2")
}

func (p *Github) UserName() *se.Element {
	return p.TextBox(se.ById, "user\\[login\\]")
}

// func OpenWanSetup(wd se.WebDriver) {
// 	page := struct{ Page }{Page{WD: wd}}
// 	k := page.link(LinkText, "WAN Setup").Click()
// 	println("====================", k)
// }

func (p *Github) Apply() *se.Element {
	return p.Button(se.ByValue, "Apply")
}

// func (p *Github) UserName() *se.Element {
// 	return p.TextBox(se.ById, "UserName")
// }

func (p *Github) Passord() *se.Element {
	return p.PasswordBox(se.ById, "Password1")
}

func (p *Github) Logout() *se.Element {
	return p.Link(se.ByLinkText, "ogout")
}

func (p *Github) HostName() *se.Element {
	return p.Table(se.ByClassName, "common_table").Tr(se.ByIndex, 1).Td(se.ByIndex, 1)
}

func (p *Github) WanSetup() *se.Element {
	return p.Link(se.ByLinkText, "Utilities")
}
