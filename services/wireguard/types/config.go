package types

import (
	"bytes"
	"crypto/rand"
	"io/ioutil"
	"math/big"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	ct = strings.TrimSpace(`
# Name of the network interface
interface = "{{ .Interface }}"

# Port number to accept the incoming connections
listen_port = {{ .ListenPort }}

# Server private key
private_key = "{{ .PrivateKey }}"
	`)

	t = func() *template.Template {
		t, err := template.New("").Parse(ct)
		if err != nil {
			panic(err)
		}

		return t
	}()
)

type Config struct {
	Interface  string `json:"interface" mapstructure:"interface"`
	ListenPort uint16 `json:"listen_port" mapstructure:"listen_port"`
	PrivateKey string `json:"private_key" mapstructure:"private_key"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Validate() error {
	if c.Interface == "" {
		return errors.New("interface cannot be empty")
	}
	if c.ListenPort == 0 {
		return errors.New("listen_port cannot be zero")
	}
	if c.PrivateKey == "" {
		return errors.New("private_key cannot be empty")
	}
	if _, err := KeyFromString(c.PrivateKey); err != nil {
		return errors.Wrap(err, "invalid private_key")
	}

	return nil
}

func (c *Config) WithDefaultValues() *Config {
	n, err := rand.Int(rand.Reader, big.NewInt(1<<16-1<<10))
	if err != nil {
		panic(err)
	}

	key, err := NewPrivateKey()
	if err != nil {
		panic(err)
	}

	c.Interface = "wg0"
	c.ListenPort = uint16(n.Int64() + 1<<10)
	c.PrivateKey = key.String()

	return c
}

func (c *Config) SaveToPath(path string) error {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		return err
	}

	return ioutil.WriteFile(path, buffer.Bytes(), 0600)
}

func (c *Config) String() string {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		panic(err)
	}

	return buffer.String()
}

func ReadInConfig(v *viper.Viper) (*Config, error) {
	cfg := NewConfig().WithDefaultValues()
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}