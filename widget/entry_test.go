package widget_test

import (
	"fmt"
	"image/color"
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/data/binding"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

func TestEntry_Binding(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("Init")
	assert.Equal(t, "Init", entry.Text)

	str := binding.NewString()
	entry.Bind(str)
	waitForBinding()
	assert.Equal(t, "", entry.Text)

	err := str.Set("Updated")
	assert.Nil(t, err)
	waitForBinding()
	assert.Equal(t, "Updated", entry.Text)

	entry.SetText("Typed")
	v, err := str.Get()
	assert.Nil(t, err)
	assert.Equal(t, "Typed", v)

	entry.Unbind()
	waitForBinding()
	assert.Equal(t, "Typed", entry.Text)
}

func TestEntry_CursorColumn(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("")
	assert.Equal(t, 0, entry.CursorColumn)

	// only 0 columns, do nothing
	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.TypedKey(right)
	assert.Equal(t, 0, entry.CursorColumn)

	// 1, this should increment
	entry.SetText("a")
	entry.TypedKey(right)
	assert.Equal(t, 1, entry.CursorColumn)

	left := &fyne.KeyEvent{Name: fyne.KeyLeft}
	entry.TypedKey(left)
	assert.Equal(t, 0, entry.CursorColumn)

	// don't go beyond left
	entry.TypedKey(left)
	assert.Equal(t, 0, entry.CursorColumn)
}

func TestEntry_CursorColumn_Jump(t *testing.T) {
	entry := widget.NewMultiLineEntry()
	entry.SetText("a\nbc")

	// go to end of text
	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.TypedKey(right)
	entry.TypedKey(right)
	entry.TypedKey(right)
	entry.TypedKey(right)
	assert.Equal(t, 1, entry.CursorRow)
	assert.Equal(t, 2, entry.CursorColumn)

	// go up, to a shorter line
	up := &fyne.KeyEvent{Name: fyne.KeyUp}
	entry.TypedKey(up)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestEntry_CursorColumn_Wrap(t *testing.T) {
	entry := widget.NewMultiLineEntry()
	entry.SetText("a\nb")
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)

	// go to end of line
	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.TypedKey(right)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	// wrap to new line
	entry.TypedKey(right)
	assert.Equal(t, 1, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)

	// and back
	left := &fyne.KeyEvent{Name: fyne.KeyLeft}
	entry.TypedKey(left)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestEntry_CursorPasswordRevealer(t *testing.T) {
	pr := widget.NewPasswordEntry().ActionItem.(desktop.Cursorable)
	assert.Equal(t, desktop.DefaultCursor, pr.Cursor())
}

func TestEntry_CursorRow(t *testing.T) {
	entry := widget.NewMultiLineEntry()
	entry.SetText("test")
	assert.Equal(t, 0, entry.CursorRow)

	// only 1 line, do nothing
	down := &fyne.KeyEvent{Name: fyne.KeyDown}
	entry.TypedKey(down)
	assert.Equal(t, 0, entry.CursorRow)

	// 2 lines, this should increment
	entry.SetText("test\nrows")
	entry.TypedKey(down)
	assert.Equal(t, 1, entry.CursorRow)

	up := &fyne.KeyEvent{Name: fyne.KeyUp}
	entry.TypedKey(up)
	assert.Equal(t, 0, entry.CursorRow)

	// don't go beyond top
	entry.TypedKey(up)
	assert.Equal(t, 0, entry.CursorRow)
}

func TestEntry_Disableable(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.False(t, entry.Disabled())

	entry.Disable()
	assert.True(t, entry.Disabled())
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="disabled" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text color="disabled" pos="4,4" size="104x21"></text>
							</widget>
							<widget size="112x29" type="*widget.textProvider">
								<text color="disabled" pos="4,4" size="104x21"></text>
							</widget>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)

	entry.Enable()
	assert.False(t, entry.Disabled())
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="shadow" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text color="placeholder" pos="4,4" size="104x21"></text>
							</widget>
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21"></text>
							</widget>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)

	entry.SetPlaceHolder("Type!")
	assert.False(t, entry.Disabled())
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="shadow" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text color="placeholder" pos="4,4" size="104x21">Type!</text>
							</widget>
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21"></text>
							</widget>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)

	entry.Disable()
	assert.True(t, entry.Disabled())
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="disabled" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text color="disabled" pos="4,4" size="104x21">Type!</text>
							</widget>
							<widget size="112x29" type="*widget.textProvider">
								<text color="disabled" pos="4,4" size="104x21"></text>
							</widget>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)

	entry.Enable()
	assert.False(t, entry.Disabled())
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="shadow" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text color="placeholder" pos="4,4" size="104x21">Type!</text>
							</widget>
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21"></text>
							</widget>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)

	entry.SetText("Hello")
	assert.False(t, entry.Disabled())
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="shadow" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Hello</text>
							</widget>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)

	entry.Disable()
	assert.True(t, entry.Disabled())
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="disabled" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text color="disabled" pos="4,4" size="104x21">Hello</text>
							</widget>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)

	entry.Enable()
	assert.False(t, entry.Disabled())
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="shadow" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Hello</text>
							</widget>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
}

func TestEntry_EmptySelection(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("text")

	// trying to select at the edge
	typeKeys(entry, keyShiftLeftDown, fyne.KeyLeft, keyShiftLeftUp)
	assert.Equal(t, "", entry.SelectedText())

	typeKeys(entry, fyne.KeyRight)
	assert.Equal(t, 1, entry.CursorColumn)

	// stop selecting at the edge when nothing is selected
	typeKeys(entry, fyne.KeyLeft, keyShiftLeftDown, fyne.KeyRight, fyne.KeyLeft, keyShiftLeftUp)
	assert.Equal(t, "", entry.SelectedText())
	assert.Equal(t, 0, entry.CursorColumn)

	// check that the selection has been removed
	typeKeys(entry, fyne.KeyRight, keyShiftLeftDown, fyne.KeyRight, fyne.KeyLeft, keyShiftLeftUp)
	assert.Equal(t, "", entry.SelectedText())
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestEntry_Focus(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.FocusGained()
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text color="placeholder" pos="4,4" size="104x21"></text>
							</widget>
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21"></text>
							</widget>
							<rectangle fillColor="focus" pos="3,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)

	entry.FocusLost()
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="shadow" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text color="placeholder" pos="4,4" size="104x21"></text>
							</widget>
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21"></text>
							</widget>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)

	window.Canvas().Focus(entry)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text color="placeholder" pos="4,4" size="104x21"></text>
							</widget>
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21"></text>
							</widget>
							<rectangle fillColor="focus" pos="3,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
}

func TestEntry_FocusWithPopUp(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	test.TapSecondaryAt(entry, fyne.NewPos(1, 1))
	fmt.Println(theme.FocusColor())
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text color="placeholder" pos="4,4" size="104x21"></text>
							</widget>
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21"></text>
							</widget>
							<rectangle fillColor="focus" pos="3,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
			<overlay>
				<widget size="150x200" type="*widget.OverlayContainer">
					<widget pos="11,11" size="79x136" type="*widget.Menu">
						<widget size="79x136" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="79x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="79,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="79,0" size="4x136" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="79,136" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,136" size="79x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,136" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x136"/>
						</widget>
						<widget size="79x136" type="*widget.Scroll">
							<widget size="79x136" type="*widget.menuBox">
								<container pos="0,4" size="79x144">
									<widget size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="25x21">Cut</text>
									</widget>
									<widget pos="0,33" size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="36x21">Copy</text>
									</widget>
									<widget pos="0,66" size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="39x21">Paste</text>
									</widget>
									<widget pos="0,99" size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Select all</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, c)

	test.TapCanvas(c, fyne.NewPos(20, 20))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text color="placeholder" pos="4,4" size="104x21"></text>
							</widget>
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21"></text>
							</widget>
							<rectangle fillColor="focus" pos="3,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)

	test.TapSecondaryAt(entry, fyne.NewPos(1, 1))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text color="placeholder" pos="4,4" size="104x21"></text>
							</widget>
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21"></text>
							</widget>
							<rectangle fillColor="focus" pos="3,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
			<overlay>
				<widget size="150x200" type="*widget.OverlayContainer">
					<widget pos="11,11" size="79x136" type="*widget.Menu">
						<widget size="79x136" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="79x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="79,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="79,0" size="4x136" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="79,136" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,136" size="79x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,136" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x136"/>
						</widget>
						<widget size="79x136" type="*widget.Scroll">
							<widget size="79x136" type="*widget.menuBox">
								<container pos="0,4" size="79x144">
									<widget size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="25x21">Cut</text>
									</widget>
									<widget pos="0,33" size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="36x21">Copy</text>
									</widget>
									<widget pos="0,66" size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="39x21">Paste</text>
									</widget>
									<widget pos="0,99" size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Select all</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, c)

	test.TapCanvas(c, fyne.NewPos(5, 5))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text color="placeholder" pos="4,4" size="104x21"></text>
							</widget>
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21"></text>
							</widget>
							<rectangle fillColor="focus" pos="3,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
}

func TestEntry_HidePopUpOnEntry(t *testing.T) {
	entry := widget.NewEntry()
	tapPos := fyne.NewPos(1, 1)
	c := fyne.CurrentApp().Driver().CanvasForObject(entry)

	assert.Nil(t, c.Overlays().Top())

	test.TapSecondaryAt(entry, tapPos)
	assert.NotNil(t, c.Overlays().Top())

	test.Type(entry, "KJGFD")
	assert.Nil(t, c.Overlays().Top())
	assert.Equal(t, "KJGFD", entry.Text)
}

func TestEntry_MinSize(t *testing.T) {
	entry := widget.NewEntry()
	min := entry.MinSize()
	entry.SetPlaceHolder("")
	assert.Equal(t, min, entry.MinSize())
	entry.SetText("")
	assert.Equal(t, min, entry.MinSize())
	entry.SetPlaceHolder("Hi")
	assert.Equal(t, entry.MinSize().Width, min.Width)
	assert.Equal(t, entry.MinSize().Height, min.Height)

	assert.True(t, min.Width > theme.Padding()*2)
	assert.True(t, min.Height > theme.Padding()*2)

	min = entry.MinSize()
	entry.ActionItem = canvas.NewCircle(color.Black)
	assert.Equal(t, min.Add(fyne.NewSize(theme.IconInlineSize()+theme.Padding(), 0)), entry.MinSize())
}

func TestEntry_MultilineSelect(t *testing.T) {
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	// Extend the selection down one row
	typeKeys(e, keyShiftLeftDown, fyne.KeyDown)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<rectangle fillColor="focus" pos="20,25" size="37x21"/>
							<rectangle fillColor="focus" pos="3,46" size="35x21"/>
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Testing</text>
								<text pos="4,25" size="104x21">Testing</text>
								<text pos="4,46" size="104x21">Testing</text>
							</widget>
							<rectangle fillColor="focus" pos="37,46" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Equal(t, "sting\nTesti", e.SelectedText())

	typeKeys(e, fyne.KeyUp)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<rectangle fillColor="focus" pos="20,25" size="18x21"/>
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Testing</text>
								<text pos="4,25" size="104x21">Testing</text>
								<text pos="4,46" size="104x21">Testing</text>
							</widget>
							<rectangle fillColor="focus" pos="37,25" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Equal(t, "sti", e.SelectedText())

	typeKeys(e, fyne.KeyUp)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<rectangle fillColor="focus" pos="37,4" size="20x21"/>
							<rectangle fillColor="focus" pos="3,25" size="18x21"/>
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Testing</text>
								<text pos="4,25" size="104x21">Testing</text>
								<text pos="4,46" size="104x21">Testing</text>
							</widget>
							<rectangle fillColor="focus" pos="37,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Equal(t, "ng\nTe", e.SelectedText())
}

func TestEntry_Notify(t *testing.T) {
	entry := widget.NewEntry()
	changed := false

	entry.OnChanged = func(string) {
		changed = true
	}
	entry.SetText("Test")

	assert.True(t, changed)
}

func TestEntry_OnCopy(t *testing.T) {
	e := widget.NewEntry()
	e.SetText("Testing")
	typeKeys(e, fyne.KeyRight, fyne.KeyRight, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight, fyne.KeyRight)

	clipboard := test.NewClipboard()
	shortcut := &fyne.ShortcutCopy{Clipboard: clipboard}
	e.TypedShortcut(shortcut)

	assert.Equal(t, "sti", clipboard.Content())
	assert.Equal(t, "Testing", e.Text)
}

func TestEntry_OnCopy_Password(t *testing.T) {
	e := widget.NewPasswordEntry()
	e.SetText("Testing")
	typeKeys(e, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight, fyne.KeyRight)

	clipboard := test.NewClipboard()
	shortcut := &fyne.ShortcutCopy{Clipboard: clipboard}
	e.TypedShortcut(shortcut)

	assert.Equal(t, "", clipboard.Content())
	assert.Equal(t, "Testing", e.Text)
}

func TestEntry_OnCut(t *testing.T) {
	e := widget.NewEntry()
	e.SetText("Testing")
	typeKeys(e, fyne.KeyRight, fyne.KeyRight, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight, fyne.KeyRight)

	clipboard := test.NewClipboard()
	shortcut := &fyne.ShortcutCut{Clipboard: clipboard}
	e.TypedShortcut(shortcut)

	assert.Equal(t, "sti", clipboard.Content())
	assert.Equal(t, "Teng", e.Text)
}

func TestEntry_OnCut_Password(t *testing.T) {
	e := widget.NewPasswordEntry()
	e.SetText("Testing")
	typeKeys(e, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight, fyne.KeyRight)

	clipboard := test.NewClipboard()
	shortcut := &fyne.ShortcutCut{Clipboard: clipboard}
	e.TypedShortcut(shortcut)

	assert.Equal(t, "", clipboard.Content())
	assert.Equal(t, "Testing", e.Text)
}

func TestEntry_OnKeyDown(t *testing.T) {
	entry := widget.NewEntry()

	test.Type(entry, "Hi")

	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_OnKeyDown_Backspace(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("Hi")
	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.TypedKey(right)
	entry.TypedKey(right)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 2, entry.CursorColumn)

	key := &fyne.KeyEvent{Name: fyne.KeyBackspace}
	entry.TypedKey(key)

	assert.Equal(t, "H", entry.Text)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestEntry_OnKeyDown_BackspaceBeyondText(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("Hi")
	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.TypedKey(right)
	entry.TypedKey(right)

	key := &fyne.KeyEvent{Name: fyne.KeyBackspace}
	entry.TypedKey(key)
	entry.TypedKey(key)
	entry.TypedKey(key)

	assert.Equal(t, "", entry.Text)
}

func TestEntry_OnKeyDown_BackspaceBeyondTextAndNewLine(t *testing.T) {
	entry := widget.NewMultiLineEntry()
	entry.SetText("H\ni")

	down := &fyne.KeyEvent{Name: fyne.KeyDown}
	entry.TypedKey(down)
	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.TypedKey(right)

	key := &fyne.KeyEvent{Name: fyne.KeyBackspace}
	entry.TypedKey(key)

	assert.Equal(t, 1, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)
	entry.TypedKey(key)

	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
	assert.Equal(t, "H", entry.Text)
}

func TestEntry_OnKeyDown_BackspaceNewline(t *testing.T) {
	entry := widget.NewMultiLineEntry()
	entry.SetText("H\ni")

	down := &fyne.KeyEvent{Name: fyne.KeyDown}
	entry.TypedKey(down)

	key := &fyne.KeyEvent{Name: fyne.KeyBackspace}
	entry.TypedKey(key)

	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_OnKeyDown_BackspaceUnicode(t *testing.T) {
	entry := widget.NewEntry()

	test.Type(entry, "è")
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	bs := &fyne.KeyEvent{Name: fyne.KeyBackspace}
	entry.TypedKey(bs)
	assert.Equal(t, "", entry.Text)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)
}

func TestEntry_OnKeyDown_Delete(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("Hi")
	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.TypedKey(right)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	key := &fyne.KeyEvent{Name: fyne.KeyDelete}
	entry.TypedKey(key)

	assert.Equal(t, "H", entry.Text)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestEntry_OnKeyDown_DeleteBeyondText(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("Hi")

	key := &fyne.KeyEvent{Name: fyne.KeyDelete}
	entry.TypedKey(key)
	entry.TypedKey(key)
	entry.TypedKey(key)

	assert.Equal(t, "", entry.Text)
}

func TestEntry_OnKeyDown_DeleteNewline(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("H\ni")

	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.TypedKey(right)

	key := &fyne.KeyEvent{Name: fyne.KeyDelete}
	entry.TypedKey(key)

	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_OnKeyDown_HomeEnd(t *testing.T) {
	entry := &widget.Entry{}
	entry.SetText("Hi")
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)

	end := &fyne.KeyEvent{Name: fyne.KeyEnd}
	entry.TypedKey(end)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 2, entry.CursorColumn)

	home := &fyne.KeyEvent{Name: fyne.KeyHome}
	entry.TypedKey(home)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)
}

func TestEntry_OnKeyDown_Insert(t *testing.T) {
	entry := widget.NewEntry()

	test.Type(entry, "Hi")
	assert.Equal(t, "Hi", entry.Text)

	left := &fyne.KeyEvent{Name: fyne.KeyLeft}
	entry.TypedKey(left)

	test.Type(entry, "o")
	assert.Equal(t, "Hoi", entry.Text)
}

func TestEntry_OnKeyDown_Newline(t *testing.T) {
	entry, window := setupImageTest(t, true)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.SetText("Hi")
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="shadow" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Hi</text>
							</widget>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)

	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.TypedKey(right)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	key := &fyne.KeyEvent{Name: fyne.KeyReturn}
	entry.TypedKey(key)

	assert.Equal(t, "H\ni", entry.Text)
	assert.Equal(t, 1, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)

	test.Type(entry, "o")
	assert.Equal(t, "H\noi", entry.Text)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">H</text>
								<text pos="4,25" size="104x21">oi</text>
							</widget>
							<rectangle fillColor="focus" pos="13,25" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
}

func TestEntry_OnPaste(t *testing.T) {
	clipboard := test.NewClipboard()
	shortcut := &fyne.ShortcutPaste{Clipboard: clipboard}
	tests := []struct {
		name             string
		entry            *widget.Entry
		clipboardContent string
		wantText         string
		wantRow, wantCol int
	}{
		{
			name:             "singleline: empty content",
			entry:            widget.NewEntry(),
			clipboardContent: "",
			wantText:         "",
			wantRow:          0,
			wantCol:          0,
		},
		{
			name:             "singleline: simple text",
			entry:            widget.NewEntry(),
			clipboardContent: "clipboard content",
			wantText:         "clipboard content",
			wantRow:          0,
			wantCol:          17,
		},
		{
			name:             "singleline: UTF8 text",
			entry:            widget.NewEntry(),
			clipboardContent: "Hié™שרה",
			wantText:         "Hié™שרה",
			wantRow:          0,
			wantCol:          7,
		},
		{
			name:             "singleline: with new line",
			entry:            widget.NewEntry(),
			clipboardContent: "clipboard\ncontent",
			wantText:         "clipboard content",
			wantRow:          0,
			wantCol:          17,
		},
		{
			name:             "singleline: with tab",
			entry:            widget.NewEntry(),
			clipboardContent: "clipboard\tcontent",
			wantText:         "clipboard\tcontent",
			wantRow:          0,
			wantCol:          17,
		},
		{
			name:             "password: with new line",
			entry:            widget.NewPasswordEntry(),
			clipboardContent: "3SB=y+)z\nkHGK(hx6 -e_\"1TZu q^bF3^$u H[:e\"1O.",
			wantText:         `3SB=y+)z kHGK(hx6 -e_"1TZu q^bF3^$u H[:e"1O.`,
			wantRow:          0,
			wantCol:          44,
		},
		{
			name:             "multiline: with new line",
			entry:            widget.NewMultiLineEntry(),
			clipboardContent: "clipboard\ncontent",
			wantText:         "clipboard\ncontent",
			wantRow:          1,
			wantCol:          7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clipboard.SetContent(tt.clipboardContent)
			tt.entry.TypedShortcut(shortcut)
			assert.Equal(t, tt.wantText, tt.entry.Text)
			assert.Equal(t, tt.wantRow, tt.entry.CursorRow)
			assert.Equal(t, tt.wantCol, tt.entry.CursorColumn)
		})
	}
}

func TestEntry_PageUpDown(t *testing.T) {
	t.Run("single line", func(*testing.T) {
		e, window := setupImageTest(t, false)
		defer teardownImageTest(window)
		c := window.Canvas()

		c.Focus(e)
		e.SetText("Testing")
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<widget pos="10,10" size="120x37" type="*widget.Entry">
						<rectangle fillColor="focus" pos="0,33" size="120x4"/>
						<widget pos="4,4" size="112x29" type="*widget.Scroll">
							<widget size="112x29" type="*widget.entryContent">
								<widget size="112x29" type="*widget.textProvider">
									<text pos="4,4" size="104x21">Testing</text>
								</widget>
								<rectangle fillColor="focus" pos="3,4" size="2x21"/>
							</widget>
						</widget>
					</widget>
				</content>
			</canvas>
		`, c)

		// move right, press & hold shift and pagedown
		typeKeys(e, fyne.KeyRight, keyShiftLeftDown, fyne.KeyPageDown)
		assert.Equal(t, "esting", e.SelectedText())
		assert.Equal(t, 0, e.CursorRow)
		assert.Equal(t, 7, e.CursorColumn)
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<widget pos="10,10" size="120x37" type="*widget.Entry">
						<rectangle fillColor="focus" pos="0,33" size="120x4"/>
						<widget pos="4,4" size="112x29" type="*widget.Scroll">
							<widget size="112x29" type="*widget.entryContent">
								<rectangle fillColor="focus" pos="12,4" size="45x21"/>
								<widget size="112x29" type="*widget.textProvider">
									<text pos="4,4" size="104x21">Testing</text>
								</widget>
								<rectangle fillColor="focus" pos="56,4" size="2x21"/>
							</widget>
						</widget>
					</widget>
				</content>
			</canvas>
		`, c)

		// while shift is held press pageup
		typeKeys(e, fyne.KeyPageUp)
		assert.Equal(t, "T", e.SelectedText())
		assert.Equal(t, 0, e.CursorRow)
		assert.Equal(t, 0, e.CursorColumn)
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<widget pos="10,10" size="120x37" type="*widget.Entry">
						<rectangle fillColor="focus" pos="0,33" size="120x4"/>
						<widget pos="4,4" size="112x29" type="*widget.Scroll">
							<widget size="112x29" type="*widget.entryContent">
								<rectangle fillColor="focus" pos="3,4" size="10x21"/>
								<widget size="112x29" type="*widget.textProvider">
									<text pos="4,4" size="104x21">Testing</text>
								</widget>
								<rectangle fillColor="focus" pos="3,4" size="2x21"/>
							</widget>
						</widget>
					</widget>
				</content>
			</canvas>
		`, c)

		// release shift and press pagedown
		typeKeys(e, keyShiftLeftUp, fyne.KeyPageDown)
		assert.Equal(t, "", e.SelectedText())
		assert.Equal(t, 0, e.CursorRow)
		assert.Equal(t, 7, e.CursorColumn)
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<widget pos="10,10" size="120x37" type="*widget.Entry">
						<rectangle fillColor="focus" pos="0,33" size="120x4"/>
						<widget pos="4,4" size="112x29" type="*widget.Scroll">
							<widget size="112x29" type="*widget.entryContent">
								<widget size="112x29" type="*widget.textProvider">
									<text pos="4,4" size="104x21">Testing</text>
								</widget>
								<rectangle fillColor="focus" pos="56,4" size="2x21"/>
							</widget>
						</widget>
					</widget>
				</content>
			</canvas>
		`, c)
	})

	t.Run("page down single line", func(*testing.T) {
		e, window := setupImageTest(t, true)
		defer teardownImageTest(window)
		c := window.Canvas()

		c.Focus(e)
		e.SetText("Testing\nTesting\nTesting")
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<widget pos="10,10" size="120x100" type="*widget.Entry">
						<rectangle fillColor="focus" pos="0,96" size="120x4"/>
						<widget pos="4,4" size="112x92" type="*widget.Scroll">
							<widget size="112x92" type="*widget.entryContent">
								<widget size="112x92" type="*widget.textProvider">
									<text pos="4,4" size="104x21">Testing</text>
									<text pos="4,25" size="104x21">Testing</text>
									<text pos="4,46" size="104x21">Testing</text>
								</widget>
								<rectangle fillColor="focus" pos="3,4" size="2x21"/>
							</widget>
						</widget>
					</widget>
				</content>
			</canvas>
		`, c)

		// move right, press & hold shift and pagedown
		typeKeys(e, fyne.KeyRight, keyShiftLeftDown, fyne.KeyPageDown)
		assert.Equal(t, "esting\nTesting\nTesting", e.SelectedText())
		assert.Equal(t, 2, e.CursorRow)
		assert.Equal(t, 7, e.CursorColumn)
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<widget pos="10,10" size="120x100" type="*widget.Entry">
						<rectangle fillColor="focus" pos="0,96" size="120x4"/>
						<widget pos="4,4" size="112x92" type="*widget.Scroll">
							<widget size="112x92" type="*widget.entryContent">
								<rectangle fillColor="focus" pos="12,4" size="45x21"/>
								<rectangle fillColor="focus" pos="3,25" size="54x21"/>
								<rectangle fillColor="focus" pos="3,46" size="54x21"/>
								<widget size="112x92" type="*widget.textProvider">
									<text pos="4,4" size="104x21">Testing</text>
									<text pos="4,25" size="104x21">Testing</text>
									<text pos="4,46" size="104x21">Testing</text>
								</widget>
								<rectangle fillColor="focus" pos="56,46" size="2x21"/>
							</widget>
						</widget>
					</widget>
				</content>
			</canvas>
		`, c)

		// while shift is held press pageup
		typeKeys(e, fyne.KeyPageUp)
		assert.Equal(t, "T", e.SelectedText())
		assert.Equal(t, 0, e.CursorRow)
		assert.Equal(t, 0, e.CursorColumn)
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<widget pos="10,10" size="120x100" type="*widget.Entry">
						<rectangle fillColor="focus" pos="0,96" size="120x4"/>
						<widget pos="4,4" size="112x92" type="*widget.Scroll">
							<widget size="112x92" type="*widget.entryContent">
								<rectangle fillColor="focus" pos="3,4" size="10x21"/>
								<widget size="112x92" type="*widget.textProvider">
									<text pos="4,4" size="104x21">Testing</text>
									<text pos="4,25" size="104x21">Testing</text>
									<text pos="4,46" size="104x21">Testing</text>
								</widget>
								<rectangle fillColor="focus" pos="3,4" size="2x21"/>
							</widget>
						</widget>
					</widget>
				</content>
			</canvas>
		`, c)

		// release shift and press pagedown
		typeKeys(e, keyShiftLeftUp, fyne.KeyPageDown)
		assert.Equal(t, "", e.SelectedText())
		assert.Equal(t, 2, e.CursorRow)
		assert.Equal(t, 7, e.CursorColumn)
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<widget pos="10,10" size="120x100" type="*widget.Entry">
						<rectangle fillColor="focus" pos="0,96" size="120x4"/>
						<widget pos="4,4" size="112x92" type="*widget.Scroll">
							<widget size="112x92" type="*widget.entryContent">
								<widget size="112x92" type="*widget.textProvider">
									<text pos="4,4" size="104x21">Testing</text>
									<text pos="4,25" size="104x21">Testing</text>
									<text pos="4,46" size="104x21">Testing</text>
								</widget>
								<rectangle fillColor="focus" pos="56,46" size="2x21"/>
							</widget>
						</widget>
					</widget>
				</content>
			</canvas>
		`, c)
	})
}

func TestEntry_PasteOverSelection(t *testing.T) {
	e := widget.NewEntry()
	e.SetText("Testing")
	typeKeys(e, fyne.KeyRight, fyne.KeyRight, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight, fyne.KeyRight)

	clipboard := test.NewClipboard()
	clipboard.SetContent("Insert")
	shortcut := &fyne.ShortcutPaste{Clipboard: clipboard}
	e.TypedShortcut(shortcut)

	assert.Equal(t, "Insert", clipboard.Content())
	assert.Equal(t, "TeInsertng", e.Text)
}

func TestEntry_PasteUnicode(t *testing.T) {
	e := widget.NewMultiLineEntry()
	e.SetText("line")
	e.CursorColumn = 4

	clipboard := test.NewClipboard()
	clipboard.SetContent("thing {\n\titem: 'val测试'\n}")
	shortcut := &fyne.ShortcutPaste{Clipboard: clipboard}
	e.TypedShortcut(shortcut)

	assert.Equal(t, "thing {\n\titem: 'val测试'\n}", clipboard.Content())
	assert.Equal(t, "linething {\n\titem: 'val测试'\n}", e.Text)

	assert.Equal(t, 2, e.CursorRow)
	assert.Equal(t, 1, e.CursorColumn)
}

func TestEntry_Placeholder(t *testing.T) {
	entry := &widget.Entry{}
	entry.Text = "Text"
	entry.PlaceHolder = "Placehold"

	window := test.NewWindow(entry)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, "Text", entry.Text)
	test.AssertRendersToMarkup(t, `
			<canvas padded size="46x45">
				<content>
					<widget pos="4,4" size="38x37" type="*widget.Entry">
						<rectangle fillColor="shadow" pos="0,33" size="38x4"/>
						<widget pos="4,4" size="30x29" type="*widget.Scroll">
							<widget size="30x29" type="*widget.entryContent">
								<widget size="30x29" type="*widget.textProvider">
									<text pos="4,4" size="22x21">Text</text>
								</widget>
							</widget>
						</widget>
					</widget>
				</content>
			</canvas>
		`, c)

	entry.SetText("")
	assert.Equal(t, "", entry.Text)
	test.AssertRendersToMarkup(t, `
			<canvas padded size="46x45">
				<content>
					<widget pos="4,4" size="38x37" type="*widget.Entry">
						<rectangle fillColor="shadow" pos="0,33" size="38x4"/>
						<widget pos="4,4" size="30x29" type="*widget.Scroll">
							<widget size="30x29" type="*widget.entryContent">
								<widget size="30x29" type="*widget.textProvider">
									<text color="placeholder" pos="4,4" size="22x21">Placehold</text>
								</widget>
								<widget size="30x29" type="*widget.textProvider">
									<text pos="4,4" size="22x21"></text>
								</widget>
							</widget>
						</widget>
					</widget>
				</content>
			</canvas>
		`, c)
}

func TestEntry_Select(t *testing.T) {
	for name, tt := range map[string]struct {
		keys          []fyne.KeyName
		text          string
		setupReverse  bool
		wantMarkup    string
		wantSelection string
		wantText      string
	}{
		"delete single-line": {
			keys:     []fyne.KeyName{fyne.KeyDelete},
			wantText: "Testing\nTeng\nTesting",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Teng</text>
      							<text pos="4,46" size="104x21">Testing</text>
      						</widget>
      						<rectangle fillColor="focus" pos="20,25" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
		"delete multi-line": {
			keys:     []fyne.KeyName{keyShiftLeftDown, fyne.KeyDown, fyne.KeyDelete},
			wantText: "Testing\nTeng",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Teng</text>
      							<text pos="4,46" size="104x21"></text>
      						</widget>
      						<rectangle fillColor="focus" pos="20,25" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
		"delete reverse multi-line": {
			keys:         []fyne.KeyName{keyShiftLeftDown, fyne.KeyDown, fyne.KeyDelete},
			setupReverse: true,
			wantText:     "Testing\nTestisting",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Testisting</text>
      							<text pos="4,46" size="104x21"></text>
      						</widget>
      						<rectangle fillColor="focus" pos="37,25" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
		"delete select down with Shift held": {
			keys:          []fyne.KeyName{keyShiftLeftDown, fyne.KeyDelete, fyne.KeyDown},
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "ng\nTe",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<rectangle fillColor="focus" pos="20,25" size="20x21"/>
      						<rectangle fillColor="focus" pos="3,46" size="18x21"/>
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Teng</text>
      							<text pos="4,46" size="104x21">Testing</text>
      						</widget>
      						<rectangle fillColor="focus" pos="20,46" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
		"delete reverse select down with Shift held": {
			keys:          []fyne.KeyName{keyShiftLeftDown, fyne.KeyDelete, fyne.KeyDown},
			setupReverse:  true,
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "ng\nTe",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<rectangle fillColor="focus" pos="20,25" size="20x21"/>
      						<rectangle fillColor="focus" pos="3,46" size="18x21"/>
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Teng</text>
      							<text pos="4,46" size="104x21">Testing</text>
      						</widget>
      						<rectangle fillColor="focus" pos="20,46" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
		"delete select up with Shift held": {
			keys:          []fyne.KeyName{keyShiftLeftDown, fyne.KeyDelete, fyne.KeyUp},
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "sting\nTe",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<rectangle fillColor="focus" pos="20,4" size="37x21"/>
      						<rectangle fillColor="focus" pos="3,25" size="18x21"/>
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Teng</text>
      							<text pos="4,46" size="104x21">Testing</text>
      						</widget>
      						<rectangle fillColor="focus" pos="20,4" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
		"delete reverse select up with Shift held": {
			keys:          []fyne.KeyName{keyShiftLeftDown, fyne.KeyDelete, fyne.KeyUp},
			setupReverse:  true,
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "sting\nTe",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<rectangle fillColor="focus" pos="20,4" size="37x21"/>
      						<rectangle fillColor="focus" pos="3,25" size="18x21"/>
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Teng</text>
      							<text pos="4,46" size="104x21">Testing</text>
      						</widget>
      						<rectangle fillColor="focus" pos="20,4" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
		// The backspace delete behaviour is the same as via delete.
		"backspace single-line": {
			keys:     []fyne.KeyName{fyne.KeyBackspace},
			wantText: "Testing\nTeng\nTesting",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Teng</text>
      							<text pos="4,46" size="104x21">Testing</text>
      						</widget>
      						<rectangle fillColor="focus" pos="20,25" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
		"backspace multi-line": {
			keys:     []fyne.KeyName{keyShiftLeftDown, fyne.KeyDown, fyne.KeyBackspace},
			wantText: "Testing\nTeng",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Teng</text>
      							<text pos="4,46" size="104x21"></text>
      						</widget>
      						<rectangle fillColor="focus" pos="20,25" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
		"backspace reverse multi-line": {
			keys:         []fyne.KeyName{keyShiftLeftDown, fyne.KeyDown, fyne.KeyBackspace},
			setupReverse: true,
			wantText:     "Testing\nTestisting",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Testisting</text>
      							<text pos="4,46" size="104x21"></text>
      						</widget>
      						<rectangle fillColor="focus" pos="37,25" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
		"backspace select down with Shift held": {
			keys:          []fyne.KeyName{keyShiftLeftDown, fyne.KeyBackspace, fyne.KeyDown},
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "ng\nTe",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<rectangle fillColor="focus" pos="20,25" size="20x21"/>
      						<rectangle fillColor="focus" pos="3,46" size="18x21"/>
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Teng</text>
      							<text pos="4,46" size="104x21">Testing</text>
      						</widget>
      						<rectangle fillColor="focus" pos="20,46" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
		"backspace reverse select down with Shift held": {
			keys:          []fyne.KeyName{keyShiftLeftDown, fyne.KeyBackspace, fyne.KeyDown},
			setupReverse:  true,
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "ng\nTe",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<rectangle fillColor="focus" pos="20,25" size="20x21"/>
      						<rectangle fillColor="focus" pos="3,46" size="18x21"/>
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Teng</text>
      							<text pos="4,46" size="104x21">Testing</text>
      						</widget>
      						<rectangle fillColor="focus" pos="20,46" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
		"backspace select up with Shift held": {
			keys:          []fyne.KeyName{keyShiftLeftDown, fyne.KeyBackspace, fyne.KeyUp},
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "sting\nTe",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<rectangle fillColor="focus" pos="20,4" size="37x21"/>
      						<rectangle fillColor="focus" pos="3,25" size="18x21"/>
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Teng</text>
      							<text pos="4,46" size="104x21">Testing</text>
      						</widget>
      						<rectangle fillColor="focus" pos="20,4" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
		"backspace reverse select up with Shift held": {
			keys:          []fyne.KeyName{keyShiftLeftDown, fyne.KeyBackspace, fyne.KeyUp},
			setupReverse:  true,
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "sting\nTe",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<rectangle fillColor="focus" pos="20,4" size="37x21"/>
      						<rectangle fillColor="focus" pos="3,25" size="18x21"/>
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Teng</text>
      							<text pos="4,46" size="104x21">Testing</text>
      						</widget>
      						<rectangle fillColor="focus" pos="20,4" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
		// Erase the selection and add a newline at selection start
		"enter": {
			keys:     []fyne.KeyName{fyne.KeyEnter},
			wantText: "Testing\nTe\nng\nTesting",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Te</text>
      							<text pos="4,46" size="104x21">ng</text>
      							<text pos="4,67" size="104x21">Testing</text>
      						</widget>
      						<rectangle fillColor="focus" pos="3,46" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
		"enter reverse": {
			keys:         []fyne.KeyName{fyne.KeyEnter},
			setupReverse: true,
			wantText:     "Testing\nTe\nng\nTesting",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Te</text>
      							<text pos="4,46" size="104x21">ng</text>
      							<text pos="4,67" size="104x21">Testing</text>
      						</widget>
      						<rectangle fillColor="focus" pos="3,46" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
		"replace": {
			text:     "hello",
			wantText: "Testing\nTehellong\nTesting",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Tehellong</text>
      							<text pos="4,46" size="104x21">Testing</text>
      						</widget>
      						<rectangle fillColor="focus" pos="55,25" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
		"replace reverse": {
			text:         "hello",
			setupReverse: true,
			wantText:     "Testing\nTehellong\nTesting",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Tehellong</text>
      							<text pos="4,46" size="104x21">Testing</text>
      						</widget>
      						<rectangle fillColor="focus" pos="55,25" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
		"deselect and delete": {
			keys:     []fyne.KeyName{fyne.KeyLeft, fyne.KeyDelete},
			wantText: "Testing\nTeting\nTesting",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Teting</text>
      							<text pos="4,46" size="104x21">Testing</text>
      						</widget>
      						<rectangle fillColor="focus" pos="20,25" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
		"deselect and delete holding Shift": {
			keys:     []fyne.KeyName{fyne.KeyLeft, keyShiftLeftDown, fyne.KeyDelete},
			wantText: "Testing\nTeting\nTesting",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Teting</text>
      							<text pos="4,46" size="104x21">Testing</text>
      						</widget>
      						<rectangle fillColor="focus" pos="20,25" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
		// ensure that backspace doesn't leave a selection start at the old cursor position
		"deselect and backspace holding Shift": {
			keys:     []fyne.KeyName{fyne.KeyLeft, keyShiftLeftDown, fyne.KeyBackspace},
			wantText: "Testing\nTsting\nTesting",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Tsting</text>
      							<text pos="4,46" size="104x21">Testing</text>
      						</widget>
      						<rectangle fillColor="focus" pos="12,25" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
		// clear selection, select a character and while holding shift issue two backspaces
		"deselect, select and double backspace": {
			keys:     []fyne.KeyName{fyne.KeyRight, fyne.KeyLeft, keyShiftLeftDown, fyne.KeyLeft, fyne.KeyBackspace, fyne.KeyBackspace},
			wantText: "Testing\nTeing\nTesting",
			wantMarkup: `
      	<canvas padded size="150x200">
      		<content>
      			<widget pos="10,10" size="120x100" type="*widget.Entry">
      				<rectangle fillColor="focus" pos="0,96" size="120x4"/>
      				<widget pos="4,4" size="112x92" type="*widget.Scroll">
      					<widget size="112x92" type="*widget.entryContent">
      						<widget size="112x92" type="*widget.textProvider">
      							<text pos="4,4" size="104x21">Testing</text>
      							<text pos="4,25" size="104x21">Teing</text>
      							<text pos="4,46" size="104x21">Testing</text>
      						</widget>
      						<rectangle fillColor="focus" pos="20,25" size="2x21"/>
      					</widget>
      				</widget>
      			</widget>
      		</content>
      	</canvas>
			`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			entry, window := setupSelection(t, tt.setupReverse)
			defer teardownImageTest(window)
			c := window.Canvas()

			if tt.text != "" {
				test.Type(entry, tt.text)
			} else {
				typeKeys(entry, tt.keys...)
			}
			assert.Equal(t, tt.wantText, entry.Text)
			assert.Equal(t, tt.wantSelection, entry.SelectedText())
			test.AssertRendersToMarkup(t, tt.wantMarkup, c)
		})
	}
}

func TestEntry_SelectAll(t *testing.T) {
	e, window := setupImageTest(t, true)
	defer teardownImageTest(window)
	c := window.Canvas()

	c.Focus(e)
	e.SetText("First Row\nSecond Row\nThird Row")
	test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<widget pos="10,10" size="120x100" type="*widget.Entry">
						<rectangle fillColor="focus" pos="0,96" size="120x4"/>
						<widget pos="4,4" size="112x92" type="*widget.Scroll">
							<widget size="112x92" type="*widget.entryContent">
								<widget size="112x92" type="*widget.textProvider">
									<text pos="4,4" size="104x21">First Row</text>
									<text pos="4,25" size="104x21">Second Row</text>
									<text pos="4,46" size="104x21">Third Row</text>
								</widget>
								<rectangle fillColor="focus" pos="3,4" size="2x21"/>
							</widget>
						</widget>
					</widget>
				</content>
			</canvas>
		`, c)

	shortcut := &fyne.ShortcutSelectAll{}
	e.TypedShortcut(shortcut)
	assert.Equal(t, 2, e.CursorRow)
	assert.Equal(t, 9, e.CursorColumn)
	test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<widget pos="10,10" size="120x100" type="*widget.Entry">
						<rectangle fillColor="focus" pos="0,96" size="120x4"/>
						<widget pos="4,4" size="112x92" type="*widget.Scroll">
							<widget size="112x92" type="*widget.entryContent">
								<rectangle fillColor="focus" pos="3,4" size="67x21"/>
								<rectangle fillColor="focus" pos="3,25" size="88x21"/>
								<rectangle fillColor="focus" pos="3,46" size="73x21"/>
								<widget size="112x92" type="*widget.textProvider">
									<text pos="4,4" size="104x21">First Row</text>
									<text pos="4,25" size="104x21">Second Row</text>
									<text pos="4,46" size="104x21">Third Row</text>
								</widget>
								<rectangle fillColor="focus" pos="75,46" size="2x21"/>
							</widget>
						</widget>
					</widget>
				</content>
			</canvas>
		`, c)
}

func TestEntry_SelectAll_EmptyEntry(t *testing.T) {
	entry := widget.NewEntry()
	entry.TypedShortcut(&fyne.ShortcutSelectAll{})

	assert.Equal(t, "", entry.SelectedText())
}

func TestEntry_SelectEndWithoutShift(t *testing.T) {
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	// end after releasing shift
	typeKeys(e, keyShiftLeftUp, fyne.KeyEnd)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Testing</text>
								<text pos="4,25" size="104x21">Testing</text>
								<text pos="4,46" size="104x21">Testing</text>
							</widget>
							<rectangle fillColor="focus" pos="56,25" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Equal(t, "", e.SelectedText())
}

func TestEntry_SelectHomeEnd(t *testing.T) {
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	// Hold shift to continue selection
	typeKeys(e, keyShiftLeftDown)

	// T e[s t i]n g -> end -> // T e[s t i n g]
	typeKeys(e, fyne.KeyEnd)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<rectangle fillColor="focus" pos="20,25" size="37x21"/>
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Testing</text>
								<text pos="4,25" size="104x21">Testing</text>
								<text pos="4,46" size="104x21">Testing</text>
							</widget>
							<rectangle fillColor="focus" pos="56,25" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Equal(t, "sting", e.SelectedText())

	// T e[s t i n g] -> home -> [T e]s t i n g
	typeKeys(e, fyne.KeyHome)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<rectangle fillColor="focus" pos="3,25" size="18x21"/>
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Testing</text>
								<text pos="4,25" size="104x21">Testing</text>
								<text pos="4,46" size="104x21">Testing</text>
							</widget>
							<rectangle fillColor="focus" pos="3,25" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Equal(t, "Te", e.SelectedText())
}

func TestEntry_SelectHomeWithoutShift(t *testing.T) {
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	// home after releasing shift
	typeKeys(e, keyShiftLeftUp, fyne.KeyHome)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Testing</text>
								<text pos="4,25" size="104x21">Testing</text>
								<text pos="4,46" size="104x21">Testing</text>
							</widget>
							<rectangle fillColor="focus" pos="3,25" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Equal(t, "", e.SelectedText())
}

func TestEntry_SelectSnapDown(t *testing.T) {
	// down snaps to end, but it also moves
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)

	typeKeys(e, keyShiftLeftUp, fyne.KeyDown)
	assert.Equal(t, 2, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Testing</text>
								<text pos="4,25" size="104x21">Testing</text>
								<text pos="4,46" size="104x21">Testing</text>
							</widget>
							<rectangle fillColor="focus" pos="37,46" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Equal(t, "", e.SelectedText())
}

func TestEntry_SelectSnapLeft(t *testing.T) {
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)

	typeKeys(e, keyShiftLeftUp, fyne.KeyLeft)
	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 2, e.CursorColumn)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Testing</text>
								<text pos="4,25" size="104x21">Testing</text>
								<text pos="4,46" size="104x21">Testing</text>
							</widget>
							<rectangle fillColor="focus" pos="20,25" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Equal(t, "", e.SelectedText())
}

func TestEntry_SelectSnapRight(t *testing.T) {
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)

	typeKeys(e, keyShiftLeftUp, fyne.KeyRight)
	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Testing</text>
								<text pos="4,25" size="104x21">Testing</text>
								<text pos="4,46" size="104x21">Testing</text>
							</widget>
							<rectangle fillColor="focus" pos="37,25" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Equal(t, "", e.SelectedText())
}

func TestEntry_SelectSnapUp(t *testing.T) {
	// up snaps to start, but it also moves
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)

	typeKeys(e, keyShiftLeftUp, fyne.KeyUp)
	assert.Equal(t, 0, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Testing</text>
								<text pos="4,25" size="104x21">Testing</text>
								<text pos="4,46" size="104x21">Testing</text>
							</widget>
							<rectangle fillColor="focus" pos="37,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Equal(t, "", e.SelectedText())
}

func TestEntry_SelectedText(t *testing.T) {
	e, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	c.Focus(e)
	e.SetText("Testing")
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Testing</text>
							</widget>
							<rectangle fillColor="focus" pos="3,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Equal(t, "", e.SelectedText())

	// move right, press & hold shift and move right
	typeKeys(e, fyne.KeyRight, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight)
	assert.Equal(t, "es", e.SelectedText())
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<rectangle fillColor="focus" pos="12,4" size="17x21"/>
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Testing</text>
							</widget>
							<rectangle fillColor="focus" pos="28,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)

	// release shift
	typeKeys(e, keyShiftLeftUp)
	// press shift and move
	typeKeys(e, keyShiftLeftDown, fyne.KeyRight)
	assert.Equal(t, "est", e.SelectedText())
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<rectangle fillColor="focus" pos="12,4" size="22x21"/>
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Testing</text>
							</widget>
							<rectangle fillColor="focus" pos="33,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)

	// release shift and move right
	typeKeys(e, keyShiftLeftUp, fyne.KeyRight)
	assert.Equal(t, "", e.SelectedText())
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Testing</text>
							</widget>
							<rectangle fillColor="focus" pos="33,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)

	// press shift and move left
	typeKeys(e, keyShiftLeftDown, fyne.KeyLeft, fyne.KeyLeft)
	assert.Equal(t, "st", e.SelectedText())
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<rectangle fillColor="focus" pos="20,4" size="14x21"/>
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Testing</text>
							</widget>
							<rectangle fillColor="focus" pos="20,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
}

func TestEntry_SelectionHides(t *testing.T) {
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	c.Unfocus()
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="shadow" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Testing</text>
								<text pos="4,25" size="104x21">Testing</text>
								<text pos="4,46" size="104x21">Testing</text>
							</widget>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Equal(t, "sti", e.SelectedText())

	c.Focus(e)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<rectangle fillColor="focus" pos="20,25" size="18x21"/>
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Testing</text>
								<text pos="4,25" size="104x21">Testing</text>
								<text pos="4,46" size="104x21">Testing</text>
							</widget>
							<rectangle fillColor="focus" pos="37,25" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Equal(t, "sti", e.SelectedText())
}

func TestEntry_SetPlaceHolder(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, 0, len(entry.Text))

	entry.SetPlaceHolder("Test")
	assert.Equal(t, 0, len(entry.Text))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="shadow" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text color="placeholder" pos="4,4" size="104x21">Test</text>
							</widget>
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21"></text>
							</widget>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)

	entry.SetText("Hi")
	assert.Equal(t, 2, len(entry.Text))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="shadow" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Hi</text>
							</widget>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
}

func TestEntry_Disable_KeyDown(t *testing.T) {
	entry := widget.NewEntry()

	test.Type(entry, "H")
	entry.Disable()
	test.Type(entry, "i")
	assert.Equal(t, "H", entry.Text)

	entry.Enable()
	test.Type(entry, "i")
	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_Disable_OnFocus(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.Disable()
	entry.FocusGained()
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="disabled" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text color="disabled" pos="4,4" size="104x21"></text>
							</widget>
							<widget size="112x29" type="*widget.textProvider">
								<text color="disabled" pos="4,4" size="104x21"></text>
							</widget>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)

	entry.Enable()
	entry.FocusGained()
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text color="placeholder" pos="4,4" size="104x21"></text>
							</widget>
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21"></text>
							</widget>
							<rectangle fillColor="focus" pos="3,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
}

func TestEntry_SetText_EmptyString(t *testing.T) {
	entry := widget.NewEntry()

	assert.Equal(t, 0, entry.CursorColumn)

	test.Type(entry, "test")
	assert.Equal(t, 4, entry.CursorColumn)
	entry.SetText("")
	assert.Equal(t, 0, entry.CursorColumn)

	entry = widget.NewMultiLineEntry()
	test.Type(entry, "test\ntest")

	down := &fyne.KeyEvent{Name: fyne.KeyDown}
	entry.TypedKey(down)

	assert.Equal(t, 4, entry.CursorColumn)
	assert.Equal(t, 1, entry.CursorRow)
	entry.SetText("")
	assert.Equal(t, 0, entry.CursorColumn)
	assert.Equal(t, 0, entry.CursorRow)
}

func TestEntry_SetText_Manual(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.Text = "Test"
	entry.Refresh()
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="shadow" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Test</text>
							</widget>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
}

func TestEntry_SetText_Overflow(t *testing.T) {
	entry := widget.NewEntry()

	assert.Equal(t, 0, entry.CursorColumn)

	test.Type(entry, "test")
	assert.Equal(t, 4, entry.CursorColumn)

	entry.SetText("x")
	assert.Equal(t, 1, entry.CursorColumn)

	key := &fyne.KeyEvent{Name: fyne.KeyDelete}
	entry.TypedKey(key)

	assert.Equal(t, 1, entry.CursorColumn)
	assert.Equal(t, "x", entry.Text)

	key = &fyne.KeyEvent{Name: fyne.KeyBackspace}
	entry.TypedKey(key)

	assert.Equal(t, 0, entry.CursorColumn)
	assert.Equal(t, "", entry.Text)
}

func TestEntry_SetText_Underflow(t *testing.T) {
	entry := widget.NewEntry()
	test.Type(entry, "test")
	assert.Equal(t, 4, entry.CursorColumn)

	entry.Text = ""
	entry.Refresh()
	assert.Equal(t, 0, entry.CursorColumn)

	key := &fyne.KeyEvent{Name: fyne.KeyBackspace}
	entry.TypedKey(key)

	assert.Equal(t, 0, entry.CursorColumn)
	assert.Equal(t, "", entry.Text)
}

func TestEntry_SetTextStyle(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.Text = "Styled Text"
	entry.TextStyle = fyne.TextStyle{Bold: true}
	entry.Refresh()
	test.AssertRendersToMarkup(t, `
	<canvas padded size="150x200">
		<content>
			<widget pos="10,10" size="120x37" type="*widget.Entry">
				<rectangle fillColor="shadow" pos="0,33" size="120x4"/>
				<widget pos="4,4" size="112x29" type="*widget.Scroll">
					<widget size="112x29" type="*widget.entryContent">
						<widget size="112x29" type="*widget.textProvider">
							<text bold pos="4,4" size="104x21">Styled Text</text>
						</widget>
					</widget>
				</widget>
			</widget>
		</content>
	</canvas>
	`, c)

	entry.TextStyle = fyne.TextStyle{Monospace: true}
	entry.Refresh()
	test.AssertRendersToMarkup(t, `
	<canvas padded size="150x200">
		<content>
			<widget pos="10,10" size="120x37" type="*widget.Entry">
				<rectangle fillColor="shadow" pos="0,33" size="120x4"/>
				<widget pos="4,4" size="112x29" type="*widget.Scroll">
					<widget size="112x29" type="*widget.entryContent">
						<widget size="112x29" type="*widget.textProvider">
							<text monospace pos="4,4" size="104x18">Styled Text</text>
						</widget>
					</widget>
				</widget>
			</widget>
		</content>
	</canvas>
	`, c)

	entry.TextStyle = fyne.TextStyle{Italic: true}
	entry.Refresh()
	test.AssertRendersToMarkup(t, `
	<canvas padded size="150x200">
		<content>
			<widget pos="10,10" size="120x37" type="*widget.Entry">
				<rectangle fillColor="shadow" pos="0,33" size="120x4"/>
				<widget pos="4,4" size="112x29" type="*widget.Scroll">
					<widget size="112x29" type="*widget.entryContent">
						<widget size="112x29" type="*widget.textProvider">
							<text italic pos="4,4" size="104x21">Styled Text</text>
						</widget>
					</widget>
				</widget>
			</widget>
		</content>
	</canvas>
	`, c)
}

func TestEntry_Tapped(t *testing.T) {
	entry, window := setupImageTest(t, true)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.SetText("MMM\nWWW\n")
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="shadow" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">MMM</text>
								<text pos="4,25" size="104x21">WWW</text>
								<text pos="4,46" size="104x21"></text>
							</widget>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)

	entry.FocusGained()
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">MMM</text>
								<text pos="4,25" size="104x21">WWW</text>
								<text pos="4,46" size="104x21"></text>
							</widget>
							<rectangle fillColor="focus" pos="3,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)

	testCharSize := theme.TextSize()
	pos := fyne.NewPos(entryOffset+theme.Padding()+testCharSize*1.5, entryOffset+theme.Padding()+testCharSize/2) // tap in the middle of the 2nd "M"
	test.TapCanvas(window.Canvas(), pos)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">MMM</text>
								<text pos="4,25" size="104x21">WWW</text>
								<text pos="4,46" size="104x21"></text>
							</widget>
							<rectangle fillColor="focus" pos="17,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	pos = fyne.NewPos(entryOffset+theme.Padding()+theme.Padding()+testCharSize*2.5, entryOffset+theme.Padding()+testCharSize/2) // tap in the middle of the 3rd "M"
	test.TapCanvas(window.Canvas(), pos)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">MMM</text>
								<text pos="4,25" size="104x21">WWW</text>
								<text pos="4,46" size="104x21"></text>
							</widget>
							<rectangle fillColor="focus" pos="31,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 2, entry.CursorColumn)

	pos = fyne.NewPos(entryOffset+theme.Padding()+testCharSize*4, entryOffset+theme.Padding()+testCharSize/2) // tap after text
	test.TapCanvas(window.Canvas(), pos)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">MMM</text>
								<text pos="4,25" size="104x21">WWW</text>
								<text pos="4,46" size="104x21"></text>
							</widget>
							<rectangle fillColor="focus" pos="45,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 3, entry.CursorColumn)

	pos = fyne.NewPos(entryOffset+testCharSize, entryOffset+testCharSize*4) // tap below rows
	test.TapCanvas(window.Canvas(), pos)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">MMM</text>
								<text pos="4,25" size="104x21">WWW</text>
								<text pos="4,46" size="104x21"></text>
							</widget>
							<rectangle fillColor="focus" pos="3,46" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Equal(t, 2, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)
}

func TestEntry_TappedSecondary(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	tapPos := fyne.NewPos(20, 10)
	test.TapSecondaryAt(entry, tapPos)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text color="placeholder" pos="4,4" size="104x21"></text>
							</widget>
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21"></text>
							</widget>
							<rectangle fillColor="focus" pos="3,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
			<overlay>
				<widget size="150x200" type="*widget.OverlayContainer">
					<widget pos="30,20" size="79x136" type="*widget.Menu">
						<widget size="79x136" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="79x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="79,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="79,0" size="4x136" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="79,136" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,136" size="79x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,136" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x136"/>
						</widget>
						<widget size="79x136" type="*widget.Scroll">
							<widget size="79x136" type="*widget.menuBox">
								<container pos="0,4" size="79x144">
									<widget size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="25x21">Cut</text>
									</widget>
									<widget pos="0,33" size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="36x21">Copy</text>
									</widget>
									<widget pos="0,66" size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="39x21">Paste</text>
									</widget>
									<widget pos="0,99" size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Select all</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, c)
	assert.Equal(t, 1, len(c.Overlays().List()))
	c.Overlays().Remove(c.Overlays().Top())

	entry.Disable()
	test.TapSecondaryAt(entry, tapPos)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="disabled" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text color="disabled" pos="4,4" size="104x21"></text>
							</widget>
							<widget size="112x29" type="*widget.textProvider">
								<text color="disabled" pos="4,4" size="104x21"></text>
							</widget>
						</widget>
					</widget>
				</widget>
			</content>
			<overlay>
				<widget size="150x200" type="*widget.OverlayContainer">
					<widget pos="30,20" size="79x70" type="*widget.Menu">
						<widget size="79x70" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="79x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="79,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="79,0" size="4x70" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="79,70" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,70" size="79x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
						</widget>
						<widget size="79x70" type="*widget.Scroll">
							<widget size="79x70" type="*widget.menuBox">
								<container pos="0,4" size="79x78">
									<widget size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="36x21">Copy</text>
									</widget>
									<widget pos="0,33" size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Select all</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, c)
	assert.Equal(t, 1, len(c.Overlays().List()))
	c.Overlays().Remove(c.Overlays().Top())

	entry.Password = true
	entry.Refresh()
	test.TapSecondaryAt(entry, tapPos)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="disabled" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text color="disabled" pos="4,4" size="104x21"></text>
							</widget>
							<widget size="112x29" type="*widget.textProvider">
								<text color="disabled" pos="4,4" size="104x21"></text>
							</widget>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Nil(t, c.Overlays().Top(), "No popup for disabled password")

	entry.Enable()
	test.TapSecondaryAt(entry, tapPos)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text color="placeholder" pos="4,4" size="104x21"></text>
							</widget>
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21"></text>
							</widget>
							<rectangle fillColor="focus" pos="3,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
			<overlay>
				<widget size="150x200" type="*widget.OverlayContainer">
					<widget pos="30,20" size="79x70" type="*widget.Menu">
						<widget size="79x70" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="79x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="79,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="79,0" size="4x70" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="79,70" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,70" size="79x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
						</widget>
						<widget size="79x70" type="*widget.Scroll">
							<widget size="79x70" type="*widget.menuBox">
								<container pos="0,4" size="79x78">
									<widget size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="39x21">Paste</text>
									</widget>
									<widget pos="0,33" size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Select all</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, c)
	assert.Equal(t, 1, len(c.Overlays().List()))
}

func TestEntry_TextWrap(t *testing.T) {
	singleLineWrapOffMarkup := `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.Scroll">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21">Testing Wrapping</text>
							</widget>
							<rectangle fillColor="focus" pos="3,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`
	multiLineWrapOffMarkup := `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x100" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,96" size="120x4"/>
					<widget pos="4,4" size="112x92" type="*widget.Scroll">
						<widget size="112x92" type="*widget.entryContent">
							<widget size="112x92" type="*widget.textProvider">
								<text pos="4,4" size="104x21">A long text on short words w/o NLs or LFs.</text>
							</widget>
							<rectangle fillColor="focus" pos="3,4" size="2x21"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`
	for name, tt := range map[string]struct {
		multiLine bool
		want      string
		wrap      fyne.TextWrap
	}{
		"single line WrapOff": {
			want: singleLineWrapOffMarkup,
		},
		// Disallowed - fallback to TextWrapOff
		"single line Truncate": {
			wrap: fyne.TextTruncate,
			want: singleLineWrapOffMarkup,
		},
		// Disallowed - fallback to TextWrapOff
		"single line WrapBreak": {
			wrap: fyne.TextWrapBreak,
			want: singleLineWrapOffMarkup,
		},
		// Disallowed - fallback to TextWrapOff
		"single line WrapWord": {
			wrap: fyne.TextWrapWord,
			want: singleLineWrapOffMarkup,
		},
		"multi line WrapOff": {
			multiLine: true,
			want:      multiLineWrapOffMarkup,
		},
		// Disallowed - fallback to TextWrapOff
		"multi line Truncate": {
			multiLine: true,
			wrap:      fyne.TextTruncate,
			want:      multiLineWrapOffMarkup,
		},
		"multi line WrapBreak": {
			multiLine: true,
			wrap:      fyne.TextWrapBreak,
			want: `
				<canvas padded size="150x200">
					<content>
						<widget pos="10,10" size="120x100" type="*widget.Entry">
							<rectangle fillColor="focus" pos="0,96" size="120x4"/>
							<widget pos="4,4" size="112x92" type="*widget.Scroll">
								<widget size="112x92" type="*widget.entryContent">
									<widget size="112x92" type="*widget.textProvider">
										<text pos="4,4" size="104x21">A long text on </text>
										<text pos="4,25" size="104x21">short words w</text>
										<text pos="4,46" size="104x21">/o NLs or LFs.</text>
									</widget>
									<rectangle fillColor="focus" pos="3,4" size="2x21"/>
								</widget>
							</widget>
						</widget>
					</content>
				</canvas>
			`,
		},
		"multi line WrapWord": {
			multiLine: true,
			wrap:      fyne.TextWrapWord,
			want: `
				<canvas padded size="150x200">
					<content>
						<widget pos="10,10" size="120x100" type="*widget.Entry">
							<rectangle fillColor="focus" pos="0,96" size="120x4"/>
							<widget pos="4,4" size="112x92" type="*widget.Scroll">
								<widget size="112x92" type="*widget.entryContent">
									<widget size="112x92" type="*widget.textProvider">
										<text pos="4,4" size="104x21">A long text on</text>
										<text pos="4,25" size="104x21">short words</text>
										<text pos="4,46" size="104x21">w/o NLs or</text>
										<text pos="4,67" size="104x21">LFs.</text>
									</widget>
									<rectangle fillColor="focus" pos="3,4" size="2x21"/>
								</widget>
							</widget>
						</widget>
					</content>
				</canvas>
			`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			e, window := setupImageTest(t, tt.multiLine)
			defer teardownImageTest(window)
			c := window.Canvas()

			c.Focus(e)
			e.Wrapping = tt.wrap
			if tt.multiLine {
				e.SetText("A long text on short words w/o NLs or LFs.")
			} else {
				e.SetText("Testing Wrapping")
			}
			test.AssertRendersToMarkup(t, tt.want, c)
		})
	}
}

func TestMultiLineEntry_MinSize(t *testing.T) {
	entry := widget.NewEntry()
	singleMin := entry.MinSize()

	multi := widget.NewMultiLineEntry()
	multiMin := multi.MinSize()

	assert.Equal(t, singleMin.Width, multiMin.Width)
	assert.True(t, multiMin.Height > singleMin.Height)

	multi.MultiLine = false
	multiMin = multi.MinSize()
	assert.Equal(t, singleMin.Height, multiMin.Height)
}

func TestNewEntryWithData(t *testing.T) {
	str := binding.NewString()
	err := str.Set("Init")
	assert.Nil(t, err)

	entry := widget.NewEntryWithData(str)
	waitForBinding()
	assert.Equal(t, "Init", entry.Text)

	entry.SetText("Typed")
	v, err := str.Get()
	assert.Nil(t, err)
	assert.Equal(t, "Typed", v)
}

func TestPasswordEntry_ActionItemSizeAndPlacement(t *testing.T) {
	e := widget.NewEntry()
	b := widget.NewButton("", func() {})
	b.Icon = theme.CancelIcon()
	e.ActionItem = b
	test.WidgetRenderer(e).Layout(e.MinSize())
	assert.Equal(t, fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()), b.Size())
	assert.Equal(t, fyne.NewPos(e.MinSize().Width-2*theme.Padding()-b.Size().Width, 2*theme.Padding()), b.Position())
}

func TestPasswordEntry_NewlineIgnored(t *testing.T) {
	entry := widget.NewPasswordEntry()
	entry.SetText("test")

	checkNewlineIgnored(t, entry)
}

func TestPasswordEntry_Obfuscation(t *testing.T) {
	entry, window := setupPasswordTest(t)
	defer teardownImageTest(window)
	c := window.Canvas()

	test.Type(entry, "Hié™שרה")
	assert.Equal(t, "Hié™שרה", entry.Text)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x100">
			<content>
				<widget pos="10,10" size="130x37" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,33" size="130x4"/>
					<widget pos="4,4" size="102x29" type="*widget.Scroll">
						<widget size="102x29" type="*widget.entryContent">
							<widget size="102x29" type="*widget.textProvider">
								<text pos="4,4" size="94x21">•••••••</text>
							</widget>
							<rectangle fillColor="focus" pos="43,4" size="2x21"/>
						</widget>
					</widget>
					<widget pos="102,8" size="20x20" type="*widget.passwordRevealer">
						<image rsc="visibilityOffIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
}

func TestPasswordEntry_Placeholder(t *testing.T) {
	entry, window := setupPasswordTest(t)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.SetPlaceHolder("Password")
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x100">
			<content>
				<widget pos="10,10" size="130x37" type="*widget.Entry">
					<rectangle fillColor="shadow" pos="0,33" size="130x4"/>
					<widget pos="4,4" size="102x29" type="*widget.Scroll">
						<widget size="102x29" type="*widget.entryContent">
							<widget size="102x29" type="*widget.textProvider">
								<text color="placeholder" pos="4,4" size="94x21">Password</text>
							</widget>
							<widget size="102x29" type="*widget.textProvider">
								<text pos="4,4" size="94x21"></text>
							</widget>
						</widget>
					</widget>
					<widget pos="102,8" size="20x20" type="*widget.passwordRevealer">
						<image rsc="visibilityOffIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)

	test.Type(entry, "Hié™שרה")
	assert.Equal(t, "Hié™שרה", entry.Text)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x100">
			<content>
				<widget pos="10,10" size="130x37" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,33" size="130x4"/>
					<widget pos="4,4" size="102x29" type="*widget.Scroll">
						<widget size="102x29" type="*widget.entryContent">
							<widget size="102x29" type="*widget.textProvider">
								<text pos="4,4" size="94x21">•••••••</text>
							</widget>
							<rectangle fillColor="focus" pos="43,4" size="2x21"/>
						</widget>
					</widget>
					<widget pos="102,8" size="20x20" type="*widget.passwordRevealer">
						<image rsc="visibilityOffIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
}

func TestPasswordEntry_Reveal(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	initial := `
		<canvas padded size="150x100">
			<content>
				<widget pos="10,10" size="130x37" type="*widget.Entry">
					<rectangle fillColor="shadow" pos="0,33" size="130x4"/>
					<widget pos="4,4" size="102x29" type="*widget.Scroll">
						<widget size="102x29" type="*widget.entryContent">
							<widget size="102x29" type="*widget.textProvider">
								<text color="placeholder" pos="4,4" size="94x21"></text>
							</widget>
							<widget size="102x29" type="*widget.textProvider">
								<text pos="4,4" size="94x21"></text>
							</widget>
						</widget>
					</widget>
					<widget pos="102,8" size="20x20" type="*widget.passwordRevealer">
						<image rsc="visibilityOffIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
		</canvas>
	`
	concealed := `
		<canvas padded size="150x100">
			<content>
				<widget pos="10,10" size="130x37" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,33" size="130x4"/>
					<widget pos="4,4" size="102x29" type="*widget.Scroll">
						<widget size="102x29" type="*widget.entryContent">
							<widget size="102x29" type="*widget.textProvider">
								<text pos="4,4" size="94x21">•••••••</text>
							</widget>
							<rectangle fillColor="focus" pos="43,4" size="2x21"/>
						</widget>
					</widget>
					<widget pos="102,8" size="20x20" type="*widget.passwordRevealer">
						<image rsc="visibilityOffIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
		</canvas>
	`
	revealed := `
		<canvas padded size="150x100">
			<content>
				<widget pos="10,10" size="130x37" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,33" size="130x4"/>
					<widget pos="4,4" size="102x29" type="*widget.Scroll">
						<widget size="102x29" type="*widget.entryContent">
							<widget size="102x29" type="*widget.textProvider">
								<text pos="4,4" size="94x21">Hié™שרה</text>
							</widget>
							<rectangle fillColor="focus" pos="66,4" size="2x21"/>
						</widget>
					</widget>
					<widget pos="102,8" size="20x20" type="*widget.passwordRevealer">
						<image rsc="visibilityIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
		</canvas>
	`
	t.Run("NewPasswordEntry constructor", func(t *testing.T) {
		entry := widget.NewPasswordEntry()
		window := test.NewWindow(entry)
		defer window.Close()
		window.Resize(fyne.NewSize(150, 100))
		entry.Resize(entry.MinSize().Max(fyne.NewSize(130, 0)))
		entry.Move(fyne.NewPos(10, 10))
		c := window.Canvas()

		test.AssertRendersToMarkup(t, initial, c)

		c.Focus(entry)
		test.Type(entry, "Hié™שרה")
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertRendersToMarkup(t, concealed, c)

		// update the Password field
		entry.Password = false
		entry.Refresh()
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertRendersToMarkup(t, revealed, c)
		assert.Equal(t, entry, c.Focused())

		// update the Password field
		entry.Password = true
		entry.Refresh()
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertRendersToMarkup(t, concealed, c)
		assert.Equal(t, entry, c.Focused())

		// tap on action icon
		tapPos := fyne.NewPos(140-theme.Padding()*2-theme.IconInlineSize()/2, 10+entry.Size().Height/2)
		test.TapCanvas(c, tapPos)
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertRendersToMarkup(t, revealed, c)
		assert.Equal(t, entry, c.Focused())

		// tap on action icon
		test.TapCanvas(c, tapPos)
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertRendersToMarkup(t, concealed, c)
		assert.Equal(t, entry, c.Focused())
	})

	// This test cover backward compatibility use case when on an Entry widget
	// the Password field is set to true.
	// In this case the action item will be set when the renderer is created.
	t.Run("Entry with Password field", func(t *testing.T) {
		entry := &widget.Entry{}
		entry.Password = true
		entry.Refresh()
		window := test.NewWindow(entry)
		defer window.Close()
		window.Resize(fyne.NewSize(150, 100))
		entry.Resize(entry.MinSize().Max(fyne.NewSize(130, 0)))
		entry.Move(fyne.NewPos(10, 10))
		c := window.Canvas()

		test.AssertRendersToMarkup(t, initial, c)

		c.Focus(entry)
		test.Type(entry, "Hié™שרה")
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertRendersToMarkup(t, concealed, c)

		// update the Password field
		entry.Password = false
		entry.Refresh()
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertRendersToMarkup(t, revealed, c)
		assert.Equal(t, entry, c.Focused())
	})
}

func TestSingleLineEntry_NewlineIgnored(t *testing.T) {
	entry := &widget.Entry{MultiLine: false}
	entry.SetText("test")

	checkNewlineIgnored(t, entry)
}

const (
	entryOffset = 10

	keyShiftLeftDown  fyne.KeyName = "LeftShiftDown"
	keyShiftLeftUp    fyne.KeyName = "LeftShiftUp"
	keyShiftRightDown fyne.KeyName = "RightShiftDown"
	keyShiftRightUp   fyne.KeyName = "RightShiftUp"
)

var typeKeys = func(e *widget.Entry, keys ...fyne.KeyName) {
	var keyDown = func(key *fyne.KeyEvent) {
		e.KeyDown(key)
		e.TypedKey(key)
	}

	for _, key := range keys {
		switch key {
		case keyShiftLeftDown:
			keyDown(&fyne.KeyEvent{Name: desktop.KeyShiftLeft})
		case keyShiftLeftUp:
			e.KeyUp(&fyne.KeyEvent{Name: desktop.KeyShiftLeft})
		case keyShiftRightDown:
			keyDown(&fyne.KeyEvent{Name: desktop.KeyShiftRight})
		case keyShiftRightUp:
			e.KeyUp(&fyne.KeyEvent{Name: desktop.KeyShiftRight})
		default:
			keyDown(&fyne.KeyEvent{Name: key})
			e.KeyUp(&fyne.KeyEvent{Name: key})
		}
	}
}

func checkNewlineIgnored(t *testing.T, entry *widget.Entry) {
	assert.Equal(t, 0, entry.CursorRow)

	// only 1 line, do nothing
	down := &fyne.KeyEvent{Name: fyne.KeyDown}
	entry.TypedKey(down)
	assert.Equal(t, 0, entry.CursorRow)

	// return is ignored, do nothing
	ret := &fyne.KeyEvent{Name: fyne.KeyReturn}
	entry.TypedKey(ret)
	assert.Equal(t, 0, entry.CursorRow)

	up := &fyne.KeyEvent{Name: fyne.KeyUp}
	entry.TypedKey(up)
	assert.Equal(t, 0, entry.CursorRow)

	// don't go beyond top
	entry.TypedKey(up)
	assert.Equal(t, 0, entry.CursorRow)
}

func setupImageTest(t *testing.T, multiLine bool) (*widget.Entry, fyne.Window) {
	test.NewApp()

	entry := &widget.Entry{MultiLine: multiLine}
	w := test.NewWindow(entry)
	w.Resize(fyne.NewSize(150, 200))

	if multiLine {
		entry.Resize(fyne.NewSize(120, 100))
	} else {
		entry.Resize(entry.MinSize().Max(fyne.NewSize(120, 0)))
	}
	entry.Move(fyne.NewPos(10, 10))

	if multiLine {
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<widget pos="10,10" size="120x100" type="*widget.Entry">
						<rectangle fillColor="shadow" pos="0,96" size="120x4"/>
						<widget pos="4,4" size="112x92" type="*widget.Scroll">
							<widget size="112x92" type="*widget.entryContent">
								<widget size="112x92" type="*widget.textProvider">
									<text color="placeholder" pos="4,4" size="104x21"></text>
								</widget>
								<widget size="112x92" type="*widget.textProvider">
									<text pos="4,4" size="104x21"></text>
								</widget>
							</widget>
						</widget>
					</widget>
				</content>
			</canvas>
		`, w.Canvas())
	} else {
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<widget pos="10,10" size="120x37" type="*widget.Entry">
						<rectangle fillColor="shadow" pos="0,33" size="120x4"/>
						<widget pos="4,4" size="112x29" type="*widget.Scroll">
							<widget size="112x29" type="*widget.entryContent">
								<widget size="112x29" type="*widget.textProvider">
									<text color="placeholder" pos="4,4" size="104x21"></text>
								</widget>
								<widget size="112x29" type="*widget.textProvider">
									<text pos="4,4" size="104x21"></text>
								</widget>
							</widget>
						</widget>
					</widget>
				</content>
			</canvas>
		`, w.Canvas())
	}

	return entry, w
}

func setupPasswordTest(t *testing.T) (*widget.Entry, fyne.Window) {
	test.NewApp()

	entry := widget.NewPasswordEntry()
	w := test.NewWindow(entry)
	w.Resize(fyne.NewSize(150, 100))

	entry.Resize(entry.MinSize().Max(fyne.NewSize(130, 0)))
	entry.Move(fyne.NewPos(entryOffset, entryOffset))

	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x100">
			<content>
				<widget pos="10,10" size="130x37" type="*widget.Entry">
					<rectangle fillColor="shadow" pos="0,33" size="130x4"/>
					<widget pos="4,4" size="102x29" type="*widget.Scroll">
						<widget size="102x29" type="*widget.entryContent">
							<widget size="102x29" type="*widget.textProvider">
								<text color="placeholder" pos="4,4" size="94x21"></text>
							</widget>
							<widget size="102x29" type="*widget.textProvider">
								<text pos="4,4" size="94x21"></text>
							</widget>
						</widget>
					</widget>
					<widget pos="102,8" size="20x20" type="*widget.passwordRevealer">
						<image rsc="visibilityOffIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
		</canvas>
	`, w.Canvas())

	return entry, w
}

// Selects "sti" on line 2 of a new multiline
// T e s t i n g
// T e[s t i]n g
// T e s t i n g
func setupSelection(t *testing.T, reverse bool) (*widget.Entry, fyne.Window) {
	e, window := setupImageTest(t, true)
	e.SetText("Testing\nTesting\nTesting")
	c := window.Canvas()
	c.Focus(e)
	if reverse {
		e.CursorRow = 1
		e.CursorColumn = 5
		typeKeys(e, keyShiftLeftDown, fyne.KeyLeft, fyne.KeyLeft, fyne.KeyLeft, keyShiftLeftUp)
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<widget pos="10,10" size="120x100" type="*widget.Entry">
						<rectangle fillColor="focus" pos="0,96" size="120x4"/>
						<widget pos="4,4" size="112x92" type="*widget.Scroll">
							<widget size="112x92" type="*widget.entryContent">
								<rectangle fillColor="focus" pos="20,25" size="18x21"/>
								<widget size="112x92" type="*widget.textProvider">
									<text pos="4,4" size="104x21">Testing</text>
									<text pos="4,25" size="104x21">Testing</text>
									<text pos="4,46" size="104x21">Testing</text>
								</widget>
								<rectangle fillColor="focus" pos="20,25" size="2x21"/>
							</widget>
						</widget>
					</widget>
				</content>
			</canvas>
		`, c)
		assert.Equal(t, "sti", e.SelectedText())
	} else {
		e.CursorRow = 1
		e.CursorColumn = 2
		typeKeys(e, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight, fyne.KeyRight, keyShiftLeftUp)
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<widget pos="10,10" size="120x100" type="*widget.Entry">
						<rectangle fillColor="focus" pos="0,96" size="120x4"/>
						<widget pos="4,4" size="112x92" type="*widget.Scroll">
							<widget size="112x92" type="*widget.entryContent">
								<rectangle fillColor="focus" pos="20,25" size="18x21"/>
								<widget size="112x92" type="*widget.textProvider">
									<text pos="4,4" size="104x21">Testing</text>
									<text pos="4,25" size="104x21">Testing</text>
									<text pos="4,46" size="104x21">Testing</text>
								</widget>
								<rectangle fillColor="focus" pos="37,25" size="2x21"/>
							</widget>
						</widget>
					</widget>
				</content>
			</canvas>
		`, c)
		assert.Equal(t, "sti", e.SelectedText())
	}

	return e, window
}

func teardownImageTest(w fyne.Window) {
	w.Close()
	test.NewApp()
}

func waitForBinding() {
	time.Sleep(time.Millisecond * 100) // data resolves on background thread
}
