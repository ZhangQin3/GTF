package arrs

import "se"

type WanSetup struct{ se.Page }

func OpenWanSetup(parent LoginPage) WanSetup {
	parent.Link(se.ByLinkText, "WAN Setup").Click()

	return WanSetup(parent)
}

// func (lp *GWTWLoginPage) Apply() *PageElem {
// 	return lp.button(Value, "Apply")
// }

// func (lp *GWTWLoginPage) UserName() *PageElem {
// 	return lp.textBox(Id, "UserName")
// }

// func (lp *GWTWLoginPage) Passord() *PageElem {
// 	return lp.passwordBox(Id, "Password")
// }

// func (lp *GWTWLoginPage) Logout() *PageElem {
// 	return lp.link(LinkText, "ogout")
// }

// func (lp *GWTWLoginPage) HostName() *PageElem {
// 	return lp.table(ClassName, "common_table").tr(Index, 1).td(Index, 1)
// }

// func (lp *GWTWLoginPage) WanSetup() *PageElem {
// 	return lp.link(LinkText, "Utilities")
// }
