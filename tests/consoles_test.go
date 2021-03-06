package integration

import . "gopkg.in/check.v1"

func (s *QemuSuite) TestCloudConfigConsole(c *C) {
	s.RunQemu(c, "--cloud-config", "./tests/assets/test_03/cloud-config.yml")

	s.CheckCall(c, "apt-get --version")
	s.CheckCall(c, `
sudo ros console list | grep default | grep disabled
sudo ros console list | grep debian | grep current`)
}

func (s *QemuSuite) TestConsoleCommand(c *C) {
	s.RunQemu(c)

	s.CheckCall(c, `
sudo ros console list | grep default | grep current
sudo ros console list | grep debian | grep disabled
(sudo ros console switch invalid 2>&1 || true) | grep "invalid is not a valid console"
(sudo ros console enable invalid 2>&1 || true) | grep "invalid is not a valid console"`)

	s.MakeCall("sudo ros console switch -f debian")
	c.Assert(s.WaitForSSH(), IsNil)

	s.CheckCall(c, "apt-get --version")
	s.CheckCall(c, `
sudo ros console list | grep default | grep disabled
sudo ros console list | grep debian | grep current`)

	s.Reboot(c)

	s.CheckCall(c, "apt-get --version")
	s.CheckCall(c, `
sudo ros console list | grep default | grep disabled
sudo ros console list | grep debian | grep current`)

	s.MakeCall("sudo ros console switch -f default")
	c.Assert(s.WaitForSSH(), IsNil)

	s.CheckCall(c, `
sudo ros console list | grep default | grep current
sudo ros console list | grep debian | grep disabled`)

	s.CheckCall(c, "sudo ros console enable debian")

	s.CheckCall(c, "sudo ros console list | grep default | grep current")
	s.CheckCall(c, "sudo ros console list | grep debian | grep enabled")

	s.Reboot(c)

	s.CheckCall(c, `
sudo ros console list | grep default | grep disabled
sudo ros console list | grep debian | grep current`)
}

func (s *QemuSuite) TestContainersRestartAfterConsoleSwitch(c *C) {
	s.RunQemu(c)

	s.CheckCall(c, `
docker run -d --restart=always nginx
docker ps | grep nginx`)

	s.MakeCall("sudo ros console switch -f debian")
	c.Assert(s.WaitForSSH(), IsNil)

	s.CheckCall(c, "docker ps | grep nginx")
	s.CheckCall(c, "sudo ros console enable default")

	s.Reboot(c)

	s.CheckCall(c, "docker ps | grep nginx")
}
