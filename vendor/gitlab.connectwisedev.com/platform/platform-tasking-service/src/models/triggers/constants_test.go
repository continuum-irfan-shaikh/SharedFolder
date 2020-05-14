package triggers

import (
	"fmt"
	"testing"
)

func TestConstants(t *testing.T) {
	c := GenericTypePrefix
	c = AlertTypePrefix
	c = LogoutTrigger
	c = LoginTrigger
	c = StartupTrigger
	c = DynamicGroupExitTrigger
	c = DynamicGroupEnterTrigger
	c = ShutdownTrigger
	c = FirstCheckInTrigger
	c = MockGeneric
	c = MockAlerting
	fmt.Println(c)
}
