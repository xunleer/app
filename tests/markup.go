package tests

import (
	"testing"

	"github.com/murlokswarm/app"
)

// TestMarkup is a test suite used to ensure that all markups implementations
// behave the same.
func TestMarkup(t *testing.T, newMarkup func(factory app.Factory) app.Markup) {
	factory := app.NewFactory()
	factory.RegisterComponent(&Foo{})
	factory.RegisterComponent(&Bar{})
	factory.RegisterComponent(&CompoWithBadTmpl{})
	factory.RegisterComponent(&CompoWithBadTag{})
	factory.RegisterComponent(&CompoWithNotRegisteredChild{})
	factory.RegisterComponent(&CompoWithBadChild{})
	factory.RegisterComponent(&Hello{})
	factory.RegisterComponent(&World{})

	tests := []struct {
		scenario string
		function func(t *testing.T, markup app.Markup)
	}{
		{
			scenario: "mounting and dismounting a component",
			function: testMarkupMountDismount,
		},
		{
			scenario: "contains the mounted component",
			function: testMarkupContains,
		},
		{
			scenario: "get the component root",
			function: testMarkupRoot,
		},
		{
			scenario: "mounting a mounted component returns an error",
			function: testMarkupMountMounted,
		},
		{
			scenario: "mounting a component with a bad template returns an error",
			function: testMarkupMountComponentWithBadTemplate,
		},
		{
			scenario: "mounting a component with a bad tag returns an error",
			function: testMarkuptMountComponentWithBadTag,
		},
		{
			scenario: "mounting a component with a not registered child returns an error",
			function: testMarkuptMountComponentWithNotRegistedChild,
		},
		{
			scenario: "mounting a component with bad attributes returns an error",
			function: testMarkuptMountComponentWithBadAttrs,
		},
		{
			scenario: "skip dismounting a dismounted component",
			function: testMarkupDismountDismounted,
		},
		{
			scenario: "skip dismounting a component with dismounted child",
			function: testMarkupDismountDismountedChild,
		},
		{
			scenario: "update does not trigger changes",
			function: testMarkupUpdateNoChanges,
		},
		{
			scenario: "updating text",
			function: testMarkupUpdateText,
		},
		{
			scenario: "updating simple tag to component",
			function: testMarkupUpdateSimpleToCompo,
		},
		{
			scenario: "updating simple tag to text",
			function: testMarkupUpdateSimpleToText,
		},
		{
			scenario: "updating text to simple tag",
			function: testMarkupUpdateTextToSimple,
		},
		{
			scenario: "updating component",
			function: testMarkupUpdateComponent,
		},
		{
			scenario: "skip updating an unchanged component",
			function: testMarkupUpdateComponentNoChange,
		},
		{
			scenario: "updating attributes",
			function: testMarkupUpdateUpdateAttributes,
		},
		{
			scenario: "updating a not mounted component returns an error",
			function: testMarkupUpdateUpdateNotMountedComponent,
		},
		{
			scenario: "updating a component with bad template returns an error",
			function: testMarkupUpdateComponentWithBadTemplate,
		},
		{
			scenario: "updating a component with bad child returns an error",
			function: testMarkupUpdateComponentWithBadChild,
		},
		{
			scenario: "updating a component with an error returns an error",
			function: testMarkupUpdateComponentWithError,
		},
		{
			scenario: "updating a tag with bad attribute returns an error",
			function: testMarkupUpdateBadAttribute,
		},
		{
			scenario: "updating a component with dismounted child returns an error",
			function: testMarkupUpdateComponentWithDismountedChild,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			test.function(t, newMarkup(factory))
		})
	}
}

func testMarkupMountDismount(t *testing.T, markup app.Markup) {
	compo := &Foo{}

	root, err := markup.Mount(compo)
	if err != nil {
		t.Fatal(err)
	}
	if count := markup.Len(); count != 2 {
		t.Fatal("markup doesn't have 2 components:", count)
	}

	barTag := root.Children[1]
	if name := barTag.Name; name != "tests.bar" {
		t.Fatalf("bar tag is not a tests.bar: %s", name)
	}
	if _, err = markup.Component(barTag.ID); err != nil {
		t.Fatal(err)
	}

	markup.Dismount(compo)
	if count := markup.Len(); count != 0 {
		t.Fatal("markup have components")
	}
}

func testMarkupContains(t *testing.T, markup app.Markup) {
	compo := &Hello{}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}
	if !markup.Contains(compo) {
		t.Error("markup doesn't contrain the mounted component")
	}
}

func testMarkupRoot(t *testing.T, markup app.Markup) {
	compo := &Hello{}
	root, err := markup.Mount(compo)
	if err != nil {
		t.Fatal(err)
	}

	var root2 app.Tag
	if root2, err = markup.Root(compo); err != nil {
		t.Fatal(err)
	}
	if root2.ID != root.ID {
		t.Error("root and root 2 doesn't have the same id")
	}
}

func testMarkupMountMounted(t *testing.T, markup app.Markup) {
	compo := &Foo{}

	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	_, err := markup.Mount(compo)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMarkupMountComponentWithBadTemplate(t *testing.T, markup app.Markup) {
	testMarkuptMountInvalidComponent(t, markup, &CompoWithBadTmpl{})

}

func testMarkuptMountInvalidComponent(t *testing.T, markup app.Markup, compo app.Component) {
	_, err := markup.Mount(compo)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMarkuptMountComponentWithBadTag(t *testing.T, markup app.Markup) {
	testMarkuptMountInvalidComponent(t, markup, &CompoWithBadTag{})

}

func testMarkuptMountComponentWithNotRegistedChild(t *testing.T, markup app.Markup) {
	testMarkuptMountInvalidComponent(t, markup, &CompoWithNotRegisteredChild{})

}

func testMarkuptMountComponentWithBadAttrs(t *testing.T, markup app.Markup) {
	testMarkuptMountInvalidComponent(t, markup, &CompoWithBadAttrs{})
}

func testMarkupDismountDismounted(t *testing.T, markup app.Markup) {
	compo := &Foo{}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}
	markup.Dismount(compo)
	markup.Dismount(compo)
}

func testMarkupDismountDismountedChild(t *testing.T, markup app.Markup) {
	compo := &Foo{}
	root, err := markup.Mount(compo)
	if err != nil {
		t.Fatal(err)
	}

	barTag := root.Children[1]

	var bar app.Component
	if bar, err = markup.Component(barTag.ID); err != nil {
		t.Fatal(err)
	}

	markup.Dismount(bar)
	markup.Dismount(compo)
}

func testMarkupUpdateNoChanges(t *testing.T, markup app.Markup) {
	compo := &Hello{}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	syncs, err := markup.Update(compo)
	if err != nil {
		t.Fatal(err)
	}
	if len(syncs) != 0 {
		t.Error("syncs is not empty:", len(syncs))
	}
}

func testMarkupUpdateText(t *testing.T, markup app.Markup) {
	compo := &Hello{Greeting: "Hi"}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	compo.Greeting = "Hello"

	syncs, err := markup.Update(compo)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(syncs); l != 1 {
		t.Fatal("syncs doesn't have 1 element:", l)
	}

	sync := syncs[0]
	if !sync.Replace {
		t.Error("sync is not a replace")
	}

	h1 := sync.Tag
	if h1.Name != "h1" {
		t.Fatal("tag updated is not a h1:", h1.Name)
	}

	if text := h1.Children[0]; text.Text != compo.Greeting {
		t.Errorf(`text is not "%s": "%s"`, compo.Greeting, text.Text)
	}
}

func testMarkupUpdateSimpleToCompo(t *testing.T, markup app.Markup) {
	compo := &Hello{}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	compo.Name = "Maxence"

	syncs, err := markup.Update(compo)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(syncs); l != 1 {
		t.Fatal("syncs doesn't have 1 element:", l)
	}

	sync := syncs[0]
	if !sync.Replace {
		t.Error("sync is not a replace")
	}

	world := sync.Tag
	if world.Name != "tests.world" {
		t.Fatal("tag updated is not a component tests.world:", world.Name)
	}
	if name := world.Attributes["name"]; name != compo.Name {
		t.Fatalf(`name is not "%s": "%s"`, compo.Name, name)
	}
	if l := len(world.Children); l != 0 {
		t.Fatal("world has children", l)
	}
}

func testMarkupUpdateSimpleToText(t *testing.T, markup app.Markup) {
	compo := &Hello{}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	compo.TextBye = true

	syncs, err := markup.Update(compo)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(syncs); l != 1 {
		t.Fatal("syncs doesn't have 1 element:", l)
	}

	sync := syncs[0]
	if !sync.Replace {
		t.Error("sync is not a replace")
	}

	root := sync.Tag
	if root.Name != "div" {
		t.Fatal("root is not a div:", root.Name)
	}
	if l := len(root.Children); l != 4 {
		t.Fatal("root doesn't have 4 children:", l)
	}
	if text := root.Children[3]; text.Text != "Goodbye" {
		t.Fatalf(`text is not "Goodbye": "%s"`, text.Text)
	}
}

func testMarkupUpdateTextToSimple(t *testing.T, markup app.Markup) {
	compo := &Hello{TextBye: true}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	compo.TextBye = false

	syncs, err := markup.Update(compo)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(syncs); l != 1 {
		t.Fatal("syncs doesn't have 1 element:", l)
	}

	sync := syncs[0]
	if !sync.Replace {
		t.Error("sync is not a replace")
	}

	root := sync.Tag
	if l := len(root.Children); l != 5 {
		t.Fatal("root doesn't have 5 children:", l)
	}
	if span := root.Children[3]; span.Name != "span" {
		t.Fatalf(`span is not a span tag: %s`, span.Name)
	}
	if p := root.Children[4]; p.Name != "p" {
		t.Fatalf(`p is not a p tag: %s`, p.Name)
	}
}

func testMarkupUpdateComponent(t *testing.T, markup app.Markup) {
	compo := &Hello{Name: "Jonhy"}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	compo.Name = "Maxence"

	syncs, err := markup.Update(compo)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(syncs); l != 1 {
		t.Fatal("syncs doesn't have 1 element:", l)
	}

	sync := syncs[0]
	if !sync.Replace {
		t.Error("sync is not a replace")
	}

	worldRoot := sync.Tag
	if worldRoot.Name != "div" {
		t.Fatal("root of world is not a div:", worldRoot.Name)
	}
	if l := len(worldRoot.Children); l != 1 {
		t.Fatal("root of world doesn't have 1 child:", l)
	}
	if text := worldRoot.Children[0]; text.Text != compo.Name {
		t.Fatalf(`text is not "%s": "%s"`, compo.Name, text.Text)
	}
}

func testMarkupUpdateComponentNoChange(t *testing.T, markup app.Markup) {
	compo := &Hello{Name: "JonhyMaxoo"}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	syncs, err := markup.Update(compo)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(syncs); l != 0 {
		t.Error("syncs is not empty:", l)
	}
}

func testMarkupUpdateUpdateAttributes(t *testing.T, markup app.Markup) {
	compo := &Hello{}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	compo.Placeholder = "Enter your name"

	syncs, err := markup.Update(compo)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(syncs); l != 1 {
		t.Fatal("syncs doesn't have 1 element:", l)
	}

	sync := syncs[0]
	if sync.Replace {
		t.Error("sync is a replace")
	}

	input := sync.Tag
	if input.Name != "input" {
		t.Error("input is not an input tag:", input.Name)
	}
	if placeholder := input.Attributes["placeholder"]; placeholder != compo.Placeholder {
		t.Errorf("input placeholder is not %s: %s", compo.Placeholder, placeholder)
	}
	if l := len(input.Children); l != 0 {
		t.Error("input has child")
	}
}

func testMarkupUpdateUpdateNotMountedComponent(t *testing.T, markup app.Markup) {
	_, err := markup.Update(&Hello{})
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMarkupUpdateComponentWithBadTemplate(t *testing.T, markup app.Markup) {
	compo := &Hello{}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	compo.TmplErr = true

	_, err := markup.Update(compo)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMarkupUpdateComponentWithBadChild(t *testing.T, markup app.Markup) {
	compo := &Hello{Name: "Max"}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	compo.ChildErr = true

	_, err := markup.Update(compo)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMarkupUpdateComponentWithError(t *testing.T, markup app.Markup) {
	compo := &Hello{}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	compo.Name = "Jonhy"
	compo.ChildErr = true

	_, err := markup.Update(compo)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMarkupUpdateBadAttribute(t *testing.T, markup app.Markup) {
	compo := &Hello{Name: "Maxoo"}
	if _, err := markup.Mount(compo); err != nil {
		t.Fatal(err)
	}

	compo.CompoFieldErr = true

	_, err := markup.Update(compo)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testMarkupUpdateComponentWithDismountedChild(t *testing.T, markup app.Markup) {
	compo := &Hello{Name: "Maxoo"}
	root, err := markup.Mount(compo)
	if err != nil {
		t.Fatal(err)
	}

	worldTag := root.Children[2].Children[0]

	var world app.Component
	if world, err = markup.Component(worldTag.ID); err != nil {
		t.Fatal(err)
	}

	markup.Dismount(world)

	compo.Name = "Jonhy"

	if _, err = markup.Update(compo); err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}