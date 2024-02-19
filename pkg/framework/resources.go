package framework

import (
	"log"

	"github.com/blushft/go-diagrams/diagram"
	"github.com/blushft/go-diagrams/nodes/k8s"
)

type Component interface {
	Render(d *diagram.Diagram, app, user *diagram.Group)
}

type Framework interface {
	Render(...Component)
}

type Opendatahub struct{}

func (f *Opendatahub) Render(components ...Component) {
	d, err := diagram.New(diagram.Filename("app"), diagram.Label("Openshift AI"), diagram.Direction("LR"))
	if err != nil {
		log.Fatal(err)
	}

	dc := diagram.NewGroup("Cluster")

	app := dc.NewGroup("app").Label("Application")
	user := dc.NewGroup("user").Label("DataScienceProject")
	d.Group(dc)
	for _, comp := range components {
		comp.Render(d, app, user)
	}

	if err := d.Render(); err != nil {
		log.Fatal(err)
	}
}

type Dashboard struct{}

func (c *Dashboard) Render(d *diagram.Diagram, app, user *diagram.Group) {
	app.Add(k8s.Compute.Pod(diagram.NodeLabel("Dashboard")))
}

type Notebook struct{}

func (c *Notebook) Render(d *diagram.Diagram, app, user *diagram.Group) {
	notebookCtrl := k8s.Compute.Pod(diagram.NodeLabel("Notebook Ctrl"))
	notebook := k8s.Compute.Pod(diagram.NodeLabel("Notebook"))
	app.Add(notebookCtrl)
	user.Add(notebook)
	d.Connect(notebookCtrl, notebook, diagram.Forward())
}

type KServe struct{}

func (c *KServe) Render(d *diagram.Diagram, app, user *diagram.Group) {

	kservingNs := d.Groups()[0].NewGroup("kserving").Label("KNative Serving")
	activator := k8s.Compute.Pod(diagram.NodeLabel("Activator"))
	kservingNs.Add(activator)

	modelCtrl := k8s.Compute.Pod(diagram.NodeLabel("Model Ctrl"))
	kServeCtrl := k8s.Compute.Pod(diagram.NodeLabel("KServe Ctrl"))
	inferenceService := k8s.Compute.Pod(diagram.NodeLabel("Inference Service"))
	servingRuntime := k8s.Compute.Pod(diagram.NodeLabel("Serving Runtime"))

	app.Add(modelCtrl)
	app.Add(kServeCtrl)
	user.Add(inferenceService)
	user.Add(servingRuntime)
}
