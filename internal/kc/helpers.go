package kc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func getConfigFile(f string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error when accessing $HOME dir: %v", err)
	}

	configfile := filepath.Join(home, f)
	if _, err := os.Stat(configfile); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Configuration file missing. Make sure %v exists\n", configfile)
	}
	return configfile
}

func readConf(f string) (*kubeconf, error) {
	buf, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	conf := &kubeconf{}
	err = yaml.Unmarshal(buf, conf)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", f, err)
	}
	return conf, nil
}

func (k *kubeconf) writeConf(f string) {
	d, err := yaml.Marshal(&k)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	file, err := os.Create(f)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(f, d, 0644)
	if err != nil {
		log.Fatal(err)
	}

	file.Close()
}

func printConf(k *kubeconf) {
	d, err := yaml.Marshal(&k)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println(string(d))
}

func (k *kubeconf) uniqueCluster(c *Clusters) bool {
	unique := true
	for _, cluster := range k.Clusters {
		if cluster.Name == c.Name && cluster.Cluster.CertificateAuthorityData == c.Cluster.CertificateAuthorityData && cluster.Cluster.Server == c.Cluster.Server {
			unique = false
		}
	}

	return unique
}

func (k *kubeconf) uniqueContext(c *Contexts) bool {
	unique := true
	for _, context := range k.Contexts {
		if context.Name == c.Name && context.Context.User == c.Context.User && context.Context.Cluster == c.Context.Cluster {
			unique = false
		}
	}

	return unique
}

func (k *kubeconf) uniqueUser(u *Users) bool {
	unique := true
	for _, user := range k.Users {
		if user.Name == u.Name && user.User.ClientCertificateData == u.User.ClientCertificateData && user.User.ClientKeyData == u.User.ClientKeyData {
			unique = false
		}
	}

	return unique
}

func (k *kubeconf) uniqueUserEKS(u *Users) bool {
	unique := true
	for _, user := range k.Users {
		if user.Name == u.Name && user.User.Exec.APIVersion == u.User.Exec.APIVersion && user.User.Exec.Command == u.User.Exec.Command {
			unique = false
		}
	}

	return unique
}

func (k *kubeconf) addCluster(name string, ca string, server string) {
	newCluster := &Clusters{}
	newCluster.Name = name
	newCluster.Cluster.CertificateAuthorityData = ca
	newCluster.Cluster.Server = server

	if k.uniqueCluster(newCluster) {
		k.Clusters = append(k.Clusters, *newCluster)
	}
}

func (k *kubeconf) addContext(name string, cluster string, user string) {
	newContext := &Contexts{}
	newContext.Name = name
	newContext.Context.Cluster = cluster
	newContext.Context.User = user

	if k.uniqueContext(newContext) {
		k.Contexts = append(k.Contexts, *newContext)
	}
}

func (k *kubeconf) addUser(name string, cert string, key string) {
	newUser := &Users{}
	newUser.Name = name
	newUser.User.ClientCertificateData = cert
	newUser.User.ClientKeyData = key
	if k.uniqueUser(newUser) {
		k.Users = append(k.Users, *newUser)
	}
}

func (k *kubeconf) addUserEKS(clustername string, rolearn string, profile string, region string) {
	var args []string
	if rolearn != "" {
		args = []string{
			"eks",
			"get-token",
			"--cluster-name",
			clustername,
			"--role-arn",
			rolearn,
		}
	} else {
		args = []string{
			"eks",
			"get-token",
			"--cluster-name",
			clustername,
		}
	}

	regionenv := &Env{
		Name:  "AWS_DEFAULT_REGION",
		Value: region,
	}
	profileenv := &Env{
		Name:  "AWS_PROFILE",
		Value: profile,
	}
	newUser := &Users{}
	newUser.Name = profile + "-" + clustername
	newUser.User.Exec.APIVersion = "client.authentication.k8s.io/v1alpha1"
	newUser.User.Exec.Args = args
	newUser.User.Exec.Command = "aws"
	newUser.User.Exec.Env = append(newUser.User.Exec.Env, *regionenv)
	newUser.User.Exec.Env = append(newUser.User.Exec.Env, *profileenv)
	if k.uniqueUserEKS(newUser) {
		k.Users = append(k.Users, *newUser)
	}
}
