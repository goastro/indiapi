// Package indiserver is used to start up and control an instance of indiserver. It can also be used to list all known INDI drivers on the current machine.
package indiserver

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path"
	"syscall"

	"github.com/rickbassham/goexec"
	"github.com/rickbassham/logging"
	"github.com/spf13/afero"
)

// Commander is an interface used to create a Command (which can be executed).
type Commander interface {
	Command(name string, args ...string) goexec.Command
}

// NewINDIServer creates a struct that can be used to get info about installed INDI drivers
// and start/stop a local indiserver.
func NewINDIServer(log logging.Logger, fs afero.Fs, port string, cmder Commander) *INDIServer {
	if len(port) == 0 {
		port = "7624"
	}

	s := &INDIServer{
		log:   log,
		fs:    fs,
		port:  port,
		cmder: cmder,
	}

	s.findDrivers()

	return s
}

// Driver represents an installed INDI driver.
type Driver struct {
	Label   string
	Driver  string
	Version string
}

type device struct {
	XMLName      xml.Name `xml:"device"`
	Label        string   `xml:"label,attr"`
	SkeletonFile string   `xml:"skel,attr"`
	Driver       string   `xml:"driver"`
	Version      string   `xml:"version"`
	Port         string   `xml:"port"`
}

type devGroup struct {
	XMLName xml.Name `xml:"devGroup"`
	Group   string   `xml:"group,attr"`
	Devices []device `xml:"device"`
}

type driversList struct {
	XMLName   xml.Name   `xml:"driversList"`
	DevGroups []devGroup `xml:"devGroup"`
}

// INDIServer is a struct to start/stop and control a local indiserver executable.
type INDIServer struct {
	log   logging.Logger
	fs    afero.Fs
	port  string
	cmder Commander

	fifoPath string
	fifo     io.WriteCloser
	cmd      goexec.Command

	drivers map[string][]Driver
}

func (s *INDIServer) findDrivers() {
	s.drivers = map[string][]Driver{}

	files, err := afero.Glob(s.fs, "/usr/share/indi/*.xml")
	if err != nil {
		s.log.WithError(err).Warn("error in afero.Glob")
		return
	}

	for _, fp := range files {
		f, err := s.fs.Open(fp)
		if err != nil {
			s.log.WithError(err).Warn("error in s.fs.Open")
			continue
		}

		var dl driversList

		err = xml.NewDecoder(f).Decode(&dl)
		if err != nil {
			s.log.WithError(err).Warn("error in xml.NewDecoder(f).Decode")
			continue
		}

		for _, dg := range dl.DevGroups {
			list, ok := s.drivers[dg.Group]
			if !ok {
				list = []Driver{}
			}

			for _, d := range dg.Devices {
				list = append(list, Driver{
					Driver:  d.Driver,
					Version: d.Version,
					Label:   d.Label,
				})
			}

			s.drivers[dg.Group] = list
		}
	}
}

// Drivers returns a list of drivers organzied by group.
func (s *INDIServer) Drivers() map[string][]Driver {
	return s.drivers
}

// StartServer starts up the indiserver. Be sure to call StopServer when you are done!
func (s *INDIServer) StartServer() error {
	if s.cmd != nil {
		return nil
	}

	dir, err := afero.TempDir(s.fs, "", "")
	if err != nil {
		s.log.WithError(err).Warn("error in afero.TempDir")
		return err
	}

	s.fifoPath = fmt.Sprintf("%s/fifo", dir)

	err = syscall.Mkfifo(s.fifoPath, 0666)
	if err != nil {
		s.log.WithError(err).Warn("error in syscall.Mkfifo")
		return err
	}

	s.cmd = s.cmder.Command("/usr/bin/indiserver", "-v", "-f", s.fifoPath, "-p", s.port)

	stdout, err := s.cmd.Stdout()
	if err != nil {
		s.log.WithError(err).Warn("error in s.cmd.Stdout")
		return err
	}

	stderr, err := s.cmd.Stderr()
	if err != nil {
		s.log.WithError(err).Warn("error in s.cmd.Stderr")
		return err
	}

	go func() {
		for line := range stdout {
			s.log.WithField("line", line).Info("from indiserver")
		}
	}()

	go func() {
		for line := range stderr {
			s.log.WithField("line", line).Info("from indiserver")
		}
	}()

	err = s.cmd.Start()
	if err != nil {
		s.log.WithError(err).Warn("error in s.cmd.Start")
		return err
	}

	s.fifo, err = s.fs.OpenFile(s.fifoPath, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		s.log.WithError(err).Warn("error in s.fs.OpenFile")
		return err
	}

	return nil
}

// StopServer stops the currently running indiserver and cleans up.
func (s *INDIServer) StopServer() error {
	if s.cmd == nil {
		return nil
	}

	defer func() {
		s.cmd = nil

		err := s.fs.RemoveAll(path.Dir(s.fifoPath))
		if err != nil {
			s.log.WithError(err).Warn("error in s.fs.RemoveAll")
		}
	}()

	err := s.cmd.Kill()
	if err != nil {
		s.log.WithError(err).Warn("error in s.cmd.Signal")
		return err
	}

	err = s.cmd.Wait()
	if err != nil {
		if err.Error() == "signal: killed" {
			// We just killed it. It's not an error.
			return nil
		}

		s.log.WithError(err).Warn("error in s.cmd.Wait")
		return err
	}

	return nil
}

// StartDriver starts up a driver on the indiserver. Note that this will NOT return an
// error if the indiserver doesn't recognize the driver or if it has any other issues.
// Watch the log for info on failures inside indiserver.
func (s *INDIServer) StartDriver(driver, name string) error {
	cmd := fmt.Sprintf("start %s -n \"%s\"\n", driver, name)

	_, err := s.fifo.Write([]byte(cmd))
	if err != nil {
		s.log.WithError(err).Warn("error in s.fifo.Write")
		return err
	}

	return nil
}

// StopDriver stops a driver on the indiserver.
func (s *INDIServer) StopDriver(driver, name string) error {
	cmd := fmt.Sprintf("stop %s \"%s\"\n", driver, name)

	_, err := s.fifo.Write([]byte(cmd))
	if err != nil {
		s.log.WithError(err).Warn("error in s.fifo.Write")
		return err
	}

	return nil
}
