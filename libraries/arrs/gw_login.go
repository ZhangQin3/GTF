package arrs

import "gse"

type LoginPage struct{ gse.Page }

func OpenLoginPage() LoginPage {
	url := "http://10.89.80.20:8080/"
	page := gse.OpenPage(url, nil)

	return LoginPage(page)
}

// func OpenWanSetup(wd gse.WebDriver) {
// 	page := struct{ Page }{Page{WD: wd}}
// 	k := page.link(LinkText, "WAN Setup").Click()
// 	println("====================", k)
// }

func (p *LoginPage) Apply() *gse.Element {
	return p.Button(gse.ByValue, "Apply")
}

func (p *LoginPage) UserName() *gse.Element {
	return p.TextBox(gse.ById, "UserName")
}

func (p *LoginPage) Passord() *gse.Element {
	return p.PasswordBox(gse.ById, "Password1")
}

func (p *LoginPage) Logout() *gse.Element {
	return p.Link(gse.ByLinkText, "ogout")
}

func (p *LoginPage) HostName() *gse.Element {
	return p.Table(gse.ByClassName, "common_table").Tr(gse.ByIndex, 1).Td(gse.ByIndex, 1)
}

func (p *LoginPage) WanSetup() *gse.Element {
	return p.Link(gse.ByLinkText, "Utilities")
}
