package openconnect

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	octypes "github.com/sentinel-official/dvpn-node/services/openconnect/types"
	"github.com/sentinel-official/dvpn-node/types"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type OpenConnect struct {
	info   []byte
	cmd    *exec.Cmd
	config *octypes.Config
	peers  *octypes.Peers
}

func NewOpenConnect() types.Service {
	return &OpenConnect{}
}

func (s *OpenConnect) configFilePath() string {
	return filepath.Join(os.TempDir(), "ocserv.conf")
}

func (s *OpenConnect) Type() uint64 {
	return octypes.Type
}

func (s *OpenConnect) Init(home string) (err error) {
	v := viper.New()
	v.SetConfigFile(filepath.Join(home, octypes.ConfigFileName))

	s.config, err = octypes.ReadInConfig(v)
	if err != nil {
		return err
	}
	if err = s.config.Validate(); err != nil {
		return err
	}

	t, err := template.New("ocserv").Parse(configTemplate)
	if err != nil {
		return err
	}

	panic("implement me") //TODO
}

func (s *OpenConnect) Info() []byte {
	return s.info
}

func (s *OpenConnect) Start() error {
	s.cmd = exec.Command("ocserv", strings.Split(
		fmt.Sprintf("--config %s", s.configFilePath()), " ")...)

	s.cmd.Env = os.Environ()

	s.cmd.Stdout = os.Stdout
	s.cmd.Stderr = os.Stderr

	return s.cmd.Start()
}

func (s *OpenConnect) Stop() error {
	if s.cmd == nil {
		return errors.New("command is nil")
	}

	return s.cmd.Process.Kill()
}

func (s *OpenConnect) AddPeer(data []byte) (result []byte, err error) {
	// TODO: Add to Radius server
	panic("implement me")
}

func (s *OpenConnect) HasPeer(data []byte) bool {
	var (
		username = base64.StdEncoding.EncodeToString(data)
		peer     = s.peers.Get(username)
	)

	return !peer.Empty()
}

func (s *OpenConnect) RemovePeer(data []byte) error {
	// TODO: Remove from Radius server
	username := base64.StdEncoding.EncodeToString(data)

	panic("implement me")
}

func (s *OpenConnect) Peers() (items []types.Peer, err error) {
	// TODO: Get them from Radius server
	panic("implement me")
}

func (s *OpenConnect) PeerCount() int {
	return s.peers.Len()
}
