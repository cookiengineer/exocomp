//go:build windows

package cli

import "golang.org/x/sys/windows"
import "fmt"

func CheckRoot() error {

	var token windows.Token

	err0  := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token)

	if err0 == nil {

		defer token.Close()

		admin_sid, err1 := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)

		if err1 == nil {

			is_member, err2 := token.IsMember(admin_sid)

			if err2 == nil {

				if is_member == true {
					return nil
				} else {
					return fmt.Errorf("Error: Run this program as admin")
				}

			} else {
				return fmt.Errorf("Error: %s", err2.Error())
			}

		} else {
			return fmt.Errorf("Error: %s", err1.Error())
		}

	} else {
		return fmt.Errorf("Error: %s", err0.Error())
	}

}
