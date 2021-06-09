package plugins

// Plugin is an executable part that can be invoked separately
// One usecase of plugin is post-deploy hook
// You may want to transfer kube config file to somewhere after
// you execute mediumkube init on a certain node
// Or you want to open a port on mediumkube bridge in order to
// start some services which listen on that bridge
type Plugin interface {
	Exec(args ...string)
	Desc()
}
