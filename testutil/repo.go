package testutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	conf "github.com/ipfs/ipfs-ds-convert/config"

	madns "gx/ipfs/QmQMRYmPn77CKRFf4YFjX3M5e6uw6DFAgsQffCX6mwZ4mA/go-multiaddr-dns"
	maddr "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	config "gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/repo/config"
	fsrepo "gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/repo/fsrepo"
)

//Hack
func init() {
	maddr.Protocols = append(maddr.Protocols, madns.DnsaddrProtocol)
}

func NewTestRepo(t *testing.T, spec map[string]interface{}) (string, func(t *testing.T)) {
	conf, err := config.Init(os.Stdout, 1024)
	if err != nil {
		t.Fatal(err)
	}

	err = config.Profiles["test"].Transform(conf)
	if err != nil {
		t.Fatal(err)
	}

	if spec != nil {
		conf.Datastore.Spec = spec
	}

	repoRoot, err := ioutil.TempDir(os.TempDir(), "ds-convert-test-")
	if err != nil {
		t.Fatal(err)
	}

	if err := fsrepo.Init(repoRoot, conf); err != nil {
		t.Fatal(err)
	}

	return repoRoot, func(t *testing.T) {
		err := os.RemoveAll(repoRoot)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func PatchConfig(t *testing.T, configPath string, newSpecPath string) {
	newSpec := make(map[string]interface{})
	err := conf.Load(newSpecPath, &newSpec)
	if err != nil {
		t.Fatal(err)
	}

	repoConfig := make(map[string]interface{})
	err = conf.Load(configPath, &repoConfig)
	if err != nil {
		t.Fatal(err)
	}

	dsConfig, ok := repoConfig["Datastore"].(map[string]interface{})
	if !ok {
		t.Fatal(fmt.Errorf("no 'Datastore' or invalid type in %s", configPath))
	}

	_, ok = dsConfig["Spec"].(map[string]interface{})
	if !ok {
		t.Fatal(fmt.Errorf("no 'Datastore.Spec' or invalid type in %s", configPath))
	}

	dsConfig["Spec"] = newSpec

	b, err := json.MarshalIndent(repoConfig, "", "  ")
	ioutil.WriteFile(configPath, b, 0660)
}
