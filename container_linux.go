package houdini

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"code.cloudfoundry.org/garden"
)

func (container *container) setup() error {
	for _, bm := range container.spec.BindMounts {
		dest := filepath.Join(container.workDir, bm.DstPath)
		err := os.MkdirAll(dest, 0755)
		if err != nil {
			return fmt.Errorf("failed to create target for bind mount: %s", err)
		}

		flags := uintptr(syscall.MS_BIND)
		if bm.Mode == garden.BindMountModeRO {
			flags |= syscall.MS_RDONLY
		}

		err = syscall.Mount(bm.SrcPath, dest, "none", flags, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func (container *container) unsetup() error {
	for i := len(container.spec.BindMounts)-1; i >= 0; i-- {
		bm := container.spec.BindMounts[i]
		dest := filepath.Join(container.workDir, bm.DstPath)
		println("going to unmount", dest)
		err := syscall.Unmount(dest, syscall.MNT_FORCE)
		if err != nil {
			print("WARNING: unmount of ", dest, " got ", err)
		}
	}

	return nil
}
