package arrs

import "se"

type LoginPage struct{ se.Page }

func OpenLoginPage() LoginPage {
	url := "http://10.89.80.20:8080/"
	page := se.OpenPage(url, nil)

	return LoginPage(page)
}

// func OpenWanSetup(wd se.WebDriver) {
// 	page := struct{ Page }{Page{WD: wd}}
// 	k := page.link(LinkText, "WAN Setup").Click()
// 	println("====================", k)
// }

func (p *LoginPage) Apply() *se.Element {
	return p.Button(se.ByValue, "Apply")
}

func (p *LoginPage) UserName() *se.Element {
	return p.TextBox(se.ById, "UserName")
}

func (p *LoginPage) Passord() *se.Element {
	return p.PasswordBox(se.ById, "Password1")
}

func (p *LoginPage) Logout() *se.Element {
	return p.Link(se.ByLinkText, "ogout")
}

func (p *LoginPage) HostName() *se.Element {
	return p.Table(se.ByClassName, "common_table").Tr(se.ByIndex, 1).Td(se.ByIndex, 1)
}

func (p *LoginPage) WanSetup() *se.Element {
	return p.Link(se.ByLinkText, "Utilities")
}
