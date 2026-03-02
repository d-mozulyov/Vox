//go:build darwin
// +build darwin

package indicator

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

// isRetinaDisplay checks if the main display has a backing scale factor > 1.0
int isRetinaDisplay() {
    if ([[NSScreen mainScreen] respondsToSelector:@selector(backingScaleFactor)]) {
        CGFloat scaleFactor = [[NSScreen mainScreen] backingScaleFactor];
        return scaleFactor > 1.0 ? 1 : 0;
    }
    return 0;
}
*/
import "C"

// isRetina returns true if the main display is a Retina display
func isRetina() bool {
	return C.isRetinaDisplay() == 1
}
