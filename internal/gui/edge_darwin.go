//go:build !nogui

package gui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>
#import <objc/runtime.h>

// A window takes a press that resizes it within 2pt of its edge, and the page,
// drawn to that edge, covers the cursor AppKit sets there. This view lies over
// the page and takes a press within 6pt of an edge, resizing the window itself:
// the corners keep AppKit's own, and the top left the traffic lights'.
enum { magpieLeft = 1, magpieRight = 2, magpieTop = 4, magpieBottom = 8 };
static const CGFloat magpieBand = 6;
static const CGFloat magpieLightsW = 90; // the traffic lights', which the top edge starts after
static const CGFloat magpieLightsH = 30; // their height, which the left edge stops below

// The edges a point is within magpieBand of: 0 for the page's own and at a corner.
static int magpieSides(NSPoint p, NSRect b) {
	CGFloat W = NSWidth(b), H = NSHeight(b);
	int s = 0;
	if (p.x <= magpieBand && p.y <= H - magpieLightsH) s |= magpieLeft;
	if (p.x >= W - magpieBand) s |= magpieRight;
	if (p.y <= magpieBand) s |= magpieBottom;
	if (p.y >= H - magpieBand && p.x >= magpieLightsW) s |= magpieTop;
	return (s & (magpieLeft | magpieRight)) && (s & (magpieTop | magpieBottom)) ? 0 : s;
}

@interface MagpieEdge : NSView {
@public
	int sides;    // the edges under the pointer, from the last hitTest
	NSPoint down; // where the press was, in screen points
	NSRect start; // the window's frame then
}
@end
@implementation MagpieEdge
// AppKit leaves a press to the view that takes it (titlebar_darwin.go firstPress).
- (BOOL)acceptsFirstMouse:(NSEvent *)e { return YES; }
- (NSView *)hitTest:(NSPoint)p {
	if (self.window == nil || (self.window.styleMask & NSWindowStyleMaskFullScreen)) return nil;
	sides = magpieSides(p, self.bounds);
	return sides == 0 ? nil : self;
}
- (void)resetCursorRects {
	if (self.window == nil || (self.window.styleMask & NSWindowStyleMaskFullScreen)) return;
	CGFloat W = NSWidth(self.bounds), H = NSHeight(self.bounds);
	NSRect r[4] = {
		NSMakeRect(0, 0, magpieBand, H - magpieLightsH), NSMakeRect(W - magpieBand, 0, magpieBand, H),
		NSMakeRect(0, 0, W, magpieBand), NSMakeRect(magpieLightsW, H - magpieBand, W - magpieLightsW, magpieBand),
	};
	for (int i = 0; i < 4; i++)
		[self addCursorRect:r[i] cursor:i < 2 ? [NSCursor resizeLeftRightCursor] : [NSCursor resizeUpDownCursor]];
}
- (void)mouseDown:(NSEvent *)e {
	if (self.window == nil || sides == 0) return;
	down = [NSEvent mouseLocation];
	start = self.window.frame;
}
- (void)mouseDragged:(NSEvent *)e {
	if (self.window == nil || sides == 0) return;
	NSPoint m = [NSEvent mouseLocation];
	CGFloat dx = m.x - down.x, dy = m.y - down.y;
	NSRect f = start;
	if (sides & magpieLeft) { f.origin.x += dx; f.size.width -= dx; }
	if (sides & magpieRight) f.size.width += dx;
	if (sides & magpieTop) f.size.height += dy;
	if (sides & magpieBottom) { f.origin.y += dy; f.size.height -= dy; }
	// The window's least size holds the side AppKit holds it to.
	NSSize min = self.window.minSize;
	if (f.size.width < min.width) { if (sides & magpieLeft) f.origin.x = NSMaxX(start) - min.width; f.size.width = min.width; }
	if (f.size.height < min.height) { if (sides & magpieBottom) f.origin.y = NSMaxY(start) - min.height; f.size.height = min.height; }
	[self.window setFrame:f display:YES];
}
@end
static char magpieEdgeKey;

// magpieEdge gives a window this view, once (plainTitlebar runs twice).
static void magpieEdge(void *p) {
	NSWindow *w = (NSWindow *)p;
	if (w == nil || w.contentView == nil || objc_getAssociatedObject(w, &magpieEdgeKey) != nil) return;
	MagpieEdge *v = [[MagpieEdge alloc] initWithFrame:w.contentView.bounds];
	v.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
	[w.contentView addSubview:v positioned:NSWindowAbove relativeTo:nil]; // above the page
	objc_setAssociatedObject(w, &magpieEdgeKey, v, OBJC_ASSOCIATION_RETAIN);
	[v release];
}
*/
import "C"

import "unsafe"

// edgeHandles gives the window's edges a press within 6pt of them, with the
// edge's cursor (see magpieEdge). It runs inside the caller's InvokeSync.
func edgeHandles(w unsafe.Pointer) {
	if w != nil {
		C.magpieEdge(w)
	}
}
