package main

func Build(args []string) {
	cliFlags := NewCliFlags(true)
	cliFlags.UsageHeader = "kubedev build <context directory>"
	cliFlags.Parse(args)
	context := cliFlags.NewContext()

	context.BuildAllDockerImages()
	// kubeCtl := context.MakeKubeCtl()

	// if kubernetesContexts, err := kubeCtl.GetContexts(); err == nil {
	// 	log.Printf("Kubernetes contexts: %v", kubernetesContexts)
	// } else {
	// 	log.Fatal(err)
	// }
	// if currentKubeContext, err := kubeCtl.GetCurrentContext(); err == nil {
	// 	log.Printf("Current kubernetes context: %v", currentKubeContext)
	// } else {
	// 	log.Fatal(err)
	// }
}
