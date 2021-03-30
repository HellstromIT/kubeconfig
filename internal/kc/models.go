package kc

type kubeconf struct {
	APIVersion     string      `yaml:"apiVersion"`
	Clusters       []Clusters  `yaml:"clusters"`
	Contexts       []Contexts  `yaml:"contexts"`
	CurrentContext string      `yaml:"current-context"`
	Kind           string      `yaml:"kind"`
	Preferences    Preferences `yaml:"preferences"`
	Users          []Users     `yaml:"users"`
}

type Clusters struct {
	Cluster Cluster `yaml:"cluster"`
	Name    string  `yaml:"name"`
}

type Cluster struct {
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
	Server                   string `yaml:"server"`
}

type Context struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

type Contexts struct {
	Context Context `yaml:"context"`
	Name    string  `yaml:"name"`
}

type Preferences struct {
}

type Env struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Exec struct {
	APIVersion         string   `yaml:"apiVersion"`
	Args               []string `yaml:"args"`
	Command            string   `yaml:"command"`
	Env                []Env    `yaml:"env"`
	ProvideClusterInfo bool     `yaml:"provideClusterInfo"`
}

type User struct {
	Exec                  Exec   `yaml:"exec,omitempty"`
	ClientCertificateData string `yaml:"client-certificate-data,omitempty"`
	ClientKeyData         string `yaml:"client-key-data,omitempty"`
}

type Users struct {
	Name string `yaml:"name"`
	User User   `yaml:"user,omitempty`
}
